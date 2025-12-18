package providers

import "context"

func init() {
	err := RegisterV2(Registration{
		ID:       "antigravity",
		Name:     "Antigravity",
		Priority: PriorityAntigravity,
		Provider: &AntigravityProvider{},
	})
	if err != nil {
		panic(err)
	}
}

// AntigravityProvider implements the ProviderV2 interface for Antigravity.
// Antigravity uses AGENTS.md and .agent/workflows/ for slash commands.
// Note: Uses prefixed command paths (spectr-proposal.md, spectr-apply.md)
// instead of subdirectory structure.
type AntigravityProvider struct{}

// Initializers returns the initializers for Antigravity.
func (p *AntigravityProvider) Initializers(ctx context.Context) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".agent/workflows"),
		NewConfigFileInitializer("AGENTS.md"),
		NewSlashCommandsInitializerWithFrontmatter(
			".agent/workflows",
			".md",
			FormatMarkdown,
			map[string]string{
				"proposal": FrontmatterProposal,
				"apply":    FrontmatterApply,
			},
		),
	}
}
