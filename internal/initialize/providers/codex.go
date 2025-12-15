package providers

// init registers the Codex CLI provider with the global registry.
func init() {
	Register(&CodexProvider{})
}

// CodexProvider implements the Provider interface for Codex CLI.
// Codex uses AGENTS.md and global ~/.codex/prompts/ for commands.
type CodexProvider struct{}

// ID returns the unique identifier for the Codex CLI provider.
func (*CodexProvider) ID() string { return "codex" }

// Name returns the display name for Codex CLI.
func (*CodexProvider) Name() string { return "Codex CLI" }

// Priority returns the display order for Codex CLI.
func (*CodexProvider) Priority() int { return PriorityCodex }

// Initializers returns the file initializers for the Codex CLI provider.
func (*CodexProvider) Initializers() []FileInitializer {
	// Codex uses global paths, not project-relative paths
	proposalPath := "~/.codex/prompts/spectr-proposal.md"
	applyPath := "~/.codex/prompts/spectr-apply.md"

	return []FileInitializer{
		NewInstructionFileInitializer("AGENTS.md"),
		NewMarkdownSlashCommandInitializer(
			proposalPath,
			"proposal",
			FrontmatterProposal,
		),
		NewMarkdownSlashCommandInitializer(
			applyPath,
			"apply",
			FrontmatterApply,
		),
	}
}

func (p *CodexProvider) IsConfigured(projectPath string) bool {
	return AreInitializersConfigured(p.Initializers(), projectPath)
}

func (p *CodexProvider) GetFilePaths() []string {
	return GetInitializerPaths(p.Initializers())
}
