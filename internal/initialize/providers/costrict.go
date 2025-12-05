package providers

func init() {
	Register(NewCostrictProvider())
}

// CostrictProvider implements the Provider interface for CoStrict.
// CoStrict uses COSTRICT.md and .costrict/commands/ for slash commands.
type CostrictProvider struct {
	BaseProvider
}

// NewCostrictProvider constructs a *CostrictProvider configured for CoStrict command files.
// The provider is initialized with id "costrict", name "CoStrict", priority PriorityCostrict,
// config file "COSTRICT.md", proposal and apply paths from StandardCommandPaths, markdown
// command format, and frontmatter from StandardFrontmatter().
func NewCostrictProvider() *CostrictProvider {
	proposalPath, applyPath := StandardCommandPaths(
		".costrict/commands", ".md",
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