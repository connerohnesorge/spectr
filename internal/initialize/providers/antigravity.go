// Package providers implements AI tool provider registration and
// initialization.
package providers

import "context"

func init() {
	if err := Register(Registration{
		ID:       "antigravity",
		Name:     "Antigravity",
		Priority: PriorityAntigravity,
		Provider: &AntigravityProvider{},
	}); err != nil {
		panic(err)
	}
}

// AntigravityProvider implements the Provider interface for Antigravity.
type AntigravityProvider struct{}

func (*AntigravityProvider) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".antigravity/commands/spectr"),
		NewConfigFileInitializer(
			"ANTIGRAVITY.md",
			"instruction_pointer",
		),
		NewSlashCommandsInitializer(
			".antigravity/commands/spectr",
			".md",
			FormatMarkdown,
		),
	}
}
