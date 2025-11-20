//nolint:revive // line-length-limit,file-length-limit,add-constant,unused-receiver - readability over strict formatting
package init

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// spectr_MARKERS for managing config file updates
const (
	SpectrStartMarker = "<!-- spectr:START -->"
	SpectrEndMarker   = "<!-- spectr:END -->"
)

// UpdateFileWithMarkers updates a file with content between markers
// If file doesn't exist, creates it with markers
// If file exists, updates content between markers
func UpdateFileWithMarkers(filePath, content, startMarker, endMarker string) error {
	expandedPath, err := ExpandPath(filePath)
	if err != nil {
		return fmt.Errorf("failed to expand path: %w", err)
	}

	var existingContent string

	if FileExists(expandedPath) {
		// Read existing content
		data, err := os.ReadFile(expandedPath)
		if err != nil {
			return fmt.Errorf("failed to read existing file: %w", err)
		}
		existingContent = string(data)

		// Find markers
		startIndex := findMarkerIndex(existingContent, startMarker, 0)
		var endIndex int
		if startIndex != -1 {
			endIndex = findMarkerIndex(existingContent, endMarker, startIndex+len(startMarker))
		} else {
			endIndex = findMarkerIndex(existingContent, endMarker, 0)
		}

		// Handle different marker states
		switch {
		case startIndex != -1 && endIndex != -1:
			// Both markers found - update content between them
			if endIndex < startIndex {
				return fmt.Errorf(
					"invalid marker state in %s: end marker appears before start marker",
					expandedPath,
				)
			}

			before := existingContent[:startIndex]
			after := existingContent[endIndex+len(endMarker):]
			existingContent = before + startMarker + "\n" + content + "\n" + endMarker + after
		case startIndex == -1 && endIndex == -1:
			// No markers found - prepend with markers
			existingContent = startMarker + "\n" + content + "\n" + endMarker + "\n\n" + existingContent
		default:
			// Only one marker found - error
			return fmt.Errorf("invalid marker state in %s: found start: %t, found end: %t",
				expandedPath, startIndex != -1, endIndex != -1)
		}
	} else {
		// File doesn't exist - create with markers
		existingContent = startMarker + "\n" + content + "\n" + endMarker
	}

	// Ensure parent directory exists
	dir := filepath.Dir(expandedPath)
	if err := EnsureDir(dir); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(expandedPath, []byte(existingContent), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// findMarkerIndex finds the index of a marker that is on its own line
// Returns -1 if not found
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

// isMarkerOnOwnLine checks if a marker is on its own line (only whitespace around it)
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

// ============================================================================
// Legacy code removed - now using ToolProvider pattern
// See interfaces.go, providers.go, slash_providers.go, and
// composite_providers.go
// ============================================================================
