// Package initializers provides built-in initializers for the provider system.
//
// This file contains the ConfigFileInitializer, which creates or updates
// instruction files (like CLAUDE.md) using marker-based content insertion.
package initializers

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/connerohnesorge/spectr/internal/initialize/providers"
	"github.com/spf13/afero"
)

const (
	// File permission for config files (rw-r--r--)
	configFilePerm = 0644

	// Markers for managed content blocks
	spectrStartMarker = "<!-- spectr:START -->"
	spectrEndMarker   = "<!-- spectr:END -->"

	// String constants
	newline       = "\n"
	newlineDouble = "\n\n"
)

// ConfigFileInitializer creates or updates instruction files with marker-based content.
//
// ConfigFileInitializer manages configuration files (like CLAUDE.md, .cursorrules) by
// inserting or updating content between special markers. This allows spectr to manage
// its section of the file while preserving user content outside the markers.
//
// # Marker System
//
// Content is wrapped with markers:
//
//	<!-- spectr:START -->
//	[spectr-managed content here]
//	<!-- spectr:END -->
//
// # Scenarios
//
// 1. **New file**: If the file doesn't exist, it's created with the content between markers.
//
// 2. **Existing file with markers**: Content between markers is replaced, preserving
// content outside the markers.
//
// 3. **Existing file without markers**: Markers and content are appended to the end
// of the file.
//
// # Execution Order
//
// ConfigFileInitializer has priority 2 in the initializer ordering, meaning it runs
// after DirectoryInitializer (priority 1) but before SlashCommandsInitializer
// (priority 3). This ensures parent directories exist before files are written.
//
// # Idempotency
//
// ConfigFileInitializer is idempotent: calling Init() multiple times updates the
// content between markers to the latest version without duplicating content or
// affecting content outside the markers.
//
// # Example Usage
//
//	func (p *ClaudeProvider) Initializers(ctx context.Context) []providers.Initializer {
//	    return []providers.Initializer{
//	        initializers.NewDirectoryInitializer(".claude/commands/spectr"),
//	        initializers.NewConfigFileInitializer("CLAUDE.md"),
//	        // ... other initializers
//	    }
//	}
//
// # Global Config Files
//
// For configuration files in the user's home directory (e.g., ~/.config/aider/),
// use NewGlobalConfigFileInitializer instead:
//
//	init := initializers.NewGlobalConfigFileInitializer(".config/aider/settings.yml")
type ConfigFileInitializer struct {
	// path is the file path (relative to fs root)
	path string
	// global indicates whether to use globalFs (true) or projectFs (false)
	global bool
}

// NewConfigFileInitializer creates a new ConfigFileInitializer for a project config file.
//
// The path should be relative to the project root. The initializer will create
// parent directories if they don't exist.
//
// Example:
//
//	init := NewConfigFileInitializer("CLAUDE.md")
//	init := NewConfigFileInitializer(".cursorrules")
//	init := NewConfigFileInitializer("docs/AI_INSTRUCTIONS.md")
func NewConfigFileInitializer(path string) *ConfigFileInitializer {
	return &ConfigFileInitializer{
		path:   path,
		global: false,
	}
}

// NewGlobalConfigFileInitializer creates a new ConfigFileInitializer for a global config file.
//
// The path should be relative to the user's home directory. This is used for
// tools that store configuration globally (e.g., Aider uses ~/.config/aider/).
//
// Example:
//
//	init := NewGlobalConfigFileInitializer(".config/aider/aider.conf.yml")
func NewGlobalConfigFileInitializer(path string) *ConfigFileInitializer {
	return &ConfigFileInitializer{
		path:   path,
		global: true,
	}
}

// Init creates or updates the config file with marker-based content.
//
// Init handles three scenarios:
//
// 1. If the file doesn't exist, it creates it with the instruction pointer content
// wrapped in markers.
//
// 2. If the file exists and contains markers, it replaces the content between
// the markers with the new instruction pointer content.
//
// 3. If the file exists but has no markers, it appends the markers and content
// to the end of the file.
//
// The content is rendered using the TemplateManager's RenderInstructionPointer method.
//
// Parameters:
//   - ctx: Context for cancellation (passed to template rendering)
//   - fs: Filesystem to operate on (project or global, based on IsGlobal())
//   - cfg: Configuration containing spectr directory paths (used for template context)
//   - tm: Template manager for rendering the instruction pointer template
//
// Returns an error if file operations fail or template rendering fails.
func (c *ConfigFileInitializer) Init(ctx context.Context, fs afero.Fs, cfg *providers.Config, tm providers.TemplateManager) error {
	// Render the instruction pointer template
	templateCtx := providers.NewTemplateContext(cfg)
	content, err := tm.RenderInstructionPointer(templateCtx)
	if err != nil {
		return fmt.Errorf("failed to render instruction pointer: %w", err)
	}

	// Ensure parent directory exists
	dir := filepath.Dir(c.path)
	if dir != "" && dir != "." {
		if err := fs.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Check if file exists
	exists, err := afero.Exists(fs, c.path)
	if err != nil {
		return fmt.Errorf("failed to check if file exists: %w", err)
	}

	if !exists {
		// Create new file with markers
		return c.createNewFile(fs, content)
	}

	// Read existing file
	existingContent, err := afero.ReadFile(fs, c.path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", c.path, err)
	}

	// Update file with markers
	return c.updateExistingFile(fs, string(existingContent), content)
}

// createNewFile creates a new config file with content wrapped in markers.
func (c *ConfigFileInitializer) createNewFile(fs afero.Fs, content string) error {
	newContent := spectrStartMarker + newline + content + newline + spectrEndMarker + newline
	if err := afero.WriteFile(fs, c.path, []byte(newContent), configFilePerm); err != nil {
		return fmt.Errorf("failed to create file %s: %w", c.path, err)
	}
	return nil
}

// updateExistingFile updates an existing config file by replacing or appending
// marker-wrapped content.
func (c *ConfigFileInitializer) updateExistingFile(fs afero.Fs, existingContent, newMarkerContent string) error {
	// Find markers
	startIndex := findMarkerIndex(existingContent, spectrStartMarker, 0)
	endIndex := -1
	if startIndex != -1 {
		searchOffset := startIndex + len(spectrStartMarker)
		endIndex = findMarkerIndex(existingContent, spectrEndMarker, searchOffset)
	}

	var finalContent string
	if startIndex == -1 || endIndex == -1 {
		// No markers found, append to end
		finalContent = existingContent + newlineDouble +
			spectrStartMarker + newline + newMarkerContent + newline + spectrEndMarker + newline
	} else {
		// Replace content between markers
		before := existingContent[:startIndex]
		after := existingContent[endIndex+len(spectrEndMarker):]
		finalContent = before + spectrStartMarker + newline +
			newMarkerContent + newline + spectrEndMarker + after
	}

	if err := afero.WriteFile(fs, c.path, []byte(finalContent), configFilePerm); err != nil {
		return fmt.Errorf("failed to update file %s: %w", c.path, err)
	}
	return nil
}

// findMarkerIndex finds the index of a marker in content, starting from offset.
// Returns -1 if the marker is not found.
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

// IsSetup returns true if the config file exists and contains the spectr markers.
//
// IsSetup checks for the presence of both the start and end markers in the file.
// If the file doesn't exist or doesn't contain both markers, it returns false.
//
// Parameters:
//   - fs: Filesystem to check (project or global, based on IsGlobal())
//   - cfg: Configuration (currently unused by ConfigFileInitializer)
func (c *ConfigFileInitializer) IsSetup(fs afero.Fs, cfg *providers.Config) bool {
	exists, err := afero.Exists(fs, c.path)
	if err != nil || !exists {
		return false
	}

	content, err := afero.ReadFile(fs, c.path)
	if err != nil {
		return false
	}

	contentStr := string(content)
	return strings.Contains(contentStr, spectrStartMarker) &&
		strings.Contains(contentStr, spectrEndMarker)
}

// Path returns the file path this initializer manages.
//
// Path returns the config file path, which is used for deduplication.
// When multiple providers return ConfigFileInitializers with the same Path(),
// only one will be executed.
//
// Example return values:
//   - "CLAUDE.md"
//   - ".cursorrules"
//   - ".config/aider/aider.conf.yml"
func (c *ConfigFileInitializer) Path() string {
	return c.path
}

// IsGlobal returns true if this initializer uses the global filesystem.
//
// Project initializers (IsGlobal() == false) operate on files within
// the project directory (e.g., CLAUDE.md).
//
// Global initializers (IsGlobal() == true) operate on files in the
// user's home directory (e.g., ~/.config/aider/aider.conf.yml).
func (c *ConfigFileInitializer) IsGlobal() bool {
	return c.global
}
