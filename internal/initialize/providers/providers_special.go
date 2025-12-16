// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file implements Antigravity, Codex, CoStrict, and Crush providers -
// tools with special configurations or instruction file requirements.
package providers

import (
	"context"
)

// Compile-time interface satisfaction checks.
var (
	_ NewProvider = (*AntigravityNewProvider)(nil)
	_ NewProvider = (*CodexNewProvider)(nil)
	_ NewProvider = (*CostrictNewProvider)(nil)
	_ NewProvider = (*CrushNewProvider)(nil)
)

// =============================================================================
// Antigravity Provider
// =============================================================================

// AntigravityNewProvider implements the NewProvider interface for
// Antigravity. It uses AGENTS.md for instructions and .agent/workflows/
// for slash commands with a prefixed naming pattern.
type AntigravityNewProvider struct {
	renderer TemplateRenderer
	factory  InitializerFactory
}

// NewAntigravityNewProvider creates a new Antigravity provider.
func NewAntigravityNewProvider(
	renderer TemplateRenderer,
	factory InitializerFactory,
) *AntigravityNewProvider {
	return &AntigravityNewProvider{
		renderer: renderer,
		factory:  factory,
	}
}

// Initializers returns the list of initializers needed to configure
// Antigravity.
func (p *AntigravityNewProvider) Initializers(
	_ context.Context,
) []Initializer {
	templateCtx := DefaultTemplateContext()
	instructionContent, err := p.renderer.RenderInstructionPointer(
		templateCtx,
	)
	if err != nil {
		instructionContent = ""
	}

	return []Initializer{
		p.factory.CreateDirectoryInitializer(
			".agent/workflows",
		),
		p.factory.CreateConfigFileInitializer(
			"AGENTS.md",
			instructionContent,
		),
		p.factory.CreateSlashCommandsInitializer(
			".agent/workflows",
			extMD,
			FormatMarkdown,
			p.renderer,
		),
	}
}

// RegisterAntigravityProvider registers the Antigravity provider with the
// given registry.
func RegisterAntigravityProvider(
	reg *ProviderRegistry,
	renderer TemplateRenderer,
	factory InitializerFactory,
) error {
	return reg.Register(Registration{
		ID:       "antigravity",
		Name:     "Antigravity",
		Priority: PriorityAntigravity,
		Provider: NewAntigravityNewProvider(
			renderer,
			factory,
		),
	})
}

// =============================================================================
// Codex Provider
// =============================================================================

// CodexNewProvider implements the NewProvider interface for Codex CLI.
// Codex uses AGENTS.md for instructions and global ~/.codex/prompts/
// for slash commands. Note: Codex uses global paths rather than
// project-relative paths.
type CodexNewProvider struct {
	renderer TemplateRenderer
	factory  InitializerFactory
}

// NewCodexNewProvider creates a new Codex CLI provider.
func NewCodexNewProvider(
	renderer TemplateRenderer,
	factory InitializerFactory,
) *CodexNewProvider {
	return &CodexNewProvider{
		renderer: renderer,
		factory:  factory,
	}
}

// Initializers returns the list of initializers needed to configure Codex.
// Note: Codex uses global paths (~/...) which the initializers need to
// handle specially.
func (p *CodexNewProvider) Initializers(
	_ context.Context,
) []Initializer {
	templateCtx := DefaultTemplateContext()
	instructionContent, err := p.renderer.RenderInstructionPointer(
		templateCtx,
	)
	if err != nil {
		instructionContent = ""
	}

	return []Initializer{
		p.factory.CreateDirectoryInitializer(
			"~/.codex/prompts",
		),
		p.factory.CreateConfigFileInitializer(
			"AGENTS.md",
			instructionContent,
		),
		p.factory.CreateSlashCommandsInitializer(
			"~/.codex/prompts",
			extMD,
			FormatMarkdown,
			p.renderer,
		),
	}
}

// RegisterCodexProvider registers the Codex CLI provider with the given
// registry.
func RegisterCodexProvider(
	reg *ProviderRegistry,
	renderer TemplateRenderer,
	factory InitializerFactory,
) error {
	return reg.Register(Registration{
		ID:       "codex",
		Name:     "Codex CLI",
		Priority: PriorityCodex,
		Provider: NewCodexNewProvider(
			renderer,
			factory,
		),
	})
}

// =============================================================================
// CoStrict Provider
// =============================================================================

// CostrictNewProvider implements the NewProvider interface for CoStrict.
// It uses COSTRICT.md for instructions and .costrict/commands/spectr/
// for slash commands.
type CostrictNewProvider struct {
	renderer TemplateRenderer
	factory  InitializerFactory
}

// NewCostrictNewProvider creates a new CoStrict provider.
func NewCostrictNewProvider(
	renderer TemplateRenderer,
	factory InitializerFactory,
) *CostrictNewProvider {
	return &CostrictNewProvider{
		renderer: renderer,
		factory:  factory,
	}
}

// Initializers returns the list of initializers needed to configure CoStrict.
func (p *CostrictNewProvider) Initializers(
	_ context.Context,
) []Initializer {
	templateCtx := DefaultTemplateContext()
	instructionContent, err := p.renderer.RenderInstructionPointer(
		templateCtx,
	)
	if err != nil {
		instructionContent = ""
	}

	return []Initializer{
		p.factory.CreateDirectoryInitializer(
			".costrict/commands/spectr",
		),
		p.factory.CreateConfigFileInitializer(
			"COSTRICT.md",
			instructionContent,
		),
		p.factory.CreateSlashCommandsInitializer(
			".costrict/commands/spectr",
			extMD,
			FormatMarkdown,
			p.renderer,
		),
	}
}

// RegisterCostrictProvider registers the CoStrict provider with the given
// registry.
func RegisterCostrictProvider(
	reg *ProviderRegistry,
	renderer TemplateRenderer,
	factory InitializerFactory,
) error {
	return reg.Register(Registration{
		ID:       "costrict",
		Name:     "CoStrict",
		Priority: PriorityCostrict,
		Provider: NewCostrictNewProvider(
			renderer,
			factory,
		),
	})
}

// =============================================================================
// Crush Provider
// =============================================================================

// CrushNewProvider implements the NewProvider interface for Crush. It uses
// CRUSH.md for instructions and .crush/commands/spectr/ for slash commands.
type CrushNewProvider struct {
	renderer TemplateRenderer
	factory  InitializerFactory
}

// NewCrushNewProvider creates a new Crush provider.
func NewCrushNewProvider(
	renderer TemplateRenderer,
	factory InitializerFactory,
) *CrushNewProvider {
	return &CrushNewProvider{
		renderer: renderer,
		factory:  factory,
	}
}

// Initializers returns the list of initializers needed to configure Crush.
func (p *CrushNewProvider) Initializers(
	_ context.Context,
) []Initializer {
	templateCtx := DefaultTemplateContext()
	instructionContent, err := p.renderer.RenderInstructionPointer(
		templateCtx,
	)
	if err != nil {
		instructionContent = ""
	}

	return []Initializer{
		p.factory.CreateDirectoryInitializer(
			".crush/commands/spectr",
		),
		p.factory.CreateConfigFileInitializer(
			"CRUSH.md",
			instructionContent,
		),
		p.factory.CreateSlashCommandsInitializer(
			".crush/commands/spectr",
			extMD,
			FormatMarkdown,
			p.renderer,
		),
	}
}

// RegisterCrushProvider registers the Crush provider with the given registry.
func RegisterCrushProvider(
	reg *ProviderRegistry,
	renderer TemplateRenderer,
	factory InitializerFactory,
) error {
	return reg.Register(Registration{
		ID:       "crush",
		Name:     "Crush",
		Priority: PriorityCrush,
		Provider: NewCrushNewProvider(
			renderer,
			factory,
		),
	})
}
