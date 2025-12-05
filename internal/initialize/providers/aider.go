package providers

func init() {
	Register(NewAiderProvider())
}

// AiderProvider implements the Provider interface for Aider.
// Aider uses .aider/commands/ for slash commands (no config file).
type AiderProvider struct {
	BaseProvider
}

// NewAiderProvider creates an AiderProvider configured to load commands from
// ".aider/commands" with Markdown command files and standard frontmatter.
// The provider is identified as "aider", has no config file, and uses PriorityAider.
func NewAiderProvider() *AiderProvider {
	proposalPath, applyPath := StandardCommandPaths(
		".aider/commands", ".md",
	)

	return &AiderProvider{
		BaseProvider: BaseProvider{
			id:            "aider",
			name:          "Aider",
			priority:      PriorityAider,
			configFile:    "",
			proposalPath:  proposalPath,
			applyPath:     applyPath,
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}