package initialize

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
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
// embedded templates loaded. Templates are named by their full path
// relative to the templates directory (e.g., "spectr/AGENTS.md.tmpl").
func NewTemplateManager() (*TemplateManager, error) {
	root := template.New("")

	// Walk through all template files and parse them with their full path as the name
	err := fs.WalkDir(templateFS, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".tmpl" {
			return nil
		}

		// Read the template content
		content, err := fs.ReadFile(templateFS, path)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", path, err)
		}

		// Create template with full path as name (e.g., "templates/spectr/AGENTS.md.tmpl")
		_, err = root.New(path).Parse(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", path, err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}

	return &TemplateManager{
		templates: root,
	}, nil
}

// resolveTemplatePath resolves the template path with provider-first lookup and fallback.
// If providerID is non-empty and a template exists at templates/{providerID}/{templateName},
// it returns that path. Otherwise, it returns templates/{fallbackDir}/{templateName}.
//
// Parameters:
//   - providerID: the provider identifier (e.g., "claude-code", "crush"). Can be empty.
//   - templateName: the template filename (e.g., "AGENTS.md.tmpl")
//   - fallbackDir: the fallback directory name (e.g., "spectr", "tools")
//
// Returns the resolved template path to use with ExecuteTemplate.
func resolveTemplatePath(providerID, templateName, fallbackDir string) string {
	// If provider ID is provided, check if provider-specific template exists
	if providerID != "" {
		providerPath := fmt.Sprintf("templates/%s/%s", providerID, templateName)
		// Check if the template file exists in the embedded FS
		if templateExists(providerPath) {
			return providerPath
		}
	}

	// Fall back to generic template
	return fmt.Sprintf("templates/%s/%s", fallbackDir, templateName)
}

// templateExists checks if a template file exists in the embedded filesystem.
func templateExists(path string) bool {
	_, err := fs.Stat(templateFS, path)
	return err == nil
}

// RenderProject renders the project.md template with the given context
func (tm *TemplateManager) RenderProject(ctx ProjectContext) (string, error) {
	var buf bytes.Buffer
	// Template names use full path (e.g., "templates/spectr/project.md.tmpl")
	err := tm.templates.ExecuteTemplate(&buf, "templates/spectr/project.md.tmpl", ctx)
	if err != nil {
		return "", fmt.Errorf("failed to render project template: %w", err)
	}

	return buf.String(), nil
}

// RenderAgents renders the AGENTS.md template with the given template context.
// The context provides path variables for dynamic directory names.
//
// If providerID is non-empty, it looks for templates/{providerID}/AGENTS.md.tmpl first,
// falling back to templates/spectr/AGENTS.md.tmpl if not found.
// If providerID is empty, it uses the generic template directly.
func (tm *TemplateManager) RenderAgents(
	ctx providers.TemplateContext,
	providerID string,
) (string, error) {
	templatePath := resolveTemplatePath(providerID, "AGENTS.md.tmpl", "spectr")

	var buf bytes.Buffer
	err := tm.templates.ExecuteTemplate(&buf, templatePath, ctx)
	if err != nil {
		return "", fmt.Errorf("failed to render agents template: %w", err)
	}

	return buf.String(), nil
}

// RenderInstructionPointer renders the instruction-pointer.md template.
// This is a short pointer that directs AI assistants to read the AGENTS.md file.
// The context provides path variables for dynamic directory names.
//
// If providerID is non-empty, it looks for templates/{providerID}/instruction-pointer.md.tmpl first,
// falling back to templates/spectr/instruction-pointer.md.tmpl if not found.
// If providerID is empty, it uses the generic template directly.
func (tm *TemplateManager) RenderInstructionPointer(
	ctx providers.TemplateContext,
	providerID string,
) (string, error) {
	templatePath := resolveTemplatePath(providerID, "instruction-pointer.md.tmpl", "spectr")

	var buf bytes.Buffer
	err := tm.templates.ExecuteTemplate(&buf, templatePath, ctx)
	if err != nil {
		return "",
			fmt.Errorf(
				"failed to render instruction pointer template: %w",
				err,
			)
	}

	return buf.String(), nil
}

// RenderSlashCommand renders a slash command template with the given context.
// commandType must be one of: "proposal", "apply", "archive".
// The context provides path variables for dynamic directory names.
//
// If providerID is non-empty, it looks for templates/{providerID}/slash-{commandType}.md.tmpl first,
// falling back to templates/tools/slash-{commandType}.md.tmpl if not found.
// If providerID is empty, it uses the generic template directly.
func (tm *TemplateManager) RenderSlashCommand(
	commandType string,
	ctx providers.TemplateContext,
	providerID string,
) (string, error) {
	templateName := fmt.Sprintf("slash-%s.md.tmpl", commandType)
	templatePath := resolveTemplatePath(providerID, templateName, "tools")

	var buf bytes.Buffer
	err := tm.templates.ExecuteTemplate(&buf, templatePath, ctx)
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
	// CI workflow template is at templates/ci/spectr-ci.yml.tmpl
	err := tm.templates.ExecuteTemplate(&buf, "templates/ci/spectr-ci.yml.tmpl", nil)
	if err != nil {
		return "", fmt.Errorf("failed to render CI workflow template: %w", err)
	}

	return buf.String(), nil
}
