package providers

// init registers the Qoder provider with the global registry.
func init() {
	Register(&QoderProvider{})
}

// QoderProvider implements the Provider interface for Qoder.
// Qoder uses QODER.md and .qoder/commands/ for slash commands.
type QoderProvider struct{}

// ID returns the unique identifier for the Qoder provider.
func (*QoderProvider) ID() string { return "qoder" }

// Name returns the display name for Qoder.
func (*QoderProvider) Name() string { return "Qoder" }

// Priority returns the display order for Qoder.
func (*QoderProvider) Priority() int { return PriorityQoder }

// Initializers returns the file initializers for the Qoder provider.
func (*QoderProvider) Initializers() []FileInitializer {
	proposalPath, applyPath := StandardCommandPaths(
		".qoder/commands",
		".md",
	)

	return []FileInitializer{
		NewInstructionFileInitializer("QODER.md"),
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

func (p *QoderProvider) IsConfigured(projectPath string) bool {
	return AreInitializersConfigured(p.Initializers(), projectPath)
}

func (p *QoderProvider) GetFilePaths() []string {
	return GetInitializerPaths(p.Initializers())
}
