// Package initializers provides initialization logic for various providers.
// This file contains utility functions for file and path operations.
package initializers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/connerohnesorge/spectr/internal/initialize/types"
	"github.com/spf13/afero"
)

// Permission constants for directory and file creation.
const (
	// dirPerm is the default permission for new directories (0755).
	dirPerm = 0755
	// filePerm is the default permission for new files (0644).
	filePerm = 0644
)

// String constants for newlines and block suffixes.
const (
	// markerBlockSuffix is the suffix added after a marker block.
	markerBlockSuffix = "\n\n"
	// newlineDouble is a double newline string.
	newlineDouble = "\n\n"
	// newline is a single newline string.
	newline = "\n"
)

// UpdateFileWithMarkers updates a file with content between markers.
// It handles different scenarios:
// 1. If the file exists and has markers, it replaces the content between them.
// 2. If the file exists but has no markers, it prepends the markers and content.
// 3. If the file does not exist, it creates it with markers and content.
//
// Parameters:
//   - fs: The afero.Fs to perform operations on.
//   - filePath: The path to the file to update.
//   - content: The content to place between markers.
//   - startMarker: The string that marks the beginning of the block.
//   - endMarker: The string that marks the end of the block.
//
//nolint:revive // argument-limit - interface defined elsewhere
func UpdateFileWithMarkers(
	fs afero.Fs,
	filePath, content, startMarker, endMarker string,
) error {
	var existingContent string
	exists, err := afero.Exists(fs, filePath)
	if err != nil {
		return fmt.Errorf(
			"failed to check file existence: %w",
			err,
		)
	}

	// If the file exists, read its content and find marker positions.
	if exists {
		data, err := afero.ReadFile(fs, filePath)
		if err != nil {
			return fmt.Errorf(
				"failed to read existing file: %w",
				err,
			)
		}
		existingContent = string(data)

		// Find where the start marker begins.
		startIndex := findMarkerIndex(
			existingContent,
			startMarker,
			0,
		)

		var endIndex int
		if startIndex != -1 {
			// Start marker exists, look for end marker after it.
			startPos := startIndex + len(
				startMarker,
			)
			endIndex = findMarkerIndex(
				existingContent,
				endMarker,
				startPos,
			)
		} else {
			// Start marker not found, check if end marker exists anywhere.
			endIndex = findMarkerIndex(existingContent, endMarker, 0)
		}

		switch {
		case startIndex != -1 && endIndex != -1:
			// Both markers found, replace what's between them.
			if endIndex < startIndex {
				return fmt.Errorf(
					"invalid marker state in %s: end before start",
					filePath,
				)
			}
			before := existingContent[:startIndex]
			after := existingContent[endIndex+len(endMarker):]
			existingContent = before + startMarker + newline +
				content + newline + endMarker + after
		case startIndex == -1 && endIndex == -1:
			// Neither marker found, prepend them to the file.
			existingContent = startMarker + newline + content + newline +
				endMarker + markerBlockSuffix + existingContent
		default:
			// Only one marker found, which is an invalid state for the file.
			return fmt.Errorf(
				"invalid marker state in %s",
				filePath,
			)
		}
	} else {
		// File does not exist, initialize it with markers and content.
		existingContent = startMarker + newline + content + newline +
			endMarker + markerBlockSuffix
	}

	// Ensure the parent directory exists before writing.
	dir := filepath.Dir(filePath)
	if err := fs.MkdirAll(dir, dirPerm); err != nil {
		return fmt.Errorf(
			"failed to create parent directory: %w",
			err,
		)
	}

	// Write the updated content to the file.
	if err := afero.WriteFile(
		fs,
		filePath,
		[]byte(existingContent),
		filePerm,
	); err != nil {
		return fmt.Errorf(
			"failed to write file: %w",
			err,
		)
	}

	return nil
}

// IsGlobalPath checks if the path is global (starts with ~/ or /).
func IsGlobalPath(path string) bool {
	return strings.HasPrefix(path, "~/") ||
		strings.HasPrefix(path, "/")
}

// ExpandPath expands the home directory in the path.
// If the path starts with ~/, it replaces it with the user's home directory.
func ExpandPath(path string) string {
	if !strings.HasPrefix(path, "~/") {
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	return filepath.Join(home, path[2:])
}

// updateSlashCommandBody updates the body of a slash command file.
// It searches for Spectr markers and replaces the content between them.
// It also handles frontmatter preservation or addition.
func updateSlashCommandBody(
	fs afero.Fs,
	filePath, body, frontmatter string,
) error {
	content, err := afero.ReadFile(fs, filePath)
	if err != nil {
		return fmt.Errorf(
			"failed to read file: %w",
			err,
		)
	}

	contentStr := string(content)

	// Find Spectr start marker in the file.
	startIndex := findMarkerIndex(
		contentStr,
		types.SpectrStartMarker,
		0,
	)
	if startIndex == -1 {
		return fmt.Errorf(
			"start marker not found in %s",
			filePath,
		)
	}

	// Find Spectr end marker in the file.
	searchOffset := startIndex + len(
		types.SpectrStartMarker,
	)
	endIndex := findMarkerIndex(
		contentStr,
		types.SpectrEndMarker,
		searchOffset,
	)
	if endIndex == -1 {
		return fmt.Errorf(
			"end marker not found in %s",
			filePath,
		)
	}

	if endIndex < startIndex {
		return fmt.Errorf(
			"end marker appears before start marker in %s",
			filePath,
		)
	}

	before := contentStr[:startIndex]
	after := contentStr[endIndex+len(types.SpectrEndMarker):]

	// Handle frontmatter if it's required but missing.
	if frontmatter != "" &&
		!strings.HasPrefix(
			strings.TrimSpace(before),
			"---",
		) {
		before = strings.TrimSpace(
			frontmatter,
		) + newlineDouble + strings.TrimLeft(
			before,
			"\n\r",
		)
	}

	// Construct and write the new file content.
	newContent := before + types.SpectrStartMarker + newline +
		body + newline + types.SpectrEndMarker + after

	return afero.WriteFile(
		fs,
		filePath,
		[]byte(newContent),
		filePerm,
	)
}

// createNewSlashCommand creates a new slash command file.
// It includes optional frontmatter and the Spectr markers.
func createNewSlashCommand(
	fs afero.Fs,
	filePath, body, frontmatter string,
) error {
	var sections []string

	// Add frontmatter at the top if provided.
	if frontmatter != "" {
		sections = append(
			sections,
			strings.TrimSpace(frontmatter),
		)
	}

	// Add the body wrapped in Spectr markers.
	sections = append(
		sections,
		types.SpectrStartMarker+newlineDouble+body+newlineDouble+
			types.SpectrEndMarker,
	)

	// Join sections with double newlines.
	content := strings.Join(
		sections,
		newlineDouble,
	) + newlineDouble

	// Create directory structure if needed.
	dir := filepath.Dir(filePath)
	if err := fs.MkdirAll(dir, dirPerm); err != nil {
		return err
	}

	// Write the initial command file.
	return afero.WriteFile(
		fs,
		filePath,
		[]byte(content),
		filePerm,
	)
}

// findMarkerIndex finds the index of a marker in the content.
// It ensures that the marker is on its own line by checking
// surrounding whitespace.
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

	// Iterate to find a marker that satisfies the "own line" condition.
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
// It looks for only whitespace characters (space, tab, carriage return)
// between the marker and the nearest newlines on both sides.
func isMarkerOnOwnLine(
	content string,
	markerIndex, markerLength int,
) bool {
	// Check the left side of the marker.
	leftIndex := markerIndex - 1
	for leftIndex >= 0 && content[leftIndex] != '\n' {
		char := content[leftIndex]
		if char != ' ' && char != '\t' &&
			char != '\r' {
			return false
		}
		leftIndex--
	}

	// Check the right side of the marker.
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