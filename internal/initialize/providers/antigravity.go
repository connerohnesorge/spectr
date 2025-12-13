//nolint:revive // unused-receiver acceptable for interface compliance
package providers

func init() {
	Register(&AntigravityProvider{})
}

// AntigravityProvider implements the Provider interface for Antigravity.
//
// Antigravity uses AGENTS.md and .agent/workflows/ for slash commands.
type AntigravityProvider struct{}

// ID returns the unique identifier for this provider.
func (*AntigravityProvider) ID() string { return "antigravity" }

// Name returns the human-readable name for display.
func (*AntigravityProvider) Name() string { return "Antigravity" }

// Priority returns the display order (lower = higher priority).
func (*AntigravityProvider) Priority() int { return PriorityAntigravity }

// Initializers returns the file initializers for this provider.
func (*AntigravityProvider) Initializers() []FileInitializer {
	proposalPath, applyPath := PrefixedCommandPaths(
		".agent/workflows",
		".md",
	)

	return []FileInitializer{
		NewInstructionFileInitializer("AGENTS.md"),
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
func (p *AntigravityProvider) IsConfigured(projectPath string) bool {
	return AreInitializersConfigured(p.Initializers(), projectPath)
}

// GetFilePaths returns the file paths managed by this provider.
func (p *AntigravityProvider) GetFilePaths() []string {
	return GetInitializerPaths(p.Initializers())
}
