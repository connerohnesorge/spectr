package providers

import (
	ini "github.com/connerohnesorge/spectr/internal/initialize/providers/initializers" //nolint:revive
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
func (*KilocodeProvider) Initializers() []types.Initializer {
	proposalPath, applyPath := StandardCommandPaths(
		".kilocode/commands",
		".md",
	)

	return []types.Initializer{
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
