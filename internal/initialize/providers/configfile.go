package providers

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

// ConfigFileInitializer creates or updates instruction files with markers.
// Implements the Initializer interface for config file management.
//
// Example usage:
//
//	init := NewConfigFileInitializer("CLAUDE.md", RenderInstructionPointer)
type ConfigFileInitializer struct {
	FilePath   string
	Renderer   TemplateRenderer
	IsGlobalFs bool
}

// NewConfigFileInitializer creates a new ConfigFileInitializer.
//
// Parameters:
//   - path: Path to the config file (e.g., "CLAUDE.md")
//   - renderer: Template renderer function (e.g., RenderInstructionPointer)
//
// Returns:
//   - *ConfigFileInitializer: A new config file initializer
func NewConfigFileInitializer(
	path string,
	renderer TemplateRenderer,
) *ConfigFileInitializer {
	return &ConfigFileInitializer{
		FilePath:   path,
		Renderer:   renderer,
		IsGlobalFs: false,
	}
}

// WithGlobal configures the initializer to use the global filesystem.
func (c *ConfigFileInitializer) WithGlobal(global bool) *ConfigFileInitializer {
	c.IsGlobalFs = global

	return c
}

// Init creates or updates the config file with content between markers.
//
// If the file doesn't exist:
//   - Creates the file with markers and content
//   - Reports file as created
//
// If the file exists:
//   - If markers exist: replaces content between markers
//   - If markers don't exist: appends markers and content to end
//   - Reports file as updated (only if content changed)
//
// Parameters:
//   - ctx: Context for cancellation
//   - fs: Filesystem abstraction
//   - cfg: Configuration with SpectrDir and derived paths
//   - tm: TemplateManager for rendering content
//
// Returns:
//   - InitResult: Contains created or updated file path
//   - error: Non-nil if initialization fails
func (c *ConfigFileInitializer) Init(
	_ context.Context,
	fs afero.Fs,
	cfg *Config,
	tm TemplateManager,
) (InitResult, error) {
	// Render the template content
	content, err := c.renderContent(cfg, tm)
	if err != nil {
		return InitResult{}, err
	}

	// Ensure parent directory exists
	if err := c.ensureParentDir(fs); err != nil {
		return InitResult{}, err
	}

	// Check if file exists
	exists, err := afero.Exists(fs, c.FilePath)
	if err != nil {
		return InitResult{}, fmt.Errorf(
			"failed to check if file exists: %w",
			err,
		)
	}

	if !exists {
		return c.createNewFile(fs, content)
	}

	return c.updateExistingFile(fs, content)
}

// renderContent renders the template content for this config file.
func (c *ConfigFileInitializer) renderContent(
	cfg *Config,
	tm TemplateManager,
) (string, error) {
	templateCtx := TemplateContext{
		BaseDir:     cfg.SpectrDir,
		SpecsDir:    cfg.SpecsDir(),
		ChangesDir:  cfg.ChangesDir(),
		ProjectFile: cfg.ProjectFile(),
		AgentsFile:  cfg.AgentsFile(),
	}

	content, err := c.Renderer(tm, templateCtx)
	if err != nil {
		return "", fmt.Errorf(
			"failed to render template: %w",
			err,
		)
	}

	return content, nil
}

// ensureParentDir creates the parent directory if needed.
func (c *ConfigFileInitializer) ensureParentDir(fs afero.Fs) error {
	dir := filepath.Dir(c.FilePath)
	if dir != "." && dir != "/" {
		if err := fs.MkdirAll(dir, dirPerm); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	return nil
}

// createNewFile creates a new file with markers and content.
func (c *ConfigFileInitializer) createNewFile(
	fs afero.Fs,
	content string,
) (InitResult, error) {
	var result InitResult

	newContent := SpectrStartMarker + "\n" + content + "\n" +
		SpectrEndMarker + "\n"

	if err := afero.WriteFile(
		fs,
		c.FilePath,
		[]byte(newContent),
		filePerm,
	); err != nil {
		return InitResult{}, fmt.Errorf(
			"failed to create file: %w",
			err,
		)
	}

	result.CreatedFiles = append(result.CreatedFiles, c.FilePath)

	return result, nil
}

// updateExistingFile updates an existing file with new content.
func (c *ConfigFileInitializer) updateExistingFile(
	fs afero.Fs,
	content string,
) (InitResult, error) {
	var result InitResult

	existingContent, err := afero.ReadFile(fs, c.FilePath)
	if err != nil {
		return InitResult{}, fmt.Errorf(
			"failed to read file: %w",
			err,
		)
	}

	contentStr := string(existingContent)
	newContent := c.mergeContent(contentStr, content)

	// Only write if content changed
	if newContent != contentStr {
		if err := afero.WriteFile(
			fs,
			c.FilePath,
			[]byte(newContent),
			filePerm,
		); err != nil {
			return InitResult{}, fmt.Errorf(
				"failed to write file: %w",
				err,
			)
		}

		result.UpdatedFiles = append(result.UpdatedFiles, c.FilePath)
	}

	return result, nil
}

// mergeContent merges new content with existing file content.
// If markers exist, replaces content between them.
// If no markers, appends to end.
func (*ConfigFileInitializer) mergeContent(
	existingContent, newContent string,
) string {
	startIndex := findMarkerIndex(existingContent, SpectrStartMarker, 0)
	endIndex := -1

	if startIndex != -1 {
		searchOffset := startIndex + len(SpectrStartMarker)
		endIndex = findMarkerIndex(
			existingContent,
			SpectrEndMarker,
			searchOffset,
		)
	}

	if startIndex == -1 || endIndex == -1 {
		// No markers found, append to end
		return fmt.Sprintf(
			"%s\n\n%s\n%s\n%s\n",
			existingContent,
			SpectrStartMarker,
			newContent,
			SpectrEndMarker,
		)
	}

	// Replace content between markers
	before := existingContent[:startIndex]
	after := existingContent[endIndex+len(SpectrEndMarker):]

	return fmt.Sprintf(
		"%s%s\n%s\n%s%s",
		before,
		SpectrStartMarker,
		newContent,
		SpectrEndMarker,
		after,
	)
}

// IsSetup returns true if the config file exists AND contains spectr markers.
//
// Parameters:
//   - fs: Filesystem abstraction
//   - cfg: Configuration (not used)
//
// Returns:
//   - bool: True if file exists and has markers
func (c *ConfigFileInitializer) IsSetup(fs afero.Fs, _ *Config) bool {
	exists, err := afero.Exists(fs, c.FilePath)
	if err != nil || !exists {
		return false
	}

	// Read file to check for markers
	content, err := afero.ReadFile(fs, c.FilePath)
	if err != nil {
		return false
	}

	contentStr := string(content)

	return strings.Contains(contentStr, SpectrStartMarker) &&
		strings.Contains(contentStr, SpectrEndMarker)
}

// Path returns the config file path for deduplication.
//
// Returns:
//   - string: The config file path
func (c *ConfigFileInitializer) Path() string {
	return c.FilePath
}

// IsGlobal returns true if this initializer uses the global filesystem.
//
// Returns:
//   - bool: True if using global filesystem, false for project-relative
func (c *ConfigFileInitializer) IsGlobal() bool {
	return c.IsGlobalFs
}
