package providers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	// File and directory permission constants.
	dirPerm  = 0o755
	filePerm = 0o644

	// Marker constants for managing config file updates.
	SpectrStartMarker = "<!-- spectr:START -->"
	SpectrEndMarker   = "<!-- spectr:END -->"

	// Common strings.
	newline       = "\n"
	newlineDouble = "\n\n"
)

// FileExists checks if a file or directory exists.
func FileExists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

// EnsureDir creates a directory and all parent directories if they don't exist.
func EnsureDir(path string) error {
	return os.MkdirAll(path, dirPerm)
}

// UpdateFileWithMarkers updates content between markers in a file,
// or creates the file with markers if it doesn't exist.
func UpdateFileWithMarkers(
	filePath, content, startMarker, endMarker string,
) error {
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := EnsureDir(dir); err != nil {
		return fmt.Errorf(
			"failed to create directory: %w",
			err,
		)
	}

	// Check if file exists
	if !FileExists(filePath) {
		// Create new file with markers
		newContent := startMarker + "\n" + content + "\n" + endMarker + "\n"

		return os.WriteFile(
			filePath,
			[]byte(newContent),
			filePerm,
		)
	}

	// Read existing file
	existingContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf(
			"failed to read file: %w",
			err,
		)
	}

	contentStr := string(existingContent)

	// Find markers
	startIndex := findMarkerIndex(
		contentStr,
		startMarker,
		0,
	)
	endIndex := -1
	if startIndex != -1 {
		searchOffset := startIndex + len(
			startMarker,
		)
		endIndex = findMarkerIndex(
			contentStr,
			endMarker,
			searchOffset,
		)
	}

	var newContent string
	if startIndex == -1 || endIndex == -1 {
		// No markers found, append to end
		newContent = contentStr + newlineDouble +
			startMarker + newline + content + newline + endMarker + newline
	} else {
		// Replace content between markers
		before := contentStr[:startIndex]
		after := contentStr[endIndex+len(endMarker):]
		newContent = before + startMarker + newline +
			content + newline + endMarker + after
	}

	return os.WriteFile(
		filePath,
		[]byte(newContent),
		filePerm,
	)
}

// findMarkerIndex finds the index of a marker in content, starting from offset.
func findMarkerIndex(
	content, marker string,
	offset int,
) int {
	if offset >= len(content) {
		return -1
	}
	idx := strings.Index(content[offset:], marker)
	if idx == -1 {
		return -1
	}

	return offset + idx
}
