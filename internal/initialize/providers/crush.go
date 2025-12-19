package providers

import "context"

func init() {
	err := Register(Registration{
		ID:       "crush",
		Name:     "Crush",
		Priority: PriorityCrush,
		Provider: &CrushProvider{},
	})
	if err != nil {
		panic(err)
	}
}

// CrushProvider implements the Provider interface for Crush.
// Crush uses CRUSH.md for instructions and .crush/commands/spectr/ for slash commands.
type CrushProvider struct{}

// Initializers returns the initializers for Crush.
func (p *CrushProvider) Initializers(ctx context.Context) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".crush/commands/spectr"),
		NewConfigFileInitializer("CRUSH.md"),
		NewSlashCommandsInitializerWithFrontmatter(
			".crush/commands/spectr",
			".md",
			FormatMarkdown,
			map[string]string{
				"proposal": FrontmatterProposal,
				"apply":    FrontmatterApply,
			},
		),
	}
}
