package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/connerohnesorge/spectr/internal/templates"
)

// GeminiProvider implements the Provider interface for Gemini CLI.
// Gemini uses .gemini/commands/spectr/ for TOML-based slash commands
// (no instruction file).
type GeminiProvider struct{}

// Initializers returns the initializers for Gemini CLI provider.
func (*GeminiProvider) Initializers(
	_ context.Context,
	tm *templates.TemplateManager,
) []Initializer {
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
