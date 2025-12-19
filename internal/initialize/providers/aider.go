// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file implements the Aider provider using the ProviderV2 interface.
// Aider uses .aider/commands/ for slash commands (no config file).
package providers

import (
	"context"
)

func init() {
	// Register with RegistryV2
	err := RegisterV2(Registration{
		ID:       "aider",
		Name:     "Aider",
		Priority: PriorityAider,
		Provider: &AiderProviderV2{},
	})
	if err != nil {
		// Panic on registration failure since this is called at init time
		// and indicates a programming error (e.g., duplicate ID)
		panic("failed to register aider provider: " + err.Error())
	}
}

// AiderProviderV2 implements the ProviderV2 interface for Aider.
//
// Aider uses:
//   - No instruction file (configFile is empty)
//   - .aider/commands/spectr/ for slash commands
//   - Markdown format for slash commands with YAML frontmatter
type AiderProviderV2 struct{}

// Initializers returns the list of Initializers needed to configure
// Aider for use with spectr.
//
// Returns:
//   - DirectoryInitializer for .aider/commands/spectr
//   - SlashCommandsInitializer for .aider/commands/spectr with Markdown format
//
// Note: No ConfigFileInitializer since Aider doesn't use an instruction file.
func (*AiderProviderV2) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		// Create the slash commands directory
		NewDirectoryInitializer(false, ".aider/commands/spectr"),

		// Create/update slash commands (proposal.md, apply.md)
		NewSlashCommandsInitializer(
			".aider/commands/spectr",
			".md",
			FormatMarkdown,
			StandardFrontmatter(),
			false,
		),
	}
}

// Ensure AiderProviderV2 implements the ProviderV2 interface.
var _ ProviderV2 = (*AiderProviderV2)(nil)
