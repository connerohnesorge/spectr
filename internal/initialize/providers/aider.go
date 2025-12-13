//nolint:revive // unused-receiver acceptable for interface compliance
package providers

func init() {
	Register(&AiderProvider{})
}

// AiderProvider implements the Provider interface for Aider.
// Aider uses .aider/commands/ for slash commands (no config file).
type AiderProvider struct{}

// ID returns the unique identifier for this provider.
func (*AiderProvider) ID() string { return "aider" }

// Name returns the human-readable name for display.
func (*AiderProvider) Name() string { return "Aider" }

// Priority returns the display order (lower = higher priority).
func (*AiderProvider) Priority() int { return PriorityAider }

// Initializers returns the file initializers for this provider.
func (*AiderProvider) Initializers() []FileInitializer {
	proposalPath, applyPath := StandardCommandPaths(
		".aider/commands",
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
func (p *AiderProvider) IsConfigured(projectPath string) bool {
	return AreInitializersConfigured(p.Initializers(), projectPath)
}

// GetFilePaths returns the file paths managed by this provider.
func (p *AiderProvider) GetFilePaths() []string {
	return GetInitializerPaths(p.Initializers())
}
