// Package initializers provides built-in initializers for the provider system.
// These initializers handle common initialization tasks like creating
// directories, writing config files, and setting up slash commands.
package initializers

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/connerohnesorge/spectr/internal/initialize/providers"
	"github.com/spf13/afero"
)

// Compile-time interface satisfaction check.
var _ providers.Initializer = (*DirectoryInitializer)(
	nil,
)

// DirectoryInitializer creates directories.
// It is idempotent - running Init multiple times produces the same result.
type DirectoryInitializer struct {
	// Paths are the directory paths to create, relative to the project root.
	Paths []string
}

// NewDirectoryInitializer creates a new DirectoryInitializer for the given
// paths. Paths should be relative to the project root.
func NewDirectoryInitializer(
	paths ...string,
) *DirectoryInitializer {
	return &DirectoryInitializer{
		Paths: paths,
	}
}

// Init creates all directories specified in Paths.
// It uses fs.MkdirAll to create directories and any necessary parents.
// Returns nil if all directories are created successfully.
func (d *DirectoryInitializer) Init(
	_ context.Context,
	fs afero.Fs,
	_ *providers.Config,
) error {
	for _, path := range d.Paths {
		if err := fs.MkdirAll(path, dirPerm); err != nil {
			return fmt.Errorf(
				"failed to create directory %s: %w",
				path,
				err,
			)
		}
	}

	return nil
}

// IsSetup returns true if all directories in Paths exist.
// Returns false if any directory is missing.
func (d *DirectoryInitializer) IsSetup(
	fs afero.Fs,
	_ *providers.Config,
) bool {
	for _, path := range d.Paths {
		info, err := fs.Stat(path)
		if err != nil {
			return false
		}
		if !info.IsDir() {
			return false
		}
	}

	return true
}

// Key returns a unique key for this initializer based on its
// configuration. This is used for deduplication when multiple providers
// use the same initializer. The key format is "dir:" followed by sorted,
// comma-separated paths. Example: "dir:.claude/commands/spectr"
func (d *DirectoryInitializer) Key() string {
	if len(d.Paths) == 0 {
		return "dir:"
	}

	// Create a sorted copy to ensure consistent keys
	sortedPaths := make([]string, len(d.Paths))
	copy(sortedPaths, d.Paths)
	sort.Strings(sortedPaths)

	return "dir:" + strings.Join(sortedPaths, ",")
}
