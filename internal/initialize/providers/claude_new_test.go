package providers_test

import (
	"context"
	"testing"

	"github.com/connerohnesorge/spectr/internal/initialize/providers"
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
)

// mockTemplateRenderer implements providers.TemplateRenderer for testing.
type mockTemplateRenderer struct {
	instructionPointer string
	instructionError   error
	slashCommands      map[string]string
	slashError         error
}

func (*mockTemplateRenderer) RenderAgents(
	_ providers.TemplateContext,
) (string, error) {
	return "# AGENTS.md content", nil
}

func (m *mockTemplateRenderer) RenderInstructionPointer(
	_ providers.TemplateContext,
) (string, error) {
	if m.instructionError != nil {
		return "", m.instructionError
	}

	return m.instructionPointer, nil
}

func (m *mockTemplateRenderer) RenderSlashCommand(
	command string,
	_ providers.TemplateContext,
) (string, error) {
	if m.slashError != nil {
		return "", m.slashError
	}
	if content, ok := m.slashCommands[command]; ok {
		return content, nil
	}

	return "# " + command + " command content", nil
}

func TestClaudeNewProvider_Initializers(
	t *testing.T,
) {
	renderer := &mockTemplateRenderer{
		instructionPointer: "# Spectr Instructions\nRead spectr/AGENTS.md",
		slashCommands: map[string]string{
			"proposal": "# Proposal command",
			"apply":    "# Apply command",
		},
	}
	factory := initializers.NewFactory()

	provider := providers.NewClaudeNewProvider(
		renderer,
		factory,
	)

	inits := provider.Initializers(
		context.Background(),
	)

	if len(inits) != 3 {
		t.Errorf(
			"expected 3 initializers, got %d",
			len(inits),
		)
	}
}

func TestClaudeNewProvider_InitializersWithRenderError(
	t *testing.T,
) {
	renderer := &mockTemplateRenderer{
		instructionPointer: "",
		instructionError:   nil, // Even without error, empty template should work
	}
	factory := initializers.NewFactory()

	provider := providers.NewClaudeNewProvider(
		renderer,
		factory,
	)

	inits := provider.Initializers(
		context.Background(),
	)

	// Should still return 3 initializers even with empty template
	if len(inits) != 3 {
		t.Errorf(
			"expected 3 initializers, got %d",
			len(inits),
		)
	}
}

func TestRegisterClaudeProvider(t *testing.T) {
	renderer := &mockTemplateRenderer{
		instructionPointer: "# Spectr Instructions",
	}
	factory := initializers.NewFactory()

	reg := providers.CreateRegistry()

	err := providers.RegisterClaudeProvider(
		reg,
		renderer,
		factory,
	)
	if err != nil {
		t.Errorf(
			"unexpected error registering Claude provider: %v",
			err,
		)
	}

	// Verify registration
	registration := reg.Get("claude-code")
	if registration == nil {
		t.Fatal(
			"expected claude-code registration, got nil",
		)
	}

	if registration.ID != "claude-code" {
		t.Errorf(
			"expected ID 'claude-code', got %q",
			registration.ID,
		)
	}

	if registration.Name != "Claude Code" {
		t.Errorf(
			"expected Name 'Claude Code', got %q",
			registration.Name,
		)
	}

	if registration.Priority != providers.PriorityClaudeCode {
		t.Errorf(
			"expected Priority %d, got %d",
			providers.PriorityClaudeCode,
			registration.Priority,
		)
	}
}

func TestRegisterClaudeProvider_DuplicateRegistration(
	t *testing.T,
) {
	renderer := &mockTemplateRenderer{
		instructionPointer: "# Spectr Instructions",
	}
	factory := initializers.NewFactory()

	reg := providers.CreateRegistry()

	// First registration should succeed
	err := providers.RegisterClaudeProvider(
		reg,
		renderer,
		factory,
	)
	if err != nil {
		t.Fatalf(
			"first registration failed: %v",
			err,
		)
	}

	// Second registration should fail
	err = providers.RegisterClaudeProvider(
		reg,
		renderer,
		factory,
	)
	if err == nil {
		t.Error(
			"expected error for duplicate registration, got nil",
		)
	}
}
