package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"
)

// ContinueProvider implements the Provider interface for Continue.
// Continue uses .continue/commands/spectr/ for slash commands (no config file).
type ContinueProvider struct{}

// Initializers returns the list of initializers for Continue.
func (*ContinueProvider) Initializers(
	_ context.Context,
	tm TemplateManager,
) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(
			".continue/commands/spectr",
		),
		NewSlashCommandsInitializer(
			".continue/commands/spectr",
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
