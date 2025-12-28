package providers

import (
	"context"
	"strings"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
	"github.com/spf13/afero"
)

// QoderProvider configures Qoder with QODER.md and .qoder/commands/spectr/.
// No init() - registration happens in RegisterAllProviders().
type QoderProvider struct{}

// Initializers returns the list of initializers for Qoder.
// Receives TemplateManager to allow passing TemplateRef directly to initializers.
func (*QoderProvider) Initializers(_ context.Context, tm any) []Initializer {
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
		initializers.NewDirectoryInitializer(".qoder/commands/spectr"),
		initializers.NewConfigFileInitializer("QODER.md", tmgr.InstructionPointer()),
		initializers.NewSlashCommandsInitializer(
			".qoder/commands/spectr",
			map[domain.SlashCommand]domain.TemplateRef{
				domain.SlashProposal: tmgr.SlashCommand(domain.SlashProposal),
				domain.SlashApply:    tmgr.SlashCommand(domain.SlashApply),
			},
		),
	}
}

// IsConfigured checks if Qoder is already configured in the project.
// Returns true if QODER.md exists with spectr markers and slash commands exist.
func (*QoderProvider) IsConfigured(projectDir string) bool {
	projectFs := afero.NewBasePathFs(afero.NewOsFs(), projectDir)

	// Check if QODER.md exists with spectr markers
	content, err := afero.ReadFile(projectFs, "QODER.md")
	if err != nil {
		return false
	}
	contentLower := strings.ToLower(string(content))
	if !strings.Contains(contentLower, "<!-- spectr:start -->") {
		return false
	}

	// Check if slash commands directory exists
	exists, err := afero.DirExists(projectFs, ".qoder/commands/spectr")
	if err != nil || !exists {
		return false
	}

	// Check if slash command files exist
	proposalExists, _ := afero.Exists(projectFs, ".qoder/commands/spectr/proposal.md")
	applyExists, _ := afero.Exists(projectFs, ".qoder/commands/spectr/apply.md")

	return proposalExists && applyExists
}
