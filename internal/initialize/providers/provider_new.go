package providers

import (
	"context"
)

// Provider represents a composable AI CLI/IDE tool provider.
//
// This is the new, simplified Provider interface that replaces
// the legacy 12-method LegacyProvider interface. Providers now
// return a list of composable Initializers that handle specific
// initialization tasks (directories, config files, slash commands).
//
// Provider metadata (ID, Name, Priority) is specified at
// registration time via the Registration struct, not in the Provider
// interface itself. This eliminates boilerplate and allows providers
// to focus on what they initialize.
//
// Example implementation:
//
//	type ClaudeProvider struct{}
//
//	func (p *ClaudeProvider) Initializers(
//	    ctx context.Context,
//	) []Initializer {
//	    return []Initializer{
//	        NewDirectoryInitializer(".claude/commands/spectr"),
//	        NewConfigFileInitializer(
//	            "CLAUDE.md",
//	            (*TemplateManager).InstructionPointer,
//	        ),
//	        NewSlashCommandsInitializer(
//	            ".claude/commands/spectr",
//	            ".md",
//	            []SlashCommand{
//	                SlashProposal,
//	                SlashApply,
//	            },
//	        ),
//	    }
//	}
//
//	func init() {
//	    providers.Register(
//	        providers.Registration{
//	            ID:       "claude-code",
//	            Name:     "Claude Code",
//	            Priority: 1,
//	            Provider: &ClaudeProvider{},
//	        },
//	    )
//	}
//
// This interface replaced the legacy LegacyProvider interface in task 7.1.
// See design.md for the complete architecture rationale.
type Provider interface {
	// Initializers returns the list of initialization tasks for this
	// provider.
	//
	// Each Initializer handles a specific aspect of setup (creating
	// directories, config files, slash commands, etc.). Initializers
	// are:
	// - Composable: Multiple providers can share the same initializer
	//   instances
	// - Deduplicated: Same path = run once, even if multiple
	//   providers return it
	// - Ordered: Executed in a guaranteed order
	//   (directories → config files → slash commands)
	// - Idempotent: Safe to run multiple times
	//
	// The context can be used for cancellation or timeout control during
	// initializer collection (not execution - that happens later in executor).
	Initializers(
		ctx context.Context,
	) []Initializer
}
