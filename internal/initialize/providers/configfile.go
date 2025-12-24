package providers

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/connerohnesorge/spectr/internal/initialize/templates"
	"github.com/spf13/afero"
)

const (
	// File permission for created config files
	configFilePerm = 0o644

	// Directory permission for creating parent directories
	dirPerm = 0o755

	// Marker constants for managing config file updates
	spectrStartMarker = "<!-- spectr:START -->"
	spectrEndMarker   = "<!-- spectr:END -->"

	// Common string constants
	newline       = "\n"
	doubleNewline = "\n\n"
)

// TemplateGetter is a function that retrieves a TemplateRef from
// TemplateManager. This allows compile-time checked template selection
// via method references.
//
// Example usage:
//
//	NewConfigFileInitializer(
//	  "CLAUDE.md",
//	  TemplateManager.InstructionPointer,
//	)
//
// The compiler will catch typos in the method name, preventing runtime
// errors. Returns any due to import cycle constraints; will be cast to
// templates.TemplateRef.
type TemplateGetter func(TemplateManager) any

// ConfigFileInitializer creates or updates a configuration file with
// marker-based updates. It uses the Spectr marker system
// (<!-- spectr:START --> ... <!-- spectr:END -->) to allow safe updates
// of managed content while preserving user modifications outside the
// markers.
//
// Behavior:
//   - If file doesn't exist: Creates file with markers
//   - If file exists without markers: Appends markers with content
//   - If file exists with markers: Updates content between markers
type ConfigFileInitializer struct {
	path        string
	getTemplate TemplateGetter
	isGlobal    bool
}

// NewConfigFileInitializer creates a ConfigFileInitializer for a
// project-relative config file.
//
// Parameters:
//   - path: Relative path to the config file (e.g., "CLAUDE.md")
//   - getTemplate: Function that gets the template reference from
//     TemplateManager
//
// Example:
//
//	NewConfigFileInitializer(
//	  "CLAUDE.md",
//	  (*TemplateManager).InstructionPointer,
//	)
func NewConfigFileInitializer(
	path string,
	getTemplate TemplateGetter,
) *ConfigFileInitializer {
	return &ConfigFileInitializer{
		path:        path,
		getTemplate: getTemplate,
		isGlobal:    false,
	}
}

// NewGlobalConfigFileInitializer creates a ConfigFileInitializer for a
// global config file.
//
// Parameters:
//   - path: Relative path to config file (relative to home directory)
//   - getTemplate: Function that gets template reference from
//     TemplateManager
func NewGlobalConfigFileInitializer(
	path string,
	getTemplate TemplateGetter,
) *ConfigFileInitializer {
	return &ConfigFileInitializer{
		path:        path,
		getTemplate: getTemplate,
		isGlobal:    true,
	}
}

// Init creates or updates the config file with marker-based content.
func (c *ConfigFileInitializer) Init(
	_ context.Context,
	fs afero.Fs,
	cfg *Config,
	tm TemplateManager,
) (InitResult, error) {
	// Get the template reference using the getter function
	templateRefRaw := c.getTemplate(tm)

	// Cast to TemplateRef (getter returns templates.TemplateRef
	// stored as interface{})
	templateRef, ok := templateRefRaw.(templates.TemplateRef)
	if !ok {
		return InitResult{}, errors.New(
			"template getter did not return a TemplateRef",
		)
	}

	// Create template context from config
	templateCtx := TemplateContext{
		BaseDir:     cfg.SpectrDir,
		SpecsDir:    cfg.SpecsDir(),
		ChangesDir:  cfg.ChangesDir(),
		ProjectFile: cfg.ProjectFile(),
		AgentsFile:  cfg.AgentsFile(),
	}

	// Render the template content
	content, err := templateRef.Render(
		templateCtx,
	)
	if err != nil {
		return InitResult{}, fmt.Errorf(
			"failed to render template: %w",
			err,
		)
	}

	// Check if file exists
	exists, err := afero.Exists(fs, c.path)
	if err != nil {
		return InitResult{}, fmt.Errorf(
			"failed to check file existence: %w",
			err,
		)
	}

	if !exists {
		// Create new file with markers
		return c.createFile(fs, content)
	}

	// Update existing file with markers
	return c.updateFile(fs, content)
}

// createFile creates a new config file with marker-wrapped content.
func (c *ConfigFileInitializer) createFile(
	fs afero.Fs,
	content string,
) (InitResult, error) {
	// Ensure parent directory exists
	dir := filepath.Dir(c.path)
	if dir != "" && dir != "." {
		if err := fs.MkdirAll(dir, dirPerm); err != nil {
			return InitResult{}, fmt.Errorf(
				"failed to create directory: %w",
				err,
			)
		}
	}

	// Create file with markers
	fileContent := spectrStartMarker + newline + content +
		newline + spectrEndMarker + newline
	err := afero.WriteFile(
		fs,
		c.path,
		[]byte(fileContent),
		configFilePerm,
	)
	if err != nil {
		return InitResult{}, fmt.Errorf(
			"failed to write file: %w",
			err,
		)
	}

	return InitResult{
		CreatedFiles: []string{c.path},
	}, nil
}

// updateFile updates an existing config file, replacing content
// between markers.
func (c *ConfigFileInitializer) updateFile(
	fs afero.Fs,
	content string,
) (InitResult, error) {
	// Read existing file
	existingBytes, err := afero.ReadFile(
		fs,
		c.path,
	)
	if err != nil {
		return InitResult{}, fmt.Errorf(
			"failed to read file: %w",
			err,
		)
	}

	existingContent := string(existingBytes)

	// Find markers
	startIdx := strings.Index(
		existingContent,
		spectrStartMarker,
	)
	endIdx := -1
	if startIdx != -1 {
		searchOffset := startIdx + len(
			spectrStartMarker,
		)
		relativeEndIdx := strings.Index(
			existingContent[searchOffset:],
			spectrEndMarker,
		)
		if relativeEndIdx != -1 {
			endIdx = searchOffset + relativeEndIdx
		}
	}

	var newContent string
	if startIdx == -1 || endIdx == -1 {
		// No markers found, append to end
		newContent = existingContent
		if !strings.HasSuffix(
			newContent,
			doubleNewline,
		) {
			if strings.HasSuffix(
				newContent,
				newline,
			) {
				newContent += newline
			} else {
				newContent += doubleNewline
			}
		}
		newContent += spectrStartMarker + newline + content +
			newline + spectrEndMarker + newline
	} else {
		// Replace content between markers
		before := existingContent[:startIdx]
		after := existingContent[endIdx+len(spectrEndMarker):]
		newContent = before + spectrStartMarker + newline + content +
			newline + spectrEndMarker + after
	}

	// Write updated content
	err = afero.WriteFile(
		fs,
		c.path,
		[]byte(newContent),
		configFilePerm,
	)
	if err != nil {
		return InitResult{}, fmt.Errorf(
			"failed to write file: %w",
			err,
		)
	}

	return InitResult{
		UpdatedFiles: []string{c.path},
	}, nil
}

// IsSetup returns true if the config file exists.
func (c *ConfigFileInitializer) IsSetup(
	fs afero.Fs,
	_ *Config,
) bool {
	exists, err := afero.Exists(fs, c.path)
	if err != nil {
		return false
	}

	return exists
}

// Path returns the config file path this initializer manages.
func (c *ConfigFileInitializer) Path() string {
	return c.path
}

// IsGlobal returns true if this initializer uses globalFs instead of projectFs.
func (c *ConfigFileInitializer) IsGlobal() bool {
	return c.isGlobal
}
