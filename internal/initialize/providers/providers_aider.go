// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file implements Aider, Continue, and Tabnine providers - tools that use
// slash commands but no instruction files.
package providers

import (
	"context"
)

// Compile-time interface satisfaction checks.
var (
	_ NewProvider = (*AiderNewProvider)(nil)
	_ NewProvider = (*ContinueNewProvider)(nil)
	_ NewProvider = (*TabnineNewProvider)(nil)
)

// =============================================================================
// Aider Provider
// =============================================================================

// AiderNewProvider implements the NewProvider interface for Aider.
// Aider uses .aider/commands/spectr/ for slash commands in Markdown format.
// It has no instruction file.
type AiderNewProvider struct {
	renderer TemplateRenderer
	factory  InitializerFactory
}

// NewAiderNewProvider creates a new Aider provider with the given renderer
// and factory.
func NewAiderNewProvider(
	renderer TemplateRenderer,
	factory InitializerFactory,
) *AiderNewProvider {
	return &AiderNewProvider{
		renderer: renderer,
		factory:  factory,
	}
}

// Initializers returns the list of initializers needed to configure Aider.
func (p *AiderNewProvider) Initializers(
	_ context.Context,
) []Initializer {
	return []Initializer{
		p.factory.CreateDirectoryInitializer(
			".aider/commands/spectr",
		),
		p.factory.CreateSlashCommandsInitializer(
			".aider/commands/spectr",
			extMD,
			FormatMarkdown,
			p.renderer,
		),
	}
}

// RegisterAiderProvider registers the Aider provider with the given registry.
func RegisterAiderProvider(
	reg *ProviderRegistry,
	renderer TemplateRenderer,
	factory InitializerFactory,
) error {
	return reg.Register(Registration{
		ID:       "aider",
		Name:     "Aider",
		Priority: PriorityAider,
		Provider: NewAiderNewProvider(
			renderer,
			factory,
		),
	})
}

// =============================================================================
// Continue Provider
// =============================================================================

// ContinueNewProvider implements the NewProvider interface for Continue.
// It uses .continue/commands/spectr/ for slash commands in Markdown
// format. It has no instruction file.
type ContinueNewProvider struct {
	renderer TemplateRenderer
	factory  InitializerFactory
}

// NewContinueNewProvider creates a new Continue provider.
func NewContinueNewProvider(
	renderer TemplateRenderer,
	factory InitializerFactory,
) *ContinueNewProvider {
	return &ContinueNewProvider{
		renderer: renderer,
		factory:  factory,
	}
}

// Initializers returns the list of initializers needed to configure Continue.
func (p *ContinueNewProvider) Initializers(
	_ context.Context,
) []Initializer {
	return []Initializer{
		p.factory.CreateDirectoryInitializer(
			".continue/commands/spectr",
		),
		p.factory.CreateSlashCommandsInitializer(
			".continue/commands/spectr",
			extMD,
			FormatMarkdown,
			p.renderer,
		),
	}
}

// RegisterContinueProvider registers the Continue provider with the given
// registry.
func RegisterContinueProvider(
	reg *ProviderRegistry,
	renderer TemplateRenderer,
	factory InitializerFactory,
) error {
	return reg.Register(Registration{
		ID:       "continue",
		Name:     "Continue",
		Priority: PriorityContinue,
		Provider: NewContinueNewProvider(
			renderer,
			factory,
		),
	})
}

// =============================================================================
// Tabnine Provider
// =============================================================================

// TabnineNewProvider implements the NewProvider interface for Tabnine.
// Tabnine uses .tabnine/commands/spectr/ for slash commands in Markdown format.
// It has no instruction file.
type TabnineNewProvider struct {
	renderer TemplateRenderer
	factory  InitializerFactory
}

// NewTabnineNewProvider creates a new Tabnine provider.
func NewTabnineNewProvider(
	renderer TemplateRenderer,
	factory InitializerFactory,
) *TabnineNewProvider {
	return &TabnineNewProvider{
		renderer: renderer,
		factory:  factory,
	}
}

// Initializers returns the list of initializers needed to configure Tabnine.
func (p *TabnineNewProvider) Initializers(
	_ context.Context,
) []Initializer {
	return []Initializer{
		p.factory.CreateDirectoryInitializer(
			".tabnine/commands/spectr",
		),
		p.factory.CreateSlashCommandsInitializer(
			".tabnine/commands/spectr",
			extMD,
			FormatMarkdown,
			p.renderer,
		),
	}
}

// RegisterTabnineProvider registers the Tabnine provider with the given
// registry.
func RegisterTabnineProvider(
	reg *ProviderRegistry,
	renderer TemplateRenderer,
	factory InitializerFactory,
) error {
	return reg.Register(Registration{
		ID:       "tabnine",
		Name:     "Tabnine",
		Priority: PriorityTabnine,
		Provider: NewTabnineNewProvider(
			renderer,
			factory,
		),
	})
}
