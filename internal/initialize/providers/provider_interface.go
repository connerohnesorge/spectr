package providers

import (
	"context"
)

// Provider represents an AI CLI/IDE tool
// (Claude Code, Gemini, Cline, etc.).
// This is the minimal provider interface that returns
// composable initializers.
//
// Example implementation:
//
//	type ClaudeProvider struct{}
//
//	func (p *ClaudeProvider) Initializers(
//	    ctx context.Context,
//	) []Initializer {
//	    return []Initializer{
//	        NewDirectoryInitializer(".claude/commands/spectr"),
//	        NewConfigFileInitializer(
//	            "CLAUDE.md",
//	            InstructionTemplate,
//	        ),
//	        NewSlashCommandsInitializer(
//	            ".claude/commands/spectr",
//	            ".md",
//	            FormatMarkdown,
//	        ),
//	    }
//	}
//
// Provider metadata (ID, name, priority) is provided at registration time
// via the Registration struct, not by the Provider interface.
type Provider interface {
	// Initializers returns the list of initializers for this provider.
	// Each initializer represents a single initialization step
	// (directory, file, etc.).
	// Initializers may be shared across providers and will be
	// deduplicated.
	Initializers(ctx context.Context) []Initializer
}
