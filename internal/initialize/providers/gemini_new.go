// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file implements the new Gemini CLI provider using the redesigned
// Provider interface with composable initializers.
package providers

import (
	"context"
)

// Compile-time interface satisfaction check.
var _ NewProvider = (*GeminiNewProvider)(nil)

// GeminiNewProvider implements the NewProvider interface for Gemini CLI.
// Gemini CLI uses .gemini/commands/spectr/ for slash commands in TOML format.
// Unlike Claude Code, Gemini does NOT have an instruction file (no GEMINI.md).
//
// This is the new implementation using the redesigned provider architecture.
// The old implementation is in gemini.go.
type GeminiNewProvider struct {
	// renderer is the template renderer for generating slash command content.
	renderer TemplateRenderer

	// initializerFactory creates the initializers for this provider. This
	// factory pattern avoids import cycles between providers and
	// initializers packages. We reuse ClaudeInitializerFactory since the
	// interface methods are the same.
	initializerFactory GeminiInitializerFactory
}

// GeminiInitializerFactory creates initializers for the Gemini provider.
// This interface allows the initializers to be created outside of the
// providers package, avoiding import cycles.
//
// Note: This interface is identical to ClaudeInitializerFactory but is
// defined separately to allow independent evolution and to maintain clear
// provider ownership. In a future refactor, these could be consolidated
// into a generic InitializerFactory.
type GeminiInitializerFactory interface {
	// CreateDirectoryInitializer creates a DirectoryInitializer for the
	// given paths.
	CreateDirectoryInitializer(
		paths ...string,
	) Initializer

	// CreateSlashCommandsInitializer creates a SlashCommandsInitializer.
	CreateSlashCommandsInitializer(
		dir, ext string,
		format CommandFormat,
		renderer TemplateRenderer,
	) Initializer
}

// NewGeminiNewProvider creates a new Gemini CLI provider with the given
// renderer and factory. The renderer is used to render slash command
// templates. The factory creates initializers without causing import cycles.
func NewGeminiNewProvider(
	renderer TemplateRenderer,
	factory GeminiInitializerFactory,
) *GeminiNewProvider {
	return &GeminiNewProvider{
		renderer:           renderer,
		initializerFactory: factory,
	}
}

// Initializers returns the list of initializers needed to configure Gemini CLI.
// Gemini CLI requires:
//   - Directory: .gemini/commands/spectr/
//   - Slash commands: proposal.toml and apply.toml in TOML format
//
// Unlike Claude Code, Gemini does NOT use an instruction file.
func (p *GeminiNewProvider) Initializers(
	_ context.Context,
) []Initializer {
	return []Initializer{
		// Create the slash commands directory
		p.initializerFactory.CreateDirectoryInitializer(
			".gemini/commands/spectr",
		),

		// Create/update the slash command files in TOML format
		p.initializerFactory.CreateSlashCommandsInitializer(
			".gemini/commands/spectr",
			".toml",
			FormatTOML,
			p.renderer,
		),
	}
}

// RegisterGeminiProvider registers the Gemini CLI provider with the given
// registry. This function should be called at application startup to make
// Gemini CLI available.
//
// Example:
//
//	reg := providers.CreateRegistry()
//	tm, _ := initialize.NewTemplateManager()
//	factory := initializers.NewGeminiFactory()
//	providers.RegisterGeminiProvider(reg, tm, factory)
func RegisterGeminiProvider(
	reg *ProviderRegistry,
	renderer TemplateRenderer,
	factory GeminiInitializerFactory,
) error {
	return reg.Register(Registration{
		ID:       "gemini",
		Name:     "Gemini CLI",
		Priority: PriorityGemini,
		Provider: NewGeminiNewProvider(
			renderer,
			factory,
		),
	})
}
