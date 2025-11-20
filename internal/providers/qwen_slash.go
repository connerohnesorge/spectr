package providers

import (
	"github.com/conneroisu/spectr/internal/providerkit"
)

func init() {
	MustRegister(
		NewSlashMetadata(
			"qwen-slash",    // ID
			"Qwen Commands", // Name
			[]string{
				".qwen/commands/spectr-proposal.md",
				".qwen/commands/spectr-apply.md",
				".qwen/commands/spectr-archive.md",
			},
			PriorityQwenSlash, // Priority
		),
		func() providerkit.Provider {
			return providerkit.NewSlashCommandConfigurator(
				providerkit.SlashCommandConfig{
					ToolID:   "qwen",
					ToolName: "Qwen Commands",
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
						"proposal": ".qwen/commands/spectr-proposal.md",
						"apply":    ".qwen/commands/spectr-apply.md",
						"archive":  ".qwen/commands/spectr-archive.md",
					},
				})
		},
	)
}
