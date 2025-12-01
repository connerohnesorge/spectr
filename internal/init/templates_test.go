package init

import (
	"strings"
	"testing"
)

func TestNewTemplateManager(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}
	if tm == nil {
		t.Fatal("NewTemplateManager() returned nil")
	}
	if tm.templates == nil {
		t.Fatal("TemplateManager.templates is nil")
	}
}

//nolint:revive // cognitive-complexity - comprehensive test coverage
func TestTemplateManager_RenderProject(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	tests := []struct {
		name    string
		ctx     ProjectContext
		want    []string // Strings that should be in the output
		wantErr bool
	}{
		{
			name: "basic project",
			ctx: ProjectContext{
				ProjectName: "MyProject",
				Description: "A test project",
				TechStack:   []string{"Go", "PostgreSQL"},
			},
			want: []string{
				"# MyProject Context",
				"A test project",
				"- Go",
				"- PostgreSQL",
				"## Project Conventions",
			},
			wantErr: false,
		},
		{
			name: "empty tech stack",
			ctx: ProjectContext{
				ProjectName: "EmptyStack",
				Description: "No tech stack",
				TechStack:   make([]string, 0),
			},
			want: []string{
				"# EmptyStack Context",
				"No tech stack",
			},
			wantErr: false,
		},
		{
			name: "single tech",
			ctx: ProjectContext{
				ProjectName: "SingleTech",
				Description: "One technology",
				TechStack:   []string{"TypeScript"},
			},
			want: []string{
				"# SingleTech Context",
				"- TypeScript",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tm.RenderProject(tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("RenderProject() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if err != nil {
				return
			}

			// Check that all expected strings are in the output
			for _, want := range tt.want {
				if !strings.Contains(got, want) {
					t.Errorf("RenderProject() missing expected string %q in output:\n%s", want, got)
				}
			}

			// Verify basic structure
			if !strings.Contains(got, "## Purpose") {
				t.Error("RenderProject() missing '## Purpose' section")
			}
			if !strings.Contains(got, "## Tech Stack") {
				t.Error("RenderProject() missing '## Tech Stack' section")
			}
		})
	}
}

func TestTemplateManager_RenderAgents(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	got, err := tm.RenderAgents()
	if err != nil {
		t.Fatalf("RenderAgents() error = %v", err)
	}

	// Check for key sections in AGENTS.md
	expectedSections := []string{
		"# Spectr Instructions",
		"## TL;DR Quick Checklist",
		"## Three-Stage Workflow",
		"### Stage 1: Creating Changes",
		"### Stage 2: Implementing Changes",
		"### Stage 3: Syncing",
		"## Directory Structure",
		"## Creating Change Proposals",
		"## Spec File Format",
		"#### Scenario:",
		"spectr validate",
		"spectr list",
		"## ADDED Requirements",
		"## MODIFIED Requirements",
		"## REMOVED Requirements",
	}

	for _, section := range expectedSections {
		if !strings.Contains(got, section) {
			t.Errorf("RenderAgents() missing expected section %q", section)
		}
	}

	// Verify it's a substantial document (should be thousands of characters)
	if len(got) < 5000 {
		t.Errorf(
			"RenderAgents() output too short: got %d characters, expected at least 5000",
			len(got),
		)
	}
}

func TestTemplateManager_RenderInstructionPointer(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	got, err := tm.RenderInstructionPointer()
	if err != nil {
		t.Fatalf("RenderInstructionPointer() error = %v", err)
	}

	// Check for key content in instruction pointer
	expectedContent := []string{
		"# Spectr Instructions",
		"spectr/AGENTS.md",
		"proposal",
		"spec",
		"change",
	}

	for _, content := range expectedContent {
		if !strings.Contains(got, content) {
			t.Errorf("RenderInstructionPointer() missing expected content %q", content)
		}
	}

	// Verify it's a concise pointer (less than 20 lines as per spec)
	lineCount := strings.Count(got, "\n") + 1
	if lineCount > 20 {
		t.Errorf(
			"RenderInstructionPointer() output too long: got %d lines, expected at most 20",
			lineCount,
		)
	}

	// Verify it does NOT contain the full workflow instructions
	fullWorkflowIndicators := []string{
		"## TL;DR Quick Checklist",
		"## Three-Stage Workflow",
		"## Directory Structure",
	}

	for _, indicator := range fullWorkflowIndicators {
		if strings.Contains(got, indicator) {
			t.Errorf(
				"RenderInstructionPointer() should not contain full workflow content %q",
				indicator,
			)
		}
	}
}

func TestTemplateManager_RenderSlashCommand(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	tests := []struct {
		name        string
		commandType string
		want        []string
		wantErr     bool
	}{
		{
			name:        "proposal command",
			commandType: "proposal",
			want: []string{
				"**Guardrails**",
				"**Steps**",
				"**Reference**",
				"spectr validate",
				"change-id",
				"proposal.md",
				"tasks.md",
				"design.md",
			},
			wantErr: false,
		},
		{
			name:        "apply command",
			commandType: "apply",
			want: []string{
				"**Guardrails**",
				"**Steps**",
				"**Reference**",
				"changes/<id>/",
				"tasks.md",
				"- [x]",
			},
			wantErr: false,
		},
		{
			name:        "sync command",
			commandType: "sync",
			want: []string{
				"### Guardrails",
				"### Steps",
				"### Reference",
				"Code is the source of truth",
				"spectr list --specs",
				"spectr validate --strict",
			},
			wantErr: false,
		},
		{
			name:        "invalid command type",
			commandType: "invalid",
			want:        nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tm.RenderSlashCommand(tt.commandType)
			if (err != nil) != tt.wantErr {
				t.Errorf("RenderSlashCommand() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if err != nil {
				return
			}

			for _, want := range tt.want {
				if !strings.Contains(got, want) {
					t.Errorf(
						"RenderSlashCommand() missing expected string %q in output:\n%s",
						want,
						got,
					)
				}
			}
		})
	}
}

func TestTemplateManager_AllTemplatesCompile(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	// Test that all templates can be rendered without errors
	t.Run("project template", func(t *testing.T) {
		ctx := ProjectContext{
			ProjectName: "Test",
			Description: "Test Description",
			TechStack:   []string{"Go"},
		}
		_, err := tm.RenderProject(ctx)
		if err != nil {
			t.Errorf("Project template failed to render: %v", err)
		}
	})

	t.Run("agents template", func(t *testing.T) {
		_, err := tm.RenderAgents()
		if err != nil {
			t.Errorf("Agents template failed to render: %v", err)
		}
	})

	t.Run("instruction pointer template", func(t *testing.T) {
		_, err := tm.RenderInstructionPointer()
		if err != nil {
			t.Errorf("Instruction pointer template failed to render: %v", err)
		}
	})

	t.Run("slash commands", func(t *testing.T) {
		commands := []string{"proposal", "apply", "sync"}
		for _, cmd := range commands {
			_, err := tm.RenderSlashCommand(cmd)
			if err != nil {
				t.Errorf("Slash command %s failed to render: %v", cmd, err)
			}
		}
	})
}

//nolint:revive // cognitive-complexity - comprehensive test coverage
func TestTemplateManager_VariableSubstitution(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	t.Run("project variables are substituted", func(t *testing.T) {
		ctx := ProjectContext{
			ProjectName: "VariableTest",
			Description: "Testing variable substitution",
			TechStack:   []string{"Go", "React", "PostgreSQL"},
		}
		got, err := tm.RenderProject(ctx)
		if err != nil {
			t.Fatalf("RenderProject() error = %v", err)
		}

		// Verify no template syntax remains
		if strings.Contains(got, "{{") || strings.Contains(got, "}}") {
			t.Error("Template contains unreplaced template syntax")
		}

		// Verify all variables were substituted
		if !strings.Contains(got, "VariableTest") {
			t.Error("ProjectName not substituted")
		}
		if !strings.Contains(got, "Testing variable substitution") {
			t.Error("Description not substituted")
		}
		if !strings.Contains(got, "Go") || !strings.Contains(got, "React") ||
			!strings.Contains(got, "PostgreSQL") {
			t.Error("TechStack items not substituted")
		}
	})
}

func TestTemplateManager_EmptyTechStack(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	// Test with nil tech stack
	ctx := ProjectContext{
		ProjectName: "NilStack",
		Description: "Test nil tech stack",
		TechStack:   nil,
	}
	got, err := tm.RenderProject(ctx)
	if err != nil {
		t.Fatalf("RenderProject() with nil TechStack error = %v", err)
	}

	// Should still have the Tech Stack section
	if !strings.Contains(got, "## Tech Stack") {
		t.Error("Missing Tech Stack section with nil slice")
	}

	// Test with empty slice
	ctx.TechStack = make([]string, 0)
	got, err = tm.RenderProject(ctx)
	if err != nil {
		t.Fatalf("RenderProject() with empty TechStack error = %v", err)
	}

	if !strings.Contains(got, "## Tech Stack") {
		t.Error("Missing Tech Stack section with empty slice")
	}
}
