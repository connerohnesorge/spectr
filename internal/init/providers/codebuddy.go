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
	proposalPath, archivePath, applyPath := StandardCommandPaths(
		".codebuddy/commands", ".md",
	)

	return &CodeBuddyProvider{
		BaseProvider: BaseProvider{
			id:            "codebuddy",
			name:          "CodeBuddy",
			priority:      PriorityCodeBuddy,
			configFile:    "CODEBUDDY.md",
			proposalPath:  proposalPath,
			archivePath:   archivePath,
			applyPath:     applyPath,
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}
