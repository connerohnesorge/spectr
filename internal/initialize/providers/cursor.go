package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"
)

// CursorProvider implements the Provider interface for Cursor.
// Cursor uses .cursorrules/commands/spectr/ for slash commands
// (no config file).
type CursorProvider struct{}

// Initializers returns the list of initializers for Cursor. //nolint:lll
func (*CursorProvider) Initializers(
	_ context.Context,
	tm TemplateManager,
) []Initializer { //nolint:lll
	return []Initializer{
		NewDirectoryInitializer(".cursorrules/commands/spectr"),
		NewSlashCommandsInitializer(
			".cursorrules/commands/spectr",
			map[domain.SlashCommand]domain.TemplateRef{
				domain.SlashProposal: tm.SlashCommand(domain.SlashProposal),
				domain.SlashApply:    tm.SlashCommand(domain.SlashApply),
			},
		),
	}
}
