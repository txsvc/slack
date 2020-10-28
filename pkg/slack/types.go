package slack

type (

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
)
