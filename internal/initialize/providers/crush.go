package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers" //nolint:revive
)

// CrushProvider configures Crush.
// Crush uses CRUSH.md and .crush/commands/ for slash commands.
type CrushProvider struct{}

func (*CrushProvider) Initializers(
	_ context.Context,
) []Initializer {
	return []Initializer{
		initializers.NewDirectoryInitializer(
			".crush/commands/spectr",
		),
		initializers.NewConfigFileInitializer(
			"CRUSH.md",
			"instruction-pointer.md.tmpl",
		),
		initializers.NewSlashCommandsInitializer(
			".crush/commands/spectr",
			".md",
			[]domain.SlashCommand{
				domain.SlashProposal,
				domain.SlashApply,
			},
		),
	}
}
