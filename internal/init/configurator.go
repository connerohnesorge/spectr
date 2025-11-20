//nolint:revive // line-length-limit - readability over strict formatting
package init

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Configurator interface for all tool configurators
type Configurator interface {
	// Configure configures a tool for the given project path
	Configure(projectPath, spectrDir string) error
	// IsConfigured checks if a tool is already configured for the given project path
	IsConfigured(projectPath string) bool
	// GetName returns the name of the tool
	GetName() string
}

// ============================================================================
// GenericConfigurator - Data-driven configurator for all tools
// ============================================================================

// GenericConfigurator is a data-driven configurator that works with both
// config-based tools (creates CLAUDE.md, CLINE.md, etc.) and slash command
// tools (creates files in .claude/commands/, .cline/commands/, etc.)
type GenericConfigurator struct {
	config ToolConfig
	tm     *TemplateManager
}

// NewGenericConfigurator creates a new generic configurator with the given tool config
func NewGenericConfigurator(config ToolConfig) (*GenericConfigurator, error) {
	tm, err := NewTemplateManager()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize template manager: %w", err)
	}

	return &GenericConfigurator{config: config, tm: tm}, nil
}

// Config returns the tool configuration
func (g *GenericConfigurator) Config() ToolConfig {
	return g.config
}

// Configure configures the tool by creating/updating the appropriate files
func (g *GenericConfigurator) Configure(projectPath, _spectrDir string) error {
	switch g.config.Type {
	case ToolTypeConfig:
		return g.configureConfigTool(projectPath)
	case ToolTypeSlash:
		return g.configureSlashTool(projectPath)
	default:
		return fmt.Errorf("unknown tool type: %s", g.config.Type)
	}
}

// configureConfigTool configures a config-based tool (creates instruction file)
func (g *GenericConfigurator) configureConfigTool(projectPath string) error {
	content, err := g.tm.RenderAgents()
	if err != nil {
		return fmt.Errorf("failed to render agents template: %w", err)
	}

	filePath := filepath.Join(projectPath, g.config.ConfigFile)

	return UpdateFileWithMarkers(filePath, content, SpectrStartMarker, SpectrEndMarker)
}

// configureSlashTool configures a slash command tool (creates 3 command files)
func (g *GenericConfigurator) configureSlashTool(projectPath string) error {
	commands := []string{"proposal", "apply", "archive"}
	for _, cmd := range commands {
		if err := g.configureSlashCommand(projectPath, cmd); err != nil {
			return err
		}
	}

	return nil
}

// configureSlashCommand configures a single slash command
func (g *GenericConfigurator) configureSlashCommand(projectPath, cmd string) error {
	relPath, ok := g.config.SlashPaths[cmd]
	if !ok {
		return fmt.Errorf("missing file path for command: %s", cmd)
	}

	filePath := filepath.Join(projectPath, relPath)

	body, err := g.tm.RenderSlashCommand(cmd)
	if err != nil {
		return fmt.Errorf("failed to render slash command %s: %w", cmd, err)
	}

	if FileExists(filePath) {
		return g.updateExistingSlashCommand(filePath, body)
	}

	return g.createNewSlashCommand(filePath, cmd, body)
}

// updateExistingSlashCommand updates an existing slash command file
func (g *GenericConfigurator) updateExistingSlashCommand(filePath, body string) error {
	if err := updateSlashCommandBody(filePath, body); err != nil {
		return fmt.Errorf("failed to update slash command file %s: %w", filePath, err)
	}

	return nil
}

// createNewSlashCommand creates a new slash command file
func (g *GenericConfigurator) createNewSlashCommand(filePath, cmd, body string) error {
	var sections []string

	if frontmatter, ok := g.config.Frontmatter[cmd]; ok && frontmatter != "" {
		sections = append(sections, strings.TrimSpace(frontmatter))
	}

	sections = append(sections, SpectrStartMarker+newlineDouble+body+newlineDouble+SpectrEndMarker)
	content := strings.Join(sections, newlineDouble) + newlineDouble

	dir := filepath.Dir(filePath)
	if err := EnsureDir(dir); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", filePath, err)
	}

	if err := os.WriteFile(filePath, []byte(content), filePerm); err != nil {
		return fmt.Errorf("failed to write slash command file %s: %w", filePath, err)
	}

	return nil
}

// IsConfigured checks if the tool is already configured
func (g *GenericConfigurator) IsConfigured(projectPath string) bool {
	switch g.config.Type {
	case ToolTypeConfig:
		filePath := filepath.Join(projectPath, g.config.ConfigFile)

		return FileExists(filePath)
	case ToolTypeSlash:
		// Check if all three slash command files exist
		commands := []string{"proposal", "apply", "archive"}
		for _, cmd := range commands {
			relPath, ok := g.config.SlashPaths[cmd]
			if !ok {
				return false
			}
			filePath := filepath.Join(projectPath, relPath)
			if !FileExists(filePath) {
				return false
			}
		}

		return true
	default:
		return false
	}
}

// GetName returns the name of the tool
func (g *GenericConfigurator) GetName() string {
	return g.config.Name
}

// GetFilePaths returns the file paths that this configurator would create/update
func (g *GenericConfigurator) GetFilePaths() []string {
	switch g.config.Type {
	case ToolTypeConfig:
		return []string{g.config.ConfigFile}
	case ToolTypeSlash:
		paths := make([]string, 0, len(g.config.SlashPaths))
		for _, path := range g.config.SlashPaths {
			paths = append(paths, path)
		}

		return paths
	default:
		return []string{}
	}
}

// GetSlashPath returns the slash command path for the given command name
// Returns the path and a boolean indicating if the command exists
func (g *GenericConfigurator) GetSlashPath(cmd string) (string, bool) {
	path, exists := g.config.SlashPaths[cmd]
	return path, exists
}

// ============================================================================
// SlashCommandConfigurator - Legacy wrapper for backward compatibility
// ============================================================================

// SlashCommandConfig holds configuration for a slash command tool
// Deprecated: Use ToolConfig and GenericConfigurator instead
type SlashCommandConfig struct {
	ToolID      string
	ToolName    string
	Frontmatter map[string]string // proposal, apply, archive frontmatter
	FilePaths   map[string]string // proposal, apply, archive paths
}

// SlashCommandConfigurator configures slash commands for a tool
// Deprecated: Use GenericConfigurator instead
type SlashCommandConfigurator struct {
	config SlashCommandConfig
}

// NewSlashCommandConfigurator creates a new slash command configurator
// Deprecated: Use NewGenericConfigurator instead
func NewSlashCommandConfigurator(config SlashCommandConfig) *SlashCommandConfigurator {
	return &SlashCommandConfigurator{config: config}
}

func (s *SlashCommandConfigurator) Configure(projectPath, _spectrDir string) error {
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

// configureCommand configures a single slash command
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
		return fmt.Errorf("failed to render slash command %s: %w", cmd, err)
	}

	if FileExists(filePath) {
		return s.updateExistingCommand(filePath, body)
	}

	return s.createNewCommand(filePath, cmd, body)
}

// updateExistingCommand updates an existing slash command file
func (s *SlashCommandConfigurator) updateExistingCommand(filePath, body string) error {
	if err := updateSlashCommandBody(filePath, body); err != nil {
		return fmt.Errorf("failed to update slash command file %s: %w", filePath, err)
	}

	return nil
}

// createNewCommand creates a new slash command file
func (s *SlashCommandConfigurator) createNewCommand(filePath, cmd, body string) error {
	var sections []string

	if frontmatter, ok := s.config.Frontmatter[cmd]; ok && frontmatter != "" {
		sections = append(sections, strings.TrimSpace(frontmatter))
	}

	sections = append(sections, SpectrStartMarker+newlineDouble+body+newlineDouble+SpectrEndMarker)
	content := strings.Join(sections, newlineDouble) + newlineDouble

	dir := filepath.Dir(filePath)
	if err := EnsureDir(dir); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", filePath, err)
	}

	if err := os.WriteFile(filePath, []byte(content), filePerm); err != nil {
		return fmt.Errorf("failed to write slash command file %s: %w", filePath, err)
	}

	return nil
}

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

func (s *SlashCommandConfigurator) GetName() string {
	return s.config.ToolName
}

// updateSlashCommandBody updates the body of a slash command file between markers
func updateSlashCommandBody(filePath, body string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	contentStr := string(content)

	startIndex := findMarkerIndex(contentStr, SpectrStartMarker, 0)
	if startIndex == -1 {
		return fmt.Errorf("start marker not found in %s", filePath)
	}

	endIndex := findMarkerIndex(contentStr, SpectrEndMarker, startIndex+len(SpectrStartMarker))
	if endIndex == -1 {
		return fmt.Errorf("end marker not found in %s", filePath)
	}

	if endIndex < startIndex {
		return fmt.Errorf("end marker appears before start marker in %s", filePath)
	}

	before := contentStr[:startIndex]
	after := contentStr[endIndex+len(SpectrEndMarker):]
	newContent := before + SpectrStartMarker + "\n" + body + "\n" + SpectrEndMarker + after

	if err := os.WriteFile(filePath, []byte(newContent), filePerm); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
