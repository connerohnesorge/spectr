// Package templates provides type-safe template reference management.
package templates

import (
	"bytes"
	"fmt"
	"text/template"
)

// TemplateRef is a type-safe reference to a parsed template.
// It provides a Render method that executes the template with context.
type TemplateRef struct {
	// template file name (e.g., "instruction-pointer.md.tmpl")
	name string
	// pre-parsed template
	template *template.Template
}

// NewTemplateRef creates a new TemplateRef with the given name and template.
// This is a package-internal constructor used by TemplateManager accessor
// methods.
func NewTemplateRef(
	name string,
	tmpl *template.Template,
) TemplateRef {
	return TemplateRef{
		name:     name,
		template: tmpl,
	}
}

// Render executes the template with the given context.
// The context can be any type compatible with the template's expectations.
// Most templates use providers.TemplateContext, but some
// (like project.md.tmpl) use other context types (e.g., ProjectContext).
// Returns the rendered template content or an error if rendering fails.
func (tr TemplateRef) Render(
	ctx any,
) (string, error) {
	var buf bytes.Buffer
	if err := tr.template.ExecuteTemplate(&buf, tr.name, ctx); err != nil {
		return "", fmt.Errorf(
			"failed to render template %s: %w",
			tr.name,
			err,
		)
	}

	return buf.String(), nil
}
