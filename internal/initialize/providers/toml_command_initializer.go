package providers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// TOMLSlashCommandInitializer handles TOML slash command files (Gemini CLI).
// These files use TOML format with description and prompt fields.
type TOMLSlashCommandInitializer struct {
	path        string // Relative path to command file
	commandName string // Command name for rendering (e.g., "proposal")
	description string // Description for the TOML file
}

// NewTOMLSlashCommandInitializer creates a new TOML slash command initializer.
// path is the relative path to the command file.
// commandName is used to render the template (e.g., "proposal", "apply").
// description is the description field in the TOML file.
func NewTOMLSlashCommandInitializer(
	path, commandName, description string,
) *TOMLSlashCommandInitializer {
	return &TOMLSlashCommandInitializer{
		path:        path,
		commandName: commandName,
		description: description,
	}
}

// ID returns the unique identifier for this initializer.
// Format: "toml-cmd:{path}"
func (i *TOMLSlashCommandInitializer) ID() string {
	return "toml-cmd:" + i.path
}

// FilePath returns the relative path this initializer manages.
func (i *TOMLSlashCommandInitializer) FilePath() string {
	return i.path
}

// Configure creates or updates the TOML slash command file.
// Note: TOML files are always overwritten (no marker-based updates).
func (i *TOMLSlashCommandInitializer) Configure(
	projectPath string,
	tm TemplateRenderer,
) error {
	prompt, err := tm.RenderSlashCommand(
		i.commandName,
		DefaultTemplateContext(),
	)
	if err != nil {
		return fmt.Errorf(
			"failed to render slash command %s: %w",
			i.commandName,
			err,
		)
	}

	fullPath := i.expandedPath(projectPath)
	content := i.generateTOMLContent(prompt)

	dir := filepath.Dir(fullPath)
	err = EnsureDir(dir)
	if err != nil {
		return fmt.Errorf(
			"failed to create directory for %s: %w",
			fullPath,
			err,
		)
	}

	err = os.WriteFile(
		fullPath,
		[]byte(content),
		filePerm,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to write TOML command file %s: %w",
			fullPath,
			err,
		)
	}

	return nil
}

// IsConfigured checks if the TOML command file exists.
func (i *TOMLSlashCommandInitializer) IsConfigured(
	projectPath string,
) bool {
	fullPath := i.expandedPath(projectPath)

	return FileExists(fullPath)
}

// expandedPath returns the full path for the command file,
// handling ~ expansion for global paths.
func (i *TOMLSlashCommandInitializer) expandedPath(
	projectPath string,
) string {
	if isGlobalPath(i.path) {
		return expandPath(i.path)
	}

	return filepath.Join(projectPath, i.path)
}

// generateTOMLContent creates TOML content for a command file.
func (i *TOMLSlashCommandInitializer) generateTOMLContent(
	prompt string,
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
		i.description,
		escapedPrompt,
	)
}
