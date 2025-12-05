package providers

func init() {
	Register(NewContinueProvider())
}

// ContinueProvider implements the Provider interface for Continue.
// Continue uses .continue/commands/ for slash commands (no config file).
type ContinueProvider struct {
	BaseProvider
}

// NewContinueProvider returns a ContinueProvider configured for the "continue" command set.
// The provider has id "continue", name "Continue", priority PriorityContinue, an empty
// configFile, commandFormat FormatMarkdown, frontmatter from StandardFrontmatter(), and
// proposalPath and applyPath derived from the standard command paths for ".continue/commands"
// with ".md" extension.
func NewContinueProvider() *ContinueProvider {
	proposalPath, applyPath := StandardCommandPaths(
		".continue/commands", ".md",
	)

	return &ContinueProvider{
		BaseProvider: BaseProvider{
			id:            "continue",
			name:          "Continue",
			priority:      PriorityContinue,
			configFile:    "",
			proposalPath:  proposalPath,
			applyPath:     applyPath,
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}