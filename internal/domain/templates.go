package domain

import "embed"

// TemplateFS contains embedded slash command templates.
// These templates are shared across all providers.
//
//go:embed templates/*.tmpl
var TemplateFS embed.FS
