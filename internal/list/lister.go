package list

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/connerohnesorge/spectr/internal/discovery"
	"github.com/connerohnesorge/spectr/internal/git"
	"github.com/connerohnesorge/spectr/internal/parsers"
)

// Lister handles listing operations for changes and specs
type Lister struct {
	projectPath string
	// rootPath is the relative path from cwd to this root (empty for single root)
	rootPath string
	// absPath is the absolute path to the project root
	absPath string
}

// NewLister creates a new Lister for the given project path
func NewLister(projectPath string) *Lister {
	return &Lister{
		projectPath: projectPath,
		absPath:     projectPath,
	}
}

// NewListerWithRoot creates a new Lister with root path information
func NewListerWithRoot(projectPath, rootPath string) *Lister {
	return &Lister{
		projectPath: projectPath,
		rootPath:    rootPath,
		absPath:     projectPath,
	}
}

// ListChanges retrieves information about all active changes
func (l *Lister) ListChanges() ([]ChangeInfo, error) {
	changeIDs, err := discovery.GetActiveChanges(
		l.projectPath,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to discover changes: %w",
			err,
		)
	}

	changes := make(
		[]ChangeInfo,
		0,
		len(changeIDs),
	)
	for _, id := range changeIDs {
		changeDir := filepath.Join(
			l.projectPath,
			"spectr",
			"changes",
			id,
		)
		proposalPath := filepath.Join(
			changeDir,
			"proposal.md",
		)

		// Extract title
		title, err := parsers.ExtractTitle(
			proposalPath,
		)
		if err != nil || title == "" {
			// Fallback to ID if title extraction fails
			title = id
		}

		// Count tasks (from tasks.json or tasks.md)
		taskStatus, err := parsers.CountTasks(
			changeDir,
		)
		if err != nil {
			// If error reading tasks, use zero status
			taskStatus = parsers.TaskStatus{
				Total: 0, Completed: 0, InProgress: 0,
			}
		}

		// Count deltas
		deltaCount, err := parsers.CountDeltas(
			changeDir,
		)
		if err != nil {
			deltaCount = 0
		}

		changes = append(changes, ChangeInfo{
			ID:          id,
			Title:       title,
			DeltaCount:  deltaCount,
			TaskStatus:  taskStatus,
			RootPath:    l.rootPath,
			RootAbsPath: l.absPath,
		})
	}

	return changes, nil
}

// ListSpecs retrieves information about all specs
func (l *Lister) ListSpecs() ([]SpecInfo, error) {
	specIDs, err := discovery.GetSpecs(
		l.projectPath,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to discover specs: %w",
			err,
		)
	}

	specs := make([]SpecInfo, 0, len(specIDs))
	for _, id := range specIDs {
		specPath := filepath.Join(
			l.projectPath,
			"spectr",
			"specs",
			id,
			"spec.md",
		)

		// Extract title
		title, err := parsers.ExtractTitle(
			specPath,
		)
		if err != nil || title == "" {
			// Fallback to ID if title extraction fails
			title = id
		}

		// Count requirements
		reqCount, err := parsers.CountRequirements(
			specPath,
		)
		if err != nil {
			reqCount = 0
		}

		specs = append(specs, SpecInfo{
			ID:               id,
			Title:            title,
			RequirementCount: reqCount,
			RootPath:         l.rootPath,
			RootAbsPath:      l.absPath,
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
		for i := range changes {
			items = append(
				items,
				NewChangeItem(&changes[i]),
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

// MultiRootLister aggregates listing results from multiple spectr roots.
type MultiRootLister struct {
	listers []*Lister
}

// NewMultiRootLister creates a lister that aggregates from multiple roots.
func NewMultiRootLister(roots []discovery.SpectrRoot) *MultiRootLister {
	listers := make([]*Lister, len(roots))
	for i, root := range roots {
		listers[i] = NewListerWithRoot(root.Path, root.RelativeTo)
	}

	return &MultiRootLister{listers: listers}
}

// ListChanges retrieves changes from all roots.
func (m *MultiRootLister) ListChanges() ([]ChangeInfo, error) {
	var allChanges []ChangeInfo

	for _, lister := range m.listers {
		changes, err := lister.ListChanges()
		if err != nil {
			return nil, fmt.Errorf(
				"failed to list changes from %s: %w",
				lister.rootPath,
				err,
			)
		}
		allChanges = append(allChanges, changes...)
	}

	// Sort by ID for consistency
	sort.Slice(allChanges, func(i, j int) bool {
		return allChanges[i].ID < allChanges[j].ID
	})

	return allChanges, nil
}

// ListSpecs retrieves specs from all roots.
func (m *MultiRootLister) ListSpecs() ([]SpecInfo, error) {
	var allSpecs []SpecInfo

	for _, lister := range m.listers {
		specs, err := lister.ListSpecs()
		if err != nil {
			return nil, fmt.Errorf(
				"failed to list specs from %s: %w",
				lister.rootPath,
				err,
			)
		}
		allSpecs = append(allSpecs, specs...)
	}

	// Sort by ID for consistency
	sort.Slice(allSpecs, func(i, j int) bool {
		return allSpecs[i].ID < allSpecs[j].ID
	})

	return allSpecs, nil
}

// ListAll retrieves all items from all roots.
func (m *MultiRootLister) ListAll(opts *ListAllOptions) (ItemList, error) {
	var items ItemList

	for _, lister := range m.listers {
		rootItems, err := lister.ListAll(opts)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to list items from %s: %w",
				lister.rootPath,
				err,
			)
		}
		items = append(items, rootItems...)
	}

	// Sort by ID for consistency (if requested by options)
	if opts == nil || opts.SortByID {
		sort.Slice(items, func(i, j int) bool {
			return items[i].ID() < items[j].ID()
		})
	}

	return items, nil
}

// HasMultipleRoots returns true if there are multiple roots.
func (m *MultiRootLister) HasMultipleRoots() bool {
	return len(m.listers) > 1
}
