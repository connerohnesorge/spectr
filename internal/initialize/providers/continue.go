package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/initialize/templates"
)

func init() {
	_ = RegisterProvider(Registration{
		ID:       "continue",
		Name:     "Continue",
		Priority: PriorityContinue,
		Provider: &ContinueProvider{},
	})
}

// ContinueProvider implements the Provider interface for Continue.
// Continue uses .continue/commands/ for slash commands (no config file).
type ContinueProvider struct{}

func (*ContinueProvider) Initializers(
	_ context.Context,
) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(
			".continue/commands/spectr",
		),
		NewSlashCommandsInitializer(
			".continue/commands/spectr",
			".md",
			[]templates.SlashCommand{
				templates.SlashProposal,
				templates.SlashApply,
			},
		),
	}
}
