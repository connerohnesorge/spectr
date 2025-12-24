package initializers

import (
	"context"

	"github.com/spf13/afero"
)

// DirectoryInitializer creates one or more directories.
// It is idempotent and will not fail if directories already exist.
type DirectoryInitializer struct {
	paths []string
}

// NewDirectoryInitializer creates a DirectoryInitializer for the given paths.
// All paths should be relative to the filesystem root.
func NewDirectoryInitializer(
	paths ...string,
) *DirectoryInitializer {
	return &DirectoryInitializer{
		paths: paths,
	}
}

// Init creates all directories specified in the initializer.
// Returns InitResult with the list of created directories.
// If a directory already exists, it is not included in the result.
func (d *DirectoryInitializer) Init(
	_ context.Context,
	fs afero.Fs,
	_ *Config,
	_ any,
) (InitResult, error) {
	var result InitResult

	for _, path := range d.paths {
		// Check if directory already exists
		exists, err := afero.DirExists(fs, path)
		if err != nil {
			return result, err
		}

		if !exists {
			// Create directory with all parent directories
			if err := fs.MkdirAll(path, DirPerm); err != nil {
				return result, err
			}

			result.CreatedFiles = append(
				result.CreatedFiles,
				path,
			)
		}
	}

	return result, nil
}

// IsSetup returns true if all directories exist.
func (d *DirectoryInitializer) IsSetup(
	fs afero.Fs,
	_ *Config,
) bool {
	for _, path := range d.paths {
		exists, err := afero.DirExists(fs, path)
		if err != nil || !exists {
			return false
		}
	}

	return true
}

// Path returns the first directory path.
// Used for deduplication - if multiple initializers manage the same path,
// only one will be executed.
func (d *DirectoryInitializer) Path() string {
	if len(d.paths) > 0 {
		return d.paths[0]
	}

	return ""
}

// IsGlobal returns false for directory initializers.
func (*DirectoryInitializer) IsGlobal() bool {
	return false
}
