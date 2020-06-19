package slack

import (
	"github.com/gin-gonic/gin"
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
