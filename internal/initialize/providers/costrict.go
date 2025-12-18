package providers

import "context"

func init() {
	err := RegisterV2(Registration{
		ID:       "costrict",
		Name:     "CoStrict",
		Priority: PriorityCostrict,
		Provider: &CostrictProvider{},
	})
	if err != nil {
		panic(err)
	}
}

// CostrictProvider implements the ProviderV2 interface for CoStrict.
// CoStrict uses COSTRICT.md and .costrict/commands/spectr/ for slash commands.
type CostrictProvider struct{}

// Initializers returns the initializers for CoStrict.
func (p *CostrictProvider) Initializers(ctx context.Context) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".costrict/commands/spectr"),
		NewConfigFileInitializer("COSTRICT.md"),
		NewSlashCommandsInitializerWithFrontmatter(
			".costrict/commands/spectr",
			".md",
			FormatMarkdown,
			map[string]string{
				"proposal": FrontmatterProposal,
				"apply":    FrontmatterApply,
			},
		),
	}
}
