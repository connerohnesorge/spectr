package markdown

// Document represents a parsed markdown document with all extracted elements.
// The document is parsed once via AST walk, and all data is available for
// O(1) lookups via the indexed maps.
type Document struct {
	// Source data
	Content []byte   // Original raw content
	Lines   []string // Content split by newline

	// Structure extracted from AST (single walk, in document order)
	Headers []Header // All headers H1-H6 in document order

	// Indexed data for O(1) lookups
	Sections     map[string]*Section     // H2 sections by name
	Requirements map[string]*Requirement // Requirements by name
	Scenarios    map[string]*Scenario    // Scenarios by name

	// Headers by level (for ordered iteration)
	H2Headers []Header
	H3Headers []Header
	H4Headers []Header

	// Tasks contains all checkbox items with hierarchical structure.
	Tasks []Task
}

// Header represents a markdown heading with source location and byte offsets.
type Header struct {
	// Level is the heading level (1-6, e.g., 1 for H1, 2 for H2).
	Level int

	// Text is the header text content without the leading # marks.
	Text string

	// Line is the source line number (1-indexed).
	Line int

	// Start is the byte offset where this header starts in content.
	Start int

	// End is the byte offset where this header ends in content.
	End int
}

// Section represents an H2 section with its content.
type Section struct {
	// Name is the section header text (e.g., "Requirements").
	Name string

	// Header is the section's H2 header.
	Header Header

	// Content is the raw markdown text between this H2 and the next.
	// This preserves original formatting including blank lines.
	Content string

	// StartLine is the line number where section content starts (after header).
	StartLine int

	// EndLine is the line number where section content ends (exclusive).
	EndLine int

	// IsDelta indicates this is a delta section.
	// Delta types: ADDED, MODIFIED, REMOVED, RENAMED.
	IsDelta bool

	// DeltaType is the delta type if IsDelta is true.
	DeltaType string
}

// Requirement represents an H3 requirement header with its content.
type Requirement struct {
	// Name is the requirement name (text after "Requirement: ").
	Name string

	// Section is the parent section name.
	Section string

	// Header is the requirement's H3 header.
	Header Header

	// Content is the raw markdown text under this requirement.
	Content string

	// Scenarios contains pointers to child scenarios.
	Scenarios []*Scenario
}

// Scenario represents an H4 scenario header with its content.
type Scenario struct {
	// Name is the scenario name (text after "Scenario: ").
	Name string

	// Requirement is the parent requirement name.
	Requirement string

	// Header is the scenario's H4 header.
	Header Header

	// Content is the raw markdown text under this scenario.
	Content string
}

// Task represents a checkbox item with hierarchy support.
type Task struct {
	// Text is the full original line text.
	Text string

	// Checked is true if [x] or [X], false if [ ].
	Checked bool

	// Line is the source line number (1-indexed).
	Line int

	// Indent is the indentation level (number of leading spaces/tabs).
	Indent int

	// Children contains nested task items.
	Children []Task
}

// NumberedTask represents a task with ID in the format "1.1 Description".
// Used for tasks.md parsing.
type NumberedTask struct {
	// ID is the task identifier (e.g., "1.1", "2.3").
	ID string

	// Description is the task description text.
	Description string

	// Checked is true if the task is marked complete.
	Checked bool

	// LineNum is the source line number (1-indexed).
	LineNum int
}

// NumberedSection represents a section header with number prefix.
// Used for tasks.md parsing (e.g., "## 1. Section Name").
type NumberedSection struct {
	// Number is the section number (e.g., "1", "2").
	Number string

	// Name is the section name without the number prefix.
	Name string

	// LineNum is the source line number (1-indexed).
	LineNum int
}

// RenamedPair represents a FROM/TO rename pair.
type RenamedPair struct {
	// From is the original requirement name.
	From string

	// To is the new requirement name.
	To string

	// FromLine is the source line number of the FROM entry.
	FromLine int

	// ToLine is the source line number of the TO entry.
	ToLine int
}
