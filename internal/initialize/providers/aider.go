package providers

func init() {
	Register(NewAiderProvider())
}

// AiderProvider implements the Provider interface for Aider.
// Aider uses .aider/commands/ for slash commands (no config file).
type AiderProvider struct {
	BaseProvider
}

// NewAiderProvider creates a new Aider provider.
func NewAiderProvider() *AiderProvider {
	proposalPath, applyPath := StandardCommandPaths(
		".aider/commands",
		".md",
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
