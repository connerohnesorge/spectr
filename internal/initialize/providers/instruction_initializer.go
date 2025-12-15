package providers

import (
	"fmt"
	"path/filepath"
)

// InstructionFileInitializer handles instruction files like CLAUDE.md,
// CURSOR.md, etc. These files contain pointers to the spectr/AGENTS.md file
// and are updated using start/end markers to preserve user content.
type InstructionFileInitializer struct {
	path string // Relative path to instruction file (e.g., "CLAUDE.md")
}

// NewInstructionFileInitializer creates a new instruction file initializer.
// path is the relative path to the instruction file (e.g., "CLAUDE.md").
// The path may include ~ for home directory paths.
func NewInstructionFileInitializer(path string) *InstructionFileInitializer {
	return &InstructionFileInitializer{
		path: path,
	}
}

// ID returns the unique identifier for this initializer.
// Format: "instruction:{path}"
func (i *InstructionFileInitializer) ID() string {
	return "instruction:" + i.path
}

// FilePath returns the relative path this initializer manages.
func (i *InstructionFileInitializer) FilePath() string {
	return i.path
}

// Configure creates or updates the instruction file.
// It renders the instruction pointer template and writes it between markers.
func (i *InstructionFileInitializer) Configure(
	projectPath string,
	tm TemplateRenderer,
) error {
	content, err := tm.RenderInstructionPointer(
		DefaultTemplateContext(),
	)
	if err != nil {
		return fmt.Errorf(
			"failed to render instruction pointer template: %w",
			err,
		)
	}

	fullPath := i.expandedPath(projectPath)

	return UpdateFileWithMarkers(
		fullPath,
		content,
		SpectrStartMarker,
		SpectrEndMarker,
	)
}

// IsConfigured checks if the instruction file exists.
func (i *InstructionFileInitializer) IsConfigured(projectPath string) bool {
	fullPath := i.expandedPath(projectPath)

	return FileExists(fullPath)
}

// expandedPath returns the full path for the instruction file,
// handling ~ expansion for global paths.
func (i *InstructionFileInitializer) expandedPath(projectPath string) string {
	if isGlobalPath(i.path) {
		return expandPath(i.path)
	}

	return filepath.Join(projectPath, i.path)
}
