package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"

	//nolint:revive // long import path is unavoidable
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
)

// GeminiProvider configures Gemini CLI.
// Gemini uses TOML format for slash commands and has no instruction file.
type GeminiProvider struct{}

func (*GeminiProvider) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		initializers.NewDirectoryInitializer(".gemini/commands/spectr"),
		initializers.NewSlashCommandsInitializer(
			".gemini/commands/spectr",
			".toml",
			[]domain.SlashCommand{
				domain.SlashProposal,
				domain.SlashApply,
			},
		),
	}
}
