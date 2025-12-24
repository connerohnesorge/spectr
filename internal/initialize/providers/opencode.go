package providers

import (
	"context"

	"github.com/connerohnesorge/spectr/internal/initialize/templates"
)

func init() {
	_ = RegisterProvider(Registration{
		ID:       "opencode",
		Name:     "OpenCode",
		Priority: PriorityOpencode,
		Provider: &OpencodeProvider{},
	})
}

// OpencodeProvider implements the Provider interface for OpenCode.
// OpenCode uses .opencode/command/spectr/ for slash commands.
// It has no instruction file as it uses JSON configuration.
type OpencodeProvider struct{}

func (*OpencodeProvider) Initializers(
	_ context.Context,
) []Initializer {
	return []Initializer{
		NewDirectoryInitializer(
			".opencode/command/spectr",
		),
		NewSlashCommandsInitializer(
			".opencode/command/spectr",
			".md",
			[]templates.SlashCommand{
				templates.SlashProposal,
				templates.SlashApply,
			},
		),
	}
}
