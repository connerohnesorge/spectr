package providers

func init() {
	Register(NewOpencodeProvider())
}

// OpencodeProvider implements the Provider interface for OpenCode.
// OpenCode uses .opencode/command/spectr/ for slash commands.
// It has no instruction file as it uses JSON configuration.
type OpencodeProvider struct {
	BaseProvider
}

// NewOpencodeProvider creates a new OpenCode provider.
func NewOpencodeProvider() *OpencodeProvider {
	proposalPath, applyPath := StandardCommandPaths(
		".opencode/command", ".md",
	)

	return &OpencodeProvider{
		BaseProvider: BaseProvider{
			id:            "opencode",
			name:          "OpenCode",
			priority:      PriorityOpencode,
			configFile:    "",
			proposalPath:  proposalPath,
			applyPath:     applyPath,
			commandFormat: FormatMarkdown,
			frontmatter:   StandardFrontmatter(),
		},
	}
}
