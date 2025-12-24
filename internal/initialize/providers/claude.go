package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers" //nolint:revive
)

// ClaudeProvider configures Claude Code.
type ClaudeProvider struct{}

func (*ClaudeProvider) Initializers(
	_ context.Context,
) []Initializer {
	return []Initializer{
		initializers.NewDirectoryInitializer(
			".claude/commands/spectr",
		),
		initializers.NewConfigFileInitializer(
			"CLAUDE.md",
			"instruction-pointer.md.tmpl",
		),
		initializers.NewSlashCommandsInitializer(
			".claude/commands/spectr",
			".md",
			[]domain.SlashCommand{
				domain.SlashProposal,
				domain.SlashApply,
			},
		),
	}
}
