//nolint:revive // unused-receiver acceptable for interface compliance
package providers

func init() {
	Register(&CostrictProvider{})
}

// CostrictProvider implements the Provider interface for CoStrict.
// CoStrict uses COSTRICT.md and .costrict/commands/ for slash commands.
type CostrictProvider struct{}

// ID returns the unique identifier for this provider.
func (p *CostrictProvider) ID() string { return "costrict" }

// Name returns the human-readable name for display.
func (p *CostrictProvider) Name() string { return "CoStrict" }

// Priority returns the display order (lower = higher priority).
func (p *CostrictProvider) Priority() int { return PriorityCostrict }

// Initializers returns the file initializers for this provider.
func (p *CostrictProvider) Initializers() []FileInitializer {
	proposalPath, applyPath := StandardCommandPaths(
		".costrict/commands",
		".md",
	)

	return []FileInitializer{
		NewInstructionFileInitializer("COSTRICT.md"),
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
func (p *CostrictProvider) IsConfigured(projectPath string) bool {
	return AreInitializersConfigured(p.Initializers(), projectPath)
}

// GetFilePaths returns the file paths managed by this provider.
func (p *CostrictProvider) GetFilePaths() []string {
	return GetInitializerPaths(p.Initializers())
}
