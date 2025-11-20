package providers

import (
	"github.com/conneroisu/spectr/internal/providerkit"
)

func init() {
	MustRegister(
		NewSlashMetadata(
			"windsurf",           // ID
			"Windsurf Workflows", // Name
			[]string{
				".windsurf/workflows/spectr-proposal.md",
				".windsurf/workflows/spectr-apply.md",
				".windsurf/workflows/spectr-archive.md",
			},
			PriorityWindsurfSlash, // Priority
		),
		func() providerkit.Provider {
			return providerkit.NewSlashCommandConfigurator(
				providerkit.SlashCommandConfig{
					ToolID:   "windsurf",
					ToolName: "Windsurf Workflows",
					Frontmatter: map[string]string{
						"proposal": "---\n" +
							"description: Scaffold a new Spectr " +
							"change and validate strictly.\n" +
							"auto_execution_mode: 3\n---",
						"apply": "---\n" +
							"description: Implement an approved " +
							"Spectr change and keep tasks in sync.\n" +
							"auto_execution_mode: 3\n---",
						"archive": "---\n" +
							"description: Archive a deployed Spectr " +
							"change and update specs.\n" +
							"auto_execution_mode: 3\n---",
					},
					FilePaths: map[string]string{
						"proposal": ".windsurf/workflows/spectr-proposal.md",
						"apply":    ".windsurf/workflows/spectr-apply.md",
						"archive":  ".windsurf/workflows/spectr-archive.md",
					},
				})
		},
	)
}
