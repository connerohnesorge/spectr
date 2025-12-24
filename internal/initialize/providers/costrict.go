package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/initialize/templates"
)

func init() {
	_ = RegisterProvider(Registration{
		ID:       "costrict",
		Name:     "CoStrict",
		Priority: PriorityCostrict,
		Provider: &CostrictProvider{},
	})
}

// CostrictProvider implements the Provider interface for CoStrict.
// CoStrict uses COSTRICT.md and .costrict/commands/ for slash commands.
type CostrictProvider struct{}

func (*CostrictProvider) Initializers(
	_ context.Context,
) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(
			".costrict/commands/spectr",
		),
		NewConfigFileInitializer(
			"COSTRICT.md",
			func(tm TemplateManager) any {
				return tm.InstructionPointer()
			},
		),
		NewSlashCommandsInitializer(
			".costrict/commands/spectr",
			".md",
			[]templates.SlashCommand{
				templates.SlashProposal,
				templates.SlashApply,
			},
		),
	}
}
