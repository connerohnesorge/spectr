// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file implements the Gemini CLI provider using the Provider interface.
// Gemini uses .gemini/commands/ for TOML-based slash commands.
package providers

import (
	"context"
)

func init() {
	// Register with Registry
	err := Register(Registration{
		ID:       "gemini",
		Name:     "Gemini CLI",
		Priority: PriorityGemini,
		Provider: &GeminiProvider{},
	})
	if err != nil {
		// Panic on registration failure since this is called at init time
		// and indicates a programming error (e.g., duplicate ID)
		panic("failed to register gemini provider: " + err.Error())
	}
}

// GeminiProvider implements the Provider interface for Gemini CLI.
//
// Gemini CLI uses:
//   - NO instruction file (configFile is empty)
//   - .gemini/commands/spectr/ for slash commands
//   - TOML format for slash commands
type GeminiProvider struct{}

// Initializers returns the list of Initializers needed to configure
// Gemini CLI for use with spectr.
//
// Returns:
//   - DirectoryInitializer for .gemini/commands/spectr
//   - SlashCommandsInitializer for .gemini/commands/spectr with TOML format
//
// Note: Gemini has no instruction file, so no ConfigFileInitializer.
func (*GeminiProvider) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		// Create the slash commands directory
		NewDirectoryInitializer(false, ".gemini/commands/spectr"),

		// Create/update slash commands (proposal.toml, apply.toml)
		// Gemini uses TOML format with no frontmatter
		NewSlashCommandsInitializer(
			".gemini/commands/spectr",
			".toml",
			FormatTOML,
			nil, // No frontmatter for TOML
			false,
		),
	}
}

// Ensure GeminiProvider implements the Provider interface.
var _ Provider = (*GeminiProvider)(nil)
