package slack

import (
	e "errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/txsvc/commons/pkg/errors"
	"github.com/txsvc/platform/pkg/platform"
)

// see https://api.slack.com/interactivity/slash-commands

type (
	// SlashCommandFunc handles a slash command
	SlashCommandFunc func(c *gin.Context, cmd *SlashCommand) (*SectionBlocks, error)

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

	cmdErrorWrapper struct {
		err error
		cmd *SlashCommand
		msg string
	}
)

var slashCommandLookup map[string]SlashCommandFunc
var defaultSlashCommandHandler SlashCommandFunc

// SlashCmdEndpoint receives callbacks from Slack command /lnkk
func SlashCmdEndpoint(c *gin.Context) {
	status := http.StatusOK

	// extract the cmd and react on it
	cmd := GetSlashCommand(c)

	// dispatch to a handler
	handler := slashCommandLookup[strings.ToLower(cmd.Command)]
	if handler == nil {
		e := errors.NewOperationError(cmd.Command, e.New(fmt.Sprintf("No handler for command '%s'", cmd.Command)))
		platform.ReportError(e)

		handler = defaultSlashCommandHandler
		status = http.StatusBadRequest
	}

	resp, err := handler(c, cmd)

	if err != nil {
		status = http.StatusOK
		if resp == nil {
			resp = genericErrorSectionBlock(c, cmd)
		}
	}

	c.JSON(status, resp)
}

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

// RegisterSlashCmdHandler adds a slash-cmd handler
func RegisterSlashCmdHandler(cmd string, h SlashCommandFunc) {
	slashCommandLookup[strings.ToLower(cmd)] = h
}

// RegisterDefaultSlashCmdHandler adds a default slash-cmd handler
func RegisterDefaultSlashCmdHandler(h SlashCommandFunc) {
	defaultSlashCommandHandler = h
}

// NewSlackCmdEror wraps an error with additional metadata
func NewSlackCmdEror(msg string, cmd *SlashCommand, e error) error {
	return &cmdErrorWrapper{cmd: cmd, msg: msg, err: e}
}

func (ee *cmdErrorWrapper) Error() string {
	return fmt.Sprintf("cmd: %s, team_id: %s, channel_id: %s, user_id: %s, cmdline: %s", ee.cmd.Command, ee.cmd.TeamID, ee.cmd.ChannelID, ee.cmd.UserID, ee.cmd.Txt)
}

func (ee *cmdErrorWrapper) Unwrap() error {
	return ee.err
}

func unknownCommandHandler(c *gin.Context, cmd *SlashCommand) (*SectionBlocks, error) {
	return genericErrorSectionBlock(c, cmd), nil
}

func genericErrorSectionBlock(c *gin.Context, cmd *SlashCommand) *SectionBlocks {
	return &SectionBlocks{
		Blocks: []SectionBlock{
			{
				Type: "section",
				Text: TextObject{
					Type: "mrkdwn",
					Text: fmt.Sprintf("Sorry, but I can't do this: %s %s", cmd.Command, cmd.Txt),
				},
			},
		},
	}
}
