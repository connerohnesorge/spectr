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
	proposalPath, syncPath, applyPath := PrefixedCommandPaths(
		".agent/workflows", ".md",
	)

	return &AntigravityProvider{
		BaseProvider: BaseProvider{
			id:            "antigravity",
			name:          "Antigravity",
			priority:      PriorityAntigravity,
			configFile:    "AGENTS.md",
			proposalPath:  proposalPath,
			syncPath:      syncPath,
			applyPath:     applyPath,
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}
