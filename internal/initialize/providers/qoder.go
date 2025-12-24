// Package providers implements AI tool provider registration and
// initialization.
package providers

import "context"

func init() {
	if err := Register(Registration{
		ID:       "qoder",
		Name:     "Qoder",
		Priority: PriorityQoder,
		Provider: &QoderProvider{},
	}); err != nil {
		panic(err)
	}
}

// QoderProvider implements the Provider interface for Qoder.
type QoderProvider struct{}

func (*QoderProvider) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".qoder/commands/spectr"),
		NewConfigFileInitializer("QODER.md", RenderInstructionPointer),
		NewSlashCommandsInitializer(
			".qoder/commands/spectr",
			".md",
			FormatMarkdown,
		),
	}
}
