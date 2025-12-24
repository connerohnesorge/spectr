package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"

	//nolint:revive // long import path is unavoidable
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
)

// AntigravityProvider configures Antigravity.
// Antigravity uses AGENTS.md and .agent/workflows/spectr/ for slash commands.
type AntigravityProvider struct{}

func (*AntigravityProvider) Initializers(
	_ context.Context,
) []Initializer {
	return []Initializer{
		initializers.NewDirectoryInitializer(
			".agent/workflows/spectr",
		),
		initializers.NewConfigFileInitializer(
			"AGENTS.md",
			"AGENTS.md.tmpl",
		),
		initializers.NewSlashCommandsInitializer(
			".agent/workflows/spectr",
			".md",
			[]domain.SlashCommand{
				domain.SlashProposal,
				domain.SlashApply,
			},
		),
	}
}
