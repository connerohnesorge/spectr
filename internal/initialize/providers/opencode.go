// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file implements the OpenCode provider using the Provider interface.
// OpenCode uses .opencode/command/spectr/ for slash commands.
package providers

import (
	"context"
)

func init() {
	// Register with Registry
	err := Register(Registration{
		ID:       "opencode",
		Name:     "OpenCode",
		Priority: PriorityOpencode,
		Provider: &OpencodeProvider{},
	})
	if err != nil {
		// Panic on registration failure since this is called at init time
		// and indicates a programming error (e.g., duplicate ID)
		panic("failed to register opencode provider: " + err.Error())
	}
}

// OpencodeProvider implements the Provider interface for OpenCode.
//
// OpenCode uses:
//   - No instruction file (uses JSON configuration)
//   - .opencode/command/spectr/ for slash commands
//   - Markdown format for slash commands with YAML frontmatter
type OpencodeProvider struct{}

// Initializers returns the list of Initializers needed to configure
// OpenCode for use with spectr.
//
// Note: OpenCode has no instruction file (it uses JSON configuration),
// so no ConfigFileInitializer is returned.
//
// Returns:
//   - DirectoryInitializer for .opencode/command/spectr
//   - SlashCommandsInitializer for .opencode/command/spectr (Markdown)
func (*OpencodeProvider) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		// Create the slash commands directory
		NewDirectoryInitializer(false, ".opencode/command/spectr"),

		// No config file for OpenCode (uses JSON configuration)

		// Create/update slash commands (proposal.md, apply.md)
		NewSlashCommandsInitializer(
			".opencode/command/spectr",
			".md",
			FormatMarkdown,
			StandardFrontmatter(),
			false,
		),
	}
}

// Ensure OpencodeProvider implements the Provider interface.
var _ Provider = (*OpencodeProvider)(nil)
