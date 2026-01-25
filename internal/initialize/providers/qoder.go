package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"
)

// QoderProvider implements the Provider interface for Qoder.
// Qoder uses QODER.md and .qoder/commands/spectr/ for slash commands.
type QoderProvider struct{}

// Initializers returns the list of initializers for Qoder.
func (*QoderProvider) Initializers(
	_ context.Context,
	tm TemplateManager,
) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(
			".qoder/commands/spectr",
		),
		NewConfigFileInitializer(
			"QODER.md",
			tm.InstructionPointer(),
		),
		NewSlashCommandsInitializer(
			".qoder/commands/spectr",
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
