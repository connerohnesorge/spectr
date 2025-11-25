package providers

func init() {
	Register(NewMentatProvider())
}

// MentatProvider implements the Provider interface for Mentat.
// Mentat uses .mentat/commands/ for slash commands (no config file).
type MentatProvider struct {
	BaseProvider
}

// NewMentatProvider creates a new Mentat provider.
func NewMentatProvider() *MentatProvider {
	return &MentatProvider{
		BaseProvider: BaseProvider{
			id:            "mentat",
			name:          "Mentat",
			priority:      PriorityMentat,
			configFile:    "",
			slashDir:      ".mentat/commands",
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}
