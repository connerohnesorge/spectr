package providers

func init() {
	Register(NewWindsurfProvider())
}

// WindsurfProvider implements the Provider interface for Windsurf.
// Windsurf uses .windsurf/commands/ for slash commands (no config file).
type WindsurfProvider struct {
	BaseProvider
}

// NewWindsurfProvider creates a new Windsurf provider.
func NewWindsurfProvider() *WindsurfProvider {
	return &WindsurfProvider{
		BaseProvider: BaseProvider{
			id:            "windsurf",
			name:          "Windsurf",
			priority:      PriorityWindsurf,
			configFile:    "",
			slashDir:      ".windsurf/commands",
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}
