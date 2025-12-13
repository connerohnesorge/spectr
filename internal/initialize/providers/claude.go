//nolint:revive // unused-receiver acceptable for interface compliance
package providers

func init() {
	Register(&ClaudeProvider{})
}

// ClaudeProvider implements the Provider interface for Claude Code.
// Claude Code uses CLAUDE.md and .claude/commands/ for slash commands.
type ClaudeProvider struct{}

// ID returns the unique identifier for this provider.
func (*ClaudeProvider) ID() string { return "claude-code" }

// Name returns the human-readable name for display.
func (*ClaudeProvider) Name() string { return "Claude Code" }

// Priority returns the display order (lower = higher priority).
func (*ClaudeProvider) Priority() int { return PriorityClaudeCode }

// Initializers returns the file initializers for this provider.
func (*ClaudeProvider) Initializers() []FileInitializer {
	proposalPath, applyPath := StandardCommandPaths(
		".claude/commands",
		".md",
	)

	return []FileInitializer{
		NewInstructionFileInitializer("CLAUDE.md"),
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
func (p *ClaudeProvider) IsConfigured(projectPath string) bool {
	return AreInitializersConfigured(p.Initializers(), projectPath)
}

// GetFilePaths returns the file paths managed by this provider.
func (p *ClaudeProvider) GetFilePaths() []string {
	return GetInitializerPaths(p.Initializers())
}
