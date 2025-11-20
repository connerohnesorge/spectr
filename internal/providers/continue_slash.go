package providers

import (
	"github.com/conneroisu/spectr/internal/providerkit"
)

func init() {
	MustRegister(
		NewSlashMetadata(
			"continue",          // ID
			"Continue Commands", // Name
			[]string{
				".continue/commands/spectr-proposal.md",
				".continue/commands/spectr-apply.md",
				".continue/commands/spectr-archive.md",
			},
			PriorityContinueSlash, // Priority
		),
		func() providerkit.Provider {
			return providerkit.NewSlashCommandConfigurator(
				providerkit.SlashCommandConfig{
					ToolID:   "continue",
					ToolName: "Continue Commands",
					Frontmatter: map[string]string{
						"proposal": `---
name: spectr-proposal
description: Scaffold a new Spectr change and validate strictly.
---`,
						"apply": `---
name: spectr-apply
description: Implement an approved Spectr change and keep tasks in sync.
---`,
						"archive": `---
name: spectr-archive
description: Archive a deployed Spectr change and update specs.
---`,
					},
					FilePaths: map[string]string{
						"proposal": ".continue/commands/spectr-proposal.md",
						"apply":    ".continue/commands/spectr-apply.md",
						"archive":  ".continue/commands/spectr-archive.md",
					},
				})
		},
	)
}
