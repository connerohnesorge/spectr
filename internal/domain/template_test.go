package domain

import (
	"strings"
	"testing"
	"text/template"
)

func TestTemplateRef_Render_NoProvider(t *testing.T) {
	base := template.Must(template.New("base.tmpl").Parse(
		`{{define "guardrails"}}base-guardrails{{end}}
{{define "steps"}}base-steps{{end}}
{{define "reference"}}base-reference{{end}}
{{define "main"}}{{template "guardrails" .}}|{{template "steps" .}}|{{template "reference" .}}{{end}}
{{template "main" .}}`,
	))

	ref := TemplateRef{
		Name:     "base.tmpl",
		Template: base,
	}

	result, err := ref.Render(&TemplateContext{})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	result = strings.TrimSpace(result)
	expected := "base-guardrails|base-steps|base-reference"
	if result != expected {
		t.Errorf("Render() = %q, want %q", result, expected)
	}
}

func TestTemplateRef_Render_WithProvider(t *testing.T) {
	base := template.Must(template.New("base.tmpl").Parse(
		`{{define "guardrails"}}base-guardrails{{end}}
{{define "steps"}}base-steps{{end}}
{{define "reference"}}base-reference{{end}}
{{define "main"}}{{template "guardrails" .}}|{{template "steps" .}}|{{template "reference" .}}{{end}}
{{template "main" .}}`,
	))
	provider := template.Must(template.New("provider.tmpl").Parse(
		`{{define "guardrails"}}provider-guardrails{{end}}`,
	))

	ref := TemplateRef{
		Name:             "base.tmpl",
		Template:         base,
		ProviderTemplate: provider,
	}

	result, err := ref.Render(&TemplateContext{})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	expected := "provider-guardrails|base-steps|base-reference"
	if result != expected {
		t.Errorf("Render() = %q, want %q", result, expected)
	}
}

func TestTemplateRef_Render_PartialOverride(t *testing.T) {
	base := template.Must(template.New("base.tmpl").Parse(
		`{{define "guardrails"}}base-guardrails{{end}}
{{define "steps"}}base-steps{{end}}
{{define "reference"}}base-reference{{end}}
{{define "main"}}{{template "guardrails" .}}|{{template "steps" .}}|{{template "reference" .}}{{end}}
{{template "main" .}}`,
	))
	provider := template.Must(template.New("provider.tmpl").Parse(
		`{{define "steps"}}provider-steps{{end}}`,
	))

	ref := TemplateRef{
		Name:             "base.tmpl",
		Template:         base,
		ProviderTemplate: provider,
	}

	result, err := ref.Render(&TemplateContext{})
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	expected := "base-guardrails|provider-steps|base-reference"
	if result != expected {
		t.Errorf("Render() = %q, want %q", result, expected)
	}
}

func TestTemplateRef_composeTemplate(t *testing.T) {
	base := template.Must(template.New("base.tmpl").Parse(
		`{{define "guardrails"}}base-guardrails{{end}}
{{define "steps"}}base-steps{{end}}
{{define "reference"}}base-reference{{end}}
{{define "main"}}{{template "guardrails" .}}|{{template "steps" .}}|{{template "reference" .}}{{end}}
{{template "main" .}}`,
	))
	provider := template.Must(template.New("provider.tmpl").Parse(
		`{{define "reference"}}provider-reference{{end}}`,
	))

	ref := TemplateRef{
		Name:             "base.tmpl",
		Template:         base,
		ProviderTemplate: provider,
	}

	composed, err := ref.composeTemplate()
	if err != nil {
		t.Fatalf("composeTemplate() error = %v", err)
	}

	var buf strings.Builder
	if err := composed.ExecuteTemplate(&buf, "main", &TemplateContext{}); err != nil {
		t.Fatalf("ExecuteTemplate() error = %v", err)
	}

	expected := "base-guardrails|base-steps|provider-reference"
	if buf.String() != expected {
		t.Errorf("composeTemplate() result = %q, want %q", buf.String(), expected)
	}
}

func TestTemplateRef_composeTemplate_Error(t *testing.T) {
	base := template.Must(template.New("base.tmpl").Parse(`{{define "main"}}ok{{end}}`))
	provider := template.Must(
		template.New("provider.tmpl").Parse(`{{define "guardrails"}}bad{{end}}`),
	)
	provider.Tree = nil
	if guardrails := provider.Lookup("guardrails"); guardrails != nil {
		guardrails.Tree = nil
	}

	ref := TemplateRef{
		Name:             "base.tmpl",
		Template:         base,
		ProviderTemplate: provider,
	}

	_, err := ref.composeTemplate()
	if err == nil {
		t.Fatal("composeTemplate() expected error, got nil")
	}
	if !strings.Contains(err.Error(), "base.tmpl") {
		t.Errorf("composeTemplate() error = %v, want template name in error", err)
	}
}

func TestTemplateRef_Render_Error(t *testing.T) {
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
