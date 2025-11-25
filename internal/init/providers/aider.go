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
	proposalPath, archivePath, applyPath := StandardCommandPaths(
		".aider/commands", ".md",
	)

	return &AiderProvider{
		BaseProvider: BaseProvider{
			id:            "aider",
			name:          "Aider",
			priority:      PriorityAider,
			configFile:    "",
			proposalPath:  proposalPath,
			archivePath:   archivePath,
			applyPath:     applyPath,
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}
