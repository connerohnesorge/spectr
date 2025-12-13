//nolint:revive // unused-receiver acceptable for interface compliance
package providers

func init() {
	Register(&ClineProvider{})
}

// ClineProvider implements the Provider interface for Cline.
// Cline uses CLINE.md and .clinerules/commands/ for slash commands.
type ClineProvider struct{}

// ID returns the unique identifier for this provider.
func (*ClineProvider) ID() string { return "cline" }

// Name returns the human-readable name for display.
func (*ClineProvider) Name() string { return "Cline" }

// Priority returns the display order (lower = higher priority).
func (*ClineProvider) Priority() int { return PriorityCline }

// Initializers returns the file initializers for this provider.
func (*ClineProvider) Initializers() []FileInitializer {
	proposalPath, applyPath := StandardCommandPaths(
		".clinerules/commands",
		".md",
	)

	return []FileInitializer{
		NewInstructionFileInitializer("CLINE.md"),
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
func (p *ClineProvider) IsConfigured(projectPath string) bool {
	return AreInitializersConfigured(p.Initializers(), projectPath)
}

// GetFilePaths returns the file paths managed by this provider.
func (p *ClineProvider) GetFilePaths() []string {
	return GetInitializerPaths(p.Initializers())
}
