package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"
)

// ClineProvider implements the Provider interface for Cline.
// Cline uses CLINE.md and .clinerules/commands/spectr/ for slash commands.
type ClineProvider struct{}

// Initializers returns the list of initializers for Cline.
func (*ClineProvider) Initializers(
	_ context.Context,
	tm TemplateManager,
) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(
			".clinerules/commands/spectr",
		),
		NewConfigFileInitializer(
			"CLINE.md",
			tm.InstructionPointer(),
		),
		NewSlashCommandsInitializer(
			".clinerules/commands/spectr",
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
