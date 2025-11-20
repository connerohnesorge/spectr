package providers

import (
	"github.com/conneroisu/spectr/internal/providerkit"
)

func init() {
	MustRegister(
		NewSlashMetadata(
			"codebuddy-slash",    // ID
			"CodeBuddy Commands", // Name
			[]string{
				".codebuddy/commands/spectr/proposal.md",
				".codebuddy/commands/spectr/apply.md",
				".codebuddy/commands/spectr/archive.md",
			},
			PriorityCodeBuddySlash, // Priority
		),
		func() providerkit.Provider {
			return providerkit.NewSlashCommandConfigurator(
				providerkit.SlashCommandConfig{
					ToolID:   "codebuddy",
					ToolName: "CodeBuddy Commands",
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
						"proposal": ".codebuddy/commands/spectr/proposal.md",
						"apply":    ".codebuddy/commands/spectr/apply.md",
						"archive":  ".codebuddy/commands/spectr/archive.md",
					},
				})
		},
	)
}
