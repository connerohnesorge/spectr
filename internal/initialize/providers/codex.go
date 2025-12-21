package providers

import (
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
	"github.com/connerohnesorge/spectr/internal/initialize/types"
)

func init() {
	Register(Registration{
		ID:       "codex",
		Name:     "Codex CLI",
		Priority: 9,
		Provider: &CodexProvider{},
	})
}

// CodexProvider implements the new Provider interface for Codex CLI.
type CodexProvider struct{}

// Initializers returns the initializers for Codex CLI.
func (p *CodexProvider) Initializers() []types.Initializer {
	// Codex uses global paths for commands
	proposalPath := "~/.codex/prompts/spectr-proposal.md"
	applyPath := "~/.codex/prompts/spectr-apply.md"

	return []types.Initializer{
		initializers.NewConfigFileInitializer("AGENTS.md"),
		initializers.NewSlashCommandsInitializer("proposal", proposalPath, FrontmatterProposal),
		initializers.NewSlashCommandsInitializer("apply", applyPath, FrontmatterApply),
	}
}