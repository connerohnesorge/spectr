package domain

import (
	"strings"
	"testing"
	"text/template"
)

func TestTemplateRef_Render(t *testing.T) {
	// Create a simple test template
	tmplContent := `Base: {{.BaseDir}}
Specs: {{.SpecsDir}}
Changes: {{.ChangesDir}}
Project: {{.ProjectFile}}
Agents: {{.AgentsFile}}`

	tmpl, err := template.New("test.tmpl").Parse(tmplContent)
	if err != nil {
		t.Fatalf("failed to parse test template: %v", err)
	}

	ref := TemplateRef{
		Name:     "test.tmpl",
		Template: tmpl,
	}

	ctx := TemplateContext{
		BaseDir:     "spectr",
		SpecsDir:    "spectr/specs",
		ChangesDir:  "spectr/changes",
		ProjectFile: "spectr/project.md",
		AgentsFile:  "spectr/AGENTS.md",
	}

	result, err := ref.Render(ctx)
	if err != nil {
		t.Fatalf("Render() failed: %v", err)
	}

	// Verify all fields are rendered correctly
	expectedLines := []string{
		"Base: spectr",
		"Specs: spectr/specs",
		"Changes: spectr/changes",
		"Project: spectr/project.md",
		"Agents: spectr/AGENTS.md",
	}

	for _, expected := range expectedLines {
		if !strings.Contains(result, expected) {
			t.Errorf("Render() result missing expected line: %s\nGot:\n%s", expected, result)
		}
	}
}

func TestTemplateRef_RenderError(t *testing.T) {
	// Create a template that expects different fields
	tmplContent := `{{.NonExistentField}}`

	tmpl, err := template.New("error.tmpl").Parse(tmplContent)
	if err != nil {
		t.Fatalf("failed to parse test template: %v", err)
	}

	ref := TemplateRef{
		Name:     "error.tmpl",
		Template: tmpl,
	}

	ctx := TemplateContext{
		BaseDir: "spectr",
	}

	_, err = ref.Render(ctx)
	if err == nil {
		t.Error("Render() should have returned an error for non-existent field")
	}
}

func TestDefaultTemplateContext(t *testing.T) {
	ctx := DefaultTemplateContext()

	tests := []struct {
		name     string
		got      string
		expected string
	}{
		{"BaseDir", ctx.BaseDir, "spectr"},
		{"SpecsDir", ctx.SpecsDir, "spectr/specs"},
		{"ChangesDir", ctx.ChangesDir, "spectr/changes"},
		{"ProjectFile", ctx.ProjectFile, "spectr/project.md"},
		{"AgentsFile", ctx.AgentsFile, "spectr/AGENTS.md"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("DefaultTemplateContext().%s = %q, want %q", tt.name, tt.got, tt.expected)
			}
		})
	}
}
