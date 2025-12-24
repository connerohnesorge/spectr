package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/initialize/templates"
)

func init() {
	_ = RegisterProvider(Registration{
		ID:       "claude-code",
		Name:     "Claude Code",
		Priority: PriorityClaudeCode,
		Provider: &ClaudeProvider{},
	})
}

// ClaudeProvider implements the Provider interface for Claude Code.
// Claude Code uses CLAUDE.md and .claude/commands/ for slash commands.
type ClaudeProvider struct{}

func (*ClaudeProvider) Initializers(
	_ context.Context,
) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(
			".claude/commands/spectr",
		),
		NewConfigFileInitializer(
			"CLAUDE.md",
			func(tm TemplateManager) any {
				return tm.InstructionPointer()
			},
		),
		NewSlashCommandsInitializer(
			".claude/commands/spectr",
			".md",
			[]templates.SlashCommand{
				templates.SlashProposal,
				templates.SlashApply,
			},
		),
	}
}
