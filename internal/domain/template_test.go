package domain

import (
	"html/template"
	"testing"
)

func TestTemplateRef_FieldAccess(t *testing.T) {
	tmpl := template.New("test")

	ref := TemplateRef{
		Name:     "test-template.md.tmpl",
		Template: tmpl,
	}

	if ref.Name != "test-template.md.tmpl" {
		t.Errorf("expected Name to be 'test-template.md.tmpl', got %q", ref.Name)
	}

	if ref.Template != tmpl {
		t.Error("expected Template to be the assigned template")
	}
}

func TestTemplateRef_ZeroValue(t *testing.T) {
	var ref TemplateRef

	if ref.Name != "" {
		t.Errorf("expected zero value Name to be empty, got %q", ref.Name)
	}

	if ref.Template != nil {
		t.Error("expected zero value Template to be nil")
	}
}

func TestTemplateContext_FieldAccess(t *testing.T) {
	ctx := TemplateContext{
		BaseDir:     "spectr",
		SpecsDir:    "spectr/specs",
		ChangesDir:  "spectr/changes",
		ProjectFile: "spectr/project.md",
		AgentsFile:  "spectr/AGENTS.md",
	}

	tests := []struct {
		field    string
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
		if tt.got != tt.expected {
			t.Errorf("%s: expected %q, got %q", tt.field, tt.expected, tt.got)
		}
	}
}

func TestTemplateContext_ZeroValue(t *testing.T) {
	var ctx TemplateContext

	if ctx.BaseDir != "" {
		t.Errorf("expected zero value BaseDir to be empty, got %q", ctx.BaseDir)
	}
	if ctx.SpecsDir != "" {
		t.Errorf("expected zero value SpecsDir to be empty, got %q", ctx.SpecsDir)
	}
	if ctx.ChangesDir != "" {
		t.Errorf("expected zero value ChangesDir to be empty, got %q", ctx.ChangesDir)
	}
	if ctx.ProjectFile != "" {
		t.Errorf("expected zero value ProjectFile to be empty, got %q", ctx.ProjectFile)
	}
	if ctx.AgentsFile != "" {
		t.Errorf("expected zero value AgentsFile to be empty, got %q", ctx.AgentsFile)
	}
}
