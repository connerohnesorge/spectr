// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file implements the new Claude Code provider using the redesigned
// Provider interface with composable initializers.
package providers

import (
	"context"
)

// Compile-time interface satisfaction check.
var _ NewProvider = (*ClaudeNewProvider)(nil)

// ClaudeNewProvider implements the NewProvider interface for Claude Code.
// Claude Code uses CLAUDE.md for instructions and .claude/commands/spectr/
// for slash commands in Markdown format.
//
// This is the new implementation using the redesigned provider architecture.
// The old implementation is in claude.go.
type ClaudeNewProvider struct {
	// renderer is the template renderer for generating content.
	renderer TemplateRenderer

	// initializerFactory creates the initializers for this provider.
	// This factory pattern avoids import cycles between providers and
	// initializers packages.
	initializerFactory ClaudeInitializerFactory
}

// ClaudeInitializerFactory creates initializers for the Claude provider.
// This interface allows the initializers to be created outside of the
// providers package, avoiding import cycles.
type ClaudeInitializerFactory interface {
	// CreateDirectoryInitializer creates a DirectoryInitializer.
	CreateDirectoryInitializer(
		paths ...string,
	) Initializer

	// CreateConfigFileInitializer creates a ConfigFileInitializer.
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

// NewClaudeNewProvider creates a new Claude Code provider with the given
// renderer and factory. The renderer is used to render the instruction
// pointer template for the config file and slash command templates.
// The factory is used to create initializers without causing import cycles.
func NewClaudeNewProvider(
	renderer TemplateRenderer,
	factory ClaudeInitializerFactory,
) *ClaudeNewProvider {
	return &ClaudeNewProvider{
		renderer:           renderer,
		initializerFactory: factory,
	}
}

// Initializers returns the list of initializers needed to configure Claude
// Code. Claude Code requires:
//   - Directory: .claude/commands/spectr/
//   - Config file: CLAUDE.md with instruction pointer content
//   - Slash commands: proposal.md and apply.md in Markdown format
func (p *ClaudeNewProvider) Initializers(
	_ context.Context,
) []Initializer {
	// Render the instruction pointer template for the config file.
	// This template directs AI assistants to read the AGENTS.md file.
	templateCtx := DefaultTemplateContext()
	instructionContent, err := p.renderer.RenderInstructionPointer(
		templateCtx,
	)
	if err != nil {
		// If rendering fails, use an empty template.
		// The error will be surfaced when the initializer runs.
		instructionContent = ""
	}

	cmdDir := ".claude/commands/spectr"

	return []Initializer{
		// Create the slash commands directory
		p.initializerFactory.CreateDirectoryInitializer(
			cmdDir,
		),

		// Create/update the CLAUDE.md instruction file
		p.initializerFactory.CreateConfigFileInitializer(
			"CLAUDE.md",
			instructionContent,
		),

		// Create/update the slash command files
		p.initializerFactory.CreateSlashCommandsInitializer(
			cmdDir,
			".md",
			FormatMarkdown,
			p.renderer,
		),
	}
}

// RegisterClaudeProvider registers the Claude Code provider with the given
// registry. This function should be called at application startup to make
// Claude Code available.
//
// Example:
//
//	reg := providers.CreateRegistry()
//	tm, _ := initialize.NewTemplateManager()
//	factory := initializers.NewFactory()
//	providers.RegisterClaudeProvider(reg, tm, factory)
func RegisterClaudeProvider(
	reg *ProviderRegistry,
	renderer TemplateRenderer,
	factory ClaudeInitializerFactory,
) error {
	return reg.Register(Registration{
		ID:       "claude-code",
		Name:     "Claude Code",
		Priority: PriorityClaudeCode,
		Provider: NewClaudeNewProvider(
			renderer,
			factory,
		),
	})
}
