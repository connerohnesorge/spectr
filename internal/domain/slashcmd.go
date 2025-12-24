// Package domain contains shared domain types used across packages.
package domain

import "fmt"

// SlashCommand represents a type-safe slash command identifier.
type SlashCommand int

const (
	SlashProposal SlashCommand = iota
	SlashApply
)

// templateNames maps slash commands to their template file names.
var templateNames = map[SlashCommand]string{
	SlashProposal: "slash-proposal.md.tmpl",
	SlashApply:    "slash-apply.md.tmpl",
}

// String returns the command name for debugging.
func (s SlashCommand) String() string {
	names := []string{"proposal", "apply"}
	if int(s) < len(names) {
		return names[s]
	}

	return "unknown"
}

// TemplateName returns the template file name for this command.
// Returns an error if the command is not recognized.
func (s SlashCommand) TemplateName() (string, error) {
	name, ok := templateNames[s]
	if !ok {
		return "", fmt.Errorf(
			"unknown slash command: %d",
			s,
		)
	}

	return name, nil
}
