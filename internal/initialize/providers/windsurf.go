package providers

// init registers the Windsurf provider with the global registry.
func init() {
	Register(&WindsurfProvider{})
}

// WindsurfProvider implements the Provider interface for Windsurf.
// Windsurf uses .windsurf/commands/ for slash commands.
// It does not use a separate instruction file.
type WindsurfProvider struct{}

// ID returns the unique identifier for the Windsurf provider.
func (*WindsurfProvider) ID() string { return "windsurf" }

// Name returns the display name for Windsurf.
func (*WindsurfProvider) Name() string { return "Windsurf" }

// Priority returns the display order for Windsurf.
func (*WindsurfProvider) Priority() int { return PriorityWindsurf }

// Initializers returns the file initializers for the Windsurf provider.
func (*WindsurfProvider) Initializers() []FileInitializer {
	proposalPath, applyPath := StandardCommandPaths(
		".windsurf/commands",
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

func (p *WindsurfProvider) IsConfigured(projectPath string) bool {
	return AreInitializersConfigured(p.Initializers(), projectPath)
}

func (p *WindsurfProvider) GetFilePaths() []string {
	return GetInitializerPaths(p.Initializers())
}
