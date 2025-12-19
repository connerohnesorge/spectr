package providers

import "context"

func init() {
	err := Register(Registration{
		ID:       "cursor",
		Name:     "Cursor",
		Priority: PriorityCursor,
		Provider: &CursorProvider{},
	})
	if err != nil {
		panic(err)
	}
}

// CursorProvider implements the Provider interface for Cursor.
// Cursor uses .cursorrules/commands/spectr/ for slash commands (no config file).
type CursorProvider struct{}

// Initializers returns the initializers for Cursor.
func (p *CursorProvider) Initializers(ctx context.Context) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".cursorrules/commands/spectr"),
		NewSlashCommandsInitializerWithFrontmatter(
			".cursorrules/commands/spectr",
			".md",
			FormatMarkdown,
			map[string]string{
				"proposal": FrontmatterProposal,
				"apply":    FrontmatterApply,
			},
		),
	}
}
