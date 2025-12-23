package providers

import "context"

func init() {
	if err := Register(Registration{
		ID:       "cursor",
		Name:     "Cursor",
		Priority: PriorityCursor,
		Provider: &CursorProvider{},
	}); err != nil {
		panic(err)
	}
}

// CursorProvider implements the Provider interface for Cursor.
type CursorProvider struct{}

func (*CursorProvider) Initializers(_ context.Context) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(".cursorrules/commands/spectr"),
		// No config file for Cursor
		NewSlashCommandsInitializer(
			".cursorrules/commands/spectr",
			".md",
			FormatMarkdown,
		),
	}
}
