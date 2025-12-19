// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file defines the Initializer interface, the core abstraction for
// creating and updating files during spectr initialization.
//
// See also:
//   - config.go: Config struct with SpectrDir and derived path methods
//   - provider.go: TemplateRenderer interface for template rendering
//
//nolint:revive // line-length-limit - interface documentation
package providers

import (
	"context"

	"github.com/spf13/afero"
)

// Note: Config is defined in config.go with SpectrDir field and derived
// path methods (SpecsDir, ChangesDir, ProjectFile, AgentsFile).

// Note: TemplateRenderer interface is defined in provider.go and is reused
// here. It provides template rendering capabilities for initializers.

// Initializer defines the interface for components that create or update
// files during spectr initialization. Each initializer is responsible for
// a single file or directory.
//
// Initializers are designed to be:
//   - Idempotent: safe to run multiple times
//   - Composable: multiple initializers can be combined
//   - Testable: using afero.Fs for filesystem abstraction
//   - Deduplicable: same Path() means run once
//
// Built-in initializers include:
//   - DirectoryInitializer: creates directories
//   - ConfigFileInitializer: creates/updates instruction files with markers
//   - SlashCommandsInitializer: creates slash command files
type Initializer interface {
	// Init creates or updates files managed by this initializer.
	//
	// Parameters:
	//   - ctx: context for cancellation and timeouts
	//   - fs: filesystem abstraction (project or global based on IsGlobal)
	//   - cfg: configuration containing spectr directory paths
	//   - tm: template manager for rendering templates
	//
	// Returns an error if initialization fails. Implementations must be
	// idempotent - calling Init multiple times produces the same result.
	Init(
		ctx context.Context,
		fs afero.Fs,
		cfg *Config,
		tm TemplateRenderer,
	) error

	// IsSetup returns true if this initializer's artifacts already exist.
	//
	// This method is used to determine if initialization has already been
	// performed. It should check for the existence of files or directories
	// that Init would create.
	//
	// Parameters:
	//   - fs: filesystem abstraction to check
	//   - cfg: configuration containing spectr directory paths
	//
	// Returns true if artifacts exist, false otherwise.
	IsSetup(fs afero.Fs, cfg *Config) bool

	// Path returns the file or directory path this initializer manages.
	//
	// This path is used for deduplication: when multiple providers return
	// initializers with the same Path(), only the first one is executed.
	//
	// The path should be relative to the filesystem root (project root for
	// project-relative initializers, home directory for global initializers).
	Path() string

	// IsGlobal returns true if this initializer uses globalFs instead of projectFs.
	//
	// Global initializers operate on paths relative to the user's home directory
	// (e.g., ~/.config/tool/commands/). Project initializers operate on paths
	// relative to the project root (e.g., .claude/commands/).
	//
	// Returns true for global paths, false for project-relative paths.
	IsGlobal() bool
}
