package providers

// init registers the Antigravity provider with the global registry.
func init() {
	Register(&AntigravityProvider{})
}

// AntigravityProvider implements the Provider interface for Antigravity.
// Antigravity uses AGENTS.md and .agent/workflows/ for slash commands.
type AntigravityProvider struct{}

// ID returns the unique identifier for the Antigravity provider.
func (*AntigravityProvider) ID() string { return "antigravity" }

// Name returns the display name for Antigravity.
func (*AntigravityProvider) Name() string { return "Antigravity" }

// Priority returns the display order for Antigravity.
func (*AntigravityProvider) Priority() int { return PriorityAntigravity }

// Initializers returns the file initializers for the Antigravity provider.
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
			FrontmatterProposal,
		),
		NewMarkdownSlashCommandInitializer(
			applyPath,
			"apply",
			FrontmatterApply,
		),
	}
}

func (p *AntigravityProvider) IsConfigured(projectPath string) bool {
	return AreInitializersConfigured(p.Initializers(), projectPath)
}

func (p *AntigravityProvider) GetFilePaths() []string {
	return GetInitializerPaths(p.Initializers())
}
