package providers

// ClaudeProvider implements the Provider interface for Claude Code.
// Claude Code uses CLAUDE.md and .claude/commands/ for slash commands.
type ClaudeProvider struct {
	BaseProvider
}

// NewClaudeProvider creates a new Claude Code provider.
func NewClaudeProvider() *ClaudeProvider {
	proposalPath, applyPath := StandardCommandPaths(
		".claude/commands",
		".md",
	)

	return &ClaudeProvider{
		BaseProvider: BaseProvider{
			id:            "claude-code",
			name:          "Claude Code",
			priority:      PriorityClaudeCode,
			configFile:    "CLAUDE.md",
			proposalPath:  proposalPath,
			applyPath:     applyPath,
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}
