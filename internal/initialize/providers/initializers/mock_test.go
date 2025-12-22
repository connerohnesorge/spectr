package initializers

import "github.com/connerohnesorge/spectr/internal/initialize/types"

type MockTemplateRenderer struct {
	RenderAgentsFunc             func(ctx types.TemplateContext) (string, error)
	RenderInstructionPointerFunc func(ctx types.TemplateContext) (string, error)
	RenderSlashCommandFunc       func(cmd string, ctx types.TemplateContext) (string, error)
}

func (m *MockTemplateRenderer) RenderAgents(
	ctx types.TemplateContext,
) (string, error) {
	if m.RenderAgentsFunc != nil {
		return m.RenderAgentsFunc(ctx)
	}

	return "", nil
}

func (m *MockTemplateRenderer) RenderInstructionPointer(
	ctx types.TemplateContext,
) (string, error) {
	if m.RenderInstructionPointerFunc != nil {
		return m.RenderInstructionPointerFunc(ctx)
	}

	return "", nil
}

func (m *MockTemplateRenderer) RenderSlashCommand(
	cmd string,
	ctx types.TemplateContext,
) (string, error) {
	if m.RenderSlashCommandFunc != nil {
		return m.RenderSlashCommandFunc(cmd, ctx)
	}

	return "", nil
}
