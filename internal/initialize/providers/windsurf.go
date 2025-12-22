package providers

import (
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
	"github.com/connerohnesorge/spectr/internal/initialize/types"
)

func init() {
	Register(Registration{
		ID:       "windsurf",
		Name:     "Windsurf",
		Priority: 13,
		Provider: &WindsurfProvider{},
	})
}

// WindsurfProvider implements the new Provider interface for Windsurf.
type WindsurfProvider struct{}

// Initializers returns the initializers for Windsurf.
func (p *WindsurfProvider) Initializers() []types.Initializer {
	proposalPath, applyPath := StandardCommandPaths(
		".windsurf/commands",
		".md",
	)

	return []types.Initializer{
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
