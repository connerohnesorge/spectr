// Package ralph provides task orchestration for Spectr change proposals.
package ralph

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/connerohnesorge/spectr/internal/parsers"
)

// ErrCircularDependency is returned when a circular dependency is detected in the task graph.
var ErrCircularDependency = errors.New("circular dependency detected in task graph")

// Task status constants
const (
	StatusPending    = "pending"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
	StatusFailed     = "failed"
)

// graph.go defines the task graph data structures and dependency resolution.
// It provides functionality for:
// - Parsing tasks.jsonc files into a structured task graph
// - Topological sorting for dependency-aware execution order
// - Detection of parallel execution opportunities (independent tasks)
// - Task dependency tracking based on hierarchical task IDs
//
// Task dependencies are inferred from task ID prefixes:
// - Task 1.2 depends on task 1.1
// - Task 2.1 is independent of task 1.x (can run in parallel)
// - Task 1.2.1 depends on task 1.2
//

// Task represents a single task from a tasks.jsonc file.
// Tasks are the fundamental units of work that can be delegated to AI agents
// for execution. Each task has a hierarchical ID that determines its dependencies.
type Task struct {
	// ID is the hierarchical task identifier (e.g., "1.2", "3.4.1").
	// Task dependencies are inferred from ID prefixes:
	//   - "1.2" depends on "1.1"
	//   - "2.1" is independent of "1.x"
	//   - "1.2.1" depends on "1.2"
	ID string `json:"id"`

	// Section is the human-readable section name this task belongs to.
	// Example: "Core Infrastructure", "TUI Implementation"
	Section string `json:"section"`

	// Description is a concise explanation of what this task entails.
	// This is used to generate the task context for AI agents.
	Description string `json:"description"`

	// Status indicates the current state of the task.
	// Valid values: "pending", "in_progress", "completed"
	// Tasks progress through states: pending -> in_progress -> completed
	Status string `json:"status"`

	// Children is a JSON reference to a child task file, if this task
	// has subtasks defined in a separate file (e.g., "$ref:tasks-2.jsonc").
	// Empty string if the task has no subtasks.
	Children string `json:"children,omitempty"`
}

// TaskGraph represents a dependency graph of tasks parsed from tasks.jsonc files.
// It provides efficient lookups for task execution order and parallelization.
type TaskGraph struct {
	// Tasks maps task ID to the task details.
	// This is the primary storage for all tasks in the graph.
	// Example: "1.2" -> Task{ID: "1.2", Section: "...", ...}
	Tasks map[string]*Task

	// Children maps a parent task ID to its list of child task IDs.
	// This enables efficient traversal of the task hierarchy.
	// Example: "1" -> ["1.1", "1.2", "1.3"]
	Children map[string][]string

	// Roots contains task IDs that have no dependencies (no parent prefix).
	// These are the entry points for task execution.
	// Tasks with different root prefixes can run in parallel.
	// Example: ["1.1", "2.1", "3.1"] - three independent task chains
	Roots []string
}

// TasksFile represents the JSONC file structure for task definitions.
// This mirrors parsers.TasksFile but exists in the ralph package for clarity.
type TasksFile struct {
	Version int    `json:"version"`
	Tasks   []Task `json:"tasks"`
	// Parent is the parent task ID for child task files (version 2)
	Parent string `json:"parent,omitempty"`
	// Includes is a list of child task file patterns (version 2 root files)
	Includes []string `json:"includes,omitempty"`
}

// ParseTaskGraph parses tasks*.jsonc files from a change directory and builds a TaskGraph.
// It discovers all tasks*.jsonc files using glob pattern matching, parses the JSONC content,
// and builds parent-child relationships based on task ID prefixes.
//
// Task ID prefix rules:
//   - Task "1" is a root
//   - Task "1.1" is a child of "1"
//   - Task "1.1.2" is a child of "1.1"
//   - Tasks "1.1" and "1.2" are siblings (same parent "1")
//   - Tasks "1" and "2" are independent (different roots)
//
// Parameters:
//   - changeDir: Path to the change directory containing tasks*.jsonc files
//
// Returns:
//   - *TaskGraph: The constructed task graph with all tasks and relationships
//   - error: Any error encountered during file discovery or parsing
func ParseTaskGraph(changeDir string) (*TaskGraph, error) {
	// Find all tasks*.jsonc files in the change directory
	pattern := filepath.Join(changeDir, "tasks*.jsonc")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to glob tasks files: %w", err)
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("no tasks*.jsonc files found in %s", changeDir)
	}

	// Initialize the graph
	graph := &TaskGraph{
		Tasks:    make(map[string]*Task),
		Children: make(map[string][]string),
		Roots:    make([]string, 0),
	}

	// Parse each tasks*.jsonc file
	for _, filePath := range matches {
		if err := parseTasksFile(filePath, graph); err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", filePath, err)
		}
	}

	// Build parent-child relationships and identify roots
	buildTaskRelationships(graph)

	return graph, nil
}

// parseTasksFile reads and parses a single tasks*.jsonc file, adding tasks to the graph.
func parseTasksFile(filePath string, graph *TaskGraph) error {
	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Strip JSONC comments to get valid JSON
	jsonData := parsers.JSONCToJSON(data)

	// Unmarshal into TasksFile
	var tasksFile TasksFile
	if err := json.Unmarshal(jsonData, &tasksFile); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// Add all tasks to the graph
	for _, task := range tasksFile.Tasks {
		// Create a copy of the task to avoid pointer issues
		taskCopy := task
		graph.Tasks[task.ID] = &taskCopy
	}

	return nil
}

// buildTaskRelationships analyzes task IDs to build parent-child relationships and identify roots.
// It uses task ID prefixes to infer dependencies:
//   - "1.2" depends on "1.1" (sequential siblings)
//   - "1.2.1" depends on "1.2" (parent-child)
//   - "2.1" is independent of "1.x" (different root)
//
// A task is a root if its parent ID doesn't exist as a task in the graph.
// For example, "1.1" is a root if there's no task with ID "1".
func buildTaskRelationships(graph *TaskGraph) {
	// First pass: build parent-child relationships
	for taskID := range graph.Tasks {
		parentID := getParentID(taskID)
		if parentID != "" {
			// Add this task as a child of its parent
			graph.Children[parentID] = append(graph.Children[parentID], taskID)
		}
	}

	// Second pass: identify roots (tasks whose parent doesn't exist)
	for taskID := range graph.Tasks {
		parentID := getParentID(taskID)
		if parentID == "" {
			// No parent at all - this is definitely a root
			graph.Roots = append(graph.Roots, taskID)
		} else if _, exists := graph.Tasks[parentID]; !exists {
			// Parent ID doesn't exist as a task - this is a root
			graph.Roots = append(graph.Roots, taskID)
		}
	}
}

// getParentID extracts the parent task ID from a hierarchical task ID.
// Examples:
//   - "1.2" -> "1"
//   - "1.2.3" -> "1.2"
//   - "1" -> "" (root task)
//   - "2.1" -> "2"
func getParentID(taskID string) string {
	// Find the last dot separator
	lastDot := strings.LastIndex(taskID, ".")
	if lastDot == -1 {
		// No dot found, this is a root task
		return ""
	}
	// Return everything before the last dot
	return taskID[:lastDot]
}

// TopologicalSort performs a topological sort on the task graph and returns
// execution stages where tasks in each stage can run in parallel.
//
// The returned 2D slice structure:
//   - Outer slice: execution stages (must be executed in order)
//   - Inner slices: task IDs that can run in parallel within that stage
//
// Example:
//   - Tasks: 1.1, 1.2, 2.1, 2.2
//   - Stage 0: ["1.1", "2.1"] - roots can run in parallel
//   - Stage 1: ["1.2", "2.2"] - children run after their parents
//
// Returns an error if a circular dependency is detected (though this shouldn't
// happen with our hierarchical task ID structure).
func (g *TaskGraph) TopologicalSort() ([][]string, error) {
	inDegree, dependents := g.computeDependencyGraph()

	return g.kahnTopologicalSort(inDegree, dependents)
}

// computeDependencyGraph builds the in-degree and dependents maps for topological sort.
func (g *TaskGraph) computeDependencyGraph() (map[string]int, map[string][]string) {
	inDegree := make(map[string]int)
	dependents := make(map[string][]string)

	// Initialize in-degree for all tasks
	for taskID := range g.Tasks {
		inDegree[taskID] = 0
	}

	// Calculate in-degrees based on dependencies
	for taskID := range g.Tasks {
		deps := g.getDependencies(taskID)
		inDegree[taskID] = len(deps)

		// For each dependency, track that this task is a dependent
		for _, depID := range deps {
			dependents[depID] = append(dependents[depID], taskID)
		}
	}

	return inDegree, dependents
}

// kahnTopologicalSort implements Kahn's algorithm for topological sorting.
func (g *TaskGraph) kahnTopologicalSort(
	inDegree map[string]int,
	dependents map[string][]string,
) ([][]string, error) {
	visited := make(map[string]bool)
	stages := make([][]string, 0)

	// Process tasks level by level
	for len(visited) < len(g.Tasks) {
		currentStage := g.findTasksWithNoDependencies(visited, inDegree)

		if len(currentStage) == 0 {
			return nil, ErrCircularDependency
		}

		g.markStageAsVisited(currentStage, visited, inDegree, dependents)
		stages = append(stages, currentStage)
	}

	return stages, nil
}

// findTasksWithNoDependencies returns tasks that have no pending dependencies.
func (g *TaskGraph) findTasksWithNoDependencies(
	visited map[string]bool,
	inDegree map[string]int,
) []string {
	currentStage := make([]string, 0)

	for taskID := range g.Tasks {
		if !visited[taskID] && inDegree[taskID] == 0 {
			currentStage = append(currentStage, taskID)
		}
	}

	return currentStage
}

// markStageAsVisited marks tasks in the current stage as visited and updates in-degrees.
func (*TaskGraph) markStageAsVisited(
	stage []string,
	visited map[string]bool,
	inDegree map[string]int,
	dependents map[string][]string,
) {
	for _, taskID := range stage {
		visited[taskID] = true

		// Reduce in-degree for all dependents
		for _, dependent := range dependents[taskID] {
			inDegree[dependent]--
		}
	}
}

// getDependencies returns the list of task IDs that the given task depends on.
// A task depends on:
// 1. Its parent task (if exists in the graph)
// 2. Its previous sibling (task with same parent but lower sequence number)
//
// Examples:
//   - "1.2" depends on "1.1" (previous sibling) and potentially "1" (parent)
//   - "1.1" depends on "1" (parent, if exists)
//   - "2.1" is independent of "1.x" tasks
func (g *TaskGraph) getDependencies(taskID string) []string {
	deps := make([]string, 0)

	// Check for parent dependency
	parentID := getParentID(taskID)
	if parentID != "" {
		if _, exists := g.Tasks[parentID]; exists {
			deps = append(deps, parentID)
		}
	}

	// Check for previous sibling dependency
	// Extract the last component to find previous sibling
	parts := strings.Split(taskID, ".")
	if len(parts) > 0 {
		lastPart := parts[len(parts)-1]
		// Try to parse as integer to find previous sibling
		var currentNum int
		if _, err := fmt.Sscanf(lastPart, "%d", &currentNum); err == nil && currentNum > 1 {
			// Build previous sibling ID
			parts[len(parts)-1] = fmt.Sprintf("%d", currentNum-1)
			prevSiblingID := strings.Join(parts, ".")

			// Only add if the previous sibling exists
			if _, exists := g.Tasks[prevSiblingID]; exists {
				deps = append(deps, prevSiblingID)
			}
		}
	}

	return deps
}

// GetRootPrefix extracts the root prefix (first component) of a task ID.
// This is used to identify which root tree a task belongs to, enabling
// parallel execution detection for tasks in different trees.
//
// Examples:
//   - "1" -> "1"
//   - "1.2" -> "1"
//   - "1.2.3" -> "1"
//   - "2.1" -> "2"
//   - "" -> "" (empty input)
func (*TaskGraph) GetRootPrefix(taskID string) string {
	if taskID == "" {
		return ""
	}

	const dotSeparator = "."
	// Find the first dot separator
	firstDot := strings.Index(taskID, dotSeparator)
	if firstDot == -1 {
		// No dot found, the entire ID is the root
		return taskID
	}

	// Return everything before the first dot
	return taskID[:firstDot]
}

// CanRunInParallel determines if two tasks can run in parallel.
// Tasks can run in parallel if:
// - They have different root prefixes (e.g., "1.x" and "2.x")
// - They are not in a parent-child relationship
//
// Tasks cannot run in parallel if:
// - They have the same root prefix (same tree, sequential execution)
// - One is a parent/ancestor of the other
// - One is a child/descendant of the other
//
// Examples of parallel tasks:
//   - "1.1" and "2.1" (different roots: 1 vs 2)
//   - "1" and "2" (different roots)
//   - "1.1.1" and "2.3.4" (different roots)
//
// Examples of non-parallel tasks:
//   - "1.1" and "1.2" (same root, sequential siblings)
//   - "1" and "1.1" (parent-child relationship)
//   - "1.1" and "1.1.2" (parent-child relationship)
func (g *TaskGraph) CanRunInParallel(taskID1, taskID2 string) bool {
	// Empty task IDs cannot be compared
	if taskID1 == "" || taskID2 == "" {
		return false
	}

	// Same task cannot run in parallel with itself
	if taskID1 == taskID2 {
		return false
	}

	// Check if they have different root prefixes
	root1 := g.GetRootPrefix(taskID1)
	root2 := g.GetRootPrefix(taskID2)

	if root1 != root2 {
		// Different roots means they're in independent trees
		return true
	}

	// Same root - check for parent-child relationship
	// If one is a prefix of the other, they have a dependency
	if isParentChild(taskID1, taskID2) {
		return false
	}

	// Same root but not parent-child means sequential siblings
	// Sequential siblings in the same tree cannot run in parallel
	return false
}

// isParentChild checks if two task IDs have a parent-child relationship.
// Returns true if one task is a parent/ancestor or child/descendant of the other.
func isParentChild(taskID1, taskID2 string) bool {
	const dotSeparator = "."
	// Check if taskID1 is a parent of taskID2
	if strings.HasPrefix(taskID2, taskID1+dotSeparator) {
		return true
	}

	// Check if taskID2 is a parent of taskID1
	if strings.HasPrefix(taskID1, taskID2+dotSeparator) {
		return true
	}

	return false
}
