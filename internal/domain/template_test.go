package domain

import (
	"strings"
	"testing"
	"text/template"
)

func TestTemplateRef_Render(t *testing.T) {
	// Create a simple template
	tmpl := template.New("test.tmpl")
	tmpl, err := tmpl.Parse("BaseDir: {{.BaseDir}}, SpecsDir: {{.SpecsDir}}")
	if err != nil {
		t.Fatalf("failed to parse template: %v", err)
	}

	ref := TemplateRef{
		Name:     "test.tmpl",
		Template: tmpl,
	}

	ctx := TemplateContext{
		BaseDir:  "spectr",
		SpecsDir: "spectr/specs",
	}

	result, err := ref.Render(&ctx)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	expected := "BaseDir: spectr, SpecsDir: spectr/specs"
	if result != expected {
		t.Errorf("Render() = %q, want %q", result, expected)
	}
}

func TestTemplateRef_Render_Error(t *testing.T) {
	// Create a template with an invalid reference
	tmpl := template.New("error.tmpl")
	tmpl, err := tmpl.Parse("{{.NonExistentField}}")
	if err != nil {
		t.Fatalf("failed to parse template: %v", err)
	}

	ref := TemplateRef{
		Name:     "error.tmpl",
		Template: tmpl,
	}

	ctx := TemplateContext{}

	_, err = ref.Render(&ctx)
	if err == nil {
		t.Error("Render() expected error for invalid template field, got nil")
	}
	if !strings.Contains(err.Error(), "failed to render template") {
		t.Errorf(
			"Render() error = %v, want error message containing 'failed to render template'",
			err,
		)
	}
}

func TestDefaultTemplateContext(t *testing.T) {
	ctx := DefaultTemplateContext()

	tests := []struct {
		name string
		got  string
		want string
	}{
		{"BaseDir", ctx.BaseDir, "spectr"},
		{"SpecsDir", ctx.SpecsDir, "spectr/specs"},
		{"ChangesDir", ctx.ChangesDir, "spectr/changes"},
		{"ProjectFile", ctx.ProjectFile, "spectr/project.md"},
		{"AgentsFile", ctx.AgentsFile, "spectr/AGENTS.md"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("DefaultTemplateContext().%s = %q, want %q", tt.name, tt.got, tt.want)
			}
		})
	}
}
