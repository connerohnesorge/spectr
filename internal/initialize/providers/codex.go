package providers

func init() {
	Register(NewCodexProvider())
}

// CodexProvider implements the Provider interface for Codex CLI.
// Codex uses AGENTS.md and global ~/.codex/prompts/spectr/ for commands.
type CodexProvider struct {
	BaseProvider
}

// NewCodexProvider creates a new Codex CLI provider.
func NewCodexProvider() *CodexProvider {
	// Codex uses global paths, not project-relative paths
	proposalPath := "~/.codex/prompts/spectr-proposal.md"
	syncPath := "~/.codex/prompts/spectr-sync.md"
	applyPath := "~/.codex/prompts/spectr-apply.md"

	return &CodexProvider{
		BaseProvider: BaseProvider{
			id:            "codex",
			name:          "Codex CLI",
			priority:      PriorityCodex,
			configFile:    "AGENTS.md",
			proposalPath:  proposalPath,
			syncPath:      syncPath,
			applyPath:     applyPath,
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}
