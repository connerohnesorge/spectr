package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"

	//nolint:revive // long import path is unavoidable
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
)

// WindsurfProvider configures Windsurf.
// Windsurf uses .windsurf/commands/ for slash commands (no config file).
type WindsurfProvider struct{}

func (*WindsurfProvider) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		initializers.NewDirectoryInitializer(".windsurf/commands/spectr"),
		initializers.NewSlashCommandsInitializer(
			".windsurf/commands/spectr",
			".md",
			[]domain.SlashCommand{
				domain.SlashProposal,
				domain.SlashApply,
			},
		),
	}
}
