package slack

import (
	"encoding/json"
	e "errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"google.golang.org/appengine"

	"github.com/txsvc/commons/pkg/errors"
	"github.com/txsvc/platform/pkg/platform"
)

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

// startAction initiates a dialog with the user
func startAction(c *gin.Context, a *ActionRequest) error {
	action := a.CallbackID
	handler := startActionLookup[action]
	if handler == nil {
		return errors.NewOperationError(action, e.New(fmt.Sprintf("No handler for action request '%s'", action)))
	}

	return handler(c, a)
}

// completeAction starts the processing of the action's result
func completeAction(c *gin.Context, s *ViewSubmission) error {
	ctx := appengine.NewContext(c.Request)

	action := LookupActionCorrelation(ctx, s.View.ID, s.Team.ID)
	if action == "" {
		return nil
	}

	handler := completeActionLookup[action]
	if handler == nil {
		return errors.NewOperationError(action, e.New(fmt.Sprintf("No handler for action response '%s'", action)))
	}

	return handler(c, s)
}

// StoreActionCorrelation is a helper to mange correlation keys
func StoreActionCorrelation(ctx context.Context, action, viewID, teamID string) error {
	err := platform.SetKV(ctx, correlationKey(viewID, teamID), action, 1800)
	if err != nil {
		platform.ReportError(err)
	}
	return err
}

// LookupActionCorrelation is a helper to mange correlation keys
func LookupActionCorrelation(ctx context.Context, viewID, teamID string) string {
	v, err := platform.GetKV(ctx, correlationKey(viewID, teamID))
	if err != nil {
		platform.ReportError(err)
		return ""
	}
	return v
}

func correlationKey(viewID, teamID string) string {
	return viewID + "." + teamID
}
