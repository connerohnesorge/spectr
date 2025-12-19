package providers

import "context"

func init() {
	err := Register(Registration{
		ID:       "qoder",
		Name:     "Qoder",
		Priority: PriorityQoder,
		Provider: &QoderProvider{},
	})
	if err != nil {
		panic(err)
	}
}

// QoderProvider implements the Provider interface for Qoder.
// Qoder uses QODER.md and .qoder/commands/spectr/ for slash commands.
type QoderProvider struct{}

// Initializers returns the initializers for Qoder.
func (p *QoderProvider) Initializers(ctx context.Context) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".qoder/commands/spectr"),
		NewConfigFileInitializer("QODER.md"),
		NewSlashCommandsInitializerWithFrontmatter(
			".qoder/commands/spectr",
			".md",
			FormatMarkdown,
			map[string]string{
				"proposal": FrontmatterProposal,
				"apply":    FrontmatterApply,
			},
		),
	}
}
