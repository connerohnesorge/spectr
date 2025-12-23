// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file implements the Kilocode provider using the Provider interface.
// Kilocode uses .kilocode/commands/ for slash commands (no instruction file).
package providers

import (
	"context"
)

func init() {
	// Register with Registry
	err := Register(Registration{
		ID:       "kilocode",
		Name:     "Kilocode",
		Priority: PriorityKilocode,
		Provider: &KilocodeProvider{},
	})
	if err != nil {
		// Panic on registration failure since this is called at init time
		// and indicates a programming error (e.g., duplicate ID)
		panic("failed to register kilocode provider: " + err.Error())
	}
}

// KilocodeProvider implements the Provider interface for Kilocode.
//
// Kilocode uses:
//   - No instruction file
//   - .kilocode/commands/spectr/ for slash commands
//   - Markdown format for slash commands with YAML frontmatter
type KilocodeProvider struct{}

// Initializers returns the list of Initializers needed to configure
// Kilocode for use with spectr.
//
// Note: Kilocode has no instruction file (configFile is empty), so no
// ConfigFileInitializer is returned.
//
// Returns:
//   - DirectoryInitializer for .kilocode/commands/spectr
//   - SlashCommandsInitializer for .kilocode/commands/spectr (Markdown)
func (*KilocodeProvider) Initializers(_ context.Context) []Initializer {
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

// Ensure KilocodeProvider implements the Provider interface.
var _ Provider = (*KilocodeProvider)(nil)
