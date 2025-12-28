package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"
)

// WindsurfProvider implements the Provider interface for Windsurf.
// Windsurf uses .windsurf/commands/spectr/ for slash commands (no config file).
type WindsurfProvider struct{}

// Initializers returns the list of initializers for Windsurf.
func (*WindsurfProvider) Initializers(
	_ context.Context,
	tm TemplateManager,
) []Initializer { //nolint:lll
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
