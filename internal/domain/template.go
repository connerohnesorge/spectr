package domain

import (
	"html/template"
)

// TemplateRef is a type-safe reference to a parsed template.
// It serves as a lightweight typed handle that can be safely passed
// between packages without creating import cycles.
// Rendering is performed by TemplateManager, not by TemplateRef itself.
type TemplateRef struct {
	Name     string             // template file name (e.g., "instruction-pointer.md.tmpl")
	Template *template.Template // pre-parsed template
}

// TemplateContext holds path-related template variables for dynamic directory names.
// Created via templateContextFromConfig(cfg) in the executor, not via defaults.
type TemplateContext struct {
	BaseDir     string // e.g., "spectr" (from cfg.SpectrDir)
	SpecsDir    string // e.g., "spectr/specs" (from cfg.SpecsDir())
	ChangesDir  string // e.g., "spectr/changes" (from cfg.ChangesDir())
	ProjectFile string // e.g., "spectr/project.md" (from cfg.ProjectFile())
	AgentsFile  string // e.g., "spectr/AGENTS.md" (from cfg.AgentsFile())
}
