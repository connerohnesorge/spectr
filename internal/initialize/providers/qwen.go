package providers

// init registers the Qwen Code provider with the global registry.
func init() {
	Register(&QwenProvider{})
}

// QwenProvider implements the Provider interface for Qwen Code.
// Qwen uses QWEN.md and .qwen/commands/ for slash commands.
type QwenProvider struct{}

// ID returns the unique identifier for the Qwen Code provider.
func (*QwenProvider) ID() string { return "qwen" }

// Name returns the display name for Qwen Code.
func (*QwenProvider) Name() string { return "Qwen Code" }

// Priority returns the display order for Qwen Code.
func (*QwenProvider) Priority() int { return PriorityQwen }

// Initializers returns the file initializers for the Qwen Code provider.
func (*QwenProvider) Initializers() []FileInitializer {
	proposalPath, applyPath := StandardCommandPaths(
		".qwen/commands",
		".md",
	)

	return []FileInitializer{
		NewInstructionFileInitializer("QWEN.md"),
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

func (p *QwenProvider) IsConfigured(projectPath string) bool {
	return AreInitializersConfigured(p.Initializers(), projectPath)
}

func (p *QwenProvider) GetFilePaths() []string {
	return GetInitializerPaths(p.Initializers())
}
