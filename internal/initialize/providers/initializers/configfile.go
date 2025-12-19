// Package initializers provides built-in Initializer implementations for the
// provider architecture.
//
// This file implements ConfigFileInitializer, which creates or updates
// instruction files (e.g., CLAUDE.md) with marker-based content sections.
//
//nolint:revive // line-length-limit, unused-parameter - interface compliance
package initializers

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/connerohnesorge/spectr/internal/initialize/providers"
	"github.com/spf13/afero"
)

// File and directory permission constants.
const (
	dirPerm  = 0o755
	filePerm = 0o644

	// Marker constants for managing config file updates.
	spectrStartMarker = "<!-- spectr:START -->"
	spectrEndMarker   = "<!-- spectr:END -->"

	// Common strings.
	newline       = "\n"
	newlineDouble = "\n\n"
)

// ConfigFileInitializer creates or updates instruction files with marker-based
// content sections (e.g., CLAUDE.md with <!-- spectr:START --> markers).
//
// It implements the providers.Initializer interface and is designed to:
//   - Create new files with marker sections if they don't exist
//   - Update existing files by replacing content between markers
//   - Be idempotent (safe to run multiple times)
//   - Work with afero.Fs for testability
type ConfigFileInitializer struct {
	// path is the file path to create/update (e.g., "CLAUDE.md")
	path string

	// templateName identifies which template to render (e.g., "instruction-pointer")
	templateName string

	// isGlobal indicates whether this initializer operates on global paths
	// (relative to home directory) instead of project-relative paths.
	isGlobal bool
}

// NewConfigFileInitializer creates a new ConfigFileInitializer that will
// create or update the specified config file with rendered template content.
//
// Parameters:
//   - path: the file path to create/update (e.g., "CLAUDE.md")
//   - templateName: identifies which template to render (e.g., "instruction-pointer")
//   - isGlobal: if true, path is relative to home directory; otherwise project-relative
//
// Returns nil if path is empty.
func NewConfigFileInitializer(
	path string,
	templateName string,
	isGlobal bool,
) *ConfigFileInitializer {
	if path == "" {
		return nil
	}

	return &ConfigFileInitializer{
		path:         path,
		templateName: templateName,
		isGlobal:     isGlobal,
	}
}

// Init creates or updates the config file with marker-based content.
//
// If the file doesn't exist, it creates a new file with markers containing
// the rendered template content. If the file exists, it replaces the content
// between the markers or appends the marker section if markers don't exist.
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
func (c *ConfigFileInitializer) Init(
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

	// Render the template content
	content, err := tm.RenderInstructionPointer(templateCtx)
	if err != nil {
		return fmt.Errorf("failed to render template %q: %w", c.templateName, err)
	}

	// Update or create the file with markers
	return c.updateFileWithMarkers(fs, c.path, content)
}

// updateFileWithMarkers updates content between markers in a file,
// or creates the file with markers if it doesn't exist.
func (c *ConfigFileInitializer) updateFileWithMarkers(fs afero.Fs, filePath, content string) error {
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := fs.MkdirAll(dir, dirPerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Check if file exists
	exists, err := afero.Exists(fs, filePath)
	if err != nil {
		return fmt.Errorf("failed to check file existence: %w", err)
	}

	if !exists {
		// Create new file with markers
		newContent := spectrStartMarker + newline + content + newline + spectrEndMarker + newline

		return afero.WriteFile(fs, filePath, []byte(newContent), filePerm)
	}

	// Read existing file
	existingContent, err := afero.ReadFile(fs, filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	contentStr := string(existingContent)

	// Find markers
	startIndex := c.findMarkerIndex(contentStr, spectrStartMarker, 0)
	endIndex := -1
	if startIndex != -1 {
		searchOffset := startIndex + len(spectrStartMarker)
		endIndex = c.findMarkerIndex(contentStr, spectrEndMarker, searchOffset)
	}

	var newContent string
	if startIndex == -1 || endIndex == -1 {
		// No markers found, append to end
		newContent = contentStr + newlineDouble +
			spectrStartMarker + newline + content + newline + spectrEndMarker + newline
	} else {
		// Replace content between markers
		before := contentStr[:startIndex]
		after := contentStr[endIndex+len(spectrEndMarker):]
		newContent = before + spectrStartMarker + newline +
			content + newline + spectrEndMarker + after
	}

	return afero.WriteFile(fs, filePath, []byte(newContent), filePerm)
}

// findMarkerIndex finds the index of a marker in content, starting from offset.
func (c *ConfigFileInitializer) findMarkerIndex(content, marker string, offset int) int {
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
//
// Optionally, this could also check if markers are present, but for simplicity
// we just check file existence. The Init method is idempotent anyway.
//
// Parameters:
//   - fs: filesystem abstraction to check
//   - cfg: configuration (not currently used but part of interface)
//
// Returns true if the file exists, false otherwise.
func (c *ConfigFileInitializer) IsSetup(fs afero.Fs, cfg *providers.Config) bool {
	exists, err := afero.Exists(fs, c.path)
	if err != nil {
		return false
	}

	return exists
}

// Path returns the file path this initializer manages.
//
// This is used for deduplication: when multiple providers return
// ConfigFileInitializers with the same Path(), only the first one is executed.
func (c *ConfigFileInitializer) Path() string {
	return c.path
}

// IsGlobal returns true if this initializer uses globalFs instead of projectFs.
//
// Global initializers operate on paths relative to the user's home directory
// (e.g., ~/.config/tool/config). Project initializers operate on paths
// relative to the project root (e.g., CLAUDE.md).
func (c *ConfigFileInitializer) IsGlobal() bool {
	return c.isGlobal
}

// TemplateName returns the template name this initializer uses.
// This is useful for testing and inspection.
func (c *ConfigFileInitializer) TemplateName() string {
	return c.templateName
}

// Ensure ConfigFileInitializer implements the Initializer interface.
var _ providers.Initializer = (*ConfigFileInitializer)(nil)
