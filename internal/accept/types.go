package accept

// Task represents a single task item from tasks.md with its completion
// status and nested subtasks.
type Task struct {
	// ID is the task identifier (e.g., "1.1", "2.3", "1.1.1.1")
	ID string `json:"id"`
	// Description is the task text, with indented detail lines appended
	Description string `json:"description"`
	// Completed indicates [x] vs [ ]
	Completed bool `json:"completed"`
	// Subtasks supports recursive unlimited nesting depth
	Subtasks []Task `json:"subtasks"`
}

// Section represents a section from tasks.md containing grouped tasks.
type Section struct {
	Number int    `json:"number"` // Section number from header
	Name   string `json:"name"`   // Section title
	Tasks  []Task `json:"tasks"`  // Tasks in this section
}

// Summary provides aggregate counts of tasks and their completion status.
type Summary struct {
	// Total tasks including all nested subtasks
	Total int `json:"total"`
	// Completed tasks including nested subtasks
	Completed int `json:"completed"`
}

// TasksJSON is the root structure for the tasks.json output file.
type TasksJSON struct {
	Version    string    `json:"version"`    // Schema version, value "1.0"
	ChangeID   string    `json:"changeId"`   // The change identifier
	AcceptedAt string    `json:"acceptedAt"` // ISO 8601 timestamp
	Sections   []Section `json:"sections"`   // Parsed sections from tasks.md
	Summary    Summary   `json:"summary"`    // Aggregate task counts
}
