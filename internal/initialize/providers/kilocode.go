package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/connerohnesorge/spectr/internal/templates"
)

// KilocodeProvider implements the Provider interface for Kilocode.
// Kilocode uses .kilocode/commands/spectr/ for slash commands (no config file).
type KilocodeProvider struct{}

// Initializers returns the initializers for Kilocode provider.
func (*KilocodeProvider) Initializers(
	_ context.Context,
	tm *templates.TemplateManager,
) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".kilocode/commands/spectr"),
		NewSlashCommandsInitializer(
			".kilocode/commands/spectr",
			map[domain.SlashCommand]domain.TemplateRef{
				domain.SlashProposal: tm.SlashCommand(domain.SlashProposal),
				domain.SlashApply:    tm.SlashCommand(domain.SlashApply),
			},
		),
	}
}
