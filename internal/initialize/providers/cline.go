// Package providers implements AI tool provider registration and
// initialization.
package providers

import "context"

func init() {
	if err := Register(Registration{
		ID:       "cline",
		Name:     "Cline",
		Priority: PriorityCline,
		Provider: &ClineProvider{},
	}); err != nil {
		panic(err)
	}
}

// ClineProvider implements the Provider interface for Cline.
type ClineProvider struct{}

func (*ClineProvider) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".clinerules/commands/spectr"),
		NewConfigFileInitializer("CLINE.md", RenderInstructionPointer),
		NewSlashCommandsInitializer(
			".clinerules/commands/spectr",
			".md",
			FormatMarkdown,
		),
	}
}
