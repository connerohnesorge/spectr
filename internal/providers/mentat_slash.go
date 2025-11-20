package providers

import (
	"github.com/conneroisu/spectr/internal/providerkit"
)

func init() {
	MustRegister(
		NewSlashMetadata(
			"mentat",          // ID
			"Mentat Commands", // Name
			[]string{
				".mentat/commands/spectr-proposal.md",
				".mentat/commands/spectr-apply.md",
				".mentat/commands/spectr-archive.md",
			},
			PriorityMentatSlash, // Priority
		),
		func() providerkit.Provider {
			return providerkit.NewSlashCommandConfigurator(
				providerkit.SlashCommandConfig{
					ToolID:   "mentat",
					ToolName: "Mentat Commands",
					// No frontmatter for Mentat
					Frontmatter: make(map[string]string),
					FilePaths: map[string]string{
						"proposal": ".mentat/commands/spectr-proposal.md",
						"apply":    ".mentat/commands/spectr-apply.md",
						"archive":  ".mentat/commands/spectr-archive.md",
					},
				})
		},
	)
}
