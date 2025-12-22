package providers

import (
	inits "github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
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