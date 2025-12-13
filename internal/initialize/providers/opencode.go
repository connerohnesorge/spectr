//nolint:revive // unused-receiver acceptable for interface compliance
package providers

func init() {
	Register(&OpencodeProvider{})
}

// OpencodeProvider implements the Provider interface for OpenCode.
// OpenCode uses .opencode/command/spectr/ for slash commands.
// It has no instruction file as it uses JSON configuration.
type OpencodeProvider struct{}

// ID returns the unique identifier for this provider.
func (*OpencodeProvider) ID() string { return "opencode" }

// Name returns the human-readable name for display.
func (*OpencodeProvider) Name() string { return "OpenCode" }

// Priority returns the display order (lower = higher priority).
func (*OpencodeProvider) Priority() int { return PriorityOpencode }

// Initializers returns the file initializers for this provider.
func (p *OpencodeProvider) Initializers() []FileInitializer {
	proposalPath, applyPath := StandardCommandPaths(
		".opencode/command",
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
func (p *OpencodeProvider) IsConfigured(projectPath string) bool {
	return AreInitializersConfigured(p.Initializers(), projectPath)
}

// GetFilePaths returns the file paths managed by this provider.
func (p *OpencodeProvider) GetFilePaths() []string {
	return GetInitializerPaths(p.Initializers())
}
