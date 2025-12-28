package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"
)

// CrushProvider implements the Provider interface for Crush.
// Crush uses CRUSH.md and .crush/commands/spectr/ for slash commands.
type CrushProvider struct{}

// Initializers returns the list of initializers for Crush.
func (*CrushProvider) Initializers(
	_ context.Context,
	tm TemplateManager,
) []Initializer { //nolint:lll
	return []Initializer{
		NewDirectoryInitializer(".crush/commands/spectr"),
		NewConfigFileInitializer("CRUSH.md", tm.InstructionPointer()),
		NewSlashCommandsInitializer(
			".crush/commands/spectr",
			map[domain.SlashCommand]domain.TemplateRef{
				domain.SlashProposal: tm.SlashCommand(domain.SlashProposal),
				domain.SlashApply:    tm.SlashCommand(domain.SlashApply),
			},
		),
	}
}
