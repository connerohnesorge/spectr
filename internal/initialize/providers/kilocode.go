// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file implements the Kilocode provider using the ProviderV2 interface.
// Kilocode uses .kilocode/commands/ for slash commands (no instruction file).
package providers

import (
	"context"
)

func init() {
	// Register with RegistryV2
	err := RegisterV2(Registration{
		ID:       "kilocode",
		Name:     "Kilocode",
		Priority: PriorityKilocode,
		Provider: &KilocodeProviderV2{},
	})
	if err != nil {
		// Panic on registration failure since this is called at init time
		// and indicates a programming error (e.g., duplicate ID)
		panic("failed to register kilocode provider: " + err.Error())
	}
}

// KilocodeProviderV2 implements the ProviderV2 interface for Kilocode.
//
// Kilocode uses:
//   - No instruction file
//   - .kilocode/commands/spectr/ for slash commands
//   - Markdown format for slash commands with YAML frontmatter
type KilocodeProviderV2 struct{}

// Initializers returns the list of Initializers needed to configure
// Kilocode for use with spectr.
//
// Note: Kilocode has no instruction file (configFile is empty), so no
// ConfigFileInitializer is returned.
//
// Returns:
//   - DirectoryInitializer for .kilocode/commands/spectr
//   - SlashCommandsInitializer for .kilocode/commands/spectr (Markdown)
func (*KilocodeProviderV2) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		// Create the slash commands directory
		NewDirectoryInitializer(false, ".kilocode/commands/spectr"),

		// No config file for Kilocode (empty configFile)

		// Create/update slash commands (proposal.md, apply.md)
		NewSlashCommandsInitializer(
			".kilocode/commands/spectr",
			".md",
			FormatMarkdown,
			StandardFrontmatter(),
			false,
		),
	}
}

// Ensure KilocodeProviderV2 implements the ProviderV2 interface.
var _ ProviderV2 = (*KilocodeProviderV2)(nil)
