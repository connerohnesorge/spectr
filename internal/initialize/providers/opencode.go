package providers

import (
	ini "github.com/connerohnesorge/spectr/internal/initialize/providers/initializers" //nolint:revive
	"github.com/connerohnesorge/spectr/internal/initialize/types"
)

func init() {
	Register(Registration{
		ID:       "opencode",
		Name:     "OpenCode",
		Priority: 10,
		Provider: &OpencodeProvider{},
	})
}

// OpencodeProvider implements the new Provider interface for OpenCode.
type OpencodeProvider struct{}

// Initializers returns the initializers for OpenCode.
func (*OpencodeProvider) Initializers() []types.Initializer {
	proposalPath, applyPath := StandardCommandPaths(
		".opencode/command",
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
