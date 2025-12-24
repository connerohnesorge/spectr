// Package domain contains shared domain types used across packages.
package domain

import (
	"embed"
	"fmt"
	"sync"
	"text/template"
)

//go:embed templates/*.tmpl
var slashTemplateFS embed.FS

var (
	slashTemplates     *template.Template
	slashTemplatesOnce sync.Once
	errSlashTemplates  error
)

// SlashCommand represents a type-safe slash command identifier.
type SlashCommand int

const (
	SlashProposal SlashCommand = iota
	SlashApply
)

// templateNames maps slash commands to their template file names.
var templateNames = map[SlashCommand]string{
	SlashProposal: "slash-proposal.md.tmpl",
	SlashApply:    "slash-apply.md.tmpl",
}

// String returns the command name for debugging.
func (s SlashCommand) String() string {
	names := []string{"proposal", "apply"}
	if int(s) < len(names) {
		return names[s]
	}

	return "unknown"
}

// TemplateRef returns a type-safe reference to the slash command's template.
// Templates are parsed once on first access.
func (s SlashCommand) TemplateRef() (TemplateRef, error) {
	slashTemplatesOnce.Do(func() {
		slashTemplates, errSlashTemplates = template.ParseFS(
			slashTemplateFS,
			"templates/*.tmpl",
		)
	})
	if errSlashTemplates != nil {
		return TemplateRef{}, fmt.Errorf(
			"failed to parse slash templates: %w",
			errSlashTemplates,
		)
	}

	name, ok := templateNames[s]
	if !ok {
		return TemplateRef{}, fmt.Errorf(
			"unknown slash command: %d",
			s,
		)
	}

	return TemplateRef{
		Name:     name,
		Template: slashTemplates,
	}, nil
}
