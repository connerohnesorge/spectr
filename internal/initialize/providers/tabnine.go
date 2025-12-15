package providers

// init registers the Tabnine provider with the global registry.
func init() {
	Register(&TabnineProvider{})
}

// TabnineProvider implements the Provider interface for Tabnine.
// Tabnine uses .tabnine/commands/ for slash commands.
// It does not use a separate instruction file.
type TabnineProvider struct{}

// ID returns the unique identifier for the Tabnine provider.
func (*TabnineProvider) ID() string { return "tabnine" }

// Name returns the display name for Tabnine.
func (*TabnineProvider) Name() string { return "Tabnine" }

// Priority returns the display order for Tabnine.
func (*TabnineProvider) Priority() int { return PriorityTabnine }

// Initializers returns the file initializers for the Tabnine provider.
func (*TabnineProvider) Initializers() []FileInitializer {
	proposalPath, applyPath := StandardCommandPaths(
		".tabnine/commands",
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

func (p *TabnineProvider) IsConfigured(projectPath string) bool {
	return AreInitializersConfigured(p.Initializers(), projectPath)
}

func (p *TabnineProvider) GetFilePaths() []string {
	return GetInitializerPaths(p.Initializers())
}
