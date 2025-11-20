// Package view provides dashboard functionality for displaying
// a comprehensive project overview including specs, changes, and tasks.
package view

import (
	"path/filepath"
	"sort"

	"github.com/connerohnesorge/spectr/internal/discovery"
	"github.com/connerohnesorge/spectr/internal/parsers"
)

// CollectData gathers all dashboard information from the project,
// including active changes, completed changes, specifications,
// and summary metrics.
//
// The function performs the following steps:
//  1. Discovers all changes in spectr/changes/ directory
//  2. Parses each change's proposal.md for title
//  3. Parses each change's tasks.md for task completion status
//  4. Categorizes changes as active (incomplete) or completed
//  5. Discovers all specs in spectr/specs/ directory
//  6. Parses each spec's spec.md for title and requirement count
//  7. Sorts results per design specification (active changes by
//     completion ascending, specs by requirement count descending)
//
// Returns DashboardData structure or error if discovery fails.
//
//nolint:revive // cognitive-complexity justified for data collection
func CollectData(projectPath string) (*DashboardData, error) {
	data := &DashboardData{
		Summary:          SummaryMetrics{},
		ActiveChanges:    []ChangeProgress{},
		CompletedChanges: []CompletedChange{},
		Specs:            []SpecInfo{},
	}

	// Discover all changes
	changeIDs, err := discovery.GetActiveChanges(projectPath)
	if err != nil {
		return nil, err
	}

	// Process each change
	for _, changeID := range changeIDs {
		changeDir := filepath.Join(projectPath, "spectr", "changes", changeID)

		// Parse title from proposal.md
		proposalPath := filepath.Join(changeDir, "proposal.md")
		title, err := parsers.ExtractTitle(proposalPath)
		if err != nil {
			// Skip changes that can't be parsed, use ID as fallback title
			title = changeID
		}

		// Parse task counts from tasks.md
		tasksPath := filepath.Join(changeDir, "tasks.md")
		taskStatus, err := parsers.CountTasks(tasksPath)
		if err != nil {
			// If tasks.md can't be read, assume zero tasks
			taskStatus = parsers.TaskStatus{Total: 0, Completed: 0}
		}

		// Calculate completion percentage
		percentage := calculatePercentage(taskStatus.Completed, taskStatus.Total)

		// Categorize as active or completed
		// A change is completed if: (completed == total AND total > 0) OR total == 0
		isCompleted := (taskStatus.Completed == taskStatus.Total && taskStatus.Total > 0) ||
			taskStatus.Total == 0

		if isCompleted {
			// Add to completed changes
			data.CompletedChanges = append(data.CompletedChanges, CompletedChange{
				ID:    changeID,
				Title: title,
			})
			data.Summary.CompletedChanges++
		} else {
			// Add to active changes
			data.ActiveChanges = append(data.ActiveChanges, ChangeProgress{
				ID:    changeID,
				Title: title,
				Progress: ProgressMetrics{
					Total:      taskStatus.Total,
					Completed:  taskStatus.Completed,
					Percentage: percentage,
				},
			})
			data.Summary.ActiveChanges++
			data.Summary.TotalTasks += taskStatus.Total
			data.Summary.CompletedTasks += taskStatus.Completed
		}
	}

	// Discover all specs
	specIDs, err := discovery.GetSpecs(projectPath)
	if err != nil {
		return nil, err
	}

	// Process each spec
	for _, specID := range specIDs {
		specPath := filepath.Join(projectPath, "spectr", "specs", specID, "spec.md")

		// Parse title from spec.md
		title, err := parsers.ExtractTitle(specPath)
		if err != nil {
			// Skip specs that can't be parsed, use ID as fallback title
			title = specID
		}

		// Count requirements
		reqCount, err := parsers.CountRequirements(specPath)
		if err != nil {
			// If requirement count fails, assume zero
			reqCount = 0
		}

		// Add to specs list
		data.Specs = append(data.Specs, SpecInfo{
			ID:               specID,
			Title:            title,
			RequirementCount: reqCount,
		})

		// Update summary metrics
		data.Summary.TotalSpecs++
		data.Summary.TotalRequirements += reqCount
	}

	// Sort active changes by completion percentage (ascending), then by ID (alphabetical)
	sort.Slice(data.ActiveChanges, func(i, j int) bool {
		// Sort by percentage first (ascending - lower completion first)
		if data.ActiveChanges[i].Progress.Percentage != data.ActiveChanges[j].Progress.Percentage {
			return data.ActiveChanges[i].Progress.Percentage < data.ActiveChanges[j].Progress.Percentage
		}
		// Tie-breaker: alphabetical by ID
		return data.ActiveChanges[i].ID < data.ActiveChanges[j].ID
	})

	// Sort completed changes alphabetically by ID
	sort.Slice(data.CompletedChanges, func(i, j int) bool {
		return data.CompletedChanges[i].ID < data.CompletedChanges[j].ID
	})

	// Sort specs by requirement count (descending), then by ID (alphabetical)
	sort.Slice(data.Specs, func(i, j int) bool {
		// Sort by requirement count first (descending - more requirements first)
		if data.Specs[i].RequirementCount != data.Specs[j].RequirementCount {
			return data.Specs[i].RequirementCount > data.Specs[j].RequirementCount
		}
		// Tie-breaker: alphabetical by ID
		return data.Specs[i].ID < data.Specs[j].ID
	})

	return data, nil
}

// calculatePercentage calculates completion percentage (0-100).
// It handles the zero division edge case by returning 0 when total is 0.
// The result is rounded to the nearest integer using standard rounding rules.
//
// Example:
//   - calculatePercentage(3, 8) returns 37
//   - calculatePercentage(0, 10) returns 0
//   - calculatePercentage(5, 0) returns 0
func calculatePercentage(completed, total int) int {
	if total == 0 {
		return 0
	}
	// Round to nearest integer
	return int(float64(completed) / float64(total) * 100.0)
}
