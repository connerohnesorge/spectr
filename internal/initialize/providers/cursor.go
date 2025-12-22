package providers

import (
	inits "github.com/connerohnesorge/spectr/internal/initialize/providers/initializers" //nolint:revive
	"github.com/connerohnesorge/spectr/internal/initialize/types"
)

func init() {
	Register(Registration{
		ID:       "cursor",
		Name:     "Cursor",
		Priority: 8,
		Provider: &CursorProvider{},
	})
}

// CursorProvider implements the new Provider interface for Cursor.
type CursorProvider struct{}

// Initializers returns the initializers for Cursor.
func (*CursorProvider) Initializers() []types.Initializer {
	proposalPath, applyPath := StandardCommandPaths(
		".cursorrules/commands",
		".md",
	)

	return []types.Initializer{
		inits.NewSlashCommandsInitializer(
			"proposal",
			proposalPath,
			FrontmatterProposal,
		),
		inits.NewSlashCommandsInitializer(
			"apply",
			applyPath,
			FrontmatterApply,
		),
	}
}