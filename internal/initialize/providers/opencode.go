package providers

import "context"

func init() {
	err := RegisterV2(Registration{
		ID:       "opencode",
		Name:     "OpenCode",
		Priority: PriorityOpencode,
		Provider: &OpencodeProvider{},
	})
	if err != nil {
		panic(err)
	}
}

// OpencodeProvider implements the ProviderV2 interface for OpenCode.
// OpenCode uses .opencode/command/spectr/ for slash commands.
// It has no instruction file as it uses JSON configuration.
type OpencodeProvider struct{}

// Initializers returns the initializers for OpenCode.
func (p *OpencodeProvider) Initializers(ctx context.Context) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".opencode/command/spectr"),
		NewSlashCommandsInitializerWithFrontmatter(
			".opencode/command/spectr",
			".md",
			FormatMarkdown,
			map[string]string{
				"proposal": FrontmatterProposal,
				"apply":    FrontmatterApply,
			},
		),
	}
}
