//nolint:revive // unused-receiver acceptable for interface compliance
package providers

func init() {
	Register(&CursorProvider{})
}

// CursorProvider implements the Provider interface for Cursor.
// Cursor uses .cursorrules/commands/ for slash commands (no config file).
type CursorProvider struct{}

// ID returns the unique identifier for this provider.
func (p *CursorProvider) ID() string { return "cursor" }

// Name returns the human-readable name for display.
func (p *CursorProvider) Name() string { return "Cursor" }

// Priority returns the display order (lower = higher priority).
func (p *CursorProvider) Priority() int { return PriorityCursor }

// Initializers returns the file initializers for this provider.
func (p *CursorProvider) Initializers() []FileInitializer {
	proposalPath, applyPath := StandardCommandPaths(
		".cursorrules/commands",
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
func (p *CursorProvider) IsConfigured(projectPath string) bool {
	return AreInitializersConfigured(p.Initializers(), projectPath)
}

// GetFilePaths returns the file paths managed by this provider.
func (p *CursorProvider) GetFilePaths() []string {
	return GetInitializerPaths(p.Initializers())
}
