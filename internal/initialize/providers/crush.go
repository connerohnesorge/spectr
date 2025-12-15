package providers

// init registers the Crush provider with the global registry.
func init() {
	Register(&CrushProvider{})
}

// CrushProvider implements the Provider interface for Crush.
// Crush uses CRUSH.md for instructions and .crush/commands/
// for slash commands.
type CrushProvider struct{}

// ID returns the unique identifier for the Crush provider.
func (*CrushProvider) ID() string { return "crush" }

// Name returns the display name for Crush.
func (*CrushProvider) Name() string { return "Crush" }

// Priority returns the display order for Crush.
func (*CrushProvider) Priority() int { return PriorityCrush }

// Initializers returns the file initializers for the Crush provider.
func (*CrushProvider) Initializers() []FileInitializer {
	proposalPath, applyPath := StandardCommandPaths(
		".crush/commands",
		".md",
	)

	return []FileInitializer{
		NewInstructionFileInitializer("CRUSH.md"),
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

func (p *CrushProvider) IsConfigured(projectPath string) bool {
	return AreInitializersConfigured(p.Initializers(), projectPath)
}

func (p *CrushProvider) GetFilePaths() []string {
	return GetInitializerPaths(p.Initializers())
}
