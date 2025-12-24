// Package domain contains shared domain types used across packages.
package domain

// SlashCommand represents a type-safe slash command identifier.
type SlashCommand int

const (
	SlashProposal SlashCommand = iota
	SlashApply
)

// String returns the command name for debugging.
func (s SlashCommand) String() string {
	names := []string{"proposal", "apply"}
	if int(s) < len(names) {
		return names[s]
	}

	return "unknown"
}

// TemplateName returns the template file name for this command.
func (s SlashCommand) TemplateName() string {
	names := map[SlashCommand]string{
		SlashProposal: "slash-proposal.md.tmpl",
		SlashApply:    "slash-apply.md.tmpl",
	}

	return names[s]
}
