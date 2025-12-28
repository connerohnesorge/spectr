package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"
)

// OpencodeProvider implements the Provider interface for OpenCode.
// OpenCode uses .opencode/commands/spectr/ for slash commands (no config file).
type OpencodeProvider struct{}

// Initializers returns the list of initializers for OpenCode.
func (*OpencodeProvider) Initializers(
	_ context.Context,
	tm TemplateManager,
) []Initializer { //nolint:lll
	return []Initializer{
		NewDirectoryInitializer(".opencode/commands/spectr"),
		NewSlashCommandsInitializer(
			".opencode/commands/spectr",
			map[domain.SlashCommand]domain.TemplateRef{
				domain.SlashProposal: tm.SlashCommand(domain.SlashProposal),
				domain.SlashApply:    tm.SlashCommand(domain.SlashApply),
			},
		),
	}
}
