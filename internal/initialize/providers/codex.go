package providers

import (
	"context"
	"os"
	"strings"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
	"github.com/spf13/afero"
)

// CodexProvider configures Codex CLI with AGENTS.md and slash commands in ~/.codex/prompts/.
// Uses HomeDirectoryInitializer and HomePrefixedSlashCommandsInitializer for home filesystem paths.
// Uses prefix "spectr-" for files like spectr-proposal.md.
// No init() - registration happens in RegisterAllProviders().
type CodexProvider struct{}

// Initializers returns the list of initializers for Codex CLI.
// Receives TemplateManager to allow passing TemplateRef directly to initializers.
func (*CodexProvider) Initializers(_ context.Context, tm any) []Initializer {
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
		// Uses HomeDirectoryInitializer for ~/.codex/prompts/
		initializers.NewHomeDirectoryInitializer(".codex/prompts"),
		// Uses config file AGENTS.md in project directory
		initializers.NewConfigFileInitializer("AGENTS.md", tmgr.InstructionPointer()),
		// Uses HomePrefixedSlashCommandsInitializer with prefix "spectr-" in home filesystem
		// Output: ~/.codex/prompts/spectr-proposal.md, ~/.codex/prompts/spectr-apply.md
		initializers.NewHomePrefixedSlashCommandsInitializer(
			".codex/prompts",
			"spectr-",
			map[domain.SlashCommand]domain.TemplateRef{
				domain.SlashProposal: tmgr.SlashCommand(domain.SlashProposal),
				domain.SlashApply:    tmgr.SlashCommand(domain.SlashApply),
			},
		),
	}
}

// IsConfigured checks if Codex CLI is already configured in the project.
// Returns true if AGENTS.md exists with spectr markers and slash commands exist in home directory.
func (*CodexProvider) IsConfigured(projectDir string) bool {
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

	// Check if home directory slash commands exist
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	homeFs := afero.NewBasePathFs(afero.NewOsFs(), homeDir)

	// Check if slash commands directory exists in home
	exists, err := afero.DirExists(homeFs, ".codex/prompts")
	if err != nil || !exists {
		return false
	}

	// Check if slash command files exist (prefixed format)
	proposalExists, _ := afero.Exists(homeFs, ".codex/prompts/spectr-proposal.md")
	applyExists, _ := afero.Exists(homeFs, ".codex/prompts/spectr-apply.md")

	return proposalExists && applyExists
}
