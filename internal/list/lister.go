package list

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/connerohnesorge/spectr/internal/discovery"
	"github.com/connerohnesorge/spectr/internal/parsers"
)

// Lister handles listing operations for changes and specs
type Lister struct {
	projectPath string
}

// NewLister creates a new Lister for the given project path
func NewLister(projectPath string) *Lister {
	return &Lister{projectPath: projectPath}
}

// ListChanges retrieves information about all active changes
func (l *Lister) ListChanges() ([]ChangeInfo, error) {
	changeIDs, err := discovery.GetActiveChanges(l.projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to discover changes: %w", err)
	}

	var changes []ChangeInfo
	for _, id := range changeIDs {
		changeDir := filepath.Join(l.projectPath, "spectr", "changes", id)
		proposalPath := filepath.Join(changeDir, "proposal.md")
		tasksPath := filepath.Join(changeDir, "tasks.md")

		// Extract title
		title, err := parsers.ExtractTitle(proposalPath)
		if err != nil || title == "" {
			// Fallback to ID if title extraction fails
			title = id
		}

		// Count tasks
		taskStatus, err := parsers.CountTasks(tasksPath)
		if err != nil {
			// If error reading tasks, use zero status
			taskStatus = parsers.TaskStatus{Total: 0, Completed: 0}
		}

		// Count deltas
		deltaCount, err := parsers.CountDeltas(changeDir)
		if err != nil {
			deltaCount = 0
		}

		changes = append(changes, ChangeInfo{
			ID:         id,
			Title:      title,
			DeltaCount: deltaCount,
			TaskStatus: taskStatus,
		})
	}

	return changes, nil
}

// ListSpecs retrieves information about all specs
func (l *Lister) ListSpecs() ([]SpecInfo, error) {
	specIDs, err := discovery.GetSpecs(l.projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to discover specs: %w", err)
	}

	var specs []SpecInfo
	for _, id := range specIDs {
		specPath := filepath.Join(
			l.projectPath,
			"spectr",
			"specs",
			id,
			"spec.md",
		)

		// Extract title
		title, err := parsers.ExtractTitle(specPath)
		if err != nil || title == "" {
			// Fallback to ID if title extraction fails
			title = id
		}

		// Count requirements
		reqCount, err := parsers.CountRequirements(specPath)
		if err != nil {
			reqCount = 0
		}

		specs = append(specs, SpecInfo{
			ID:               id,
			Title:            title,
			RequirementCount: reqCount,
		})
	}

	return specs, nil
}

// ListAllOptions contains optional filtering and sorting parameters for ListAll
type ListAllOptions struct {
	// FilterType specifies which types to include (nil = all types)
	FilterType *ItemType
	// SortByID sorts items alphabetically by ID (default: true)
	SortByID bool
}

// ListAll retrieves all changes and specs as a unified ItemList
func (l *Lister) ListAll(opts *ListAllOptions) (ItemList, error) {
	// Use default options if none provided
	options := opts
	if options == nil {
		options = &ListAllOptions{
			SortByID: true,
		}
	}

	var items ItemList

	// Load changes if not filtered out
	if options.FilterType == nil || *options.FilterType == ItemTypeChange {
		changes, err := l.ListChanges()
		if err != nil {
			return nil, fmt.Errorf("failed to list changes: %w", err)
		}
		for _, change := range changes {
			items = append(items, NewChangeItem(change))
		}
	}

	// Load specs if not filtered out
	if options.FilterType == nil || *options.FilterType == ItemTypeSpec {
		specs, err := l.ListSpecs()
		if err != nil {
			return nil, fmt.Errorf("failed to list specs: %w", err)
		}
		for _, spec := range specs {
			items = append(items, NewSpecItem(spec))
		}
	}

	// Sort by ID if requested
	if options.SortByID {
		sort.Slice(items, func(i, j int) bool {
			return items[i].ID() < items[j].ID()
		})
	}

	return items, nil
}
