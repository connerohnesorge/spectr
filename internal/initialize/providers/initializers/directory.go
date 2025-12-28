// Package initializers provides built-in initializer implementations for provider configuration.
package initializers

import (
	"context"
	"path/filepath"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/spf13/afero"
)

// DirectoryInitializer creates directories in the project filesystem.
type DirectoryInitializer struct {
	paths []string
}

// NewDirectoryInitializer creates an initializer that creates the specified directories
// in the project filesystem.
func NewDirectoryInitializer(paths ...string) domain.Initializer {
	return &DirectoryInitializer{paths: paths}
}

// Directory permission constant.
const dirPerm = 0o755

// Init creates directories in the project filesystem.
// Uses MkdirAll for recursive creation, silent success if directory exists.
//
//nolint:revive // argument-limit - interface signature requires these parameters
func (d *DirectoryInitializer) Init(
	_ context.Context,
	projectFs, _ afero.Fs,
	_ *domain.Config,
	_ any,
) (domain.ExecutionResult, error) {
	var created []string

	for _, path := range d.paths {
		// Check if directory already exists
		exists, err := afero.DirExists(projectFs, path)
		if err != nil {
			return domain.ExecutionResult{CreatedFiles: created}, err
		}

		if !exists {
			// Create directory with parents
			if err := projectFs.MkdirAll(path, dirPerm); err != nil {
				return domain.ExecutionResult{CreatedFiles: created}, err
			}
			created = append(created, path)
		}
		// Silent success if directory already exists
	}

	return domain.ExecutionResult{CreatedFiles: created}, nil
}

// IsSetup returns true if all directories already exist.
func (d *DirectoryInitializer) IsSetup(projectFs, _ afero.Fs, _ *domain.Config) bool {
	for _, path := range d.paths {
		exists, err := afero.DirExists(projectFs, path)
		if err != nil || !exists {
			return false
		}
	}

	return true
}

// DedupeKey returns a unique key for deduplication.
// Uses the first path with filepath.Clean for normalization.
// Exported to allow deduplication from the executor package.
func (d *DirectoryInitializer) DedupeKey() string {
	if len(d.paths) == 0 {
		return "DirectoryInitializer:"
	}

	return "DirectoryInitializer:" + filepath.Clean(d.paths[0])
}

// Ensure DirectoryInitializer implements the Deduplicatable interface.
var _ Deduplicatable = (*DirectoryInitializer)(nil)

// HomeDirectoryInitializer creates directories in the home filesystem.
type HomeDirectoryInitializer struct {
	paths []string
}

// NewHomeDirectoryInitializer creates an initializer that creates the specified directories
// in the home filesystem (e.g., ~/.config/tool/).
func NewHomeDirectoryInitializer(paths ...string) domain.Initializer {
	return &HomeDirectoryInitializer{paths: paths}
}

// Init creates directories in the home filesystem.
// Uses MkdirAll for recursive creation, silent success if directory exists.
//
//nolint:revive // argument-limit - interface signature requires these parameters
func (d *HomeDirectoryInitializer) Init(
	_ context.Context,
	_, homeFs afero.Fs,
	_ *domain.Config,
	_ any,
) (domain.ExecutionResult, error) {
	var created []string

	for _, path := range d.paths {
		// Check if directory already exists
		exists, err := afero.DirExists(homeFs, path)
		if err != nil {
			return domain.ExecutionResult{CreatedFiles: created}, err
		}

		if !exists {
			// Create directory with parents
			if err := homeFs.MkdirAll(path, dirPerm); err != nil {
				return domain.ExecutionResult{CreatedFiles: created}, err
			}
			created = append(created, path)
		}
		// Silent success if directory already exists
	}

	return domain.ExecutionResult{CreatedFiles: created}, nil
}

// IsSetup returns true if all directories already exist.
func (d *HomeDirectoryInitializer) IsSetup(_, homeFs afero.Fs, _ *domain.Config) bool {
	for _, path := range d.paths {
		exists, err := afero.DirExists(homeFs, path)
		if err != nil || !exists {
			return false
		}
	}

	return true
}

// DedupeKey returns a unique key for deduplication.
// Uses the first path with filepath.Clean for normalization.
// Exported to allow deduplication from the executor package.
func (d *HomeDirectoryInitializer) DedupeKey() string {
	if len(d.paths) == 0 {
		return "HomeDirectoryInitializer:"
	}

	return "HomeDirectoryInitializer:" + filepath.Clean(d.paths[0])
}

// Ensure HomeDirectoryInitializer implements the Deduplicatable interface.
var _ Deduplicatable = (*HomeDirectoryInitializer)(nil)

// Deduplicatable is an optional interface for initializers that support deduplication.
// Initializers that implement this interface can be deduplicated based on their key.
// Exported to allow deduplication from the executor package.
type Deduplicatable interface {
	// DedupeKey returns a unique key for deduplication.
	DedupeKey() string
}
