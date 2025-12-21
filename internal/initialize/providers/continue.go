package providers

import (
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
	"github.com/connerohnesorge/spectr/internal/initialize/types"
)

func init() {
	Register(Registration{
		ID:       "continue",
		Name:     "Continue",
		Priority: 15,
		Provider: &ContinueProvider{},
	})
}

// ContinueProvider implements the new Provider interface for Continue.
type ContinueProvider struct{}

// Initializers returns the initializers for Continue.
func (p *ContinueProvider) Initializers() []types.Initializer {
	proposalPath, applyPath := StandardCommandPaths(
		".continue/commands",
		".md",
	)

	return []types.Initializer{
		initializers.NewSlashCommandsInitializer("proposal", proposalPath, FrontmatterProposal),
		initializers.NewSlashCommandsInitializer("apply", applyPath, FrontmatterApply),
	}
}