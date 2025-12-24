package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"

	//nolint:revive // long import path is unavoidable
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
)

// ClineProvider configures Cline.
// Cline uses CLINE.md and .clinerules/commands/ for slash commands.
type ClineProvider struct{}

func (*ClineProvider) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		initializers.NewDirectoryInitializer(".clinerules/commands/spectr"),
		initializers.NewConfigFileInitializer(
			"CLINE.md",
			"instruction-pointer.md.tmpl",
		),
		initializers.NewSlashCommandsInitializer(
			".clinerules/commands/spectr",
			".md",
			[]domain.SlashCommand{
				domain.SlashProposal,
				domain.SlashApply,
			},
		),
	}
}
