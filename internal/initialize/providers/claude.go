package providers

func init() {
	Register(NewClaudeProvider())
}

// ClaudeProvider implements the Provider interface for Claude Code.
// Claude Code uses CLAUDE.md and .claude/commands/ for slash commands.
type ClaudeProvider struct {
	BaseProvider
}

// NewClaudeProvider constructs a Claude Code provider configured with standard
// command paths, Markdown command format, and default frontmatter.
//
// The returned *ClaudeProvider has id "claude-code", name "Claude Code",
// priority PriorityClaudeCode, config file "CLAUDE.md", proposal and apply
// paths from StandardCommandPaths(".claude/commands", ".md"), commandFormat
// FormatMarkdown, and frontmatter from StandardFrontmatter().
func NewClaudeProvider() *ClaudeProvider {
	proposalPath, applyPath := StandardCommandPaths(
		".claude/commands", ".md",
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