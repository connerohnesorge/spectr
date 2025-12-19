// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file implements the Qoder provider using the ProviderV2 interface.
// Qoder uses QODER.md for instructions and .qoder/commands/ for slash commands.
package providers

import (
	"context"
)

func init() {
	// Register with RegistryV2
	err := RegisterV2(Registration{
		ID:       "qoder",
		Name:     "Qoder",
		Priority: PriorityQoder,
		Provider: &QoderProviderV2{},
	})
	if err != nil {
		// Panic on registration failure since this is called at init time
		// and indicates a programming error (e.g., duplicate ID)
		panic("failed to register qoder provider: " + err.Error())
	}
}

// QoderProviderV2 implements the ProviderV2 interface for Qoder.
//
// Qoder uses:
//   - QODER.md for instruction file (with spectr markers)
//   - .qoder/commands/spectr/ for slash commands
//   - Markdown format for slash commands with YAML frontmatter
type QoderProviderV2 struct{}

// Initializers returns the list of Initializers needed to configure
// Qoder for use with spectr.
//
// Returns:
//   - DirectoryInitializer for .qoder/commands/spectr
//   - ConfigFileInitializer for QODER.md
//   - SlashCommandsInitializer for .qoder/commands/spectr (Markdown)
func (*QoderProviderV2) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		// Create the slash commands directory
		NewDirectoryInitializer(false, ".qoder/commands/spectr"),

		// Create/update the QODER.md instruction file
		NewConfigFileInitializer("QODER.md", "instruction-pointer", false),

		// Create/update slash commands (proposal.md, apply.md)
		NewSlashCommandsInitializer(
			".qoder/commands/spectr",
			".md",
			FormatMarkdown,
			StandardFrontmatter(),
			false,
		),
	}
}

// Ensure QoderProviderV2 implements the ProviderV2 interface.
var _ ProviderV2 = (*QoderProviderV2)(nil)
