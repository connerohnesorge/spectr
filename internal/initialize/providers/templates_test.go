package providers

import (
	"testing"
)

func TestRenderInstructionPointer(t *testing.T) {
	tm := &mockTemplateManager{
		content: "Test instruction content",
	}
	ctx := TemplateContext{
		BaseDir:     "spectr",
		SpecsDir:    "spectr/specs",
		ChangesDir:  "spectr/changes",
		ProjectFile: "spectr/project.md",
		AgentsFile:  "spectr/AGENTS.md",
	}

	content, err := RenderInstructionPointer(tm, ctx)
	if err != nil {
		t.Fatalf("RenderInstructionPointer() error = %v", err)
	}

	if content != "Test instruction content" {
		t.Errorf(
			"RenderInstructionPointer() = %v, want Test instruction content",
			content,
		)
	}
}

func TestSlashProposal(t *testing.T) {
	if SlashProposal.Name != "proposal" {
		t.Errorf("SlashProposal.Name = %v, want proposal", SlashProposal.Name)
	}

	expectedDesc := "Scaffold a new Spectr change and validate strictly."
	if SlashProposal.Description != expectedDesc {
		t.Errorf(
			"SlashProposal.Description = %v, want %v",
			SlashProposal.Description,
			expectedDesc,
		)
	}

	tm := &mockTemplateManager{
		content: "Proposal content",
	}
	ctx := TemplateContext{
		BaseDir: "spectr",
	}

	content, err := SlashProposal.Renderer(tm, ctx)
	if err != nil {
		t.Fatalf("SlashProposal.Renderer() error = %v", err)
	}

	if content != "Proposal content" {
		t.Errorf(
			"SlashProposal.Renderer() = %v, want Proposal content",
			content,
		)
	}
}

func TestSlashApply(t *testing.T) {
	if SlashApply.Name != "apply" {
		t.Errorf("SlashApply.Name = %v, want apply", SlashApply.Name)
	}

	expectedDesc := "Implement an approved Spectr change and keep tasks in sync."
	if SlashApply.Description != expectedDesc {
		t.Errorf(
			"SlashApply.Description = %v, want %v",
			SlashApply.Description,
			expectedDesc,
		)
	}

	tm := &mockTemplateManager{
		content: "Apply content",
	}
	ctx := TemplateContext{
		BaseDir: "spectr",
	}

	content, err := SlashApply.Renderer(tm, ctx)
	if err != nil {
		t.Fatalf("SlashApply.Renderer() error = %v", err)
	}

	if content != "Apply content" {
		t.Errorf("SlashApply.Renderer() = %v, want Apply content", content)
	}
}

func TestDefaultSlashCommands(t *testing.T) {
	commands := DefaultSlashCommands()

	if len(commands) != 2 {
		t.Fatalf("DefaultSlashCommands() returned %d commands, want 2", len(commands))
	}

	if commands[0].Name != "proposal" {
		t.Errorf("commands[0].Name = %v, want proposal", commands[0].Name)
	}

	if commands[1].Name != "apply" {
		t.Errorf("commands[1].Name = %v, want apply", commands[1].Name)
	}
}
