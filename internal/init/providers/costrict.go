package providers

func init() {
	Register(NewCostrictProvider())
}

// CostrictProvider implements the Provider interface for CoStrict.
// CoStrict uses COSTRICT.md and .costrict/commands/ for slash commands.
type CostrictProvider struct {
	BaseProvider
}

// NewCostrictProvider creates a new CoStrict provider.
func NewCostrictProvider() *CostrictProvider {
	proposalPath, archivePath, applyPath := StandardCommandPaths(
		".costrict/commands", ".md",
	)

	return &CostrictProvider{
		BaseProvider: BaseProvider{
			id:            "costrict",
			name:          "CoStrict",
			priority:      PriorityCostrict,
			configFile:    "COSTRICT.md",
			proposalPath:  proposalPath,
			archivePath:   archivePath,
			applyPath:     applyPath,
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}
