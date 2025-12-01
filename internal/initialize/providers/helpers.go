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
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Check if file exists
	if !FileExists(filePath) {
		// Create new file with markers
		newContent := startMarker + "\n" + content + "\n" + endMarker + "\n"

		return os.WriteFile(filePath, []byte(newContent), filePerm)
	}

	// Read existing file
	existingContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	contentStr := string(existingContent)

	// Find markers
	startIndex := findMarkerIndex(contentStr, startMarker, 0)
	endIndex := -1
	if startIndex != -1 {
		searchOffset := startIndex + len(startMarker)
		endIndex = findMarkerIndex(contentStr, endMarker, searchOffset)
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

	return os.WriteFile(filePath, []byte(newContent), filePerm)
}

// findMarkerIndex finds the index of a marker in content, starting from offset.
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

// updateSlashCommandBody updates the body of a slash command file.
func updateSlashCommandBody(filePath, body, frontmatter string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	contentStr := string(content)

	startIndex := findMarkerIndex(contentStr, SpectrStartMarker, 0)
	if startIndex == -1 {
		return fmt.Errorf("start marker not found in %s", filePath)
	}

	searchOffset := startIndex + len(SpectrStartMarker)
	endIndex := findMarkerIndex(contentStr, SpectrEndMarker, searchOffset)
	if endIndex == -1 {
		return fmt.Errorf("end marker not found in %s", filePath)
	}

	if endIndex < startIndex {
		return fmt.Errorf(
			"end marker appears before start marker in %s", filePath)
	}

	before := contentStr[:startIndex]
	after := contentStr[endIndex+len(SpectrEndMarker):]

	hasFrontmatter := strings.HasPrefix(strings.TrimSpace(before), "---")
	if frontmatter != "" && !hasFrontmatter {
		before = strings.TrimSpace(frontmatter) + newlineDouble +
			strings.TrimLeft(before, "\n\r")
	}

	newContent := before + SpectrStartMarker + newline +
		body + newline + SpectrEndMarker + after

	if err := os.WriteFile(filePath, []byte(newContent), filePerm); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// expandPath expands ~ to the user's home directory.
// If path starts with ~/, replace ~ with the result of os.UserHomeDir().
// If UserHomeDir returns an error, return the original path unchanged.
// For paths not starting with ~/, return unchanged.
func expandPath(path string) string {
	if !strings.HasPrefix(path, "~/") {
		return path
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	return filepath.Join(homeDir, path[2:])
}

// isGlobalPath returns true if the path is a global path
// (starts with ~/ or /).
func isGlobalPath(path string) bool {
	return strings.HasPrefix(path, "~/") || strings.HasPrefix(path, "/")
}
