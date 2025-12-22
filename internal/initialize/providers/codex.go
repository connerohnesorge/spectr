package providers

import (
	ini "github.com/connerohnesorge/spectr/internal/initialize/providers/initializers" //nolint:revive
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
func (*CodexProvider) Initializers() []types.Initializer {
	// Codex uses global paths for commands
	proposalPath := "~/.codex/prompts/spectr-proposal.md"
	applyPath := "~/.codex/prompts/spectr-apply.md"

	return []types.Initializer{
		ini.NewConfigFileInitializer(
			"AGENTS.md",
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