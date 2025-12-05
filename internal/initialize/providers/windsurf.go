package providers

func init() {
	Register(NewWindsurfProvider())
}

// WindsurfProvider implements the Provider interface for Windsurf.
// Windsurf uses .windsurf/commands/ for slash commands (no config file).
type WindsurfProvider struct {
	BaseProvider
}

// NewWindsurfProvider returns a WindsurfProvider configured with the Windsurf provider defaults.
//
// The provider is initialized with the id "windsurf", name "Windsurf", priority PriorityWindsurf,
// proposal and apply command paths located under ".windsurf/commands" using ".md" files,
// the Markdown command format, and the standard frontmatter.
func NewWindsurfProvider() *WindsurfProvider {
	proposalPath, applyPath := StandardCommandPaths(
		".windsurf/commands", ".md",
	)

	return &WindsurfProvider{
		BaseProvider: BaseProvider{
			id:            "windsurf",
			name:          "Windsurf",
			priority:      PriorityWindsurf,
			configFile:    "",
			proposalPath:  proposalPath,
			applyPath:     applyPath,
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}