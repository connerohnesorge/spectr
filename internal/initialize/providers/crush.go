// Package providers implements AI tool provider registration and
// initialization.
package providers

import "context"

func init() {
	if err := Register(Registration{
		ID:       "crush",
		Name:     "Crush",
		Priority: PriorityCrush,
		Provider: &CrushProvider{},
	}); err != nil {
		panic(err)
	}
}

// CrushProvider implements the Provider interface for Crush.
type CrushProvider struct{}

func (*CrushProvider) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".crush/commands/spectr"),
		NewConfigFileInitializer("CRUSH.md", "instruction_pointer"),
		NewSlashCommandsInitializer(
			".crush/commands/spectr",
			".md",
			FormatMarkdown,
		),
	}
}
