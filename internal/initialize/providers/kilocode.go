package providers

import (
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
	"github.com/connerohnesorge/spectr/internal/initialize/types"
)

func init() {
	Register(Registration{
		ID:       "kilocode",
		Name:     "Kilocode",
		Priority: 14,
		Provider: &KilocodeProvider{},
	})
}

// KilocodeProvider implements the new Provider interface for Kilocode.
type KilocodeProvider struct{}

// Initializers returns the initializers for Kilocode.
func (p *KilocodeProvider) Initializers() []types.Initializer {
	proposalPath, applyPath := StandardCommandPaths(
		".kilocode/commands",
		".md",
	)

	return []types.Initializer{
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
