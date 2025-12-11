// Package markdown provides AST-based markdown parsing using blackfriday v2.
// This package replaces regex-based markdown parsing throughout the spectr
// codebase, providing more robust and maintainable parsing of spec documents.
package markdown

// Header represents a markdown header (H1-H6)
type Header struct {
	Level int    // Header level (1-6)
	Text  string // Header text content
	Line  int    // Line number in the source document (1-indexed)
}

// Section represents a markdown section with its content
type Section struct {
	Name      string // Section header text
	Content   string // Full content under the section header
	Level     int    // Header level (1-6)
	StartLine int    // Line number where section starts (1-indexed)
}

// Task represents a markdown task checkbox item
type Task struct {
	Checked bool   // true if [x] or [X], false if [ ]
	Text    string // Task text content
	Line    int    // Line number in the source document (1-indexed)
}

// RequirementBlock represents a parsed requirement with its content
// This mirrors the parsers.RequirementBlock type for compatibility
type RequirementBlock struct {
	Name       string   // Requirement name (from "### Requirement: Name")
	HeaderLine string   // The full header line
	Raw        string   // Full block content including header
	Scenarios  []string // Extracted scenario names
}

// DeltaType represents the type of delta section
type DeltaType string

const (
	DeltaAdded    DeltaType = "ADDED"
	DeltaModified DeltaType = "MODIFIED"
	DeltaRemoved  DeltaType = "REMOVED"
	DeltaRenamed  DeltaType = "RENAMED"
)

// ValidDeltaTypes returns all valid delta type strings
func ValidDeltaTypes() []string {
	return []string{
		string(DeltaAdded),
		string(DeltaModified),
		string(DeltaRemoved),
		string(DeltaRenamed),
	}
}
