package providers

import "context"

func init() {
	err := RegisterV2(Registration{
		ID:       "claude-code",
		Name:     "Claude Code",
		Priority: PriorityClaudeCode,
		Provider: &ClaudeProvider{},
	})
	if err != nil {
		panic(err)
	}
}

// ClaudeProvider implements the ProviderV2 interface for Claude Code.
// Claude Code uses CLAUDE.md and .claude/commands/spectr/ for slash commands.
type ClaudeProvider struct{}

// Initializers returns the initializers for Claude Code.
func (p *ClaudeProvider) Initializers(ctx context.Context) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".claude/commands/spectr"),
		NewConfigFileInitializer("CLAUDE.md"),
		NewSlashCommandsInitializerWithFrontmatter(
			".claude/commands/spectr",
			".md",
			FormatMarkdown,
			map[string]string{
				"proposal": FrontmatterProposal,
				"apply":    FrontmatterApply,
			},
		),
	}
}
