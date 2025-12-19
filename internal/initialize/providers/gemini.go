// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file implements the Gemini CLI provider using the ProviderV2 interface.
// Gemini uses .gemini/commands/ for TOML-based slash commands.
package providers

import (
	"context"
)

func init() {
	// Register with RegistryV2
	err := RegisterV2(Registration{
		ID:       "gemini",
		Name:     "Gemini CLI",
		Priority: PriorityGemini,
		Provider: &GeminiProviderV2{},
	})
	if err != nil {
		// Panic on registration failure since this is called at init time
		// and indicates a programming error (e.g., duplicate ID)
		panic("failed to register gemini provider: " + err.Error())
	}
}

// GeminiProviderV2 implements the ProviderV2 interface for Gemini CLI.
//
// Gemini CLI uses:
//   - NO instruction file (configFile is empty)
//   - .gemini/commands/spectr/ for slash commands
//   - TOML format for slash commands
type GeminiProviderV2 struct{}

// Initializers returns the list of Initializers needed to configure
// Gemini CLI for use with spectr.
//
// Returns:
//   - DirectoryInitializer for .gemini/commands/spectr
//   - SlashCommandsInitializer for .gemini/commands/spectr with TOML format
//
// Note: Gemini has no instruction file, so no ConfigFileInitializer.
func (*GeminiProviderV2) Initializers(_ context.Context) []Initializer {
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

// Ensure GeminiProviderV2 implements the ProviderV2 interface.
var _ ProviderV2 = (*GeminiProviderV2)(nil)
