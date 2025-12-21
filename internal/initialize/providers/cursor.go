package providers

import (
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
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
func (p *CursorProvider) Initializers() []types.Initializer {
	proposalPath, applyPath := StandardCommandPaths(
		".cursorrules/commands",
		".md",
	)

	return []types.Initializer{
		initializers.NewSlashCommandsInitializer("proposal", proposalPath, FrontmatterProposal),
		initializers.NewSlashCommandsInitializer("apply", applyPath, FrontmatterApply),
	}
}