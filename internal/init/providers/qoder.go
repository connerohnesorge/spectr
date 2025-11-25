package providers

func init() {
	Register(NewQoderProvider())
}

// QoderProvider implements the Provider interface for Qoder.
// Qoder uses QODER.md and .qoder/commands/ for slash commands.
type QoderProvider struct {
	BaseProvider
}

// NewQoderProvider creates a new Qoder provider.
func NewQoderProvider() *QoderProvider {
	return &QoderProvider{
		BaseProvider: BaseProvider{
			id:            "qoder",
			name:          "Qoder",
			priority:      PriorityQoder,
			configFile:    "QODER.md",
			slashDir:      ".qoder/commands",
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}
