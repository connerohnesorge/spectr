// Package providers implements AI tool provider registration and
// initialization.
package providers

import (
	"context"
)

func init() {
	if err := Register(Registration{
		ID:       "claude-code",
		Name:     "Claude Code",
		Priority: PriorityClaudeCode,
		Provider: &ClaudeProvider{},
	}); err != nil {
		panic(err)
	}
}

// ClaudeProvider implements the Provider interface for Claude Code.
type ClaudeProvider struct{}

func (*ClaudeProvider) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".claude/commands/spectr"),
		NewConfigFileInitializer("CLAUDE.md", "instruction_pointer"),
		NewSlashCommandsInitializer(
			".claude/commands/spectr",
			".md",
			FormatMarkdown,
		),
	}
}
