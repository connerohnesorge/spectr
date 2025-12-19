// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file implements the Cursor provider using the ProviderV2 interface.
// Cursor uses .cursorrules/commands/ for slash commands (no config file).
package providers

import (
	"context"
)

func init() {
	// Register with RegistryV2
	err := RegisterV2(Registration{
		ID:       "cursor",
		Name:     "Cursor",
		Priority: PriorityCursor,
		Provider: &CursorProviderV2{},
	})
	if err != nil {
		// Panic on registration failure since this is called at init time
		// and indicates a programming error (e.g., duplicate ID)
		panic("failed to register cursor provider: " + err.Error())
	}
}

// CursorProviderV2 implements the ProviderV2 interface for Cursor.
//
// Cursor uses:
//   - No instruction file (configFile is empty)
//   - .cursorrules/commands/spectr/ for slash commands
//   - Markdown format for slash commands with YAML frontmatter
type CursorProviderV2 struct{}

// Initializers returns the list of Initializers needed to configure
// Cursor for use with spectr.
//
// Returns:
//   - DirectoryInitializer for .cursorrules/commands/spectr
//   - SlashCommandsInitializer for .cursorrules/commands/spectr (Markdown)
//
// Note: No ConfigFileInitializer since Cursor doesn't use an instruction file.
func (*CursorProviderV2) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		// Create the slash commands directory
		NewDirectoryInitializer(false, ".cursorrules/commands/spectr"),

		// Create/update slash commands (proposal.md, apply.md)
		NewSlashCommandsInitializer(
			".cursorrules/commands/spectr",
			".md",
			FormatMarkdown,
			StandardFrontmatter(),
			false,
		),
	}
}

// Ensure CursorProviderV2 implements the ProviderV2 interface.
var _ ProviderV2 = (*CursorProviderV2)(nil)
