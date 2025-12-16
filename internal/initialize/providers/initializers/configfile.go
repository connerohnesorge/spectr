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

const (
	// File permission for created files.
	filePerm = 0644

	// Marker constants for managing config file updates.
	spectrStartMarker = "<!-- spectr:START -->"
	spectrEndMarker   = "<!-- spectr:END -->"

	// Common strings.
	newline       = "\n"
	newlineDouble = "\n\n"
)

// Compile-time interface satisfaction check.
var _ providers.Initializer = (*ConfigFileInitializer)(nil)

// ConfigFileInitializer creates or updates config files with marker-based
// content. It is idempotent - running Init multiple times produces the
// same result.
//
// The initializer handles three scenarios:
//  1. File doesn't exist: Create with markers around content
//  2. File exists without markers: Append content with markers at end
//  3. File exists with markers: Replace content between markers
type ConfigFileInitializer struct {
	// Path is the config file path relative to the project root.
	Path string

	// Template is the content to insert between markers.
	Template string
}

// NewConfigFileInitializer creates a new ConfigFileInitializer for the given
// path and template. Path should be relative to the project root.
func NewConfigFileInitializer(path, template string) *ConfigFileInitializer {
	return &ConfigFileInitializer{Path: path, Template: template}
}

// Init creates or updates the config file with marker-based content.
func (c *ConfigFileInitializer) Init(
	_ context.Context,
	fs afero.Fs,
	_ *providers.Config,
) error {
	if err := c.ensureParentDir(fs); err != nil {
		return err
	}

	exists, err := afero.Exists(fs, c.Path)
	if err != nil {
		return fmt.Errorf(
			"failed to check if file exists %s: %w",
			c.Path,
			err,
		)
	}

	if !exists {
		return c.createNewFile(fs)
	}

	return c.updateExistingFile(fs)
}

// ensureParentDir creates the parent directory if needed.
func (c *ConfigFileInitializer) ensureParentDir(fs afero.Fs) error {
	dir := filepath.Dir(c.Path)
	if dir == "." || dir == "" {
		return nil
	}
	if err := fs.MkdirAll(dir, dirPerm); err != nil {
		return fmt.Errorf(
			"failed to create directory %s: %w",
			dir,
			err,
		)
	}

	return nil
}

// createNewFile creates a new config file with markers.
func (c *ConfigFileInitializer) createNewFile(fs afero.Fs) error {
	content := spectrStartMarker + newline +
		c.Template + newline + spectrEndMarker + newline
	err := afero.WriteFile(fs, c.Path, []byte(content), filePerm)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", c.Path, err)
	}

	return nil
}

// updateExistingFile updates an existing config file with markers.
func (c *ConfigFileInitializer) updateExistingFile(fs afero.Fs) error {
	existingContent, err := afero.ReadFile(fs, c.Path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", c.Path, err)
	}

	newContent := c.generateUpdatedContent(string(existingContent))
	err = afero.WriteFile(fs, c.Path, []byte(newContent), filePerm)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", c.Path, err)
	}

	return nil
}

// generateUpdatedContent generates updated content with markers.
func (c *ConfigFileInitializer) generateUpdatedContent(
	contentStr string,
) string {
	startIndex := findMarkerIndex(contentStr, spectrStartMarker, 0)
	endIndex := -1
	if startIndex != -1 {
		offset := startIndex + len(spectrStartMarker)
		endIndex = findMarkerIndex(contentStr, spectrEndMarker, offset)
	}

	if startIndex == -1 || endIndex == -1 {
		return contentStr + newlineDouble + spectrStartMarker + newline +
			c.Template + newline + spectrEndMarker + newline
	}

	before := contentStr[:startIndex]
	after := contentStr[endIndex+len(spectrEndMarker):]

	return before + spectrStartMarker + newline +
		c.Template + newline + spectrEndMarker + after
}

// IsSetup returns true if the config file exists and contains both markers.
func (c *ConfigFileInitializer) IsSetup(
	fs afero.Fs,
	_ *providers.Config,
) bool {
	exists, err := afero.Exists(fs, c.Path)
	if err != nil || !exists {
		return false
	}

	content, err := afero.ReadFile(fs, c.Path)
	if err != nil {
		return false
	}

	contentStr := string(content)

	return strings.Contains(contentStr, spectrStartMarker) &&
		strings.Contains(contentStr, spectrEndMarker)
}

// Key returns a unique key for this initializer based on its configuration.
func (c *ConfigFileInitializer) Key() string {
	return "config:" + c.Path
}

// findMarkerIndex finds the index of a marker in content, starting from
// offset. Returns -1 if the marker is not found.
func findMarkerIndex(content, marker string, offset int) int {
	if offset >= len(content) {
		return -1
	}
	idx := strings.Index(content[offset:], marker)
	if idx == -1 {
		return -1
	}

	return offset + idx
}
