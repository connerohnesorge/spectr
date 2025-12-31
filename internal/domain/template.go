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
// ProviderTemplate is optional and allows provider-specific section overrides.
type TemplateRef struct {
	Name             string             // template file name (e.g., "instruction-pointer.md.tmpl")
	Template         *template.Template // pre-parsed template
	ProviderTemplate *template.Template // optional provider override template
}

// Render executes the template with the given context.
func (tr TemplateRef) Render(ctx *TemplateContext) (string, error) {
	if tr.ProviderTemplate == nil {
		var buf bytes.Buffer
		if err := tr.Template.ExecuteTemplate(&buf, tr.Name, ctx); err != nil {
			return "", fmt.Errorf("failed to render template %s: %w", tr.Name, err)
		}

		return buf.String(), nil
	}

	composed, err := tr.composeTemplate()
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := composed.ExecuteTemplate(&buf, "main", ctx); err != nil {
		return "", fmt.Errorf("failed to render template %s: %w", tr.Name, err)
	}

	return buf.String(), nil
}

func (tr TemplateRef) composeTemplate() (*template.Template, error) {
	composed, err := tr.Template.Clone()
	if err != nil {
		return nil, fmt.Errorf("failed to clone template %s: %w", tr.Name, err)
	}

	for _, tmpl := range tr.ProviderTemplate.Templates() {
		if tmpl.Tree == nil {
			return nil, fmt.Errorf(
				"provider template %s has no parse tree for %s",
				tmpl.Name(),
				tr.Name,
			)
		}
		if _, err := composed.AddParseTree(tmpl.Name(), tmpl.Tree); err != nil {
			return nil, fmt.Errorf("failed to merge provider template %s: %w", tr.Name, err)
		}
	}

	return composed, nil
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
