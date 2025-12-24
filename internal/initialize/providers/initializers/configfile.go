package initializers

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"

	"github.com/connerohnesorge/spectr/internal/domain"
)

// ConfigFileInitializer creates or updates a configuration file.
//
// The file content is managed using spectr markers.
type ConfigFileInitializer struct {
	path         string
	templateName string
	isGlobal     bool
}

// NewConfigFileInitializer creates a ConfigFileInitializer.
//
// The templateName should be the template file name.
//
// Example:
//
//	NewConfigFileInitializer("CLAUDE.md", "instruction-pointer.md.tmpl")
func NewConfigFileInitializer(
	path, templateName string,
) *ConfigFileInitializer {
	return &ConfigFileInitializer{
		path:         path,
		templateName: templateName,
		isGlobal:     false,
	}
}

// Init creates or updates the config file with content from the template.
// If the file exists, it updates the content between spectr markers.
// If the file doesn't exist, it creates it with markers.
func (c *ConfigFileInitializer) Init(
	_ context.Context,
	fs afero.Fs,
	cfg *Config,
	tm any,
) (InitResult, error) {
	var result InitResult

	content, err := c.renderContent(cfg, tm)
	if err != nil {
		return result, err
	}

	// Ensure parent directory exists
	dir := filepath.Dir(c.path)
	if err := fs.MkdirAll(dir, DirPerm); err != nil {
		return result, fmt.Errorf("failed to create directory: %w", err)
	}

	exists, err := afero.Exists(fs, c.path)
	if err != nil {
		return result, err
	}

	if !exists {
		return c.createNewFile(fs, content)
	}

	return c.updateExistingFile(fs, content)
}

// renderContent renders the template content.
func (c *ConfigFileInitializer) renderContent(
	cfg *Config,
	tm any,
) (string, error) {
	templateProvider, ok := tm.(TemplateProvider)
	if !ok {
		return "", fmt.Errorf("expected TemplateProvider, got %T", tm)
	}

	templateCtx := domain.TemplateContext{
		BaseDir:     cfg.SpectrDir,
		SpecsDir:    cfg.SpecsDir(),
		ChangesDir:  cfg.ChangesDir(),
		ProjectFile: cfg.ProjectFile(),
		AgentsFile:  cfg.AgentsFile(),
	}

	templateRef := domain.TemplateRef{
		Name:     c.templateName,
		Template: templateProvider.GetTemplates(),
	}

	content, err := templateRef.Render(templateCtx)
	if err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}

	return content, nil
}

// createNewFile creates a new config file with spectr markers.
func (c *ConfigFileInitializer) createNewFile(
	fs afero.Fs,
	content string,
) (InitResult, error) {
	var result InitResult

	newContent := SpectrStartMarker + Newline +
		content + Newline +
		SpectrEndMarker + Newline

	err := afero.WriteFile(fs, c.path, []byte(newContent), FilePerm)
	if err != nil {
		return result, fmt.Errorf("failed to write file: %w", err)
	}

	result.CreatedFiles = append(result.CreatedFiles, c.path)

	return result, nil
}

// updateExistingFile updates an existing config file between markers.
func (c *ConfigFileInitializer) updateExistingFile(
	fs afero.Fs,
	content string,
) (InitResult, error) {
	var result InitResult

	existingContent, err := afero.ReadFile(fs, c.path)
	if err != nil {
		return result, fmt.Errorf("failed to read file: %w", err)
	}

	contentStr := string(existingContent)
	newContentStr, wasUpdated := updateBetweenMarkers(
		contentStr,
		content,
		SpectrStartMarker,
		SpectrEndMarker,
	)

	if wasUpdated {
		err := afero.WriteFile(fs, c.path, []byte(newContentStr), FilePerm)
		if err != nil {
			return result, fmt.Errorf("failed to write file: %w", err)
		}

		result.UpdatedFiles = append(result.UpdatedFiles, c.path)
	}

	return result, nil
}

// IsSetup returns true if the file exists and contains spectr markers.
func (c *ConfigFileInitializer) IsSetup(fs afero.Fs, _ *Config) bool {
	exists, err := afero.Exists(fs, c.path)
	if err != nil || !exists {
		return false
	}

	// Check if file has markers
	content, err := afero.ReadFile(fs, c.path)
	if err != nil {
		return false
	}

	contentStr := string(content)
	hasStartMarker := strings.Contains(contentStr, SpectrStartMarker)
	hasEndMarker := strings.Contains(contentStr, SpectrEndMarker)

	return hasStartMarker && hasEndMarker
}

// Path returns the config file path.
func (c *ConfigFileInitializer) Path() string {
	return c.path
}

// IsGlobal returns whether this initializer operates on global files.
func (c *ConfigFileInitializer) IsGlobal() bool {
	return c.isGlobal
}

// updateBetweenMarkers updates content between markers in a string.
func updateBetweenMarkers(
	contentStr, newContent, startMarker, endMarker string,
) (string, bool) {
	startIndex := strings.Index(contentStr, startMarker)
	endIndex := -1

	if startIndex != -1 {
		searchOffset := startIndex + len(startMarker)
		endIndex = strings.Index(contentStr[searchOffset:], endMarker)

		if endIndex != -1 {
			endIndex += searchOffset
		}
	}

	if startIndex == -1 || endIndex == -1 {
		// No markers found - append to end
		updated := contentStr + Newline + Newline +
			startMarker + Newline +
			newContent + Newline +
			endMarker + Newline

		return updated, true
	}

	// Replace content between markers
	before := contentStr[:startIndex]
	after := contentStr[endIndex+len(endMarker):]
	updated := before +
		startMarker + Newline +
		newContent + Newline +
		endMarker +
		after

	// Check if content actually changed
	wasUpdated := updated != contentStr

	return updated, wasUpdated
}
