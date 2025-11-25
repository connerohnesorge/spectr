package providers

func init() {
	Register(NewClaudeProvider())
}

// ClaudeProvider implements the Provider interface for Claude Code.
// Claude Code uses CLAUDE.md and .claude/commands/ for slash commands.
type ClaudeProvider struct {
	BaseProvider
}

// NewClaudeProvider creates a new Claude Code provider.
func NewClaudeProvider() *ClaudeProvider {
	return &ClaudeProvider{
		BaseProvider: BaseProvider{
			id:            "claude-code",
			name:          "Claude Code",
			priority:      PriorityClaudeCode,
			configFile:    "CLAUDE.md",
			slashDir:      ".claude/commands",
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}
