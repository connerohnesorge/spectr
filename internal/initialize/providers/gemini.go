package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
	"github.com/spf13/afero"
)

// GeminiProvider configures Gemini CLI with TOML slash commands in .gemini/commands/spectr/.
// No config file - uses TOML slash commands only.
// No init() - registration happens in RegisterAllProviders().
type GeminiProvider struct{}

// Initializers returns the list of initializers for Gemini CLI.
// Receives TemplateManager to allow passing TemplateRef directly to initializers.
func (*GeminiProvider) Initializers(_ context.Context, tm any) []Initializer {
	// Type assert tm to get TemplateManager methods
	type templateManager interface {
		TOMLSlashCommand(cmd domain.SlashCommand) domain.TemplateRef
	}

	tmgr, ok := tm.(templateManager)
	if !ok {
		return nil
	}

	return []Initializer{
		initializers.NewDirectoryInitializer(".gemini/commands/spectr"),
		// No config file for Gemini - uses TOML slash commands only
		initializers.NewTOMLSlashCommandsInitializer(
			".gemini/commands/spectr",
			map[domain.SlashCommand]domain.TemplateRef{
				domain.SlashProposal: tmgr.TOMLSlashCommand(domain.SlashProposal),
				domain.SlashApply:    tmgr.TOMLSlashCommand(domain.SlashApply),
			},
		),
	}
}

// IsConfigured checks if Gemini CLI is already configured in the project.
// Returns true if slash commands directory and files exist.
func (*GeminiProvider) IsConfigured(projectDir string) bool {
	projectFs := afero.NewBasePathFs(afero.NewOsFs(), projectDir)

	// Check if slash commands directory exists
	exists, err := afero.DirExists(projectFs, ".gemini/commands/spectr")
	if err != nil || !exists {
		return false
	}

	// Check if slash command files exist
	proposalExists, _ := afero.Exists(projectFs, ".gemini/commands/spectr/proposal.toml")
	applyExists, _ := afero.Exists(projectFs, ".gemini/commands/spectr/apply.toml")

	return proposalExists && applyExists
}
