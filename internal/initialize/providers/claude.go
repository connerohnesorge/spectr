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
) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".claude/commands/spectr"),
		NewDirectoryInitializer(".claude/skills"),
		NewConfigFileInitializer("CLAUDE.md", tm.InstructionPointer()),
		NewSlashCommandsInitializer(
			".claude/commands/spectr",
			map[domain.SlashCommand]domain.TemplateRef{
				domain.SlashProposal: tm.SlashCommand(domain.SlashProposal),
				domain.SlashApply:    tm.SlashCommand(domain.SlashApply),
			},
		),
		NewAgentSkillsInitializer(
			"spectr-accept-wo-spectr-bin",
			".claude/skills/spectr-accept-wo-spectr-bin",
			tm,
		),
		NewAgentSkillsInitializer(
			"spectr-validate-wo-spectr-bin",
			".claude/skills/spectr-validate-wo-spectr-bin",
			tm,
		),
		NewAgentSkillsInitializer(
			"spectr-validate-w-spectr-bin",
			".claude/skills/spectr-validate-w-spectr-bin",
			tm,
		),
	}
}
