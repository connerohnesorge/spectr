package initializers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/spf13/afero"
)

// ConfigFileInitializer creates or updates config files with marker-based content.
// Uses <!-- spectr:start --> and <!-- spectr:end --> markers for content sections.
type ConfigFileInitializer struct {
	path     string
	template domain.TemplateRef
}

// NewConfigFileInitializer creates an initializer that creates or updates a config file
// with marker-based content sections.
func NewConfigFileInitializer(path string, template domain.TemplateRef) domain.Initializer {
	return &ConfigFileInitializer{
		path:     path,
		template: template,
	}
}

// Init creates or updates the config file with template content between markers.
// Case-insensitive marker matching for reading, always writes lowercase markers.
//
//nolint:revive // argument-limit - interface signature requires these parameters
func (c *ConfigFileInitializer) Init(
	_ context.Context,
	projectFs, _ afero.Fs,
	cfg *domain.Config,
	_ any,
) (domain.ExecutionResult, error) {
	// Create template context from config
	templateCtx := domain.TemplateContext{
		BaseDir:     cfg.SpectrDir,
		SpecsDir:    cfg.SpecsDir(),
		ChangesDir:  cfg.ChangesDir(),
		ProjectFile: cfg.ProjectFile(),
		AgentsFile:  cfg.AgentsFile(),
	}

	// Render template content
	var buf bytes.Buffer
	if err := c.template.Template.ExecuteTemplate(&buf, c.template.Name, templateCtx); err != nil {
		return domain.ExecutionResult{}, fmt.Errorf(
			"failed to render template %s: %w",
			c.template.Name,
			err,
		)
	}
	newContent := buf.String()

	// Check if file exists
	exists, err := afero.Exists(projectFs, c.path)
	if err != nil {
		return domain.ExecutionResult{}, err
	}

	if !exists {
		// Create new file with markers
		content := startMarker + newline + newContent + newline + endMarker
		if err := afero.WriteFile(projectFs, c.path, []byte(content), filePerm); err != nil {
			return domain.ExecutionResult{}, err
		}

		return domain.ExecutionResult{CreatedFiles: []string{c.path}}, nil
	}

	// Read existing file
	existingBytes, err := afero.ReadFile(projectFs, c.path)
	if err != nil {
		return domain.ExecutionResult{}, err
	}
	existingContent := string(existingBytes)

	// Update content with markers
	updatedContent, err := updateWithMarkers(existingContent, newContent)
	if err != nil {
		return domain.ExecutionResult{}, fmt.Errorf("failed to update %s: %w", c.path, err)
	}

	// Write updated content
	if err := afero.WriteFile(projectFs, c.path, []byte(updatedContent), filePerm); err != nil {
		return domain.ExecutionResult{}, err
	}

	return domain.ExecutionResult{UpdatedFiles: []string{c.path}}, nil
}

// IsSetup returns true if the config file exists with markers.
func (c *ConfigFileInitializer) IsSetup(projectFs, _ afero.Fs, _ *domain.Config) bool {
	exists, err := afero.Exists(projectFs, c.path)
	if err != nil || !exists {
		return false
	}

	content, err := afero.ReadFile(projectFs, c.path)
	if err != nil {
		return false
	}

	// Check for start marker (case-insensitive)
	contentLower := strings.ToLower(string(content))

	return strings.Contains(contentLower, strings.ToLower(startMarker))
}

// DedupeKey returns a unique key for deduplication.
// Exported to allow deduplication from the executor package.
func (c *ConfigFileInitializer) DedupeKey() string {
	return "ConfigFileInitializer:" + filepath.Clean(c.path)
}

// Ensure ConfigFileInitializer implements the Deduplicatable interface.
var _ Deduplicatable = (*ConfigFileInitializer)(nil)

// File permission constants.
const filePerm = 0o644

// Marker constants - always write lowercase.
const (
	startMarker = "<!-- spectr:start -->"
	endMarker   = "<!-- spectr:end -->"
	newline     = "\n"
)

// findMarkerCaseInsensitive returns the index and the actual marker text found.
// Returns -1 if not found.
func findMarkerCaseInsensitive(content, marker string) int {
	lower := strings.ToLower(content)
	lowerMarker := strings.ToLower(marker)

	return strings.Index(lower, lowerMarker)
}

// updateWithMarkers updates content between markers or appends if markers don't exist.
// Case-insensitive marker matching for reading, always writes lowercase markers.
func updateWithMarkers(content, newContent string) (string, error) {
	// Case-insensitive search for existing markers
	startIdx := findMarkerCaseInsensitive(content, startMarker)

	if startIdx == -1 {
		// No start marker - check for orphaned end marker (case-insensitive)
		endIdx := findMarkerCaseInsensitive(content, endMarker)
		if endIdx != -1 {
			return "", errors.New("orphaned end marker without start marker")
		}
		// No markers exist - append new block at end with lowercase markers
		result := content
		if result != "" && !strings.HasSuffix(result, newline) {
			result += newline
		}

		return result + newline + startMarker + newline + newContent + newline + endMarker, nil
	}

	// Start marker found - look for end marker AFTER the start (case-insensitive)
	searchFrom := startIdx + len(startMarker)
	afterStart := content[searchFrom:]
	endIdx := findMarkerCaseInsensitive(afterStart, endMarker)

	if endIdx != -1 {
		// Normal case: both markers present and properly paired
		endIdxAbsolute := searchFrom + endIdx

		// Check for nested start marker before end (case-insensitive)
		betweenMarkers := content[searchFrom:endIdxAbsolute]
		nextStartIdx := findMarkerCaseInsensitive(betweenMarkers, startMarker)
		if nextStartIdx != -1 {
			return "", errors.New("nested start marker before end marker")
		}

		before := content[:startIdx]
		after := content[endIdxAbsolute+len(endMarker):]
		// Always write lowercase markers
		return before + startMarker + newline + newContent + newline + endMarker + after, nil
	}

	// Start marker exists but no end marker immediately after
	// Check for multiple start markers without end (case-insensitive)
	nextStartIdx := findMarkerCaseInsensitive(afterStart, startMarker)
	if nextStartIdx != -1 {
		return "", errors.New("multiple start markers without end markers")
	}

	// No end marker anywhere after start - orphaned start marker
	// Replace everything from start marker onward with new block
	before := content[:startIdx]

	return before + startMarker + newline + newContent + newline + endMarker, nil
}
