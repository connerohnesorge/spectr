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
	proposalPath, archivePath, applyPath := StandardCommandPaths(
		".kilocode/commands", ".md",
	)

	return &KilocodeProvider{
		BaseProvider: BaseProvider{
			id:            "kilocode",
			name:          "Kilocode",
			priority:      PriorityKilocode,
			configFile:    "",
			proposalPath:  proposalPath,
			archivePath:   archivePath,
			applyPath:     applyPath,
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}
