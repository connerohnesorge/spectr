// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
package providers

import (
	"fmt"
	"path/filepath"
)

// InstructionFileInitializer manages instruction files like CLAUDE.md, etc.
//
// These files contain instructions for AI assistants and are typically placed
// at the root of a project or in a global configuration directory.
//
// The initializer uses spectr markers (<!-- spectr:START --> and
// <!-- spectr:END -->) to manage a section of the file, allowing users to add
// their own content before or after the managed section.
type InstructionFileInitializer struct {
	// path is the file path (e.g., "CLAUDE.md", "~/.codex/AGENTS.md")
	path string
}

// NewInstructionFileInitializer creates a new InstructionFileInitializer.
//
// The path parameter specifies the file location:
//   - Project-relative paths (e.g., "CLAUDE.md") are joined with projectPath
//   - Global paths starting with ~/ are expanded to the user's home directory
//
// Example usage:
//
//	NewInstructionFileInitializer("CLAUDE.md")
//	NewInstructionFileInitializer("~/.codex/AGENTS.md")
func NewInstructionFileInitializer(path string) *InstructionFileInitializer {
	return &InstructionFileInitializer{
		path: path,
	}
}

// ID returns a unique identifier for this initializer.
//
// Format: "instruction:{path}" e.g., "instruction:CLAUDE.md"
func (i *InstructionFileInitializer) ID() string {
	return "instruction:" + i.path
}

// FilePath returns the path this initializer manages.
//
// The returned path may contain ~ for home directory paths.
// Path expansion is handled internally during Configure and IsConfigured.
func (i *InstructionFileInitializer) FilePath() string {
	return i.path
}

// Configure creates or updates the instruction file.
//
// For project-relative paths, the file is created at
// filepath.Join(projectPath, path). For global paths (starting with ~/ or /),
// the path is expanded independently.
//
// The file content is rendered using TemplateRenderer.RenderInstructionPointer,
// which generates a short pointer directing AI assistants to the
// spectr/AGENTS.md file.
//
// Content is managed using spectr markers, allowing users to add custom content
// before or after the managed section.
func (i *InstructionFileInitializer) Configure(
	projectPath string,
	tm TemplateRenderer,
) error {
	// Render the instruction pointer content
	content, err := tm.RenderInstructionPointer(
		DefaultTemplateContext(),
	)
	if err != nil {
		return fmt.Errorf(
			"failed to render instruction pointer template: %w",
			err,
		)
	}

	// Determine the full file path
	fullPath := i.resolvePath(projectPath)

	// Create or update the file with markers
	return UpdateFileWithMarkers(
		fullPath,
		content,
		SpectrStartMarker,
		SpectrEndMarker,
	)
}

// IsConfigured checks if the instruction file exists.
//
// Path resolution follows the same rules as Configure:
//   - Project-relative paths are joined with projectPath
//   - Global paths are expanded independently
func (i *InstructionFileInitializer) IsConfigured(
	projectPath string,
) bool {
	fullPath := i.resolvePath(projectPath)

	return FileExists(fullPath)
}

// resolvePath returns the full path for the instruction file.
//
// For global paths (starting with ~/ or /), the path is expanded.
// For project-relative paths, the path is joined with projectPath.
func (i *InstructionFileInitializer) resolvePath(
	projectPath string,
) string {
	if isGlobalPath(i.path) {
		return expandPath(i.path)
	}

	return filepath.Join(projectPath, i.path)
}
