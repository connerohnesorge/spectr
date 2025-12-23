// Package providers implements AI tool provider registration and
// initialization.
package providers

import "context"

func init() {
	if err := Register(Registration{
		ID:       "windsurf",
		Name:     "Windsurf",
		Priority: PriorityWindsurf,
		Provider: &WindsurfProvider{},
	}); err != nil {
		panic(err)
	}
}

// WindsurfProvider implements the Provider interface for Windsurf.
type WindsurfProvider struct{}

func (*WindsurfProvider) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".windsurf/commands/spectr"),
		NewSlashCommandsInitializer(
			".windsurf/commands/spectr",
			".md",
			FormatMarkdown,
		),
	}
}
