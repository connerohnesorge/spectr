package markdown

// Document represents a parsed markdown document with all extracted elements.
type Document struct {
	Headers  []Header           // All headers H1-H4
	Sections map[string]Section // Sections keyed by header text
	Tasks    []Task             // All task checkboxes (hierarchical)
}

// Header represents a markdown heading with source location.
type Header struct {
	Level int    // 1-6
	Text  string // Header text content
	Line  int    // Source line number (1-indexed)
}

// Section represents content between headers.
type Section struct {
	Header  Header // The section's header
	Content string // Raw markdown text between this header and next
}

// Task represents a checkbox item with hierarchy support.
type Task struct {
	Line     string // Full original line text
	Checked  bool   // true if [x] or [X]
	LineNum  int    // Source line number (1-indexed)
	Children []Task // Nested task items
}
