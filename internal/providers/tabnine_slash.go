package providers

import (
	"github.com/conneroisu/spectr/internal/providerkit"
)

func init() {
	MustRegister(
		NewSlashMetadata(
			"tabnine",          // ID
			"Tabnine Commands", // Name
			[]string{
				".tabnine/commands/spectr-proposal.md",
				".tabnine/commands/spectr-apply.md",
				".tabnine/commands/spectr-archive.md",
			},
			PriorityTabnineSlash, // Priority
		),
		func() providerkit.Provider {
			return providerkit.NewSlashCommandConfigurator(
				providerkit.SlashCommandConfig{
					ToolID:   "tabnine",
					ToolName: "Tabnine Commands",
					// No frontmatter for Tabnine
					Frontmatter: make(map[string]string),
					FilePaths: map[string]string{
						"proposal": ".tabnine/commands/spectr-proposal.md",
						"apply":    ".tabnine/commands/spectr-apply.md",
						"archive":  ".tabnine/commands/spectr-archive.md",
					},
				})
		},
	)
}
