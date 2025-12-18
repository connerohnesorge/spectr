// Package initializers provides built-in initializers for the provider system.
//
// This package contains reusable initializers that providers can compose to set up
// their configuration. The three main initializers are:
//   - DirectoryInitializer: Creates directories (e.g., .claude/commands/spectr/)
//   - ConfigFileInitializer: Creates/updates instruction files with markers
//   - SlashCommandsInitializer: Creates slash commands from templates
//
// These initializers implement the providers.Initializer interface and can be
// combined in any provider's Initializers() method.
package initializers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/initialize/providers"
	"github.com/spf13/afero"
)

// DirectoryInitializer creates one or more directories.
//
// DirectoryInitializer ensures that specified directories exist, creating them
// (including parent directories) if necessary. This is typically the first
// initializer to run, as other initializers may depend on directories existing.
//
// # Execution Order
//
// DirectoryInitializer has priority 1 in the initializer ordering, meaning it
// runs before ConfigFileInitializer (priority 2) and SlashCommandsInitializer
// (priority 3). This ensures directories exist before files are written.
//
// # Idempotency
//
// DirectoryInitializer is idempotent: calling Init() multiple times has the
// same effect as calling it once. If directories already exist, no action is
// taken.
//
// # Example Usage
//
//	func (p *ClaudeProvider) Initializers(ctx context.Context) []providers.Initializer {
//	    return []providers.Initializer{
//	        initializers.NewDirectoryInitializer(".claude/commands/spectr"),
//	        // ... other initializers
//	    }
//	}
//
// # Multiple Directories
//
// You can create multiple directories with a single initializer:
//
//	init := initializers.NewDirectoryInitializer(
//	    ".claude/commands/spectr",
//	    ".claude/settings",
//	)
//
// Note: When using multiple paths, only the first path is used for deduplication
// via the Path() method. Consider using separate DirectoryInitializer instances
// if deduplication for each path is important.
//
// # Global Directories
//
// For directories in the user's home directory (e.g., ~/.config/aider/commands/),
// use NewGlobalDirectoryInitializer instead:
//
//	init := initializers.NewGlobalDirectoryInitializer(".config/aider/commands/spectr")
type DirectoryInitializer struct {
	// paths contains the directories to create (relative to fs root)
	paths []string
	// global indicates whether to use globalFs (true) or projectFs (false)
	global bool
}

// NewDirectoryInitializer creates a new DirectoryInitializer for project directories.
//
// The paths should be relative to the project root. Parent directories will be
// created automatically if they don't exist.
//
// Example:
//
//	init := NewDirectoryInitializer(".claude/commands/spectr")
//	init := NewDirectoryInitializer(".gemini/commands/spectr", ".gemini/settings")
func NewDirectoryInitializer(paths ...string) *DirectoryInitializer {
	return &DirectoryInitializer{
		paths:  paths,
		global: false,
	}
}

// NewGlobalDirectoryInitializer creates a new DirectoryInitializer for global directories.
//
// The paths should be relative to the user's home directory. This is used for
// tools that store configuration globally (e.g., Aider uses ~/.config/aider/).
//
// Example:
//
//	init := NewGlobalDirectoryInitializer(".config/aider/commands/spectr")
func NewGlobalDirectoryInitializer(paths ...string) *DirectoryInitializer {
	return &DirectoryInitializer{
		paths:  paths,
		global: true,
	}
}

// Init creates all directories specified in the initializer.
//
// Init uses fs.MkdirAll to create directories and all necessary parent directories.
// The directories are created with permission mode 0755 (rwxr-xr-x).
//
// Init is idempotent: if directories already exist, no error is returned and
// no action is taken.
//
// Parameters:
//   - ctx: Context for cancellation (currently unused but reserved for future use)
//   - fs: Filesystem to operate on (project or global, based on IsGlobal())
//   - cfg: Configuration (currently unused by DirectoryInitializer)
//   - tm: Template manager (currently unused by DirectoryInitializer)
//
// Returns an error if any directory creation fails.
func (d *DirectoryInitializer) Init(ctx context.Context, fs afero.Fs, cfg *providers.Config, tm providers.TemplateManager) error {
	for _, p := range d.paths {
		if err := fs.MkdirAll(p, 0755); err != nil {
			return err
		}
	}
	return nil
}

// IsSetup returns true if all directories already exist.
//
// IsSetup checks each directory path and returns true only if ALL directories
// exist. If any directory is missing, it returns false.
//
// Parameters:
//   - fs: Filesystem to check (project or global, based on IsGlobal())
//   - cfg: Configuration (currently unused by DirectoryInitializer)
func (d *DirectoryInitializer) IsSetup(fs afero.Fs, cfg *providers.Config) bool {
	for _, p := range d.paths {
		exists, err := afero.DirExists(fs, p)
		if err != nil || !exists {
			return false
		}
	}
	return true
}

// Path returns the primary path this initializer manages.
//
// Path returns the first directory path, which is used for deduplication.
// When multiple providers return DirectoryInitializers with the same Path(),
// only one will be executed.
//
// Returns an empty string if no paths were specified.
func (d *DirectoryInitializer) Path() string {
	if len(d.paths) == 0 {
		return ""
	}
	return d.paths[0]
}

// IsGlobal returns true if this initializer uses the global filesystem.
//
// Project initializers (IsGlobal() == false) operate on files within
// the project directory (e.g., .claude/commands/).
//
// Global initializers (IsGlobal() == true) operate on files in the
// user's home directory (e.g., ~/.config/aider/commands/).
func (d *DirectoryInitializer) IsGlobal() bool {
	return d.global
}

// Paths returns all directory paths this initializer manages.
//
// This is useful for debugging, logging, or when you need to know all
// directories that will be created (not just the primary path).
func (d *DirectoryInitializer) Paths() []string {
	return d.paths
}
