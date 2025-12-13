//nolint:revive // unused-receiver acceptable for interface compliance
package providers

func init() {
	Register(&CodeBuddyProvider{})
}

// CodeBuddyProvider implements the Provider interface for CodeBuddy.
// CodeBuddy uses CODEBUDDY.md and .codebuddy/commands/ for slash commands.
type CodeBuddyProvider struct{}

// ID returns the unique identifier for this provider.
func (*CodeBuddyProvider) ID() string { return "codebuddy" }

// Name returns the human-readable name for display.
func (*CodeBuddyProvider) Name() string { return "CodeBuddy" }

// Priority returns the display order (lower = higher priority).
func (*CodeBuddyProvider) Priority() int { return PriorityCodeBuddy }

// Initializers returns the file initializers for this provider.
func (*CodeBuddyProvider) Initializers() []FileInitializer {
	proposalPath, applyPath := StandardCommandPaths(
		".codebuddy/commands",
		".md",
	)

	return []FileInitializer{
		NewInstructionFileInitializer("CODEBUDDY.md"),
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
func (p *CodeBuddyProvider) IsConfigured(projectPath string) bool {
	return AreInitializersConfigured(p.Initializers(), projectPath)
}

// GetFilePaths returns the file paths managed by this provider.
func (p *CodeBuddyProvider) GetFilePaths() []string {
	return GetInitializerPaths(p.Initializers())
}
