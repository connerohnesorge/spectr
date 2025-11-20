package providers

import (
	"github.com/conneroisu/spectr/internal/providerkit"
)

func init() {
	MustRegister(
		NewSlashMetadata(
			"qoder-slash",    // ID
			"Qoder Commands", // Name
			[]string{
				".qoder/commands/spectr/proposal.md",
				".qoder/commands/spectr/apply.md",
				".qoder/commands/spectr/archive.md",
			},
			PriorityQoderSlash, // Priority
		),
		func() providerkit.Provider {
			return providerkit.NewSlashCommandConfigurator(
				providerkit.SlashCommandConfig{
					ToolID:   "qoder",
					ToolName: "Qoder Commands",
					Frontmatter: map[string]string{
						"proposal": `---
name: Spectr: Proposal
description: Scaffold a new Spectr change and validate strictly.
category: Spectr
tags: [spectr, change]
---`,
						"apply": `---
name: Spectr: Apply
description: Implement an approved Spectr change and keep tasks in sync.
category: Spectr
tags: [spectr, apply]
---`,
						"archive": `---
name: Spectr: Archive
description: Archive a deployed Spectr change and update specs.
category: Spectr
tags: [spectr, archive]
---`,
					},
					FilePaths: map[string]string{
						"proposal": ".qoder/commands/spectr/proposal.md",
						"apply":    ".qoder/commands/spectr/apply.md",
						"archive":  ".qoder/commands/spectr/archive.md",
					},
				})
		},
	)
}
