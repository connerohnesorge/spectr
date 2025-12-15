package providers

// init registers the CoStrict provider with the global registry.
func init() {
	Register(&CostrictProvider{})
}

// CostrictProvider implements the Provider interface for CoStrict.
// CoStrict uses COSTRICT.md and .costrict/commands/ for slash commands.
type CostrictProvider struct{}

// ID returns the unique identifier for the CoStrict provider.
func (*CostrictProvider) ID() string { return "costrict" }

// Name returns the display name for CoStrict.
func (*CostrictProvider) Name() string { return "CoStrict" }

// Priority returns the display order for CoStrict.
func (*CostrictProvider) Priority() int { return PriorityCostrict }

// Initializers returns the file initializers for the CoStrict provider.
func (*CostrictProvider) Initializers() []FileInitializer {
	proposalPath, applyPath := StandardCommandPaths(
		".costrict/commands",
		".md",
	)

	return []FileInitializer{
		NewInstructionFileInitializer("COSTRICT.md"),
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

func (p *CostrictProvider) IsConfigured(projectPath string) bool {
	return AreInitializersConfigured(p.Initializers(), projectPath)
}

func (p *CostrictProvider) GetFilePaths() []string {
	return GetInitializerPaths(p.Initializers())
}
