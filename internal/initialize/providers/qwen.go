// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file implements the Qwen Code provider using the Provider interface.
// Qwen uses QWEN.md for instructions and .qwen/commands/ for slash commands.
package providers

import (
	"context"
)

func init() {
	// Register with Registry
	err := Register(Registration{
		ID:       "qwen",
		Name:     "Qwen Code",
		Priority: PriorityQwen,
		Provider: &QwenProvider{},
	})
	if err != nil {
		// Panic on registration failure since this is called at init time
		// and indicates a programming error (e.g., duplicate ID)
		panic("failed to register qwen provider: " + err.Error())
	}
}

// QwenProvider implements the Provider interface for Qwen Code.
//
// Qwen Code uses:
//   - QWEN.md for instruction file (with spectr markers)
//   - .qwen/commands/spectr/ for slash commands
//   - Markdown format for slash commands with YAML frontmatter
type QwenProvider struct{}

// Initializers returns the list of Initializers needed to configure
// Qwen Code for use with spectr.
//
// Returns:
//   - DirectoryInitializer for .qwen/commands/spectr
//   - ConfigFileInitializer for QWEN.md
//   - SlashCommandsInitializer for .qwen/commands/spectr (Markdown)
func (*QwenProvider) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		// Create the slash commands directory
		NewDirectoryInitializer(false, ".qwen/commands/spectr"),

		// Create/update the QWEN.md instruction file
		NewConfigFileInitializer("QWEN.md", "instruction-pointer", false),

		// Create/update slash commands (proposal.md, apply.md)
		NewSlashCommandsInitializer(
			".qwen/commands/spectr",
			".md",
			FormatMarkdown,
			StandardFrontmatter(),
			false,
		),
	}
}

// Ensure QwenProvider implements the Provider interface.
var _ Provider = (*QwenProvider)(nil)
