package providers

// init registers the CodeBuddy provider with the global registry.
func init() {
	Register(&CodeBuddyProvider{})
}

// CodeBuddyProvider implements the Provider interface for CodeBuddy.
// CodeBuddy uses CODEBUDDY.md and .codebuddy/commands/ for
// slash commands.
type CodeBuddyProvider struct{}

// ID returns the unique identifier for the CodeBuddy provider.
func (*CodeBuddyProvider) ID() string { return "codebuddy" }

// Name returns the display name for CodeBuddy.
func (*CodeBuddyProvider) Name() string { return "CodeBuddy" }

// Priority returns the display order for CodeBuddy.
func (*CodeBuddyProvider) Priority() int { return PriorityCodeBuddy }

// Initializers returns the file initializers for the CodeBuddy provider.
func (*CodeBuddyProvider) Initializers() []FileInitializer {
	proposalPath, applyPath := StandardCommandPaths(
		".codebuddy/commands",
		".md",
	)

	return []FileInitializer{
		NewInstructionFileInitializer("CODEBUDDY.md"),
		NewMarkdownSlashCommandInitializer(
			proposalPath,
			"proposal",
			FrontmatterProposal,
		),
		NewMarkdownSlashCommandInitializer(
			applyPath,
			"apply",
			FrontmatterApply,
		),
	}
}

func (p *CodeBuddyProvider) IsConfigured(projectPath string) bool {
	return AreInitializersConfigured(p.Initializers(), projectPath)
}

func (p *CodeBuddyProvider) GetFilePaths() []string {
	return GetInitializerPaths(p.Initializers())
}
