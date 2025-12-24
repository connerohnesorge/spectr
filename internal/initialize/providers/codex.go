package providers

import (
	"context"

	//nolint:revive // long import path is unavoidable
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
)

// CodexProvider configures Codex CLI.
// Codex uses AGENTS.md and global ~/.codex/prompts/ for commands.
type CodexProvider struct{}

func (*CodexProvider) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		initializers.NewConfigFileInitializer("AGENTS.md", "AGENTS.md.tmpl"),
		// TODO: Add global slash commands for ~/.codex/prompts/
		// Requires a custom global filesystem initializer.
	}
}
