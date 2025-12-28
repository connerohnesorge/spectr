package providers

// CrushProvider implements the Provider interface for Crush.
// Crush uses CRUSH.md for instructions and .crush/commands/ for slash commands.
type CrushProvider struct {
	BaseProvider
}

// NewCrushProvider creates a new Crush provider.
func NewCrushProvider() *CrushProvider {
	proposalPath, applyPath := StandardCommandPaths(
		".crush/commands",
		".md",
	)

	return &CrushProvider{
		BaseProvider: BaseProvider{
			id:            "crush",
			name:          "Crush",
			priority:      PriorityCrush,
			configFile:    "CRUSH.md",
			proposalPath:  proposalPath,
			applyPath:     applyPath,
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}
