package providers

func init() {
	Register(NewCursorProvider())
}

// CursorProvider implements the Provider interface for Cursor.
// Cursor uses .cursorrules/commands/ for slash commands (no config file).
type CursorProvider struct {
	BaseProvider
}

// NewCursorProvider creates a new Cursor provider.
func NewCursorProvider() *CursorProvider {
	proposalPath, applyPath := StandardCommandPaths(
		".cursorrules/commands",
		".md",
	)

	return &CursorProvider{
		BaseProvider: BaseProvider{
			id:            "cursor",
			name:          "Cursor",
			priority:      PriorityCursor,
			configFile:    "",
			proposalPath:  proposalPath,
			applyPath:     applyPath,
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}
