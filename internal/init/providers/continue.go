package providers

func init() {
	Register(NewContinueProvider())
}

// ContinueProvider implements the Provider interface for Continue.
// Continue uses .continue/commands/ for slash commands (no config file).
type ContinueProvider struct {
	BaseProvider
}

// NewContinueProvider creates a new Continue provider.
func NewContinueProvider() *ContinueProvider {
	return &ContinueProvider{
		BaseProvider: BaseProvider{
			id:            "continue",
			name:          "Continue",
			priority:      PriorityContinue,
			configFile:    "",
			slashDir:      ".continue/commands",
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}
