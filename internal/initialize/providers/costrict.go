// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file implements the CoStrict provider using the Provider interface.
// CoStrict uses COSTRICT.md and .costrict/commands/ for slash commands.
package providers

import (
	"context"
)

func init() {
	// Register with Registry
	err := Register(Registration{
		ID:       "costrict",
		Name:     "CoStrict",
		Priority: PriorityCostrict,
		Provider: &CostrictProvider{},
	})
	if err != nil {
		// Panic on registration failure since this is called at init time
		// and indicates a programming error (e.g., duplicate ID)
		panic("failed to register costrict provider: " + err.Error())
	}
}

// CostrictProvider implements the Provider interface for CoStrict.
//
// CoStrict uses:
//   - COSTRICT.md for instruction file (with spectr markers)
//   - .costrict/commands/spectr/ for slash commands
//   - Markdown format for slash commands with YAML frontmatter
type CostrictProvider struct{}

// Initializers returns the list of Initializers needed to configure
// CoStrict for use with spectr.
//
// Returns:
//   - DirectoryInitializer for .costrict/commands/spectr
//   - ConfigFileInitializer for COSTRICT.md
//   - SlashCommandsInitializer for .costrict/commands/spectr (Markdown)
func (*CostrictProvider) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		// Create the slash commands directory
		NewDirectoryInitializer(false, ".costrict/commands/spectr"),

		// Create/update the COSTRICT.md instruction file
		NewConfigFileInitializer("COSTRICT.md", "instruction-pointer", false),

		// Create/update slash commands (proposal.md, apply.md)
		NewSlashCommandsInitializer(
			".costrict/commands/spectr",
			".md",
			FormatMarkdown,
			StandardFrontmatter(),
			false,
		),
	}
}

// Ensure CostrictProvider implements the Provider interface.
var _ Provider = (*CostrictProvider)(nil)
