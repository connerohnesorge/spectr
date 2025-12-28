package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
	"github.com/spf13/afero"
)

// CursorProvider configures Cursor with slash commands in .cursorrules/commands/spectr/.
// No config file for Cursor.
// No init() - registration happens in RegisterAllProviders().
type CursorProvider struct{}

// Initializers returns the list of initializers for Cursor.
// Receives TemplateManager to allow passing TemplateRef directly to initializers.
func (*CursorProvider) Initializers(_ context.Context, tm any) []Initializer {
	// Type assert tm to get TemplateManager methods
	type templateManager interface {
		SlashCommand(cmd domain.SlashCommand) domain.TemplateRef
	}

	tmgr, ok := tm.(templateManager)
	if !ok {
		return nil
	}

	return []Initializer{
		initializers.NewDirectoryInitializer(".cursorrules/commands/spectr"),
		// No config file for Cursor
		initializers.NewSlashCommandsInitializer(
			".cursorrules/commands/spectr",
			map[domain.SlashCommand]domain.TemplateRef{
				domain.SlashProposal: tmgr.SlashCommand(domain.SlashProposal),
				domain.SlashApply:    tmgr.SlashCommand(domain.SlashApply),
			},
		),
	}
}

// IsConfigured checks if Cursor is already configured in the project.
// Returns true if slash commands directory and files exist.
func (*CursorProvider) IsConfigured(projectDir string) bool {
	projectFs := afero.NewBasePathFs(afero.NewOsFs(), projectDir)

	// Check if slash commands directory exists
	exists, err := afero.DirExists(projectFs, ".cursorrules/commands/spectr")
	if err != nil || !exists {
		return false
	}

	// Check if slash command files exist
	proposalExists, _ := afero.Exists(projectFs, ".cursorrules/commands/spectr/proposal.md")
	applyExists, _ := afero.Exists(projectFs, ".cursorrules/commands/spectr/apply.md")

	return proposalExists && applyExists
}
