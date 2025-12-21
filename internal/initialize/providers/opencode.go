package providers

import (
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
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
func (p *OpencodeProvider) Initializers() []types.Initializer {
	proposalPath, applyPath := StandardCommandPaths(
		".opencode/command",
		".md",
	)

	return []types.Initializer{
		initializers.NewSlashCommandsInitializer("proposal", proposalPath, FrontmatterProposal),
		initializers.NewSlashCommandsInitializer("apply", applyPath, FrontmatterApply),
	}
}