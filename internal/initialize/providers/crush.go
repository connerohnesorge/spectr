package providers

import (
	ini "github.com/connerohnesorge/spectr/internal/initialize/providers/initializers" //nolint:revive
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
func (*CrushProvider) Initializers() []types.Initializer {
	proposalPath, applyPath := StandardCommandPaths(
		".crush/commands",
		".md",
	)

	return []types.Initializer{
		ini.NewConfigFileInitializer(
			"CRUSH.md",
		),
		ini.NewSlashCommandsInitializer(
			"proposal",
			proposalPath,
			FrontmatterProposal,
		),
		ini.NewSlashCommandsInitializer(
			"apply",
			applyPath,
			FrontmatterApply,
		),
	}
}