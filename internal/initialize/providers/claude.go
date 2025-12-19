// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file implements the Claude Code provider using the ProviderV2 interface.
// Claude uses CLAUDE.md for instructions and .claude/commands/ for commands.
package providers

import (
	"context"
)

func init() {
	// Register with RegistryV2
	err := RegisterV2(Registration{
		ID:       "claude-code",
		Name:     "Claude Code",
		Priority: PriorityClaudeCode,
		Provider: &ClaudeProviderV2{},
	})
	if err != nil {
		// Panic on registration failure since this is called at init time
		// and indicates a programming error (e.g., duplicate ID)
		panic("failed to register claude-code provider: " + err.Error())
	}
}

// ClaudeProviderV2 implements the ProviderV2 interface for Claude Code.
//
// Claude Code uses:
//   - CLAUDE.md for instruction file (with spectr markers)
//   - .claude/commands/spectr/ for slash commands
//   - Markdown format for slash commands with YAML frontmatter
type ClaudeProviderV2 struct{}

// Initializers returns the list of Initializers needed to configure
// Claude Code for use with spectr.
//
// Returns:
//   - DirectoryInitializer for .claude/commands/spectr
//   - ConfigFileInitializer for CLAUDE.md
//   - SlashCommandsInitializer for .claude/commands/spectr (Markdown)
func (*ClaudeProviderV2) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		// Create the slash commands directory
		NewDirectoryInitializer(false, ".claude/commands/spectr"),

		// Create/update the CLAUDE.md instruction file
		NewConfigFileInitializer("CLAUDE.md", "instruction-pointer", false),

		// Create/update slash commands (proposal.md, apply.md)
		NewSlashCommandsInitializer(
			".claude/commands/spectr",
			".md",
			FormatMarkdown,
			StandardFrontmatter(),
			false,
		),
	}
}

// Ensure ClaudeProviderV2 implements the ProviderV2 interface.
var _ ProviderV2 = (*ClaudeProviderV2)(nil)
