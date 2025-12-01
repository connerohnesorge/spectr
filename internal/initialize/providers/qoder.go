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
	proposalPath, syncPath, applyPath := StandardCommandPaths(
		".qoder/commands", ".md",
	)

	return &QoderProvider{
		BaseProvider: BaseProvider{
			id:            "qoder",
			name:          "Qoder",
			priority:      PriorityQoder,
			configFile:    "QODER.md",
			proposalPath:  proposalPath,
			syncPath:      syncPath,
			applyPath:     applyPath,
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}
