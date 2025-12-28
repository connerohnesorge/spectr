package initialize

import (
	"strings"
	"testing"

	"github.com/connerohnesorge/spectr/internal/domain"
)

// testTemplateContext returns a default TemplateContext for testing.
// This is similar to testTemplateContext() but for test files.
func testTemplateContext() *domain.TemplateContext {
	return &domain.TemplateContext{
		BaseDir:     "spectr",
		SpecsDir:    "spectr/specs",
		ChangesDir:  "spectr/changes",
		ProjectFile: "spectr/project.md",
		AgentsFile:  "spectr/AGENTS.md",
	}
}

func TestNewTemplateManager(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf(
			"NewTemplateManager() error = %v",
			err,
		)
	}
	if tm == nil {
		t.Fatal(
			"NewTemplateManager() returned nil",
		)
	}
	if tm.templates == nil {
		t.Fatal(
			"TemplateManager.templates is nil",
		)
	}
}

//
//nolint:revive // cognitive-complexity - comprehensive test coverage
func TestTemplateManager_RenderProject(
	t *testing.T,
) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf(
			"NewTemplateManager() error = %v",
			err,
		)
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
				TechStack: []string{
					"Go",
					"PostgreSQL",
				},
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
				TechStack: []string{
					"TypeScript",
				},
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
				t.Errorf(
					"RenderProject() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)

				return
			}
			if err != nil {
				return
			}

			// Check that all expected strings are in the output
			for _, want := range tt.want {
				if !strings.Contains(got, want) {
					t.Errorf(
						"RenderProject() missing expected string %q in output:\n%s",
						want,
						got,
					)
				}
			}

			// Verify basic structure
			if !strings.Contains(
				got,
				"## Tech Stack",
			) {
				t.Error(
					"RenderProject() missing '## Tech Stack' section",
				)
			}
		})
	}
}

func TestTemplateManager_RenderAgents(
	t *testing.T,
) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf(
			"NewTemplateManager() error = %v",
			err,
		)
	}

	got, err := tm.RenderAgents(
		testTemplateContext(),
	)
	if err != nil {
		t.Fatalf("RenderAgents() error = %v", err)
	}

	// Check for key sections in AGENTS.md
	expectedSections := []string{
		"# Spectr Instructions",
		"## TL;DR Quick Checklist",
		"## Two-Stage Workflow",
		"### Stage 1: Creating Changes",
		"### Stage 2: Implementing Changes",
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
			t.Errorf(
				"RenderAgents() missing expected section %q",
				section,
			)
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

func TestTemplateManager_RenderInstructionPointer(
	t *testing.T,
) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf(
			"NewTemplateManager() error = %v",
			err,
		)
	}

	got, err := tm.RenderInstructionPointer(
		testTemplateContext(),
	)
	if err != nil {
		t.Fatalf(
			"RenderInstructionPointer() error = %v",
			err,
		)
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
			t.Errorf(
				"RenderInstructionPointer() missing expected content %q",
				content,
			)
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

func TestTemplateManager_RenderSlashCommand(
	t *testing.T,
) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf(
			"NewTemplateManager() error = %v",
			err,
		)
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
				"# Guardrails",
				"# Steps",
				"# Reference",
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
				"# Guardrails",
				"# Steps",
				"# Reference",
				"changes/<id>/",
				"tasks.json",
				"spectr accept",
				"pending",
				"completed",
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
			got, err := tm.RenderSlashCommand(
				tt.commandType,
				testTemplateContext(),
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"RenderSlashCommand() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)

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

func TestTemplateManager_RenderCIWorkflow(
	t *testing.T,
) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf(
			"NewTemplateManager() error = %v",
			err,
		)
	}

	got, err := tm.RenderCIWorkflow()
	if err != nil {
		t.Fatalf(
			"RenderCIWorkflow() error = %v",
			err,
		)
	}

	// Check for key content in CI workflow
	expectedContent := []string{
		"name: Spectr Validation",
		"push:",
		"pull_request:",
		"branches: [main]",
		"spectr-validate:",
		"runs-on: ubuntu-latest",
		"actions/checkout@v4",
		"connerohnesorge/spectr-action@v0.0.2",
		"strict: false",
	}

	for _, content := range expectedContent {
		if !strings.Contains(got, content) {
			t.Errorf(
				"RenderCIWorkflow() missing expected content %q",
				content,
			)
		}
	}
}

func TestTemplateManager_AllTemplatesCompile(
	t *testing.T,
) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf(
			"NewTemplateManager() error = %v",
			err,
		)
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
			t.Errorf(
				"Project template failed to render: %v",
				err,
			)
		}
	})

	t.Run("agents template", func(t *testing.T) {
		_, err := tm.RenderAgents(
			testTemplateContext(),
		)
		if err != nil {
			t.Errorf(
				"Agents template failed to render: %v",
				err,
			)
		}
	})

	t.Run(
		"instruction pointer template",
		func(t *testing.T) {
			_, err := tm.RenderInstructionPointer(
				testTemplateContext(),
			)
			if err != nil {
				t.Errorf(
					"Instruction pointer template failed to render: %v",
					err,
				)
			}
		},
	)

	t.Run("slash commands", func(t *testing.T) {
		commands := []string{"proposal", "apply"}
		for _, cmd := range commands {
			_, err := tm.RenderSlashCommand(
				cmd,
				testTemplateContext(),
			)
			if err != nil {
				t.Errorf(
					"Slash command %s failed to render: %v",
					cmd,
					err,
				)
			}
		}
	})

	t.Run(
		"ci workflow template",
		func(t *testing.T) {
			_, err := tm.RenderCIWorkflow()
			if err != nil {
				t.Errorf(
					"CI workflow template failed to render: %v",
					err,
				)
			}
		},
	)
}

//
//nolint:revive // cognitive-complexity - comprehensive test coverage
func TestTemplateManager_VariableSubstitution(
	t *testing.T,
) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf(
			"NewTemplateManager() error = %v",
			err,
		)
	}

	t.Run(
		"project variables are substituted",
		func(t *testing.T) {
			ctx := ProjectContext{
				ProjectName: "VariableTest",
				Description: "Testing variable substitution",
				TechStack: []string{
					"Go",
					"React",
					"PostgreSQL",
				},
			}
			got, err := tm.RenderProject(ctx)
			if err != nil {
				t.Fatalf(
					"RenderProject() error = %v",
					err,
				)
			}

			// Verify no template syntax remains
			if strings.Contains(got, "{{") ||
				strings.Contains(got, "}}") {
				t.Error(
					"Template contains unreplaced template syntax",
				)
			}

			// Verify all variables were substituted
			if !strings.Contains(
				got,
				"VariableTest",
			) {
				t.Error(
					"ProjectName not substituted",
				)
			}
			if !strings.Contains(
				got,
				"Testing variable substitution",
			) {
				t.Error(
					"Description not substituted",
				)
			}
			if !strings.Contains(got, "Go") ||
				!strings.Contains(got, "React") ||
				!strings.Contains(
					got,
					"PostgreSQL",
				) {
				t.Error(
					"TechStack items not substituted",
				)
			}
		},
	)
}

func TestTemplateManager_EmptyTechStack(
	t *testing.T,
) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf(
			"NewTemplateManager() error = %v",
			err,
		)
	}

	// Test with nil tech stack
	ctx := ProjectContext{
		ProjectName: "NilStack",
		Description: "Test nil tech stack",
		TechStack:   nil,
	}
	got, err := tm.RenderProject(ctx)
	if err != nil {
		t.Fatalf(
			"RenderProject() with nil TechStack error = %v",
			err,
		)
	}

	// Should still have the Tech Stack section
	if !strings.Contains(got, "## Tech Stack") {
		t.Error(
			"Missing Tech Stack section with nil slice",
		)
	}

	// Test with empty slice
	ctx.TechStack = make([]string, 0)
	got, err = tm.RenderProject(ctx)
	if err != nil {
		t.Fatalf(
			"RenderProject() with empty TechStack error = %v",
			err,
		)
	}

	if !strings.Contains(got, "## Tech Stack") {
		t.Error(
			"Missing Tech Stack section with empty slice",
		)
	}
}

// TestTemplateManager_InstructionPointer tests the InstructionPointer accessor method.
func TestTemplateManager_InstructionPointer(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	ref := tm.InstructionPointer()

	// Check that the TemplateRef is properly constructed
	if ref.Name != "instruction-pointer.md.tmpl" {
		t.Errorf(
			"InstructionPointer().Name = %q, want %q",
			ref.Name,
			"instruction-pointer.md.tmpl",
		)
	}
	if ref.Template == nil {
		t.Error("InstructionPointer().Template is nil")
	}

	// Test that Render() works with the returned TemplateRef
	got, err := tm.Render(ref, testTemplateContext())
	if err != nil {
		t.Fatalf("Render(InstructionPointer()) error = %v", err)
	}

	// Verify content
	if !strings.Contains(got, "Spectr") {
		t.Error("Rendered instruction pointer missing 'Spectr' content")
	}
}

// TestTemplateManager_Agents tests the Agents accessor method.
func TestTemplateManager_Agents(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	ref := tm.Agents()

	// Check that the TemplateRef is properly constructed
	if ref.Name != "AGENTS.md.tmpl" {
		t.Errorf(
			"Agents().Name = %q, want %q",
			ref.Name,
			"AGENTS.md.tmpl",
		)
	}
	if ref.Template == nil {
		t.Error("Agents().Template is nil")
	}

	// Test that Render() works with the returned TemplateRef
	got, err := tm.Render(ref, testTemplateContext())
	if err != nil {
		t.Fatalf("Render(Agents()) error = %v", err)
	}

	// Verify content
	if !strings.Contains(got, "Spectr Instructions") {
		t.Error("Rendered agents template missing 'Spectr Instructions' content")
	}
}

// TestTemplateManager_SlashCommand tests the SlashCommand accessor method.
func TestTemplateManager_SlashCommand(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	tests := []struct {
		name         string
		cmd          domain.SlashCommand
		expectedName string
		expectedText string
	}{
		{
			name:         "proposal command",
			cmd:          domain.SlashProposal,
			expectedName: "slash-proposal.md.tmpl",
			expectedText: "Guardrails",
		},
		{
			name:         "apply command",
			cmd:          domain.SlashApply,
			expectedName: "slash-apply.md.tmpl",
			expectedText: "Guardrails",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref := tm.SlashCommand(tt.cmd)

			// Check that the TemplateRef is properly constructed
			if ref.Name != tt.expectedName {
				t.Errorf(
					"SlashCommand(%v).Name = %q, want %q",
					tt.cmd,
					ref.Name,
					tt.expectedName,
				)
			}
			if ref.Template == nil {
				t.Errorf("SlashCommand(%v).Template is nil", tt.cmd)
			}

			// Test that Render() works with the returned TemplateRef
			got, err := tm.Render(ref, testTemplateContext())
			if err != nil {
				t.Fatalf("Render(SlashCommand(%v)) error = %v", tt.cmd, err)
			}

			// Verify content
			if !strings.Contains(got, tt.expectedText) {
				t.Errorf(
					"Rendered %s missing %q content",
					tt.name,
					tt.expectedText,
				)
			}
		})
	}
}

// TestTemplateManager_TOMLSlashCommand tests the TOMLSlashCommand accessor method.
func TestTemplateManager_TOMLSlashCommand(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	tests := []struct {
		name         string
		cmd          domain.SlashCommand
		expectedName string
		expectedText string
	}{
		{
			name:         "proposal TOML command",
			cmd:          domain.SlashProposal,
			expectedName: "slash-proposal.toml.tmpl",
			expectedText: "description",
		},
		{
			name:         "apply TOML command",
			cmd:          domain.SlashApply,
			expectedName: "slash-apply.toml.tmpl",
			expectedText: "description",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref := tm.TOMLSlashCommand(tt.cmd)

			// Check that the TemplateRef is properly constructed
			if ref.Name != tt.expectedName {
				t.Errorf(
					"TOMLSlashCommand(%v).Name = %q, want %q",
					tt.cmd,
					ref.Name,
					tt.expectedName,
				)
			}
			if ref.Template == nil {
				t.Errorf("TOMLSlashCommand(%v).Template is nil", tt.cmd)
			}

			// Test that Render() works with the returned TemplateRef
			got, err := tm.Render(ref, testTemplateContext())
			if err != nil {
				t.Fatalf("Render(TOMLSlashCommand(%v)) error = %v", tt.cmd, err)
			}

			// Verify TOML content
			if !strings.Contains(got, tt.expectedText) {
				t.Errorf(
					"Rendered %s missing %q content",
					tt.name,
					tt.expectedText,
				)
			}
			// TOML files should have prompt field
			if !strings.Contains(got, "prompt") {
				t.Errorf(
					"Rendered %s missing 'prompt' field",
					tt.name,
				)
			}
		})
	}
}

// TestTemplateManager_Render tests the Render method with various TemplateRefs.
func TestTemplateManager_Render(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	ctx := testTemplateContext()

	// Test rendering each accessor's template
	t.Run("render instruction pointer", func(t *testing.T) {
		ref := tm.InstructionPointer()
		got, err := tm.Render(ref, ctx)
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}
		if got == "" {
			t.Error("Render() returned empty string")
		}
	})

	t.Run("render agents", func(t *testing.T) {
		ref := tm.Agents()
		got, err := tm.Render(ref, ctx)
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}
		if len(got) < 1000 {
			t.Error("Render(Agents) returned too short content")
		}
	})

	t.Run("render slash commands", func(t *testing.T) {
		for _, cmd := range []domain.SlashCommand{domain.SlashProposal, domain.SlashApply} {
			ref := tm.SlashCommand(cmd)
			got, err := tm.Render(ref, ctx)
			if err != nil {
				t.Fatalf("Render(SlashCommand(%v)) error = %v", cmd, err)
			}
			if got == "" {
				t.Errorf("Render(SlashCommand(%v)) returned empty string", cmd)
			}
		}
	})

	t.Run("render TOML slash commands", func(t *testing.T) {
		for _, cmd := range []domain.SlashCommand{domain.SlashProposal, domain.SlashApply} {
			ref := tm.TOMLSlashCommand(cmd)
			got, err := tm.Render(ref, ctx)
			if err != nil {
				t.Fatalf("Render(TOMLSlashCommand(%v)) error = %v", cmd, err)
			}
			if got == "" {
				t.Errorf("Render(TOMLSlashCommand(%v)) returned empty string", cmd)
			}
		}
	})
}

// TestTemplateManager_MergesDomainTemplates tests that templates from domain package are merged.
func TestTemplateManager_MergesDomainTemplates(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	// Verify slash command templates from domain are available
	domainTemplates := []string{
		"slash-proposal.md.tmpl",
		"slash-apply.md.tmpl",
		"slash-proposal.toml.tmpl",
		"slash-apply.toml.tmpl",
	}

	for _, name := range domainTemplates {
		t.Run("template_"+name, func(t *testing.T) {
			// Look up template by name
			tmpl := tm.templates.Lookup(name)
			if tmpl == nil {
				t.Errorf("Template %q not found in merged template set", name)
			}
		})
	}

	// Verify main templates are still available
	mainTemplates := []string{
		"AGENTS.md.tmpl",
		"instruction-pointer.md.tmpl",
		"project.md.tmpl",
	}

	for _, name := range mainTemplates {
		t.Run("template_"+name, func(t *testing.T) {
			tmpl := tm.templates.Lookup(name)
			if tmpl == nil {
				t.Errorf("Template %q not found in merged template set", name)
			}
		})
	}
}
