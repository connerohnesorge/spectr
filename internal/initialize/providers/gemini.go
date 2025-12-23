package providers

import (
	"context"
)

func init() {
	if err := Register(Registration{
		ID:       "gemini",
		Name:     "Gemini CLI",
		Priority: PriorityGemini,
		Provider: &GeminiProvider{},
	}); err != nil {
		panic(err)
	}
}

// GeminiProvider implements the Provider interface for Gemini CLI.
type GeminiProvider struct{}

func (*GeminiProvider) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".gemini/commands/spectr"),
		// No config file for Gemini
		NewSlashCommandsInitializer(
			".gemini/commands/spectr",
			".toml",
			FormatTOML,
		),
	}
}
