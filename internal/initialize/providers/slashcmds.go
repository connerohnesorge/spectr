package providers

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

// SlashCommandsInitializer creates slash command files
// (proposal.ext, apply.ext).
// Implements the Initializer interface for slash command creation.
//
// Example usage (Markdown):
//
//	init := NewSlashCommandsInitializer(
//	    ".claude/commands/spectr",
//	    ".md",
//	    FormatMarkdown,
//	)
//
// Example usage (TOML):
//
//	init := NewSlashCommandsInitializer(
//	    ".gemini/commands/spectr",
//	    ".toml",
//	    FormatTOML,
//	)
type SlashCommandsInitializer struct {
	Dir        string
	Ext        string
	Format     CommandFormat
	IsGlobalFs bool
}

// NewSlashCommandsInitializer creates a new
// SlashCommandsInitializer.
//
// Parameters:
//   - dir: Directory where command files will be created
//     (e.g., ".claude/commands/spectr")
//   - ext: File extension (e.g., ".md", ".toml")
//   - format: Command format (FormatMarkdown or FormatTOML)
//
// Returns:
//   - *SlashCommandsInitializer: A new slash commands initializer
func NewSlashCommandsInitializer(
	dir, ext string,
	format CommandFormat,
) *SlashCommandsInitializer {
	return &SlashCommandsInitializer{
		Dir:        dir,
		Ext:        ext,
		Format:     format,
		IsGlobalFs: false,
	}
}

// WithGlobal configures the initializer to use the global
// filesystem.
func (s *SlashCommandsInitializer) WithGlobal(
	global bool,
) *SlashCommandsInitializer {
	s.IsGlobalFs = global

	return s
}

// Init creates or updates slash command files (proposal and apply).
//
// For Markdown format:
//   - Creates files with YAML frontmatter and template content
//   - Frontmatter includes description field
//
// For TOML format:
//   - Creates TOML files with description and prompt fields
//   - Escapes template content for TOML multiline strings
//
// Parameters:
//   - ctx: Context for cancellation
//   - fs: Filesystem abstraction
//   - cfg: Configuration with SpectrDir and derived paths
//   - tm: TemplateManager for rendering command templates
//
// Returns:
//   - InitResult: Contains created or updated command files
//   - error: Non-nil if initialization fails
func (s *SlashCommandsInitializer) Init(
	_ context.Context,
	fs afero.Fs,
	cfg *Config,
	tm TemplateManager,
) (InitResult, error) {
	var result InitResult

	// Create template context from config
	templateCtx := TemplateContext{
		BaseDir:     cfg.SpectrDir,
		SpecsDir:    cfg.SpecsDir(),
		ChangesDir:  cfg.ChangesDir(),
		ProjectFile: cfg.ProjectFile(),
		AgentsFile:  cfg.AgentsFile(),
	}

	// Ensure directory exists
	if err := fs.MkdirAll(s.Dir, dirPerm); err != nil {
		return InitResult{}, fmt.Errorf(
			"failed to create directory %s: %w",
			s.Dir,
			err,
		)
	}

	// Create proposal and apply commands
	commands := []struct {
		name        string
		description string
	}{
		{
			"proposal",
			"Scaffold a new Spectr change and validate strictly.",
		},
		{
			"apply",
			"Implement an approved Spectr change and keep tasks in sync.",
		},
	}

	for _, cmd := range commands {
		filePath := filepath.Join(s.Dir, cmd.name+s.Ext)

		// Check if file already exists
		exists, err := afero.Exists(fs, filePath)
		if err != nil {
			return InitResult{}, fmt.Errorf(
				"failed to check if file exists %s: %w",
				filePath,
				err,
			)
		}

		// Render template content
		content, err := tm.RenderSlashCommand(cmd.name, templateCtx)
		if err != nil {
			return InitResult{}, fmt.Errorf(
				"failed to render slash command %s: %w",
				cmd.name,
				err,
			)
		}

		var fileContent string
		switch s.Format {
		case FormatMarkdown:
			fileContent = s.generateMarkdownContent(cmd.description, content)
		case FormatTOML:
			fileContent = s.generateTOMLContent(cmd.description, content)
		default:
			return InitResult{}, fmt.Errorf(
				"unsupported command format: %d",
				s.Format,
			)
		}

		// Write file
		if err := afero.WriteFile(
			fs,
			filePath,
			[]byte(fileContent),
			filePerm,
		); err != nil {
			return InitResult{}, fmt.Errorf(
				"failed to write command file %s: %w",
				filePath,
				err,
			)
		}

		if exists {
			result.UpdatedFiles = append(result.UpdatedFiles, filePath)
		} else {
			result.CreatedFiles = append(result.CreatedFiles, filePath)
		}
	}

	return result, nil
}

// generateMarkdownContent creates Markdown content with YAML
// frontmatter.
func (*SlashCommandsInitializer) generateMarkdownContent(
	description, content string,
) string {
	frontmatter := fmt.Sprintf("---\ndescription: %s\n---", description)

	return frontmatter + "\n\n" + content + "\n"
}

// generateTOMLContent creates TOML content for a command.
// Escapes the content for TOML multiline strings.
func (*SlashCommandsInitializer) generateTOMLContent(
	description, prompt string,
) string {
	// Escape the prompt for TOML multiline string
	escapedPrompt := strings.ReplaceAll(prompt, `\`, `\\`)
	escapedPrompt = strings.ReplaceAll(escapedPrompt, `"`, `\"`)

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

// IsSetup returns true if both command files exist.
//
// Parameters:
//   - fs: Filesystem abstraction
//   - cfg: Configuration (not used)
//
// Returns:
//   - bool: True if proposal and apply command files exist
func (s *SlashCommandsInitializer) IsSetup(fs afero.Fs, _ *Config) bool {
	proposalPath := filepath.Join(s.Dir, "proposal"+s.Ext)
	applyPath := filepath.Join(s.Dir, "apply"+s.Ext)

	proposalExists, err := afero.Exists(fs, proposalPath)
	if err != nil || !proposalExists {
		return false
	}

	applyExists, err := afero.Exists(fs, applyPath)
	if err != nil || !applyExists {
		return false
	}

	return true
}

// Path returns the directory path for deduplication.
// This ensures multiple providers using the same directory share
// one initializer.
//
// Returns:
//   - string: The command directory path
func (s *SlashCommandsInitializer) Path() string {
	return s.Dir
}

// IsGlobal returns true if this initializer uses the global filesystem.
//
// Returns:
//   - bool: True if using global filesystem, false for project-relative
func (s *SlashCommandsInitializer) IsGlobal() bool {
	return s.IsGlobalFs
}
