package initializers

import (
	"context"
	"fmt"

	"github.com/connerohnesorge/spectr/internal/initialize/types"
	"github.com/spf13/afero"
)

// SlashCommandsInitializer initializes slash commands.
type SlashCommandsInitializer struct {
	cmd         string // "proposal", "apply"
	path        string
	frontmatter string
}

// NewSlashCommandsInitializer creates a new SlashCommandsInitializer.
func NewSlashCommandsInitializer(
	cmd, path, frontmatter string,
) *SlashCommandsInitializer {
	return &SlashCommandsInitializer{
		cmd:         cmd,
		path:        path,
		frontmatter: frontmatter,
	}
}

// Init initializes the slash command.
//
//nolint:revive // argument-limit - interface defined elsewhere
func (s *SlashCommandsInitializer) Init(
	_ context.Context,
	projectFs, globalFs afero.Fs,
	_ *types.Config,
	tm types.TemplateRenderer,
) error {
	fs := projectFs
	isGlobal := IsGlobalPath(s.path)
	if isGlobal {
		fs = globalFs
	}

	body, err := tm.RenderSlashCommand(
		s.cmd,
		types.DefaultTemplateContext(),
	)
	if err != nil {
		return fmt.Errorf(
			"failed to render slash command %s: %w",
			s.cmd,
			err,
		)
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
		return updateSlashCommandBody(
			fs,
			targetPath,
			body,
			s.frontmatter,
		)
	}

	return createNewSlashCommand(
		fs,
		targetPath,
		body,
		s.frontmatter,
	)
}

// IsSetup checks if the slash command is already set up.
func (s *SlashCommandsInitializer) IsSetup(
	projectFs, globalFs afero.Fs,
	_ *types.Config,
) (bool, error) {
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

// Path returns the path of the slash command.
func (s *SlashCommandsInitializer) Path() string {
	return s.path
}