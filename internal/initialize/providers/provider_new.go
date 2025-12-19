// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file defines the new minimal ProviderV2 interface. Unlike the legacy
// 12-method Provider interface, this interface has a single method.
//
// The new design separates concerns:
//   - ProviderV2: Returns the list of initializers needed for a tool
//   - Registration: Contains metadata (ID, Name, Priority) at registration time
//   - Initializer: Handles actual file creation/update operations
//
// This reduces provider authoring from ~50 lines of boilerplate to ~10 lines
// of registration code, as providers only need to return their initializers.
//
// Note: This interface is named ProviderV2 during the migration period to
// coexist with the legacy Provider interface. Once migration is complete
// (task 7.1), the legacy interface will be removed and ProviderV2 can be
// renamed to Provider.
//
// Example usage:
//
//	type ClaudeProvider struct{}
//
//	func (p *ClaudeProvider) Initializers(ctx context.Context) []Initializer {
//	    return []Initializer{
//	        NewDirectoryInitializer(".claude/commands/spectr"),
//	        NewConfigFileInitializer("CLAUDE.md", InstructionTemplate),
//	        NewSlashCommandsInitializer(
//	            ".claude/commands/spectr", ".md", FormatMarkdown,
//	        ),
//	    }
//	}
//
// See also:
//   - initializer.go: Initializer interface for file operations
//   - config.go: Config struct with SpectrDir and derived path methods
//   - registration.go: Registration struct for provider metadata
//
//nolint:revive // line-length-limit - documentation
package providers

import "context"

// ProviderV2 defines the minimal interface for AI CLI/IDE tool providers.
//
// This is the new provider interface that replaces the legacy 12-method
// Provider interface. It is named ProviderV2 during the migration period
// to allow both interfaces to coexist. After migration is complete, the
// legacy Provider interface will be removed and this can be renamed.
//
// A ProviderV2's sole responsibility is to return the list of Initializers
// needed to configure a tool for use with spectr. All metadata (ID, Name,
// Priority) is provided at registration time via the Registration struct,
// not through the ProviderV2 interface itself.
//
// This design enables:
//   - Minimal boilerplate: providers only implement one method
//   - Composable initializers: share and reuse common initialization logic
//   - Testable: each initializer can be tested independently
//   - Deduplicable: initializers with same Path() run only once
//
// The returned initializers are:
//   - Collected from all selected providers
//   - Deduplicated by Path() (same path = run once)
//   - Sorted by type (directories first, then config files, then commands)
//   - Executed in order with proper filesystem abstraction
type ProviderV2 interface {
	// Initializers returns the list of Initializers needed to configure
	// this provider's tool for use with spectr.
	//
	// The context can be used for cancellation or to pass request-scoped
	// values. Currently unused but included for future extensibility.
	//
	// Returned initializers will be:
	//   - Deduplicated by Path() across all providers
	//   - Sorted by type (Directory -> ConfigFile -> SlashCommands)
	//   - Executed with either projectFs or globalFs based on IsGlobal()
	//
	// Returns an empty slice if the provider requires no initialization.
	Initializers(ctx context.Context) []Initializer
}
