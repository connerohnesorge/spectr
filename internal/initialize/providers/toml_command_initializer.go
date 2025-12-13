// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
package providers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// TOMLSlashCommandInitializer manages TOML-based slash command files.
//
// These files are used by Gemini CLI which uses TOML format for command
// definitions. Unlike markdown commands, TOML files are always completely
// overwritten (no marker-based updates).
//
// TOML format:
//
//	# Spectr command for Gemini CLI
//	description = "{description}"
//	prompt = """
//	{escapedPrompt}
//	"""
type TOMLSlashCommandInitializer struct {
	// path is the file path (e.g., ".gemini/commands/spectr/proposal.toml")
	path string
	// commandName is the command name (e.g., "proposal", "apply")
	commandName string
	// description is the command description for TOML metadata
	description string
}

// NewTOMLSlashCommandInitializer creates a new TOMLSlashCommandInitializer.
//
// Parameters:
//   - path: The file path (e.g., ".gemini/commands/spectr/proposal.toml")
//   - commandName: The command name (e.g., "proposal", "apply")
//   - description: The command description for TOML metadata
//
// Example usage:
//
//	NewTOMLSlashCommandInitializer(
//	    ".gemini/commands/spectr/proposal.toml",
//	    "proposal",
//	    "Scaffold a new Spectr change and validate strictly.",
//	)
func NewTOMLSlashCommandInitializer(
	path, commandName, description string,
) *TOMLSlashCommandInitializer {
	return &TOMLSlashCommandInitializer{
		path:        path,
		commandName: commandName,
		description: description,
	}
}

// ID returns a unique identifier for this initializer.
//
// Format: "toml-cmd:{path}"
// e.g., "toml-cmd:.gemini/commands/spectr/proposal.toml"
func (t *TOMLSlashCommandInitializer) ID() string {
	return "toml-cmd:" + t.path
}

// FilePath returns the path this initializer manages.
//
// The returned path may contain ~ for home directory paths.
// Path expansion is handled internally during Configure and IsConfigured.
func (t *TOMLSlashCommandInitializer) FilePath() string {
	return t.path
}

// Configure creates or updates the TOML slash command file.
//
// For project-relative paths, the file is created at
// filepath.Join(projectPath, path). For global paths (starting with ~/ or /),
// the path is expanded independently.
//
// The file content is rendered using TemplateRenderer.RenderSlashCommand,
// which generates the command prompt content. The prompt is then wrapped in
// TOML format with description and proper escaping.
//
// Unlike markdown commands, TOML files are always completely overwritten
// (no marker-based updates).
func (t *TOMLSlashCommandInitializer) Configure(
	projectPath string,
	tm TemplateRenderer,
) error {
	// Render the slash command content
	prompt, err := tm.RenderSlashCommand(
		t.commandName,
		DefaultTemplateContext(),
	)
	if err != nil {
		return fmt.Errorf(
			"failed to render slash command %s: %w",
			t.commandName,
			err,
		)
	}

	// Determine the full file path
	fullPath := t.resolvePath(projectPath)

	// Generate TOML content
	content := t.generateTOMLContent(prompt)

	// Ensure directory exists
	dir := filepath.Dir(fullPath)
	err = EnsureDir(dir)
	if err != nil {
		return fmt.Errorf(
			"failed to create directory for %s: %w",
			fullPath,
			err,
		)
	}

	// Write the file (always overwrite)
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

// IsConfigured checks if the TOML slash command file exists.
//
// Path resolution follows the same rules as Configure:
//   - Project-relative paths are joined with projectPath
//   - Global paths are expanded independently
func (t *TOMLSlashCommandInitializer) IsConfigured(
	projectPath string,
) bool {
	fullPath := t.resolvePath(projectPath)

	return FileExists(fullPath)
}

// resolvePath returns the full path for the TOML command file.
//
// For global paths (starting with ~/ or /), the path is expanded.
// For project-relative paths, the path is joined with projectPath.
func (t *TOMLSlashCommandInitializer) resolvePath(
	projectPath string,
) string {
	if isGlobalPath(t.path) {
		return expandPath(t.path)
	}

	return filepath.Join(projectPath, t.path)
}

// generateTOMLContent creates TOML content for the command.
//
// The prompt is escaped for TOML multiline strings:
//   - Backslashes are doubled (\ -> \\)
//   - Double quotes are escaped (" -> \")
func (t *TOMLSlashCommandInitializer) generateTOMLContent(
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
		t.description,
		escapedPrompt,
	)
}
