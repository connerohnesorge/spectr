package templates_test

import (
	"testing"
	"text/template"

	"github.com/connerohnesorge/spectr/internal/initialize/providers"
	"github.com/connerohnesorge/spectr/internal/initialize/templates"
)

func TestTemplateRefRender(t *testing.T) {
	// Create a simple test template
	tmpl, err := template.New("test.tmpl").
		Parse("Hello {{.BaseDir}}")
	if err != nil {
		t.Fatalf(
			"Failed to create test template: %v",
			err,
		)
	}

	// Create TemplateRef
	ref := templates.NewTemplateRef(
		"test.tmpl",
		tmpl,
	)

	// Create context
	ctx := providers.TemplateContext{
		BaseDir:     "mydir",
		SpecsDir:    "mydir/specs",
		ChangesDir:  "mydir/changes",
		ProjectFile: "mydir/project.md",
		AgentsFile:  "mydir/AGENTS.md",
	}

	// Render
	content, err := ref.Render(ctx)
	if err != nil {
		t.Fatalf(
			"Failed to render template: %v",
			err,
		)
	}

	expected := "Hello mydir"
	if content != expected {
		t.Errorf(
			"Render() = %q, want %q",
			content,
			expected,
		)
	}
}

func TestSlashCommandString(t *testing.T) {
	tests := []struct {
		cmd  templates.SlashCommand
		want string
	}{
		{templates.SlashProposal, "proposal"},
		{templates.SlashApply, "apply"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.cmd.String(); got != tt.want {
				t.Errorf(
					"SlashCommand.String() = %q, want %q",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestSlashCommandTemplateName(t *testing.T) {
	tests := []struct {
		cmd  templates.SlashCommand
		want string
	}{
		{
			templates.SlashProposal,
			"slash-proposal.md.tmpl",
		},
		{
			templates.SlashApply,
			"slash-apply.md.tmpl",
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.cmd.String(),
			func(t *testing.T) {
				if got := tt.cmd.TemplateName(); got != tt.want {
					t.Errorf(
						"SlashCommand.TemplateName() = %q, want %q",
						got,
						tt.want,
					)
				}
			},
		)
	}
}
