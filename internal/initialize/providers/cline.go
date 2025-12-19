package providers

import "context"

func init() {
	err := Register(Registration{
		ID:       "cline",
		Name:     "Cline",
		Priority: PriorityCline,
		Provider: &ClineProvider{},
	})
	if err != nil {
		panic(err)
	}
}

// ClineProvider implements the Provider interface for Cline.
// Cline uses CLINE.md and .clinerules/commands/spectr/ for slash commands.
type ClineProvider struct{}

// Initializers returns the initializers for Cline.
func (p *ClineProvider) Initializers(ctx context.Context) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".clinerules/commands/spectr"),
		NewConfigFileInitializer("CLINE.md"),
		NewSlashCommandsInitializerWithFrontmatter(
			".clinerules/commands/spectr",
			".md",
			FormatMarkdown,
			map[string]string{
				"proposal": FrontmatterProposal,
				"apply":    FrontmatterApply,
			},
		),
	}
}
