package initializers

import (
	"context"

	"github.com/spf13/afero"
	"github.com/connerohnesorge/spectr/internal/initialize/types"
)

type DirectoryInitializer struct {
	path string
}

func NewDirectoryInitializer(path string) *DirectoryInitializer {
	return &DirectoryInitializer{path: path}
}

func (d *DirectoryInitializer) Init(ctx context.Context, projectFs, globalFs afero.Fs, cfg *types.Config, tm types.TemplateRenderer) error {
	fs := projectFs
	if IsGlobalPath(d.path) {
		fs = globalFs
	}
	return fs.MkdirAll(d.path, 0755)
}

func (d *DirectoryInitializer) IsSetup(projectFs, globalFs afero.Fs, cfg *types.Config) (bool, error) {
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

func (d *DirectoryInitializer) Path() string {
	return d.path
}