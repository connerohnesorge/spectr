package providers

import (
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
	"github.com/connerohnesorge/spectr/internal/initialize/types"
)

func init() {
	Register(Registration{
		ID:       "claude-code",
		Name:     "Claude Code",
		Priority: 1,
		Provider: &ClaudeProvider{},
	})
}

// ClaudeProvider implements the new Provider interface for Claude Code.
type ClaudeProvider struct{}

// Initializers returns the initializers for Claude Code.
func (*ClaudeProvider) Initializers() []types.Initializer {
	proposalPath, applyPath := StandardCommandPaths(
		".claude/commands",
		".md",
	)

	return []types.Initializer{
		initializers.NewConfigFileInitializer(
			"CLAUDE.md",
		),
		initializers.NewSlashCommandsInitializer(
			"proposal",
			proposalPath,
			FrontmatterProposal,
		),
		initializers.NewSlashCommandsInitializer(
			"apply",
			applyPath,
			FrontmatterApply,
		),
	}
}
