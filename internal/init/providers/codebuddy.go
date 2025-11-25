package providers

func init() {
	Register(NewCodeBuddyProvider())
}

// CodeBuddyProvider implements the Provider interface for CodeBuddy.
// CodeBuddy uses CODEBUDDY.md and .codebuddy/commands/ for slash commands.
type CodeBuddyProvider struct {
	BaseProvider
}

// NewCodeBuddyProvider creates a new CodeBuddy provider.
func NewCodeBuddyProvider() *CodeBuddyProvider {
	return &CodeBuddyProvider{
		BaseProvider: BaseProvider{
			id:            "codebuddy",
			name:          "CodeBuddy",
			priority:      PriorityCodeBuddy,
			configFile:    "CODEBUDDY.md",
			slashDir:      ".codebuddy/commands",
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}
