// Package initializers provides built-in initializers for the provider system.
package initializers

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/connerohnesorge/spectr/internal/initialize/providers"
	"github.com/spf13/afero"
)

// Directory permission constant.
const dirPerm = 0755

// Compile-time interface satisfaction check.
var _ providers.Initializer = (*SlashCommandsInitializer)(nil)

// SlashCommandsInitializer creates slash command files from templates.
// It supports both Markdown (with YAML frontmatter) and TOML formats.
// It is idempotent - running Init multiple times produces the same result.
//
// For Markdown format:
//   - Adds YAML frontmatter at the top
//   - Wraps body content with spectr markers
//
// For TOML format:
//   - Uses TOML structure with description and prompt fields
type SlashCommandsInitializer struct {
	// Dir is the directory for slash commands, relative to the project root.
	// Example: ".claude/commands/spectr", ".gemini/commands/spectr"
	Dir string

	// Extension is the file extension for command files.
	// Example: ".md", ".toml"
	Extension string

	// Format specifies the command file format (FormatMarkdown or FormatTOML).
	Format providers.CommandFormat

	// Renderer is the template renderer for command content.
	Renderer providers.TemplateRenderer
}

// NewSlashCommandsInitializer creates a new SlashCommandsInitializer.
// dir should be relative to the project root.
// ext should include the leading dot (e.g., ".md", ".toml").
// format specifies Markdown or TOML format.
// renderer is used to render template content for commands.
func NewSlashCommandsInitializer(
	dir, ext string,
	format providers.CommandFormat,
	renderer providers.TemplateRenderer,
) *SlashCommandsInitializer {
	return &SlashCommandsInitializer{
		Dir:       dir,
		Extension: ext,
		Format:    format,
		Renderer:  renderer,
	}
}

// commandInfo holds information about a slash command.
type commandInfo struct {
	name        string
	description string
}

// slashCommands returns the list of commands to create.
func slashCommands() []commandInfo {
	return []commandInfo{
		{
			name:        "proposal",
			description: "Scaffold a new Spectr change and validate strictly.",
		},
		{
			name: "apply",
			description: "Implement an approved Spectr change " +
				"and keep tasks in sync.",
		},
	}
}

// Init creates the slash command files (proposal and apply).
// It ensures the directory exists before writing.
// Returns nil if all files are created/updated successfully.
func (s *SlashCommandsInitializer) Init(
	_ context.Context,
	fs afero.Fs,
	cfg *providers.Config,
) error {
	// Ensure directory exists
	if err := fs.MkdirAll(s.Dir, dirPerm); err != nil {
		return fmt.Errorf(
			"failed to create directory %s: %w",
			s.Dir,
			err,
		)
	}

	for _, cmd := range slashCommands() {
		if err := s.createCommand(fs, cfg, cmd); err != nil {
			return err
		}
	}

	return nil
}

// createCommand creates or updates a single slash command file.
func (s *SlashCommandsInitializer) createCommand(
	fs afero.Fs,
	cfg *providers.Config,
	cmd commandInfo,
) error {
	filePath := s.commandPath(cmd.name)

	// Render the command content
	templateCtx := providers.DefaultTemplateContext()
	if cfg != nil && cfg.SpectrDir != "" {
		templateCtx.BaseDir = cfg.SpectrDir
		templateCtx.SpecsDir = cfg.SpectrDir + "/specs"
		templateCtx.ChangesDir = cfg.SpectrDir + "/changes"
		templateCtx.ProjectFile = cfg.SpectrDir + "/project.md"
		templateCtx.AgentsFile = cfg.SpectrDir + "/AGENTS.md"
	}

	body, err := s.Renderer.RenderSlashCommand(
		cmd.name,
		templateCtx,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to render slash command %s: %w",
			cmd.name,
			err,
		)
	}

	var content string
	switch s.Format {
	case providers.FormatMarkdown:
		content, err = s.generateMarkdownContent(
			fs,
			filePath,
			cmd,
			body,
		)
	case providers.FormatTOML:
		content = s.generateTOMLContent(
			cmd.description,
			body,
		)
	default:
		return fmt.Errorf(
			"unknown command format: %d",
			s.Format,
		)
	}

	if err != nil {
		return err
	}

	err = afero.WriteFile(
		fs,
		filePath,
		[]byte(content),
		filePerm,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to write command file %s: %w",
			filePath,
			err,
		)
	}

	return nil
}

// commandPath returns the full path for a command file.
func (s *SlashCommandsInitializer) commandPath(
	name string,
) string {
	return filepath.Join(s.Dir, name+s.Extension)
}

// generateMarkdownContent generates Markdown content with frontmatter.
// If the file already exists with markers, it updates content between markers.
func (s *SlashCommandsInitializer) generateMarkdownContent(
	fs afero.Fs,
	filePath string,
	cmd commandInfo,
	body string,
) (string, error) {
	frontmatter := getFrontmatter(cmd.name)

	// Check if file exists
	exists, err := afero.Exists(fs, filePath)
	if err != nil {
		return "", fmt.Errorf(
			"failed to check if file exists %s: %w",
			filePath,
			err,
		)
	}

	if exists {
		return s.updateMarkdownContent(
			fs,
			filePath,
			frontmatter,
			body,
		)
	}

	// Create new file with frontmatter and markers
	return createMarkdownContent(
		frontmatter,
		body,
	), nil
}

// createMarkdownContent creates new Markdown content with frontmatter.
func createMarkdownContent(
	frontmatter, body string,
) string {
	var sections []string

	if frontmatter != "" {
		sections = append(
			sections,
			strings.TrimSpace(frontmatter),
		)
	}

	markedBody := spectrStartMarker + newlineDouble +
		body + newlineDouble + spectrEndMarker
	sections = append(sections, markedBody)

	return strings.Join(
		sections,
		newlineDouble,
	) + newlineDouble
}

// updateMarkdownContent updates existing Markdown content.
// Replaces content between markers.
func (*SlashCommandsInitializer) updateMarkdownContent(
	fs afero.Fs,
	filePath, frontmatter, body string,
) (string, error) {
	existingContent, err := afero.ReadFile(
		fs,
		filePath,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to read file %s: %w",
			filePath,
			err,
		)
	}

	contentStr := string(existingContent)

	// Find markers
	startIndex := findMarkerIndex(
		contentStr,
		spectrStartMarker,
		0,
	)
	if startIndex == -1 {
		// No start marker, treat as new file but keep existing content
		// Append markers at the end
		newContent := createMarkdownContent(
			frontmatter,
			body,
		)

		return contentStr + newlineDouble + newContent, nil
	}

	searchOffset := startIndex + len(
		spectrStartMarker,
	)
	endIndex := findMarkerIndex(
		contentStr,
		spectrEndMarker,
		searchOffset,
	)
	if endIndex == -1 {
		return "", fmt.Errorf(
			"end marker not found in %s",
			filePath,
		)
	}

	// Replace content between markers
	before := contentStr[:startIndex]
	after := contentStr[endIndex+len(spectrEndMarker):]

	// Check if frontmatter exists
	hasFrontmatter := strings.HasPrefix(
		strings.TrimSpace(before),
		"---",
	)
	if frontmatter != "" && !hasFrontmatter {
		before = strings.TrimSpace(
			frontmatter,
		) + newlineDouble +
			strings.TrimLeft(
				before,
				"\n\r",
			)
	}

	result := before + spectrStartMarker + newline +
		body + newline + spectrEndMarker + after

	return result, nil
}

// getFrontmatter returns the YAML frontmatter for a command.
func getFrontmatter(command string) string {
	switch command {
	case "proposal":
		return `---
description: Scaffold a new Spectr change and validate strictly.
---`
	case "apply":
		return `---
description: Implement an approved Spectr change and keep tasks in sync.
---`
	default:
		return ""
	}
}

// generateTOMLContent generates TOML content for a Gemini command.
func (*SlashCommandsInitializer) generateTOMLContent(
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

// IsSetup returns true if both proposal and apply command files exist.
// Returns false if any command file is missing.
func (s *SlashCommandsInitializer) IsSetup(
	fs afero.Fs,
	_ *providers.Config,
) bool {
	for _, cmd := range slashCommands() {
		filePath := s.commandPath(cmd.name)
		exists, err := afero.Exists(fs, filePath)
		if err != nil || !exists {
			return false
		}
	}

	return true
}

// Key returns a unique key for this initializer based on its configuration.
// Used for deduplication when multiple providers use the same initializer.
// The key format is "slashcmds:<dir>:<ext>:<format>".
// Example: "slashcmds:.claude/commands/spectr:.md:0"
func (s *SlashCommandsInitializer) Key() string {
	return fmt.Sprintf(
		"slashcmds:%s:%s:%d",
		s.Dir,
		s.Extension,
		s.Format,
	)
}
