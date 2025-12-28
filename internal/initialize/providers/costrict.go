package providers

// CostrictProvider implements the Provider interface for CoStrict.
// CoStrict uses COSTRICT.md and .costrict/commands/ for slash commands.
type CostrictProvider struct {
	BaseProvider
}

// NewCostrictProvider creates a new CoStrict provider.
func NewCostrictProvider() *CostrictProvider {
	proposalPath, applyPath := StandardCommandPaths(
		".costrict/commands",
		".md",
	)

	return &CostrictProvider{
		BaseProvider: BaseProvider{
			id:            "costrict",
			name:          "CoStrict",
			priority:      PriorityCostrict,
			configFile:    "COSTRICT.md",
			proposalPath:  proposalPath,
			applyPath:     applyPath,
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}
