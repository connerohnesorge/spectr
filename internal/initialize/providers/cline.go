package providers

// init registers the Cline provider with the global registry.
func init() {
	Register(&ClineProvider{})
}

// ClineProvider implements the Provider interface for Cline.
// Cline uses CLINE.md and .clinerules/commands/ for slash commands.
type ClineProvider struct{}

// ID returns the unique identifier for the Cline provider.
func (*ClineProvider) ID() string { return "cline" }

// Name returns the display name for Cline.
func (*ClineProvider) Name() string { return "Cline" }

// Priority returns the display order for Cline.
func (*ClineProvider) Priority() int { return PriorityCline }

// Initializers returns the file initializers for the Cline provider.
func (*ClineProvider) Initializers() []FileInitializer {
	proposalPath, applyPath := StandardCommandPaths(
		".clinerules/commands",
		".md",
	)

	return []FileInitializer{
		NewInstructionFileInitializer("CLINE.md"),
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

func (p *ClineProvider) IsConfigured(projectPath string) bool {
	return AreInitializersConfigured(p.Initializers(), projectPath)
}

func (p *ClineProvider) GetFilePaths() []string {
	return GetInitializerPaths(p.Initializers())
}
