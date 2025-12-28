package providers

import (
	"context"
	"os"
	"strings"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
	"github.com/spf13/afero"
)

// ClaudeProvider configures Claude Code with CLAUDE.md and .claude/commands/spectr/.
// No init() - registration happens in RegisterAllProviders().
type ClaudeProvider struct{}

// Initializers returns the list of initializers for Claude Code.
// Receives TemplateManager to allow passing TemplateRef directly to initializers.
func (*ClaudeProvider) Initializers(_ context.Context, tm any) []Initializer {
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
		initializers.NewDirectoryInitializer(".claude/commands/spectr"),
		initializers.NewConfigFileInitializer("CLAUDE.md", tmgr.InstructionPointer()),
		initializers.NewSlashCommandsInitializer(
			".claude/commands/spectr",
			map[domain.SlashCommand]domain.TemplateRef{
				domain.SlashProposal: tmgr.SlashCommand(domain.SlashProposal),
				domain.SlashApply:    tmgr.SlashCommand(domain.SlashApply),
			},
		),
	}
}

// IsConfigured checks if Claude Code is already configured in the project.
// Returns true if CLAUDE.md exists with spectr markers and slash commands exist.
func (*ClaudeProvider) IsConfigured(projectDir string) bool {
	projectFs := afero.NewBasePathFs(afero.NewOsFs(), projectDir)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	homeFs := afero.NewBasePathFs(afero.NewOsFs(), homeDir)
	cfg := &domain.Config{SpectrDir: "spectr"}

	// Check if CLAUDE.md exists with spectr markers
	content, err := afero.ReadFile(projectFs, "CLAUDE.md")
	if err != nil {
		return false
	}
	contentLower := strings.ToLower(string(content))
	if !strings.Contains(contentLower, "<!-- spectr:start -->") {
		return false
	}

	// Check if slash commands directory exists
	exists, err := afero.DirExists(projectFs, ".claude/commands/spectr")
	if err != nil || !exists {
		return false
	}

	// Check if slash command files exist
	proposalExists, _ := afero.Exists(projectFs, ".claude/commands/spectr/proposal.md")
	applyExists, _ := afero.Exists(projectFs, ".claude/commands/spectr/apply.md")
	if !proposalExists || !applyExists {
		return false
	}

	// Suppress unused variable warnings
	_ = homeFs
	_ = cfg

	return true
}
