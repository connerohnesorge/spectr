// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file implements the Crush provider using the Provider interface.
// Crush uses CRUSH.md for instructions and .crush/commands/ for slash commands.
package providers

import (
	"context"
)

func init() {
	// Register with Registry
	err := Register(Registration{
		ID:       "crush",
		Name:     "Crush",
		Priority: PriorityCrush,
		Provider: &CrushProvider{},
	})
	if err != nil {
		// Panic on registration failure since this is called at init time
		// and indicates a programming error (e.g., duplicate ID)
		panic("failed to register crush provider: " + err.Error())
	}
}

// CrushProvider implements the Provider interface for Crush.
//
// Crush uses:
//   - CRUSH.md for instruction file (with spectr markers)
//   - .crush/commands/spectr/ for slash commands
//   - Markdown format for slash commands with YAML frontmatter
type CrushProvider struct{}

// Initializers returns the list of Initializers needed to configure
// Crush for use with spectr.
//
// Returns:
//   - DirectoryInitializer for .crush/commands/spectr
//   - ConfigFileInitializer for CRUSH.md
//   - SlashCommandsInitializer for .crush/commands/spectr (Markdown)
func (*CrushProvider) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		// Create the slash commands directory
		NewDirectoryInitializer(false, ".crush/commands/spectr"),

		// Create/update the CRUSH.md instruction file
		NewConfigFileInitializer("CRUSH.md", "instruction-pointer", false),

		// Create/update slash commands (proposal.md, apply.md)
		NewSlashCommandsInitializer(
			".crush/commands/spectr",
			".md",
			FormatMarkdown,
			StandardFrontmatter(),
			false,
		),
	}
}

// Ensure CrushProvider implements the Provider interface.
var _ Provider = (*CrushProvider)(nil)
