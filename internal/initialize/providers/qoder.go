//nolint:revive // unused-receiver acceptable for interface compliance
package providers

func init() {
	Register(&QoderProvider{})
}

// QoderProvider implements the Provider interface for Qoder.
// Qoder uses QODER.md and .qoder/commands/ for slash commands.
type QoderProvider struct{}

// ID returns the unique identifier for this provider.
func (p *QoderProvider) ID() string { return "qoder" }

// Name returns the human-readable name for display.
func (p *QoderProvider) Name() string { return "Qoder" }

// Priority returns the display order (lower = higher priority).
func (p *QoderProvider) Priority() int { return PriorityQoder }

// Initializers returns the file initializers for this provider.
func (p *QoderProvider) Initializers() []FileInitializer {
	proposalPath, applyPath := StandardCommandPaths(
		".qoder/commands",
		".md",
	)

	return []FileInitializer{
		NewInstructionFileInitializer("QODER.md"),
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
func (p *QoderProvider) IsConfigured(projectPath string) bool {
	return AreInitializersConfigured(p.Initializers(), projectPath)
}

// GetFilePaths returns the file paths managed by this provider.
func (p *QoderProvider) GetFilePaths() []string {
	return GetInitializerPaths(p.Initializers())
}
