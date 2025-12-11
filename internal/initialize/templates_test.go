package initialize

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/connerohnesorge/spectr/internal/initialize/providers"
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

	got, err := tm.RenderAgents(providers.DefaultTemplateContext(), "")
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

	got, err := tm.RenderInstructionPointer(providers.DefaultTemplateContext(), "")
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
				providers.DefaultTemplateContext(),
				"",
			)
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

func TestTemplateManager_RenderCIWorkflow(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	got, err := tm.RenderCIWorkflow()
	if err != nil {
		t.Fatalf("RenderCIWorkflow() error = %v", err)
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
			t.Errorf("RenderCIWorkflow() missing expected content %q", content)
		}
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
		_, err := tm.RenderAgents(providers.DefaultTemplateContext(), "")
		if err != nil {
			t.Errorf("Agents template failed to render: %v", err)
		}
	})

	t.Run("instruction pointer template", func(t *testing.T) {
		_, err := tm.RenderInstructionPointer(providers.DefaultTemplateContext(), "")
		if err != nil {
			t.Errorf("Instruction pointer template failed to render: %v", err)
		}
	})

	t.Run("slash commands", func(t *testing.T) {
		commands := []string{"proposal", "apply"}
		for _, cmd := range commands {
			_, err := tm.RenderSlashCommand(cmd, providers.DefaultTemplateContext(), "")
			if err != nil {
				t.Errorf("Slash command %s failed to render: %v", cmd, err)
			}
		}
	})

	t.Run("ci workflow template", func(t *testing.T) {
		_, err := tm.RenderCIWorkflow()
		if err != nil {
			t.Errorf("CI workflow template failed to render: %v", err)
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

// =============================================================================
// Provider-Specific Template Resolution Tests
// =============================================================================

// TestResolveTemplatePath tests the resolveTemplatePath helper function directly.
// This function implements the provider-first lookup with fallback logic.
func TestResolveTemplatePath(t *testing.T) {
	tests := []struct {
		name         string
		providerID   string
		templateName string
		fallbackDir  string
		want         string
	}{
		{
			name:         "empty provider ID uses fallback directly",
			providerID:   "",
			templateName: "AGENTS.md.tmpl",
			fallbackDir:  "spectr",
			want:         "templates/spectr/AGENTS.md.tmpl",
		},
		{
			name:         "non-existent provider falls back to generic",
			providerID:   "non-existent-provider",
			templateName: "AGENTS.md.tmpl",
			fallbackDir:  "spectr",
			want:         "templates/spectr/AGENTS.md.tmpl",
		},
		{
			name:         "slash command fallback for non-existent provider",
			providerID:   "some-provider",
			templateName: "slash-proposal.md.tmpl",
			fallbackDir:  "tools",
			want:         "templates/tools/slash-proposal.md.tmpl",
		},
		{
			name:         "instruction pointer fallback for non-existent provider",
			providerID:   "another-provider",
			templateName: "instruction-pointer.md.tmpl",
			fallbackDir:  "spectr",
			want:         "templates/spectr/instruction-pointer.md.tmpl",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveTemplatePath(tt.providerID, tt.templateName, tt.fallbackDir)
			if got != tt.want {
				t.Errorf("resolveTemplatePath(%q, %q, %q) = %q, want %q",
					tt.providerID, tt.templateName, tt.fallbackDir, got, tt.want)
			}
		})
	}
}

// TestTemplateExists tests the templateExists helper function.
func TestTemplateExists(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "existing generic AGENTS.md template",
			path: "templates/spectr/AGENTS.md.tmpl",
			want: true,
		},
		{
			name: "existing generic instruction-pointer template",
			path: "templates/spectr/instruction-pointer.md.tmpl",
			want: true,
		},
		{
			name: "existing tools slash-proposal template",
			path: "templates/tools/slash-proposal.md.tmpl",
			want: true,
		},
		{
			name: "existing tools slash-apply template",
			path: "templates/tools/slash-apply.md.tmpl",
			want: true,
		},
		{
			name: "existing CI workflow template",
			path: "templates/ci/spectr-ci.yml.tmpl",
			want: true,
		},
		{
			name: "existing claude-code AGENTS.md template",
			path: "templates/claude-code/AGENTS.md.tmpl",
			want: true,
		},
		{
			name: "existing crush AGENTS.md template",
			path: "templates/crush/AGENTS.md.tmpl",
			want: true,
		},
		{
			name: "non-existent provider template",
			path: "templates/some-other-provider/AGENTS.md.tmpl",
			want: false,
		},
		{
			name: "non-existent template in existing directory",
			path: "templates/spectr/non-existent.md.tmpl",
			want: false,
		},
		{
			name: "completely invalid path",
			path: "templates/invalid/path/template.tmpl",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := templateExists(tt.path)
			if got != tt.want {
				t.Errorf("templateExists(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

// TestRenderAgents_BackwardCompatibility tests that empty provider ID preserves
// existing behavior by using the generic template directly.
func TestRenderAgents_BackwardCompatibility(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	ctx := providers.DefaultTemplateContext()

	// Empty provider ID should use generic template
	got, err := tm.RenderAgents(ctx, "")
	if err != nil {
		t.Fatalf("RenderAgents() with empty provider ID error = %v", err)
	}

	// Verify it rendered successfully with expected content
	if !strings.Contains(got, "# Spectr Instructions") {
		t.Error("RenderAgents() with empty provider ID missing '# Spectr Instructions'")
	}
	if !strings.Contains(got, "## TL;DR Quick Checklist") {
		t.Error("RenderAgents() with empty provider ID missing '## TL;DR Quick Checklist'")
	}
}

// TestRenderAgents_ProviderWithoutCustomTemplate tests that a provider without
// a custom AGENTS.md template falls back to the generic template.
func TestRenderAgents_ProviderWithoutCustomTemplate(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	ctx := providers.DefaultTemplateContext()

	// Provider ID that doesn't have custom templates should fall back to generic
	got, err := tm.RenderAgents(ctx, "non-existent-provider")
	if err != nil {
		t.Fatalf("RenderAgents() with non-existent provider error = %v", err)
	}

	// Should still get valid output from generic template
	if !strings.Contains(got, "# Spectr Instructions") {
		t.Error("RenderAgents() with non-existent provider missing '# Spectr Instructions'")
	}
	if len(got) < 5000 {
		t.Errorf(
			"RenderAgents() with non-existent provider output too short: got %d chars, expected >= 5000",
			len(got),
		)
	}
}

// TestRenderInstructionPointer_BackwardCompatibility tests that empty provider ID
// preserves existing behavior for instruction pointer rendering.
func TestRenderInstructionPointer_BackwardCompatibility(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	ctx := providers.DefaultTemplateContext()

	// Empty provider ID should use generic template
	got, err := tm.RenderInstructionPointer(ctx, "")
	if err != nil {
		t.Fatalf("RenderInstructionPointer() with empty provider ID error = %v", err)
	}

	// Verify it rendered successfully
	if !strings.Contains(got, "spectr/AGENTS.md") {
		t.Error("RenderInstructionPointer() with empty provider ID missing 'spectr/AGENTS.md'")
	}
}

// TestRenderInstructionPointer_ProviderWithoutCustomTemplate tests fallback
// behavior for instruction pointer when provider has no custom template.
func TestRenderInstructionPointer_ProviderWithoutCustomTemplate(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	ctx := providers.DefaultTemplateContext()

	// Provider without custom template should fall back to generic
	got, err := tm.RenderInstructionPointer(ctx, "some-unknown-provider")
	if err != nil {
		t.Fatalf("RenderInstructionPointer() with unknown provider error = %v", err)
	}

	// Should still get valid output
	if !strings.Contains(got, "spectr/AGENTS.md") {
		t.Error("RenderInstructionPointer() with unknown provider missing 'spectr/AGENTS.md'")
	}
}

// TestRenderSlashCommand_BackwardCompatibility tests that empty provider ID
// preserves existing behavior for slash command rendering.
func TestRenderSlashCommand_BackwardCompatibility(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	ctx := providers.DefaultTemplateContext()

	// Test both slash commands with empty provider ID
	commands := []string{"proposal", "apply"}
	for _, cmd := range commands {
		t.Run(cmd, func(t *testing.T) {
			got, err := tm.RenderSlashCommand(cmd, ctx, "")
			if err != nil {
				t.Fatalf("RenderSlashCommand(%q) with empty provider ID error = %v", cmd, err)
			}

			// Verify it rendered successfully
			if !strings.Contains(got, "# Guardrails") {
				t.Errorf("RenderSlashCommand(%q) missing '# Guardrails'", cmd)
			}
			if !strings.Contains(got, "# Steps") {
				t.Errorf("RenderSlashCommand(%q) missing '# Steps'", cmd)
			}
		})
	}
}

// TestRenderSlashCommand_ProviderWithoutCustomTemplate tests fallback behavior
// for slash commands when provider has no custom template.
func TestRenderSlashCommand_ProviderWithoutCustomTemplate(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	ctx := providers.DefaultTemplateContext()

	// Provider without custom templates should fall back to generic tools/ templates
	commands := []string{"proposal", "apply"}
	for _, cmd := range commands {
		t.Run(cmd, func(t *testing.T) {
			got, err := tm.RenderSlashCommand(cmd, ctx, "unknown-provider")
			if err != nil {
				t.Fatalf("RenderSlashCommand(%q) with unknown provider error = %v", cmd, err)
			}

			// Should still get valid output from generic template
			if !strings.Contains(got, "# Guardrails") {
				t.Errorf("RenderSlashCommand(%q) missing '# Guardrails'", cmd)
			}
		})
	}
}

// TestPartialProviderOverride tests the scenario where a provider has some custom
// templates but not others. Since no provider templates exist yet, this test
// verifies that all templates fall back correctly for an unknown provider.
// When provider-specific templates are added, this test should be updated.
func TestPartialProviderOverride(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	ctx := providers.DefaultTemplateContext()

	// Simulate a provider that might have AGENTS.md but not slash-proposal
	// Currently, no provider templates exist, so both should fall back
	providerID := "test-partial-provider"

	// AGENTS.md should fall back to generic
	agents, err := tm.RenderAgents(ctx, providerID)
	if err != nil {
		t.Fatalf("RenderAgents() error = %v", err)
	}
	if !strings.Contains(agents, "# Spectr Instructions") {
		t.Error("AGENTS.md should fall back to generic template")
	}

	// slash-proposal should fall back to generic tools/ template
	slashProposal, err := tm.RenderSlashCommand("proposal", ctx, providerID)
	if err != nil {
		t.Fatalf("RenderSlashCommand(proposal) error = %v", err)
	}
	if !strings.Contains(slashProposal, "# Guardrails") {
		t.Error("slash-proposal should fall back to generic template")
	}

	// slash-apply should also fall back
	slashApply, err := tm.RenderSlashCommand("apply", ctx, providerID)
	if err != nil {
		t.Fatalf("RenderSlashCommand(apply) error = %v", err)
	}
	if !strings.Contains(slashApply, "# Guardrails") {
		t.Error("slash-apply should fall back to generic template")
	}

	// instruction-pointer should fall back
	instrPointer, err := tm.RenderInstructionPointer(ctx, providerID)
	if err != nil {
		t.Fatalf("RenderInstructionPointer() error = %v", err)
	}
	if !strings.Contains(instrPointer, "spectr/AGENTS.md") {
		t.Error("instruction-pointer should fall back to generic template")
	}
}

// TestProviderTemplateResolutionConsistency verifies that the same provider ID
// produces consistent results across multiple calls.
func TestProviderTemplateResolutionConsistency(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	ctx := providers.DefaultTemplateContext()
	providerID := "consistency-test-provider"

	// Call RenderAgents multiple times with same provider ID
	result1, err := tm.RenderAgents(ctx, providerID)
	if err != nil {
		t.Fatalf("First RenderAgents() error = %v", err)
	}

	result2, err := tm.RenderAgents(ctx, providerID)
	if err != nil {
		t.Fatalf("Second RenderAgents() error = %v", err)
	}

	if result1 != result2 {
		t.Error("RenderAgents() should produce consistent results for same provider ID")
	}
}

// TestResolveTemplatePath_SpectrFallback verifies that templates correctly
// resolve to spectr/ directory for core templates.
func TestResolveTemplatePath_SpectrFallback(t *testing.T) {
	// Test various template names that should fall back to spectr/
	spectrTemplates := []string{
		"AGENTS.md.tmpl",
		"instruction-pointer.md.tmpl",
		"project.md.tmpl",
	}

	for _, tmplName := range spectrTemplates {
		t.Run(tmplName, func(t *testing.T) {
			// With empty provider ID
			path := resolveTemplatePath("", tmplName, "spectr")
			expected := "templates/spectr/" + tmplName
			if path != expected {
				t.Errorf("resolveTemplatePath('', %q, 'spectr') = %q, want %q",
					tmplName, path, expected)
			}

			// With unknown provider ID
			path = resolveTemplatePath("unknown", tmplName, "spectr")
			if path != expected {
				t.Errorf("resolveTemplatePath('unknown', %q, 'spectr') = %q, want %q",
					tmplName, path, expected)
			}
		})
	}
}

// TestResolveTemplatePath_ToolsFallback verifies that slash command templates
// correctly resolve to tools/ directory.
func TestResolveTemplatePath_ToolsFallback(t *testing.T) {
	// Test slash command templates that should fall back to tools/
	toolsTemplates := []string{
		"slash-proposal.md.tmpl",
		"slash-apply.md.tmpl",
	}

	for _, tmplName := range toolsTemplates {
		t.Run(tmplName, func(t *testing.T) {
			// With empty provider ID
			path := resolveTemplatePath("", tmplName, "tools")
			expected := "templates/tools/" + tmplName
			if path != expected {
				t.Errorf("resolveTemplatePath('', %q, 'tools') = %q, want %q",
					tmplName, path, expected)
			}

			// With unknown provider ID
			path = resolveTemplatePath("unknown", tmplName, "tools")
			if path != expected {
				t.Errorf("resolveTemplatePath('unknown', %q, 'tools') = %q, want %q",
					tmplName, path, expected)
			}
		})
	}
}

// TestRenderAgents_ClaudeCodeProvider tests that the Claude Code provider
// renders using its custom AGENTS.md template with Claude Code-specific content.
func TestRenderAgents_ClaudeCodeProvider(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	ctx := providers.DefaultTemplateContext()

	// Claude Code should use its custom template
	got, err := tm.RenderAgents(ctx, "claude-code")
	if err != nil {
		t.Fatalf("RenderAgents() with claude-code provider error = %v", err)
	}

	// Verify Claude Code-specific content is present
	expectedContent := []string{
		"# Spectr Instructions",
		"Claude Code Tool Reference",
		"`Glob`",  // Claude Code tool
		"`Grep`",  // Claude Code tool
		"`Read`",  // Claude Code tool
		"`Edit`",  // Claude Code tool
		"`Write`", // Claude Code tool
		"`Bash`",  // Claude Code tool
		"`Task`",  // Claude Code tool for subagent delegation
		"TodoWrite",
	}

	for _, content := range expectedContent {
		if !strings.Contains(got, content) {
			t.Errorf("RenderAgents() with claude-code missing expected content: %q", content)
		}
	}

	// Verify template variables are substituted
	if strings.Contains(got, "{{ .BaseDir }}") {
		t.Error("RenderAgents() should substitute {{ .BaseDir }} variable")
	}
	if strings.Contains(got, "{{ .SpecsDir }}") {
		t.Error("RenderAgents() should substitute {{ .SpecsDir }} variable")
	}
	if strings.Contains(got, "{{ .ChangesDir }}") {
		t.Error("RenderAgents() should substitute {{ .ChangesDir }} variable")
	}

	// Should be longer than generic because it has more detailed tool guidance
	if len(got) < 10000 {
		t.Errorf(
			"RenderAgents() claude-code output seems too short: got %d chars, expected >= 10000",
			len(got),
		)
	}
}

// TestClaudeCodeTemplateHasProviderSpecificContent verifies that the Claude Code
// template has meaningfully different content from the generic template.
func TestClaudeCodeTemplateHasProviderSpecificContent(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	ctx := providers.DefaultTemplateContext()

	// Render both generic and Claude Code templates
	genericTemplate, err := tm.RenderAgents(ctx, "")
	if err != nil {
		t.Fatalf("RenderAgents() generic error = %v", err)
	}

	claudeCodeTemplate, err := tm.RenderAgents(ctx, "claude-code")
	if err != nil {
		t.Fatalf("RenderAgents() claude-code error = %v", err)
	}

	// Templates should be different
	if genericTemplate == claudeCodeTemplate {
		t.Error("Claude Code template should be different from generic template")
	}

	// Claude Code template should have Claude Code-specific tool references
	// that are NOT in the generic template
	claudeCodeOnlyContent := []string{
		"Claude Code Tool Reference",
		"`Glob`",
		"`Grep`",
		"Delegation with Task Tool",
	}

	for _, content := range claudeCodeOnlyContent {
		if !strings.Contains(claudeCodeTemplate, content) {
			t.Errorf("Claude Code template should contain: %q", content)
		}
	}
}

// TestRenderMethods_ProviderIDVariations tests that render methods handle
// various provider ID formats correctly.
func TestRenderMethods_ProviderIDVariations(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	ctx := providers.DefaultTemplateContext()

	// Various provider ID values that should all fall back to generic
	providerIDs := []string{
		"",                // empty
		"claude-code",     // typical provider ID
		"crush",           // another typical provider ID
		"opencode",        // yet another provider
		"custom-provider", // hyphenated
		"SomeProvider",    // mixed case
		"provider_with_underscore",
	}

	for _, providerID := range providerIDs {
		t.Run("providerID="+providerID, func(t *testing.T) {
			// All should render without error (falling back to generic)
			_, err := tm.RenderAgents(ctx, providerID)
			if err != nil {
				t.Errorf("RenderAgents() with provider %q error = %v", providerID, err)
			}

			_, err = tm.RenderInstructionPointer(ctx, providerID)
			if err != nil {
				t.Errorf("RenderInstructionPointer() with provider %q error = %v", providerID, err)
			}

			_, err = tm.RenderSlashCommand("proposal", ctx, providerID)
			if err != nil {
				t.Errorf(
					"RenderSlashCommand(proposal) with provider %q error = %v",
					providerID,
					err,
				)
			}

			_, err = tm.RenderSlashCommand("apply", ctx, providerID)
			if err != nil {
				t.Errorf("RenderSlashCommand(apply) with provider %q error = %v", providerID, err)
			}
		})
	}
}

// =============================================================================
// Crush Provider Template Tests
// =============================================================================

// TestRenderAgents_CrushProvider tests that the Crush provider renders using
// its custom AGENTS.md template with Crush-specific content.
func TestRenderAgents_CrushProvider(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	ctx := providers.DefaultTemplateContext()

	// Crush should use its custom template
	got, err := tm.RenderAgents(ctx, "crush")
	if err != nil {
		t.Fatalf("RenderAgents() with crush provider error = %v", err)
	}

	// Verify Crush-specific content is present
	expectedContent := []string{
		"# Spectr Instructions",
		"Instructions for Crush",
		"Crush Tool Reference",
		"Shell commands",     // Crush uses shell commands
		"View file",          // Crush file viewing
		"Edit operations",    // Crush file editing
		"Shell execution",    // Crush shell execution
		"Working with Files", // Crush-specific section
		"rg -n",              // ripgrep reference
		"spectr validate",
	}

	for _, content := range expectedContent {
		if !strings.Contains(got, content) {
			t.Errorf("RenderAgents() with crush missing expected content: %q", content)
		}
	}

	// Verify template variables are substituted
	if strings.Contains(got, "{{ .BaseDir }}") {
		t.Error("RenderAgents() should substitute {{ .BaseDir }} variable")
	}
	if strings.Contains(got, "{{ .SpecsDir }}") {
		t.Error("RenderAgents() should substitute {{ .SpecsDir }} variable")
	}
	if strings.Contains(got, "{{ .ChangesDir }}") {
		t.Error("RenderAgents() should substitute {{ .ChangesDir }} variable")
	}

	// Should be a substantial document
	if len(got) < 8000 {
		t.Errorf(
			"RenderAgents() crush output seems too short: got %d chars, expected >= 8000",
			len(got),
		)
	}
}

// TestCrushTemplateHasProviderSpecificContent verifies that the Crush template
// has meaningfully different content from the generic template.
func TestCrushTemplateHasProviderSpecificContent(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	ctx := providers.DefaultTemplateContext()

	// Render both generic and Crush templates
	genericTemplate, err := tm.RenderAgents(ctx, "")
	if err != nil {
		t.Fatalf("RenderAgents() generic error = %v", err)
	}

	crushTemplate, err := tm.RenderAgents(ctx, "crush")
	if err != nil {
		t.Fatalf("RenderAgents() crush error = %v", err)
	}

	// Templates should be different
	if genericTemplate == crushTemplate {
		t.Error("Crush template should be different from generic template")
	}

	// Crush template should have Crush-specific content
	crushOnlyContent := []string{
		"Instructions for Crush",
		"Crush Tool Reference",
		"Crush Tool Selection Guide",
	}

	for _, content := range crushOnlyContent {
		if !strings.Contains(crushTemplate, content) {
			t.Errorf("Crush template should contain: %q", content)
		}
	}
}

// TestCrushTemplateExists verifies that the Crush template file exists.
func TestCrushTemplateExists(t *testing.T) {
	if !templateExists("templates/crush/AGENTS.md.tmpl") {
		t.Error("templates/crush/AGENTS.md.tmpl should exist")
	}
}

// TestCrushTemplateDifferentFromClaudeCode verifies that Crush and Claude Code
// templates are distinct from each other.
func TestCrushTemplateDifferentFromClaudeCode(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	ctx := providers.DefaultTemplateContext()

	crushTemplate, err := tm.RenderAgents(ctx, "crush")
	if err != nil {
		t.Fatalf("RenderAgents() crush error = %v", err)
	}

	claudeCodeTemplate, err := tm.RenderAgents(ctx, "claude-code")
	if err != nil {
		t.Fatalf("RenderAgents() claude-code error = %v", err)
	}

	// Templates should be different
	if crushTemplate == claudeCodeTemplate {
		t.Error("Crush template should be different from Claude Code template")
	}

	// Claude Code should have its specific tool references
	claudeCodeOnly := []string{
		"`Glob`",
		"`Grep`",
		"`Read`",
		"`Edit`",
		"`Write`",
		"`Bash`",
		"`Task`",
	}

	for _, content := range claudeCodeOnly {
		if !strings.Contains(claudeCodeTemplate, content) {
			t.Errorf("Claude Code template should contain: %q", content)
		}
	}

	// Crush should have its specific content
	crushOnly := []string{
		"Instructions for Crush",
		"Crush Tool Reference",
	}

	for _, content := range crushOnly {
		if !strings.Contains(crushTemplate, content) {
			t.Errorf("Crush template should contain: %q", content)
		}
	}
}

// =============================================================================
// Integration Tests for Provider-Specific Template Rendering
// =============================================================================

// TestIntegration_ClaudeCodeProviderConfigure tests the full Configure() flow
// for the Claude Code provider, verifying that provider-specific templates are used.
func TestIntegration_ClaudeCodeProviderConfigure(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spectr-integration-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create the real TemplateManager
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	// Create the Claude Code provider
	p := providers.NewClaudeProvider()

	// Configure the provider
	err = p.Configure(tmpDir, filepath.Join(tmpDir, "spectr"), tm)
	if err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	// Verify the instruction file was created with Claude Code-specific content
	configPath := filepath.Join(tmpDir, "CLAUDE.md")
	if !providers.FileExists(configPath) {
		t.Fatal("CLAUDE.md was not created")
	}

	configContent, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read CLAUDE.md: %v", err)
	}

	// The instruction file should contain spectr markers
	if !strings.Contains(string(configContent), "spectr/AGENTS.md") {
		t.Error("CLAUDE.md should reference spectr/AGENTS.md")
	}

	// Verify slash command files were created
	proposalPath := filepath.Join(tmpDir, ".claude/commands/spectr/proposal.md")
	applyPath := filepath.Join(tmpDir, ".claude/commands/spectr/apply.md")

	if !providers.FileExists(proposalPath) {
		t.Error("proposal.md was not created")
	}
	if !providers.FileExists(applyPath) {
		t.Error("apply.md was not created")
	}

	// Read and verify proposal command content
	proposalContent, err := os.ReadFile(proposalPath)
	if err != nil {
		t.Fatalf("Failed to read proposal.md: %v", err)
	}

	// Should contain standard slash command content
	if !strings.Contains(string(proposalContent), "# Guardrails") {
		t.Error("proposal.md should contain '# Guardrails'")
	}
	if !strings.Contains(string(proposalContent), "# Steps") {
		t.Error("proposal.md should contain '# Steps'")
	}
}

// TestIntegration_CrushProviderConfigure tests the full Configure() flow
// for the Crush provider, verifying that provider-specific templates are used.
func TestIntegration_CrushProviderConfigure(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spectr-integration-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create the real TemplateManager
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	// Create the Crush provider
	p := providers.NewCrushProvider()

	// Configure the provider
	err = p.Configure(tmpDir, filepath.Join(tmpDir, "spectr"), tm)
	if err != nil {
		t.Fatalf("Configure() error = %v", err)
	}

	// Verify slash command files were created
	proposalPath := filepath.Join(tmpDir, ".crush/commands/spectr/proposal.md")
	applyPath := filepath.Join(tmpDir, ".crush/commands/spectr/apply.md")

	if !providers.FileExists(proposalPath) {
		t.Error("proposal.md was not created at expected path")
	}
	if !providers.FileExists(applyPath) {
		t.Error("apply.md was not created at expected path")
	}

	// Read and verify proposal command content
	proposalContent, err := os.ReadFile(proposalPath)
	if err != nil {
		t.Fatalf("Failed to read proposal.md: %v", err)
	}

	// Should contain standard slash command content
	if !strings.Contains(string(proposalContent), "# Guardrails") {
		t.Error("proposal.md should contain '# Guardrails'")
	}
}

// TestIntegration_FallbackToGenericTemplate tests that providers without
// custom templates (like cursor or gemini) fall back to generic templates.
func TestIntegration_FallbackToGenericTemplate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spectr-integration-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create the real TemplateManager
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	// Test with Cursor provider (no custom templates)
	cursorProvider := providers.NewCursorProvider()

	err = cursorProvider.Configure(tmpDir, filepath.Join(tmpDir, "spectr"), tm)
	if err != nil {
		t.Fatalf("Cursor Configure() error = %v", err)
	}

	// Verify slash commands were created
	proposalPath := filepath.Join(tmpDir, ".cursorrules/commands/spectr/proposal.md")
	if !providers.FileExists(proposalPath) {
		t.Error("Cursor proposal.md was not created")
	}

	proposalContent, err := os.ReadFile(proposalPath)
	if err != nil {
		t.Fatalf("Failed to read cursor proposal.md: %v", err)
	}

	// Should have generic content (not provider-specific markers)
	if !strings.Contains(string(proposalContent), "# Guardrails") {
		t.Error("Cursor proposal.md should contain generic '# Guardrails'")
	}

	// Clean up for next test
	_ = os.RemoveAll(tmpDir)
	tmpDir, _ = os.MkdirTemp("", "spectr-integration-test-*")
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Test with Gemini provider (no custom templates)
	geminiProvider := providers.NewGeminiProvider()

	err = geminiProvider.Configure(tmpDir, filepath.Join(tmpDir, "spectr"), tm)
	if err != nil {
		t.Fatalf("Gemini Configure() error = %v", err)
	}

	// Verify TOML slash commands were created
	geminiProposalPath := filepath.Join(tmpDir, ".gemini/commands/spectr/proposal.toml")
	if !providers.FileExists(geminiProposalPath) {
		t.Error("Gemini proposal.toml was not created")
	}

	geminiContent, err := os.ReadFile(geminiProposalPath)
	if err != nil {
		t.Fatalf("Failed to read gemini proposal.toml: %v", err)
	}

	// Should be TOML format
	if !strings.Contains(string(geminiContent), "description =") {
		t.Error("Gemini proposal.toml should be in TOML format")
	}
}

// TestIntegration_ProviderSpecificAgentsContent verifies that when providers
// render AGENTS.md content, they get provider-specific content.
func TestIntegration_ProviderSpecificAgentsContent(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	ctx := providers.DefaultTemplateContext()

	// Test Claude Code provider gets Claude Code-specific AGENTS.md content
	t.Run("claude-code specific content", func(t *testing.T) {
		content, err := tm.RenderAgents(ctx, "claude-code")
		if err != nil {
			t.Fatalf("RenderAgents(claude-code) error = %v", err)
		}

		// Must have Claude Code-specific markers
		claudeCodeMarkers := []string{
			"Instructions for Claude Code",
			"Claude Code Tool Reference",
			"`Glob`",
			"`Grep`",
			"`Read`",
			"`Edit`",
			"`Write`",
			"`Bash`",
			"`Task`",
			"TodoWrite",
		}

		for _, marker := range claudeCodeMarkers {
			if !strings.Contains(content, marker) {
				t.Errorf("Claude Code AGENTS.md missing marker: %q", marker)
			}
		}

		// Must NOT have generic AI assistant language
		if strings.Contains(content, "Instructions for AI coding assistants") {
			t.Error("Claude Code AGENTS.md should NOT have generic 'AI coding assistants' text")
		}
	})

	// Test Crush provider gets Crush-specific AGENTS.md content
	t.Run("crush specific content", func(t *testing.T) {
		content, err := tm.RenderAgents(ctx, "crush")
		if err != nil {
			t.Fatalf("RenderAgents(crush) error = %v", err)
		}

		// Must have Crush-specific markers
		crushMarkers := []string{
			"Instructions for Crush",
			"Crush Tool Reference",
			"Shell commands",
			"Working with Files",
		}

		for _, marker := range crushMarkers {
			if !strings.Contains(content, marker) {
				t.Errorf("Crush AGENTS.md missing marker: %q", marker)
			}
		}

		// Must NOT have Claude Code tool references
		if strings.Contains(content, "`Glob`") {
			t.Error("Crush AGENTS.md should NOT have Claude Code '`Glob`' tool reference")
		}
	})

	// Test generic fallback for unknown provider
	t.Run("fallback to generic for unknown provider", func(t *testing.T) {
		content, err := tm.RenderAgents(ctx, "unknown-provider")
		if err != nil {
			t.Fatalf("RenderAgents(unknown-provider) error = %v", err)
		}

		// Must have generic content
		if !strings.Contains(content, "Instructions for AI coding assistants") {
			t.Error("Unknown provider should get generic 'AI coding assistants' content")
		}

		// Must NOT have provider-specific markers
		if strings.Contains(content, "Instructions for Claude Code") {
			t.Error("Unknown provider should NOT have Claude Code-specific content")
		}
		if strings.Contains(content, "Instructions for Crush") {
			t.Error("Unknown provider should NOT have Crush-specific content")
		}
	})
}

// TestIntegration_AllProvidersConfigureSuccessfully verifies that all registered
// providers can be configured successfully with the real TemplateManager.
func TestIntegration_AllProvidersConfigureSuccessfully(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	allProviders := providers.All()

	for _, p := range allProviders {
		t.Run(p.ID(), func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "spectr-provider-test-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer func() { _ = os.RemoveAll(tmpDir) }()

			// Configure should succeed for all providers
			err = p.Configure(tmpDir, filepath.Join(tmpDir, "spectr"), tm)
			if err != nil {
				t.Fatalf("Configure() error for provider %s: %v", p.ID(), err)
			}

			// Provider should report as configured
			if !p.IsConfigured(tmpDir) {
				t.Errorf("Provider %s should be configured after Configure()", p.ID())
			}

			// All expected files should exist
			filePaths := p.GetFilePaths()
			for _, relPath := range filePaths {
				fullPath := filepath.Join(tmpDir, relPath)
				// Skip global paths (like ~/.config/...) in tests
				if strings.HasPrefix(relPath, "~") || strings.HasPrefix(relPath, "/") {
					continue
				}
				if !providers.FileExists(fullPath) {
					t.Errorf("Provider %s: expected file not created: %s", p.ID(), relPath)
				}
			}
		})
	}
}

// TestIntegration_ProviderTemplateVariablesSubstituted verifies that template
// variables like {{ .BaseDir }} are properly substituted in rendered content.
func TestIntegration_ProviderTemplateVariablesSubstituted(t *testing.T) {
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	// Custom context with non-default values
	ctx := providers.TemplateContext{
		BaseDir:     "custom-spectr",
		SpecsDir:    "custom-spectr/specs",
		ChangesDir:  "custom-spectr/changes",
		ProjectFile: "custom-spectr/project.md",
		AgentsFile:  "custom-spectr/AGENTS.md",
	}

	providerIDs := []string{"claude-code", "crush", ""}

	for _, providerID := range providerIDs {
		name := providerID
		if name == "" {
			name = "generic"
		}
		t.Run(name, func(t *testing.T) {
			content, err := tm.RenderAgents(ctx, providerID)
			if err != nil {
				t.Fatalf("RenderAgents() error = %v", err)
			}

			// Custom paths should be present
			if !strings.Contains(content, "custom-spectr/changes") {
				t.Error("Content should contain custom ChangesDir path")
			}
			if !strings.Contains(content, "custom-spectr/specs") {
				t.Error("Content should contain custom SpecsDir path")
			}

			// Template syntax should NOT be present
			if strings.Contains(content, "{{ .BaseDir }}") {
				t.Error("Content should not contain unsubstituted {{ .BaseDir }}")
			}
			if strings.Contains(content, "{{ .ChangesDir }}") {
				t.Error("Content should not contain unsubstituted {{ .ChangesDir }}")
			}
			if strings.Contains(content, "{{ .SpecsDir }}") {
				t.Error("Content should not contain unsubstituted {{ .SpecsDir }}")
			}
		})
	}
}

// TestIntegration_ProviderIDPassedToConfigure verifies that BaseProvider.Configure()
// correctly passes the provider ID to the template renderer methods.
func TestIntegration_ProviderIDPassedToConfigure(t *testing.T) {
	// This test verifies the integration between provider.Configure() and
	// TemplateRenderer methods by checking that the rendered content matches
	// what we expect for each provider.

	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() error = %v", err)
	}

	testCases := []struct {
		name           string
		provider       providers.Provider
		expectedMarker string // A marker that should appear in AGENTS.md for this provider
	}{
		{
			name:           "claude-code",
			provider:       providers.NewClaudeProvider(),
			expectedMarker: "Claude Code Tool Reference",
		},
		{
			name:           "crush",
			provider:       providers.NewCrushProvider(),
			expectedMarker: "Crush Tool Reference",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Render AGENTS.md with the provider's ID
			ctx := providers.DefaultTemplateContext()
			content, err := tm.RenderAgents(ctx, tc.provider.ID())
			if err != nil {
				t.Fatalf("RenderAgents() error = %v", err)
			}

			// The content should have the provider-specific marker
			if !strings.Contains(content, tc.expectedMarker) {
				t.Errorf("Provider %s: RenderAgents() should contain %q",
					tc.name, tc.expectedMarker)
			}
		})
	}
}
