package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"
)

// GeminiProvider implements the Provider interface for Gemini CLI.
// Gemini uses GEMINI.md for instructions and .gemini/commands/spectr/ for TOML-based slash commands. //nolint:lll
type GeminiProvider struct{}

// Initializers returns the list of initializers for Gemini CLI.
func (*GeminiProvider) Initializers(
	_ context.Context,
	tm TemplateManager,
) []Initializer { //nolint:lll // Constructor calls with template refs exceed line limit
	return []Initializer{
		NewDirectoryInitializer(".gemini/commands/spectr"),
		NewDirectoryInitializer(".gemini/skills"),
		NewConfigFileInitializer("GEMINI.md", tm.InstructionPointer()),
		NewTOMLSlashCommandsInitializer(
			".gemini/commands/spectr",
			map[domain.SlashCommand]domain.TemplateRef{
				domain.SlashProposal: tm.TOMLSlashCommand(
					domain.SlashProposal,
				),
				domain.SlashApply: tm.TOMLSlashCommand(
					domain.SlashApply,
				),
			},
		),
		NewAgentSkillsInitializer(
			"spectr-accept-wo-spectr-bin",
			".gemini/skills/spectr-accept-wo-spectr-bin",
			tm,
		),
		NewAgentSkillsInitializer(
			"spectr-validate-wo-spectr-bin",
			".gemini/skills/spectr-validate-wo-spectr-bin",
			tm,
		),
	}
}
