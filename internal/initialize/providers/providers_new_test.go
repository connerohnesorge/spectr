package providers_test

import (
	"context"
	"testing"

	"github.com/connerohnesorge/spectr/internal/initialize/providers"
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
)

// mockRenderer implements providers.TemplateRenderer for testing.
// This is a simplified mock that returns predictable content.
type mockRenderer struct{}

func (*mockRenderer) RenderAgents(
	_ providers.TemplateContext,
) (string, error) {
	return "mock agents content", nil
}

func (*mockRenderer) RenderInstructionPointer(
	_ providers.TemplateContext,
) (string, error) {
	return "mock instruction content", nil
}

func (*mockRenderer) RenderSlashCommand(
	command string,
	_ providers.TemplateContext,
) (string, error) {
	return "mock command content for " + command, nil
}

// TestClaudeProviderInitializers verifies Claude returns the correct initializers.
// Claude should return 3 initializers: Directory, ConfigFile, SlashCommands.
func TestClaudeProviderInitializers(
	t *testing.T,
) {
	renderer := &mockRenderer{}
	factory := initializers.NewFactory()

	provider := providers.NewClaudeNewProvider(
		renderer,
		factory,
	)
	inits := provider.Initializers(
		context.Background(),
	)

	// Verify count
	if len(inits) != 3 {
		t.Fatalf(
			"expected 3 initializers for Claude, got %d",
			len(inits),
		)
	}

	// Verify types using type assertions
	// First should be DirectoryInitializer
	if _, ok := inits[0].(*initializers.DirectoryInitializer); !ok {
		t.Errorf(
			"initializer[0] should be *DirectoryInitializer, got %T",
			inits[0],
		)
	}

	// Second should be ConfigFileInitializer
	if _, ok := inits[1].(*initializers.ConfigFileInitializer); !ok {
		t.Errorf(
			"initializer[1] should be *ConfigFileInitializer, got %T",
			inits[1],
		)
	}

	// Third should be SlashCommandsInitializer
	if _, ok := inits[2].(*initializers.SlashCommandsInitializer); !ok {
		t.Errorf(
			"initializer[2] should be *SlashCommandsInitializer, got %T",
			inits[2],
		)
	}
}

// TestGeminiProviderInitializers verifies Gemini returns the correct initializers.
// Gemini should return 2 initializers: Directory, SlashCommands (no ConfigFile).
func TestGeminiProviderInitializers(
	t *testing.T,
) {
	renderer := &mockRenderer{}
	factory := initializers.NewFactory()

	provider := providers.NewGeminiNewProvider(
		renderer,
		factory,
	)
	inits := provider.Initializers(
		context.Background(),
	)

	// Verify count
	if len(inits) != 2 {
		t.Fatalf(
			"expected 2 initializers for Gemini (no config file), got %d",
			len(inits),
		)
	}

	// Verify types using type assertions
	// First should be DirectoryInitializer
	if _, ok := inits[0].(*initializers.DirectoryInitializer); !ok {
		t.Errorf(
			"initializer[0] should be *DirectoryInitializer, got %T",
			inits[0],
		)
	}

	// Second should be SlashCommandsInitializer (no ConfigFile)
	if _, ok := inits[1].(*initializers.SlashCommandsInitializer); !ok {
		t.Errorf(
			"initializer[1] should be *SlashCommandsInitializer, got %T",
			inits[1],
		)
	}

	// Verify there is NO ConfigFileInitializer
	for i, init := range inits {
		if _, ok := init.(*initializers.ConfigFileInitializer); ok {
			t.Errorf(
				"Gemini should NOT have ConfigFileInitializer, but found one at index %d",
				i,
			)
		}
	}
}

// TestProvidersWithConfigFile tests all providers that should have 3 initializers (Dir, Config, Slash).
func TestProvidersWithConfigFile(t *testing.T) {
	renderer := &mockRenderer{}
	factory := initializers.NewFactory()

	// Providers that have config files (3 initializers each)
	providersWithConfig := []struct {
		name     string
		provider providers.NewProvider
	}{
		{
			"Antigravity",
			providers.NewAntigravityNewProvider(
				renderer,
				factory,
			),
		},
		{
			"Cline",
			providers.NewClineNewProvider(
				renderer,
				factory,
			),
		},
		{
			"CodeBuddy",
			providers.NewCodeBuddyNewProvider(
				renderer,
				factory,
			),
		},
		{
			"Codex",
			providers.NewCodexNewProvider(
				renderer,
				factory,
			),
		},
		{
			"CoStrict",
			providers.NewCostrictNewProvider(
				renderer,
				factory,
			),
		},
		{
			"Crush",
			providers.NewCrushNewProvider(
				renderer,
				factory,
			),
		},
		{
			"Qoder",
			providers.NewQoderNewProvider(
				renderer,
				factory,
			),
		},
		{
			"Qwen",
			providers.NewQwenNewProvider(
				renderer,
				factory,
			),
		},
	}

	for _, tc := range providersWithConfig {
		t.Run(tc.name, func(t *testing.T) {
			inits := tc.provider.Initializers(
				context.Background(),
			)

			// Should have 3 initializers
			if len(inits) != 3 {
				t.Fatalf(
					"expected 3 initializers for %s (has config file), got %d",
					tc.name,
					len(inits),
				)
			}

			// Verify types in order: Directory, ConfigFile, SlashCommands
			if _, ok := inits[0].(*initializers.DirectoryInitializer); !ok {
				t.Errorf(
					"%s initializer[0] should be *DirectoryInitializer, got %T",
					tc.name,
					inits[0],
				)
			}
			if _, ok := inits[1].(*initializers.ConfigFileInitializer); !ok {
				t.Errorf(
					"%s initializer[1] should be *ConfigFileInitializer, got %T",
					tc.name,
					inits[1],
				)
			}
			if _, ok := inits[2].(*initializers.SlashCommandsInitializer); !ok {
				t.Errorf(
					"%s initializer[2] should be *SlashCommandsInitializer, got %T",
					tc.name,
					inits[2],
				)
			}
		})
	}
}

// TestProvidersWithoutConfigFile tests all providers that should have 2 initializers (Dir, Slash only).
func TestProvidersWithoutConfigFile(
	t *testing.T,
) {
	renderer := &mockRenderer{}
	factory := initializers.NewFactory()

	// Providers without config files (2 initializers each)
	providersWithoutConfig := []struct {
		name     string
		provider providers.NewProvider
	}{
		{
			"Aider",
			providers.NewAiderNewProvider(
				renderer,
				factory,
			),
		},
		{
			"Continue",
			providers.NewContinueNewProvider(
				renderer,
				factory,
			),
		},
		{
			"Cursor",
			providers.NewCursorNewProvider(
				renderer,
				factory,
			),
		},
		{
			"Kilocode",
			providers.NewKilocodeNewProvider(
				renderer,
				factory,
			),
		},
		{
			"OpenCode",
			providers.NewOpencodeNewProvider(
				renderer,
				factory,
			),
		},
		{
			"Tabnine",
			providers.NewTabnineNewProvider(
				renderer,
				factory,
			),
		},
		{
			"Windsurf",
			providers.NewWindsurfNewProvider(
				renderer,
				factory,
			),
		},
	}

	for _, tc := range providersWithoutConfig {
		t.Run(tc.name, func(t *testing.T) {
			inits := tc.provider.Initializers(
				context.Background(),
			)

			// Should have 2 initializers (no config file)
			if len(inits) != 2 {
				t.Fatalf(
					"expected 2 initializers for %s (no config file), got %d",
					tc.name,
					len(inits),
				)
			}

			// Verify types in order: Directory, SlashCommands
			if _, ok := inits[0].(*initializers.DirectoryInitializer); !ok {
				t.Errorf(
					"%s initializer[0] should be *DirectoryInitializer, got %T",
					tc.name,
					inits[0],
				)
			}
			if _, ok := inits[1].(*initializers.SlashCommandsInitializer); !ok {
				t.Errorf(
					"%s initializer[1] should be *SlashCommandsInitializer, got %T",
					tc.name,
					inits[1],
				)
			}

			// Verify there is NO ConfigFileInitializer
			for i, init := range inits {
				if _, ok := init.(*initializers.ConfigFileInitializer); ok {
					t.Errorf(
						"%s should NOT have ConfigFileInitializer, but found one at index %d",
						tc.name,
						i,
					)
				}
			}
		})
	}
}

// TestRegisterAllProviders verifies that RegisterAllProviders registers exactly 15 providers.
// This excludes Claude and Gemini which are registered separately.
func TestRegisterAllProviders(t *testing.T) {
	renderer := &mockRenderer{}
	factory := initializers.NewFactory()
	reg := providers.CreateRegistry()

	err := providers.RegisterAllProviders(
		reg,
		renderer,
		factory,
	)
	if err != nil {
		t.Fatalf(
			"RegisterAllProviders failed: %v",
			err,
		)
	}

	// Should have 15 providers (excludes Claude and Gemini)
	expectedCount := 15
	if reg.Count() != expectedCount {
		t.Errorf(
			"expected %d providers from RegisterAllProviders, got %d",
			expectedCount,
			reg.Count(),
		)
	}

	// Verify specific providers are registered
	expectedProviders := []string{
		"aider",
		"antigravity",
		"cline",
		"codebuddy",
		"codex",
		"continue",
		"costrict",
		"crush",
		"cursor",
		"kilocode",
		"opencode",
		"qoder",
		"qwen",
		"tabnine",
		"windsurf",
	}

	for _, id := range expectedProviders {
		if reg.Get(id) == nil {
			t.Errorf(
				"expected provider %q to be registered, but it was not",
				id,
			)
		}
	}

	// Verify Claude and Gemini are NOT registered
	if reg.Get("claude-code") != nil {
		t.Error(
			"claude-code should NOT be registered by RegisterAllProviders",
		)
	}
	if reg.Get("gemini") != nil {
		t.Error(
			"gemini should NOT be registered by RegisterAllProviders",
		)
	}
}

// TestRegisterAllProvidersIncludingBase verifies that RegisterAllProvidersIncludingBase
// registers all 17 providers including Claude and Gemini.
func TestRegisterAllProvidersIncludingBase(
	t *testing.T,
) {
	renderer := &mockRenderer{}
	factory := initializers.NewFactory()
	reg := providers.CreateRegistry()

	err := providers.RegisterAllProvidersIncludingBase(
		reg,
		renderer,
		factory,
	)
	if err != nil {
		t.Fatalf(
			"RegisterAllProvidersIncludingBase failed: %v",
			err,
		)
	}

	// Should have 17 providers total
	expectedCount := 17
	if reg.Count() != expectedCount {
		t.Errorf(
			"expected %d providers from RegisterAllProvidersIncludingBase, got %d",
			expectedCount,
			reg.Count(),
		)
	}

	// Verify Claude and Gemini ARE registered
	if reg.Get("claude-code") == nil {
		t.Error(
			"claude-code should be registered by RegisterAllProvidersIncludingBase",
		)
	}
	if reg.Get("gemini") == nil {
		t.Error(
			"gemini should be registered by RegisterAllProvidersIncludingBase",
		)
	}

	// Verify all 17 providers
	allProviders := []string{
		"claude-code",
		"gemini",
		"aider",
		"antigravity",
		"cline",
		"codebuddy",
		"codex",
		"continue",
		"costrict",
		"crush",
		"cursor",
		"kilocode",
		"opencode",
		"qoder",
		"qwen",
		"tabnine",
		"windsurf",
	}

	for _, id := range allProviders {
		if reg.Get(id) == nil {
			t.Errorf(
				"expected provider %q to be registered, but it was not",
				id,
			)
		}
	}
}

// TestProviderInitializerKeys verifies that each provider's initializers have unique keys.
// Keys are used for deduplication when multiple providers have overlapping initializers.
func TestProviderInitializerKeys(t *testing.T) {
	renderer := &mockRenderer{}
	factory := initializers.NewFactory()

	// Test a few providers to ensure keys are unique within each provider
	testCases := []struct {
		name     string
		provider providers.NewProvider
	}{
		{
			"Claude",
			providers.NewClaudeNewProvider(
				renderer,
				factory,
			),
		},
		{
			"Gemini",
			providers.NewGeminiNewProvider(
				renderer,
				factory,
			),
		},
		{
			"Aider",
			providers.NewAiderNewProvider(
				renderer,
				factory,
			),
		},
		{
			"Cline",
			providers.NewClineNewProvider(
				renderer,
				factory,
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inits := tc.provider.Initializers(
				context.Background(),
			)
			keys := make(map[string]bool)

			for i, init := range inits {
				// Get key based on concrete type
				var key string
				switch v := init.(type) {
				case *initializers.DirectoryInitializer:
					key = v.Key()
				case *initializers.ConfigFileInitializer:
					key = v.Key()
				case *initializers.SlashCommandsInitializer:
					key = v.Key()
				default:
					t.Errorf(
						"%s initializer[%d] has unknown type %T",
						tc.name,
						i,
						init,
					)

					continue
				}

				if key == "" {
					t.Errorf(
						"%s initializer[%d] has empty key",
						tc.name,
						i,
					)
				}
				if keys[key] {
					t.Errorf(
						"%s has duplicate initializer key: %q",
						tc.name,
						key,
					)
				}
				keys[key] = true
			}
		})
	}
}

// TestProviderPriorities verifies that providers are registered with correct priorities.
func TestProviderPriorities(t *testing.T) {
	renderer := &mockRenderer{}
	factory := initializers.NewFactory()
	reg := providers.CreateRegistry()

	err := providers.RegisterAllProvidersIncludingBase(
		reg,
		renderer,
		factory,
	)
	if err != nil {
		t.Fatalf("registration failed: %v", err)
	}

	// Verify some key priorities
	expectedPriorities := map[string]int{
		"claude-code": providers.PriorityClaudeCode,
		"gemini":      providers.PriorityGemini,
		"costrict":    providers.PriorityCostrict,
		"cursor":      providers.PriorityCursor,
		"cline":       providers.PriorityCline,
	}

	for id, expectedPriority := range expectedPriorities {
		registration := reg.Get(id)
		if registration == nil {
			t.Errorf("provider %q not found", id)

			continue
		}
		if registration.Priority != expectedPriority {
			t.Errorf(
				"provider %q has priority %d, expected %d",
				id,
				registration.Priority,
				expectedPriority,
			)
		}
	}

	// Verify All() returns in priority order
	all := reg.All()
	for i := 1; i < len(all); i++ {
		if all[i-1].Priority > all[i].Priority {
			t.Errorf(
				"providers not sorted by priority: %s (priority %d) comes before %s (priority %d)",
				all[i-1].ID,
				all[i-1].Priority,
				all[i].ID,
				all[i].Priority,
			)
		}
	}
}

// TestAllProvidersReturnNonNilInitializers verifies that all providers return non-nil initializers.
func TestAllProvidersReturnNonNilInitializers(
	t *testing.T,
) {
	renderer := &mockRenderer{}
	factory := initializers.NewFactory()

	allProviders := []struct {
		name     string
		provider providers.NewProvider
	}{
		{
			"Claude",
			providers.NewClaudeNewProvider(
				renderer,
				factory,
			),
		},
		{
			"Gemini",
			providers.NewGeminiNewProvider(
				renderer,
				factory,
			),
		},
		{
			"Aider",
			providers.NewAiderNewProvider(
				renderer,
				factory,
			),
		},
		{
			"Antigravity",
			providers.NewAntigravityNewProvider(
				renderer,
				factory,
			),
		},
		{
			"Cline",
			providers.NewClineNewProvider(
				renderer,
				factory,
			),
		},
		{
			"CodeBuddy",
			providers.NewCodeBuddyNewProvider(
				renderer,
				factory,
			),
		},
		{
			"Codex",
			providers.NewCodexNewProvider(
				renderer,
				factory,
			),
		},
		{
			"Continue",
			providers.NewContinueNewProvider(
				renderer,
				factory,
			),
		},
		{
			"CoStrict",
			providers.NewCostrictNewProvider(
				renderer,
				factory,
			),
		},
		{
			"Crush",
			providers.NewCrushNewProvider(
				renderer,
				factory,
			),
		},
		{
			"Cursor",
			providers.NewCursorNewProvider(
				renderer,
				factory,
			),
		},
		{
			"Kilocode",
			providers.NewKilocodeNewProvider(
				renderer,
				factory,
			),
		},
		{
			"OpenCode",
			providers.NewOpencodeNewProvider(
				renderer,
				factory,
			),
		},
		{
			"Qoder",
			providers.NewQoderNewProvider(
				renderer,
				factory,
			),
		},
		{
			"Qwen",
			providers.NewQwenNewProvider(
				renderer,
				factory,
			),
		},
		{
			"Tabnine",
			providers.NewTabnineNewProvider(
				renderer,
				factory,
			),
		},
		{
			"Windsurf",
			providers.NewWindsurfNewProvider(
				renderer,
				factory,
			),
		},
	}

	for _, tc := range allProviders {
		t.Run(tc.name, func(t *testing.T) {
			inits := tc.provider.Initializers(
				context.Background(),
			)

			if inits == nil {
				t.Errorf(
					"%s returned nil initializers slice",
					tc.name,
				)

				return
			}

			if len(inits) == 0 {
				t.Errorf(
					"%s returned empty initializers slice",
					tc.name,
				)

				return
			}

			for i, init := range inits {
				if init == nil {
					t.Errorf(
						"%s initializer[%d] is nil",
						tc.name,
						i,
					)
				}
			}
		})
	}
}

// TestDuplicateRegistrationFails verifies that registering the same provider twice fails.
func TestDuplicateRegistrationFails(
	t *testing.T,
) {
	renderer := &mockRenderer{}
	factory := initializers.NewFactory()
	reg := providers.CreateRegistry()

	// Register all providers
	err := providers.RegisterAllProvidersIncludingBase(
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

	// Try to register again - should fail
	err = providers.RegisterAllProvidersIncludingBase(
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
