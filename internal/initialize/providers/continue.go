package providers

import "context"

func init() {
	err := Register(Registration{
		ID:       "continue",
		Name:     "Continue",
		Priority: PriorityContinue,
		Provider: &ContinueProvider{},
	})
	if err != nil {
		panic(err)
	}
}

// ContinueProvider implements the Provider interface for Continue.
// Continue uses .continue/commands/spectr/ for slash commands (no config file).
type ContinueProvider struct{}

// Initializers returns the initializers for Continue.
func (p *ContinueProvider) Initializers(ctx context.Context) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".continue/commands/spectr"),
		NewSlashCommandsInitializerWithFrontmatter(
			".continue/commands/spectr",
			".md",
			FormatMarkdown,
			map[string]string{
				"proposal": FrontmatterProposal,
				"apply":    FrontmatterApply,
			},
		),
	}
}
