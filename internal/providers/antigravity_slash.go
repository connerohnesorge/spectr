package providers

import (
	"github.com/conneroisu/spectr/internal/providerkit"
)

func init() {
	MustRegister(
		NewSlashMetadata(
			"antigravity-slash",     // ID
			"Antigravity Workflows", // Name
			[]string{
				".agent/workflows/spectr-proposal.md",
				".agent/workflows/spectr-apply.md",
				".agent/workflows/spectr-archive.md",
			},
			PriorityAntigravitySlash, // Priority
		),
		func() providerkit.Provider {
			return providerkit.NewSlashCommandConfigurator(
				providerkit.SlashCommandConfig{
					ToolID:   "antigravity",
					ToolName: "Antigravity Workflows",
					// No frontmatter for Antigravity
					Frontmatter: make(map[string]string),
					FilePaths: map[string]string{
						"proposal": ".agent/workflows/spectr-proposal.md",
						"apply":    ".agent/workflows/spectr-apply.md",
						"archive":  ".agent/workflows/spectr-archive.md",
					},
				})
		},
	)
}
