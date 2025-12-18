package providers

import "context"

func init() {
	err := RegisterV2(Registration{
		ID:       "windsurf",
		Name:     "Windsurf",
		Priority: PriorityWindsurf,
		Provider: &WindsurfProvider{},
	})
	if err != nil {
		panic(err)
	}
}

// WindsurfProvider implements the ProviderV2 interface for Windsurf.
// Windsurf uses .windsurf/commands/spectr/ for slash commands (no config file).
type WindsurfProvider struct{}

// Initializers returns the initializers for Windsurf.
func (p *WindsurfProvider) Initializers(ctx context.Context) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".windsurf/commands/spectr"),
		NewSlashCommandsInitializerWithFrontmatter(
			".windsurf/commands/spectr",
			".md",
			FormatMarkdown,
			map[string]string{
				"proposal": FrontmatterProposal,
				"apply":    FrontmatterApply,
			},
		),
	}
}
