// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file contains built-in initializer implementations that can be used
// by providers within this package. These implementations avoid import cycles
// by being defined in the same package as the Initializer interface.
//
// Note: There is also an initializers sub-package that contains the same
// implementations but with its own package. The sub-package is used for
// external code that needs to create initializers directly. This file exists
// to allow providers in this package to create initializers without cycles.
//
//nolint:revive // file-length-limit - logically cohesive initializer implementations
package providers

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

// File and directory permission constants.
const (
	dirPermBuiltin  = 0o755
	filePermBuiltin = 0o644

	// Marker constants for managing config file updates.
	spectrStartMarkerBuiltin = "<!-- spectr:START -->"
	spectrEndMarkerBuiltin   = "<!-- spectr:END -->"

	// Common strings.
	newlineBuiltin       = "\n"
	newlineDoubleBuiltin = "\n\n"
)

// -----------------------------------------------------------------------------
// DirectoryInitializer
// -----------------------------------------------------------------------------

// DirectoryInitializerBuiltin creates directories needed by providers.
// This is the in-package version to avoid import cycles.
type DirectoryInitializerBuiltin struct {
	paths    []string
	isGlobal bool
}

// NewDirectoryInitializer creates a new DirectoryInitializer that will create
// the specified directories.
//
// Parameters:
//   - isGlobal: if true, paths are relative to home directory; otherwise project-relative
//   - paths: one or more directory paths to create
//
// Returns nil if no paths are provided.
func NewDirectoryInitializer(isGlobal bool, paths ...string) Initializer {
	if len(paths) == 0 {
		return nil
	}

	return &DirectoryInitializerBuiltin{
		paths:    paths,
		isGlobal: isGlobal,
	}
}

// Init creates all directories specified in the initializer.
func (d *DirectoryInitializerBuiltin) Init(
	_ context.Context,
	fs afero.Fs,
	_ *Config,
	_ TemplateRenderer,
) error {
	for _, path := range d.paths {
		if err := fs.MkdirAll(path, dirPermBuiltin); err != nil {
			return err
		}
	}

	return nil
}

// IsSetup returns true if all directories exist.
func (d *DirectoryInitializerBuiltin) IsSetup(fs afero.Fs, _ *Config) bool {
	for _, path := range d.paths {
		info, err := fs.Stat(path)
		if err != nil || !info.IsDir() {
			return false
		}
	}

	return true
}

// Path returns the primary directory path for deduplication.
func (d *DirectoryInitializerBuiltin) Path() string {
	if len(d.paths) == 0 {
		return ""
	}

	return d.paths[0]
}

// IsGlobal returns true if this initializer uses globalFs.
func (d *DirectoryInitializerBuiltin) IsGlobal() bool {
	return d.isGlobal
}

// Ensure DirectoryInitializerBuiltin implements the Initializer interface.
var _ Initializer = (*DirectoryInitializerBuiltin)(nil)

// -----------------------------------------------------------------------------
// ConfigFileInitializer
// -----------------------------------------------------------------------------

// ConfigFileInitializerBuiltin creates or updates instruction files with
// marker-based content sections.
type ConfigFileInitializerBuiltin struct {
	path         string
	templateName string
	isGlobal     bool
}

// NewConfigFileInitializer creates a new ConfigFileInitializer.
//
// Parameters:
//   - path: the file path to create/update (e.g., "CLAUDE.md")
//   - templateName: identifies which template to render (e.g., "instruction-pointer")
//   - isGlobal: if true, path is relative to home directory; otherwise project-relative
//
// Returns nil if path is empty.
func NewConfigFileInitializer(path, templateName string, isGlobal bool) Initializer {
	if path == "" {
		return nil
	}

	return &ConfigFileInitializerBuiltin{
		path:         path,
		templateName: templateName,
		isGlobal:     isGlobal,
	}
}

// Init creates or updates the config file with marker-based content.
func (c *ConfigFileInitializerBuiltin) Init(
	_ context.Context,
	fs afero.Fs,
	cfg *Config,
	tm TemplateRenderer,
) error {
	// Build template context from config
	templateCtx := TemplateContext{
		BaseDir:     cfg.SpectrDir,
		SpecsDir:    cfg.SpecsDir(),
		ChangesDir:  cfg.ChangesDir(),
		ProjectFile: cfg.ProjectFile(),
		AgentsFile:  cfg.AgentsFile(),
	}

	// Render the template content
	content, err := tm.RenderInstructionPointer(templateCtx)
	if err != nil {
		return fmt.Errorf("failed to render template %q: %w", c.templateName, err)
	}

	// Update or create the file with markers
	return c.updateFileWithMarkers(fs, c.path, content)
}

// updateFileWithMarkers updates content between markers in a file.
func (c *ConfigFileInitializerBuiltin) updateFileWithMarkers(
	fs afero.Fs,
	filePath, content string,
) error {
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := fs.MkdirAll(dir, dirPermBuiltin); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Check if file exists
	exists, err := afero.Exists(fs, filePath)
	if err != nil {
		return fmt.Errorf("failed to check file existence: %w", err)
	}

	if !exists {
		// Create new file with markers
		newContent := spectrStartMarkerBuiltin + newlineBuiltin + content + newlineBuiltin + spectrEndMarkerBuiltin + newlineBuiltin

		return afero.WriteFile(fs, filePath, []byte(newContent), filePermBuiltin)
	}

	// Read existing file
	existingContent, err := afero.ReadFile(fs, filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	contentStr := string(existingContent)

	// Find markers
	startIndex := c.findMarkerIndex(contentStr, spectrStartMarkerBuiltin, 0)
	endIndex := -1
	if startIndex != -1 {
		searchOffset := startIndex + len(spectrStartMarkerBuiltin)
		endIndex = c.findMarkerIndex(contentStr, spectrEndMarkerBuiltin, searchOffset)
	}

	var newContent string
	if startIndex == -1 || endIndex == -1 {
		// No markers found, append to end
		newContent = contentStr + newlineDoubleBuiltin +
			spectrStartMarkerBuiltin + newlineBuiltin + content + newlineBuiltin + spectrEndMarkerBuiltin + newlineBuiltin
	} else {
		// Replace content between markers
		before := contentStr[:startIndex]
		after := contentStr[endIndex+len(spectrEndMarkerBuiltin):]
		newContent = before + spectrStartMarkerBuiltin + newlineBuiltin +
			content + newlineBuiltin + spectrEndMarkerBuiltin + after
	}

	return afero.WriteFile(fs, filePath, []byte(newContent), filePermBuiltin)
}

// findMarkerIndex finds the index of a marker in content.
func (c *ConfigFileInitializerBuiltin) findMarkerIndex(content, marker string, offset int) int {
	if offset >= len(content) {
		return -1
	}
	idx := strings.Index(content[offset:], marker)
	if idx == -1 {
		return -1
	}

	return offset + idx
}

// IsSetup returns true if the config file exists.
func (c *ConfigFileInitializerBuiltin) IsSetup(fs afero.Fs, _ *Config) bool {
	exists, err := afero.Exists(fs, c.path)

	return err == nil && exists
}

// Path returns the file path for deduplication.
func (c *ConfigFileInitializerBuiltin) Path() string {
	return c.path
}

// IsGlobal returns true if this initializer uses globalFs.
func (c *ConfigFileInitializerBuiltin) IsGlobal() bool {
	return c.isGlobal
}

// Ensure ConfigFileInitializerBuiltin implements the Initializer interface.
var _ Initializer = (*ConfigFileInitializerBuiltin)(nil)

// -----------------------------------------------------------------------------
// SlashCommandsInitializer
// -----------------------------------------------------------------------------

// SlashCommandsInitializerBuiltin creates or updates slash command files.
type SlashCommandsInitializerBuiltin struct {
	dir         string
	ext         string
	format      CommandFormat
	frontmatter map[string]string
	isGlobal    bool
}

// NewSlashCommandsInitializer creates a new SlashCommandsInitializer.
//
// Parameters:
//   - dir: the directory where command files are created (e.g., ".claude/commands/spectr")
//   - ext: the file extension for command files (e.g., ".md", ".toml")
//   - format: Markdown or TOML format for command files
//   - frontmatter: maps command names to their frontmatter content
//   - isGlobal: if true, dir is relative to home directory; otherwise project-relative
//
// Returns nil if dir is empty.
//
//nolint:revive // argument-limit - all params are necessary for initialization
func NewSlashCommandsInitializer(
	dir, ext string,
	format CommandFormat,
	frontmatter map[string]string,
	isGlobal bool,
) Initializer {
	if dir == "" {
		return nil
	}

	return &SlashCommandsInitializerBuiltin{
		dir:         dir,
		ext:         ext,
		format:      format,
		frontmatter: frontmatter,
		isGlobal:    isGlobal,
	}
}

// Init creates or updates the slash command files.
func (s *SlashCommandsInitializerBuiltin) Init(
	_ context.Context,
	fs afero.Fs,
	cfg *Config,
	tm TemplateRenderer,
) error {
	// Build template context from config
	templateCtx := TemplateContext{
		BaseDir:     cfg.SpectrDir,
		SpecsDir:    cfg.SpecsDir(),
		ChangesDir:  cfg.ChangesDir(),
		ProjectFile: cfg.ProjectFile(),
		AgentsFile:  cfg.AgentsFile(),
	}

	// Ensure directory exists
	if err := fs.MkdirAll(s.dir, dirPermBuiltin); err != nil {
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
//
//nolint:revive // argument-limit - all params are contextually related
func (s *SlashCommandsInitializerBuiltin) configureCommand(
	fs afero.Fs,
	filePath, cmd string,
	templateCtx TemplateContext,
	tm TemplateRenderer,
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
func (s *SlashCommandsInitializerBuiltin) updateExistingCommand(
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
	startIndex := s.findMarkerIndex(contentStr, spectrStartMarkerBuiltin, 0)
	if startIndex == -1 {
		return fmt.Errorf("start marker not found in %s", filePath)
	}

	searchOffset := startIndex + len(spectrStartMarkerBuiltin)
	endIndex := s.findMarkerIndex(contentStr, spectrEndMarkerBuiltin, searchOffset)
	if endIndex == -1 {
		return fmt.Errorf("end marker not found in %s", filePath)
	}

	if endIndex < startIndex {
		return fmt.Errorf("end marker appears before start marker in %s", filePath)
	}

	before := contentStr[:startIndex]
	after := contentStr[endIndex+len(spectrEndMarkerBuiltin):]

	// Check if frontmatter needs to be added
	frontmatter := s.frontmatter[cmd]
	hasFrontmatter := strings.HasPrefix(strings.TrimSpace(before), "---")
	if frontmatter != "" && !hasFrontmatter {
		before = strings.TrimSpace(frontmatter) + newlineDoubleBuiltin +
			strings.TrimLeft(before, "\n\r")
	}

	newContent := before + spectrStartMarkerBuiltin + newlineBuiltin +
		body + newlineBuiltin + spectrEndMarkerBuiltin + after

	if err := afero.WriteFile(fs, filePath, []byte(newContent), filePermBuiltin); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filePath, err)
	}

	return nil
}

// createNewCommand creates a new slash command file.
func (s *SlashCommandsInitializerBuiltin) createNewCommand(
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
	sections = append(
		sections,
		spectrStartMarkerBuiltin+newlineDoubleBuiltin+body+newlineDoubleBuiltin+spectrEndMarkerBuiltin,
	)

	content := strings.Join(sections, newlineDoubleBuiltin) + newlineDoubleBuiltin

	if err := afero.WriteFile(fs, filePath, []byte(content), filePermBuiltin); err != nil {
		return fmt.Errorf("failed to write slash command file %s: %w", filePath, err)
	}

	return nil
}

// findMarkerIndex finds the index of a marker in content.
func (s *SlashCommandsInitializerBuiltin) findMarkerIndex(content, marker string, offset int) int {
	if offset >= len(content) {
		return -1
	}
	idx := strings.Index(content[offset:], marker)
	if idx == -1 {
		return -1
	}

	return offset + idx
}

// IsSetup returns true if all command files exist.
func (s *SlashCommandsInitializerBuiltin) IsSetup(fs afero.Fs, _ *Config) bool {
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

// Path returns the directory path for deduplication.
func (s *SlashCommandsInitializerBuiltin) Path() string {
	return s.dir
}

// IsGlobal returns true if this initializer uses globalFs.
func (s *SlashCommandsInitializerBuiltin) IsGlobal() bool {
	return s.isGlobal
}

// Ensure SlashCommandsInitializerBuiltin implements the Initializer interface.
var _ Initializer = (*SlashCommandsInitializerBuiltin)(nil)
