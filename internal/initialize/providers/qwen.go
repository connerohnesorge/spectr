package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/connerohnesorge/spectr/internal/templates"
)

// QwenProvider implements the Provider interface for Qwen Code.
// Qwen uses QWEN.md and .qwen/commands/spectr/ for slash commands.
type QwenProvider struct{}

// Initializers returns the initializers for Qwen Code provider.
func (*QwenProvider) Initializers(
	_ context.Context,
	tm *templates.TemplateManager,
) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".qwen/commands/spectr"),
		NewConfigFileInitializer("QWEN.md", tm.InstructionPointer()),
		NewSlashCommandsInitializer(
			".qwen/commands/spectr",
			map[domain.SlashCommand]domain.TemplateRef{
				domain.SlashProposal: tm.SlashCommand(domain.SlashProposal),
				domain.SlashApply:    tm.SlashCommand(domain.SlashApply),
			},
		),
	}
}
