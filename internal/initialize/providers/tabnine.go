//nolint:revive // unused-receiver acceptable for interface compliance
package providers

func init() {
	Register(&TabnineProvider{})
}

// TabnineProvider implements the Provider interface for Tabnine.
// Tabnine uses .tabnine/commands/ for slash commands (no config file).
type TabnineProvider struct{}

// ID returns the unique identifier for this provider.
func (*TabnineProvider) ID() string { return "tabnine" }

// Name returns the human-readable name for display.
func (*TabnineProvider) Name() string { return "Tabnine" }

// Priority returns the display order (lower = higher priority).
func (*TabnineProvider) Priority() int { return PriorityTabnine }

// Initializers returns the file initializers for this provider.
func (*TabnineProvider) Initializers() []FileInitializer {
	proposalPath, applyPath := StandardCommandPaths(
		".tabnine/commands",
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
func (p *TabnineProvider) IsConfigured(projectPath string) bool {
	return AreInitializersConfigured(p.Initializers(), projectPath)
}

// GetFilePaths returns the file paths managed by this provider.
func (p *TabnineProvider) GetFilePaths() []string {
	return GetInitializerPaths(p.Initializers())
}
