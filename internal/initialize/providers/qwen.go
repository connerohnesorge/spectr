package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/initialize/templates"
)

func init() {
	_ = RegisterProvider(Registration{
		ID:       "qwen",
		Name:     "Qwen Code",
		Priority: PriorityQwen,
		Provider: &QwenProvider{},
	})
}

// QwenProvider implements the Provider interface for Qwen Code.
// Qwen uses QWEN.md and .qwen/commands/ for slash commands.
type QwenProvider struct{}

func (*QwenProvider) Initializers(
	_ context.Context,
) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(
			".qwen/commands/spectr",
		),
		NewConfigFileInitializer(
			"QWEN.md",
			func(tm TemplateManager) any {
				return tm.InstructionPointer()
			},
		),
		NewSlashCommandsInitializer(
			".qwen/commands/spectr",
			".md",
			[]templates.SlashCommand{
				templates.SlashProposal,
				templates.SlashApply,
			},
		),
	}
}
