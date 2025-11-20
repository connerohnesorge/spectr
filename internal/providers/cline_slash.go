package providers

import (
	"github.com/conneroisu/spectr/internal/providerkit"
)

func init() {
	MustRegister(
		NewSlashMetadata(
			"cline-slash", // ID
			"Cline Rules", // Name
			[]string{
				".clinerules/spectr-proposal.md",
				".clinerules/spectr-apply.md",
				".clinerules/spectr-archive.md",
			},
			PriorityClineSlash, // Priority
		),
		func() providerkit.Provider {
			return providerkit.NewSlashCommandConfigurator(
				providerkit.SlashCommandConfig{
					ToolID:   "cline",
					ToolName: "Cline Rules",
					Frontmatter: map[string]string{
						"proposal": "# Spectr: Proposal\n\n" +
							"Scaffold a new Spectr change and " +
							"validate strictly.",
						"apply": "# Spectr: Apply\n\n" +
							"Implement an approved Spectr change " +
							"and keep tasks in sync.",
						"archive": "# Spectr: Archive\n\n" +
							"Archive a deployed Spectr change " +
							"and update specs.",
					},
					FilePaths: map[string]string{
						"proposal": ".clinerules/spectr-proposal.md",
						"apply":    ".clinerules/spectr-apply.md",
						"archive":  ".clinerules/spectr-archive.md",
					},
				})
		},
	)
}
