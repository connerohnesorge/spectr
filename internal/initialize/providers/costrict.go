package providers

import (
	ini "github.com/connerohnesorge/spectr/internal/initialize/providers/initializers" //nolint:revive
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
func (*CostrictProvider) Initializers() []types.Initializer {
	proposalPath, applyPath := StandardCommandPaths(
		".costrict/commands",
		".md",
	)

	return []types.Initializer{
		ini.NewConfigFileInitializer(
			"COSTRICT.md",
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
