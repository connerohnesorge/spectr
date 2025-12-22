package providers

import (
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
	"github.com/connerohnesorge/spectr/internal/initialize/types"
)

func init() {
	Register(Registration{
		ID:       "antigravity",
		Name:     "Antigravity",
		Priority: 6,
		Provider: &AntigravityProvider{},
	})
}

// AntigravityProvider implements the new Provider interface for Antigravity.
type AntigravityProvider struct{}

// Initializers returns the initializers for Antigravity.
func (p *AntigravityProvider) Initializers() []types.Initializer {
	proposalPath, applyPath := PrefixedCommandPaths(
		".agent/workflows",
		".md",
	)

	return []types.Initializer{
		initializers.NewConfigFileInitializer(
			"AGENTS.md",
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
