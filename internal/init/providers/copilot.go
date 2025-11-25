package providers

func init() {
	Register(NewCopilotProvider())
}

// CopilotProvider implements the Provider interface for GitHub Copilot.
// Copilot uses .github/copilot/commands/ for slash commands (no config file).
type CopilotProvider struct {
	BaseProvider
}

// NewCopilotProvider creates a new GitHub Copilot provider.
func NewCopilotProvider() *CopilotProvider {
	return &CopilotProvider{
		BaseProvider: BaseProvider{
			id:            "copilot",
			name:          "Copilot",
			priority:      PriorityCopilot,
			configFile:    ".github/copilot-instructions.md",
			slashDir:      ".github/copilot/commands",
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}
