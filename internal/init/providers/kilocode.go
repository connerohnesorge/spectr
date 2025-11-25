package providers

func init() {
	Register(NewKilocodeProvider())
}

// KilocodeProvider implements the Provider interface for Kilocode.
// Kilocode uses .kilocode/commands/ for slash commands (no config file).
type KilocodeProvider struct {
	BaseProvider
}

// NewKilocodeProvider creates a new Kilocode provider.
func NewKilocodeProvider() *KilocodeProvider {
	return &KilocodeProvider{
		BaseProvider: BaseProvider{
			id:            "kilocode",
			name:          "Kilocode",
			priority:      PriorityKilocode,
			configFile:    "",
			slashDir:      ".kilocode/commands",
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}
