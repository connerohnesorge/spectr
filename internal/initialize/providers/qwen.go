package providers

func init() {
	Register(NewQwenProvider())
}

// QwenProvider implements the Provider interface for Qwen Code.
// Qwen uses QWEN.md and .qwen/commands/ for slash commands.
type QwenProvider struct {
	BaseProvider
}

// NewQwenProvider constructs a QwenProvider preconfigured for Qwen Code commands.
// The provider is initialized with id "qwen", name "Qwen Code", PriorityQwen, config file "QWEN.md", proposal and apply paths derived from StandardCommandPaths(".qwen/commands", ".md"), markdown command format, and standard frontmatter.
func NewQwenProvider() *QwenProvider {
	proposalPath, applyPath := StandardCommandPaths(
		".qwen/commands", ".md",
	)

	return &QwenProvider{
		BaseProvider: BaseProvider{
			id:            "qwen",
			name:          "Qwen Code",
			priority:      PriorityQwen,
			configFile:    "QWEN.md",
			proposalPath:  proposalPath,
			applyPath:     applyPath,
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}