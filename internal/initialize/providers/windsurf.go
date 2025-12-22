package providers

import (
	inits "github.com/connerohnesorge/spectr/internal/initialize/providers/initializers" //nolint:revive
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
func (*WindsurfProvider) Initializers() []types.Initializer {
	proposalPath, applyPath := StandardCommandPaths(
		".windsurf/commands",
		".md",
	)

	return []types.Initializer{
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