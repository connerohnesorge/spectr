package providers

import (
	"github.com/conneroisu/spectr/internal/providerkit"
)

func init() {
	MustRegister(
		NewSlashMetadata(
			"copilot",                     // ID
			"GitHub Copilot Instructions", // Name
			[]string{
				".github/copilot/spectr-proposal.md",
				".github/copilot/spectr-apply.md",
				".github/copilot/spectr-archive.md",
			},
			PriorityCopilotSlash, // Priority
		),
		func() providerkit.Provider {
			return providerkit.NewSlashCommandConfigurator(
				providerkit.SlashCommandConfig{
					ToolID:   "copilot",
					ToolName: "GitHub Copilot Instructions",
					// No frontmatter for Copilot
					Frontmatter: make(map[string]string),
					FilePaths: map[string]string{
						"proposal": ".github/copilot/spectr-proposal.md",
						"apply":    ".github/copilot/spectr-apply.md",
						"archive":  ".github/copilot/spectr-archive.md",
					},
				})
		},
	)
}
