package templates

import "fmt"

// SlashCommand type for compile-time checked slash command types
type SlashCommand int

const (
	// SlashProposal represents the proposal slash command
	SlashProposal SlashCommand = iota
	// SlashApply represents the apply slash command
	SlashApply
)

// String returns the command name (e.g., "proposal", "apply")
func (sc SlashCommand) String() string {
	switch sc {
	case SlashProposal:
		return "proposal"
	case SlashApply:
		return "apply"
	default:
		return fmt.Sprintf("unknown(%d)", sc)
	}
}

// TemplateName returns the template filename (e.g., "slash-proposal.md.tmpl")
func (sc SlashCommand) TemplateName() string {
	switch sc {
	case SlashProposal:
		return "slash-proposal.md.tmpl"
	case SlashApply:
		return "slash-apply.md.tmpl"
	default:
		return ""
	}
}
