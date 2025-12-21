package initializers

import (
	"context"
	"fmt"

	"github.com/spf13/afero"
	"github.com/connerohnesorge/spectr/internal/initialize/types"
)

type SlashCommandsInitializer struct {
	cmd         string // "proposal", "apply"
	path        string
	frontmatter string
}

func NewSlashCommandsInitializer(cmd, path, frontmatter string) *SlashCommandsInitializer {
	return &SlashCommandsInitializer{
		cmd:         cmd,
		path:        path,
		frontmatter: frontmatter,
	}
}

func (s *SlashCommandsInitializer) Init(ctx context.Context, projectFs, globalFs afero.Fs, cfg *types.Config, tm types.TemplateRenderer) error {
	fs := projectFs
	isGlobal := IsGlobalPath(s.path)
	if isGlobal {
		fs = globalFs
	}

	body, err := tm.RenderSlashCommand(s.cmd, types.DefaultTemplateContext())
	if err != nil {
		return fmt.Errorf("failed to render slash command %s: %w", s.cmd, err)
	}

    // Use ExpandPath only if global, otherwise keep relative for BasePathFs
    targetPath := s.path
    if isGlobal {
        targetPath = ExpandPath(s.path)
    }

	exists, err := afero.Exists(fs, targetPath)
	if err != nil {
		return err
	}

	if exists {
		return updateSlashCommandBody(fs, targetPath, body, s.frontmatter)
	}

	return createNewSlashCommand(fs, targetPath, s.cmd, body, s.frontmatter)
}

func (s *SlashCommandsInitializer) IsSetup(projectFs, globalFs afero.Fs, cfg *types.Config) (bool, error) {
	fs := projectFs
    targetPath := s.path
	if IsGlobalPath(s.path) {
		fs = globalFs
        targetPath = ExpandPath(s.path)
	}
	exists, err := afero.Exists(fs, targetPath)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *SlashCommandsInitializer) Path() string {
	return s.path
}
