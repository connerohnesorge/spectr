//nolint:revive // unused-receiver acceptable for interface compliance
package providers

func init() {
	Register(&QwenProvider{})
}

// QwenProvider implements the Provider interface for Qwen Code.
// Qwen uses QWEN.md and .qwen/commands/ for slash commands.
type QwenProvider struct{}

// ID returns the unique identifier for this provider.
func (*QwenProvider) ID() string { return "qwen" }

// Name returns the human-readable name for display.
func (*QwenProvider) Name() string { return "Qwen Code" }

// Priority returns the display order (lower = higher priority).
func (*QwenProvider) Priority() int { return PriorityQwen }

// Initializers returns the file initializers for this provider.
func (p *QwenProvider) Initializers() []FileInitializer {
	proposalPath, applyPath := StandardCommandPaths(
		".qwen/commands",
		".md",
	)

	return []FileInitializer{
		NewInstructionFileInitializer("QWEN.md"),
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
func (p *QwenProvider) IsConfigured(projectPath string) bool {
	return AreInitializersConfigured(p.Initializers(), projectPath)
}

// GetFilePaths returns the file paths managed by this provider.
func (p *QwenProvider) GetFilePaths() []string {
	return GetInitializerPaths(p.Initializers())
}
