package initialize_test

import (
	"testing"

	"github.com/connerohnesorge/spectr/internal/initialize"
	"github.com/connerohnesorge/spectr/internal/initialize/providers"
	"github.com/connerohnesorge/spectr/internal/initialize/templates"
)

func TestTemplateManagerAccessors(t *testing.T) {
	// Create template manager
	tm, err := initialize.NewTemplateManager()
	if err != nil {
		t.Fatalf(
			"NewTemplateManager() error = %v",
			err,
		)
	}

	ctx := providers.DefaultTemplateContext()

	t.Run(
		"InstructionPointer",
		func(t *testing.T) {
			ref := tm.InstructionPointer()
			content, err := ref.Render(ctx)
			if err != nil {
				t.Errorf(
					"InstructionPointer().Render() error = %v",
					err,
				)
			}
			if len(content) == 0 {
				t.Error(
					"InstructionPointer().Render() returned empty content",
				)
			}
		},
	)

	t.Run("Agents", func(t *testing.T) {
		ref := tm.Agents()
		content, err := ref.Render(ctx)
		if err != nil {
			t.Errorf(
				"Agents().Render() error = %v",
				err,
			)
		}
		if len(content) == 0 {
			t.Error(
				"Agents().Render() returned empty content",
			)
		}
	})

	t.Run("Project", func(t *testing.T) {
		ref := tm.Project()
		// Project template uses ProjectContext, not TemplateContext
		projectCtx := initialize.ProjectContext{
			ProjectName: "test-project",
			Description: "A test project",
			TechStack:   []string{"Go", "Python"},
			Conventions: "Test conventions",
		}
		content, err := ref.Render(projectCtx)
		if err != nil {
			t.Errorf(
				"Project().Render() error = %v",
				err,
			)
		}
		if len(content) == 0 {
			t.Error(
				"Project().Render() returned empty content",
			)
		}
	})

	t.Run("CIWorkflow", func(t *testing.T) {
		ref := tm.CIWorkflow()
		// CIWorkflow template has no variables, pass empty context
		content, err := ref.Render(
			providers.TemplateContext{},
		)
		if err != nil {
			t.Errorf(
				"CIWorkflow().Render() error = %v",
				err,
			)
		}
		if len(content) == 0 {
			t.Error(
				"CIWorkflow().Render() returned empty content",
			)
		}
	})

	t.Run(
		"SlashCommand_Proposal",
		func(t *testing.T) {
			ref := tm.SlashCommand(
				templates.SlashProposal,
			)
			content, err := ref.Render(ctx)
			if err != nil {
				t.Errorf(
					"SlashCommand(SlashProposal).Render() error = %v",
					err,
				)
			}
			if len(content) == 0 {
				t.Error(
					"SlashCommand(SlashProposal).Render() returned empty content",
				)
			}
		},
	)

	t.Run(
		"SlashCommand_Apply",
		func(t *testing.T) {
			ref := tm.SlashCommand(
				templates.SlashApply,
			)
			content, err := ref.Render(ctx)
			if err != nil {
				t.Errorf(
					"SlashCommand(SlashApply).Render() error = %v",
					err,
				)
			}
			if len(content) == 0 {
				t.Error(
					"SlashCommand(SlashApply).Render() returned empty content",
				)
			}
		},
	)
}
