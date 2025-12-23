// Package providers implements AI tool provider registration and
// initialization.
package providers

import "context"

func init() {
	if err := Register(Registration{
		ID:       "continue",
		Name:     "Continue",
		Priority: PriorityContinue,
		Provider: &ContinueProvider{},
	}); err != nil {
		panic(err)
	}
}

// ContinueProvider implements the Provider interface for Continue.
type ContinueProvider struct{}

func (*ContinueProvider) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".continue/commands/spectr"),
		NewSlashCommandsInitializer(
			".continue/commands/spectr",
			".md",
			FormatMarkdown,
		),
	}
}
