package providers

func init() {
	Register(NewTabnineProvider())
}

// TabnineProvider implements the Provider interface for Tabnine.
// Tabnine uses .tabnine/commands/ for slash commands (no config file).
type TabnineProvider struct {
	BaseProvider
}

// NewTabnineProvider creates a new Tabnine provider.
func NewTabnineProvider() *TabnineProvider {
	proposalPath, syncPath, applyPath := StandardCommandPaths(
		".tabnine/commands", ".md",
	)

	return &TabnineProvider{
		BaseProvider: BaseProvider{
			id:            "tabnine",
			name:          "Tabnine",
			priority:      PriorityTabnine,
			configFile:    "",
			proposalPath:  proposalPath,
			syncPath:      syncPath,
			applyPath:     applyPath,
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}
