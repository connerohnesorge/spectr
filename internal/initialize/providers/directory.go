package providers

import (
	"context"
	"fmt"

	"github.com/spf13/afero"
)

// DirectoryInitializer creates one or more directories.
// Implements the Initializer interface for directory creation.
//
// Example usage:
//
//	init := NewDirectoryInitializer(
//	    ".claude/commands/spectr",
//	    ".claude/contexts",
//	)
type DirectoryInitializer struct {
	Paths      []string
	IsGlobalFs bool
}

// NewDirectoryInitializer creates a new DirectoryInitializer
// for the given paths.
// All directories will be created with 0755 permissions.
//
// Parameters:
//   - paths: One or more directory paths to create
//
// Returns:
//   - *DirectoryInitializer: A new directory initializer
func NewDirectoryInitializer(paths ...string) *DirectoryInitializer {
	return &DirectoryInitializer{
		Paths:      paths,
		IsGlobalFs: false,
	}
}

// WithGlobal configures the initializer to use the global filesystem.
// This is useful for tools that install commands in ~/.config/
// directories.
func (d *DirectoryInitializer) WithGlobal(global bool) *DirectoryInitializer {
	d.IsGlobalFs = global

	return d
}

// Init creates all directories configured in this initializer.
// Directories are created with 0755 permissions using MkdirAll
// (creates parent directories).
//
// Parameters:
//   - ctx: Context for cancellation
//   - fs: Filesystem abstraction (project-relative or global)
//   - cfg: Configuration (not used for directory creation)
//   - tm: TemplateManager (not used for directory creation)
//
// Returns:
//   - InitResult: Contains list of created directories
//   - error: Non-nil if directory creation fails
func (d *DirectoryInitializer) Init(
	_ context.Context,
	fs afero.Fs,
	_ *Config,
	_ TemplateManager,
) (InitResult, error) {
	var result InitResult

	for _, path := range d.Paths {
		// Check if directory already exists
		exists, err := afero.DirExists(fs, path)
		if err != nil {
			return InitResult{}, fmt.Errorf(
				"failed to check if directory exists %s: %w",
				path,
				err,
			)
		}

		if !exists {
			// Create directory with all parent directories
			if err := fs.MkdirAll(path, dirPerm); err != nil {
				return InitResult{}, fmt.Errorf(
					"failed to create directory %s: %w",
					path,
					err,
				)
			}
			result.CreatedFiles = append(result.CreatedFiles, path)
		}
	}

	return result, nil
}

// IsSetup returns true if ALL configured directories exist.
// If any directory is missing, returns false.
//
// Parameters:
//   - fs: Filesystem abstraction
//   - cfg: Configuration (not used)
//
// Returns:
//   - bool: True if all directories exist
func (d *DirectoryInitializer) IsSetup(fs afero.Fs, _ *Config) bool {
	for _, path := range d.Paths {
		exists, err := afero.DirExists(fs, path)
		if err != nil || !exists {
			return false
		}
	}

	return true
}

// Path returns the first directory path for deduplication.
// When multiple providers share the same directory, only one
// initializer runs.
//
// Returns:
//   - string: The first directory path, or empty string if no
//     paths configured
func (d *DirectoryInitializer) Path() string {
	if len(d.Paths) == 0 {
		return ""
	}

	return d.Paths[0]
}

// IsGlobal returns true if this initializer uses the global filesystem.
//
// Returns:
//   - bool: True if using global filesystem, false for project-relative
func (d *DirectoryInitializer) IsGlobal() bool {
	return d.IsGlobalFs
}
