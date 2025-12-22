// Package initializers provides initialization logic for various providers.
package initializers

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/connerohnesorge/spectr/internal/initialize/types"
	"github.com/spf13/afero"
)

const (
	dirPerm  = 0755
	filePerm = 0644
)

type TOMLCommandInitializer struct {
	cmd         string
	path        string
	description string
}

func NewTOMLCommandInitializer(
	cmd, path, description string,
) *TOMLCommandInitializer {
	return &TOMLCommandInitializer{
		cmd:         cmd,
		path:        path,
		description: description,
	}
}

//nolint:revive // argument-limit - interface defined elsewhere
func (t *TOMLCommandInitializer) Init(
	ctx context.Context,
	projectFs, globalFs afero.Fs,
	cfg *types.Config,
	tm types.TemplateRenderer,
) error {
	fs := projectFs
	if IsGlobalPath(t.path) {
		fs = globalFs
	}

	prompt, err := tm.RenderSlashCommand(
		t.cmd,
		types.DefaultTemplateContext(),
	)
	if err != nil {
		return fmt.Errorf(
			"failed to render slash command %s: %w",
			t.cmd,
			err,
		)
	}

	content := generateTOMLContent(
		t.description,
		prompt,
	)

	dir := filepath.Dir(t.path)
	if err := fs.MkdirAll(dir, dirPerm); err != nil {
		return fmt.Errorf(
			"failed to create directory for %s: %w",
			t.path,
			err,
		)
	}

	return afero.WriteFile(
		fs,
		t.path,
		[]byte(content),
		filePerm,
	)
}

func (t *TOMLCommandInitializer) IsSetup(
	projectFs, globalFs afero.Fs,
	cfg *types.Config,
) (bool, error) {
	fs := projectFs
	if IsGlobalPath(t.path) {
		fs = globalFs
	}
	return afero.Exists(fs, t.path)
}

func (t *TOMLCommandInitializer) Path() string {
	return t.path
}

func generateTOMLContent(
	description, prompt string,
) string {
	// Escape the prompt for TOML multiline string
	escapedPrompt := strings.ReplaceAll(
		prompt,
		`\`,
		`\\`,
	)
	escapedPrompt = strings.ReplaceAll(
		escapedPrompt,
		`"`,
		`\"`,
	)

	return fmt.Sprintf(
		`# Spectr command for Gemini CLI
description = "%s"
prompt = """
%s
"""
`,
		description,
		escapedPrompt,
	)
}