// Package initializers provides built-in Initializer implementations for the
// provider architecture.
//
// This file implements DirectoryInitializer, which creates directories needed
// by providers (e.g., .claude/commands/spectr/).
//
//nolint:revive // line-length-limit, unused-parameter - interface compliance
package initializers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/initialize/providers"
	"github.com/spf13/afero"
)

// DirectoryInitializer creates directories needed by providers.
//
// It implements the providers.Initializer interface and is designed to:
//   - Accept one or more directory paths
//   - Create all parent directories as needed (like mkdir -p)
//   - Be idempotent (safe to run multiple times)
//   - Work with afero.Fs for testability
type DirectoryInitializer struct {
	// paths contains the directories to create. The first path is used
	// as the primary path for deduplication via Path().
	paths []string

	// isGlobal indicates whether this initializer operates on global paths
	// (relative to home directory) instead of project-relative paths.
	isGlobal bool
}

// NewDirectoryInitializer creates a new DirectoryInitializer that will create
// the specified directories.
//
// Parameters:
//   - isGlobal: if true, paths are relative to home directory; otherwise project-relative
//   - paths: one or more directory paths to create
//
// Returns nil if no paths are provided.
func NewDirectoryInitializer(isGlobal bool, paths ...string) *DirectoryInitializer {
	if len(paths) == 0 {
		return nil
	}

	return &DirectoryInitializer{
		paths:    paths,
		isGlobal: isGlobal,
	}
}

// Init creates all directories specified in the initializer.
//
// It uses MkdirAll to create parent directories as needed, similar to mkdir -p.
// This operation is idempotent - running it multiple times has the same effect
// as running it once.
//
// Parameters:
//   - ctx: context for cancellation (not currently used but part of interface)
//   - fs: filesystem abstraction to create directories on
//   - cfg: configuration (not currently used but part of interface)
//   - tm: template manager (not used for directory creation)
//
// Returns an error if any directory creation fails.
func (d *DirectoryInitializer) Init(
	ctx context.Context,
	fs afero.Fs,
	cfg *providers.Config,
	tm providers.TemplateRenderer,
) error {
	for _, path := range d.paths {
		// Use 0755 permissions (rwxr-xr-x) for directories
		// MkdirAll creates all parent directories as needed
		if err := fs.MkdirAll(path, 0755); err != nil {
			return err
		}
	}

	return nil
}

// IsSetup returns true if all directories managed by this initializer exist.
//
// Parameters:
//   - fs: filesystem abstraction to check
//   - cfg: configuration (not currently used but part of interface)
//
// Returns true only if ALL directories exist; false if any are missing.
func (d *DirectoryInitializer) IsSetup(fs afero.Fs, cfg *providers.Config) bool {
	for _, path := range d.paths {
		info, err := fs.Stat(path)
		if err != nil {
			// Directory doesn't exist or error accessing it
			return false
		}
		if !info.IsDir() {
			// Path exists but is not a directory
			return false
		}
	}

	return true
}

// Path returns the primary directory path this initializer manages.
//
// This is used for deduplication: when multiple providers return
// DirectoryInitializers with the same Path(), only the first one is executed.
//
// Returns the first path in the list of directories to create.
func (d *DirectoryInitializer) Path() string {
	if len(d.paths) == 0 {
		return ""
	}

	return d.paths[0]
}

// IsGlobal returns true if this initializer uses globalFs instead of projectFs.
//
// Global initializers operate on paths relative to the user's home directory
// (e.g., ~/.config/tool/commands/). Project initializers operate on paths
// relative to the project root (e.g., .claude/commands/).
func (d *DirectoryInitializer) IsGlobal() bool {
	return d.isGlobal
}

// Paths returns all directory paths this initializer will create.
// This is useful for testing and inspection.
func (d *DirectoryInitializer) Paths() []string {
	return d.paths
}

// Ensure DirectoryInitializer implements the Initializer interface.
var _ providers.Initializer = (*DirectoryInitializer)(nil)
