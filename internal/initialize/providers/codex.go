package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/initialize/templates"
)

func init() {
	_ = RegisterProvider(Registration{
		ID:       "codex",
		Name:     "Codex CLI",
		Priority: PriorityCodex,
		Provider: &CodexProvider{},
	})
}

// CodexProvider implements the Provider interface for Codex CLI.
// Codex uses AGENTS.md and global ~/.codex/prompts/ for commands.
type CodexProvider struct{}

func (*CodexProvider) Initializers(
	_ context.Context,
) []Initializer {
	return []Initializer{
		NewGlobalDirectoryInitializer(
			".codex/prompts",
		),
		NewConfigFileInitializer(
			"AGENTS.md",
			func(tm TemplateManager) any {
				return tm.Agents()
			},
		),
		NewGlobalSlashCommandsInitializer(
			".codex/prompts",
			".md",
			[]templates.SlashCommand{
				templates.SlashProposal,
				templates.SlashApply,
			},
		),
	}
}
