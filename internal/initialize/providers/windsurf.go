package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/initialize/templates"
)

func init() {
	_ = RegisterProvider(Registration{
		ID:       "windsurf",
		Name:     "Windsurf",
		Priority: PriorityWindsurf,
		Provider: &WindsurfProvider{},
	})
}

// WindsurfProvider implements the Provider interface for Windsurf.
// Windsurf uses .windsurf/commands/ for slash commands (no config file).
type WindsurfProvider struct{}

func (*WindsurfProvider) Initializers(
	_ context.Context,
) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(
			".windsurf/commands/spectr",
		),
		NewSlashCommandsInitializer(
			".windsurf/commands/spectr",
			".md",
			[]templates.SlashCommand{
				templates.SlashProposal,
				templates.SlashApply,
			},
		),
	}
}
