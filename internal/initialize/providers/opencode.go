package providers

import "context"

func init() {
	if err := Register(Registration{
		ID:       "opencode",
		Name:     "OpenCode",
		Priority: PriorityOpencode,
		Provider: &OpencodeProvider{},
	}); err != nil {
		panic(err)
	}
}

// OpencodeProvider implements the Provider interface for OpenCode.
type OpencodeProvider struct{}

func (*OpencodeProvider) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".opencode/command/spectr"),
		// No config file for OpenCode
		NewSlashCommandsInitializer(
			".opencode/command/spectr",
			".md",
			FormatMarkdown,
		),
	}
}
