package providers

import (
	inits "github.com/connerohnesorge/spectr/internal/initialize/providers/initializers" //nolint:revive
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
func (*AntigravityProvider) Initializers() []types.Initializer {
	proposalPath, applyPath := PrefixedCommandPaths(
		".agent/workflows",
		".md",
	)

	return []types.Initializer{
		inits.NewConfigFileInitializer(
			"AGENTS.md",
		),
		inits.NewSlashCommandsInitializer(
			"proposal",
			proposalPath,
			FrontmatterProposal,
		),
		inits.NewSlashCommandsInitializer(
			"apply",
			applyPath,
			FrontmatterApply,
		),
	}
}
