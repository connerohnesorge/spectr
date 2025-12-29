package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"
)

// GeminiProvider implements the Provider interface for Gemini CLI.
// Gemini uses .gemini/commands/spectr/ for TOML-based slash commands (no instruction file). //nolint:lll
type GeminiProvider struct{}

// Initializers returns the list of initializers for Gemini CLI.
func (*GeminiProvider) Initializers(
	_ context.Context,
	tm TemplateManager,
) []Initializer { //nolint:lll // Constructor calls with template refs exceed line limit
	return []Initializer{
		NewDirectoryInitializer(".gemini/commands/spectr"),
		NewTOMLSlashCommandsInitializer(
			".gemini/commands/spectr",
			map[domain.SlashCommand]domain.TemplateRef{
				domain.SlashProposal: tm.TOMLSlashCommand(domain.SlashProposal),
				domain.SlashApply:    tm.TOMLSlashCommand(domain.SlashApply),
			},
		),
	}
}
