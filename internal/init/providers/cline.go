package providers

func init() {
	Register(NewClineProvider())
}

// ClineProvider implements the Provider interface for Cline.
// Cline uses CLINE.md and .clinerules/commands/ for slash commands.
type ClineProvider struct {
	BaseProvider
}

// NewClineProvider creates a new Cline provider.
func NewClineProvider() *ClineProvider {
	proposalPath, syncPath, applyPath := StandardCommandPaths(
		".clinerules/commands", ".md",
	)

	return &ClineProvider{
		BaseProvider: BaseProvider{
			id:            "cline",
			name:          "Cline",
			priority:      PriorityCline,
			configFile:    "CLINE.md",
			proposalPath:  proposalPath,
			syncPath:      syncPath,
			applyPath:     applyPath,
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}
