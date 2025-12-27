package providers

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/connerohnesorge/spectr/internal/templates"
	"github.com/spf13/afero"
)

// ConfigFileInitializer creates or updates instruction files with
// marker-based updates. Uses case-insensitive marker matching for reading,
// always writes lowercase markers.
type ConfigFileInitializer struct {
	path     string             // file path (relative to project root)
	template domain.TemplateRef // template reference for rendering content
}

// NewConfigFileInitializer creates a new ConfigFileInitializer.
func NewConfigFileInitializer(
	path string,
	template domain.TemplateRef,
) *ConfigFileInitializer {
	return &ConfigFileInitializer{
		path:     path,
		template: template,
	}
}

// Init creates or updates the config file with marker-based updates.
//
//nolint:revive // argument-limit: interface requires 5 args for dual-fs support
func (c *ConfigFileInitializer) Init(
	_ context.Context,
	projectFs, _ afero.Fs,
	cfg *Config,
	_ *templates.TemplateManager,
) (InitResult, error) {
	// Create template context from config
	templateCtx := domain.TemplateContext{
		BaseDir:     cfg.SpectrDir,
		SpecsDir:    cfg.SpecsDir(),
		ChangesDir:  cfg.ChangesDir(),
		ProjectFile: cfg.ProjectFile(),
		AgentsFile:  cfg.AgentsFile(),
	}

	// Render template
	content, err := c.template.Render(templateCtx)
	if err != nil {
		return InitResult{}, fmt.Errorf("failed to render template: %w", err)
	}

	// Check if file exists
	exists, err := afero.Exists(projectFs, c.path)
	if err != nil {
		return InitResult{}, fmt.Errorf("failed to check file existence: %w", err)
	}

	if !exists {
		// Create new file with markers
		startMarker := "<!-- spectr:start -->"
		endMarker := "<!-- spectr:end -->"
		newContent := startMarker + "\n" + content + "\n" + endMarker + "\n"

		if err := afero.WriteFile(projectFs, c.path, []byte(newContent), 0644); err != nil { //nolint:revive
			return InitResult{}, fmt.Errorf("failed to create file: %w", err)
		}

		return InitResult{
			CreatedFiles: []string{c.path},
			UpdatedFiles: nil,
		}, nil
	}

	// File exists - update with markers
	existingBytes, err := afero.ReadFile(projectFs, c.path)
	if err != nil {
		return InitResult{}, fmt.Errorf("failed to read file: %w", err)
	}

	existingContent := string(existingBytes)
	updatedContent, err := updateWithMarkers(existingContent, content)
	if err != nil {
		return InitResult{}, err
	}

	// Only write if content changed
	if updatedContent != existingContent {
		if err := afero.WriteFile(projectFs, c.path, []byte(updatedContent), 0644); err != nil { //nolint:revive
			return InitResult{}, fmt.Errorf("failed to update file: %w", err)
		}

		return InitResult{
			CreatedFiles: nil,
			UpdatedFiles: []string{c.path},
		}, nil
	}

	// No changes needed
	return InitResult{
		CreatedFiles: nil,
		UpdatedFiles: nil,
	}, nil
}

// IsSetup returns true if the config file exists in the project filesystem.
func (c *ConfigFileInitializer) IsSetup(projectFs, _ afero.Fs, _ *Config) bool {
	exists, err := afero.Exists(projectFs, c.path)

	return err == nil && exists
}

// dedupeKey returns a unique key for deduplication.
// Format: "ConfigFileInitializer:<path>"
func (c *ConfigFileInitializer) dedupeKey() string {
	return fmt.Sprintf("ConfigFileInitializer:%s", filepath.Clean(c.path))
}

// updateWithMarkers updates file content between spectr markers.
// Marker matching is case-insensitive for reading, always writes lowercase markers.
//
//nolint:revive // function-length: complex marker handling requires multiple branches
func updateWithMarkers(content, newContent string) (string, error) {
	// Always write lowercase markers
	startMarker := "<!-- spectr:start -->"
	endMarker := "<!-- spectr:end -->"

	// Case-insensitive search for existing markers
	startIdx, startLen := findMarkerCaseInsensitive(content, startMarker)

	if startIdx == -1 {
		// No start marker - check for orphaned end marker (case-insensitive)
		endIdx, _ := findMarkerCaseInsensitive(content, endMarker)
		if endIdx != -1 {
			return "", fmt.Errorf("orphaned end marker at position %d without start marker", endIdx)
		}
		// No markers exist - append new block at end with lowercase markers
		result := content
		if !strings.HasSuffix(result, "\n") && len(result) > 0 {
			result += "\n"
		}

		return result + "\n" + startMarker + "\n" + newContent + "\n" + endMarker, nil
	}

	// Start marker found - look for end marker AFTER the start (case-insensitive)
	searchFrom := startIdx + startLen
	endIdx, endLen := findMarkerCaseInsensitive(content[searchFrom:], endMarker)

	if endIdx != -1 {
		// Normal case: both markers present and properly paired
		endIdx += searchFrom // Adjust to absolute position

		// Check for nested start marker before end (case-insensitive)
		nextStartIdx, _ := findMarkerCaseInsensitive(content[searchFrom:endIdx], startMarker)
		if nextStartIdx != -1 {
			return "", fmt.Errorf(
				"nested start marker at position %d before end marker at %d",
				searchFrom+nextStartIdx,
				endIdx,
			)
		}

		before := content[:startIdx]
		after := content[endIdx+endLen:]
		// Always write lowercase markers
		return before + startMarker + "\n" + newContent + "\n" + endMarker + after, nil
	}

	// Start marker exists but no end marker immediately after
	// Search for trailing end marker anywhere in the file (case-insensitive)
	trailingEndIdx, trailingEndLen := findMarkerCaseInsensitive(content[searchFrom:], endMarker)
	if trailingEndIdx != -1 {
		trailingEndIdx += searchFrom
		// Found end marker after start - use it
		before := content[:startIdx]
		after := content[trailingEndIdx+trailingEndLen:]
		// Always write lowercase markers
		return before + startMarker + "\n" + newContent + "\n" + endMarker + after, nil
	}

	// Check for multiple start markers without end (case-insensitive)
	nextStartIdx, _ := findMarkerCaseInsensitive(content[searchFrom:], startMarker)
	if nextStartIdx != -1 {
		return "", fmt.Errorf(
			"multiple start markers at positions %d and %d without end markers",
			startIdx,
			searchFrom+nextStartIdx,
		)
	}

	// No end marker anywhere after start - orphaned start marker
	// Replace everything from start marker onward with new block
	before := content[:startIdx]

	return before + startMarker + "\n" + newContent + "\n" + endMarker, nil
}

// findMarkerCaseInsensitive returns the index and length of a marker using
// case-insensitive matching. Returns -1, 0 if marker is not found.
func findMarkerCaseInsensitive(content, marker string) (index, length int) {
	lower := strings.ToLower(content)
	lowerMarker := strings.ToLower(marker)
	idx := strings.Index(lower, lowerMarker)
	if idx == -1 {
		return -1, 0
	}
	// Return the actual length of the matched text
	// (preserves original case). We need to find the actual marker in the
	// original content at this position.
	actualMarker := content[idx : idx+len(marker)]

	return idx, len(actualMarker)
}
