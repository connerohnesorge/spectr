package providers

import (
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
	"github.com/connerohnesorge/spectr/internal/initialize/types"
)

func init() {
	Register(Registration{
		ID:       "gemini",
		Name:     "Gemini CLI",
		Priority: 2,
		Provider: &GeminiProvider{},
	})
}

// GeminiProvider implements the new Provider interface for Gemini CLI.
type GeminiProvider struct{}

// Initializers returns the initializers for Gemini CLI.
func (p *GeminiProvider) Initializers() []types.Initializer {
	proposalPath, applyPath := StandardCommandPaths(
		".gemini/commands",
		".toml",
	)

	return []types.Initializer{
		initializers.NewTOMLCommandInitializer("proposal", proposalPath, "Scaffold a new Spectr change and validate strictly."),
		initializers.NewTOMLCommandInitializer("apply", applyPath, "Implement an approved Spectr change and keep tasks in sync."),
	}
}