package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"

	//nolint:revive // long import path is unavoidable
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
)

// CostrictProvider configures CoStrict.
// CoStrict uses COSTRICT.md and .costrict/commands/ for slash commands.
type CostrictProvider struct{}

func (*CostrictProvider) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		initializers.NewDirectoryInitializer(".costrict/commands/spectr"),
		initializers.NewConfigFileInitializer(
			"COSTRICT.md",
			"instruction-pointer.md.tmpl",
		),
		initializers.NewSlashCommandsInitializer(
			".costrict/commands/spectr",
			".md",
			[]domain.SlashCommand{
				domain.SlashProposal,
				domain.SlashApply,
			},
		),
	}
}
