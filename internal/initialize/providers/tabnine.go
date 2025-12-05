package providers

func init() {
	Register(NewTabnineProvider())
}

// TabnineProvider implements the Provider interface for Tabnine.
// Tabnine uses .tabnine/commands/ for slash commands (no config file).
type TabnineProvider struct {
	BaseProvider
}

// NewTabnineProvider returns a configured *TabnineProvider for the Tabnine backend.
// The provider is initialized with id "tabnine", name "Tabnine", priority PriorityTabnine,
// an empty config file, command paths derived from StandardCommandPaths(".tabnine/commands", ".md"),
// a Markdown command format (FormatMarkdown), and frontmatter from StandardFrontmatter().
func NewTabnineProvider() *TabnineProvider {
	proposalPath, applyPath := StandardCommandPaths(
		".tabnine/commands", ".md",
	)

	return &TabnineProvider{
		BaseProvider: BaseProvider{
			id:            "tabnine",
			name:          "Tabnine",
			priority:      PriorityTabnine,
			configFile:    "",
			proposalPath:  proposalPath,
			applyPath:     applyPath,
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}