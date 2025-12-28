package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"
)

// ClaudeProvider implements the Provider interface for Claude Code.
// Claude Code uses CLAUDE.md and .claude/commands/spectr/ for slash commands.
type ClaudeProvider struct{}

// Initializers returns the list of initializers for Claude Code.
func (*ClaudeProvider) Initializers(
	_ context.Context,
	tm TemplateManager,
) []Initializer { //nolint:lll
	return []Initializer{
		NewDirectoryInitializer(".claude/commands/spectr"),
		NewConfigFileInitializer("CLAUDE.md", tm.InstructionPointer()),
		NewSlashCommandsInitializer(
			".claude/commands/spectr",
			map[domain.SlashCommand]domain.TemplateRef{
				domain.SlashProposal: tm.SlashCommand(domain.SlashProposal),
				domain.SlashApply:    tm.SlashCommand(domain.SlashApply),
			},
		),
	}
}
