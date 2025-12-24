package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"

	//nolint:revive // long import path is unavoidable
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
)

// AiderProvider configures Aider.
// Aider uses .aider/commands/ for slash commands (no config file).
type AiderProvider struct{}

func (*AiderProvider) Initializers(
	_ context.Context,
) []Initializer {
	return []Initializer{
		initializers.NewDirectoryInitializer(
			".aider/commands/spectr",
		),
		initializers.NewSlashCommandsInitializer(
			".aider/commands/spectr",
			".md",
			[]domain.SlashCommand{
				domain.SlashProposal,
				domain.SlashApply,
			},
		),
	}
}
