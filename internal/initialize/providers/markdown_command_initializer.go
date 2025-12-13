// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
package providers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// MarkdownSlashCommandInitializer manages markdown-based slash command files.
//
// These files are used by AI tools like Claude, Cline, Cursor, Aider,
// etc. They contain command instructions with optional YAML frontmatter
// for metadata.
//
// The initializer uses spectr markers (<!-- spectr:START --> and
// <!-- spectr:END -->) to manage a section of the file, allowing users to add
// their own content before or after the managed section.
type MarkdownSlashCommandInitializer struct {
	// path is the file path (e.g., ".claude/commands/spectr/proposal.md")
	path string
	// commandName is the command name (e.g., "proposal", "apply")
	commandName string
	// frontmatter is the YAML frontmatter for new files
	frontmatter string
}

// NewMarkdownSlashCommandInitializer creates a new initializer.
//
// Parameters:
//   - path: The file path (e.g., ".claude/commands/spectr/proposal.md")
//   - commandName: The command name (e.g., "proposal", "apply")
//   - frontmatter: The YAML frontmatter for new files
//     (e.g., StandardProposalFrontmatter)
//
// Example usage:
//
//	NewMarkdownSlashCommandInitializer(
//	    ".claude/commands/spectr/proposal.md",
//	    "proposal",
//	    StandardProposalFrontmatter,
//	)
func NewMarkdownSlashCommandInitializer(
	path, commandName, frontmatter string,
) *MarkdownSlashCommandInitializer {
	return &MarkdownSlashCommandInitializer{
		path:        path,
		commandName: commandName,
		frontmatter: frontmatter,
	}
}

// ID returns a unique identifier for this initializer.
//
// Format: "markdown-cmd:{path}"
// e.g., "markdown-cmd:.claude/commands/spectr/proposal.md"
func (m *MarkdownSlashCommandInitializer) ID() string {
	return "markdown-cmd:" + m.path
}

// FilePath returns the path this initializer manages.
//
// The returned path may contain ~ for home directory paths.
// Path expansion is handled internally during Configure and IsConfigured.
func (m *MarkdownSlashCommandInitializer) FilePath() string {
	return m.path
}

// Configure creates or updates the markdown slash command file.
//
// For project-relative paths, the file is created at
// filepath.Join(projectPath, path). For global paths (starting with ~/ or /),
// the path is expanded independently.
//
// The file content is rendered using TemplateRenderer.RenderSlashCommand,
// which generates the command body content.
//
// Content is managed using spectr markers, allowing users to add custom content
// before or after the managed section. If the file exists, only the content
// between markers is updated. If the file doesn't exist, it's created with
// frontmatter (if provided) and markers.
func (m *MarkdownSlashCommandInitializer) Configure(
	projectPath string,
	tm TemplateRenderer,
) error {
	// Render the slash command content
	body, err := tm.RenderSlashCommand(
		m.commandName,
		DefaultTemplateContext(),
	)
	if err != nil {
		return fmt.Errorf(
			"failed to render slash command %s: %w",
			m.commandName,
			err,
		)
	}

	// Determine the full file path
	fullPath := m.resolvePath(projectPath)

	// Update existing file or create new one
	if FileExists(fullPath) {
		return m.updateExisting(fullPath, body)
	}

	return m.createNew(fullPath, body)
}

// IsConfigured checks if the markdown slash command file exists.
//
// Path resolution follows the same rules as Configure:
//   - Project-relative paths are joined with projectPath
//   - Global paths are expanded independently
func (m *MarkdownSlashCommandInitializer) IsConfigured(
	projectPath string,
) bool {
	fullPath := m.resolvePath(projectPath)

	return FileExists(fullPath)
}

// resolvePath returns the full path for the slash command file.
//
// For global paths (starting with ~/ or /), the path is expanded.
// For project-relative paths, the path is joined with projectPath.
func (m *MarkdownSlashCommandInitializer) resolvePath(
	projectPath string,
) string {
	if isGlobalPath(m.path) {
		return expandPath(m.path)
	}

	return filepath.Join(projectPath, m.path)
}

// updateExisting updates an existing slash command file.
//
// Uses the updateSlashCommandBody helper to update content between markers.
// If frontmatter is provided and the file doesn't have frontmatter, it's added.
func (m *MarkdownSlashCommandInitializer) updateExisting(
	filePath, body string,
) error {
	err := updateSlashCommandBody(
		filePath,
		body,
		m.frontmatter,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to update slash command file %s: %w",
			filePath,
			err,
		)
	}

	return nil
}

// createNew creates a new slash command file with frontmatter and markers.
//
// The file structure is:
//
//	{frontmatter (if provided)}
//
//	<!-- spectr:START -->
//
//	{body content}
//
//	<!-- spectr:END -->
func (m *MarkdownSlashCommandInitializer) createNew(
	filePath, body string,
) error {
	var sections []string

	// Add frontmatter if provided
	if m.frontmatter != "" {
		sections = append(
			sections,
			strings.TrimSpace(m.frontmatter),
		)
	}

	// Add markers with body content
	sections = append(
		sections,
		SpectrStartMarker+newlineDouble+body+newlineDouble+SpectrEndMarker,
	)

	// Join sections and ensure trailing newline
	content := strings.Join(
		sections,
		newlineDouble,
	) + newlineDouble

	// Ensure directory exists
	dir := filepath.Dir(filePath)
	err := EnsureDir(dir)
	if err != nil {
		return fmt.Errorf(
			"failed to create directory for %s: %w",
			filePath,
			err,
		)
	}

	// Write the file
	err = os.WriteFile(
		filePath,
		[]byte(content),
		filePerm,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to write slash command file %s: %w",
			filePath,
			err,
		)
	}

	return nil
}
