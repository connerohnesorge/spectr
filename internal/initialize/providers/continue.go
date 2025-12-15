package providers

// init registers the Continue provider with the global registry.
func init() {
	Register(&ContinueProvider{})
}

// ContinueProvider implements the Provider interface for Continue.
// Continue uses .continue/commands/ for slash commands.
// It does not use a separate instruction file.
type ContinueProvider struct{}

// ID returns the unique identifier for the Continue provider.
func (*ContinueProvider) ID() string { return "continue" }

// Name returns the display name for Continue.
func (*ContinueProvider) Name() string { return "Continue" }

// Priority returns the display order for Continue.
func (*ContinueProvider) Priority() int { return PriorityContinue }

// Initializers returns the file initializers for the Continue provider.
func (*ContinueProvider) Initializers() []FileInitializer {
	proposalPath, applyPath := StandardCommandPaths(
		".continue/commands",
		".md",
	)

	return []FileInitializer{
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

func (p *ContinueProvider) IsConfigured(projectPath string) bool {
	return AreInitializersConfigured(p.Initializers(), projectPath)
}

func (p *ContinueProvider) GetFilePaths() []string {
	return GetInitializerPaths(p.Initializers())
}
