package cmd

import (
	"fmt"
	"os"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/connerohnesorge/spectr/internal/hooks"
)

// HooksCmd represents the hooks command for processing Claude Code hook events.
type HooksCmd struct {
	HookType string `arg:""           help:"Hook event type (PreToolUse, Stop, etc.)"`  //nolint:lll,revive // Kong struct tag
	Command  string `name:"command" short:"c" help:"Slash command context" required:""` //nolint:lll,revive // Kong struct tag
}

// Run executes the hooks command by parsing the hook type and delegating
// to the hooks package for processing.
func (c *HooksCmd) Run() error {
	ht, ok := domain.ParseHookType(c.HookType)
	if !ok {
		return fmt.Errorf(
			"unknown hook type: %s",
			c.HookType,
		)
	}

	return hooks.Handle(ht, c.Command, os.Stdin, os.Stdout)
}
