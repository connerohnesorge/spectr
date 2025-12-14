package regex

import "regexp"

// Rename patterns for RENAMED Requirements section parsing.
// Both backtick-wrapped and non-backtick variants are supported.
var (
	// RenamedFrom matches "- FROM: `### Requirement: Name`" (with backticks).
	// Captures the requirement name.
	RenamedFrom = regexp.MustCompile(
		`^-\s*FROM:\s*` + "`" + `###\s+Requirement:\s*(.+?)` + "`" + `\s*$`,
	)

	// RenamedTo matches "- TO: `### Requirement: Name`" (with backticks).
	// Captures the requirement name.
	RenamedTo = regexp.MustCompile(
		`^-\s*TO:\s*` + "`" + `###\s+Requirement:\s*(.+?)` + "`" + `\s*$`,
	)

	// RenamedFromAlt matches "- FROM: ### Requirement: Name"
	// (without backticks). Captures the requirement name.
	RenamedFromAlt = regexp.MustCompile(
		`^\s*-\s*FROM:\s*###\s*Requirement:\s*(.+)$`,
	)

	// RenamedToAlt matches "- TO: ### Requirement: Name"
	// (without backticks). Captures the requirement name.
	RenamedToAlt = regexp.MustCompile(
		`^\s*-\s*TO:\s*###\s*Requirement:\s*(.+)$`,
	)
)

// MatchRenamedFrom checks if a line matches the backtick-wrapped FROM
// format. Returns the requirement name and true if matched, or empty
// string and false otherwise.
//
// Example input: "- FROM: `### Requirement: Old Name`"
// Returns: "Old Name", true
func MatchRenamedFrom(line string) (name string, ok bool) {
	matches := RenamedFrom.FindStringSubmatch(line)
	if len(matches) < 2 {
		return "", false
	}

	return matches[1], true
}

// MatchRenamedFromAlt checks if a line matches the non-backtick FROM
// format. Returns the requirement name and true if matched, or empty
// string and false otherwise.
//
// Example input: "- FROM: ### Requirement: Old Name"
// Returns: "Old Name", true
func MatchRenamedFromAlt(line string) (name string, ok bool) {
	matches := RenamedFromAlt.FindStringSubmatch(line)
	if len(matches) < 2 {
		return "", false
	}

	return matches[1], true
}

// MatchRenamedTo checks if a line matches the backtick-wrapped TO
// format. Returns the requirement name and true if matched, or empty
// string and false otherwise.
//
// Example input: "- TO: `### Requirement: New Name`"
// Returns: "New Name", true
func MatchRenamedTo(line string) (name string, ok bool) {
	matches := RenamedTo.FindStringSubmatch(line)
	if len(matches) < 2 {
		return "", false
	}

	return matches[1], true
}

// MatchRenamedToAlt checks if a line matches the non-backtick TO format.
// Returns the requirement name and true if matched, or empty string and
// false otherwise.
//
// Example input: "- TO: ### Requirement: New Name"
// Returns: "New Name", true
func MatchRenamedToAlt(line string) (name string, ok bool) {
	matches := RenamedToAlt.FindStringSubmatch(line)
	if len(matches) < 2 {
		return "", false
	}

	return matches[1], true
}

// MatchAnyRenamedFrom tries both backtick and non-backtick FROM formats.
// Returns the requirement name and true if either format matches.
// Tries backtick format first, then falls back to non-backtick.
func MatchAnyRenamedFrom(line string) (name string, ok bool) {
	if name, ok = MatchRenamedFrom(line); ok {
		return name, true
	}

	return MatchRenamedFromAlt(line)
}

// MatchAnyRenamedTo tries both backtick and non-backtick TO formats.
// Returns the requirement name and true if either format matches.
// Tries backtick format first, then falls back to non-backtick.
func MatchAnyRenamedTo(line string) (name string, ok bool) {
	if name, ok = MatchRenamedTo(line); ok {
		return name, true
	}

	return MatchRenamedToAlt(line)
}
