// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file implements the Codex CLI provider using the Provider interface.
// Codex uses AGENTS.md and global ~/.codex/prompts/ for commands.
package providers

import (
	"context"
)

func init() {
	// Register with Registry
	err := Register(Registration{
		ID:       "codex",
		Name:     "Codex CLI",
		Priority: PriorityCodex,
		Provider: &CodexProvider{},
	})
	if err != nil {
		// Panic on registration failure since this is called at init time
		// and indicates a programming error (e.g., duplicate ID)
		panic("failed to register codex provider: " + err.Error())
	}
}

// CodexProvider implements the Provider interface for Codex CLI.
//
// Codex uses:
//   - AGENTS.md for instruction file (with spectr markers) - project-relative
//   - ~/.codex/prompts/ for slash commands - global (home directory)
//   - Markdown format for slash commands with YAML frontmatter
//   - Prefixed command files (spectr-proposal.md, spectr-apply.md)
type CodexProvider struct{}

// Initializers returns the list of Initializers needed to configure
// Codex CLI for use with spectr.
//
// Returns:
//   - DirectoryInitializer for ~/.codex/prompts (global)
//   - ConfigFileInitializer for AGENTS.md (project-relative)
//   - SlashCommandsInitializer for ~/.codex/prompts (Markdown, global)
//
// Note: Codex uses global paths for commands, project-relative for config.
func (*CodexProvider) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		// Create the global prompts directory
		// Note: Path is relative to home directory when isGlobal=true
		NewDirectoryInitializer(true, ".codex/prompts"),

		// Create/update the AGENTS.md instruction file (project-relative)
		NewConfigFileInitializer("AGENTS.md", "instruction-pointer", false),

		// Create/update slash commands in global prompts
		// Note: Using prefixed command pattern for Codex
		NewSlashCommandsInitializer(
			".codex/prompts",
			".md",
			FormatMarkdown,
			StandardFrontmatter(),
			true, // isGlobal - commands are in home directory
		),
	}
}

// Ensure CodexProvider implements the Provider interface.
var _ Provider = (*CodexProvider)(nil)
