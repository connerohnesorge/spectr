package providers

import (
	"github.com/conneroisu/spectr/internal/providerkit"
)

func init() {
	MustRegister(
		NewSlashMetadata(
			"aider",          // ID
			"Aider Commands", // Name
			[]string{
				".aider/commands/spectr-proposal.md",
				".aider/commands/spectr-apply.md",
				".aider/commands/spectr-archive.md",
			},
			PriorityAiderSlash, // Priority
		),
		func() providerkit.Provider {
			return providerkit.NewSlashCommandConfigurator(
				providerkit.SlashCommandConfig{
					ToolID:   "aider",
					ToolName: "Aider Commands",
					// No frontmatter for Aider
					Frontmatter: make(map[string]string),
					FilePaths: map[string]string{
						"proposal": ".aider/commands/spectr-proposal.md",
						"apply":    ".aider/commands/spectr-apply.md",
						"archive":  ".aider/commands/spectr-archive.md",
					},
				})
		},
	)
}
