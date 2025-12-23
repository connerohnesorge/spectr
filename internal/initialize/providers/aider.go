// Package providers implements AI tool provider registration and
// initialization.
package providers

import "context"

func init() {
	if err := Register(Registration{
		ID:       "aider",
		Name:     "Aider",
		Priority: PriorityAider,
		Provider: &AiderProvider{},
	}); err != nil {
		panic(err)
	}
}

// AiderProvider implements the Provider interface for Aider.
type AiderProvider struct{}

func (*AiderProvider) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".aider/prompts/spectr"),
		NewSlashCommandsInitializer(
			".aider/prompts/spectr",
			".md",
			FormatMarkdown,
		),
	}
}
