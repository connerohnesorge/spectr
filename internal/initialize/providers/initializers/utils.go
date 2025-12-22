package initializers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/connerohnesorge/spectr/internal/initialize/types"
	"github.com/spf13/afero"
)

const (
	markerBlockSuffix = "\n\n"
	newlineDouble     = "\n\n"
	newline           = "\n"
)

// UpdateFileWithMarkers updates a file with content between markers.
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
	if exists {
		data, err := afero.ReadFile(fs, filePath)
		if err != nil {
			return fmt.Errorf(
				"failed to read existing file: %w",
				err,
			)
		}
		existingContent = string(data)
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
		switch {
		case startIndex != -1 && endIndex != -1:
			if endIndex < startIndex {
				return fmt.Errorf(
					"invalid marker state in %s: end before start",
					filePath,
				)
			}
			before := existingContent[:startIndex]
			after := existingContent[endIndex+len(endMarker):]
			existingContent = before + startMarker + "\n" +
				content + "\n" + endMarker + after
		case startIndex == -1 && endIndex == -1:
			existingContent = startMarker + "\n" + content + "\n" +
				endMarker + markerBlockSuffix + existingContent
		default:
			return fmt.Errorf(
				"invalid marker state in %s",
				filePath,
			)
		}
	} else {
		existingContent = startMarker + "\n" + content + "\n" +
			endMarker + markerBlockSuffix
	}
	dir := filepath.Dir(filePath)
	if err := fs.MkdirAll(dir, dirPerm); err != nil {
		return fmt.Errorf(
			"failed to create parent directory: %w",
			err,
		)
	}
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
	newContent := before + types.SpectrStartMarker + newline + 
		body + newline + types.SpectrEndMarker + after
	return afero.WriteFile(
		fs,
		filePath,
		[]byte(newContent),
		filePerm,
	)
}

func createNewSlashCommand(
	fs afero.Fs,
	filePath, body, frontmatter string,
) error {
	var sections []string
	if frontmatter != "" {
		sections = append(
			sections,
			strings.TrimSpace(frontmatter),
		)
	}
	sections = append(
		sections,
		types.SpectrStartMarker+newlineDouble+body+newlineDouble+
			types.SpectrEndMarker,
	)
	content := strings.Join(
		sections,
		newlineDouble,
	) + newlineDouble
	dir := filepath.Dir(filePath)
	if err := fs.MkdirAll(dir, dirPerm); err != nil {
		return err
	}
	return afero.WriteFile(
		fs,
		filePath,
		[]byte(content),
		filePerm,
	)
}

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

func isMarkerOnOwnLine(
	content string,
	markerIndex, markerLength int,
) bool {
	leftIndex := markerIndex - 1
	for leftIndex >= 0 && content[leftIndex] != '\n' {
		char := content[leftIndex]
		if char != ' ' && char != '\t' &&
			char != '\r' {
			return false
		}
		leftIndex--
	}
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