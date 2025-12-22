package providers

import (
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
	"github.com/connerohnesorge/spectr/internal/initialize/types"
)

func init() {
	Register(Registration{
		ID:       "qwen",
		Name:     "Qwen Code",
		Priority: 5,
		Provider: &QwenProvider{},
	})
}

// QwenProvider implements the new Provider interface for Qwen Code.
type QwenProvider struct{}

// Initializers returns the initializers for Qwen Code.
func (p *QwenProvider) Initializers() []types.Initializer {
	proposalPath, applyPath := StandardCommandPaths(
		".qwen/commands",
		".md",
	)

	return []types.Initializer{
		initializers.NewConfigFileInitializer(
			"QWEN.md",
		),
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
