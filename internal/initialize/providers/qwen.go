package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"
)

// QwenProvider implements the Provider interface for Qwen Code.
// Qwen Code uses QWEN.md and .qwen/commands/spectr/ for slash commands.
type QwenProvider struct{}

// Initializers returns the list of initializers for Qwen Code.
func (*QwenProvider) Initializers(
	_ context.Context,
	tm TemplateManager,
) []Initializer { //nolint:lll
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
