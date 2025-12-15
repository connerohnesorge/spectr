package providers

// init registers the OpenCode provider with the global registry.
func init() {
	Register(&OpencodeProvider{})
}

// OpencodeProvider implements the Provider interface for OpenCode.
// OpenCode uses .opencode/command/spectr/ for slash commands.
// It does not use a separate instruction file.
type OpencodeProvider struct{}

// ID returns the unique identifier for the OpenCode provider.
func (*OpencodeProvider) ID() string { return "opencode" }

// Name returns the display name for OpenCode.
func (*OpencodeProvider) Name() string { return "OpenCode" }

// Priority returns the display order for OpenCode.
func (*OpencodeProvider) Priority() int { return PriorityOpencode }

// Initializers returns the file initializers for the OpenCode provider.
func (*OpencodeProvider) Initializers() []FileInitializer {
	proposalPath, applyPath := StandardCommandPaths(
		".opencode/command",
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

func (p *OpencodeProvider) IsConfigured(projectPath string) bool {
	return AreInitializersConfigured(p.Initializers(), projectPath)
}

func (p *OpencodeProvider) GetFilePaths() []string {
	return GetInitializerPaths(p.Initializers())
}
