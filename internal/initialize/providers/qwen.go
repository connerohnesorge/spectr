package providers

func init() {
	Register(NewQwenProvider())
}

// QwenProvider implements the Provider interface for Qwen Code.
// Qwen uses QWEN.md and .qwen/commands/ for slash commands.
type QwenProvider struct {
	BaseProvider
}

// NewQwenProvider creates a new Qwen Code provider.
func NewQwenProvider() *QwenProvider {
	proposalPath, applyPath := StandardCommandPaths(
		".qwen/commands",
		".md",
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
