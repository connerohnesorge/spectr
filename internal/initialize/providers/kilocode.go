// Package providers implements AI tool provider registration and
// initialization.
package providers

import "context"

func init() {
	if err := Register(Registration{
		ID:       "kilocode",
		Name:     "Kilocode",
		Priority: PriorityKilocode,
		Provider: &KilocodeProvider{},
	}); err != nil {
		panic(err)
	}
}

// KilocodeProvider implements the Provider interface for Kilocode.
type KilocodeProvider struct{}

func (*KilocodeProvider) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".kilocode/commands/spectr"),
		NewSlashCommandsInitializer(
			".kilocode/commands/spectr",
			".md",
			FormatMarkdown,
		),
	}
}
