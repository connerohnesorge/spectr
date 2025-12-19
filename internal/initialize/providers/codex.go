package providers

import "context"

func init() {
	err := Register(Registration{
		ID:       "codex",
		Name:     "Codex CLI",
		Priority: PriorityCodex,
		Provider: &CodexProvider{},
	})
	if err != nil {
		panic(err)
	}
}

// CodexProvider implements the Provider interface for Codex CLI.
// Codex uses AGENTS.md and global ~/.codex/prompts/ for commands.
type CodexProvider struct{}

// Initializers returns the initializers for Codex CLI.
func (p *CodexProvider) Initializers(ctx context.Context) []Initializer {
	return []Initializer{
		NewGlobalDirectoryInitializer(".codex/prompts"),
		NewConfigFileInitializer("AGENTS.md"),
		NewGlobalSlashCommandsInitializerWithFrontmatter(
			".codex/prompts",
			".md",
			FormatMarkdown,
			map[string]string{
				"proposal": FrontmatterProposal,
				"apply":    FrontmatterApply,
			},
		),
	}
}
