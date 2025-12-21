package initializers

import (
	"context"
	"fmt"

	"github.com/spf13/afero"
	"github.com/connerohnesorge/spectr/internal/initialize/types"
)

// ConfigFileInitializer creates or updates an instruction file (e.g. CLAUDE.md)
type ConfigFileInitializer struct {
	path string
}

func NewConfigFileInitializer(path string) *ConfigFileInitializer {
	return &ConfigFileInitializer{path: path}
}

func (c *ConfigFileInitializer) Init(ctx context.Context, projectFs, globalFs afero.Fs, cfg *types.Config, tm types.TemplateRenderer) error {
	fs := projectFs
    isGlobal := IsGlobalPath(c.path)
	if isGlobal {
		fs = globalFs
	}

	content, err := tm.RenderInstructionPointer(types.DefaultTemplateContext())
	if err != nil {
		return fmt.Errorf("failed to render instruction pointer: %w", err)
	}

    targetPath := c.path
    if isGlobal {
        targetPath = ExpandPath(c.path)
    }

	return UpdateFileWithMarkers(
		fs,
		targetPath,
		content,
		types.SpectrStartMarker,
		types.SpectrEndMarker,
	)
}

func (c *ConfigFileInitializer) IsSetup(projectFs, globalFs afero.Fs, cfg *types.Config) (bool, error) {
	fs := projectFs
    targetPath := c.path
	if IsGlobalPath(c.path) {
		fs = globalFs
        targetPath = ExpandPath(c.path)
	}
	exists, err := afero.Exists(fs, targetPath)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (c *ConfigFileInitializer) Path() string {
	return c.path
}
