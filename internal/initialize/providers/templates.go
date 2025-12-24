package providers

// TemplateRenderer is a type-safe function reference for template rendering.
// Using a function type ensures compile-time validation of template references.
type TemplateRenderer func(
	tm TemplateManager,
	ctx TemplateContext,
) (string, error)

// RenderInstructionPointer renders the instruction pointer template
// that directs AI assistants to read spectr/AGENTS.md.
var RenderInstructionPointer TemplateRenderer = func(
	tm TemplateManager,
	ctx TemplateContext,
) (string, error) {
	return tm.RenderInstructionPointer(ctx)
}

// SlashCommand is a type-safe slash command definition.
type SlashCommand struct {
	Name        string
	Description string
	Renderer    TemplateRenderer
}

// SlashProposal is the /proposal command for creating change proposals.
var SlashProposal = SlashCommand{
	Name:        "proposal",
	Description: "Scaffold a new Spectr change and validate strictly.",
	Renderer: func(tm TemplateManager, ctx TemplateContext) (string, error) {
		return tm.RenderSlashCommand("proposal", ctx)
	},
}

// SlashApply is the /apply command for implementing approved changes.
var SlashApply = SlashCommand{
	Name:        "apply",
	Description: "Implement an approved Spectr change and keep tasks in sync.",
	Renderer: func(tm TemplateManager, ctx TemplateContext) (string, error) {
		return tm.RenderSlashCommand("apply", ctx)
	},
}

// DefaultSlashCommands returns the standard set of slash commands.
func DefaultSlashCommands() []SlashCommand {
	return []SlashCommand{SlashProposal, SlashApply}
}
