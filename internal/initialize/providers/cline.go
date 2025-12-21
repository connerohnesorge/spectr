package providers

import (
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
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
func (p *ClineProvider) Initializers() []types.Initializer {
	proposalPath, applyPath := StandardCommandPaths(
		".cline/commands",
		".md",
	)

	return []types.Initializer{
		initializers.NewSlashCommandsInitializer("proposal", proposalPath, FrontmatterProposal),
		initializers.NewSlashCommandsInitializer("apply", applyPath, FrontmatterApply),
	}
}