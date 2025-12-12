package initialize

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// markerBlockSuffix enforces consistent spacing after marker blocks.
const markerBlockSuffix = "\n\n"

// UpdateFileWithMarkers updates a file with content between markers.
// If the file doesn't exist, creates it with markers.
// If the file exists, updates content between markers.
//
// Marker behavior:
//   - Both markers found: Updates content between them
//   - No markers found: Prepends content with markers
//   - Only one marker found: Returns error (invalid state)
//
// The markers must be on their own line (only whitespace around them).
//
// Example usage:
//
//	err := UpdateFileWithMarkers(
//	    "CLAUDE.md",
//	    "# Instructions\nFollow these rules...",
//	    SpectrStartMarker,
//	    SpectrEndMarker,
//	)
//
//nolint:revive // line-length-limit - string literals need clarity
func UpdateFileWithMarkers(
	filePath, content, startMarker, endMarker string,
) error {
	expandedPath, err := ExpandPath(filePath)
	if err != nil {
		return fmt.Errorf(
			"failed to expand path: %w",
			err,
		)
	}

	var existingContent string

	if FileExists(expandedPath) {
		// Read existing content
		data, err := os.ReadFile(expandedPath)
		if err != nil {
			return fmt.Errorf(
				"failed to read existing file: %w",
				err,
			)
		}
		existingContent = string(data)

		// Find markers
		startIndex := findMarkerIndex(
			existingContent,
			startMarker,
			0,
		)
		var endIndex int
		if startIndex != -1 {
			startPos := startIndex + len(
				startMarker,
			)
			endIndex = findMarkerIndex(
				existingContent,
				endMarker,
				startPos,
			)
		} else {
			endIndex = findMarkerIndex(existingContent, endMarker, 0)
		}

		// Handle different marker states
		switch {
		case startIndex != -1 && endIndex != -1:
			// Both markers found - update content between them
			if endIndex < startIndex {
				return fmt.Errorf(
					"invalid marker state in %s: "+
						"end marker appears before start marker",
					expandedPath,
				)
			}

			before := existingContent[:startIndex]
			after := existingContent[endIndex+len(endMarker):]
			newContent := before + startMarker + "\n" + content +
				"\n" + endMarker + after
			existingContent = newContent
		case startIndex == -1 && endIndex == -1:
			// No markers found - prepend with markers
			newContent := startMarker + "\n" + content + "\n" +
				endMarker + markerBlockSuffix + existingContent
			existingContent = newContent
		default:
			// Only one marker found - error
			return fmt.Errorf(
				"invalid marker state in %s: found start: %t, found end: %t",
				expandedPath,
				startIndex != -1,
				endIndex != -1,
			)
		}
	} else {
		// File doesn't exist - create with markers
		existingContent = startMarker + "\n" + content + "\n" +
			endMarker + markerBlockSuffix
	}

	// Ensure parent directory exists
	dir := filepath.Dir(expandedPath)
	if err := EnsureDir(dir); err != nil {
		return fmt.Errorf(
			"failed to create parent directory: %w",
			err,
		)
	}

	// Write file
	fileData := []byte(existingContent)
	if err := os.WriteFile(expandedPath, fileData, filePerm); err != nil {
		return fmt.Errorf(
			"failed to write file: %w",
			err,
		)
	}

	return nil
}

// findMarkerIndex finds the index of a marker that is on its own line.
// Only returns markers that have only whitespace around them on the same line.
// Returns -1 if not found.
//
// Parameters:
//   - content: The file content to search
//   - marker: The marker string to find
//   - fromIndex: The index to start searching from
func findMarkerIndex(
	content, marker string,
	fromIndex int,
) int {
	currentIndex := strings.Index(
		content[fromIndex:],
		marker,
	)
	if currentIndex == -1 {
		return -1
	}
	currentIndex += fromIndex

	for currentIndex != -1 {
		if isMarkerOnOwnLine(
			content,
			currentIndex,
			len(marker),
		) {
			return currentIndex
		}

		nextIndex := strings.Index(
			content[currentIndex+len(marker):],
			marker,
		)
		if nextIndex == -1 {
			return -1
		}
		currentIndex = currentIndex + len(
			marker,
		) + nextIndex
	}

	return -1
}

// isMarkerOnOwnLine checks if a marker is on its own line.
// Returns true if there is only whitespace (spaces, tabs, \r) around the marker
// on the same line.
//
// Parameters:
//   - content: The file content
//   - markerIndex: The index where the marker starts
//   - markerLength: The length of the marker string
func isMarkerOnOwnLine(
	content string,
	markerIndex, markerLength int,
) bool {
	// Check left side (from marker backwards to newline)
	leftIndex := markerIndex - 1
	for leftIndex >= 0 && content[leftIndex] != '\n' {
		char := content[leftIndex]
		if char != ' ' && char != '\t' &&
			char != '\r' {
			return false
		}
		leftIndex--
	}

	// Check right side (from marker forwards to newline)
	rightIndex := markerIndex + markerLength
	for rightIndex < len(content) && content[rightIndex] != '\n' {
		char := content[rightIndex]
		if char != ' ' && char != '\t' &&
			char != '\r' {
			return false
		}
		rightIndex++
	}

	return true
}
