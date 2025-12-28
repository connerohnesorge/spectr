// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// # Overview
//
// This package defines the Provider interface that all AI CLI tools
// (Claude Code, Gemini CLI, Cline, Cursor, etc.) must implement.
//
// Each provider returns a list of Initializers that handle the creation
// of instruction files (e.g., CLAUDE.md) and slash commands
// (e.g., .claude/commands/).
//
// # Adding a New Provider
//
// To add a new AI CLI provider:
//  1. Create a new file (e.g., providers/mytool.go)
//  2. Implement the Provider interface with an Initializers() method
//  3. Add the provider to RegisterAllProviders() in registry.go
//
// Example:
//
//	package providers
//
//	import "context"
//
//	type MyToolProvider struct{}
//
//	func (p *MyToolProvider) Initializers(ctx context.Context, tm any) []domain.Initializer {
//		return []domain.Initializer{
//			NewDirectoryInitializer(".mytool/commands/spectr"),
//			NewConfigFileInitializer("MYTOOL.md", tm.InstructionPointer()),
//			NewSlashCommandsInitializer(".mytool/commands/spectr", ...),
//		}
//	}
//
// The Initializer types handle all common logic for creating directories,
// config files, and slash commands.
package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/domain"
)

// Initializer is an alias for domain.Initializer for convenience.
// This allows provider files to use Initializer instead of domain.Initializer.
type Initializer = domain.Initializer

// Provider represents an AI CLI tool (Claude Code, Gemini, Cline, etc.).
// Each provider returns a list of Initializers that handle the creation
// of instruction files and slash commands.
type Provider interface {
	// Initializers returns the list of initializers for this provider.
	// Receives TemplateManager to allow passing TemplateRef directly to initializers.
	Initializers(ctx context.Context, tm any) []Initializer

	// IsConfigured checks if this provider is already configured in the project.
	IsConfigured(projectDir string) bool
}
