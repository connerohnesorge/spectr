package providers

// init registers the Cursor provider with the global registry.
func init() {
	Register(&CursorProvider{})
}

// CursorProvider implements the Provider interface for Cursor.
// Cursor uses .cursorrules/commands/ for slash commands.
// It does not use a separate instruction file.
type CursorProvider struct{}

// ID returns the unique identifier for the Cursor provider.
func (*CursorProvider) ID() string { return "cursor" }

// Name returns the display name for Cursor.
func (*CursorProvider) Name() string { return "Cursor" }

// Priority returns the display order for Cursor.
func (*CursorProvider) Priority() int { return PriorityCursor }

// Initializers returns the file initializers for the Cursor provider.
func (*CursorProvider) Initializers() []FileInitializer {
	proposalPath, applyPath := StandardCommandPaths(
		".cursorrules/commands",
		".md",
	)

	return []FileInitializer{
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

func (p *CursorProvider) IsConfigured(projectPath string) bool {
	return AreInitializersConfigured(p.Initializers(), projectPath)
}

func (p *CursorProvider) GetFilePaths() []string {
	return GetInitializerPaths(p.Initializers())
}
