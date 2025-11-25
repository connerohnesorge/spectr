package providers

func init() {
	Register(NewSmolProvider())
}

// SmolProvider implements the Provider interface for Smol.
// Smol uses .smol/commands/ for slash commands (no config file).
type SmolProvider struct {
	BaseProvider
}

// NewSmolProvider creates a new Smol provider.
func NewSmolProvider() *SmolProvider {
	return &SmolProvider{
		BaseProvider: BaseProvider{
			id:            "smol",
			name:          "Smol",
			priority:      PrioritySmol,
			configFile:    "",
			slashDir:      ".smol/commands",
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}
