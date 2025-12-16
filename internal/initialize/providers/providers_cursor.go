// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file implements Cursor, Kilocode, OpenCode, and Windsurf providers -
// IDE extensions that use slash commands without instruction files.
package providers

import (
	"context"
)

// Compile-time interface satisfaction checks.
var (
	_ NewProvider = (*CursorNewProvider)(nil)
	_ NewProvider = (*KilocodeNewProvider)(nil)
	_ NewProvider = (*OpencodeNewProvider)(nil)
	_ NewProvider = (*WindsurfNewProvider)(nil)
)

// =============================================================================
// Cursor Provider
// =============================================================================

// CursorNewProvider implements the NewProvider interface for Cursor. It
// uses .cursorrules/commands/spectr/ for slash commands in Markdown
// format. It has no instruction file.
type CursorNewProvider struct {
	renderer TemplateRenderer
	factory  InitializerFactory
}

// NewCursorNewProvider creates a new Cursor provider.
func NewCursorNewProvider(
	renderer TemplateRenderer,
	factory InitializerFactory,
) *CursorNewProvider {
	return &CursorNewProvider{
		renderer: renderer,
		factory:  factory,
	}
}

// Initializers returns the list of initializers needed to configure Cursor.
func (p *CursorNewProvider) Initializers(
	_ context.Context,
) []Initializer {
	return []Initializer{
		p.factory.CreateDirectoryInitializer(
			".cursorrules/commands/spectr",
		),
		p.factory.CreateSlashCommandsInitializer(
			".cursorrules/commands/spectr",
			extMD,
			FormatMarkdown,
			p.renderer,
		),
	}
}

// RegisterCursorProvider registers the Cursor provider with the given registry.
func RegisterCursorProvider(
	reg *ProviderRegistry,
	renderer TemplateRenderer,
	factory InitializerFactory,
) error {
	return reg.Register(Registration{
		ID:       "cursor",
		Name:     "Cursor",
		Priority: PriorityCursor,
		Provider: NewCursorNewProvider(
			renderer,
			factory,
		),
	})
}

// =============================================================================
// Kilocode Provider
// =============================================================================

// KilocodeNewProvider implements the NewProvider interface for Kilocode.
// It uses .kilocode/commands/spectr/ for slash commands in Markdown
// format. It has no instruction file.
type KilocodeNewProvider struct {
	renderer TemplateRenderer
	factory  InitializerFactory
}

// NewKilocodeNewProvider creates a new Kilocode provider.
func NewKilocodeNewProvider(
	renderer TemplateRenderer,
	factory InitializerFactory,
) *KilocodeNewProvider {
	return &KilocodeNewProvider{
		renderer: renderer,
		factory:  factory,
	}
}

// Initializers returns the list of initializers needed to configure Kilocode.
func (p *KilocodeNewProvider) Initializers(
	_ context.Context,
) []Initializer {
	return []Initializer{
		p.factory.CreateDirectoryInitializer(
			".kilocode/commands/spectr",
		),
		p.factory.CreateSlashCommandsInitializer(
			".kilocode/commands/spectr",
			extMD,
			FormatMarkdown,
			p.renderer,
		),
	}
}

// RegisterKilocodeProvider registers the Kilocode provider with the given
// registry.
func RegisterKilocodeProvider(
	reg *ProviderRegistry,
	renderer TemplateRenderer,
	factory InitializerFactory,
) error {
	return reg.Register(Registration{
		ID:       "kilocode",
		Name:     "Kilocode",
		Priority: PriorityKilocode,
		Provider: NewKilocodeNewProvider(
			renderer,
			factory,
		),
	})
}

// =============================================================================
// OpenCode Provider
// =============================================================================

// OpencodeNewProvider implements the NewProvider interface for OpenCode.
// It uses .opencode/command/spectr/ for slash commands in Markdown
// format. It has no instruction file (uses JSON configuration).
type OpencodeNewProvider struct {
	renderer TemplateRenderer
	factory  InitializerFactory
}

// NewOpencodeNewProvider creates a new OpenCode provider.
func NewOpencodeNewProvider(
	renderer TemplateRenderer,
	factory InitializerFactory,
) *OpencodeNewProvider {
	return &OpencodeNewProvider{
		renderer: renderer,
		factory:  factory,
	}
}

// Initializers returns the list of initializers needed to configure OpenCode.
func (p *OpencodeNewProvider) Initializers(
	_ context.Context,
) []Initializer {
	return []Initializer{
		p.factory.CreateDirectoryInitializer(
			".opencode/command/spectr",
		),
		p.factory.CreateSlashCommandsInitializer(
			".opencode/command/spectr",
			extMD,
			FormatMarkdown,
			p.renderer,
		),
	}
}

// RegisterOpencodeProvider registers the OpenCode provider with the given
// registry.
func RegisterOpencodeProvider(
	reg *ProviderRegistry,
	renderer TemplateRenderer,
	factory InitializerFactory,
) error {
	return reg.Register(Registration{
		ID:       "opencode",
		Name:     "OpenCode",
		Priority: PriorityOpencode,
		Provider: NewOpencodeNewProvider(
			renderer,
			factory,
		),
	})
}

// =============================================================================
// Windsurf Provider
// =============================================================================

// WindsurfNewProvider implements the NewProvider interface for Windsurf.
// It uses .windsurf/commands/spectr/ for slash commands in Markdown
// format. It has no instruction file.
type WindsurfNewProvider struct {
	renderer TemplateRenderer
	factory  InitializerFactory
}

// NewWindsurfNewProvider creates a new Windsurf provider.
func NewWindsurfNewProvider(
	renderer TemplateRenderer,
	factory InitializerFactory,
) *WindsurfNewProvider {
	return &WindsurfNewProvider{
		renderer: renderer,
		factory:  factory,
	}
}

// Initializers returns the list of initializers needed to configure Windsurf.
func (p *WindsurfNewProvider) Initializers(
	_ context.Context,
) []Initializer {
	return []Initializer{
		p.factory.CreateDirectoryInitializer(
			".windsurf/commands/spectr",
		),
		p.factory.CreateSlashCommandsInitializer(
			".windsurf/commands/spectr",
			extMD,
			FormatMarkdown,
			p.renderer,
		),
	}
}

// RegisterWindsurfProvider registers the Windsurf provider with the given
// registry.
func RegisterWindsurfProvider(
	reg *ProviderRegistry,
	renderer TemplateRenderer,
	factory InitializerFactory,
) error {
	return reg.Register(Registration{
		ID:       "windsurf",
		Name:     "Windsurf",
		Priority: PriorityWindsurf,
		Provider: NewWindsurfNewProvider(
			renderer,
			factory,
		),
	})
}
