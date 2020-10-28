package slack

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// SlackEndpoint is the URI of the Slack API
	SlackEndpoint string = "https://slack.com/api/"
	// SlackClientID is the Client ID of App
	SlackClientID string = "SLACK_CLIENT_ID"
	// SlackClientSecret is the Client Secret of the App
	SlackClientSecret string = "SLACK_CLIENT_SECRET"
	// SlackOAuthToken is the default OAuth token
	SlackOAuthToken string = "SLACK_OAUTH_TOKEN"
	// SlackVerificationToken is a secret token used to verify requests from Slack
	SlackVerificationToken string = "SLACK_VERIFICATION_TOKEN"
	// SlackResponseTypeChannel is used to send messages to channels that are visible to evryone
	SlackResponseTypeChannel string = "in_channel"
	// SlackResponseTypeEphemeral is used to send a message to a channel that is only visible to the current user
	SlackResponseTypeEphemeral string = "ephemeral"
)

// StartActionFunc is a callback for starting a action
type StartActionFunc func(*gin.Context, *ActionRequest) error

// CompleteActionFunc is a callback for completing an action
type CompleteActionFunc func(*gin.Context, *ViewSubmission) error

// SlashCommandFunc handles a slash command
type SlashCommandFunc func(c *gin.Context, cmd *SlashCommand) (*SectionBlocks, error)

// actions callback lookups
var startActionLookup map[string]StartActionFunc
var completeActionLookup map[string]CompleteActionFunc
var slashCommandLookup map[string]SlashCommandFunc
var defaultSlashCommandHandler SlashCommandFunc

func init() {
	// initialize the action lookup table
	startActionLookup = make(map[string]StartActionFunc)
	completeActionLookup = make(map[string]CompleteActionFunc)
	// initialize the slash-command lookup table
	slashCommandLookup = make(map[string]SlashCommandFunc)
	RegisterDefaultSlashCmdHandler(errorSlashCmdHandler)
}

// RegisterStartAction adds a start action handler
func RegisterStartAction(action string, h StartActionFunc) {
	startActionLookup[action] = h
}

// RegisterCompleteAction adds a completion action handler
func RegisterCompleteAction(action string, h CompleteActionFunc) {
	completeActionLookup[action] = h
}

// RegisterSlashCmdHandler adds a slash-cmd handler
func RegisterSlashCmdHandler(cmd string, h SlashCommandFunc) {
	slashCommandLookup[cmd] = h
}

// RegisterDefaultSlashCmdHandler adds a default slash-cmd handler
func RegisterDefaultSlashCmdHandler(h SlashCommandFunc) {
	defaultSlashCommandHandler = h
}

// Timestamp returns the seconds part of a Slack timestamp
// Example: "1533028651.000211" -> 1533028651
func Timestamp(ts string) int64 {
	s := strings.Split(ts, ".")
	i, _ := strconv.ParseInt(s[0], 10, 64)
	return i
}

// TimestampNano returns a Slack timestamp as nanoseconds
func TimestampNano(ts string) int64 {
	s := strings.Split(ts, ".")
	i, _ := strconv.ParseInt(s[0], 10, 64)
	j, _ := strconv.ParseInt(s[1], 10, 64)

	return (i * 1000000) + j
}

// TimestampNanoString converts a Slack TS in nanoseconds into Slack's string representation
func TimestampNanoString(ts int64) string {
	_p1 := ts / 1000000
	_p2 := ts % 1000000
	return fmt.Sprintf("%d.%06d", _p1, _p2)
}
