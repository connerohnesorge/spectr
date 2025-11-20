package providers

import (
	"github.com/conneroisu/spectr/internal/providerkit"
)

func init() {
	MustRegister(
		NewSlashMetadata(
			"claude",                // ID
			"Claude Slash Commands", // Name
			[]string{ // File paths
				".claude/commands/spectr/proposal.md",
				".claude/commands/spectr/apply.md",
				".claude/commands/spectr/archive.md",
			},
			PriorityClaudeSlash, // Priority
		),
		func() providerkit.Provider {
			return providerkit.NewSlashCommandConfigurator(
				providerkit.SlashCommandConfig{
					ToolID:   "claude",
					ToolName: "Claude Slash Commands",
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
						"proposal": ".claude/commands/spectr/proposal.md",
						"apply":    ".claude/commands/spectr/apply.md",
						"archive":  ".claude/commands/spectr/archive.md",
					},
				})
		},
	)
}
