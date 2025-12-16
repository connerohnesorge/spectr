// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file contains the InitializerFactory interface and registration helpers.
// Provider implementations are split across separate files by category.
package providers

// InitializerFactory is a generic factory interface for creating initializers.
// This consolidates the ClaudeInitializerFactory and GeminiInitializerFactory
// interfaces into a single interface that can be used by all providers.
type InitializerFactory interface {
	// CreateDirectoryInitializer creates a DirectoryInitializer for the
	// given paths.
	CreateDirectoryInitializer(
		paths ...string,
	) Initializer

	// CreateConfigFileInitializer creates a ConfigFileInitializer for the
	// given path and template.
	CreateConfigFileInitializer(
		path, template string,
	) Initializer

	// CreateSlashCommandsInitializer creates a SlashCommandsInitializer.
	CreateSlashCommandsInitializer(
		dir, ext string,
		format CommandFormat,
		renderer TemplateRenderer,
	) Initializer
}

// File extension constant.
const extMD = ".md"

// =============================================================================
// Registration Helper Functions
// =============================================================================

// RegisterAllProviders registers all providers with the given registry.
// This is a convenience function for registering all providers at once.
//
// The function registers providers in the following order (by priority):
//   - Claude Code (1) - registered separately via RegisterClaudeProvider
//   - Gemini CLI (2) - registered separately via RegisterGeminiProvider
//   - CoStrict (3)
//   - Qoder (4)
//   - CodeBuddy (5)
//   - Qwen Code (6)
//   - Antigravity (7)
//   - Cline (8)
//   - Cursor (9)
//   - Codex CLI (10)
//   - Aider (11)
//   - Tabnine (12)
//   - Windsurf (13)
//   - Kilocode (14)
//   - Continue (15)
//   - Crush (16)
//   - OpenCode (17)
//
// Note: This function does NOT register Claude and Gemini providers. Use
// RegisterClaudeProvider and RegisterGeminiProvider separately for those.
func RegisterAllProviders(
	reg *ProviderRegistry,
	renderer TemplateRenderer,
	factory InitializerFactory,
) error {
	type regFunc func(
		*ProviderRegistry, TemplateRenderer, InitializerFactory,
	) error
	registrations := []regFunc{
		RegisterCostrictProvider,
		RegisterQoderProvider,
		RegisterCodeBuddyProvider,
		RegisterQwenProvider,
		RegisterAntigravityProvider,
		RegisterClineProvider,
		RegisterCursorProvider,
		RegisterCodexProvider,
		RegisterAiderProvider,
		RegisterTabnineProvider,
		RegisterWindsurfProvider,
		RegisterKilocodeProvider,
		RegisterContinueProvider,
		RegisterCrushProvider,
		RegisterOpencodeProvider,
	}

	for _, register := range registrations {
		if err := register(reg, renderer, factory); err != nil {
			return err
		}
	}

	return nil
}

// RegisterAllProvidersIncludingBase registers all providers including
// Claude and Gemini. This is the most complete registration function
// that sets up all 17 providers.
//
// Note: This function requires the factory to implement
// InitializerFactory, which is a superset of ClaudeInitializerFactory
// and GeminiInitializerFactory.
func RegisterAllProvidersIncludingBase(
	reg *ProviderRegistry,
	renderer TemplateRenderer,
	factory InitializerFactory,
) error {
	// Register Claude (uses ClaudeInitializerFactory interface, which
	// InitializerFactory satisfies)
	if err := RegisterClaudeProvider(reg, renderer, factory); err != nil {
		return err
	}

	// Register Gemini (uses GeminiInitializerFactory interface, which
	// InitializerFactory satisfies)
	if err := RegisterGeminiProvider(reg, renderer, factory); err != nil {
		return err
	}

	// Register all other providers
	return RegisterAllProviders(
		reg,
		renderer,
		factory,
	)
}
