package providers

// init registers the Kilocode provider with the global registry.
func init() {
	Register(&KilocodeProvider{})
}

// KilocodeProvider implements the Provider interface for Kilocode.
// Kilocode uses .kilocode/commands/ for slash commands.
// It does not use a separate instruction file.
type KilocodeProvider struct{}

// ID returns the unique identifier for the Kilocode provider.
func (*KilocodeProvider) ID() string { return "kilocode" }

// Name returns the display name for Kilocode.
func (*KilocodeProvider) Name() string { return "Kilocode" }

// Priority returns the display order for Kilocode.
func (*KilocodeProvider) Priority() int { return PriorityKilocode }

// Initializers returns the file initializers for the Kilocode provider.
func (*KilocodeProvider) Initializers() []FileInitializer {
	proposalPath, applyPath := StandardCommandPaths(
		".kilocode/commands",
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

func (p *KilocodeProvider) IsConfigured(projectPath string) bool {
	return AreInitializersConfigured(p.Initializers(), projectPath)
}

func (p *KilocodeProvider) GetFilePaths() []string {
	return GetInitializerPaths(p.Initializers())
}
