package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/initialize/templates"
)

func init() {
	_ = RegisterProvider(Registration{
		ID:       "qoder",
		Name:     "Qoder",
		Priority: PriorityQoder,
		Provider: &QoderProvider{},
	})
}

// QoderProvider implements the Provider interface for Qoder.
// Qoder uses QODER.md and .qoder/commands/ for slash commands.
type QoderProvider struct{}

func (*QoderProvider) Initializers(
	_ context.Context,
) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(
			".qoder/commands/spectr",
		),
		NewConfigFileInitializer(
			"QODER.md",
			func(tm TemplateManager) any {
				return tm.InstructionPointer()
			},
		),
		NewSlashCommandsInitializer(
			".qoder/commands/spectr",
			".md",
			[]templates.SlashCommand{
				templates.SlashProposal,
				templates.SlashApply,
			},
		),
	}
}
