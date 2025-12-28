package providers

import (
	"context"
	"strings"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
	"github.com/spf13/afero"
)

// CrushProvider configures Crush with CRUSH.md and .crush/commands/spectr/.
// No init() - registration happens in RegisterAllProviders().
type CrushProvider struct{}

// Initializers returns the list of initializers for Crush.
// Receives TemplateManager to allow passing TemplateRef directly to initializers.
func (*CrushProvider) Initializers(_ context.Context, tm any) []Initializer {
	// Type assert tm to get TemplateManager methods
	type templateManager interface {
		InstructionPointer() domain.TemplateRef
		SlashCommand(cmd domain.SlashCommand) domain.TemplateRef
	}

	tmgr, ok := tm.(templateManager)
	if !ok {
		return nil
	}

	return []Initializer{
		initializers.NewDirectoryInitializer(".crush/commands/spectr"),
		initializers.NewConfigFileInitializer("CRUSH.md", tmgr.InstructionPointer()),
		initializers.NewSlashCommandsInitializer(
			".crush/commands/spectr",
			map[domain.SlashCommand]domain.TemplateRef{
				domain.SlashProposal: tmgr.SlashCommand(domain.SlashProposal),
				domain.SlashApply:    tmgr.SlashCommand(domain.SlashApply),
			},
		),
	}
}

// IsConfigured checks if Crush is already configured in the project.
// Returns true if CRUSH.md exists with spectr markers and slash commands exist.
func (*CrushProvider) IsConfigured(projectDir string) bool {
	projectFs := afero.NewBasePathFs(afero.NewOsFs(), projectDir)

	// Check if CRUSH.md exists with spectr markers
	content, err := afero.ReadFile(projectFs, "CRUSH.md")
	if err != nil {
		return false
	}
	contentLower := strings.ToLower(string(content))
	if !strings.Contains(contentLower, "<!-- spectr:start -->") {
		return false
	}

	// Check if slash commands directory exists
	exists, err := afero.DirExists(projectFs, ".crush/commands/spectr")
	if err != nil || !exists {
		return false
	}

	// Check if slash command files exist
	proposalExists, _ := afero.Exists(projectFs, ".crush/commands/spectr/proposal.md")
	applyExists, _ := afero.Exists(projectFs, ".crush/commands/spectr/apply.md")

	return proposalExists && applyExists
}
