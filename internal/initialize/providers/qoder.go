package providers

import (
	inits "github.com/connerohnesorge/spectr/internal/initialize/providers/initializers" //nolint:revive
	"github.com/connerohnesorge/spectr/internal/initialize/types"
)

func init() {
	Register(Registration{
		ID:       "qoder",
		Name:     "Qoder",
		Priority: 4,
		Provider: &QoderProvider{},
	})
}

// QoderProvider implements the new Provider interface for Qoder.
type QoderProvider struct{}

// Initializers returns the initializers for Qoder.
func (*QoderProvider) Initializers() []types.Initializer {
	proposalPath, applyPath := StandardCommandPaths(
		".qoder/commands",
		".md",
	)

	return []types.Initializer{
		inits.NewConfigFileInitializer(
			"QODER.md",
		),
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
