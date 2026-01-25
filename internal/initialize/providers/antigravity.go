package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"
)

// AntigravityProvider implements the Provider interface for Antigravity.
// Antigravity uses AGENTS.md and .agent/workflows/ for prefixed
// slash commands (spectr-proposal.md, spectr-apply.md).
type AntigravityProvider struct{}

// Initializers returns the list of initializers for Antigravity.
func (*AntigravityProvider) Initializers(
	_ context.Context,
	tm TemplateManager,
) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(
			".agent/workflows",
		),
		NewConfigFileInitializer(
			"AGENTS.md",
			tm.Agents(),
		),
		NewPrefixedSlashCommandsInitializer(
			".agent/workflows",
			"spectr-",
			map[domain.SlashCommand]domain.TemplateRef{
				domain.SlashProposal: tm.SlashCommand(
					domain.SlashProposal,
				),
				domain.SlashApply: tm.SlashCommand(
					domain.SlashApply,
				),
				domain.SlashNext: tm.SlashCommand(
					domain.SlashNext,
				),
			},
		),
	}
}
