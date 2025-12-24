package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/initialize/templates"
)

func init() {
	_ = RegisterProvider(Registration{
		ID:       "cline",
		Name:     "Cline",
		Priority: PriorityCline,
		Provider: &ClineProvider{},
	})
}

// ClineProvider implements the Provider interface for Cline.
// Cline uses CLINE.md and .clinerules/commands/ for slash commands.
type ClineProvider struct{}

func (*ClineProvider) Initializers(
	_ context.Context,
) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(
			".clinerules/commands/spectr",
		),
		NewConfigFileInitializer(
			"CLINE.md",
			func(tm TemplateManager) any {
				return tm.InstructionPointer()
			},
		),
		NewSlashCommandsInitializer(
			".clinerules/commands/spectr",
			".md",
			[]templates.SlashCommand{
				templates.SlashProposal,
				templates.SlashApply,
			},
		),
	}
}
