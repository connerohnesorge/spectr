package initializers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/initialize/types"
	"github.com/spf13/afero"
)

// DirectoryInitializer initializes a directory.
type DirectoryInitializer struct {
	path string
}

// NewDirectoryInitializer creates a new DirectoryInitializer.
func NewDirectoryInitializer(
	path string,
) *DirectoryInitializer {
	return &DirectoryInitializer{path: path}
}

// Init initializes the directory.
//
//nolint:revive // argument-limit - interface defined elsewhere
func (d *DirectoryInitializer) Init(
	_ context.Context,
	projectFs, globalFs afero.Fs,
	_ *types.Config,
	_ types.TemplateRenderer,
) error {
	fs := projectFs
	if IsGlobalPath(d.path) {
		fs = globalFs
	}

	return fs.MkdirAll(d.path, dirPerm)
}

// IsSetup checks if the directory is already set up.
func (d *DirectoryInitializer) IsSetup(
	projectFs, globalFs afero.Fs,
	_ *types.Config,
) (bool, error) {
	fs := projectFs
	if IsGlobalPath(d.path) {
		fs = globalFs
	}

	exists, err := afero.DirExists(fs, d.path)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// Path returns the path of the directory.
func (d *DirectoryInitializer) Path() string {
	return d.path
}
