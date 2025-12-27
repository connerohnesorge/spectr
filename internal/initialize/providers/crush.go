package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/connerohnesorge/spectr/internal/templates"
)

// CrushProvider implements the Provider interface for Crush.
// Crush uses CRUSH.md for instructions and .crush/commands/spectr/ for
// slash commands.
type CrushProvider struct{}

// Initializers returns the initializers for Crush provider.
func (*CrushProvider) Initializers(
	_ context.Context,
	tm *templates.TemplateManager,
) []Initializer {
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
