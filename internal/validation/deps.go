// Package validation provides validation logic for specs and changes.
//
//nolint:revive // This file has many function arguments and minor style issues that are acceptable
package validation

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/connerohnesorge/spectr/internal/discovery"
	"github.com/connerohnesorge/spectr/internal/domain"
)

// Common path components for spectr directory structure.
const (
	spectrDir    = "spectr"
	changesDir   = "changes"
	proposalFile = "proposal.md"
)

// DependencyGraph represents a directed graph of proposal dependencies.
type DependencyGraph struct {
	// Nodes maps change ID to ProposalMetadata
	Nodes map[string]*domain.ProposalMetadata
	// Edges maps change ID to list of required change IDs
	Edges map[string][]string
}

// NewDependencyGraph creates an empty dependency graph.
func NewDependencyGraph() *DependencyGraph {
	return &DependencyGraph{
		Nodes: make(map[string]*domain.ProposalMetadata),
		Edges: make(map[string][]string),
	}
}

// AddNode adds a proposal to the graph.
func (g *DependencyGraph) AddNode(changeID string, meta *domain.ProposalMetadata) {
	g.Nodes[changeID] = meta
	if meta != nil {
		g.Edges[changeID] = meta.RequiredIDs()
	}
}

// DependencyValidationResult contains the results of dependency validation.
type DependencyValidationResult struct {
	// UnmetDependencies maps change ID to list of unmet dependency IDs
	UnmetDependencies map[string][]UnmetDependency
	// Cycles contains detected dependency cycles
	Cycles [][]string
	// Issues contains validation issues (errors and warnings)
	Issues []ValidationIssue
}

// UnmetDependency represents a dependency that is not yet archived.
type UnmetDependency struct {
	ID     string
	Reason string
	Status discovery.ChangeStatus
}

// HasErrors returns true if there are blocking errors (cycles or self-references).
func (r *DependencyValidationResult) HasErrors() bool {
	return len(r.Cycles) > 0
}

// HasWarnings returns true if there are unmet dependencies.
func (r *DependencyValidationResult) HasWarnings() bool {
	for _, deps := range r.UnmetDependencies {
		if len(deps) > 0 {
			return true
		}
	}

	return false
}

// BuildDependencyGraph constructs a dependency graph from all active proposals.
func BuildDependencyGraph(projectRoot string) (*DependencyGraph, error) {
	graph := NewDependencyGraph()

	// Get all active changes
	changeIDs, err := discovery.GetActiveChangeIDs(projectRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to get active changes: %w", err)
	}

	// Parse frontmatter from each proposal
	for _, changeID := range changeIDs {
		proposalPath := filepath.Join(
			projectRoot,
			"spectr",
			"changes",
			changeID,
			"proposal.md",
		)

		meta, err := domain.ParseProposalFrontmatterFromFile(proposalPath)
		if err != nil {
			// Log warning but continue with other proposals
			graph.AddNode(changeID, &domain.ProposalMetadata{})

			continue
		}

		graph.AddNode(changeID, meta)
	}

	return graph, nil
}

// DetectCycles finds all cycles in the dependency graph using DFS with coloring.
// Returns a list of cycles, where each cycle is a list of change IDs forming the cycle.
func DetectCycles(graph *DependencyGraph) [][]string {
	// Color states: 0=white (unvisited), 1=gray (in progress), 2=black (complete)
	colors := make(map[string]int)
	parent := make(map[string]string)
	var cycles [][]string

	var dfs func(node string)
	dfs = func(node string) {
		colors[node] = 1 // Mark as in-progress (gray)

		for _, dep := range graph.Edges[node] {
			switch colors[dep] {
			case 1:
				// Found a cycle - trace back to find the cycle path
				cycle := []string{dep}
				curr := node
				for curr != dep {
					cycle = append([]string{curr}, cycle...)
					curr = parent[curr]
				}
				cycle = append([]string{dep}, cycle...)
				cycles = append(cycles, cycle)
			case 0:
				parent[dep] = node
				dfs(dep)
			}
		}

		colors[node] = 2 // Mark as complete (black)
	}

	// Run DFS from each unvisited node
	for node := range graph.Nodes {
		if colors[node] == 0 {
			dfs(node)
		}
	}

	return cycles
}

// ValidateDependencies checks if all required dependencies are archived.
// Returns validation issues for unmet dependencies and cycles.
func ValidateDependencies(
	changeID string,
	projectRoot string,
) (*DependencyValidationResult, error) {
	result := &DependencyValidationResult{
		UnmetDependencies: make(map[string][]UnmetDependency),
		Issues:            make([]ValidationIssue, 0),
	}

	// Get proposal metadata
	proposalPath := filepath.Join(
		projectRoot,
		"spectr",
		"changes",
		changeID,
		"proposal.md",
	)

	meta, err := domain.ParseProposalFrontmatterFromFile(proposalPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse proposal frontmatter: %w", err)
	}

	// Check for self-reference
	if err := domain.ValidateProposalMetadata(meta, changeID); err != nil {
		result.Issues = append(result.Issues, ValidationIssue{
			Level:   LevelError,
			Path:    proposalFile,
			Message: err.Error(),
		})

		return result, nil
	}

	// Check each required dependency
	for _, dep := range meta.Requires {
		status, err := discovery.GetChangeStatus(dep.ID, projectRoot)
		if err != nil {
			return nil, fmt.Errorf("failed to get status for %s: %w", dep.ID, err)
		}

		switch status {
		case discovery.ChangeStatusArchived:
			// Dependency is met - no issue
		case discovery.ChangeStatusActive:
			// Dependency exists but not archived - warning
			result.UnmetDependencies[changeID] = append(
				result.UnmetDependencies[changeID],
				UnmetDependency{
					ID:     dep.ID,
					Reason: dep.Reason,
					Status: status,
				},
			)
			result.Issues = append(result.Issues, ValidationIssue{
				Level: LevelWarning,
				Path:  proposalFile,
				Message: fmt.Sprintf(
					"Dependency '%s' is not yet archived (currently active)",
					dep.ID,
				),
			})
		case discovery.ChangeStatusUnknown:
			// Dependency not found anywhere - warning
			result.UnmetDependencies[changeID] = append(
				result.UnmetDependencies[changeID],
				UnmetDependency{
					ID:     dep.ID,
					Reason: dep.Reason,
					Status: status,
				},
			)
			result.Issues = append(result.Issues, ValidationIssue{
				Level:   LevelWarning,
				Path:    proposalFile,
				Message: fmt.Sprintf("Dependency '%s' not found", dep.ID),
			})
		}
	}

	// Build full graph and check for cycles
	graph, err := BuildDependencyGraph(projectRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to build dependency graph: %w", err)
	}

	cycles := DetectCycles(graph)

	// Filter cycles that include this change
	for _, cycle := range cycles {
		for _, node := range cycle {
			if node == changeID {
				result.Cycles = append(result.Cycles, cycle)
				cyclePath := strings.Join(cycle, " → ")
				result.Issues = append(result.Issues, ValidationIssue{
					Level:   LevelError,
					Path:    proposalFile,
					Message: fmt.Sprintf("Circular dependency detected: %s", cyclePath),
				})

				break
			}
		}
	}

	return result, nil
}

// ValidateDependenciesForAccept performs strict dependency validation for accept.
// Returns an error if any required dependencies are not archived.
func ValidateDependenciesForAccept(
	changeID string,
	projectRoot string,
) error {
	proposalPath := filepath.Join(
		projectRoot,
		"spectr",
		"changes",
		changeID,
		"proposal.md",
	)

	meta, err := domain.ParseProposalFrontmatterFromFile(proposalPath)
	if err != nil {
		return fmt.Errorf("failed to parse proposal frontmatter: %w", err)
	}

	// Check for self-reference
	if err := domain.ValidateProposalMetadata(meta, changeID); err != nil {
		return err
	}

	// No dependencies to check
	if len(meta.Requires) == 0 {
		return nil
	}

	var unmetDeps []string

	for _, dep := range meta.Requires {
		archived, err := discovery.IsChangeArchived(dep.ID, projectRoot)
		if err != nil {
			return fmt.Errorf("failed to check archive status for %s: %w", dep.ID, err)
		}

		if !archived {
			unmetDeps = append(unmetDeps, dep.ID)
		}
	}

	if len(unmetDeps) > 0 {
		return &UnmetDependenciesError{
			ChangeID:     changeID,
			Dependencies: unmetDeps,
		}
	}

	return nil
}

// UnmetDependenciesError is returned when required dependencies are not archived.
type UnmetDependenciesError struct {
	ChangeID     string
	Dependencies []string
}

func (e *UnmetDependenciesError) Error() string {
	if len(e.Dependencies) == 1 {
		return fmt.Sprintf(
			"cannot accept '%s': required dependency '%s' is not archived",
			e.ChangeID,
			e.Dependencies[0],
		)
	}

	return fmt.Sprintf(
		"cannot accept '%s': required dependencies not archived: %s",
		e.ChangeID,
		strings.Join(e.Dependencies, ", "),
	)
}

// GetProposalDependencies returns the parsed dependencies for a proposal.
func GetProposalDependencies(
	changeID string,
	projectRoot string,
) (*domain.ProposalMetadata, error) {
	proposalPath := filepath.Join(
		projectRoot,
		"spectr",
		"changes",
		changeID,
		"proposal.md",
	)

	return domain.ParseProposalFrontmatterFromFile(proposalPath)
}
