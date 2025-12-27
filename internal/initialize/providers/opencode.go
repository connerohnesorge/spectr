package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/connerohnesorge/spectr/internal/templates"
)

// OpencodeProvider implements the Provider interface for OpenCode.
// OpenCode uses .opencode/command/spectr/ for slash commands.
// It has no instruction file as it uses JSON configuration.
type OpencodeProvider struct{}

// Initializers returns the initializers for OpenCode provider.
func (*OpencodeProvider) Initializers(
	_ context.Context,
	tm *templates.TemplateManager,
) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".opencode/command/spectr"),
		NewSlashCommandsInitializer(
			".opencode/command/spectr",
			map[domain.SlashCommand]domain.TemplateRef{
				domain.SlashProposal: tm.SlashCommand(domain.SlashProposal),
				domain.SlashApply:    tm.SlashCommand(domain.SlashApply),
			},
		),
	}
}
