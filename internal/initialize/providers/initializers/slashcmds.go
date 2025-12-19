// Package initializers provides built-in Initializer implementations for the
// provider architecture.
//
// This file implements SlashCommandsInitializer, which creates or updates
// slash command files (e.g., .claude/commands/spectr/proposal.md).
//
//nolint:revive // line-length-limit, argument-limit - interface compliance
package initializers

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/connerohnesorge/spectr/internal/initialize/providers"
	"github.com/spf13/afero"
)

// SlashCommandsInitializer creates or updates slash command files (proposal and apply).
//
// It implements the providers.Initializer interface and is designed to:
//   - Create proposal and apply command files in the specified directory
//   - Support both Markdown (with YAML frontmatter) and TOML formats
//   - Use markers (spectr:START, spectr:END) for the command body
//   - Be idempotent (safe to run multiple times)
//   - Work with afero.Fs for testability
type SlashCommandsInitializer struct {
	// dir is the directory where command files are created
	// (e.g., ".claude/commands/spectr")
	dir string

	// ext is the file extension for command files (e.g., ".md", ".toml")
	ext string

	// format specifies Markdown or TOML format for command files
	format providers.CommandFormat

	// frontmatter maps command names to their frontmatter content
	// (e.g., "proposal" -> "---\ndescription: ...\n---")
	frontmatter map[string]string

	// isGlobal indicates whether this initializer operates on global paths
	// (relative to home directory) instead of project-relative paths.
	isGlobal bool
}

// NewSlashCommandsInitializer creates a new SlashCommandsInitializer that will
// create or update slash command files in the specified directory.
//
// Parameters:
//   - dir: the directory where command files are created (e.g., ".claude/commands/spectr")
//   - ext: the file extension for command files (e.g., ".md", ".toml")
//   - format: Markdown or TOML format for command files
//   - frontmatter: maps command names to their frontmatter content
//   - isGlobal: if true, dir is relative to home directory; otherwise project-relative
//
// Returns nil if dir is empty.
func NewSlashCommandsInitializer(
	dir string,
	ext string,
	format providers.CommandFormat,
	frontmatter map[string]string,
	isGlobal bool,
) *SlashCommandsInitializer {
	if dir == "" {
		return nil
	}

	return &SlashCommandsInitializer{
		dir:         dir,
		ext:         ext,
		format:      format,
		frontmatter: frontmatter,
		isGlobal:    isGlobal,
	}
}

// Init creates or updates the slash command files (proposal and apply).
//
// For each command (proposal, apply), it:
//   - Renders the command body using the template manager
//   - If the file doesn't exist, creates it with frontmatter and markers
//   - If the file exists, updates the content between markers
//
// This operation is idempotent - running it multiple times has the same effect
// as running it once.
//
// Parameters:
//   - ctx: context for cancellation (not currently used but part of interface)
//   - fs: filesystem abstraction to create/update files on
//   - cfg: configuration containing spectr directory paths
//   - tm: template manager for rendering template content
//
// Returns an error if file creation/update or template rendering fails.
func (s *SlashCommandsInitializer) Init(
	ctx context.Context,
	fs afero.Fs,
	cfg *providers.Config,
	tm providers.TemplateRenderer,
) error {
	// Build template context from config
	templateCtx := providers.TemplateContext{
		BaseDir:     cfg.SpectrDir,
		SpecsDir:    cfg.SpecsDir(),
		ChangesDir:  cfg.ChangesDir(),
		ProjectFile: cfg.ProjectFile(),
		AgentsFile:  cfg.AgentsFile(),
	}

	// Ensure directory exists
	if err := fs.MkdirAll(s.dir, dirPerm); err != nil {
		return fmt.Errorf("failed to create command directory %s: %w", s.dir, err)
	}

	// Commands to create
	commands := []string{"proposal", "apply"}

	for _, cmd := range commands {
		filePath := filepath.Join(s.dir, cmd+s.ext)
		if err := s.configureCommand(fs, filePath, cmd, templateCtx, tm); err != nil {
			return err
		}
	}

	return nil
}

// configureCommand creates or updates a single slash command file.
func (s *SlashCommandsInitializer) configureCommand(
	fs afero.Fs,
	filePath string,
	cmd string,
	templateCtx providers.TemplateContext,
	tm providers.TemplateRenderer,
) error {
	// Render the command body using the template manager
	body, err := tm.RenderSlashCommand(cmd, templateCtx)
	if err != nil {
		return fmt.Errorf("failed to render slash command %s: %w", cmd, err)
	}

	// Check if file exists
	exists, err := afero.Exists(fs, filePath)
	if err != nil {
		return fmt.Errorf("failed to check file existence for %s: %w", filePath, err)
	}

	if exists {
		return s.updateExistingCommand(fs, filePath, body, cmd)
	}

	return s.createNewCommand(fs, filePath, cmd, body)
}

// updateExistingCommand updates an existing slash command file.
func (s *SlashCommandsInitializer) updateExistingCommand(
	fs afero.Fs,
	filePath string,
	body string,
	cmd string,
) error {
	content, err := afero.ReadFile(fs, filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	contentStr := string(content)

	// Find markers
	startIndex := s.findMarkerIndex(contentStr, spectrStartMarker, 0)
	if startIndex == -1 {
		return fmt.Errorf("start marker not found in %s", filePath)
	}

	searchOffset := startIndex + len(spectrStartMarker)
	endIndex := s.findMarkerIndex(contentStr, spectrEndMarker, searchOffset)
	if endIndex == -1 {
		return fmt.Errorf("end marker not found in %s", filePath)
	}

	if endIndex < startIndex {
		return fmt.Errorf("end marker appears before start marker in %s", filePath)
	}

	before := contentStr[:startIndex]
	after := contentStr[endIndex+len(spectrEndMarker):]

	// Check if frontmatter needs to be added
	frontmatter := s.frontmatter[cmd]
	hasFrontmatter := strings.HasPrefix(strings.TrimSpace(before), "---")
	if frontmatter != "" && !hasFrontmatter {
		before = strings.TrimSpace(frontmatter) + newlineDouble +
			strings.TrimLeft(before, "\n\r")
	}

	newContent := before + spectrStartMarker + newline +
		body + newline + spectrEndMarker + after

	if err := afero.WriteFile(fs, filePath, []byte(newContent), filePerm); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filePath, err)
	}

	return nil
}

// createNewCommand creates a new slash command file.
func (s *SlashCommandsInitializer) createNewCommand(
	fs afero.Fs,
	filePath string,
	cmd string,
	body string,
) error {
	var sections []string

	// Add frontmatter if provided
	if frontmatter, ok := s.frontmatter[cmd]; ok && frontmatter != "" {
		sections = append(sections, strings.TrimSpace(frontmatter))
	}

	// Add body with markers
	sections = append(sections,
		spectrStartMarker+newlineDouble+body+newlineDouble+spectrEndMarker)

	content := strings.Join(sections, newlineDouble) + newlineDouble

	if err := afero.WriteFile(fs, filePath, []byte(content), filePerm); err != nil {
		return fmt.Errorf("failed to write slash command file %s: %w", filePath, err)
	}

	return nil
}

// findMarkerIndex finds the index of a marker in content, starting from offset.
func (s *SlashCommandsInitializer) findMarkerIndex(content, marker string, offset int) int {
	if offset >= len(content) {
		return -1
	}
	idx := strings.Index(content[offset:], marker)
	if idx == -1 {
		return -1
	}

	return offset + idx
}

// IsSetup returns true if all command files managed by this initializer exist.
//
// Parameters:
//   - fs: filesystem abstraction to check
//   - cfg: configuration (not currently used but part of interface)
//
// Returns true if both proposal and apply command files exist, false otherwise.
func (s *SlashCommandsInitializer) IsSetup(fs afero.Fs, cfg *providers.Config) bool {
	commands := []string{"proposal", "apply"}
	for _, cmd := range commands {
		filePath := filepath.Join(s.dir, cmd+s.ext)
		exists, err := afero.Exists(fs, filePath)
		if err != nil || !exists {
			return false
		}
	}

	return true
}

// Path returns the directory path this initializer manages.
//
// This is used for deduplication: when multiple providers return
// SlashCommandsInitializers with the same Path(), only the first one is executed.
func (s *SlashCommandsInitializer) Path() string {
	return s.dir
}

// IsGlobal returns true if this initializer uses globalFs instead of projectFs.
//
// Global initializers operate on paths relative to the user's home directory
// (e.g., ~/.config/tool/commands/). Project initializers operate on paths
// relative to the project root (e.g., .claude/commands/).
func (s *SlashCommandsInitializer) IsGlobal() bool {
	return s.isGlobal
}

// Dir returns the directory path for command files.
// This is useful for testing and inspection.
func (s *SlashCommandsInitializer) Dir() string {
	return s.dir
}

// Ext returns the file extension for command files.
// This is useful for testing and inspection.
func (s *SlashCommandsInitializer) Ext() string {
	return s.ext
}

// Format returns the command format (Markdown or TOML).
// This is useful for testing and inspection.
func (s *SlashCommandsInitializer) Format() providers.CommandFormat {
	return s.format
}

// Frontmatter returns the frontmatter map for command files.
// This is useful for testing and inspection.
func (s *SlashCommandsInitializer) Frontmatter() map[string]string {
	return s.frontmatter
}

// Ensure SlashCommandsInitializer implements the Initializer interface.
var _ providers.Initializer = (*SlashCommandsInitializer)(nil)
