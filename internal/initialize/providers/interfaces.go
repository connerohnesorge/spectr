// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file defines the new provider interfaces for the redesigned arch.
// The new design reduces provider boilerplate by separating concerns:
//   - Provider: Returns a list of initializers
//   - Initializer: Handles a single initialization step
//   - Registration: Contains provider metadata (ID, Name, Priority)
//   - Config: Holds configuration passed to initializers
//
// See the old Provider interface in provider.go for the legacy implementation.
package providers

import (
	"context"

	"github.com/spf13/afero"
)

// NewProvider is the new provider interface that returns a list of
// initializers. Providers no longer contain metadata (ID, Name, Priority) -
// that lives in Registration.
//
// This interface is intentionally minimal. Providers compose behavior by
// returning appropriate initializers for their requirements.
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
//	        NewConfigFileInitializer("CLAUDE.md", template),
//	        NewSlashCommandsInitializer(
//	            ".claude/commands/spectr", ".md", FormatMarkdown,
//	        ),
//	    }
//	}
type NewProvider interface {
	// Initializers returns the list of initializers needed to configure
	// this provider. The context can be used for cancellation or deadline
	// propagation.
	Initializers(
		ctx context.Context,
	) []Initializer
}

// Initializer represents a single initialization step. Each initializer
// handles one specific task (create directory, write config file, etc.).
//
// Initializers must be idempotent - running Init multiple times should
// produce the same result as running it once.
type Initializer interface {
	// Init performs the initialization step. The filesystem (fs) is rooted
	// at the project directory, so all paths are project-relative. Returns
	// an error if initialization fails.
	Init(
		ctx context.Context,
		fs afero.Fs,
		cfg *Config,
	) error

	// IsSetup returns true if this initializer's work is already complete.
	// This is used to determine if initialization is needed.
	IsSetup(fs afero.Fs, cfg *Config) bool
}

// Config holds configuration passed to initializers. This struct contains
// values that initializers need to know about the project structure.
type Config struct {
	// SpectrDir is the spectr directory relative to the project root.
	// Example: "spectr" (the default)
	SpectrDir string
}

// NewConfig creates a new Config with default values.
func NewConfig() *Config {
	return &Config{
		SpectrDir: "spectr",
	}
}

// Registration holds provider metadata for the registry.
// This separates the "what does this provider do" (Provider interface)
// from the "how is it identified and ordered" (Registration).
type Registration struct {
	// ID is the unique provider identifier (kebab-case).
	// Example: "claude-code", "gemini-cli", "cline"
	ID string

	// Name is the human-readable provider name for display.
	// Example: "Claude Code", "Gemini CLI", "Cline"
	Name string

	// Priority is the display/processing order (lower = higher priority).
	// Claude Code should be 1, other major tools 2-10, etc.
	Priority int

	// Provider is the provider implementation.
	Provider NewProvider
}

// Note: CommandFormat (FormatMarkdown, FormatTOML) is kept in provider.go
// as it is used by both old and new implementations.

// Note: TemplateContext and DefaultTemplateContext() are kept in provider.go
// as they are used for template rendering across the codebase.
