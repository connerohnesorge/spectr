package initializers

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"github.com/connerohnesorge/spectr/internal/initialize/types"
)

type TOMLCommandInitializer struct {
	cmd         string
	path        string
	description string
}

func NewTOMLCommandInitializer(cmd, path, description string) *TOMLCommandInitializer {
	return &TOMLCommandInitializer{
		cmd:         cmd,
		path:        path,
		description: description,
	}
}

func (t *TOMLCommandInitializer) Init(ctx context.Context, projectFs, globalFs afero.Fs, cfg *types.Config, tm types.TemplateRenderer) error {
	fs := projectFs
	if IsGlobalPath(t.path) {
		fs = globalFs
	}

	prompt, err := tm.RenderSlashCommand(t.cmd, types.DefaultTemplateContext())
	if err != nil {
		return fmt.Errorf("failed to render slash command %s: %w", t.cmd, err)
	}

	content := t.generateTOMLContent(t.description, prompt)

	dir := filepath.Dir(t.path)
	if err := fs.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", t.path, err)
	}

	return afero.WriteFile(fs, t.path, []byte(content), 0644)
}

func (t *TOMLCommandInitializer) IsSetup(projectFs, globalFs afero.Fs, cfg *types.Config) (bool, error) {
	fs := projectFs
	if IsGlobalPath(t.path) {
		fs = globalFs
	}
	return afero.Exists(fs, t.path)
}

func (t *TOMLCommandInitializer) Path() string {
	return t.path
}

func (t *TOMLCommandInitializer) generateTOMLContent(description, prompt string) string {
	// Escape the prompt for TOML multiline string
	escapedPrompt := strings.ReplaceAll(prompt, `\`, `\\`)
	escapedPrompt = strings.ReplaceAll(escapedPrompt, `"`, `"`)

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
