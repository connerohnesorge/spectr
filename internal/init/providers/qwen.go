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
	return &QwenProvider{
		BaseProvider: BaseProvider{
			id:            "qwen",
			name:          "Qwen Code",
			priority:      PriorityQwen,
			configFile:    "QWEN.md",
			slashDir:      ".qwen/commands",
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}
