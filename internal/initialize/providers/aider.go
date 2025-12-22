// Package providers implements the different AI provider configurations.
package providers

import (
	inits "github.com/connerohnesorge/spectr/internal/initialize/providers/initializers" //nolint:revive
	"github.com/connerohnesorge/spectr/internal/initialize/types"
)

func init() {
	Register(Registration{
		ID:       "aider",
		Name:     "Aider",
		Priority: 11,
		Provider: &AiderProvider{},
	})
}

// AiderProvider implements the new Provider interface for Aider.
type AiderProvider struct{}

// Initializers returns the initializers for Aider.
func (*AiderProvider) Initializers() []types.Initializer {
	proposalPath, applyPath := StandardCommandPaths(
		".aider/commands",
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