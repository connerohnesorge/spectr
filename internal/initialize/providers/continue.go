//nolint:revive // unused-receiver acceptable for interface compliance
package providers

func init() {
	Register(&ContinueProvider{})
}

// ContinueProvider implements the Provider interface for Continue.
// Continue uses .continue/commands/ for slash commands (no config file).
type ContinueProvider struct{}

// ID returns the unique identifier for this provider.
func (*ContinueProvider) ID() string { return "continue" }

// Name returns the human-readable name for display.
func (*ContinueProvider) Name() string { return "Continue" }

// Priority returns the display order (lower = higher priority).
func (*ContinueProvider) Priority() int { return PriorityContinue }

// Initializers returns the file initializers for this provider.
func (*ContinueProvider) Initializers() []FileInitializer {
	proposalPath, applyPath := StandardCommandPaths(
		".continue/commands",
		".md",
	)

	return []FileInitializer{
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
func (p *ContinueProvider) IsConfigured(projectPath string) bool {
	return AreInitializersConfigured(p.Initializers(), projectPath)
}

// GetFilePaths returns the file paths managed by this provider.
func (p *ContinueProvider) GetFilePaths() []string {
	return GetInitializerPaths(p.Initializers())
}
