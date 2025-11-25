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
	return &TabnineProvider{
		BaseProvider: BaseProvider{
			id:            "tabnine",
			name:          "Tabnine",
			priority:      PriorityTabnine,
			configFile:    "",
			slashDir:      ".tabnine/commands",
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}
