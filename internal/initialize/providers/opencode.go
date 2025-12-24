package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"

	//nolint:revive // long import path is unavoidable
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
)

// OpencodeProvider configures OpenCode.
// OpenCode uses .opencode/command/spectr/ for slash commands (no config file).
type OpencodeProvider struct{}

func (*OpencodeProvider) Initializers(
	_ context.Context,
) []Initializer {
	return []Initializer{
		initializers.NewDirectoryInitializer(
			".opencode/command/spectr",
		),
		initializers.NewSlashCommandsInitializer(
			".opencode/command/spectr",
			".md",
			[]domain.SlashCommand{
				domain.SlashProposal,
				domain.SlashApply,
			},
		),
	}
}
