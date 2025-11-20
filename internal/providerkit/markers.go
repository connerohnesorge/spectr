package providerkit

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Spectr marker constants for managing config file updates
const (
	SpectrStartMarker = "<!-- spectr:START -->"
	SpectrEndMarker   = "<!-- spectr:END -->"
	newline           = "\n"
)

// UpdateFileWithMarkers updates a file with content between markers.
// If the file doesn't exist, creates it with markers.
// If the file exists, updates content between markers.
//
// This function is idempotent and handles three scenarios:
//  1. File doesn't exist: Creates file with markers and content
//  2. File exists with both markers: Updates content between markers
//  3. File exists without markers: Prepends markers and content
//
// Returns an error if:
//   - Only one marker is found (invalid state)
//   - End marker appears before start marker
//   - File operations fail
func UpdateFileWithMarkers(
	filePath, content, startMarker, endMarker string,
) error {
	expandedPath, err := ExpandPath(filePath)
	if err != nil {
		return fmt.Errorf("failed to expand path: %w", err)
	}

	var existingContent string

	if FileExists(expandedPath) {
		existingContent, err = updateExistingFileMarkers(
			expandedPath,
			content,
			startMarker,
			endMarker,
		)
		if err != nil {
			return err
		}
	} else {
		// File doesn't exist - create with markers
		existingContent = startMarker + newline + content + newline + endMarker
	}

	// Ensure parent directory exists
	dir := filepath.Dir(expandedPath)
	if err := EnsureDir(dir); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(
		expandedPath,
		[]byte(existingContent),
		defaultFilePerm,
	); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// updateExistingFileMarkers updates markers in an existing file
func updateExistingFileMarkers(
	expandedPath, content, startMarker, endMarker string,
) (string, error) {
	// Read existing content
	data, err := os.ReadFile(expandedPath)
	if err != nil {
		return "", fmt.Errorf("failed to read existing file: %w", err)
	}
	existingContent := string(data)

	// Find markers
	startIndex := findMarkerIndex(existingContent, startMarker, 0)
	var endIndex int
	if startIndex != -1 {
		endIndex = findMarkerIndex(
			existingContent,
			endMarker,
			startIndex+len(startMarker),
		)
	} else {
		endIndex = findMarkerIndex(existingContent, endMarker, 0)
	}

	// Handle different marker states
	switch {
	case startIndex != -1 && endIndex != -1:
		// Both markers found - update content between them
		if endIndex < startIndex {
			return "", fmt.Errorf(
				"invalid marker state in %s: end marker before start",
				expandedPath,
			)
		}

		before := existingContent[:startIndex]
		after := existingContent[endIndex+len(endMarker):]
		existingContent = before + startMarker + newline +
			content + newline + endMarker + after
	case startIndex == -1 && endIndex == -1:
		// No markers found - prepend with markers
		existingContent = startMarker + newline + content + newline +
			endMarker + newline + newline + existingContent
	default:
		// Only one marker found - error
		return "", fmt.Errorf(
			"invalid marker state in %s: found start: %t, found end: %t",
			expandedPath,
			startIndex != -1,
			endIndex != -1,
		)
	}

	return existingContent, nil
}

// findMarkerIndex finds the index of a marker that is on its own line.
// Returns -1 if not found.
//
// This function ensures that markers are only recognized when they appear
// on their own line (with only whitespace around them), preventing
// false matches in commented or quoted strings.
func findMarkerIndex(content, marker string, fromIndex int) int {
	currentIndex := strings.Index(content[fromIndex:], marker)
	if currentIndex == -1 {
		return -1
	}
	currentIndex += fromIndex

	for currentIndex != -1 {
		if isMarkerOnOwnLine(content, currentIndex, len(marker)) {
			return currentIndex
		}

		nextIndex := strings.Index(content[currentIndex+len(marker):], marker)
		if nextIndex == -1 {
			return -1
		}
		currentIndex = currentIndex + len(marker) + nextIndex
	}

	return -1
}

// isMarkerOnOwnLine checks if a marker is on its own line
// (only whitespace around it). This prevents false matches when
// the marker text appears in comments or strings.
func isMarkerOnOwnLine(content string, markerIndex, markerLength int) bool {
	// Check left side
	leftIndex := markerIndex - 1
	for leftIndex >= 0 && content[leftIndex] != '\n' {
		char := content[leftIndex]
		if char != ' ' && char != '\t' && char != '\r' {
			return false
		}
		leftIndex--
	}

	// Check right side
	rightIndex := markerIndex + markerLength
	for rightIndex < len(content) && content[rightIndex] != '\n' {
		char := content[rightIndex]
		if char != ' ' && char != '\t' && char != '\r' {
			return false
		}
		rightIndex++
	}

	return true
}
