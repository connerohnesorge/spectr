// Package providers implements AI tool provider registration and
// initialization.
package providers

import "context"

func init() {
	if err := Register(Registration{
		ID:       "costrict",
		Name:     "Costrict",
		Priority: PriorityCostrict,
		Provider: &CostrictProvider{},
	}); err != nil {
		panic(err)
	}
}

// CostrictProvider implements the Provider interface for Costrict.
type CostrictProvider struct{}

func (*CostrictProvider) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".costrict/commands/spectr"),
		NewConfigFileInitializer("COSTRICT.md", "instruction_pointer"),
		NewSlashCommandsInitializer(
			".costrict/commands/spectr",
			".md",
			FormatMarkdown,
		),
	}
}
