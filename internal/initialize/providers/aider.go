package providers

// init registers the Aider provider with the global registry.
func init() {
	Register(&AiderProvider{})
}

// AiderProvider implements the Provider interface for Aider.
// Aider uses .aider/commands/ for slash commands (no instruction file).
type AiderProvider struct{}

// ID returns the unique identifier for the Aider provider.
func (*AiderProvider) ID() string { return "aider" }

// Name returns the display name for Aider.
func (*AiderProvider) Name() string { return "Aider" }

// Priority returns the display order for Aider.
func (*AiderProvider) Priority() int { return PriorityAider }

// Initializers returns the file initializers for the Aider provider.
func (*AiderProvider) Initializers() []FileInitializer {
	proposalPath, applyPath := StandardCommandPaths(
		".aider/commands",
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

func (p *AiderProvider) IsConfigured(projectPath string) bool {
	return AreInitializersConfigured(p.Initializers(), projectPath)
}

func (p *AiderProvider) GetFilePaths() []string {
	return GetInitializerPaths(p.Initializers())
}
