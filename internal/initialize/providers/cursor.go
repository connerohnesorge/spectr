package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/initialize/templates"
)

func init() {
	_ = RegisterProvider(Registration{
		ID:       "cursor",
		Name:     "Cursor",
		Priority: PriorityCursor,
		Provider: &CursorProvider{},
	})
}

// CursorProvider implements the Provider interface for Cursor.
// Cursor uses .cursorrules/commands/ for slash commands (no config file).
type CursorProvider struct{}

func (*CursorProvider) Initializers(
	_ context.Context,
) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(
			".cursorrules/commands/spectr",
		),
		NewSlashCommandsInitializer(
			".cursorrules/commands/spectr",
			".md",
			[]templates.SlashCommand{
				templates.SlashProposal,
				templates.SlashApply,
			},
		),
	}
}
