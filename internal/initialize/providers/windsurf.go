//nolint:revive // unused-receiver acceptable for interface compliance
package providers

func init() {
	Register(&WindsurfProvider{})
}

// WindsurfProvider implements the Provider interface for Windsurf.
// Windsurf uses .windsurf/commands/ for slash commands (no config file).
type WindsurfProvider struct{}

// ID returns the unique identifier for this provider.
func (*WindsurfProvider) ID() string { return "windsurf" }

// Name returns the human-readable name for display.
func (*WindsurfProvider) Name() string { return "Windsurf" }

// Priority returns the display order (lower = higher priority).
func (*WindsurfProvider) Priority() int { return PriorityWindsurf }

// Initializers returns the file initializers for this provider.
func (*WindsurfProvider) Initializers() []FileInitializer {
	proposalPath, applyPath := StandardCommandPaths(
		".windsurf/commands",
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
func (p *WindsurfProvider) IsConfigured(projectPath string) bool {
	return AreInitializersConfigured(p.Initializers(), projectPath)
}

// GetFilePaths returns the file paths managed by this provider.
func (p *WindsurfProvider) GetFilePaths() []string {
	return GetInitializerPaths(p.Initializers())
}
