package initialize

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"

	"github.com/connerohnesorge/spectr/internal/initialize/providers"
)

//go:embed templates/**/*.tmpl
var templateFS embed.FS

// TemplateManager manages embedded templates for initialization
type TemplateManager struct {
	templates *template.Template
}

// NewTemplateManager creates a new template manager with all
// embedded templates loaded
func NewTemplateManager() (*TemplateManager, error) {
	// Parse all embedded templates
	tmpl, err := template.ParseFS(
		templateFS,
		"templates/**/*.tmpl",
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse templates: %w",
			err,
		)
	}

	return &TemplateManager{
		templates: tmpl,
	}, nil
}

// RenderProject renders the project.md template with the given context
func (tm *TemplateManager) RenderProject(
	ctx ProjectContext,
) (string, error) {
	var buf bytes.Buffer
	// Template names in ParseFS include the full path from the embed directive
	err := tm.templates.ExecuteTemplate(
		&buf,
		"project.md.tmpl",
		ctx,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to render project template: %w",
			err,
		)
	}

	return buf.String(), nil
}

// RenderAgents renders the AGENTS.md template with the given template context
// The context provides path variables for dynamic directory names
func (tm *TemplateManager) RenderAgents(
	ctx *providers.TemplateContext,
) (string, error) {
	var buf bytes.Buffer
	err := tm.templates.ExecuteTemplate(
		&buf,
		"AGENTS.md.tmpl",
		ctx,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to render agents template: %w",
			err,
		)
	}

	return buf.String(), nil
}

// RenderInstructionPointer renders the instruction-pointer.md template
// This is a short pointer that directs AI assistants to read the AGENTS.md file
// The context provides path variables for dynamic directory names
func (tm *TemplateManager) RenderInstructionPointer(
	ctx *providers.TemplateContext,
) (string, error) {
	var buf bytes.Buffer
	err := tm.templates.ExecuteTemplate(
		&buf,
		"instruction-pointer.md.tmpl",
		ctx,
	)
	if err != nil {
		return "",
			fmt.Errorf(
				"failed to render instruction pointer template: %w",
				err,
			)
	}

	return buf.String(), nil
}

// RenderSlashCommand renders a slash command template with the given context
// commandType must be one of: "proposal", "apply", "archive"
// The context provides path variables for dynamic directory names
func (tm *TemplateManager) RenderSlashCommand(
	commandType string,
	ctx *providers.TemplateContext,
) (string, error) {
	templateName := fmt.Sprintf(
		"slash-%s.md.tmpl",
		commandType,
	)
	var buf bytes.Buffer
	err := tm.templates.ExecuteTemplate(
		&buf,
		templateName,
		ctx,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to render slash command template %s: %w",
			commandType,
			err,
		)
	}

	return buf.String(), nil
}

// RenderCIWorkflow renders the spectr-ci.yml template for GitHub Actions
// This template has no variables and returns the CI workflow configuration
func (tm *TemplateManager) RenderCIWorkflow() (string, error) {
	var buf bytes.Buffer
	err := tm.templates.ExecuteTemplate(
		&buf,
		"spectr-ci.yml.tmpl",
		nil,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to render CI workflow template: %w",
			err,
		)
	}

	return buf.String(), nil
}
