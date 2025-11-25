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
	return &CostrictProvider{
		BaseProvider: BaseProvider{
			id:            "costrict",
			name:          "CoStrict",
			priority:      PriorityCostrict,
			configFile:    "COSTRICT.md",
			slashDir:      ".costrict/commands",
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}
