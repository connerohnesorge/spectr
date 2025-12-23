// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file implements the Qoder provider using the Provider interface.
// Qoder uses QODER.md for instructions and .qoder/commands/ for slash commands.
package providers

import (
	"context"
)

func init() {
	// Register with Registry
	err := Register(Registration{
		ID:       "qoder",
		Name:     "Qoder",
		Priority: PriorityQoder,
		Provider: &QoderProvider{},
	})
	if err != nil {
		// Panic on registration failure since this is called at init time
		// and indicates a programming error (e.g., duplicate ID)
		panic("failed to register qoder provider: " + err.Error())
	}
}

// QoderProvider implements the Provider interface for Qoder.
//
// Qoder uses:
//   - QODER.md for instruction file (with spectr markers)
//   - .qoder/commands/spectr/ for slash commands
//   - Markdown format for slash commands with YAML frontmatter
type QoderProvider struct{}

// Initializers returns the list of Initializers needed to configure
// Qoder for use with spectr.
//
// Returns:
//   - DirectoryInitializer for .qoder/commands/spectr
//   - ConfigFileInitializer for QODER.md
//   - SlashCommandsInitializer for .qoder/commands/spectr (Markdown)
func (*QoderProvider) Initializers(_ context.Context) []Initializer {
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

// Ensure QoderProvider implements the Provider interface.
var _ Provider = (*QoderProvider)(nil)
