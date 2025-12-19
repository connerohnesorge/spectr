package providers

import "context"

func init() {
	err := Register(Registration{
		ID:       "gemini",
		Name:     "Gemini CLI",
		Priority: PriorityGemini,
		Provider: &GeminiProvider{},
	})
	if err != nil {
		panic(err)
	}
}

// GeminiProvider implements the Provider interface for Gemini CLI.
// Gemini uses .gemini/commands/spectr/ for TOML-based slash commands
// (no instruction file).
type GeminiProvider struct{}

// Initializers returns the initializers for Gemini CLI.
func (p *GeminiProvider) Initializers(ctx context.Context) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".gemini/commands/spectr"),
		NewSlashCommandsInitializer(
			".gemini/commands/spectr",
			".toml",
			FormatTOML,
		),
	}
}
