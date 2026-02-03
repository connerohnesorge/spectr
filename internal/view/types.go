// Package view provides dashboard functionality for displaying
// a comprehensive project overview including specs, changes, and tasks.
package view

// DashboardData represents the complete dashboard data structure
// containing summary metrics, active changes, completed changes,
// and specifications.
type DashboardData struct {
	Summary          SummaryMetrics    `json:"summary"`
	ActiveChanges    []ChangeProgress  `json:"activeChanges"`
	CompletedChanges []CompletedChange `json:"completedChanges"`
	Specs            []SpecInfo        `json:"specs"`
	HasMultipleRoots bool              `json:"hasMultipleRoots,omitempty"`
}

// SummaryMetrics represents aggregate metrics across the project
type SummaryMetrics struct {
	// Total number of specifications
	TotalSpecs int `json:"totalSpecs"`
	// Total requirements across all specs
	TotalRequirements int `json:"totalRequirements"`
	// Number of active changes (incomplete)
	ActiveChanges int `json:"activeChanges"`
	// Number of completed changes
	CompletedChanges int `json:"completedChanges"`
	// Total tasks across all active changes
	TotalTasks int `json:"totalTasks"`
	// Completed tasks across all active changes
	CompletedTasks int `json:"completedTasks"`
}

// ChangeProgress represents an active change with task completion progress
type ChangeProgress struct {
	ID       string          `json:"id"`                 // Change ID (directory name)
	Title    string          `json:"title"`              // Change title from proposal.md
	Progress ProgressMetrics `json:"progress"`           // Task completion metrics
	RootPath string          `json:"rootPath,omitempty"` // Relative path to root (multi-root only)
}

// ProgressMetrics represents task completion statistics for a change
type ProgressMetrics struct {
	Total      int `json:"total"`      // Total number of tasks
	Completed  int `json:"completed"`  // Number of completed tasks
	Percentage int `json:"percentage"` // Completion percentage (0-100)
}

// CompletedChange represents a change that has all tasks completed
type CompletedChange struct {
	ID       string `json:"id"`                 // Change ID (directory name)
	Title    string `json:"title"`              // Change title from proposal.md
	RootPath string `json:"rootPath,omitempty"` // Relative path to root (multi-root only)
}

// SpecInfo represents a specification with metadata
type SpecInfo struct {
	// Spec ID (directory name)
	ID string `json:"id"`
	// Spec title from spec.md
	Title string `json:"title"`
	// Number of requirements in spec
	RequirementCount int `json:"requirementCount"`
	// RootPath is the relative path to the root (multi-root only)
	RootPath string `json:"rootPath,omitempty"`
}
