package providers

import (
	"github.com/conneroisu/spectr/internal/providerkit"
)

func init() {
	MustRegister(
		NewSlashMetadata(
			"smol",          // ID
			"Smol Commands", // Name
			[]string{
				".smol/commands/spectr-proposal.md",
				".smol/commands/spectr-apply.md",
				".smol/commands/spectr-archive.md",
			},
			PrioritySmolSlash, // Priority
		),
		func() providerkit.Provider {
			return providerkit.NewSlashCommandConfigurator(
				providerkit.SlashCommandConfig{
					ToolID:   "smol",
					ToolName: "Smol Commands",
					// No frontmatter for Smol
					Frontmatter: make(map[string]string),
					FilePaths: map[string]string{
						"proposal": ".smol/commands/spectr-proposal.md",
						"apply":    ".smol/commands/spectr-apply.md",
						"archive":  ".smol/commands/spectr-archive.md",
					},
				})
		},
	)
}
