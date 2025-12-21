package providers

import (
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
	"github.com/connerohnesorge/spectr/internal/initialize/types"
)

func init() {
	Register(Registration{
		ID:       "crush",
		Name:     "Crush",
		Priority: 16,
		Provider: &CrushProvider{},
	})
}

// CrushProvider implements the new Provider interface for Crush.
type CrushProvider struct{}

// Initializers returns the initializers for Crush.
func (p *CrushProvider) Initializers() []types.Initializer {
	proposalPath, applyPath := StandardCommandPaths(
		".crush/commands",
		".md",
	)

	return []types.Initializer{
		initializers.NewConfigFileInitializer("CRUSH.md"),
		initializers.NewSlashCommandsInitializer("proposal", proposalPath, FrontmatterProposal),
		initializers.NewSlashCommandsInitializer("apply", applyPath, FrontmatterApply),
	}
}