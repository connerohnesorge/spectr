// Package providers defines the core interfaces for the provider architecture.
//
// This file contains the Initializer interface, which is the fundamental
// building block of the new provider system. Providers return a list of
// initializers that handle specific initialization tasks (creating directories,
// writing config files, creating slash commands).
package providers

import (
	"context"

	"github.com/spf13/afero"
)

// Initializer defines the interface for a single initialization step.
//
// Initializers are the core building blocks of the provider system. Each
// initializer is responsible for creating or updating a specific file or
// directory. Providers compose multiple initializers to set up their
// complete configuration.
//
// # Design Principles
//
// 1. **Idempotent**: Init() must be safe to run multiple times. Running it
// twice should produce the same result as running it once.
//
// 2. **Single Responsibility**: Each initializer manages exactly one path
// (file or directory). This enables deduplication when multiple providers
// need the same artifact.
//
// 3. **Filesystem Abstraction**: All file operations use afero.Fs, allowing
// for easy testing with in-memory filesystems.
//
// # Execution Order
//
// Initializers are sorted by type before execution (guaranteed order):
//  1. DirectoryInitializer - Create directories first
//  2. ConfigFileInitializer - Then config files (may need directories)
//  3. SlashCommandsInitializer - Then slash commands (may need directories)
//
// # Deduplication
//
// When multiple providers return initializers with the same Path(), only
// the first one runs. This allows providers to share common initializers
// (e.g., multiple providers using the same CLAUDE.md file).
//
// # Example Implementation
//
//	type MyInitializer struct {
//	    path string
//	}
//
//	func (i *MyInitializer) Init(ctx context.Context, fs afero.Fs, cfg *Config, tm *TemplateManager) error {
//	    // Create or update files in fs
//	    return nil
//	}
//
//	func (i *MyInitializer) IsSetup(fs afero.Fs, cfg *Config) bool {
//	    exists, _ := afero.Exists(fs, i.path)
//	    return exists
//	}
//
//	func (i *MyInitializer) Path() string {
//	    return i.path
//	}
//
//	func (i *MyInitializer) IsGlobal() bool {
//	    return false // Uses project filesystem
//	}
type Initializer interface {
	// Init creates or updates files managed by this initializer.
	//
	// Init must be idempotent: calling it multiple times should produce the
	// same result as calling it once. This allows users to safely re-run
	// `spectr init` without breaking existing configurations.
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeouts
	//   - fs: Filesystem to operate on (project or global, based on IsGlobal())
	//   - cfg: Configuration containing spectr directory paths
	//   - tm: Template manager for rendering templates
	//
	// Returns an error if initialization fails. Partial failures should be
	// avoided - either complete the initialization or return an error without
	// modifying the filesystem.
	Init(ctx context.Context, fs afero.Fs, cfg *Config, tm *TemplateManager) error

	// IsSetup returns true if this initializer's artifacts already exist.
	//
	// This is used to determine if initialization can be skipped and to
	// report the current state to users. The check should be lightweight
	// (e.g., file existence) rather than deep validation.
	//
	// Parameters:
	//   - fs: Filesystem to check (project or global, based on IsGlobal())
	//   - cfg: Configuration containing spectr directory paths
	IsSetup(fs afero.Fs, cfg *Config) bool

	// Path returns the file or directory path this initializer manages.
	//
	// This path is used for:
	//   1. Deduplication: initializers with the same path run only once
	//   2. Reporting: showing users what files will be created/updated
	//
	// The path should be relative to the filesystem root (either project
	// root for project files, or home directory for global files).
	//
	// Examples:
	//   - "CLAUDE.md" (project config file)
	//   - ".claude/commands/spectr" (project directory)
	//   - ".config/aider/commands/spectr" (global directory)
	Path() string

	// IsGlobal returns true if this initializer uses the global filesystem.
	//
	// Project initializers (IsGlobal() == false) operate on files within
	// the project directory (e.g., CLAUDE.md, .claude/commands/).
	//
	// Global initializers (IsGlobal() == true) operate on files in the
	// user's home directory (e.g., ~/.config/aider/commands/).
	//
	// The executor provides the appropriate filesystem based on this flag.
	IsGlobal() bool
}

// Config is defined in config.go and contains initialization configuration.
// It provides the spectr directory path and methods for computing derived paths.

// TemplateManager is a forward declaration for the template manager type.
// The actual implementation is in internal/initialize/templates.go.
// This type alias allows initializers to receive the template manager
// without creating import cycles.
type TemplateManager = interface {
	// RenderInstructionPointer renders the instruction pointer template.
	RenderInstructionPointer(ctx TemplateContext) (string, error)

	// RenderSlashCommand renders a slash command template.
	// commandType must be one of: "proposal", "apply", "archive"
	RenderSlashCommand(commandType string, ctx TemplateContext) (string, error)
}
