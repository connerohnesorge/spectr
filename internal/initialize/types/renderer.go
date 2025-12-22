package types

// TemplateRenderer provides template rendering capabilities.
type TemplateRenderer interface {
	// RenderAgents renders the AGENTS.md template content.
	RenderAgents(
		ctx TemplateContext,
	) (string, error)
	// RenderInstructionPointer renders a short pointer template that directs
	// AI assistants to read spectr/AGENTS.md for full instructions.
	RenderInstructionPointer(
		ctx TemplateContext,
	) (string, error)
	// RenderSlashCommand renders a slash command template (proposal or apply).
	RenderSlashCommand(
		command string,
		ctx TemplateContext,
	) (string, error)
}
