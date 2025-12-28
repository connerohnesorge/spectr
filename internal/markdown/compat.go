// Package markdown provides a markdown parser and AST for spectr specifications.
//
// compat.go provides line-by-line matchers for compatibility with the old
// regex package API. These functions work on single lines and are efficient
// for line-by-line processing during migration from regex-based parsing.
//
// These helpers use simple string operations instead of regex for better
// performance and maintainability.
//
//nolint:revive // file-length-limit: compatibility layer requires many matcher functions
package markdown

import (
	"bytes"
	"strings"
	"unicode"
)

// MatchRequirementHeader checks if a line is a requirement header
// ("### Requirement: Name") and extracts the name.
// Returns the requirement name and true if matched, or empty string and false.
// This function allows flexible whitespace between ### and "Requirement:"
// and after the colon, matching the behavior of the original regex.
//
// Example:
//
//	name, ok := MatchRequirementHeader("### Requirement: User Authentication")
//	// name = "User Authentication", ok = true
func MatchRequirementHeader(
	line string,
) (name string, ok bool) {
	// Must start with "###"
	if !strings.HasPrefix(line, "###") {
		return "", false
	}

	// Skip "###" and any whitespace
	rest := strings.TrimLeft(line[3:], " \t")

	// Must be followed by "Requirement:"
	if !strings.HasPrefix(rest, "Requirement:") {
		return "", false
	}

	// Skip "Requirement:" and any whitespace
	rest = strings.TrimLeft(rest[12:], " \t")

	if rest == "" {
		return "", false
	}

	return rest, true
}

// MatchScenarioHeader checks if a line is a scenario header
// ("#### Scenario: Name") and extracts the name.
// Returns the scenario name and true if matched, or empty string and false.
// This function allows flexible whitespace between #### and "Scenario:"
// and after the colon, matching the behavior of the original regex.
//
// Example:
//
//	name, ok := MatchScenarioHeader("#### Scenario: User logs in")
//	// name = "User logs in", ok = true
func MatchScenarioHeader(
	line string,
) (name string, ok bool) {
	// Must start with "####"
	if !strings.HasPrefix(line, "####") {
		return "", false
	}

	// Skip "####" and any whitespace
	rest := strings.TrimLeft(line[4:], " \t")

	// Must be followed by "Scenario:"
	if !strings.HasPrefix(rest, "Scenario:") {
		return "", false
	}

	// Skip "Scenario:" and any whitespace
	rest = strings.TrimLeft(rest[9:], " \t")

	if rest == "" {
		return "", false
	}

	return rest, true
}

// IsH2Header checks if a line starts with "## " (any H2 header).
func IsH2Header(line string) bool {
	return strings.HasPrefix(line, "## ")
}

// IsH3Header checks if a line starts with "### " (any H3 header).
func IsH3Header(line string) bool {
	return strings.HasPrefix(line, "### ")
}

// IsH4Header checks if a line starts with "#### " (any H4 header).
func IsH4Header(line string) bool {
	return strings.HasPrefix(line, "#### ")
}

// MatchTaskCheckbox checks if a line contains a task checkbox and
// extracts the state. The line must start with optional whitespace,
// followed by "- [x]" or "- [ ]".
// Returns the checkbox state (' ' for unchecked, 'x' or 'X' for checked)
// and true if matched, or 0 and false otherwise.
//
// Example:
//
//	state, ok := MatchTaskCheckbox("- [ ] Unchecked task")
//	// state = ' ', ok = true
//
//	state, ok := MatchTaskCheckbox("- [x] Checked task")
//	// state = 'x', ok = true
func MatchTaskCheckbox(
	line string,
) (status rune, ok bool) {
	// Skip leading whitespace
	trimmed := strings.TrimLeft(line, " \t")

	// Must start with "- ["
	if !strings.HasPrefix(trimmed, "- [") {
		return 0, false
	}

	// Check if we have enough characters for "- [x]"
	if len(trimmed) < 5 {
		return 0, false
	}

	// Extract the checkbox character
	checkChar := trimmed[3]

	// Must be space, x, or X
	if checkChar != ' ' && checkChar != 'x' &&
		checkChar != 'X' {
		return 0, false
	}

	// Must be followed by "]"
	if trimmed[4] != ']' {
		return 0, false
	}

	return rune(checkChar), true
}

// NumberedTaskMatch holds the parsed result of a numbered task line.
type NumberedTaskMatch struct {
	Section string // The section number (e.g., "1")
	Number  string // The task number (e.g., "1.1")
	Status  rune   // ' ' for unchecked, 'x' or 'X' for checked
	Content string // The task description
}

// MatchNumberedTask parses a numbered task line from tasks.md format.
// Format: "- [ ] 1.1 Task description" or "- [x] 2.3 Another task"
// Also accepts simpler format: "- [ ] 1. Task description" (no digits after dot)
// Returns the parsed task match and true if matched, or nil and false.
//
// Example:
//
//	match, ok := MatchNumberedTask("- [ ] 1.1 Create the parser")
//	// match.Number = "1.1", match.Status = ' ', match.Content = "Create the parser"
//
//	match, ok := MatchNumberedTask("- [ ] 1. Simple task")
//	// match.Number = "1.", match.Status = ' ', match.Content = "Simple task"
func MatchNumberedTask(
	line string,
) (*NumberedTaskMatch, bool) {
	// Must start with "- ["
	if !strings.HasPrefix(line, "- [") {
		return nil, false
	}

	// Check minimum length for "- [x] 1.1 X"
	if len(line) < 11 {
		return nil, false
	}

	// Extract checkbox state
	checkChar := line[3]
	if checkChar != ' ' && checkChar != 'x' &&
		checkChar != 'X' {
		return nil, false
	}

	// Must be followed by "] "
	if line[4] != ']' || line[5] != ' ' {
		return nil, false
	}

	// Rest of line after "- [x] "
	rest := line[6:]

	// Parse the task number (e.g., "1.1", "12.34", or "1.")
	// Must be digits, dot, optionally more digits
	numEnd := 0
	dotPos := -1
	dotSeen := false

parseLoop:
	for i, c := range rest {
		switch {
		case c >= '0' && c <= '9':
			numEnd = i + 1
		case c == '.' && !dotSeen:
			dotSeen = true
			dotPos = i
		default:
			break parseLoop
		}
	}

	// If dot was seen but no digits after it, include the dot in numEnd
	// This handles "1. Task" format where numEnd would be 1 but dot is at position 1
	if dotSeen && dotPos >= 0 &&
		dotPos >= numEnd {
		numEnd = dotPos + 1
	}

	// Validate we got a proper number format (digits followed by dot, optionally more digits)
	// Accept both "1.1" and "1." formats
	if !dotSeen || numEnd == 0 {
		return nil, false
	}

	taskNum := rest[:numEnd]

	// Must be followed by space and content
	if numEnd >= len(rest) ||
		rest[numEnd] != ' ' {
		return nil, false
	}

	content := strings.TrimLeft(
		rest[numEnd+1:],
		" \t",
	)
	if content == "" {
		return nil, false
	}

	// Extract section number (part before the dot)
	dotIdx := strings.Index(taskNum, ".")
	section := ""
	if dotIdx > 0 {
		section = taskNum[:dotIdx]
	}

	return &NumberedTaskMatch{
		Section: section,
		Number:  taskNum,
		Status:  rune(checkChar),
		Content: content,
	}, true
}

// FlexibleTaskMatch holds the parsed result of any task checkbox line.
// Unlike NumberedTaskMatch, this accepts tasks with or without numbers.
type FlexibleTaskMatch struct {
	Number  string // Explicit number if present (e.g., "1.1", "1.", "1"), empty otherwise
	Status  rune   // ' ' for unchecked, 'x' or 'X' for checked
	Content string // The task description
}

// MatchFlexibleTask parses any task checkbox line from tasks.md format.
// Accepts all formats:
//   - "- [ ] 1.1 Task description" (decimal)
//   - "- [ ] 1. Task description" (simple dot)
//   - "- [ ] 1 Task description" (number only)
//   - "- [ ] Task description" (no number)
//
// Returns the parsed task match and true if matched, or nil and false.
func MatchFlexibleTask(
	line string,
) (*FlexibleTaskMatch, bool) {
	// Must start with "- ["
	if !strings.HasPrefix(line, "- [") {
		return nil, false
	}

	// Need at least "- [x] X" (7 chars)
	if len(line) < 7 {
		return nil, false
	}

	// Extract checkbox state
	checkChar := line[3]
	if checkChar != ' ' && checkChar != 'x' &&
		checkChar != 'X' {
		return nil, false
	}

	// Must be followed by "] "
	if line[4] != ']' || line[5] != ' ' {
		return nil, false
	}

	// Rest of line after "- [x] "
	rest := line[6:]
	if rest == "" {
		return nil, false
	}

	// Try to parse optional number prefix
	// Number can be: digits, optionally followed by dot, optionally followed by more digits
	var number string
	var contentStart int

	// Check if starts with digit
	if rest != "" && rest[0] >= '0' &&
		rest[0] <= '9' {
		// Parse number: digits, optional dot, optional more digits
		i := 0
		for i < len(rest) && rest[i] >= '0' && rest[i] <= '9' {
			i++
		}
		// Check for optional dot
		if i < len(rest) && rest[i] == '.' {
			i++
			// Check for optional digits after dot
			for i < len(rest) && rest[i] >= '0' && rest[i] <= '9' {
				i++
			}
		}
		// Number must be followed by space to be valid
		if i < len(rest) && rest[i] == ' ' {
			number = rest[:i]
			contentStart = i + 1
		}
	}

	// Extract content
	var content string
	if number != "" {
		content = strings.TrimLeft(
			rest[contentStart:],
			" \t",
		)
	} else {
		content = strings.TrimLeft(rest, " \t")
	}

	if content == "" {
		return nil, false
	}

	return &FlexibleTaskMatch{
		Number:  number,
		Status:  rune(checkChar),
		Content: content,
	}, true
}

// MatchNumberedSection parses a numbered section header from tasks.md format.
// Format: "## 1. Section Name"
// Returns the section name (without number prefix) and true if matched,
// or empty string and false otherwise.
//
// Example:
//
//	name, ok := MatchNumberedSection("## 1. Core Accept Command")
//	// name = "Core Accept Command", ok = true
func MatchNumberedSection(
	line string,
) (name string, ok bool) {
	// Must start with "## "
	if !strings.HasPrefix(line, "## ") {
		return "", false
	}

	// Get the part after "## "
	rest := line[3:]

	// Must start with digits
	if rest == "" || rest[0] < '0' ||
		rest[0] > '9' {
		return "", false
	}

	// Find the end of the number and dot pattern (e.g., "1." or "12.")
	i := 0
	for i < len(rest) && rest[i] >= '0' && rest[i] <= '9' {
		i++
	}

	// Must be followed by ". "
	if i >= len(rest) || rest[i] != '.' {
		return "", false
	}
	i++

	if i >= len(rest) || rest[i] != ' ' {
		return "", false
	}
	i++

	// Extract the section name
	name = rest[i:]
	if name == "" {
		return "", false
	}

	return name, true
}

// MatchAnySection parses any H2 section header from tasks.md format.
// Accepts both numbered and unnumbered formats:
//   - "## 1. Setup" -> name="Setup", number="1", ok=true
//   - "## Implementation" -> name="Implementation", number="", ok=true
//
// Returns the section name, optional number, and true if matched.
func MatchAnySection(
	line string,
) (name, number string, ok bool) {
	// Must start with "## "
	if !strings.HasPrefix(line, "## ") {
		return "", "", false
	}

	// Get the part after "## "
	rest := line[3:]
	if rest == "" {
		return "", "", false
	}

	// Try to match numbered format first: "N. Name"
	if rest[0] >= '0' && rest[0] <= '9' {
		// Find the end of the number
		i := 0
		for i < len(rest) && rest[i] >= '0' && rest[i] <= '9' {
			i++
		}

		// Check for ". " after number
		if i < len(rest) && rest[i] == '.' {
			i++
			if i < len(rest) && rest[i] == ' ' {
				i++
				// Extract number and name
				number = rest[:i-2] // digits only
				name = strings.TrimSpace(rest[i:])
				if name != "" {
					return name, number, true
				}
				// Numbered format matched but no content after "N. "
				// This is invalid, return false
				return "", "", false
			}
		}
	}

	// Not numbered format, treat as plain section name
	name = strings.TrimSpace(rest)
	if name == "" {
		return "", "", false
	}

	return name, "", true
}

// ExtractHeaderLevel returns the header level (1-6) for a line that starts
// with hash characters followed by a space. Returns 0 if not a valid header.
//
// Example:
//
//	ExtractHeaderLevel("# Title")    // returns 1
//	ExtractHeaderLevel("### Section") // returns 3
//	ExtractHeaderLevel("Not a header") // returns 0
func ExtractHeaderLevel(line string) int {
	if line == "" || line[0] != '#' {
		return 0
	}

	level := 0
	for i := 0; i < len(line) && i < 6; i++ {
		if line[i] == '#' {
			level++
		} else {
			break
		}
	}

	// Must be followed by space
	if level == 0 || level > 6 ||
		level >= len(line) ||
		line[level] != ' ' {
		return 0
	}

	return level
}

// ExtractHeaderText returns the header text without the "#" prefix and
// leading/trailing whitespace. Returns empty string if not a valid header.
//
// Example:
//
//	ExtractHeaderText("## My Section")  // returns "My Section"
//	ExtractHeaderText("# Title  ")       // returns "Title"
func ExtractHeaderText(line string) string {
	level := ExtractHeaderLevel(line)
	if level == 0 {
		return ""
	}

	// Skip the hashes and the space
	text := line[level+1:]

	// Trim whitespace
	return strings.TrimSpace(text)
}

// IsBlankLine checks if a line is empty or contains only whitespace.
func IsBlankLine(line string) bool {
	return strings.TrimSpace(line) == ""
}

// IsListItem checks if a line starts with a list marker.
// Supports unordered markers (-, *, +) and ordered markers (1., 2., etc.)
// after optional leading whitespace.
func IsListItem(line string) bool {
	trimmed := strings.TrimLeft(line, " \t")
	if trimmed == "" {
		return false
	}

	// Check unordered list markers
	if len(trimmed) >= 2 {
		first := trimmed[0]
		if (first == '-' || first == '*' || first == '+') &&
			trimmed[1] == ' ' {
			return true
		}
	}

	// Check ordered list markers (digits followed by . and space)
	for i := range len(trimmed) {
		c := trimmed[i]
		if c >= '0' && c <= '9' {
			continue
		}
		if c == '.' && i > 0 &&
			i+1 < len(trimmed) &&
			trimmed[i+1] == ' ' {
			return true
		}

		break
	}

	return false
}

// IsCodeFence checks if a line starts a code fence (``` or ~~~).
// Returns true and the fence character ('`' or '~') if it's a fence,
// or false and 0 otherwise.
func IsCodeFence(line string) (isFence bool, delimiter rune) {
	trimmed := strings.TrimLeft(line, " \t")

	if strings.HasPrefix(trimmed, "```") {
		return true, '`'
	}
	if strings.HasPrefix(trimmed, "~~~") {
		return true, '~'
	}

	return false, 0
}

// IsBlockquote checks if a line starts with a blockquote marker (>).
func IsBlockquote(line string) bool {
	trimmed := strings.TrimLeft(line, " \t")

	return strings.HasPrefix(trimmed, ">")
}

// MatchH2SectionHeader checks if a line is an H2 section header
// and extracts the name. Returns the section name and true if matched,
// or empty string and false otherwise.
// This is a compatibility wrapper matching the old regex API.
//
// Example:
//
//	name, ok := MatchH2SectionHeader("## Purpose")
//	// name = "Purpose", ok = true
func MatchH2SectionHeader(
	line string,
) (name string, ok bool) {
	if !strings.HasPrefix(line, "## ") {
		return "", false
	}

	name = strings.TrimPrefix(line, "## ")
	if name == "" {
		return "", false
	}

	return name, true
}

// MatchH2DeltaSection checks if a line is a delta section header
// (## ADDED|MODIFIED|REMOVED|RENAMED Requirements).
// Returns the delta type and true if matched, or empty string and false.
//
// Example:
//
//	deltaType, ok := MatchH2DeltaSection("## ADDED Requirements")
//	// deltaType = "ADDED", ok = true
func MatchH2DeltaSection(
	line string,
) (deltaType string, ok bool) {
	if !strings.HasPrefix(line, "## ") {
		return "", false
	}

	rest := strings.TrimPrefix(line, "## ")
	rest = strings.TrimSpace(rest)

	// Check for each delta type
	deltaTypes := []string{
		"ADDED",
		"MODIFIED",
		"REMOVED",
		"RENAMED",
	}
	for _, dt := range deltaTypes {
		suffix := dt + " Requirements"
		if rest == suffix {
			return dt, true
		}
	}

	return "", false
}

// IsTaskChecked returns true if the checkbox state indicates completion.
// Accepts 'x' or 'X' as checked states.
func IsTaskChecked(state rune) bool {
	return state == 'x' || state == 'X'
}

// ExtractListMarker extracts the list marker from a line if present.
// Returns the marker string and true if found, or empty and false otherwise.
// For unordered lists, returns "-", "*", or "+".
// For ordered lists, returns the number with dot (e.g., "1.", "23.").
func ExtractListMarker(
	line string,
) (marker string, ok bool) {
	trimmed := strings.TrimLeft(line, " \t")
	if trimmed == "" {
		return "", false
	}

	// Check unordered list markers
	if len(trimmed) >= 2 {
		first := trimmed[0]
		if (first == '-' || first == '*' || first == '+') &&
			trimmed[1] == ' ' {
			return string(first), true
		}
	}

	// Check ordered list markers
	var numEnd int
	for numEnd = 0; numEnd < len(trimmed); numEnd++ { //nolint:intrange // numEnd needed after loop
		c := trimmed[numEnd]
		if c < '0' || c > '9' {
			break
		}
	}

	if numEnd > 0 && numEnd < len(trimmed) &&
		trimmed[numEnd] == '.' {
		if numEnd+1 < len(trimmed) &&
			trimmed[numEnd+1] == ' ' {
			return trimmed[:numEnd+1], true
		}
	}

	return "", false
}

// CountLeadingSpaces returns the number of leading space characters in a line.
// Tabs are counted as 1 character (not expanded).
func CountLeadingSpaces(line string) int {
	count := 0
	for _, c := range line {
		if c == ' ' || c == '\t' {
			count++
		} else {
			break
		}
	}

	return count
}

// TrimLeadingHashes removes leading hash characters and the following space
// from a header line, returning the header content.
func TrimLeadingHashes(line string) string {
	i := 0
	for i < len(line) && line[i] == '#' {
		i++
	}
	if i < len(line) && line[i] == ' ' {
		i++
	}

	return line[i:]
}

// ContainsKeyword checks if a line contains one of the Spectr keywords
// (WHEN, THEN, AND, GIVEN) typically used in scenario descriptions.
// Returns the keyword found and true, or empty and false.
func ContainsKeyword(
	line string,
) (keyword string, ok bool) {
	keywords := []string{
		"WHEN",
		"THEN",
		"AND",
		"GIVEN",
	}

	// Look for **KEYWORD** pattern (bold keywords)
	for _, kw := range keywords {
		boldKw := "**" + kw + "**"
		if strings.Contains(line, boldKw) {
			return kw, true
		}
	}

	return "", false
}

// IsHorizontalRule checks if a line is a horizontal rule (---, ***, ___).
func IsHorizontalRule(line string) bool {
	trimmed := strings.TrimSpace(line)
	if len(trimmed) < 3 {
		return false
	}

	// Check for ---, ***, or ___
	first := trimmed[0]
	if first != '-' && first != '*' &&
		first != '_' {
		return false
	}

	// All characters must be the same (with optional spaces)
	for _, c := range trimmed {
		if c != rune(first) &&
			!unicode.IsSpace(c) {
			return false
		}
	}

	// Count the rule characters
	count := 0
	for _, c := range trimmed {
		if c == rune(first) {
			count++
		}
	}

	return count >= 3
}

// FindDeltaSectionContent extracts the content from a delta section
// (ADDED, MODIFIED, REMOVED, RENAMED) without the header line.
// This is a compatibility function that matches the behavior of the
// old regex.FindDeltaSectionContent function.
//
// Returns the content between the section header and the next H2 header,
// or empty string if the section is not found.
//
// Example:
//
//	content := `# Delta
//	## ADDED Requirements
//	### Requirement: New Feature
//	Content here.
//	## MODIFIED Requirements
//	...`
//	FindDeltaSectionContent([]byte(content), DeltaAdded)
//	// returns "\n### Requirement: New Feature\nContent here.\n"
func FindDeltaSectionContent(
	content []byte,
	deltaType DeltaType,
) string {
	// Build the section header to search for
	sectionHeader := "## " + string(
		deltaType,
	) + " Requirements"

	// Find the start of the section
	headerStart := bytes.Index(
		content,
		[]byte(sectionHeader),
	)
	if headerStart == -1 {
		return ""
	}

	// Find the end of the header line
	headerEnd := headerStart + len(sectionHeader)
	for headerEnd < len(content) && content[headerEnd] != '\n' {
		headerEnd++
	}
	if headerEnd < len(content) {
		headerEnd++ // Include the newline
	}

	// Find the next H2 header (## followed by space)
	rest := content[headerEnd:]
	nextH2 := -1
	lines := bytes.Split(rest, []byte("\n"))
	offset := 0
	for _, line := range lines {
		trimmed := bytes.TrimLeft(line, " \t")
		if bytes.HasPrefix(
			trimmed,
			[]byte("## "),
		) {
			nextH2 = offset

			break
		}
		offset += len(line) + 1 // +1 for newline
	}

	if nextH2 != -1 {
		return string(rest[:nextH2])
	}

	return string(rest)
}

// MatchRenamedFrom checks if a line matches the backtick-wrapped FROM format
// for renamed requirements. Returns the requirement name and true if matched,
// or empty string and false otherwise.
//
// Expected format: "- FROM: `### Requirement: Old Name`"
//
// Example:
//
//	name, ok := MatchRenamedFrom("- FROM: `### Requirement: Old Name`")
//	// name = "Old Name", ok = true
func MatchRenamedFrom(
	line string,
) (name string, ok bool) {
	trimmed := strings.TrimSpace(line)

	// Must start with "- FROM:" (case-insensitive check)
	if !strings.HasPrefix(trimmed, "-") {
		return "", false
	}

	// Skip the bullet and whitespace
	rest := strings.TrimSpace(trimmed[1:])

	// Check for "FROM:" prefix (case-insensitive)
	upper := strings.ToUpper(rest)
	if !strings.HasPrefix(upper, "FROM:") {
		return "", false
	}

	// Get the content after "FROM:"
	rest = strings.TrimSpace(rest[5:])

	// Check for backtick-wrapped content
	if !strings.HasPrefix(rest, "`") ||
		!strings.HasSuffix(rest, "`") {
		return "", false
	}

	// Extract content within backticks
	content := rest[1 : len(rest)-1]

	// Must match "### Requirement: Name"
	const prefix = "### Requirement:"
	upperContent := strings.ToUpper(content)
	if !strings.HasPrefix(
		upperContent,
		"### REQUIREMENT:",
	) {
		return "", false
	}

	// Extract the name after the prefix (preserving original case)
	nameStart := len(prefix)
	if nameStart >= len(content) {
		return "", false
	}

	name = strings.TrimSpace(content[nameStart:])
	if name == "" {
		return "", false
	}

	return name, true
}

// MatchRenamedTo checks if a line matches the backtick-wrapped TO format
// for renamed requirements. Returns the requirement name and true if matched,
// or empty string and false otherwise.
//
// Expected format: "- TO: `### Requirement: New Name`"
//
// Example:
//
//	name, ok := MatchRenamedTo("- TO: `### Requirement: New Name`")
//	// name = "New Name", ok = true
func MatchRenamedTo(
	line string,
) (name string, ok bool) {
	trimmed := strings.TrimSpace(line)

	// Must start with "- TO:" (case-insensitive check)
	if !strings.HasPrefix(trimmed, "-") {
		return "", false
	}

	// Skip the bullet and whitespace
	rest := strings.TrimSpace(trimmed[1:])

	// Check for "TO:" prefix (case-insensitive)
	upper := strings.ToUpper(rest)
	if !strings.HasPrefix(upper, "TO:") {
		return "", false
	}

	// Get the content after "TO:"
	rest = strings.TrimSpace(rest[3:])

	// Check for backtick-wrapped content
	if !strings.HasPrefix(rest, "`") ||
		!strings.HasSuffix(rest, "`") {
		return "", false
	}

	// Extract content within backticks
	content := rest[1 : len(rest)-1]

	// Must match "### Requirement: Name"
	const prefix = "### Requirement:"
	upperContent := strings.ToUpper(content)
	if !strings.HasPrefix(
		upperContent,
		"### REQUIREMENT:",
	) {
		return "", false
	}

	// Extract the name after the prefix (preserving original case)
	nameStart := len(prefix)
	if nameStart >= len(content) {
		return "", false
	}

	name = strings.TrimSpace(content[nameStart:])
	if name == "" {
		return "", false
	}

	return name, true
}

// MatchRenamedFromAlt checks if a line matches the non-backtick FROM format
// for renamed requirements. Returns the requirement name and true if matched,
// or empty string and false otherwise.
//
// Expected format: "- FROM: ### Requirement: Old Name"
//
// Example:
//
//	name, ok := MatchRenamedFromAlt("- FROM: ### Requirement: Old Name")
//	// name = "Old Name", ok = true
func MatchRenamedFromAlt(
	line string,
) (name string, ok bool) {
	trimmed := strings.TrimSpace(line)

	// Must start with "-"
	if !strings.HasPrefix(trimmed, "-") {
		return "", false
	}

	// Skip the bullet and whitespace
	rest := strings.TrimSpace(trimmed[1:])

	// Check for "FROM:" prefix (case-insensitive)
	upper := strings.ToUpper(rest)
	if !strings.HasPrefix(upper, "FROM:") {
		return "", false
	}

	// Get the content after "FROM:"
	rest = strings.TrimSpace(rest[5:])

	// Must match "### Requirement: Name" (without backticks)
	// Skip any leading ### and whitespace
	if !strings.HasPrefix(rest, "###") {
		return "", false
	}

	// Skip "###" and any whitespace
	rest = strings.TrimLeft(rest[3:], " \t")

	// Check for "Requirement:" prefix (case-insensitive)
	upper = strings.ToUpper(rest)
	if !strings.HasPrefix(upper, "REQUIREMENT:") {
		return "", false
	}

	// Get the name after "Requirement:"
	name = strings.TrimSpace(rest[12:])
	if name == "" {
		return "", false
	}

	return name, true
}

// MatchRenamedToAlt checks if a line matches the non-backtick TO format
// for renamed requirements. Returns the requirement name and true if matched,
// or empty string and false otherwise.
//
// Expected format: "- TO: ### Requirement: New Name"
//
// Example:
//
//	name, ok := MatchRenamedToAlt("- TO: ### Requirement: New Name")
//	// name = "New Name", ok = true
func MatchRenamedToAlt(
	line string,
) (name string, ok bool) {
	trimmed := strings.TrimSpace(line)

	// Must start with "-"
	if !strings.HasPrefix(trimmed, "-") {
		return "", false
	}

	// Skip the bullet and whitespace
	rest := strings.TrimSpace(trimmed[1:])

	// Check for "TO:" prefix (case-insensitive)
	upper := strings.ToUpper(rest)
	if !strings.HasPrefix(upper, "TO:") {
		return "", false
	}

	// Get the content after "TO:"
	rest = strings.TrimSpace(rest[3:])

	// Must match "### Requirement: Name" (without backticks)
	// Skip any leading ### and whitespace
	if !strings.HasPrefix(rest, "###") {
		return "", false
	}

	// Skip "###" and any whitespace
	rest = strings.TrimLeft(rest[3:], " \t")

	// Check for "Requirement:" prefix (case-insensitive)
	upper = strings.ToUpper(rest)
	if !strings.HasPrefix(upper, "REQUIREMENT:") {
		return "", false
	}

	// Get the name after "Requirement:"
	name = strings.TrimSpace(rest[12:])
	if name == "" {
		return "", false
	}

	return name, true
}

// MatchAnyRenamedFrom tries both backtick and non-backtick FROM formats.
// Returns the requirement name and true if either format matches.
// Tries backtick format first, then falls back to non-backtick.
func MatchAnyRenamedFrom(
	line string,
) (name string, ok bool) {
	if name, ok = MatchRenamedFrom(line); ok {
		return name, true
	}

	return MatchRenamedFromAlt(line)
}

// MatchAnyRenamedTo tries both backtick and non-backtick TO formats.
// Returns the requirement name and true if either format matches.
// Tries backtick format first, then falls back to non-backtick.
func MatchAnyRenamedTo(
	line string,
) (name string, ok bool) {
	if name, ok = MatchRenamedTo(line); ok {
		return name, true
	}

	return MatchRenamedToAlt(line)
}

// FindH2RequirementsSection finds the "## Requirements" section header in content.
// Returns the start and end indices of the header match, or nil if not found.
// The end index points to the end of the header line (after "## Requirements").
//
// Example:
//
//	indices := FindH2RequirementsSection("# Title\n\n## Requirements\n\nContent")
//	// indices = [9, 24] (pointing to "## Requirements")
func FindH2RequirementsSection(
	content string,
) []int {
	lines := strings.Split(content, "\n")
	offset := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Match "## Requirements" with flexible whitespace
		if strings.HasPrefix(trimmed, "##") {
			rest := strings.TrimLeft(
				trimmed[2:],
				" \t",
			)
			if rest == "Requirements" {
				// Found it - return indices [start, end]
				return []int{
					offset,
					offset + len(line),
				}
			}
		}
		offset += len(line) + 1 // +1 for newline
	}

	return nil
}

// FindNextH2Section finds the next "## " header in content starting from startPos.
// Returns the start and end indices of the header match relative to the beginning
// of the content, or nil if no H2 header is found after startPos.
//
// Example:
//
//	content := "## First\n\nSome text\n\n## Second\n"
//	indices := FindNextH2Section(content, 9) // Search after "## First\n"
//	// indices = [22, 31] (pointing to "## Second")
func FindNextH2Section(
	content string,
	startPos int,
) []int {
	if startPos >= len(content) {
		return nil
	}

	remaining := content[startPos:]
	lines := strings.Split(remaining, "\n")
	offset := 0

	for _, line := range lines {
		trimmed := strings.TrimLeft(line, " \t")
		// Match any "## " header
		if strings.HasPrefix(trimmed, "## ") {
			// Found it - return indices relative to original content
			return []int{
				startPos + offset,
				startPos + offset + len(line),
			}
		}
		offset += len(line) + 1 // +1 for newline
	}

	return nil
}

// FindAllH3Requirements finds all requirement headers in content and returns
// their names. This function scans for "### Requirement: Name" patterns and
// extracts the names.
//
// Example:
//
//	content := "### Requirement: Auth\n\nContent\n\n### Requirement: Logging\n"
//	names := FindAllH3Requirements(content)
//	// names = ["Auth", "Logging"]
func FindAllH3Requirements(
	content string,
) []string {
	var names []string
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		if name, ok := MatchRequirementHeader(line); ok {
			names = append(names, name)
		}
	}

	return names
}

// NormalizeMultipleNewlines collapses sequences of 3 or more consecutive
// newlines into exactly 2 newlines. This replaces the inline regex pattern
// `\n{3,}` with a string-based implementation.
//
// Example:
//
//	result := NormalizeMultipleNewlines("text\n\n\n\nmore")
//	// result = "text\n\nmore"
func NormalizeMultipleNewlines(s string) string {
	var result strings.Builder
	result.Grow(len(s))

	newlineCount := 0
	for i := range len(s) {
		if s[i] == '\n' {
			newlineCount++
			// Only write up to 2 consecutive newlines
			if newlineCount <= 2 {
				result.WriteByte('\n')
			}
		} else {
			newlineCount = 0
			result.WriteByte(s[i])
		}
	}

	return result.String()
}
