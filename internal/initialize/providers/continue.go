package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"

	//nolint:revive // long import path is unavoidable
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
)

// ContinueProvider configures Continue.
// Continue uses .continue/commands/ for slash commands (no config file).
type ContinueProvider struct{}

func (*ContinueProvider) Initializers(
	_ context.Context,
) []Initializer {
	return []Initializer{
		initializers.NewDirectoryInitializer(
			".continue/commands/spectr",
		),
		initializers.NewSlashCommandsInitializer(
			".continue/commands/spectr",
			".md",
			[]domain.SlashCommand{
				domain.SlashProposal,
				domain.SlashApply,
			},
		),
	}
}
