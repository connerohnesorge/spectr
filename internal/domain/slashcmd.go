// Package domain contains shared domain types that can be used across packages
// without creating import cycles. This includes types like TemplateRef, SlashCommand,
// and TemplateContext that are used by both initialize/templates and initialize/providers.
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
