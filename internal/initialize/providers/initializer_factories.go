// Package providers defines the core interfaces for the provider architecture.
//
// This file provides factory functions for creating initializers.
// Provider implementations should use these factories to create initializers
// without needing to define their own initializer structs.
package providers

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

// directoryInitializer creates one or more directories.
type directoryInitializer struct {
	paths  []string
	global bool
}

// NewDirectoryInitializer creates a new DirectoryInitializer for project directories.
func NewDirectoryInitializer(paths ...string) Initializer {
	return &directoryInitializer{
		paths:  paths,
		global: false,
	}
}

// NewGlobalDirectoryInitializer creates a new DirectoryInitializer for global directories.
func NewGlobalDirectoryInitializer(paths ...string) Initializer {
	return &directoryInitializer{
		paths:  paths,
		global: true,
	}
}

func (d *directoryInitializer) Init(
	ctx context.Context,
	fs afero.Fs,
	cfg *Config,
	tm TemplateManager,
) error {
	for _, p := range d.paths {
		if err := fs.MkdirAll(p, 0755); err != nil {
			return err
		}
	}
	return nil
}

func (d *directoryInitializer) IsSetup(fs afero.Fs, cfg *Config) bool {
	for _, p := range d.paths {
		exists, err := afero.DirExists(fs, p)
		if err != nil || !exists {
			return false
		}
	}
	return true
}

func (d *directoryInitializer) Path() string {
	if len(d.paths) == 0 {
		return ""
	}
	return d.paths[0]
}

func (d *directoryInitializer) IsGlobal() bool {
	return d.global
}

// configFileInitializer creates or updates instruction files with marker-based content.
type configFileInitializer struct {
	path   string
	global bool
}

// NewConfigFileInitializer creates a new ConfigFileInitializer for a project config file.
func NewConfigFileInitializer(path string) Initializer {
	return &configFileInitializer{
		path:   path,
		global: false,
	}
}

// NewGlobalConfigFileInitializer creates a new ConfigFileInitializer for a global config file.
func NewGlobalConfigFileInitializer(path string) Initializer {
	return &configFileInitializer{
		path:   path,
		global: true,
	}
}

func (c *configFileInitializer) Init(
	ctx context.Context,
	fs afero.Fs,
	cfg *Config,
	tm TemplateManager,
) error {
	// Render the instruction pointer template
	templateCtx := NewTemplateContext(cfg)
	content, err := tm.RenderInstructionPointer(templateCtx)
	if err != nil {
		return err
	}

	// Ensure parent directory exists
	dir := filepath.Dir(c.path)
	if dir != "" && dir != "." {
		if err := fs.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// Check if file exists
	exists, err := afero.Exists(fs, c.path)
	if err != nil {
		return err
	}

	if !exists {
		// Create new file with markers
		newContent := SpectrStartMarker + newline + content + newline + SpectrEndMarker + newline
		return afero.WriteFile(fs, c.path, []byte(newContent), filePerm)
	}

	// Read existing file
	existingContent, err := afero.ReadFile(fs, c.path)
	if err != nil {
		return err
	}

	// Update file with markers
	return c.updateExistingFile(fs, string(existingContent), content)
}

func (c *configFileInitializer) updateExistingFile(
	fs afero.Fs,
	existingContent, newMarkerContent string,
) error {
	startIndex := findMarkerIndex(existingContent, SpectrStartMarker, 0)
	endIndex := -1
	if startIndex != -1 {
		searchOffset := startIndex + len(SpectrStartMarker)
		endIndex = findMarkerIndex(existingContent, SpectrEndMarker, searchOffset)
	}

	var finalContent string
	if startIndex == -1 || endIndex == -1 {
		// No markers found, append to end
		finalContent = existingContent + newlineDouble +
			SpectrStartMarker + newline + newMarkerContent + newline + SpectrEndMarker + newline
	} else {
		// Replace content between markers
		before := existingContent[:startIndex]
		after := existingContent[endIndex+len(SpectrEndMarker):]
		finalContent = before + SpectrStartMarker + newline +
			newMarkerContent + newline + SpectrEndMarker + after
	}

	return afero.WriteFile(fs, c.path, []byte(finalContent), filePerm)
}

func (c *configFileInitializer) IsSetup(fs afero.Fs, cfg *Config) bool {
	exists, err := afero.Exists(fs, c.path)
	if err != nil || !exists {
		return false
	}

	content, err := afero.ReadFile(fs, c.path)
	if err != nil {
		return false
	}

	contentStr := string(content)
	return strings.Contains(contentStr, SpectrStartMarker) &&
		strings.Contains(contentStr, SpectrEndMarker)
}

func (c *configFileInitializer) Path() string {
	return c.path
}

func (c *configFileInitializer) IsGlobal() bool {
	return c.global
}

// slashCommandsInitializer creates slash command files from templates.
type slashCommandsInitializer struct {
	dir         string
	ext         string
	format      CommandFormat
	frontmatter map[string]string
	global      bool
}

// NewSlashCommandsInitializer creates a new SlashCommandsInitializer for project commands.
func NewSlashCommandsInitializer(dir, ext string, format CommandFormat) Initializer {
	return &slashCommandsInitializer{
		dir:         dir,
		ext:         ext,
		format:      format,
		frontmatter: nil,
		global:      false,
	}
}

// NewSlashCommandsInitializerWithFrontmatter creates a SlashCommandsInitializer with
// custom YAML frontmatter for Markdown format commands.
func NewSlashCommandsInitializerWithFrontmatter(
	dir, ext string,
	format CommandFormat,
	frontmatter map[string]string,
) Initializer {
	return &slashCommandsInitializer{
		dir:         dir,
		ext:         ext,
		format:      format,
		frontmatter: frontmatter,
		global:      false,
	}
}

// NewGlobalSlashCommandsInitializer creates a SlashCommandsInitializer for global commands.
func NewGlobalSlashCommandsInitializer(dir, ext string, format CommandFormat) Initializer {
	return &slashCommandsInitializer{
		dir:         dir,
		ext:         ext,
		format:      format,
		frontmatter: nil,
		global:      true,
	}
}

// NewGlobalSlashCommandsInitializerWithFrontmatter creates a global SlashCommandsInitializer
// with custom YAML frontmatter for Markdown format commands.
func NewGlobalSlashCommandsInitializerWithFrontmatter(
	dir, ext string,
	format CommandFormat,
	frontmatter map[string]string,
) Initializer {
	return &slashCommandsInitializer{
		dir:         dir,
		ext:         ext,
		format:      format,
		frontmatter: frontmatter,
		global:      true,
	}
}

const (
	slashCommandProposal    = "proposal"
	slashCommandApply       = "apply"
	tomlDescriptionProposal = "Scaffold a new Spectr change and validate strictly."
	tomlDescriptionApply    = "Implement an approved Spectr change and keep tasks in sync."
)

func (s *slashCommandsInitializer) Init(
	ctx context.Context,
	fs afero.Fs,
	cfg *Config,
	tm TemplateManager,
) error {
	// Ensure directory exists
	if err := fs.MkdirAll(s.dir, 0755); err != nil {
		return err
	}

	// Create both proposal and apply commands
	commands := []string{slashCommandProposal, slashCommandApply}
	for _, cmd := range commands {
		if err := s.initCommand(ctx, fs, cfg, tm, cmd); err != nil {
			return err
		}
	}

	return nil
}

func (s *slashCommandsInitializer) initCommand(
	ctx context.Context,
	fs afero.Fs,
	cfg *Config,
	tm TemplateManager,
	cmd string,
) error {
	// Render the command template
	templateCtx := NewTemplateContext(cfg)
	body, err := tm.RenderSlashCommand(cmd, templateCtx)
	if err != nil {
		return err
	}

	// Build file path
	filePath := filepath.Join(s.dir, cmd+s.ext)

	// Create based on format
	switch s.format {
	case FormatMarkdown:
		return s.initMarkdownCommand(fs, filePath, cmd, body)
	case FormatTOML:
		return s.initTOMLCommand(fs, filePath, cmd, body)
	default:
		return fmt.Errorf("unsupported command format: %d", s.format)
	}
}

func (s *slashCommandsInitializer) initMarkdownCommand(
	fs afero.Fs,
	filePath, cmd, body string,
) error {
	// Check if file exists
	exists, err := afero.Exists(fs, filePath)
	if err != nil {
		return err
	}

	if exists {
		return s.updateExistingMarkdownCommand(fs, filePath, cmd, body)
	}

	return s.createNewMarkdownCommand(fs, filePath, cmd, body)
}

func (s *slashCommandsInitializer) createNewMarkdownCommand(
	fs afero.Fs,
	filePath, cmd, body string,
) error {
	var sections []string

	// Add frontmatter if provided
	if frontmatter, ok := s.frontmatter[cmd]; ok && frontmatter != "" {
		sections = append(sections, strings.TrimSpace(frontmatter))
	}

	// Add body wrapped in markers
	markerContent := SpectrStartMarker + newlineDouble + body + newlineDouble + SpectrEndMarker
	sections = append(sections, markerContent)

	content := strings.Join(sections, newlineDouble) + newlineDouble

	return afero.WriteFile(fs, filePath, []byte(content), filePerm)
}

func (s *slashCommandsInitializer) updateExistingMarkdownCommand(
	fs afero.Fs,
	filePath, cmd, body string,
) error {
	// Read existing file
	existingContent, err := afero.ReadFile(fs, filePath)
	if err != nil {
		return err
	}

	existingStr := string(existingContent)

	// Find markers
	startIndex := findMarkerIndex(existingStr, SpectrStartMarker, 0)
	endIndex := -1
	if startIndex != -1 {
		searchOffset := startIndex + len(SpectrStartMarker)
		endIndex = findMarkerIndex(existingStr, SpectrEndMarker, searchOffset)
	}

	var finalContent string
	if startIndex == -1 || endIndex == -1 {
		// No markers found, append to end
		markerContent := SpectrStartMarker + newline + body + newline + SpectrEndMarker + newline
		finalContent = existingStr + newlineDouble + markerContent
	} else {
		// Replace content between markers
		before := existingStr[:startIndex]
		after := existingStr[endIndex+len(SpectrEndMarker):]
		finalContent = before + SpectrStartMarker + newlineDouble + body + newlineDouble + SpectrEndMarker + after
	}

	return afero.WriteFile(fs, filePath, []byte(finalContent), filePerm)
}

func (s *slashCommandsInitializer) initTOMLCommand(fs afero.Fs, filePath, cmd, body string) error {
	// Get description for this command
	description := s.getTOMLDescription(cmd)

	// Generate TOML content
	content := generateSlashTOMLContent(description, body)

	return afero.WriteFile(fs, filePath, []byte(content), filePerm)
}

func (s *slashCommandsInitializer) getTOMLDescription(cmd string) string {
	switch cmd {
	case slashCommandProposal:
		return tomlDescriptionProposal
	case slashCommandApply:
		return tomlDescriptionApply
	default:
		return fmt.Sprintf("Spectr %s command.", cmd)
	}
}

func generateSlashTOMLContent(description, prompt string) string {
	// Escape the prompt for TOML multiline string
	escapedPrompt := strings.ReplaceAll(prompt, `\`, `\\`)
	escapedPrompt = strings.ReplaceAll(escapedPrompt, `"`, `\"`)

	return fmt.Sprintf(`# Spectr command for Gemini CLI
description = "%s"
prompt = """
%s
"""
`, description, escapedPrompt)
}

func (s *slashCommandsInitializer) IsSetup(fs afero.Fs, cfg *Config) bool {
	proposalPath := filepath.Join(s.dir, slashCommandProposal+s.ext)
	applyPath := filepath.Join(s.dir, slashCommandApply+s.ext)

	proposalExists, err := afero.Exists(fs, proposalPath)
	if err != nil || !proposalExists {
		return false
	}

	applyExists, err := afero.Exists(fs, applyPath)
	if err != nil || !applyExists {
		return false
	}

	return true
}

func (s *slashCommandsInitializer) Path() string {
	return s.dir
}

func (s *slashCommandsInitializer) IsGlobal() bool {
	return s.global
}
