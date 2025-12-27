package domain

import (
	"strings"
	"testing"
	"text/template"
)

func TestTemplateRef_Render(t *testing.T) {
	tests := []struct {
		name         string
		templateText string
		ctx          TemplateContext
		want         string
		wantErr      bool
	}{
		{
			name:         "renders basic template",
			templateText: "Base: {{ .BaseDir }}",
			ctx: TemplateContext{
				BaseDir: "spectr",
			},
			want:    "Base: spectr",
			wantErr: false,
		},
		{
			name:         "renders all fields",
			templateText: "{{ .BaseDir }}/{{ .SpecsDir }}/{{ .ChangesDir }}/{{ .ProjectFile }}/{{ .AgentsFile }}",
			ctx: TemplateContext{
				BaseDir:     "spectr",
				SpecsDir:    "spectr/specs",
				ChangesDir:  "spectr/changes",
				ProjectFile: "spectr/project.md",
				AgentsFile:  "spectr/AGENTS.md",
			},
			want:    "spectr/spectr/specs/spectr/changes/spectr/project.md/spectr/AGENTS.md",
			wantErr: false,
		},
		{
			name:         "handles empty values",
			templateText: "{{ .BaseDir }}",
			ctx: TemplateContext{
				BaseDir: "",
			},
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := template.New("test").Parse(tt.templateText)
			if err != nil {
				t.Fatalf("failed to parse template: %v", err)
			}

			tr := TemplateRef{
				Name:     "test",
				Template: tmpl,
			}

			got, err := tr.Render(tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("TemplateRef.Render() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if got != tt.want {
				t.Errorf("TemplateRef.Render() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultTemplateContext(t *testing.T) {
	ctx := DefaultTemplateContext()

	tests := []struct {
		name  string
		field string
		want  string
	}{
		{"BaseDir", ctx.BaseDir, "spectr"},
		{"SpecsDir", ctx.SpecsDir, "spectr/specs"},
		{"ChangesDir", ctx.ChangesDir, "spectr/changes"},
		{"ProjectFile", ctx.ProjectFile, "spectr/project.md"},
		{"AgentsFile", ctx.AgentsFile, "spectr/AGENTS.md"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.field != tt.want {
				t.Errorf("DefaultTemplateContext().%s = %v, want %v", tt.name, tt.field, tt.want)
			}
		})
	}
}

func TestTemplateRef_RenderError(t *testing.T) {
	// Test that invalid template execution returns an error
	tmpl, err := template.New("test").Parse("{{ .NonExistentField }}")
	if err != nil {
		t.Fatalf("failed to parse template: %v", err)
	}

	tr := TemplateRef{
		Name:     "test",
		Template: tmpl,
	}

	ctx := DefaultTemplateContext()
	_, err = tr.Render(ctx)
	if err == nil {
		t.Error("TemplateRef.Render() expected error for non-existent field, got nil")
	}
	if !strings.Contains(err.Error(), "failed to render template") {
		t.Errorf(
			"TemplateRef.Render() error message = %v, want to contain 'failed to render template'",
			err,
		)
	}
}
