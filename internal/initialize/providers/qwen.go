package providers

import "context"

func init() {
	err := RegisterV2(Registration{
		ID:       "qwen",
		Name:     "Qwen Code",
		Priority: PriorityQwen,
		Provider: &QwenProvider{},
	})
	if err != nil {
		panic(err)
	}
}

// QwenProvider implements the ProviderV2 interface for Qwen Code.
// Qwen uses QWEN.md and .qwen/commands/spectr/ for slash commands.
type QwenProvider struct{}

// Initializers returns the initializers for Qwen Code.
func (p *QwenProvider) Initializers(ctx context.Context) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".qwen/commands/spectr"),
		NewConfigFileInitializer("QWEN.md"),
		NewSlashCommandsInitializerWithFrontmatter(
			".qwen/commands/spectr",
			".md",
			FormatMarkdown,
			map[string]string{
				"proposal": FrontmatterProposal,
				"apply":    FrontmatterApply,
			},
		),
	}
}
