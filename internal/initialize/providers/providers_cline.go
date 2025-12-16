// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file implements Cline, CodeBuddy, Qoder, and Qwen providers - tools that
// use instruction files with slash commands.
package providers

import (
	"context"
)

// Compile-time interface satisfaction checks.
var (
	_ NewProvider = (*ClineNewProvider)(nil)
	_ NewProvider = (*CodeBuddyNewProvider)(nil)
	_ NewProvider = (*QoderNewProvider)(nil)
	_ NewProvider = (*QwenNewProvider)(nil)
)

// =============================================================================
// Cline Provider
// =============================================================================

// ClineNewProvider implements the NewProvider interface for Cline. Cline
// uses CLINE.md for instructions and .clinerules/commands/spectr/ for
// slash commands.
type ClineNewProvider struct {
	renderer TemplateRenderer
	factory  InitializerFactory
}

// NewClineNewProvider creates a new Cline provider.
func NewClineNewProvider(
	renderer TemplateRenderer,
	factory InitializerFactory,
) *ClineNewProvider {
	return &ClineNewProvider{
		renderer: renderer,
		factory:  factory,
	}
}

// Initializers returns the list of initializers needed to configure Cline.
func (p *ClineNewProvider) Initializers(
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
			".clinerules/commands/spectr",
		),
		p.factory.CreateConfigFileInitializer(
			"CLINE.md",
			instructionContent,
		),
		p.factory.CreateSlashCommandsInitializer(
			".clinerules/commands/spectr",
			extMD,
			FormatMarkdown,
			p.renderer,
		),
	}
}

// RegisterClineProvider registers the Cline provider with the given registry.
func RegisterClineProvider(
	reg *ProviderRegistry,
	renderer TemplateRenderer,
	factory InitializerFactory,
) error {
	return reg.Register(Registration{
		ID:       "cline",
		Name:     "Cline",
		Priority: PriorityCline,
		Provider: NewClineNewProvider(
			renderer,
			factory,
		),
	})
}

// =============================================================================
// CodeBuddy Provider
// =============================================================================

// CodeBuddyNewProvider implements the NewProvider interface for CodeBuddy.
// It uses CODEBUDDY.md for instructions and .codebuddy/commands/spectr/
// for slash commands.
type CodeBuddyNewProvider struct {
	renderer TemplateRenderer
	factory  InitializerFactory
}

// NewCodeBuddyNewProvider creates a new CodeBuddy provider.
func NewCodeBuddyNewProvider(
	renderer TemplateRenderer,
	factory InitializerFactory,
) *CodeBuddyNewProvider {
	return &CodeBuddyNewProvider{
		renderer: renderer,
		factory:  factory,
	}
}

// Initializers returns the list of initializers needed to configure CodeBuddy.
func (p *CodeBuddyNewProvider) Initializers(
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
			".codebuddy/commands/spectr",
		),
		p.factory.CreateConfigFileInitializer(
			"CODEBUDDY.md",
			instructionContent,
		),
		p.factory.CreateSlashCommandsInitializer(
			".codebuddy/commands/spectr",
			extMD,
			FormatMarkdown,
			p.renderer,
		),
	}
}

// RegisterCodeBuddyProvider registers the CodeBuddy provider with the given
// registry.
func RegisterCodeBuddyProvider(
	reg *ProviderRegistry,
	renderer TemplateRenderer,
	factory InitializerFactory,
) error {
	return reg.Register(Registration{
		ID:       "codebuddy",
		Name:     "CodeBuddy",
		Priority: PriorityCodeBuddy,
		Provider: NewCodeBuddyNewProvider(
			renderer,
			factory,
		),
	})
}

// =============================================================================
// Qoder Provider
// =============================================================================

// QoderNewProvider implements the NewProvider interface for Qoder. It
// uses QODER.md for instructions and .qoder/commands/spectr/ for slash
// commands.
type QoderNewProvider struct {
	renderer TemplateRenderer
	factory  InitializerFactory
}

// NewQoderNewProvider creates a new Qoder provider.
func NewQoderNewProvider(
	renderer TemplateRenderer,
	factory InitializerFactory,
) *QoderNewProvider {
	return &QoderNewProvider{
		renderer: renderer,
		factory:  factory,
	}
}

// Initializers returns the list of initializers needed to configure Qoder.
func (p *QoderNewProvider) Initializers(
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
			".qoder/commands/spectr",
		),
		p.factory.CreateConfigFileInitializer(
			"QODER.md",
			instructionContent,
		),
		p.factory.CreateSlashCommandsInitializer(
			".qoder/commands/spectr",
			extMD,
			FormatMarkdown,
			p.renderer,
		),
	}
}

// RegisterQoderProvider registers the Qoder provider with the given registry.
func RegisterQoderProvider(
	reg *ProviderRegistry,
	renderer TemplateRenderer,
	factory InitializerFactory,
) error {
	return reg.Register(Registration{
		ID:       "qoder",
		Name:     "Qoder",
		Priority: PriorityQoder,
		Provider: NewQoderNewProvider(
			renderer,
			factory,
		),
	})
}

// =============================================================================
// Qwen Provider
// =============================================================================

// QwenNewProvider implements the NewProvider interface for Qwen Code. It
// uses QWEN.md for instructions and .qwen/commands/spectr/ for slash
// commands.
type QwenNewProvider struct {
	renderer TemplateRenderer
	factory  InitializerFactory
}

// NewQwenNewProvider creates a new Qwen Code provider.
func NewQwenNewProvider(
	renderer TemplateRenderer,
	factory InitializerFactory,
) *QwenNewProvider {
	return &QwenNewProvider{
		renderer: renderer,
		factory:  factory,
	}
}

// Initializers returns the list of initializers needed to configure Qwen.
func (p *QwenNewProvider) Initializers(
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
			".qwen/commands/spectr",
		),
		p.factory.CreateConfigFileInitializer(
			"QWEN.md",
			instructionContent,
		),
		p.factory.CreateSlashCommandsInitializer(
			".qwen/commands/spectr",
			extMD,
			FormatMarkdown,
			p.renderer,
		),
	}
}

// RegisterQwenProvider registers the Qwen Code provider with the given
// registry.
func RegisterQwenProvider(
	reg *ProviderRegistry,
	renderer TemplateRenderer,
	factory InitializerFactory,
) error {
	return reg.Register(Registration{
		ID:       "qwen",
		Name:     "Qwen Code",
		Priority: PriorityQwen,
		Provider: NewQwenNewProvider(
			renderer,
			factory,
		),
	})
}
