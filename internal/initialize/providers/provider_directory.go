package providers

import (
	"context"

	"github.com/spf13/afero"
)

// DirectoryInitializer creates directories.
// It is idempotent (safe to run multiple times) and returns the directory path
// in CreatedFiles only if the directory was created by this call.
type DirectoryInitializer struct {
	path     string
	isGlobal bool
}

// NewDirectoryInitializer creates a DirectoryInitializer for a
// project-relative directory.
// The path should be relative to the project root
// (e.g., ".claude/commands/spectr").
func NewDirectoryInitializer(
	path string,
) *DirectoryInitializer {
	return &DirectoryInitializer{
		path:     path,
		isGlobal: false,
	}
}

// NewGlobalDirectoryInitializer creates a DirectoryInitializer for a
// global directory.
// The path should be relative to the user's home directory
// (e.g., ".config/aider/commands").
func NewGlobalDirectoryInitializer(
	path string,
) *DirectoryInitializer {
	return &DirectoryInitializer{
		path:     path,
		isGlobal: true,
	}
}

// Init creates the directory if it doesn't exist.
// Returns the directory path in CreatedFiles if it was created.
func (d *DirectoryInitializer) Init(
	_ context.Context,
	fs afero.Fs,
	_ *Config,
	_ TemplateManager,
) (InitResult, error) {
	// Check if directory already exists
	exists, err := afero.DirExists(fs, d.path)
	if err != nil {
		return InitResult{}, err
	}

	// Directory already exists, nothing to do
	if exists {
		return InitResult{}, nil
	}

	// Create directory and all parent directories
	if err := fs.MkdirAll(d.path, dirPerm); err != nil {
		return InitResult{}, err
	}

	// Return result indicating directory was created
	return InitResult{
		CreatedFiles: []string{d.path},
	}, nil
}

// IsSetup returns true if the directory exists.
func (d *DirectoryInitializer) IsSetup(
	fs afero.Fs,
	_ *Config,
) bool {
	exists, err := afero.DirExists(fs, d.path)
	if err != nil {
		return false
	}

	return exists
}

// Path returns the directory path this initializer manages.
func (d *DirectoryInitializer) Path() string {
	return d.path
}

// IsGlobal returns true if this initializer uses globalFs instead of projectFs.
func (d *DirectoryInitializer) IsGlobal() bool {
	return d.isGlobal
}
