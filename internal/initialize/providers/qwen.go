// Package providers implements AI tool provider registration and
// initialization.
package providers

import "context"

func init() {
	if err := Register(Registration{
		ID:       "qwen",
		Name:     "Qwen",
		Priority: PriorityQwen,
		Provider: &QwenProvider{},
	}); err != nil {
		panic(err)
	}
}

// QwenProvider implements the Provider interface for Qwen.
type QwenProvider struct{}

func (*QwenProvider) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".qwen/commands/spectr"),
		NewConfigFileInitializer("QWEN.md", RenderInstructionPointer),
		NewSlashCommandsInitializer(
			".qwen/commands/spectr",
			".md",
			FormatMarkdown,
		),
	}
}
