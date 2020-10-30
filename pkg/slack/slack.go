package slack

import (
	"fmt"
	"strconv"
	"strings"
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

func init() {
	// initialize the action lookup table
	startActionLookup = make(map[string]StartActionFunc)
	completeActionLookup = make(map[string]CompleteActionFunc)
	// initialize the slash-command lookup table
	slashCommandLookup = make(map[string]SlashCommandFunc)
	RegisterDefaultSlashCmdHandler(unknownCommandHandler)
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
