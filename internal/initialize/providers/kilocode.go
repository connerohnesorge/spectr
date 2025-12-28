package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
	"github.com/spf13/afero"
)

// KilocodeProvider configures Kilocode with slash commands in .kilocode/commands/spectr/.
// No config file for Kilocode.
// No init() - registration happens in RegisterAllProviders().
type KilocodeProvider struct{}

// Initializers returns the list of initializers for Kilocode.
// Receives TemplateManager to allow passing TemplateRef directly to initializers.
func (*KilocodeProvider) Initializers(_ context.Context, tm any) []Initializer {
	// Type assert tm to get TemplateManager methods
	type templateManager interface {
		SlashCommand(cmd domain.SlashCommand) domain.TemplateRef
	}

	tmgr, ok := tm.(templateManager)
	if !ok {
		return nil
	}

	return []Initializer{
		initializers.NewDirectoryInitializer(".kilocode/commands/spectr"),
		// No config file for Kilocode
		initializers.NewSlashCommandsInitializer(
			".kilocode/commands/spectr",
			map[domain.SlashCommand]domain.TemplateRef{
				domain.SlashProposal: tmgr.SlashCommand(domain.SlashProposal),
				domain.SlashApply:    tmgr.SlashCommand(domain.SlashApply),
			},
		),
	}
}

// IsConfigured checks if Kilocode is already configured in the project.
// Returns true if slash commands directory and files exist.
func (*KilocodeProvider) IsConfigured(projectDir string) bool {
	projectFs := afero.NewBasePathFs(afero.NewOsFs(), projectDir)

	// Check if slash commands directory exists
	exists, err := afero.DirExists(projectFs, ".kilocode/commands/spectr")
	if err != nil || !exists {
		return false
	}

	// Check if slash command files exist
	proposalExists, _ := afero.Exists(projectFs, ".kilocode/commands/spectr/proposal.md")
	applyExists, _ := afero.Exists(projectFs, ".kilocode/commands/spectr/apply.md")

	return proposalExists && applyExists
}
