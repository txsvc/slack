package slack

import (
	"encoding/json"
	e "errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"google.golang.org/appengine"

	"github.com/txsvc/commons/pkg/errors"
	"github.com/txsvc/platform/pkg/platform"
)

type (
	// StartActionFunc is a callback for starting a action
	StartActionFunc func(*gin.Context, *ActionRequest) error

	// CompleteActionFunc is a callback for completing an action
	CompleteActionFunc func(*gin.Context, *ViewSubmission) error

	// ViewSubmission see https://api.slack.com/reference/interaction-payloads/views#view_submission
	// type == view_submission
	ViewSubmission struct {
		Type      string             `json:"type,omitempty"`
		Team      *MessageActionTeam `json:"team,omitempty"`
		User      *MessageActionUser `json:"user,omitempty"`
		Token     string             `json:"token,omitempty"`
		TriggerID string             `json:"trigger_id,omitempty"`
		View      *ViewElement       `json:"view,omitempty"`
	}

	// message_action -> ActionRequest
	// view_submission -> ViewSubmission

	// ActionRequestPeek is used to determin the type of request
	ActionRequestPeek struct {
		Type string `json:"type,omitempty"`
	}

	// See https://api.slack.com/reference/interaction-payloads/actions

	// ActionRequest is the payload received from Slack when the user triggers a custom message action
	// type == message_action
	ActionRequest struct {
		Type             string                `json:"type,omitempty"`
		Token            string                `json:"token,omitempty"`
		ActionTimestamp  string                `json:"action_ts,omitempty"`
		Team             *MessageActionTeam    `json:"team,omitempty"`
		User             *MessageActionUser    `json:"user,omitempty"`
		Channel          *MessageActionChannel `json:"channel,omitempty"`
		CallbackID       string                `json:"callback_id,omitempty"`
		TriggerID        string                `json:"trigger_id,omitempty"`
		MessageTimestamp string                `json:"message_ts,omitempty"`
		Message          *ActionRequestMessage `json:"message,omitempty"`
		ResponseURL      string                `json:"response_url,omitempty"`
		Submission       map[string]string     `json:"submission,omitempty"`
	}

	// ActionRequestMessage is the message's main content
	ActionRequestMessage struct {
		Type         string                     `json:"type,omitempty"`
		User         string                     `json:"user,omitempty"`
		Text         string                     `json:"text,omitempty"`
		Attachements []ActionRequestAttachement `json:"attachments,omitempty"`
		Timestamp    string                     `json:"ts,omitempty"`
	}

	// ActionRequestAttachement describes message attachements such as links or files
	ActionRequestAttachement struct {
		ServiceName string `json:"service_name,omitempty"`
		Title       string `json:"title,omitempty"`
		TitleLink   string `json:"title_link,omitempty"`
		Text        string `json:"text,omitempty"`
		Fallback    string `json:"fallback,omitempty"`
		ImageURL    string `json:"image_url,omitempty"`
		FromURL     string `json:"from_url,omitempty"`
		ImageWidth  int    `json:"image_width,omitempty"`
		ImageHeight int    `json:"image_height,omitempty"`
		ImageBytes  int    `json:"image_bytes,omitempty"`
		ServiceIcon string `json:"service_icon,omitempty"`
		ID          int    `json:"id,omitempty"`
		OriginalURL string `json:"original_url,omitempty"`
	}

	// MessageActionTeam identifies the Slack workspace the message originates from
	MessageActionTeam struct {
		ID     string `json:"id,omitempty"`
		Domain string `json:"domain,omitempty"`
	}

	// MessageActionUser identifies the user who triggered the custom action
	MessageActionUser struct {
		ID   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	}

	// MessageActionChannel identifies the channel the custom action was triggered from
	MessageActionChannel struct {
		ID   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	}
)

// actions callback lookups
var startActionLookup map[string]StartActionFunc
var completeActionLookup map[string]CompleteActionFunc

// ActionRequestEndpoint receives callbacks from Slack
func ActionRequestEndpoint(c *gin.Context) {
	var peek ActionRequestPeek

	err := json.Unmarshal([]byte(c.Request.FormValue("payload")), &peek)
	if err != nil {
		platform.ReportError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "msg": err.Error()})
		return
	}

	if peek.Type == "message_action" {
		var action ActionRequest
		err := json.Unmarshal([]byte(c.Request.FormValue("payload")), &action)
		if err != nil {
			platform.ReportError(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "msg": err.Error()})
			return
		}

		err = startAction(c, &action)
		if err != nil {
			platform.ReportError(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "msg": err.Error()})
			return
		}

	} else if peek.Type == "view_submission" {
		var submission ViewSubmission
		err := json.Unmarshal([]byte(c.Request.FormValue("payload")), &submission)
		if err != nil {
			platform.ReportError(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "msg": err.Error()})
			return
		}

		err = completeAction(c, &submission)
		if err != nil {
			platform.ReportError(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "msg": err.Error()})
			return
		}
	} else {
		platform.ReportError(fmt.Errorf("Unknown action request: '%s'", peek.Type))
	}
}

// RegisterStartAction adds a start action handler
func RegisterStartAction(action string, h StartActionFunc) {
	startActionLookup[strings.ToLower(action)] = h
}

// RegisterCompleteAction adds a completion action handler
func RegisterCompleteAction(action string, h CompleteActionFunc) {
	completeActionLookup[strings.ToLower(action)] = h
}

// StoreActionCorrelation is a helper to mange correlation keys
func StoreActionCorrelation(ctx context.Context, action, viewID, teamID string) error {
	err := platform.SetKV(ctx, correlationKey(viewID, teamID), strings.ToLower(action), 1800)
	if err != nil {
		platform.ReportError(err)
	}
	return err
}

// startAction initiates a dialog with the user
func startAction(c *gin.Context, a *ActionRequest) error {
	action := a.CallbackID
	handler := startActionLookup[strings.ToLower(action)]
	if handler == nil {
		return errors.NewOperationError(action, e.New(fmt.Sprintf("No handler for action request '%s'", action)))
	}

	return handler(c, a)
}

// completeAction starts the processing of the action's result
func completeAction(c *gin.Context, s *ViewSubmission) error {
	ctx := appengine.NewContext(c.Request)

	action := lookupActionCorrelation(ctx, s.View.ID, s.Team.ID)
	if action == "" {
		return nil
	}

	handler := completeActionLookup[action]
	if handler == nil {
		return errors.NewOperationError(action, e.New(fmt.Sprintf("No handler for action response '%s'", action)))
	}

	return handler(c, s)
}

func lookupActionCorrelation(ctx context.Context, viewID, teamID string) string {
	v, err := platform.GetKV(ctx, correlationKey(viewID, teamID))
	if err != nil {
		return ""
	}
	return v
}

func correlationKey(viewID, teamID string) string {
	return viewID + "." + teamID
}
