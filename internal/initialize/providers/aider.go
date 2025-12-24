package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/initialize/templates"
)

func init() {
	_ = RegisterProvider(Registration{
		ID:       "aider",
		Name:     "Aider",
		Priority: PriorityAider,
		Provider: &AiderProvider{},
	})
}

// AiderProvider implements the Provider interface for Aider.
// Aider uses .aider/commands/ for slash commands (no config file).
type AiderProvider struct{}

func (*AiderProvider) Initializers(
	_ context.Context,
) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(
			".aider/commands/spectr",
		),
		NewSlashCommandsInitializer(
			".aider/commands/spectr",
			".md",
			[]templates.SlashCommand{
				templates.SlashProposal,
				templates.SlashApply,
			},
		),
	}
}
