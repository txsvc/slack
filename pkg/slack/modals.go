package slack

type (

	// ViewElement defines a modal view. See https://api.slack.com/surfaces/modals/using#composing_views
	ViewElement struct {
		ID                 string              `json:"id,omitempty"`
		TeamID             string              `json:"team_id,omitempty"`
		Type               string              `json:"type"`
		Title              DefaultViewElement  `json:"title"`
		Submit             *DefaultViewElement `json:"submit,omitempty"`
		Close              *DefaultViewElement `json:"close,omitempty"`
		Blocks             []interface{}       `json:"blocks"`
		PrivateMetadata    string              `json:"private_metadata,omitempty"`
		CallbackID         string              `json:"callback_id,omitempty"`
		State              *StateValues        `json:"state,omitempty"`
		Hash               string              `json:"hash,omitempty"`
		AppID              string              `json:"app_id,omitempty"`
		ExternalID         string              `json:"external_id,omitempty"`
		AppInstalledTeamID string              `json:"app_installed_team_id,omitempty"`
		BotID              string              `json:"bot_id,omitempty"`
	}

	// StateValues holds the selected data
	StateValues struct {
		Values map[string]map[string]ValueObject `json:"values,omitempty"`
	}

	// ValueObject see https://api.slack.com/reference/interaction-payloads/views#view_submission_fields
	ValueObject struct {
		Type    string          `json:"type"`
		Option  *OptionsObject  `json:"selected_option,omitempty"`
		Options []OptionsObject `json:"selected_options,omitempty"`
	}

	// ResponseMetadata map[string]string `json:"response_metadata,omitempty"`

	// DefaultViewElement see https://api.slack.com/surfaces/modals/using#composing_views
	DefaultViewElement struct {
		Type  string `json:"type,omitempty"`
		Text  string `json:"text,omitempty"`
		Emoji bool   `json:"emoji,omitempty"`
	}

	// SectionBlock see https://api.slack.com/reference/block-kit/blocks#section
	// type == section
	SectionBlock struct {
		Type      string       `json:"type"`
		BlockID   string       `json:"block_id,omitempty"`
		Text      TextObject   `json:"text"`
		Fields    []TextObject `json:"fields,omitempty"`
		Accessory interface{}  `json:"accessory,omitempty"`
	}

	// SectionBlocks is an array of SectionBlocks
	SectionBlocks struct {
		Blocks []SectionBlock `json:"blocks"`
	}

	// DividerBlock see https://api.slack.com/reference/block-kit/blocks#divider
	// type == divider
	DividerBlock struct {
		Type    string `json:"type"`
		BlockID string `json:"block_id,omitempty"`
	}

	// InputBlock see https://api.slack.com/reference/block-kit/blocks#input
	// type == input
	InputBlock struct {
		Type     string      `json:"type"`
		BlockID  string      `json:"block_id,omitempty"`
		Label    TextObject  `json:"label"`
		Element  interface{} `json:"element"`
		Hint     *TextObject `json:"hint,omitempty"`
		Optional bool        `json:"optional,omitempty"`
	}

	// Checkboxes see https://api.slack.com/reference/block-kit/block-elements#checkboxes
	// type == checkboxes
	Checkboxes struct {
		Type           string          `json:"type"`
		ActionID       string          `json:"action_id"`
		Options        []OptionsObject `json:"options"`
		InitialOptions []OptionsObject `json:"initial_options,omitempty"`
		Confirm        *ConfirmObject  `json:"confirm,omitempty"`
	}

	// Radiobuttons see https://api.slack.com/reference/block-kit/block-elements#radio
	// type == radio_buttons
	Radiobuttons struct {
		Type          string          `json:"type"`
		ActionID      string          `json:"action_id"`
		Options       []OptionsObject `json:"options"`
		InitialOption *OptionsObject  `json:"initial_option,omitempty"`
		Confirm       *ConfirmObject  `json:"confirm,omitempty"`
	}

	// TextObject see https://api.slack.com/reference/block-kit/composition-objects#text
	TextObject struct {
		Type     string `json:"type"`
		Text     string `json:"text"`
		Emoji    bool   `json:"emoji,omitempty"`
		Verbatim bool   `json:"verbatim,omitempty"`
	}

	// OptionsObject see https://api.slack.com/reference/block-kit/composition-objects#option
	OptionsObject struct {
		Text        TextObject  `json:"text"`
		Value       string      `json:"value"`
		Description *TextObject `json:"description,omitempty"`
		URL         string      `json:"url,omitempty"`
	}

	// ConfirmObject see https://api.slack.com/reference/block-kit/composition-objects#confirm
	ConfirmObject struct {
		Title   *TextObject `json:"title"`
		Text    *TextObject `json:"text"`
		Confirm *TextObject `json:"confirm"`
		Deny    *TextObject `json:"deny"`
	}
)
