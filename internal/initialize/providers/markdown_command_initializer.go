package providers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// MarkdownSlashCommandInitializer handles markdown slash command files.
// These files use YAML frontmatter for metadata and contain the command
// prompt between spectr markers.
type MarkdownSlashCommandInitializer struct {
	path        string // Relative path to command file
	commandName string // Command name for rendering (e.g., "proposal")
	frontmatter string // YAML frontmatter content
}

// NewMarkdownSlashCommandInitializer creates a new markdown slash command
// initializer. path is the relative path to the command file. commandName is
// used to render the appropriate template (e.g., "proposal", "apply").
// frontmatter is the YAML frontmatter to include at the top of new files.
func NewMarkdownSlashCommandInitializer(
	path, commandName, frontmatter string,
) *MarkdownSlashCommandInitializer {
	return &MarkdownSlashCommandInitializer{
		path:        path,
		commandName: commandName,
		frontmatter: frontmatter,
	}
}

// ID returns the unique identifier for this initializer.
// Format: "markdown-cmd:{path}"
func (i *MarkdownSlashCommandInitializer) ID() string {
	return "markdown-cmd:" + i.path
}

// FilePath returns the relative path this initializer manages.
func (i *MarkdownSlashCommandInitializer) FilePath() string {
	return i.path
}

// Configure creates or updates the markdown slash command file.
func (i *MarkdownSlashCommandInitializer) Configure(
	projectPath string,
	tm TemplateRenderer,
) error {
	body, err := tm.RenderSlashCommand(
		i.commandName,
		DefaultTemplateContext(),
	)
	if err != nil {
		return fmt.Errorf(
			"failed to render slash command %s: %w",
			i.commandName,
			err,
		)
	}

	fullPath := i.expandedPath(projectPath)

	if FileExists(fullPath) {
		return i.updateExisting(fullPath, body)
	}

	return i.createNew(fullPath, body)
}

// IsConfigured checks if the slash command file exists.
func (i *MarkdownSlashCommandInitializer) IsConfigured(
	projectPath string,
) bool {
	fullPath := i.expandedPath(projectPath)

	return FileExists(fullPath)
}

// expandedPath returns the full path for the command file,
// handling ~ expansion for global paths.
func (i *MarkdownSlashCommandInitializer) expandedPath(
	projectPath string,
) string {
	if isGlobalPath(i.path) {
		return expandPath(i.path)
	}

	return filepath.Join(projectPath, i.path)
}

// updateExisting updates an existing slash command file.
func (i *MarkdownSlashCommandInitializer) updateExisting(
	filePath, body string,
) error {
	err := updateSlashCommandBody(
		filePath,
		body,
		i.frontmatter,
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

// createNew creates a new slash command file.
func (i *MarkdownSlashCommandInitializer) createNew(
	filePath, body string,
) error {
	var sections []string

	if i.frontmatter != "" {
		sections = append(
			sections,
			strings.TrimSpace(i.frontmatter),
		)
	}

	sections = append(
		sections,
		SpectrStartMarker+newlineDouble+body+newlineDouble+SpectrEndMarker,
	)
	content := strings.Join(
		sections,
		newlineDouble,
	) + newlineDouble

	dir := filepath.Dir(filePath)
	err := EnsureDir(dir)
	if err != nil {
		return fmt.Errorf(
			"failed to create directory for %s: %w",
			filePath,
			err,
		)
	}

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
