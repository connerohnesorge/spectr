// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file implements the Aider provider using the Provider interface.
// Aider uses .aider/commands/ for slash commands (no config file).
package providers

import (
	"context"
)

func init() {
	// Register with Registry
	err := Register(Registration{
		ID:       "aider",
		Name:     "Aider",
		Priority: PriorityAider,
		Provider: &AiderProvider{},
	})
	if err != nil {
		// Panic on registration failure since this is called at init time
		// and indicates a programming error (e.g., duplicate ID)
		panic("failed to register aider provider: " + err.Error())
	}
}

// AiderProvider implements the Provider interface for Aider.
//
// Aider uses:
//   - No instruction file (configFile is empty)
//   - .aider/commands/spectr/ for slash commands
//   - Markdown format for slash commands with YAML frontmatter
type AiderProvider struct{}

// Initializers returns the list of Initializers needed to configure
// Aider for use with spectr.
//
// Returns:
//   - DirectoryInitializer for .aider/commands/spectr
//   - SlashCommandsInitializer for .aider/commands/spectr with Markdown format
//
// Note: No ConfigFileInitializer since Aider doesn't use an instruction file.
func (*AiderProvider) Initializers(_ context.Context) []Initializer {
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

// Ensure AiderProvider implements the Provider interface.
var _ Provider = (*AiderProvider)(nil)
