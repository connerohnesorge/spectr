package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/connerohnesorge/spectr/internal/templates"
)

// WindsurfProvider implements the Provider interface for Windsurf.
// Windsurf uses .windsurf/commands/spectr/ for slash commands (no config file).
type WindsurfProvider struct{}

// Initializers returns the initializers for Windsurf provider.
func (*WindsurfProvider) Initializers(
	_ context.Context,
	tm *templates.TemplateManager,
) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".windsurf/commands/spectr"),
		NewSlashCommandsInitializer(
			".windsurf/commands/spectr",
			map[domain.SlashCommand]domain.TemplateRef{
				domain.SlashProposal: tm.SlashCommand(domain.SlashProposal),
				domain.SlashApply:    tm.SlashCommand(domain.SlashApply),
			},
		),
	}
}
