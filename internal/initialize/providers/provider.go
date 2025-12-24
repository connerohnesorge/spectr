// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// # Overview
//
// This package defines the Provider interface that all AI CLI tools
// (Claude Code, Gemini CLI, Cline, Cursor, etc.) must implement.
//
// Each provider returns a list of initializers that handle configuration
// (directories, config files, slash commands).
//
// # Adding a New Provider
//
// To add a new AI CLI provider:
//
// 1. Create a new file (e.g., providers/mytool.go)
// 2. Implement the Provider interface with Initializers() method
// 3. Register in RegisterAllProviders() in registry.go
//
// Example:
//
//	package providers
//
//	import (
//	    "context"
//	    "github.com/connerohnesorge/spectr/internal/domain"
//	    "github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
//	)
//
//	type MyToolProvider struct{}
//
//	func (p *MyToolProvider) Initializers(_ context.Context) []Initializer {
//	    return []Initializer{
//	        initializers.NewDirectoryInitializer(".mytool/commands/spectr"),
//	        initializers.NewConfigFileInitializer("MYTOOL.md", func(tm *initialize.TemplateManager) domain.TemplateRef {
//	            return tm.InstructionPointer()
//	        }),
//	        initializers.NewSlashCommandsInitializer(".mytool/commands/spectr", ".md", []domain.SlashCommand{
//	            domain.SlashProposal,
//	            domain.SlashApply,
//	        }),
//	    }
//	}
//
//nolint:revive // File length acceptable for provider interface definition
package providers

import (
	"context"
	"text/template"
)

// TemplateProvider provides access to templates for rendering.
// This interface avoids import cycles by not depending on the initialize package.
type TemplateProvider interface {
	// GetTemplates returns the underlying template set for rendering.
	GetTemplates() *template.Template
}

// Provider is the minimal provider interface for the redesigned architecture.
// Providers return a list of initializers; all metadata (ID, Name, Priority) is
// stored in the Registration struct, not in the provider itself.
type Provider interface {
	// Initializers returns the list of initialization steps for this provider.
	// Each initializer handles one aspect of configuration (directories, files, commands).
	Initializers(_ context.Context) []Initializer
}
