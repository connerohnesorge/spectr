package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"
)

// CostrictProvider implements the Provider interface for CoStrict.
// CoStrict uses COSTRICT.md and .costrict/commands/spectr/ for slash commands.
type CostrictProvider struct{}

// Initializers returns the list of initializers for CoStrict.
func (*CostrictProvider) Initializers(
	_ context.Context,
	tm TemplateManager,
) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(
			".costrict/commands/spectr",
		),
		NewConfigFileInitializer(
			"COSTRICT.md",
			tm.InstructionPointer(),
		),
		NewSlashCommandsInitializer(
			".costrict/commands/spectr",
			map[domain.SlashCommand]domain.TemplateRef{
				domain.SlashProposal: tm.SlashCommand(
					domain.SlashProposal,
				),
				domain.SlashApply: tm.SlashCommand(
					domain.SlashApply,
				),
			},
		),
	}
}
