package slack

import (
	e "errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/txsvc/commons/pkg/errors"
	"github.com/txsvc/platform/pkg/platform"
)

// see https://api.slack.com/interactivity/slash-commands

type (
	// SlashCommand encapsulates the payload sent from Slack as a result of invoking a /slash command
	SlashCommand struct {
		TeamID         string
		TeamDomain     string
		EnterpriseID   string
		EnterpriseName string
		ChannelID      string
		ChannelName    string
		UserID         string
		UserName       string
		Command        string
		Txt            string
		ResponseURL    string
		TriggerID      string
		Token          string // DEPRECATED
	}
)

// GetSlashCommand extracts the payload from a POST received by a slash command
func GetSlashCommand(c *gin.Context) *SlashCommand {
	return &SlashCommand{
		TeamID:         c.PostForm("team_id"),
		TeamDomain:     c.PostForm("team_domain"),
		EnterpriseID:   c.PostForm("enterprise_id"),
		EnterpriseName: c.PostForm("enterprise_name"),
		ChannelID:      c.PostForm("channel_id"),
		ChannelName:    c.PostForm("channel_name"),
		UserID:         c.PostForm("user_id"),
		UserName:       c.PostForm("user_name"),
		Command:        c.PostForm("command"),
		Txt:            c.PostForm("text"),
		ResponseURL:    c.PostForm("response_url"),
		TriggerID:      c.PostForm("trigger_id"),
		Token:          c.PostForm("token"), // DEPRECATED
	}
}

// SlashCmdEndpoint receives callbacks from Slack command /lnkk
func SlashCmdEndpoint(c *gin.Context) {
	status := http.StatusOK

	// extract the cmd and react on it
	cmd := GetSlashCommand(c)

	// dispatch to a handler
	handler := slashCommandLookup[cmd.Command]
	if handler == nil {
		e := errors.NewOperationError(cmd.Command, e.New(fmt.Sprintf("No handler for slash command '%s'", cmd.Command)))
		platform.Report(e)
		handler = defaultSlashCommandHandler
	}

	resp, err := handler(c, cmd)

	if err != nil {
		status = http.StatusBadRequest
		if resp == nil {
			response := SectionBlocks{
				Blocks: []SectionBlock{
					{
						Type: "section",
						Text: TextObject{
							Type: "mrkdwn",
							Text: fmt.Sprintf("Something went wrong: '%s': %s", cmd.Command, err.Error()),
						},
					},
				},
			}
			resp = &response
		}
	}

	c.JSON(status, resp)
}

func errorSlashCmdHandler(c *gin.Context, cmd *SlashCommand) (*SectionBlocks, error) {
	return &SectionBlocks{
		Blocks: []SectionBlock{
			{
				Type: "section",
				Text: TextObject{
					Type: "mrkdwn",
					Text: fmt.Sprintf("Sorry, but I can't do this: '%s'", cmd.Command),
				},
			},
		},
	}, nil
}
