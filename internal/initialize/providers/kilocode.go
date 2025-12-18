package providers

import "context"

func init() {
	err := RegisterV2(Registration{
		ID:       "kilocode",
		Name:     "Kilocode",
		Priority: PriorityKilocode,
		Provider: &KilocodeProvider{},
	})
	if err != nil {
		panic(err)
	}
}

// KilocodeProvider implements the ProviderV2 interface for Kilocode.
// Kilocode uses .kilocode/commands/spectr/ for slash commands (no config file).
type KilocodeProvider struct{}

// Initializers returns the initializers for Kilocode.
func (p *KilocodeProvider) Initializers(ctx context.Context) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".kilocode/commands/spectr"),
		NewSlashCommandsInitializerWithFrontmatter(
			".kilocode/commands/spectr",
			".md",
			FormatMarkdown,
			map[string]string{
				"proposal": FrontmatterProposal,
				"apply":    FrontmatterApply,
			},
		),
	}
}
