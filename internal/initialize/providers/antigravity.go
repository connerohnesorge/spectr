package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/initialize/templates"
)

func init() {
	_ = RegisterProvider(Registration{
		ID:       "antigravity",
		Name:     "Antigravity",
		Priority: PriorityAntigravity,
		Provider: &AntigravityProvider{},
	})
}

// AntigravityProvider implements the Provider interface for Antigravity.
// Antigravity uses AGENTS.md and .agent/workflows/ for slash commands.
type AntigravityProvider struct{}

func (*AntigravityProvider) Initializers(
	_ context.Context,
) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(
			".agent/workflows",
		),
		NewConfigFileInitializer(
			"AGENTS.md",
			func(tm TemplateManager) any {
				return tm.Agents()
			},
		),
		NewSlashCommandsInitializer(
			".agent/workflows",
			".md",
			[]templates.SlashCommand{
				templates.SlashProposal,
				templates.SlashApply,
			},
		),
	}
}
