// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file implements the Cline provider using the Provider interface.
// Cline uses CLINE.md and .clinerules/commands/ for slash commands.
package providers

import (
	"context"
)

func init() {
	// Register with Registry
	err := Register(Registration{
		ID:       "cline",
		Name:     "Cline",
		Priority: PriorityCline,
		Provider: &ClineProvider{},
	})
	if err != nil {
		// Panic on registration failure since this is called at init time
		// and indicates a programming error (e.g., duplicate ID)
		panic("failed to register cline provider: " + err.Error())
	}
}

// ClineProvider implements the Provider interface for Cline.
//
// Cline uses:
//   - CLINE.md for instruction file (with spectr markers)
//   - .clinerules/commands/spectr/ for slash commands
//   - Markdown format for slash commands with YAML frontmatter
type ClineProvider struct{}

// Initializers returns the list of Initializers needed to configure
// Cline for use with spectr.
//
// Returns:
//   - DirectoryInitializer for .clinerules/commands/spectr
//   - ConfigFileInitializer for CLINE.md
//   - SlashCommandsInitializer for .clinerules/commands/spectr (Markdown)
func (*ClineProvider) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		// Create the slash commands directory
		NewDirectoryInitializer(false, ".clinerules/commands/spectr"),

		// Create/update the CLINE.md instruction file
		NewConfigFileInitializer("CLINE.md", "instruction-pointer", false),

		// Create/update slash commands (proposal.md, apply.md)
		NewSlashCommandsInitializer(
			".clinerules/commands/spectr",
			".md",
			FormatMarkdown,
			StandardFrontmatter(),
			false,
		),
	}
}

// Ensure ClineProvider implements the Provider interface.
var _ Provider = (*ClineProvider)(nil)
