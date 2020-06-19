package slack

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
