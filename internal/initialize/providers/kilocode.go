//nolint:revive // unused-receiver acceptable for interface compliance
package providers

func init() {
	Register(&KilocodeProvider{})
}

// KilocodeProvider implements the Provider interface for Kilocode.
// Kilocode uses .kilocode/commands/ for slash commands (no config file).
type KilocodeProvider struct{}

// ID returns the unique identifier for this provider.
func (p *KilocodeProvider) ID() string { return "kilocode" }

// Name returns the human-readable name for display.
func (p *KilocodeProvider) Name() string { return "Kilocode" }

// Priority returns the display order (lower = higher priority).
func (p *KilocodeProvider) Priority() int { return PriorityKilocode }

// Initializers returns the file initializers for this provider.
func (p *KilocodeProvider) Initializers() []FileInitializer {
	proposalPath, applyPath := StandardCommandPaths(
		".kilocode/commands",
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
func (p *KilocodeProvider) IsConfigured(projectPath string) bool {
	return AreInitializersConfigured(p.Initializers(), projectPath)
}

// GetFilePaths returns the file paths managed by this provider.
func (p *KilocodeProvider) GetFilePaths() []string {
	return GetInitializerPaths(p.Initializers())
}
