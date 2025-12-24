package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/initialize/templates"
)

func init() {
	_ = RegisterProvider(Registration{
		ID:       "crush",
		Name:     "Crush",
		Priority: PriorityCrush,
		Provider: &CrushProvider{},
	})
}

// CrushProvider implements the Provider interface for Crush.
// Crush uses CRUSH.md for instructions and .crush/commands/ for slash commands.
type CrushProvider struct{}

func (*CrushProvider) Initializers(
	_ context.Context,
) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(
			".crush/commands/spectr",
		),
		NewConfigFileInitializer(
			"CRUSH.md",
			func(tm TemplateManager) any {
				return tm.InstructionPointer()
			},
		),
		NewSlashCommandsInitializer(
			".crush/commands/spectr",
			".md",
			[]templates.SlashCommand{
				templates.SlashProposal,
				templates.SlashApply,
			},
		),
	}
}
