//nolint:revive // unused-receiver acceptable for interface compliance
package providers

func init() {
	Register(&CodexProvider{})
}

// CodexProvider implements the Provider interface for Codex CLI.
// Codex uses AGENTS.md and global ~/.codex/prompts/spectr/ for commands.
type CodexProvider struct{}

// ID returns the unique identifier for this provider.
func (*CodexProvider) ID() string { return "codex" }

// Name returns the human-readable name for display.
func (*CodexProvider) Name() string { return "Codex CLI" }

// Priority returns the display order (lower = higher priority).
func (*CodexProvider) Priority() int { return PriorityCodex }

// Initializers returns the file initializers for this provider.
func (*CodexProvider) Initializers() []FileInitializer {
	// Codex uses global paths, not project-relative paths
	proposalPath := "~/.codex/prompts/spectr-proposal.md"
	applyPath := "~/.codex/prompts/spectr-apply.md"

	return []FileInitializer{
		NewInstructionFileInitializer("AGENTS.md"),
		NewMarkdownSlashCommandInitializer(
			proposalPath,
			"proposal",
			StandardProposalFrontmatter,
		),
		NewMarkdownSlashCommandInitializer(
			applyPath,
			"apply",
			StandardApplyFrontmatter,
		),
	}
}

// IsConfigured checks if all files for this provider exist.
func (p *CodexProvider) IsConfigured(projectPath string) bool {
	return AreInitializersConfigured(p.Initializers(), projectPath)
}

// GetFilePaths returns the file paths managed by this provider.
func (p *CodexProvider) GetFilePaths() []string {
	return GetInitializerPaths(p.Initializers())
}
