package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/connerohnesorge/spectr/internal/templates"
)

// AiderProvider implements the Provider interface for Aider.
// Aider uses .aider/commands/spectr/ for slash commands (no config file).
type AiderProvider struct{}

// Initializers returns the initializers for Aider provider.
func (*AiderProvider) Initializers(
	_ context.Context,
	tm *templates.TemplateManager,
) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".aider/commands/spectr"),
		NewSlashCommandsInitializer(
			".aider/commands/spectr",
			map[domain.SlashCommand]domain.TemplateRef{
				domain.SlashProposal: tm.SlashCommand(domain.SlashProposal),
				domain.SlashApply:    tm.SlashCommand(domain.SlashApply),
			},
		),
	}
}
