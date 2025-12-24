package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/initialize/templates"
)

func init() {
	_ = RegisterProvider(Registration{
		ID:       "gemini",
		Name:     "Gemini CLI",
		Priority: PriorityGemini,
		Provider: &GeminiProvider{},
	})
}

// GeminiProvider implements the Provider interface for Gemini CLI.
// Gemini uses .gemini/commands/ for TOML-based slash commands
// (no instruction file).
type GeminiProvider struct{}

func (*GeminiProvider) Initializers(
	_ context.Context,
) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(
			".gemini/commands/spectr",
		),
		NewSlashCommandsInitializer(
			".gemini/commands/spectr",
			".toml",
			[]templates.SlashCommand{
				templates.SlashProposal,
				templates.SlashApply,
			},
		),
	}
}
