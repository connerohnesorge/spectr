package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/connerohnesorge/spectr/internal/templates"
)

// ContinueProvider implements the Provider interface for Continue.
// Continue uses .continue/commands/spectr/ for slash commands (no config file).
type ContinueProvider struct{}

// Initializers returns the initializers for Continue provider.
func (*ContinueProvider) Initializers(
	_ context.Context,
	tm *templates.TemplateManager,
) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".continue/commands/spectr"),
		NewSlashCommandsInitializer(
			".continue/commands/spectr",
			map[domain.SlashCommand]domain.TemplateRef{
				domain.SlashProposal: tm.SlashCommand(domain.SlashProposal),
				domain.SlashApply:    tm.SlashCommand(domain.SlashApply),
			},
		),
	}
}
