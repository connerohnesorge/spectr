package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"

	//nolint:revive // long import path is unavoidable
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
)

// QoderProvider configures Qoder.
// Qoder uses QODER.md and .qoder/commands/ for slash commands.
type QoderProvider struct{}

func (*QoderProvider) Initializers(
	_ context.Context,
) []Initializer {
	return []Initializer{
		initializers.NewDirectoryInitializer(
			".qoder/commands/spectr",
		),
		initializers.NewConfigFileInitializer(
			"QODER.md",
			"instruction-pointer.md.tmpl",
		),
		initializers.NewSlashCommandsInitializer(
			".qoder/commands/spectr",
			".md",
			[]domain.SlashCommand{
				domain.SlashProposal,
				domain.SlashApply,
			},
		),
	}
}
