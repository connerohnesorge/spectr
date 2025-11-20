package providers

import (
	"github.com/conneroisu/spectr/internal/providerkit"
)

func init() {
	MustRegister(
		NewSlashMetadata(
			"cursor",          // ID
			"Cursor Commands", // Name
			[]string{
				".cursor/commands/spectr-proposal.md",
				".cursor/commands/spectr-apply.md",
				".cursor/commands/spectr-archive.md",
			},
			PriorityCursorSlash, // Priority
		),
		func() providerkit.Provider {
			return providerkit.NewSlashCommandConfigurator(
				providerkit.SlashCommandConfig{
					ToolID:   "cursor",
					ToolName: "Cursor Commands",
					Frontmatter: map[string]string{
						"proposal": `---
name: /spectr-proposal
id: spectr-proposal
category: Spectr
description: Scaffold a new Spectr change and validate strictly.
---`,
						"apply": `---
name: /spectr-apply
id: spectr-apply
category: Spectr
description: Implement an approved Spectr change and keep tasks in sync.
---`,
						"archive": `---
name: /spectr-archive
id: spectr-archive
category: Spectr
description: Archive a deployed Spectr change and update specs.
---`,
					},
					FilePaths: map[string]string{
						"proposal": ".cursor/commands/spectr-proposal.md",
						"apply":    ".cursor/commands/spectr-apply.md",
						"archive":  ".cursor/commands/spectr-archive.md",
					},
				})
		},
	)
}
