package providers

import (
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
	"github.com/connerohnesorge/spectr/internal/initialize/types"
)

func init() {
	Register(Registration{
		ID:       "costrict",
		Name:     "CoStrict",
		Priority: 3,
		Provider: &CostrictProvider{},
	})
}

// CostrictProvider implements the new Provider interface for CoStrict.
type CostrictProvider struct{}

// Initializers returns the initializers for CoStrict.
func (p *CostrictProvider) Initializers() []types.Initializer {
	proposalPath, applyPath := StandardCommandPaths(
		".costrict/commands",
		".md",
	)

	return []types.Initializer{
		initializers.NewConfigFileInitializer("COSTRICT.md"),
		initializers.NewSlashCommandsInitializer("proposal", proposalPath, FrontmatterProposal),
		initializers.NewSlashCommandsInitializer("apply", applyPath, FrontmatterApply),
	}
}