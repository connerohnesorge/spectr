// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file implements the Continue provider using the ProviderV2 interface.
// Continue uses .continue/commands/ for slash commands (no instruction file).
package providers

import (
	"context"
)

func init() {
	// Register with RegistryV2
	err := RegisterV2(Registration{
		ID:       "continue",
		Name:     "Continue",
		Priority: PriorityContinue,
		Provider: &ContinueProviderV2{},
	})
	if err != nil {
		// Panic on registration failure since this is called at init time
		// and indicates a programming error (e.g., duplicate ID)
		panic("failed to register continue provider: " + err.Error())
	}
}

// ContinueProviderV2 implements the ProviderV2 interface for Continue.
//
// Continue uses:
//   - No instruction file
//   - .continue/commands/spectr/ for slash commands
//   - Markdown format for slash commands with YAML frontmatter
type ContinueProviderV2 struct{}

// Initializers returns the list of Initializers needed to configure
// Continue for use with spectr.
//
// Note: Continue has no instruction file (configFile is empty), so no
// ConfigFileInitializer is returned.
//
// Returns:
//   - DirectoryInitializer for .continue/commands/spectr
//   - SlashCommandsInitializer for .continue/commands/spectr (Markdown)
func (*ContinueProviderV2) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		// Create the slash commands directory
		NewDirectoryInitializer(false, ".continue/commands/spectr"),

		// No config file for Continue (empty configFile)

		// Create/update slash commands (proposal.md, apply.md)
		NewSlashCommandsInitializer(
			".continue/commands/spectr",
			".md",
			FormatMarkdown,
			StandardFrontmatter(),
			false,
		),
	}
}

// Ensure ContinueProviderV2 implements the ProviderV2 interface.
var _ ProviderV2 = (*ContinueProviderV2)(nil)
