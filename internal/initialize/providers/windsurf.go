// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file implements the Windsurf provider using the Provider interface.
// Windsurf uses .windsurf/commands/ for slash commands (no instruction file).
package providers

import (
	"context"
)

func init() {
	// Register with Registry
	err := Register(Registration{
		ID:       "windsurf",
		Name:     "Windsurf",
		Priority: PriorityWindsurf,
		Provider: &WindsurfProvider{},
	})
	if err != nil {
		// Panic on registration failure since this is called at init time
		// and indicates a programming error (e.g., duplicate ID)
		panic("failed to register windsurf provider: " + err.Error())
	}
}

// WindsurfProvider implements the Provider interface for Windsurf.
//
// Windsurf uses:
//   - No instruction file
//   - .windsurf/commands/spectr/ for slash commands
//   - Markdown format for slash commands with YAML frontmatter
type WindsurfProvider struct{}

// Initializers returns the list of Initializers needed to configure
// Windsurf for use with spectr.
//
// Note: Windsurf has no instruction file (configFile is empty), so no
// ConfigFileInitializer is returned.
//
// Returns:
//   - DirectoryInitializer for .windsurf/commands/spectr
//   - SlashCommandsInitializer for .windsurf/commands/spectr (Markdown)
func (*WindsurfProvider) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		// Create the slash commands directory
		NewDirectoryInitializer(false, ".windsurf/commands/spectr"),

		// No config file for Windsurf (empty configFile)

		// Create/update slash commands (proposal.md, apply.md)
		NewSlashCommandsInitializer(
			".windsurf/commands/spectr",
			".md",
			FormatMarkdown,
			StandardFrontmatter(),
			false,
		),
	}
}

// Ensure WindsurfProvider implements the Provider interface.
var _ Provider = (*WindsurfProvider)(nil)
