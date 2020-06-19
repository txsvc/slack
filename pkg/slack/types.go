package slack

type (
	// OAuthResponse is used to give a simple reponse to the user as feedback to a custom action reuqest
	OAuthResponse struct {
		OK              bool            `json:"ok,omitempty"`
		AccessToken     string          `json:"access_token,omitempty"`
		TokenType       string          `json:"token_type,omitempty"`
		Scope           string          `json:"scope,omitempty"`
		AppID           string          `json:"app_id,omitempty"`
		BotUserID       string          `json:"bot_user_id,omitempty"`
		IncomingWebhook *WebhookElement `json:"incoming_webhook,omitempty"`
	}

	// StandardResponse is the generic response received after a Web API request.
	// See https://api.slack.com/web#responses
	StandardResponse struct {
		OK               bool         `json:"ok"`
		Stuff            string       `json:"stuff,omitempty"`
		Warning          string       `json:"warning,omitempty"`
		Error            string       `json:"error,omitempty"`
		ResponseMetadata MessageArray `json:"response_metadata,omitempty"`
	}

	// MessageArray is a container for an array of error strings
	MessageArray struct {
		Messages []string `json:"messages,omitempty"`
	}

	// WebhookElement not sure?
	WebhookElement struct {
		URL              string `json:"url,omitempty"`
		Channel          string `json:"channel,omitempty"`
		ConfigurationURL string `json:"configuration_url,omitempty"`
	}

	// TeamInfo see https://api.slack.com/methods/team.info
	TeamInfo struct {
		OK    bool        `json:"ok"`
		Error string      `json:"error,omitempty"`
		Team  TeamElement `json:"team"`
	}

	// TeamElement see https://api.slack.com/methods/team.info
	TeamElement struct {
		ID             string       `json:"id,omitempty"`
		Name           string       `json:"name,omitempty"`
		Domain         string       `json:"domain,omitempty"`
		EmailDomain    string       `json:"email_domain,omitempty"`
		Icon           *IconElement `json:"icon,omitempty"`
		EnterpriseID   string       `json:"enterprise_id,omitempty"`
		EnterpriseName string       `json:"enterprise_name,omitempty"`
	}

	// IconElement see https://api.slack.com/methods/team.info
	IconElement struct {
		Image44      string `json:"image_34,omitempty"`
		Image68      string `json:"image_68,omitempty"`
		Image88      string `json:"image_88,omitempty"`
		Image102     string `json:"image_102,omitempty"`
		Image132     string `json:"image_132,omitempty"`
		ImageDefault bool   `json:"image_default,omitempty"`
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

	// ModalRequest creates a modal dialog as a response to an action request
	// see https://api.slack.com/methods/views.open#response
	ModalRequest struct {
		TriggerID string      `json:"trigger_id"`
		View      ViewElement `json:"view"`
	}

	// ModalResponse is the reply to a ModalRequest
	ModalResponse struct {
		OK               bool         `json:"ok"`
		View             *ViewElement `json:"view,omitempty"`
		Error            string       `json:"error,omitempty"`
		ResponseMetadata MessageArray `json:"response_metadata,omitempty"`
	}

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
)
