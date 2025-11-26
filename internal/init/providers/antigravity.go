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
			proposalPath:  ".agent/workflows/spectr-proposal.md",
			syncPath:      ".agent/workflows/spectr-sync.md",
			applyPath:     ".agent/workflows/spectr-apply.md",
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}
