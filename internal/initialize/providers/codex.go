package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/connerohnesorge/spectr/internal/templates"
)

// CodexProvider implements the Provider interface for Codex CLI.
// Codex uses AGENTS.md and global ~/.codex/prompts/ for commands.
type CodexProvider struct{}

// Initializers returns the initializers for Codex CLI provider.
func (*CodexProvider) Initializers(
	_ context.Context,
	tm *templates.TemplateManager,
) []Initializer {
	return []Initializer{
		NewHomeDirectoryInitializer(".codex/prompts"),
		NewConfigFileInitializer("AGENTS.md", tm.Agents()),
		NewHomePrefixedSlashCommandsInitializer(
			".codex/prompts",
			"spectr-",
			map[domain.SlashCommand]domain.TemplateRef{
				domain.SlashProposal: tm.SlashCommand(domain.SlashProposal),
				domain.SlashApply:    tm.SlashCommand(domain.SlashApply),
			},
		),
	}
}
