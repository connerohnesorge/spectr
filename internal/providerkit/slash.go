package providerkit

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	newlineDouble = "\n\n"
)

// SlashCommandConfig holds configuration for a slash command tool
type SlashCommandConfig struct {
	// ToolID is the unique identifier for the tool (e.g., "claude", "cursor")
	ToolID string
	// ToolName is the human-readable name (e.g., "Claude Slash Commands")
	ToolName string
	// Frontmatter maps command types to their frontmatter content
	// Keys: "proposal", "apply", "archive"
	// Values: YAML frontmatter or markdown headers
	Frontmatter map[string]string
	// FilePaths maps command types to their relative file paths
	// Keys: "proposal", "apply", "archive"
	// Values: Paths like ".claude/commands/spectr/proposal.md"
	FilePaths map[string]string
}

// SlashCommandConfigurator configures slash commands for a tool.
// It implements the Provider interface.
type SlashCommandConfigurator struct {
	config SlashCommandConfig
}

// NewSlashCommandConfigurator creates a new slash command configurator
func NewSlashCommandConfigurator(
	config SlashCommandConfig,
) *SlashCommandConfigurator {
	return &SlashCommandConfigurator{config: config}
}

// Configure configures all three slash command files (proposal, apply, archive)
// for the tool in the given project path.
func (s *SlashCommandConfigurator) Configure(
	projectPath,
	_spectrDir string,
) error {
	tm, err := NewTemplateManager()
	if err != nil {
		return err
	}

	commands := []string{"proposal", "apply", "archive"}

	for _, cmd := range commands {
		if err := s.configureCommand(tm, projectPath, cmd); err != nil {
			return err
		}
	}

	return nil
}

// configureCommand configures a single slash command file
func (s *SlashCommandConfigurator) configureCommand(
	tm *TemplateManager,
	projectPath, cmd string,
) error {
	relPath, ok := s.config.FilePaths[cmd]
	if !ok {
		return fmt.Errorf("missing file path for command: %s", cmd)
	}

	filePath := filepath.Join(projectPath, relPath)

	body, err := tm.RenderSlashCommand(cmd)
	if err != nil {
		return fmt.Errorf(
			"failed to render slash command %s: %w",
			cmd,
			err,
		)
	}

	if FileExists(filePath) {
		return s.updateExistingCommand(filePath, body)
	}

	return s.createNewCommand(filePath, cmd, body)
}

// updateExistingCommand updates an existing slash command file
func (*SlashCommandConfigurator) updateExistingCommand(
	filePath, body string,
) error {
	if err := updateSlashCommandBody(filePath, body); err != nil {
		return fmt.Errorf(
			"failed to update slash command file %s: %w",
			filePath,
			err,
		)
	}

	return nil
}

// createNewCommand creates a new slash command file
func (s *SlashCommandConfigurator) createNewCommand(
	filePath, cmd, body string,
) error {
	var sections []string

	if frontmatter, ok := s.config.Frontmatter[cmd]; ok && frontmatter != "" {
		sections = append(sections, strings.TrimSpace(frontmatter))
	}

	sections = append(
		sections,
		SpectrStartMarker+newlineDouble+body+newlineDouble+SpectrEndMarker,
	)

	content := strings.Join(sections, newlineDouble) + newlineDouble

	dir := filepath.Dir(filePath)
	if err := EnsureDir(dir); err != nil {
		return fmt.Errorf(
			"failed to create directory for %s: %w",
			filePath,
			err,
		)
	}

	if err := os.WriteFile(
		filePath,
		[]byte(content),
		defaultFilePerm,
	); err != nil {
		return fmt.Errorf(
			"failed to write slash command file %s: %w",
			filePath,
			err,
		)
	}

	return nil
}

// IsConfigured checks if all three slash command files exist for this tool
func (s *SlashCommandConfigurator) IsConfigured(projectPath string) bool {
	// Check if all three slash command files exist
	commands := []string{"proposal", "apply", "archive"}
	for _, cmd := range commands {
		relPath, ok := s.config.FilePaths[cmd]
		if !ok {
			return false
		}

		filePath := filepath.Join(projectPath, relPath)
		if !FileExists(filePath) {
			return false
		}
	}

	return true
}

// GetName returns the human-readable name of the tool
func (s *SlashCommandConfigurator) GetName() string {
	return s.config.ToolName
}

// GetConfig returns the underlying configuration
// (useful for extracting file paths)
func (s *SlashCommandConfigurator) GetConfig() SlashCommandConfig {
	return s.config
}

// updateSlashCommandBody updates the body of a slash command file
// between markers
func updateSlashCommandBody(filePath, body string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	contentStr := string(content)

	startIndex := strings.Index(contentStr, SpectrStartMarker)
	endIndex := strings.Index(contentStr, SpectrEndMarker)

	if startIndex == -1 || endIndex == -1 || endIndex <= startIndex {
		return fmt.Errorf("missing Spectr markers in %s", filePath)
	}

	before := contentStr[:startIndex+len(SpectrStartMarker)]
	after := contentStr[endIndex:]
	updatedContent := before + "\n" + body + "\n" + after

	if err := os.WriteFile(
		filePath,
		[]byte(updatedContent),
		defaultFilePerm,
	); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
