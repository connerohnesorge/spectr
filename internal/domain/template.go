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
	Name     string             // template file name (e.g., "instruction-pointer.md.tmpl")
	Template *template.Template // pre-parsed template
}

// Render executes the template with the given context.
func (tr TemplateRef) Render(ctx TemplateContext) (string, error) {
	var buf bytes.Buffer
	if err := tr.Template.ExecuteTemplate(&buf, tr.Name, ctx); err != nil {
		return "", fmt.Errorf("failed to render template %s: %w", tr.Name, err)
	}

	return buf.String(), nil
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
