package providers

import "context"

func init() {
	err := RegisterV2(Registration{
		ID:       "aider",
		Name:     "Aider",
		Priority: PriorityAider,
		Provider: &AiderProvider{},
	})
	if err != nil {
		panic(err)
	}
}

// AiderProvider implements the ProviderV2 interface for Aider.
// Aider uses .aider/commands/spectr/ for slash commands (no config file).
type AiderProvider struct{}

// Initializers returns the initializers for Aider.
func (p *AiderProvider) Initializers(ctx context.Context) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".aider/commands/spectr"),
		NewSlashCommandsInitializerWithFrontmatter(
			".aider/commands/spectr",
			".md",
			FormatMarkdown,
			map[string]string{
				"proposal": FrontmatterProposal,
				"apply":    FrontmatterApply,
			},
		),
	}
}
