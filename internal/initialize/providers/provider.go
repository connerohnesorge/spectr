// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// # Overview
//
// This package defines the Provider interface and registration system that all
// AI CLI tools (Claude Code, Gemini CLI, Cline, Cursor, etc.) must implement.
//
// Each provider handles both its instruction file (e.g., CLAUDE.md) and slash
// commands (e.g., .claude/commands/) through composable Initializers.
//
// # Adding a New Provider
//
// To add a new AI CLI provider, create a new file (e.g., providers/mytools.go):
//
// Example:
//
//	package providers
//
//	func init() {
//	    err := Register(Registration{
//	        ID:       "mytool",
//	        Name:     "MyTool",
//	        Priority: 100,
//	        Provider: &MyToolProvider{},
//	    })
//	    if err != nil {
//	        panic(err)
//	    }
//	}
//
//	type MyToolProvider struct{}
//
//	func (p *MyToolProvider) Initializers(ctx context.Context) []Initializer {
//	    return []Initializer{
//	        NewDirectoryInitializer(".mytool/commands/spectr"),
//	        NewConfigFileInitializer("MYTOOL.md"),
//	        NewSlashCommandsInitializerWithFrontmatter(
//	            ".mytool/commands/spectr",
//	            ".md",
//	            FormatMarkdown,
//	            map[string]string{
//	                "proposal": FrontmatterProposal,
//	                "apply":    FrontmatterApply,
//	            },
//	        ),
//	    }
//	}
//
// Each provider returns a list of Initializers that handle directory creation,
// config file management, and slash command generation.
package providers

import "context"

// CommandFormat specifies the format for slash command files.
type CommandFormat int

const (
	// FormatMarkdown uses markdown files with
	// YAML frontmatter (Claude, Cline, etc.)
	FormatMarkdown CommandFormat = iota
	// FormatTOML uses TOML files (Gemini CLI)
	FormatTOML
)

// Provider represents an AI CLI/IDE tool that can be initialized by spectr.
//
// The Provider interface is intentionally minimal - it has only ONE method.
// All metadata (ID, Name, Priority) is provided at registration time via
// the Registration struct, not through the Provider interface.
//
// # Design Philosophy
//
// The design recognizes that:
//
//  1. Metadata belongs to the registry, not the provider
//  2. Initialization logic is composable via Initializers
//  3. Each provider is just a factory for its initializers
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
type Provider interface {
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

// ProviderFunc is a function type that implements the Provider interface.
//
// This allows simple providers to be defined without creating a struct:
//
//	providers.Register(Registration{
//	    ID:       "simple-tool",
//	    Name:     "Simple Tool",
//	    Priority: 50,
//	    Provider: ProviderFunc(func(ctx context.Context) []Initializer {
//	        return []Initializer{
//	            NewConfigFileInitializer("SIMPLE.md"),
//	        }
//	    }),
//	})
type ProviderFunc func(ctx context.Context) []Initializer

// Initializers implements the Provider interface for ProviderFunc.
func (f ProviderFunc) Initializers(ctx context.Context) []Initializer {
	return f(ctx)
}
