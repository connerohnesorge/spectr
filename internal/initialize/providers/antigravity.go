// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file implements the Antigravity provider using the ProviderV2 interface.
// Antigravity uses AGENTS.md for instructions and .agent/workflows/ for
// slash commands.
//
//nolint:revive // line-length-limit - provider documentation
package providers

import (
	"context"
)

func init() {
	// Register with RegistryV2
	err := RegisterV2(Registration{
		ID:       "antigravity",
		Name:     "Antigravity",
		Priority: PriorityAntigravity,
		Provider: &AntigravityProviderV2{},
	})
	if err != nil {
		// Panic on registration failure since this is called at init time
		// and indicates a programming error (e.g., duplicate ID)
		panic("failed to register antigravity provider: " + err.Error())
	}
}

// AntigravityProviderV2 implements the ProviderV2 interface for Antigravity.
//
// Antigravity uses:
//   - AGENTS.md for instruction file (with spectr markers)
//   - .agent/workflows/ for slash commands (prefixed pattern)
//   - Markdown format for slash commands with YAML frontmatter
type AntigravityProviderV2 struct{}

// Initializers returns the list of Initializers needed to configure
// Antigravity for use with spectr.
//
// Note: Antigravity uses the prefixed command paths pattern
// (.agent/workflows/spectr-proposal.md) rather than the subdirectory pattern.
// The SlashCommandsInitializer handles this by placing files directly in the
// workflows directory with the spectr- prefix.
//
// Returns:
//   - DirectoryInitializer for .agent/workflows
//   - ConfigFileInitializer for AGENTS.md
//   - SlashCommandsInitializer for .agent/workflows (Markdown, prefixed)
func (p *AntigravityProviderV2) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		// Create the workflows directory
		NewDirectoryInitializer(false, ".agent/workflows"),

		// Create/update the AGENTS.md instruction file
		NewConfigFileInitializer("AGENTS.md", "instruction-pointer", false),

		// Create/update slash commands (spectr-proposal.md, spectr-apply.md)
		// Note: Using empty dir with full paths for prefixed pattern
		NewSlashCommandsInitializer(
			".agent/workflows",
			".md",
			FormatMarkdown,
			StandardFrontmatter(),
			false,
		),
	}
}

// Ensure AntigravityProviderV2 implements the ProviderV2 interface.
var _ ProviderV2 = (*AntigravityProviderV2)(nil)
