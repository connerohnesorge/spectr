// Package initializers provides built-in initializers for the provider system.
//
// This file contains the SlashCommandsInitializer, which creates slash command
// files (proposal and apply) for AI tools in either Markdown or TOML format.
package initializers

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/connerohnesorge/spectr/internal/initialize/providers"
	"github.com/spf13/afero"
)

// CommandFormat specifies the format for slash command files.
type CommandFormat int

const (
	// FormatMarkdown uses markdown files with YAML frontmatter (Claude, Cline, etc.)
	FormatMarkdown CommandFormat = iota
	// FormatTOML uses TOML files (Gemini CLI)
	FormatTOML
)

// Command definitions for slash commands.
const (
	// CommandProposal is the proposal command name.
	CommandProposal = "proposal"
	// CommandApply is the apply command name.
	CommandApply = "apply"
)

// Default TOML descriptions for commands.
const (
	// TomlDescriptionProposal is the TOML description for the proposal command.
	TomlDescriptionProposal = "Scaffold a new Spectr change and validate strictly."
	// TomlDescriptionApply is the TOML description for the apply command.
	TomlDescriptionApply = "Implement an approved Spectr change and keep tasks in sync."
)

// SlashCommandsInitializer creates slash command files from templates.
//
// SlashCommandsInitializer manages the creation of proposal and apply command files
// for AI tools. It supports two formats:
//
//   - FormatMarkdown: Creates markdown files with optional YAML frontmatter and
//     content wrapped in spectr markers (used by Claude, Cline, Cursor, etc.)
//
//   - FormatTOML: Creates TOML files with description and prompt fields
//     (used by Gemini CLI)
//
// # Execution Order
//
// SlashCommandsInitializer has priority 3 in the initializer ordering, meaning it
// runs after DirectoryInitializer (priority 1) and ConfigFileInitializer
// (priority 2). This ensures parent directories exist before files are written.
//
// # Idempotency
//
// For Markdown format: Existing files are updated using marker-based replacement,
// preserving content outside the markers.
//
// For TOML format: Files are overwritten completely since TOML files don't support
// marker-based updates.
//
// # Example Usage
//
//	func (p *ClaudeProvider) Initializers(ctx context.Context) []providers.Initializer {
//	    return []providers.Initializer{
//	        initializers.NewDirectoryInitializer(".claude/commands/spectr"),
//	        initializers.NewSlashCommandsInitializer(
//	            ".claude/commands/spectr",
//	            ".md",
//	            initializers.FormatMarkdown,
//	        ),
//	    }
//	}
//
// # With Frontmatter (for Markdown format)
//
//	frontmatter := map[string]string{
//	    "proposal": "---\ndescription: Scaffold a new Spectr change.\n---",
//	    "apply":    "---\ndescription: Implement an approved Spectr change.\n---",
//	}
//	init := initializers.NewSlashCommandsInitializerWithFrontmatter(
//	    ".claude/commands/spectr",
//	    ".md",
//	    initializers.FormatMarkdown,
//	    frontmatter,
//	)
//
// # TOML Format Example
//
//	init := initializers.NewSlashCommandsInitializer(
//	    ".gemini/commands/spectr",
//	    ".toml",
//	    initializers.FormatTOML,
//	)
//
// # Global Commands
//
// For slash commands stored globally (e.g., in the user's home directory),
// use NewGlobalSlashCommandsInitializer:
//
//	init := initializers.NewGlobalSlashCommandsInitializer(
//	    ".config/gemini/commands/spectr",
//	    ".toml",
//	    initializers.FormatTOML,
//	)
type SlashCommandsInitializer struct {
	// dir is the directory where command files are stored (relative to fs root)
	dir string
	// ext is the file extension (e.g., ".md", ".toml")
	ext string
	// format is the command file format (Markdown or TOML)
	format CommandFormat
	// frontmatter is optional YAML frontmatter for Markdown commands (command name -> frontmatter)
	frontmatter map[string]string
	// global indicates whether to use globalFs (true) or projectFs (false)
	global bool
}

// NewSlashCommandsInitializer creates a new SlashCommandsInitializer for project commands.
//
// Parameters:
//   - dir: Directory path relative to project root (e.g., ".claude/commands/spectr")
//   - ext: File extension including the dot (e.g., ".md", ".toml")
//   - format: Command file format (FormatMarkdown or FormatTOML)
//
// Example:
//
//	init := NewSlashCommandsInitializer(".claude/commands/spectr", ".md", FormatMarkdown)
//	init := NewSlashCommandsInitializer(".gemini/commands/spectr", ".toml", FormatTOML)
func NewSlashCommandsInitializer(dir, ext string, format CommandFormat) *SlashCommandsInitializer {
	return &SlashCommandsInitializer{
		dir:         dir,
		ext:         ext,
		format:      format,
		frontmatter: nil,
		global:      false,
	}
}

// NewSlashCommandsInitializerWithFrontmatter creates a SlashCommandsInitializer with
// custom YAML frontmatter for Markdown format commands.
//
// Parameters:
//   - dir: Directory path relative to project root
//   - ext: File extension including the dot (e.g., ".md")
//   - format: Command file format (should be FormatMarkdown for frontmatter to apply)
//   - frontmatter: Map of command name to YAML frontmatter content
//
// The frontmatter map keys should be "proposal" and/or "apply". The values should
// include the YAML delimiters, for example:
//
//	frontmatter := map[string]string{
//	    "proposal": "---\ndescription: Scaffold a new Spectr change.\n---",
//	    "apply":    "---\ndescription: Implement an approved Spectr change.\n---",
//	}
//
// Example:
//
//	init := NewSlashCommandsInitializerWithFrontmatter(
//	    ".claude/commands/spectr",
//	    ".md",
//	    FormatMarkdown,
//	    map[string]string{
//	        "proposal": "---\ndescription: Scaffold a new Spectr change.\n---",
//	        "apply":    "---\ndescription: Implement an approved Spectr change.\n---",
//	    },
//	)
func NewSlashCommandsInitializerWithFrontmatter(dir, ext string, format CommandFormat, frontmatter map[string]string) *SlashCommandsInitializer {
	return &SlashCommandsInitializer{
		dir:         dir,
		ext:         ext,
		format:      format,
		frontmatter: frontmatter,
		global:      false,
	}
}

// NewGlobalSlashCommandsInitializer creates a SlashCommandsInitializer for global commands.
//
// Global commands are stored in the user's home directory rather than the project
// directory. This is used for tools that read commands from a global location.
//
// Parameters:
//   - dir: Directory path relative to user's home directory (e.g., ".config/gemini/commands/spectr")
//   - ext: File extension including the dot (e.g., ".toml")
//   - format: Command file format (FormatMarkdown or FormatTOML)
//
// Example:
//
//	init := NewGlobalSlashCommandsInitializer(".config/gemini/commands/spectr", ".toml", FormatTOML)
func NewGlobalSlashCommandsInitializer(dir, ext string, format CommandFormat) *SlashCommandsInitializer {
	return &SlashCommandsInitializer{
		dir:         dir,
		ext:         ext,
		format:      format,
		frontmatter: nil,
		global:      true,
	}
}

// NewGlobalSlashCommandsInitializerWithFrontmatter creates a global SlashCommandsInitializer
// with custom YAML frontmatter for Markdown format commands.
//
// Parameters:
//   - dir: Directory path relative to user's home directory
//   - ext: File extension including the dot (e.g., ".md")
//   - format: Command file format (should be FormatMarkdown for frontmatter to apply)
//   - frontmatter: Map of command name to YAML frontmatter content
//
// Example:
//
//	init := NewGlobalSlashCommandsInitializerWithFrontmatter(
//	    ".config/tool/commands/spectr",
//	    ".md",
//	    FormatMarkdown,
//	    map[string]string{
//	        "proposal": "---\ndescription: Scaffold a new Spectr change.\n---",
//	    },
//	)
func NewGlobalSlashCommandsInitializerWithFrontmatter(dir, ext string, format CommandFormat, frontmatter map[string]string) *SlashCommandsInitializer {
	return &SlashCommandsInitializer{
		dir:         dir,
		ext:         ext,
		format:      format,
		frontmatter: frontmatter,
		global:      true,
	}
}

// Init creates or updates the slash command files.
//
// Init creates both the proposal and apply command files in the specified directory.
// The file format depends on the CommandFormat specified at construction time.
//
// For Markdown format:
//   - New files are created with optional frontmatter followed by content in markers
//   - Existing files with markers have their content between markers replaced
//   - Existing files without markers have markers and content appended
//
// For TOML format:
//   - Files are always overwritten with new TOML content
//   - TOML format includes description and prompt fields
//
// The content is rendered using the TemplateManager's RenderSlashCommand method.
//
// Parameters:
//   - ctx: Context for cancellation (passed to template rendering)
//   - fs: Filesystem to operate on (project or global, based on IsGlobal())
//   - cfg: Configuration containing spectr directory paths (used for template context)
//   - tm: Template manager for rendering the slash command templates
//
// Returns an error if file operations fail or template rendering fails.
func (s *SlashCommandsInitializer) Init(ctx context.Context, fs afero.Fs, cfg *providers.Config, tm providers.TemplateManager) error {
	// Ensure directory exists
	if err := fs.MkdirAll(s.dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", s.dir, err)
	}

	// Create both proposal and apply commands
	commands := []string{CommandProposal, CommandApply}
	for _, cmd := range commands {
		if err := s.initCommand(ctx, fs, cfg, tm, cmd); err != nil {
			return err
		}
	}

	return nil
}

// initCommand creates or updates a single command file.
func (s *SlashCommandsInitializer) initCommand(ctx context.Context, fs afero.Fs, cfg *providers.Config, tm providers.TemplateManager, cmd string) error {
	// Render the command template
	templateCtx := providers.NewTemplateContext(cfg)
	body, err := tm.RenderSlashCommand(cmd, templateCtx)
	if err != nil {
		return fmt.Errorf("failed to render slash command %s: %w", cmd, err)
	}

	// Build file path
	filePath := filepath.Join(s.dir, cmd+s.ext)

	// Create based on format
	switch s.format {
	case FormatMarkdown:
		return s.initMarkdownCommand(fs, filePath, cmd, body)
	case FormatTOML:
		return s.initTOMLCommand(fs, filePath, cmd, body)
	default:
		return fmt.Errorf("unsupported command format: %d", s.format)
	}
}

// initMarkdownCommand creates or updates a Markdown command file.
func (s *SlashCommandsInitializer) initMarkdownCommand(fs afero.Fs, filePath, cmd, body string) error {
	// Check if file exists
	exists, err := afero.Exists(fs, filePath)
	if err != nil {
		return fmt.Errorf("failed to check if file exists: %w", err)
	}

	if exists {
		return s.updateExistingMarkdownCommand(fs, filePath, cmd, body)
	}

	return s.createNewMarkdownCommand(fs, filePath, cmd, body)
}

// createNewMarkdownCommand creates a new Markdown command file.
func (s *SlashCommandsInitializer) createNewMarkdownCommand(fs afero.Fs, filePath, cmd, body string) error {
	var sections []string

	// Add frontmatter if provided
	if frontmatter, ok := s.frontmatter[cmd]; ok && frontmatter != "" {
		sections = append(sections, strings.TrimSpace(frontmatter))
	}

	// Add body wrapped in markers
	markerContent := spectrStartMarker + newlineDouble + body + newlineDouble + spectrEndMarker
	sections = append(sections, markerContent)

	content := strings.Join(sections, newlineDouble) + newlineDouble

	if err := afero.WriteFile(fs, filePath, []byte(content), configFilePerm); err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}

	return nil
}

// updateExistingMarkdownCommand updates an existing Markdown command file.
func (s *SlashCommandsInitializer) updateExistingMarkdownCommand(fs afero.Fs, filePath, cmd, body string) error {
	// Read existing file
	existingContent, err := afero.ReadFile(fs, filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	existingStr := string(existingContent)

	// Find markers
	startIndex := findMarkerIndex(existingStr, spectrStartMarker, 0)
	endIndex := -1
	if startIndex != -1 {
		searchOffset := startIndex + len(spectrStartMarker)
		endIndex = findMarkerIndex(existingStr, spectrEndMarker, searchOffset)
	}

	var finalContent string
	if startIndex == -1 || endIndex == -1 {
		// No markers found, append to end
		markerContent := spectrStartMarker + newline + body + newline + spectrEndMarker + newline
		finalContent = existingStr + newlineDouble + markerContent
	} else {
		// Replace content between markers
		before := existingStr[:startIndex]
		after := existingStr[endIndex+len(spectrEndMarker):]
		finalContent = before + spectrStartMarker + newlineDouble + body + newlineDouble + spectrEndMarker + after
	}

	if err := afero.WriteFile(fs, filePath, []byte(finalContent), configFilePerm); err != nil {
		return fmt.Errorf("failed to update file %s: %w", filePath, err)
	}

	return nil
}

// initTOMLCommand creates or updates a TOML command file.
func (s *SlashCommandsInitializer) initTOMLCommand(fs afero.Fs, filePath, cmd, body string) error {
	// Get description for this command
	description := s.getTOMLDescription(cmd)

	// Generate TOML content
	content := generateTOMLContent(description, body)

	if err := afero.WriteFile(fs, filePath, []byte(content), configFilePerm); err != nil {
		return fmt.Errorf("failed to write TOML command file %s: %w", filePath, err)
	}

	return nil
}

// getTOMLDescription returns the description for a TOML command.
func (s *SlashCommandsInitializer) getTOMLDescription(cmd string) string {
	switch cmd {
	case CommandProposal:
		return TomlDescriptionProposal
	case CommandApply:
		return TomlDescriptionApply
	default:
		return fmt.Sprintf("Spectr %s command.", cmd)
	}
}

// generateTOMLContent creates TOML content for a command.
func generateTOMLContent(description, prompt string) string {
	// Escape the prompt for TOML multiline string
	escapedPrompt := strings.ReplaceAll(prompt, `\`, `\\`)
	escapedPrompt = strings.ReplaceAll(escapedPrompt, `"`, `\"`)

	return fmt.Sprintf(`# Spectr command for Gemini CLI
description = "%s"
prompt = """
%s
"""
`, description, escapedPrompt)
}

// IsSetup returns true if both proposal and apply command files exist.
//
// IsSetup checks for the presence of both command files. It does not verify
// the content of the files, only their existence.
//
// Parameters:
//   - fs: Filesystem to check (project or global, based on IsGlobal())
//   - cfg: Configuration (currently unused by SlashCommandsInitializer)
func (s *SlashCommandsInitializer) IsSetup(fs afero.Fs, cfg *providers.Config) bool {
	proposalPath := filepath.Join(s.dir, CommandProposal+s.ext)
	applyPath := filepath.Join(s.dir, CommandApply+s.ext)

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

// Path returns the directory path this initializer manages.
//
// Path returns the slash commands directory, which is used for deduplication.
// When multiple providers return SlashCommandsInitializers with the same Path(),
// only one will be executed.
//
// Example return values:
//   - ".claude/commands/spectr"
//   - ".gemini/commands/spectr"
//   - ".config/tool/commands/spectr"
func (s *SlashCommandsInitializer) Path() string {
	return s.dir
}

// IsGlobal returns true if this initializer uses the global filesystem.
//
// Project initializers (IsGlobal() == false) operate on files within
// the project directory (e.g., .claude/commands/spectr/).
//
// Global initializers (IsGlobal() == true) operate on files in the
// user's home directory (e.g., ~/.config/gemini/commands/spectr/).
func (s *SlashCommandsInitializer) IsGlobal() bool {
	return s.global
}

// Dir returns the directory path for the slash commands.
//
// This is useful for debugging or logging the specific directory
// where command files will be created.
func (s *SlashCommandsInitializer) Dir() string {
	return s.dir
}

// Ext returns the file extension used for command files.
//
// This includes the dot, e.g., ".md" or ".toml".
func (s *SlashCommandsInitializer) Ext() string {
	return s.ext
}

// Format returns the command file format (Markdown or TOML).
func (s *SlashCommandsInitializer) Format() CommandFormat {
	return s.format
}
