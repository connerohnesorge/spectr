package providers

func init() {
	Register(NewAntigravityProvider())
}

// AntigravityProvider implements the Provider interface for Antigravity.
// Antigravity uses AGENTS.md and .agent/workflows/ for slash commands.
type AntigravityProvider struct {
	BaseProvider
}

// NewAntigravityProvider creates a new Antigravity provider.
func NewAntigravityProvider() *AntigravityProvider {
	return &AntigravityProvider{
		BaseProvider: BaseProvider{
			id:            "antigravity",
			name:          "Antigravity",
			priority:      PriorityAntigravity,
			configFile:    "AGENTS.md",
			slashDir:      ".agent/workflows",
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}
