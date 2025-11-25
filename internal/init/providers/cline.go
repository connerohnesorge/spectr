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
	return &ClineProvider{
		BaseProvider: BaseProvider{
			id:            "cline",
			name:          "Cline",
			priority:      PriorityCline,
			configFile:    "CLINE.md",
			slashDir:      ".clinerules/commands",
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}
