package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers" //nolint:revive
)

// KilocodeProvider configures Kilocode.
// Kilocode uses .kilocode/commands/ for slash commands (no config file).
type KilocodeProvider struct{}

func (*KilocodeProvider) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		initializers.NewDirectoryInitializer(".kilocode/commands/spectr"),
		initializers.NewSlashCommandsInitializer(
			".kilocode/commands/spectr",
			".md",
			[]domain.SlashCommand{
				domain.SlashProposal,
				domain.SlashApply,
			},
		),
	}
}
