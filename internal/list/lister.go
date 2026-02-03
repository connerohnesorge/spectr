package list

import (
	"fmt"
	"path/filepath"
	"sort"
	"sync"

	"github.com/connerohnesorge/spectr/internal/discovery"
	"github.com/connerohnesorge/spectr/internal/git"
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

// processChange extracts info for a single change directory
func (l *Lister) processChange(changeID string) ChangeInfo {
	changeDir := filepath.Join(
		l.projectPath,
		"spectr",
		"changes",
		changeID,
	)
	proposalPath := filepath.Join(changeDir, "proposal.md")

	// Extract title
	title, err := parsers.ExtractTitle(proposalPath)
	if err != nil || title == "" {
		title = changeID
	}

	// Count tasks
	taskStatus, err := parsers.CountTasks(changeDir)
	if err != nil {
		taskStatus = parsers.TaskStatus{}
	}

	// Count deltas
	deltaCount, _ := parsers.CountDeltas(changeDir)

	return ChangeInfo{
		ID:         changeID,
		Title:      title,
		DeltaCount: deltaCount,
		TaskStatus: taskStatus,
	}
}

// ListChanges retrieves information about all active changes
func (l *Lister) ListChanges() ([]ChangeInfo, error) {
	changeIDs, err := discovery.GetActiveChanges(l.projectPath)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to discover changes: %w",
			err,
		)
	}

	if len(changeIDs) == 0 {
		return nil, nil
	}

	// Process changes in parallel
	type changeResult struct {
		idx    int
		change ChangeInfo
	}

	results := make(chan changeResult, len(changeIDs))
	var wg sync.WaitGroup

	for i, id := range changeIDs {
		wg.Add(1)
		go func(idx int, changeID string) {
			defer wg.Done()
			results <- changeResult{
				idx:    idx,
				change: l.processChange(changeID),
			}
		}(i, id)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results maintaining original order
	changes := make([]ChangeInfo, len(changeIDs))
	for r := range results {
		changes[r.idx] = r.change
	}

	return changes, nil
}

// processSpec extracts info for a single spec directory
func (l *Lister) processSpec(specID string) SpecInfo {
	specPath := filepath.Join(
		l.projectPath,
		"spectr",
		"specs",
		specID,
		"spec.md",
	)

	// Extract title
	title, err := parsers.ExtractTitle(specPath)
	if err != nil || title == "" {
		title = specID
	}

	// Count requirements
	reqCount, _ := parsers.CountRequirements(specPath)

	return SpecInfo{
		ID:               specID,
		Title:            title,
		RequirementCount: reqCount,
	}
}

// ListSpecs retrieves information about all specs
func (l *Lister) ListSpecs() ([]SpecInfo, error) {
	specIDs, err := discovery.GetSpecs(l.projectPath)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to discover specs: %w",
			err,
		)
	}

	if len(specIDs) == 0 {
		return nil, nil
	}

	// Process specs in parallel
	type specResult struct {
		idx  int
		spec SpecInfo
	}

	results := make(chan specResult, len(specIDs))
	var wg sync.WaitGroup

	for i, id := range specIDs {
		wg.Add(1)
		go func(idx int, specID string) {
			defer wg.Done()
			results <- specResult{
				idx:  idx,
				spec: l.processSpec(specID),
			}
		}(i, id)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results maintaining original order
	specs := make([]SpecInfo, len(specIDs))
	for r := range results {
		specs[r.idx] = r.spec
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
func (l *Lister) ListAll(
	opts *ListAllOptions,
) (ItemList, error) {
	// Use default options if none provided
	options := opts
	if options == nil {
		options = &ListAllOptions{
			SortByID: true,
		}
	}

	var items ItemList

	// Load changes if not filtered out
	if options.FilterType == nil ||
		*options.FilterType == ItemTypeChange {
		changes, err := l.ListChanges()
		if err != nil {
			return nil, fmt.Errorf(
				"failed to list changes: %w",
				err,
			)
		}
		for _, change := range changes {
			items = append(
				items,
				NewChangeItem(change),
			)
		}
	}

	// Load specs if not filtered out
	if options.FilterType == nil ||
		*options.FilterType == ItemTypeSpec {
		specs, err := l.ListSpecs()
		if err != nil {
			return nil, fmt.Errorf(
				"failed to list specs: %w",
				err,
			)
		}
		for _, spec := range specs {
			items = append(
				items,
				NewSpecItem(spec),
			)
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

// FilterChangesNotOnRef filters the given changes to only include those
// whose paths do NOT exist on the specified git ref. This is useful for
// identifying unmerged changes that haven't been merged to the main branch yet.
// The ref should be a full ref like "origin/main" or "origin/master".
func FilterChangesNotOnRef(
	changes []ChangeInfo,
	ref string,
) ([]ChangeInfo, error) {
	var unmerged []ChangeInfo

	for _, change := range changes {
		changePath := filepath.Join(
			"spectr",
			"changes",
			change.ID,
		)
		exists, err := git.PathExistsOnRef(
			ref,
			changePath,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to check if change %q exists on %s: %w",
				change.ID,
				ref,
				err,
			)
		}

		// Only include changes that do NOT exist on the ref
		if !exists {
			unmerged = append(unmerged, change)
		}
	}

	return unmerged, nil
}
