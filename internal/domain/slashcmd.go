package domain

// SlashCommand represents a type-safe slash command identifier.
type SlashCommand int

const (
	SlashProposal SlashCommand = iota
	SlashApply
	SlashNext
)

// String returns the command name for debugging.
func (s SlashCommand) String() string {
	names := []string{"proposal", "apply", "next"}
	if int(s) < len(names) {
		return names[s]
	}

	return "unknown"
}
