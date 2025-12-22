package providers

import (
	ini "github.com/connerohnesorge/spectr/internal/initialize/providers/initializers" //nolint:revive
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
func (*QwenProvider) Initializers() []types.Initializer {
	proposalPath, applyPath := StandardCommandPaths(
		".qwen/commands",
		".md",
	)

	return []types.Initializer{
		ini.NewConfigFileInitializer(
			"QWEN.md",
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