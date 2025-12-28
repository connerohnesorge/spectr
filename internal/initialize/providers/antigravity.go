package providers

import (
	"context"
	"strings"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
	"github.com/spf13/afero"
)

// AntigravityProvider configures Antigravity with AGENTS.md and .agent/workflows/ for slash commands.
// Uses PrefixedSlashCommandsInitializer with prefix "spectr-" for files like spectr-proposal.md.
// No init() - registration happens in RegisterAllProviders().
type AntigravityProvider struct{}

// Initializers returns the list of initializers for Antigravity.
// Receives TemplateManager to allow passing TemplateRef directly to initializers.
func (*AntigravityProvider) Initializers(_ context.Context, tm any) []Initializer {
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
		initializers.NewDirectoryInitializer(".agent/workflows"),
		initializers.NewConfigFileInitializer("AGENTS.md", tmgr.InstructionPointer()),
		// Uses PrefixedSlashCommandsInitializer with prefix "spectr-"
		// Output: .agent/workflows/spectr-proposal.md, .agent/workflows/spectr-apply.md
		initializers.NewPrefixedSlashCommandsInitializer(
			".agent/workflows",
			"spectr-",
			map[domain.SlashCommand]domain.TemplateRef{
				domain.SlashProposal: tmgr.SlashCommand(domain.SlashProposal),
				domain.SlashApply:    tmgr.SlashCommand(domain.SlashApply),
			},
		),
	}
}

// IsConfigured checks if Antigravity is already configured in the project.
// Returns true if AGENTS.md exists with spectr markers and slash commands exist.
func (*AntigravityProvider) IsConfigured(projectDir string) bool {
	projectFs := afero.NewBasePathFs(afero.NewOsFs(), projectDir)

	// Check if AGENTS.md exists with spectr markers
	content, err := afero.ReadFile(projectFs, "AGENTS.md")
	if err != nil {
		return false
	}
	contentLower := strings.ToLower(string(content))
	if !strings.Contains(contentLower, "<!-- spectr:start -->") {
		return false
	}

	// Check if slash commands directory exists
	exists, err := afero.DirExists(projectFs, ".agent/workflows")
	if err != nil || !exists {
		return false
	}

	// Check if slash command files exist (prefixed format)
	proposalExists, _ := afero.Exists(projectFs, ".agent/workflows/spectr-proposal.md")
	applyExists, _ := afero.Exists(projectFs, ".agent/workflows/spectr-apply.md")

	return proposalExists && applyExists
}
