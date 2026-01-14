// Package domain contains shared domain types used across the Spectr codebase.
// This package has no internal dependencies to avoid import cycles.
package domain

import (
	"bytes"
	"fmt"
	"text/template"
)

// TemplateRef is a type-safe reference to a parsed template.
// It can be safely passed between packages without creating import cycles.
type TemplateRef struct {
	Name      string               // template file name (e.g., "instruction-pointer.md.tmpl")
	Template  *template.Template   // pre-parsed template
	Command   *SlashCommand        // slash command type (for frontmatter lookup), nil if not a slash command
	Overrides *FrontmatterOverride // optional frontmatter modifications
}

// Render executes the template with the given context.
// If Command is set, it assembles frontmatter from BaseSlashCommandFrontmatter
// and applies any Overrides before prepending to the template body.
// Returns an error if Overrides contains unknown frontmatter keys.
func (tr TemplateRef) Render(ctx *TemplateContext) (string, error) {
	// Validate overrides before rendering to catch typos early
	if err := ValidateFrontmatterOverride(tr.Overrides); err != nil {
		return "", fmt.Errorf("invalid frontmatter overrides for template %s: %w", tr.Name, err)
	}

	var buf bytes.Buffer
	if err := tr.Template.ExecuteTemplate(&buf, tr.Name, ctx); err != nil {
		return "", fmt.Errorf("failed to render template %s: %w", tr.Name, err)
	}

	body := buf.String()

	// If this is a slash command template, assemble frontmatter
	if tr.Command != nil {
		// Get base frontmatter for this command
		fm := GetBaseFrontmatter(*tr.Command)

		// Apply overrides if present
		if tr.Overrides != nil {
			fm = ApplyFrontmatterOverrides(fm, tr.Overrides)
		}

		// Only render frontmatter if we have fields
		if len(fm) > 0 {
			return RenderFrontmatter(fm, body)
		}
	}

	return body, nil
}

// TemplateContext holds path-related template variables for dynamic directory names.
type TemplateContext struct {
	BaseDir     string // e.g., "spectr"
	SpecsDir    string // e.g., "spectr/specs"
	ChangesDir  string // e.g., "spectr/changes"
	ProjectFile string // e.g., "spectr/project.md"
	AgentsFile  string // e.g., "spectr/AGENTS.md"
}

// DefaultTemplateContext returns a TemplateContext with default values.
func DefaultTemplateContext() TemplateContext {
	return TemplateContext{
		BaseDir:     "spectr",
		SpecsDir:    "spectr/specs",
		ChangesDir:  "spectr/changes",
		ProjectFile: "spectr/project.md",
		AgentsFile:  "spectr/AGENTS.md",
	}
}
