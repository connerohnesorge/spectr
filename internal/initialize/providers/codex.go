package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"
)

// CodexProvider implements the Provider interface for Codex CLI.
// Codex uses AGENTS.md and global ~/.codex/prompts/ for prefixed
// slash commands (spectr-proposal.md, spectr-apply.md).
type CodexProvider struct{}

// Initializers returns the list of initializers for Codex CLI. //nolint:lll // Function signature defined by Provider interface
func (*CodexProvider) Initializers(
	_ context.Context,
	tm TemplateManager,
) []Initializer { //nolint:lll // Constructor calls with template refs exceed line limit
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
