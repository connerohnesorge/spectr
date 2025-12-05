package providers

func init() {
	Register(NewCursorProvider())
}

// CursorProvider implements the Provider interface for Cursor.
// Cursor uses .cursorrules/commands/ for slash commands (no config file).
type CursorProvider struct {
	BaseProvider
}

// NewCursorProvider returns a CursorProvider configured for the "cursor" provider with id "cursor", display name "Cursor", priority PriorityCursor, no config file, proposal and apply command paths under ".cursorrules/commands" for ".md" files, Markdown command format, and standard frontmatter.
func NewCursorProvider() *CursorProvider {
	proposalPath, applyPath := StandardCommandPaths(
		".cursorrules/commands", ".md",
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