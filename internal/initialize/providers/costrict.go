// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file implements the CoStrict provider using the ProviderV2 interface.
// CoStrict uses COSTRICT.md and .costrict/commands/ for slash commands.
package providers

import (
	"context"
)

func init() {
	// Register with RegistryV2
	err := RegisterV2(Registration{
		ID:       "costrict",
		Name:     "CoStrict",
		Priority: PriorityCostrict,
		Provider: &CostrictProviderV2{},
	})
	if err != nil {
		// Panic on registration failure since this is called at init time
		// and indicates a programming error (e.g., duplicate ID)
		panic("failed to register costrict provider: " + err.Error())
	}
}

// CostrictProviderV2 implements the ProviderV2 interface for CoStrict.
//
// CoStrict uses:
//   - COSTRICT.md for instruction file (with spectr markers)
//   - .costrict/commands/spectr/ for slash commands
//   - Markdown format for slash commands with YAML frontmatter
type CostrictProviderV2 struct{}

// Initializers returns the list of Initializers needed to configure
// CoStrict for use with spectr.
//
// Returns:
//   - DirectoryInitializer for .costrict/commands/spectr
//   - ConfigFileInitializer for COSTRICT.md
//   - SlashCommandsInitializer for .costrict/commands/spectr (Markdown)
func (*CostrictProviderV2) Initializers(_ context.Context) []Initializer {
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

// Ensure CostrictProviderV2 implements the ProviderV2 interface.
var _ ProviderV2 = (*CostrictProviderV2)(nil)
