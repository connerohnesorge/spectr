package providers

// init registers the Claude Code provider with the global registry.
func init() {
	Register(&ClaudeProvider{})
}

// ClaudeProvider implements the Provider interface for Claude Code.
// Claude Code uses CLAUDE.md and .claude/commands/ for slash commands.
type ClaudeProvider struct{}

// ID returns the unique identifier for the Claude Code provider.
func (*ClaudeProvider) ID() string { return "claude-code" }

// Name returns the display name for Claude Code.
func (*ClaudeProvider) Name() string { return "Claude Code" }

// Priority returns the display order for Claude Code.
func (*ClaudeProvider) Priority() int { return PriorityClaudeCode }

// Initializers returns the file initializers for the Claude Code provider.
func (*ClaudeProvider) Initializers() []FileInitializer {
	proposalPath, applyPath := StandardCommandPaths(
		".claude/commands",
		".md",
	)

	return []FileInitializer{
		NewInstructionFileInitializer("CLAUDE.md"),
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

func (p *ClaudeProvider) IsConfigured(projectPath string) bool {
	return AreInitializersConfigured(p.Initializers(), projectPath)
}

func (p *ClaudeProvider) GetFilePaths() []string {
	return GetInitializerPaths(p.Initializers())
}
