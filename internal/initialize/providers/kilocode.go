package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/initialize/templates"
)

func init() {
	_ = RegisterProvider(Registration{
		ID:       "kilocode",
		Name:     "Kilocode",
		Priority: PriorityKilocode,
		Provider: &KilocodeProvider{},
	})
}

// KilocodeProvider implements the Provider interface for Kilocode.
// Kilocode uses .kilocode/commands/ for slash commands (no config file).
type KilocodeProvider struct{}

func (*KilocodeProvider) Initializers(
	_ context.Context,
) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(
			".kilocode/commands/spectr",
		),
		NewSlashCommandsInitializer(
			".kilocode/commands/spectr",
			".md",
			[]templates.SlashCommand{
				templates.SlashProposal,
				templates.SlashApply,
			},
		),
	}
}
