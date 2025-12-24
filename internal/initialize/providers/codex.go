// Package providers implements AI tool provider registration and
// initialization.
package providers

import "context"

func init() {
	if err := Register(Registration{
		ID:       "codex",
		Name:     "Codex",
		Priority: PriorityCodex,
		Provider: &CodexProvider{},
	}); err != nil {
		panic(err)
	}
}

// CodexProvider implements the Provider interface for Codex.
type CodexProvider struct{}

func (*CodexProvider) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".codex/commands/spectr"),
		NewConfigFileInitializer("AGENTS.md", RenderInstructionPointer),
		NewSlashCommandsInitializer(
			".codex/commands/spectr",
			".md",
			FormatMarkdown,
		),
	}
}
