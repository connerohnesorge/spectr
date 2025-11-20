package providers

import (
	"github.com/conneroisu/spectr/internal/providerkit"
)

func init() {
	MustRegister(
		NewSlashMetadata(
			"kilocode",           // ID
			"Kilocode Workflows", // Name
			[]string{
				".kilocode/workflows/spectr-proposal.md",
				".kilocode/workflows/spectr-apply.md",
				".kilocode/workflows/spectr-archive.md",
			},
			PriorityKilocodeSlash, // Priority
		),
		func() providerkit.Provider {
			return providerkit.NewSlashCommandConfigurator(
				providerkit.SlashCommandConfig{
					ToolID:   "kilocode",
					ToolName: "Kilocode Workflows",
					// No frontmatter for Kilocode
					Frontmatter: make(map[string]string),
					FilePaths: map[string]string{
						"proposal": ".kilocode/workflows/spectr-proposal.md",
						"apply":    ".kilocode/workflows/spectr-apply.md",
						"archive":  ".kilocode/workflows/spectr-archive.md",
					},
				})
		},
	)
}
