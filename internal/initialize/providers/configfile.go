package providers

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/spf13/afero"
)

// newline is the newline character used in marker formatting.
const newline = "\n"

// ConfigFileInitializer creates or updates instruction files with marker-based content. //nolint:lll
// Uses <!-- spectr:start --> and <!-- spectr:end --> markers (case-insensitive read, lowercase write). //nolint:lll
type ConfigFileInitializer struct {
	path     string             // Relative path from project root
	template domain.TemplateRef // Template to render between markers
}

// NewConfigFileInitializer creates an initializer for instruction files.
// Path is relative to the project filesystem root.
// Template is rendered with TemplateContext and inserted between markers.
// Example: NewConfigFileInitializer("CLAUDE.md", tm.InstructionPointer())
func NewConfigFileInitializer(path string, template domain.TemplateRef) Initializer {
	return &ConfigFileInitializer{
		path:     path,
		template: template,
	}
}

// Init creates or updates the instruction file.
// If file doesn't exist, creates it with rendered content between markers.

// If file exists, updates content between markers (case-insensitive marker search).
//
//nolint:revive // Init signature is defined by Initializer interface
func (c *ConfigFileInitializer) Init(
	ctx context.Context,
	projectFs, homeFs afero.Fs,
	cfg *Config,
	tm TemplateManager,
) (InitResult, error) {
	// Derive template context from config
	tmplCtx := &domain.TemplateContext{
		BaseDir:     cfg.SpectrDir,
		SpecsDir:    cfg.SpecsDir(),
		ChangesDir:  cfg.ChangesDir(),
		ProjectFile: cfg.ProjectFile(),
		AgentsFile:  cfg.AgentsFile(),
	}

	// Render template content
	content, err := c.template.Render(tmplCtx)
	if err != nil {
		return InitResult{}, fmt.Errorf("failed to render template: %w", err)
	}

	// Check if file exists
	exists, err := afero.Exists(projectFs, c.path)
	if err != nil {
		return InitResult{}, fmt.Errorf("failed to check file %s: %w", c.path, err)
	}

	if !exists {
		// Create new file with markers
		newContent := fmt.Sprintf("<!-- spectr:start -->\n%s\n<!-- spectr:end -->", content)
		if err := afero.WriteFile(projectFs, c.path, []byte(newContent), 0o644); err != nil {
			return InitResult{}, fmt.Errorf("failed to create file %s: %w", c.path, err)
		}

		return InitResult{
			CreatedFiles: []string{c.path},
			UpdatedFiles: nil,
		}, nil
	}

	// Read existing file
	data, err := afero.ReadFile(projectFs, c.path)
	if err != nil {
		return InitResult{}, fmt.Errorf("failed to read file %s: %w", c.path, err)
	}

	// Update content between markers
	updated, err := updateWithMarkers(string(data), content)
	if err != nil {
		return InitResult{}, fmt.Errorf("failed to update markers in %s: %w", c.path, err)
	}

	// Write updated content
	if err := afero.WriteFile(projectFs, c.path, []byte(updated), 0o644); err != nil {
		return InitResult{}, fmt.Errorf("failed to write file %s: %w", c.path, err)
	}

	return InitResult{
		CreatedFiles: nil,
		UpdatedFiles: []string{c.path},
	}, nil
}

// IsSetup checks if the instruction file exists in the project filesystem.
func (c *ConfigFileInitializer) IsSetup(projectFs, _ afero.Fs, _ *Config) bool {
	exists, err := afero.Exists(projectFs, c.path)

	return err == nil && exists
}

// dedupeKey returns a unique key for deduplication.
// Uses type name + normalized path to prevent duplicate config file updates.
func (c *ConfigFileInitializer) dedupeKey() string {
	return fmt.Sprintf("ConfigFileInitializer:%s", filepath.Clean(c.path))
}

// updateWithMarkers updates content between markers.
// Marker matching is case-insensitive for reading, always writes lowercase.
// Edge cases handled:
// - Missing markers: append new block at end with markers
// - Orphaned end marker: error (corrupted file)
// - Nested markers: error (not supported)
// - Multiple start markers: error (ambiguous)
// - Orphaned start marker: replace from start marker to end of file
// - Normal case: replace content between existing markers
func updateWithMarkers(existing, newContent string) (string, error) {
	const (
		startMarker = "<!-- spectr:start -->"
		endMarker   = "<!-- spectr:end -->"
	)

	startIdx, startLen := findMarkerCaseInsensitive(existing, startMarker)

	if startIdx == -1 {
		return handleNoStartMarker(existing, newContent, startMarker, endMarker)
	}

	return handleStartMarkerFound(
		existing,
		newContent,
		startMarker,
		endMarker,
		startIdx,
		startLen,
	)
}

// handleNoStartMarker handles the case when no start marker is found.
func handleNoStartMarker(
	existing, newContent, startMarker, endMarker string,
) (string, error) {
	// Check for orphaned end marker
	endIdx, _ := findMarkerCaseInsensitive(existing, endMarker)
	if endIdx != -1 {
		return "", fmt.Errorf(
			"orphaned end marker at position %d without start marker",
			endIdx,
		)
	}

	// No markers exist - append new block at end with lowercase markers
	result := existing
	if existing != "" && !strings.HasSuffix(existing, newline) {
		result += newline
	}

	return result + newline + startMarker + newline + newContent + newline + endMarker, nil
}

// handleStartMarkerFound handles the case when a start marker exists.
//
//nolint:revive // Helper function needs multiple params for clarity
func handleStartMarkerFound(
	existing, newContent, startMarker, endMarker string,
	startIdx, startLen int,
) (string, error) {
	searchFrom := startIdx + startLen
	endIdx, endLen := findMarkerCaseInsensitive(existing[searchFrom:], endMarker)

	if endIdx != -1 {
		return replaceMarkedSection(
			existing,
			newContent,
			startMarker,
			endMarker,
			startIdx,
			searchFrom,
			endIdx,
			endLen,
		)
	}

	return handleOrphanedStartMarker(
		existing,
		newContent,
		startMarker,
		endMarker,
		startIdx,
		searchFrom,
	)
}

// replaceMarkedSection replaces content between properly paired markers.
//
//nolint:revive // Helper function needs params; modifies endIdx locally
func replaceMarkedSection(
	existing, newContent, startMarker, endMarker string,
	startIdx, searchFrom, endIdx, endLen int,
) (string, error) {
	endIdx += searchFrom // Adjust to absolute position

	// Check for nested start marker before end
	nextStartIdx, _ := findMarkerCaseInsensitive(
		existing[searchFrom:endIdx],
		startMarker,
	)
	if nextStartIdx != -1 {
		return "", fmt.Errorf(
			"nested start marker at position %d before end marker at %d",
			searchFrom+nextStartIdx,
			endIdx,
		)
	}

	before := existing[:startIdx]
	after := existing[endIdx+endLen:]

	return before + startMarker + "\n" + newContent + "\n" + endMarker + after, nil
}

// handleOrphanedStartMarker handles start marker without matching end.
//
//nolint:revive // Helper function needs multiple params for clarity
func handleOrphanedStartMarker(
	existing, newContent, startMarker, endMarker string,
	startIdx, searchFrom int,
) (string, error) {
	// Check for multiple start markers without end
	nextStartIdx, _ := findMarkerCaseInsensitive(existing[searchFrom:], startMarker)
	if nextStartIdx != -1 {
		return "", fmt.Errorf(
			"multiple start markers at positions %d and %d without end markers",
			startIdx,
			searchFrom+nextStartIdx,
		)
	}

	// Replace everything from start marker onward with new block
	before := existing[:startIdx]

	return before + startMarker + "\n" + newContent + "\n" + endMarker, nil
}

// findMarkerCaseInsensitive finds a marker in content using case-insensitive matching.
// Returns the index and length of the matched marker, or (-1, 0) if not found.
func findMarkerCaseInsensitive(content, marker string) (index, length int) {
	lower := strings.ToLower(content)
	lowerMarker := strings.ToLower(marker)
	idx := strings.Index(lower, lowerMarker)
	if idx == -1 {
		return -1, 0
	}

	return idx, len(marker)
}
