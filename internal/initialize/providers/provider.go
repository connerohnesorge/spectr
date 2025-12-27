// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// # Overview
//
// This package defines the Provider interface that all AI CLI tools
// (Claude Code, Gemini CLI, Cline, Cursor, etc.) must implement.
//
// Each provider returns a list of initializers that handle configuration setup.
package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/templates"
)

// Provider represents an AI CLI tool (Claude Code, Gemini, Cline, etc.).
// Each provider returns a list of initializers for its configuration.
type Provider interface {
	// Initializers returns the list of initializers for this provider.
	// Receives TemplateManager to allow passing TemplateRef directly to
	// initializers.
	Initializers(
		ctx context.Context,
		tm *templates.TemplateManager,
	) []Initializer
}
