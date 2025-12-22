package providers

import (
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
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
func (p *QoderProvider) Initializers() []types.Initializer {
	proposalPath, applyPath := StandardCommandPaths(
		".qoder/commands",
		".md",
	)

	return []types.Initializer{
		initializers.NewConfigFileInitializer(
			"QODER.md",
		),
		initializers.NewSlashCommandsInitializer(
			"proposal",
			proposalPath,
			FrontmatterProposal,
		),
		initializers.NewSlashCommandsInitializer(
			"apply",
			applyPath,
			FrontmatterApply,
		),
	}
}
