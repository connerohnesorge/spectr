package initialize

import (
	"strings"
	"testing"
	"testing/fstest"
	"text/template"

	"github.com/connerohnesorge/spectr/internal/domain"
)

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

	ctx := domain.DefaultTemplateContext()
	got, err := tm.RenderAgents(
		&ctx,
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

	ctx := domain.DefaultTemplateContext()
	got, err := tm.RenderInstructionPointer(
		&ctx,
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
			ctx := domain.DefaultTemplateContext()
			got, err := tm.RenderSlashCommand(
				tt.commandType,
				&ctx,
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
		ctx := domain.DefaultTemplateContext()
		_, err := tm.RenderAgents(
			&ctx,
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
			ctx := domain.DefaultTemplateContext()
			_, err := tm.RenderInstructionPointer(
				&ctx,
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
			ctx := domain.DefaultTemplateContext()
			_, err := tm.RenderSlashCommand(
				cmd,
				&ctx,
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

func TestTemplateManager_InstructionPointer(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	// Get the TemplateRef
	ref := tm.InstructionPointer()

	// Verify TemplateRef fields
	if ref.Name != "instruction-pointer.md.tmpl" {
		t.Errorf("InstructionPointer().Name = %q, want %q", ref.Name, "instruction-pointer.md.tmpl")
	}
	if ref.Template == nil {
		t.Error("InstructionPointer().Template is nil")
	}

	// Verify Render() works
	ctx := domain.DefaultTemplateContext()
	rendered, err := ref.Render(&ctx)
	if err != nil {
		t.Fatalf("InstructionPointer().Render() error = %v", err)
	}

	// Check for expected content
	expectedContent := []string{
		"# Spectr Instructions",
		"spectr/AGENTS.md",
	}
	for _, content := range expectedContent {
		if !strings.Contains(rendered, content) {
			t.Errorf("InstructionPointer().Render() missing expected content %q", content)
		}
	}
}

func TestTemplateManager_Agents(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	// Get the TemplateRef
	ref := tm.Agents()

	// Verify TemplateRef fields
	if ref.Name != "AGENTS.md.tmpl" {
		t.Errorf("Agents().Name = %q, want %q", ref.Name, "AGENTS.md.tmpl")
	}
	if ref.Template == nil {
		t.Error("Agents().Template is nil")
	}

	// Verify Render() works
	ctx := domain.DefaultTemplateContext()
	rendered, err := ref.Render(&ctx)
	if err != nil {
		t.Fatalf("Agents().Render() error = %v", err)
	}

	// Check for expected content
	expectedContent := []string{
		"# Spectr Instructions",
		"## TL;DR Quick Checklist",
	}
	for _, content := range expectedContent {
		if !strings.Contains(rendered, content) {
			t.Errorf("Agents().Render() missing expected content %q", content)
		}
	}
}

func TestTemplateManager_SlashCommand(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	tests := []struct {
		name        string
		cmd         domain.SlashCommand
		wantName    string
		wantContent []string
	}{
		{
			name:     "proposal command",
			cmd:      domain.SlashProposal,
			wantName: "slash-proposal.md.tmpl",
			wantContent: []string{
				"# Guardrails",
				"proposal.md",
			},
		},
		{
			name:     "apply command",
			cmd:      domain.SlashApply,
			wantName: "slash-apply.md.tmpl",
			wantContent: []string{
				"# Guardrails",
				"tasks.json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get the TemplateRef
			ref := tm.SlashCommand(tt.cmd)

			// Verify TemplateRef fields
			if ref.Name != tt.wantName {
				t.Errorf("SlashCommand(%v).Name = %q, want %q", tt.cmd, ref.Name, tt.wantName)
			}
			if ref.Template == nil {
				t.Error("SlashCommand().Template is nil")
			}

			// Verify Render() works
			ctx := domain.DefaultTemplateContext()
			rendered, err := ref.Render(&ctx)
			if err != nil {
				t.Fatalf("SlashCommand(%v).Render() error = %v", tt.cmd, err)
			}

			// Check for expected content
			for _, content := range tt.wantContent {
				if !strings.Contains(rendered, content) {
					t.Errorf(
						"SlashCommand(%v).Render() missing expected content %q",
						tt.cmd,
						content,
					)
				}
			}
		})
	}
}

func TestTemplateManager_TOMLSlashCommand(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	tests := []struct {
		name        string
		cmd         domain.SlashCommand
		wantName    string
		wantContent []string
	}{
		{
			name:     "TOML proposal command",
			cmd:      domain.SlashProposal,
			wantName: "slash-proposal.toml.tmpl",
			wantContent: []string{
				"description =",
				"prompt =",
			},
		},
		{
			name:     "TOML apply command",
			cmd:      domain.SlashApply,
			wantName: "slash-apply.toml.tmpl",
			wantContent: []string{
				"description =",
				"prompt =",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get the TemplateRef
			ref := tm.TOMLSlashCommand(tt.cmd)

			// Verify TemplateRef fields
			if ref.Name != tt.wantName {
				t.Errorf("TOMLSlashCommand(%v).Name = %q, want %q", tt.cmd, ref.Name, tt.wantName)
			}
			if ref.Template == nil {
				t.Error("TOMLSlashCommand().Template is nil")
			}

			// Verify Render() works
			ctx := domain.DefaultTemplateContext()
			rendered, err := ref.Render(&ctx)
			if err != nil {
				t.Fatalf("TOMLSlashCommand(%v).Render() error = %v", tt.cmd, err)
			}

			// Check for expected content
			for _, content := range tt.wantContent {
				if !strings.Contains(rendered, content) {
					t.Errorf(
						"TOMLSlashCommand(%v).Render() missing expected content %q",
						tt.cmd,
						content,
					)
				}
			}
		})
	}
}

func TestTemplateManager_ParseProviderTemplates(t *testing.T) {
	baseTemplates := map[string]*template.Template{
		"slash-proposal.md.tmpl": template.New("slash-proposal.md.tmpl"),
	}
	providerFS := fstest.MapFS{
		"templates/providers/claude-code/slash-proposal.md.tmpl": &fstest.MapFile{
			Data: []byte(`{{define "guardrails"}}claude{{end}}`),
		},
		"templates/providers/codex/slash-proposal.md.tmpl": &fstest.MapFile{
			Data: []byte(`{{define "guardrails"}}codex{{end}}`),
		},
	}

	providerTemplates, err := loadProviderTemplates(baseTemplates, providerFS)
	if err != nil {
		t.Fatalf("loadProviderTemplates() error = %v", err)
	}
	if providerTemplates["claude-code"]["slash-proposal.md.tmpl"] == nil {
		t.Error("expected claude-code proposal template override")
	}
	if providerTemplates["codex"]["slash-proposal.md.tmpl"] == nil {
		t.Error("expected codex proposal template override")
	}
}

func TestTemplateManager_MissingProvidersDirectory(t *testing.T) {
	providerTemplates, err := loadProviderTemplates(
		make(map[string]*template.Template),
		fstest.MapFS{},
	)
	if err != nil {
		t.Fatalf("loadProviderTemplates() error = %v", err)
	}
	if len(providerTemplates) != 0 {
		t.Errorf("loadProviderTemplates() expected empty map, got %d", len(providerTemplates))
	}
}

func TestTemplateManager_ProviderSlashCommand(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	ref := tm.ProviderSlashCommand("claude-code", domain.SlashProposal)
	if ref.ProviderTemplate == nil {
		t.Fatal("ProviderSlashCommand() expected provider template for claude-code")
	}

	ctx := domain.DefaultTemplateContext()
	rendered, err := ref.Render(&ctx)
	if err != nil {
		t.Fatalf("ProviderSlashCommand().Render() error = %v", err)
	}
	if !strings.Contains(rendered, "orchestrator pattern") {
		t.Errorf(
			"ProviderSlashCommand() missing provider-specific content in output:\n%s",
			rendered,
		)
	}
}

func TestTemplateManager_ProviderSlashCommand_Unknown(t *testing.T) {
	tm := &TemplateManager{
		slashTemplates: map[string]*template.Template{
			"slash-proposal.md.tmpl": template.New("slash-proposal.md.tmpl"),
		},
		providerTemplates: make(map[string]map[string]*template.Template),
	}

	ref := tm.ProviderSlashCommand("unknown-provider", domain.SlashProposal)
	if ref.ProviderTemplate != nil {
		t.Error("ProviderSlashCommand() expected nil ProviderTemplate for unknown provider")
	}
}

func TestTemplateManager_BackwardCompatibility(t *testing.T) {
	tm := &TemplateManager{
		slashTemplates: map[string]*template.Template{
			"slash-proposal.md.tmpl": template.New("slash-proposal.md.tmpl"),
			"slash-apply.toml.tmpl":  template.New("slash-apply.toml.tmpl"),
		},
		providerTemplates: make(map[string]map[string]*template.Template),
	}

	ref := tm.SlashCommand(domain.SlashProposal)
	if ref.ProviderTemplate != nil {
		t.Error("SlashCommand() expected nil ProviderTemplate")
	}

	ref = tm.TOMLSlashCommand(domain.SlashApply)
	if ref.ProviderTemplate != nil {
		t.Error("TOMLSlashCommand() expected nil ProviderTemplate")
	}
}

func TestTemplateManager_ValidationError(t *testing.T) {
	base := template.New("slash-proposal.md.tmpl")
	providerFS := fstest.MapFS{
		"templates/providers/bad/slash-proposal.md.tmpl": &fstest.MapFile{
			Data: []byte(`{{define "unknown"}}bad{{end}}`),
		},
	}

	_, err := loadProviderTemplates(
		map[string]*template.Template{
			"slash-proposal.md.tmpl": base,
		},
		providerFS,
	)
	if err == nil {
		t.Fatal("loadProviderTemplates() expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unknown section") {
		t.Errorf("loadProviderTemplates() error = %v, want unknown section error", err)
	}
}
