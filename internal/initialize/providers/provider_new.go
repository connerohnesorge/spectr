// Package providers defines the core interfaces for the provider architecture.
//
// This file contains the new minimal Provider interface (ProviderV2) that will
// replace the complex 12-method Provider interface in provider.go.
// The new design separates concerns:
//   - Metadata (ID, Name, Priority) is provided at registration time
//   - Providers only need to implement one method: Initializers()
//
// This reduces provider implementation from ~50 lines to ~10 lines.
//
// # Migration Note
//
// During the migration period, this interface is named ProviderV2 to avoid
// conflicts with the existing Provider interface in provider.go. Once all
// providers are migrated (tasks 5.1-5.16) and the old interface is removed
// (task 7.1), ProviderV2 will be renamed to Provider.
package providers

import "context"

// ProviderV2 represents an AI CLI/IDE tool that can be initialized by spectr.
//
// The ProviderV2 interface is intentionally minimal - it has only ONE method.
// All metadata (ID, Name, Priority) is provided at registration time via
// the Registration struct, not through the ProviderV2 interface.
//
// # Migration Note
//
// This interface is named ProviderV2 during the migration period. After the
// old Provider interface is removed, this will be renamed to Provider.
//
// # Design Philosophy
//
// The previous Provider interface had 12 methods, most of which were
// boilerplate handled by BaseProvider. The new design recognizes that:
//
//  1. Metadata belongs to the registry, not the provider
//  2. Initialization logic is composable via Initializers
//  3. Each provider is just a factory for its initializers
//
// # Initializers vs Provider Methods
//
// Instead of implementing Configure(), IsConfigured(), GetFilePaths(), etc.,
// providers now return a list of composable Initializers. Each initializer
// handles one concern (directory creation, config file, slash commands).
//
// # Example: Claude Code Provider
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
//	// Registration happens separately in init()
//	func init() {
//	    providers.Register(Registration{
//	        ID:       "claude-code",
//	        Name:     "Claude Code",
//	        Priority: 1,
//	        Provider: &ClaudeProvider{},
//	    })
//	}
//
// # Example: Gemini Provider (TOML format, no config file)
//
//	type GeminiProvider struct{}
//
//	func (p *GeminiProvider) Initializers(ctx context.Context) []Initializer {
//	    return []Initializer{
//	        NewDirectoryInitializer(".gemini/commands/spectr"),
//	        // Note: No config file initializer - Gemini doesn't use one
//	        NewSlashCommandsInitializer(
//	            ".gemini/commands/spectr", ".toml", FormatTOML,
//	        ),
//	    }
//	}
//
// # Context Usage
//
// The context parameter allows for:
//   - Cancellation during long-running initialization
//   - Deadline propagation from parent commands
//   - Passing request-scoped values (though rarely needed)
//
// Providers may ignore the context if they don't need it, but the parameter
// is provided for consistency and future extensibility.
//
// # Initializer Ordering
//
// Providers do NOT need to order their initializers correctly. The executor
// automatically sorts initializers by type before execution:
//
//  1. DirectoryInitializer - Create directories first
//  2. ConfigFileInitializer - Then config files
//  3. SlashCommandsInitializer - Then slash commands
//
// # Deduplication
//
// When multiple providers return initializers with the same Path(), only
// the first one runs. This allows providers to share common initializers
// without duplicate execution.
type ProviderV2 interface {
	// Initializers returns the list of initializers for this provider.
	//
	// Each initializer handles a specific initialization task:
	//   - DirectoryInitializer: Creates directories (e.g.,
	//     .claude/commands/spectr)
	//   - ConfigFileInitializer: Creates/updates config files (e.g., CLAUDE.md)
	//   - SlashCommandsInitializer: Creates slash command files
	//
	// The returned initializers are:
	//   1. Collected from all selected providers
	//   2. Deduplicated by Path() (same path = run once)
	//   3. Sorted by type (directories before files)
	//   4. Executed in order
	//
	// Providers may return an empty slice if they have no initialization needs,
	// though this is uncommon in practice.
	//
	// Parameters:
	//   - ctx: Context for cancellation and deadlines
	//     (may be ignored if not needed)
	//
	// Returns:
	//   - A slice of Initializers that will create/update this provider's files
	Initializers(ctx context.Context) []Initializer
}

// ProviderV2Func is a function type that implements the ProviderV2 interface.
//
// This allows simple providers to be defined without creating a struct:
//
//	providers.Register(Registration{
//	    ID:       "simple-tool",
//	    Name:     "Simple Tool",
//	    Priority: 50,
//	    Provider: ProviderV2Func(func(ctx context.Context) []Initializer {
//	        return []Initializer{
//	            NewConfigFileInitializer("SIMPLE.md", template),
//	        }
//	    }),
//	})
//
// # Migration Note
//
// This type will be renamed to ProviderFunc after the old Provider interface
// is removed.
type ProviderV2Func func(ctx context.Context) []Initializer

// Initializers implements the ProviderV2 interface for ProviderV2Func.
func (f ProviderV2Func) Initializers(ctx context.Context) []Initializer {
	return f(ctx)
}
