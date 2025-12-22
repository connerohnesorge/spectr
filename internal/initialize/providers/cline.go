package providers

import (
	inits "github.com/connerohnesorge/spectr/internal/initialize/providers/initializers" //nolint:revive
	"github.com/connerohnesorge/spectr/internal/initialize/types"
)

func init() {
	Register(Registration{
		ID:       "cline",
		Name:     "Cline",
		Priority: 7,
		Provider: &ClineProvider{},
	})
}

// ClineProvider implements the new Provider interface for Cline.
type ClineProvider struct{}

// Initializers returns the initializers for Cline.
func (*ClineProvider) Initializers() []types.Initializer {
	proposalPath, applyPath := StandardCommandPaths(
		".cline/commands",
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
