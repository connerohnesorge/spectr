package providers

func init() {
	Register(NewAntigravityProvider())
}

// AntigravityProvider implements the Provider interface for Antigravity.
// Antigravity uses AGENTS.md and .agent/workflows/ for slash commands.
type AntigravityProvider struct {
	BaseProvider
}

// NewAntigravityProvider creates and returns an AntigravityProvider configured for Antigravity integrations.
// The provider is initialized with id "antigravity", name "Antigravity", PriorityAntigravity, config file "AGENTS.md", Markdown command format, standard frontmatter, and proposal/apply paths derived from ".agent/workflows" and ".md".
func NewAntigravityProvider() *AntigravityProvider {
	proposalPath, applyPath := PrefixedCommandPaths(
		".agent/workflows", ".md",
	)

	return &AntigravityProvider{
		BaseProvider: BaseProvider{
			id:            "antigravity",
			name:          "Antigravity",
			priority:      PriorityAntigravity,
			configFile:    "AGENTS.md",
			proposalPath:  proposalPath,
			applyPath:     applyPath,
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}