//nolint:revive // unused-receiver acceptable for interface compliance
package providers

func init() {
	Register(&CrushProvider{})
}

// CrushProvider implements the Provider interface for Crush.
// Crush uses CRUSH.md for instructions and .crush/commands/ for slash commands.
type CrushProvider struct{}

// ID returns the unique identifier for this provider.
func (p *CrushProvider) ID() string { return "crush" }

// Name returns the human-readable name for display.
func (p *CrushProvider) Name() string { return "Crush" }

// Priority returns the display order (lower = higher priority).
func (p *CrushProvider) Priority() int { return PriorityCrush }

// Initializers returns the file initializers for this provider.
func (p *CrushProvider) Initializers() []FileInitializer {
	proposalPath, applyPath := StandardCommandPaths(
		".crush/commands",
		".md",
	)

	return []FileInitializer{
		NewInstructionFileInitializer("CRUSH.md"),
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
func (p *CrushProvider) IsConfigured(projectPath string) bool {
	return AreInitializersConfigured(p.Initializers(), projectPath)
}

// GetFilePaths returns the file paths managed by this provider.
func (p *CrushProvider) GetFilePaths() []string {
	return GetInitializerPaths(p.Initializers())
}
