package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"

	//nolint:revive // long import path is unavoidable
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
)

// CursorProvider configures Cursor.
// Cursor uses .cursorrules/commands/ for slash commands (no config file).
type CursorProvider struct{}

func (*CursorProvider) Initializers(
	_ context.Context,
) []Initializer {
	return []Initializer{
		initializers.NewDirectoryInitializer(
			".cursorrules/commands/spectr",
		),
		initializers.NewSlashCommandsInitializer(
			".cursorrules/commands/spectr",
			".md",
			[]domain.SlashCommand{
				domain.SlashProposal,
				domain.SlashApply,
			},
		),
	}
}
