package providerkit

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
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
	tmpl, err := template.ParseFS(templateFS, "templates/**/*.tmpl")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	return &TemplateManager{
		templates: tmpl,
	}, nil
}

// ProjectContext holds template variables for rendering project.md
type ProjectContext struct {
	// ProjectName is the name of the project
	ProjectName string
	// Description is the project description/purpose
	Description string
	// TechStack is the list of technologies used
	TechStack []string
	// Conventions are the project conventions (unused in template currently)
	Conventions string
}

// SpecContext holds template variables for rendering spec.md
type SpecContext struct {
	// CapabilityName is the name of the capability/spec
	CapabilityName string
}

// ProposalContext holds template variables for rendering proposal.md
type ProposalContext struct {
	// ChangeName is the name of the change proposal
	ChangeName string
}

// RenderProject renders the project.md template with the given context
func (tm *TemplateManager) RenderProject(ctx ProjectContext) (string, error) {
	var buf bytes.Buffer
	// Template names in ParseFS include the full path from the embed directive
	err := tm.templates.ExecuteTemplate(&buf, "project.md.tmpl", ctx)
	if err != nil {
		return "", fmt.Errorf("failed to render project template: %w", err)
	}

	return buf.String(), nil
}

// RenderAgents renders the AGENTS.md template (static, no variables)
func (tm *TemplateManager) RenderAgents() (string, error) {
	var buf bytes.Buffer
	err := tm.templates.ExecuteTemplate(&buf, "AGENTS.md.tmpl", nil)
	if err != nil {
		return "", fmt.Errorf("failed to render agents template: %w", err)
	}

	return buf.String(), nil
}

// RenderSpec renders a spec.md template with the given context
func (tm *TemplateManager) RenderSpec(ctx SpecContext) (string, error) {
	var buf bytes.Buffer
	err := tm.templates.ExecuteTemplate(&buf, "spec.md.tmpl", ctx)
	if err != nil {
		return "", fmt.Errorf("failed to render spec template: %w", err)
	}

	return buf.String(), nil
}

// RenderProposal renders a proposal.md template with the given context
func (tm *TemplateManager) RenderProposal(ctx ProposalContext) (string, error) {
	var buf bytes.Buffer
	err := tm.templates.ExecuteTemplate(&buf, "proposal.md.tmpl", ctx)
	if err != nil {
		return "", fmt.Errorf("failed to render proposal template: %w", err)
	}

	return buf.String(), nil
}

// RenderSlashCommand renders a slash command template
// commandType must be one of: "proposal", "apply", "archive"
func (tm *TemplateManager) RenderSlashCommand(
	commandType string,
) (string, error) {
	templateName := fmt.Sprintf("slash-%s.md.tmpl", commandType)
	var buf bytes.Buffer
	err := tm.templates.ExecuteTemplate(&buf, templateName, nil)
	if err != nil {
		return "", fmt.Errorf(
			"failed to render slash command template %s: %w",
			commandType,
			err,
		)
	}

	return buf.String(), nil
}
