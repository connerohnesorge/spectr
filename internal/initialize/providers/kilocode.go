package providers

func init() {
	Register(NewKilocodeProvider())
}

// KilocodeProvider implements the Provider interface for Kilocode.
// Kilocode uses .kilocode/commands/ for slash commands (no config file).
type KilocodeProvider struct {
	BaseProvider
}

// NewKilocodeProvider returns a KilocodeProvider configured for Kilocode-specific slash commands located under ".kilocode/commands" using Markdown command files and no separate config file.
// The provider's proposal and apply paths are derived from StandardCommandPaths and its frontmatter is initialized with StandardFrontmatter.
func NewKilocodeProvider() *KilocodeProvider {
	proposalPath, applyPath := StandardCommandPaths(
		".kilocode/commands", ".md",
	)

	return &KilocodeProvider{
		BaseProvider: BaseProvider{
			id:            "kilocode",
			name:          "Kilocode",
			priority:      PriorityKilocode,
			configFile:    "",
			proposalPath:  proposalPath,
			applyPath:     applyPath,
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}