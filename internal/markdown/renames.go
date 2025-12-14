package markdown

import (
	"strings"
)

// h3RequirementPfx is the prefix for requirement headers in renamed entries.
const h3ReqPrefix = "### Requirement:"

// Renamed Entry Matchers

// MatchRenamedFrom checks if a line matches the backtick-wrapped FROM format.
// Returns the requirement name and true if matched.
//
// Example input: "- FROM: `### Requirement: Old Name`"
func MatchRenamedFrom(line string) (name string, ok bool) {
	return matchRenamedEntry(line, "FROM:")
}

// MatchRenamedFromAlt checks if a line matches the non-backtick FROM format.
//
// Example input: "- FROM: ### Requirement: Old Name"
func MatchRenamedFromAlt(line string) (name string, ok bool) {
	return matchRenamedEntryAlt(line, "FROM:")
}

// MatchRenamedTo checks if a line matches the backtick-wrapped TO format.
//
// Example input: "- TO: `### Requirement: New Name`"
func MatchRenamedTo(line string) (name string, ok bool) {
	return matchRenamedEntry(line, "TO:")
}

// MatchRenamedToAlt checks if a line matches the non-backtick TO format.
//
// Example input: "- TO: ### Requirement: New Name"
func MatchRenamedToAlt(line string) (name string, ok bool) {
	return matchRenamedEntryAlt(line, "TO:")
}

func matchRenamedEntry(line, prefix string) (name string, ok bool) {
	trimmed := strings.TrimSpace(line)

	if !strings.HasPrefix(trimmed, "- "+prefix) {
		return "", false
	}

	rest := strings.TrimSpace(strings.TrimPrefix(trimmed, "- "+prefix))

	if !strings.HasPrefix(rest, "`") || !strings.HasSuffix(rest, "`") {
		return "", false
	}

	inner := rest[1 : len(rest)-1]

	if !strings.HasPrefix(inner, h3ReqPrefix) {
		return "", false
	}

	name = strings.TrimSpace(strings.TrimPrefix(inner, h3ReqPrefix))

	return name, name != ""
}

func matchRenamedEntryAlt(line, prefix string) (name string, ok bool) {
	trimmed := strings.TrimSpace(line)

	if !strings.HasPrefix(trimmed, "- "+prefix) {
		return "", false
	}

	rest := strings.TrimSpace(strings.TrimPrefix(trimmed, "- "+prefix))

	if !strings.HasPrefix(rest, h3ReqPrefix) {
		if !strings.HasPrefix(rest, "###Requirement:") {
			return "", false
		}

		rest = strings.TrimPrefix(rest, "###Requirement:")
	} else {
		rest = strings.TrimPrefix(rest, h3ReqPrefix)
	}

	name = strings.TrimSpace(rest)

	return name, name != ""
}

// MatchAnyRenamedFrom tries both backtick and non-backtick FROM formats.
func MatchAnyRenamedFrom(line string) (name string, ok bool) {
	if name, ok = MatchRenamedFrom(line); ok {
		return name, true
	}

	return MatchRenamedFromAlt(line)
}

// MatchAnyRenamedTo tries both backtick and non-backtick TO formats.
func MatchAnyRenamedTo(line string) (name string, ok bool) {
	if name, ok = MatchRenamedTo(line); ok {
		return name, true
	}

	return MatchRenamedToAlt(line)
}
