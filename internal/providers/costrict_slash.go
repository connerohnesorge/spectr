package providers

import (
	"github.com/conneroisu/spectr/internal/providerkit"
)

func init() {
	MustRegister(
		NewSlashMetadata(
			"costrict-slash",    // ID
			"CoStrict Commands", // Name
			[]string{
				".cospec/spectr/commands/spectr-proposal.md",
				".cospec/spectr/commands/spectr-apply.md",
				".cospec/spectr/commands/spectr-archive.md",
			},
			PriorityCostrictSlash, // Priority
		),
		func() providerkit.Provider {
			return providerkit.NewSlashCommandConfigurator(
				providerkit.SlashCommandConfig{
					ToolID:   "costrict",
					ToolName: "CoStrict Commands",
					Frontmatter: map[string]string{
						"proposal": `---
description: "Scaffold a new Spectr change and validate strictly."
argument-hint: feature description or request
---`,
						"apply": `---
description: "Implement an approved Spectr change and keep tasks in sync."
argument-hint: change-id
---`,
						"archive": `---
description: "Archive a deployed Spectr change and update specs."
argument-hint: change-id
---`,
					},
					FilePaths: map[string]string{
						"proposal": ".cospec/spectr/commands/" +
							"spectr-proposal.md",
						"apply": ".cospec/spectr/commands/" +
							"spectr-apply.md",
						"archive": ".cospec/spectr/commands/" +
							"spectr-archive.md",
					},
				})
		},
	)
}
