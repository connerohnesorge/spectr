package validation

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/connerohnesorge/spectr/internal/domain"
)

func TestNewDependencyGraph(t *testing.T) {
	graph := NewDependencyGraph()
	if graph.Nodes == nil {
		t.Error("expected Nodes map to be initialized")
	}
	if graph.Edges == nil {
		t.Error("expected Edges map to be initialized")
	}
}

func TestDependencyGraph_AddNode(t *testing.T) {
	graph := NewDependencyGraph()

	meta := &domain.ProposalMetadata{
		ID: "test",
		Requires: []domain.Dependency{
			{ID: "dep1"},
			{ID: "dep2"},
		},
	}

	graph.AddNode("test-change", meta)

	if graph.Nodes["test-change"] != meta {
		t.Error("expected node to be added to graph")
	}

	if len(graph.Edges["test-change"]) != 2 {
		t.Errorf("expected 2 edges, got %d", len(graph.Edges["test-change"]))
	}
}

func TestDetectCycles_NoCycles(t *testing.T) {
	graph := NewDependencyGraph()

	// A -> B -> C (linear, no cycles)
	graph.AddNode("A", &domain.ProposalMetadata{
		Requires: []domain.Dependency{{ID: "B"}},
	})
	graph.AddNode("B", &domain.ProposalMetadata{
		Requires: []domain.Dependency{{ID: "C"}},
	})
	graph.AddNode("C", &domain.ProposalMetadata{})

	cycles := DetectCycles(graph)

	if len(cycles) != 0 {
		t.Errorf("expected no cycles, got %v", cycles)
	}
}

func TestDetectCycles_SimpleCycle(t *testing.T) {
	graph := NewDependencyGraph()

	// A -> B -> A (simple cycle)
	graph.AddNode("A", &domain.ProposalMetadata{
		Requires: []domain.Dependency{{ID: "B"}},
	})
	graph.AddNode("B", &domain.ProposalMetadata{
		Requires: []domain.Dependency{{ID: "A"}},
	})

	cycles := DetectCycles(graph)

	if len(cycles) == 0 {
		t.Error("expected to detect cycle")
	}
}

func TestDetectCycles_LongerCycle(t *testing.T) {
	graph := NewDependencyGraph()

	// A -> B -> C -> A (3-node cycle)
	graph.AddNode("A", &domain.ProposalMetadata{
		Requires: []domain.Dependency{{ID: "B"}},
	})
	graph.AddNode("B", &domain.ProposalMetadata{
		Requires: []domain.Dependency{{ID: "C"}},
	})
	graph.AddNode("C", &domain.ProposalMetadata{
		Requires: []domain.Dependency{{ID: "A"}},
	})

	cycles := DetectCycles(graph)

	if len(cycles) == 0 {
		t.Error("expected to detect cycle")
	}
}

func TestValidateDependencies_NoFrontmatter(t *testing.T) {
	tmpDir := t.TempDir()
	changeDir := filepath.Join(tmpDir, "spectr", "changes", "test-change")
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create proposal without frontmatter
	proposalContent := `# Test Proposal

Just a regular proposal without frontmatter.
`
	if err := os.WriteFile(
		filepath.Join(changeDir, "proposal.md"),
		[]byte(proposalContent),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	result, err := ValidateDependencies("test-change", tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Issues) != 0 {
		t.Errorf("expected no issues for proposal without frontmatter, got %v", result.Issues)
	}
}

func TestValidateDependencies_SelfReference(t *testing.T) {
	tmpDir := t.TempDir()
	changeDir := filepath.Join(tmpDir, "spectr", "changes", "test-change")
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	proposalContent := `---
requires:
  - id: test-change
---

# Test Proposal
`
	if err := os.WriteFile(
		filepath.Join(changeDir, "proposal.md"),
		[]byte(proposalContent),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	result, err := ValidateDependencies("test-change", tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check for error-level issues (self-reference creates an ERROR issue)
	hasErrorIssue := false
	foundSelfRef := false
	for _, issue := range result.Issues {
		if issue.Level != LevelError {
			continue
		}
		hasErrorIssue = true
		if issue.Message == "proposal cannot require itself: test-change" {
			foundSelfRef = true
		}
	}

	if !hasErrorIssue {
		t.Error("expected error-level issue for self-reference")
	}
	if !foundSelfRef {
		t.Error("expected self-reference error in issues")
	}
}

func TestValidateDependencies_UnmetDependency(t *testing.T) {
	tmpDir := t.TempDir()
	changeDir := filepath.Join(tmpDir, "spectr", "changes", "test-change")
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	proposalContent := `---
requires:
  - id: nonexistent-dep
---

# Test Proposal
`
	if err := os.WriteFile(
		filepath.Join(changeDir, "proposal.md"),
		[]byte(proposalContent),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	result, err := ValidateDependencies("test-change", tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.HasWarnings() {
		t.Error("expected warning for unmet dependency")
	}

	if len(result.UnmetDependencies["test-change"]) != 1 {
		t.Errorf(
			"expected 1 unmet dependency, got %d",
			len(result.UnmetDependencies["test-change"]),
		)
	}
}

func TestValidateDependencies_ArchivedDependency(t *testing.T) {
	tmpDir := t.TempDir()

	// Create change directory
	changeDir := filepath.Join(tmpDir, "spectr", "changes", "test-change")
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create archive with the required dependency
	archiveDir := filepath.Join(tmpDir, "spectr", "changes", "archive", "2024-01-15-required-dep")
	if err := os.MkdirAll(archiveDir, 0o755); err != nil {
		t.Fatal(err)
	}

	proposalContent := `---
requires:
  - id: required-dep
---

# Test Proposal
`
	if err := os.WriteFile(
		filepath.Join(changeDir, "proposal.md"),
		[]byte(proposalContent),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	result, err := ValidateDependencies("test-change", tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.HasWarnings() {
		t.Errorf("expected no warnings for archived dependency, got %v", result.Issues)
	}
}

func TestValidateDependencies_ActiveDependency(t *testing.T) {
	tmpDir := t.TempDir()

	// Create change directory
	changeDir := filepath.Join(tmpDir, "spectr", "changes", "test-change")
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create active dependency (not archived)
	depDir := filepath.Join(tmpDir, "spectr", "changes", "active-dep")
	if err := os.MkdirAll(depDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(
		filepath.Join(depDir, "proposal.md"),
		[]byte("# Active Dep"),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	proposalContent := `---
requires:
  - id: active-dep
---

# Test Proposal
`
	if err := os.WriteFile(
		filepath.Join(changeDir, "proposal.md"),
		[]byte(proposalContent),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	result, err := ValidateDependencies("test-change", tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.HasWarnings() {
		t.Error("expected warning for active (non-archived) dependency")
	}
}

func TestValidateDependenciesForAccept_NoDependencies(t *testing.T) {
	tmpDir := t.TempDir()
	changeDir := filepath.Join(tmpDir, "spectr", "changes", "test-change")
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	proposalContent := `# Test Proposal

No dependencies.
`
	if err := os.WriteFile(
		filepath.Join(changeDir, "proposal.md"),
		[]byte(proposalContent),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	err := ValidateDependenciesForAccept("test-change", tmpDir)
	if err != nil {
		t.Errorf("expected no error for proposal without dependencies, got: %v", err)
	}
}

func TestValidateDependenciesForAccept_UnmetDependency(t *testing.T) {
	tmpDir := t.TempDir()
	changeDir := filepath.Join(tmpDir, "spectr", "changes", "test-change")
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	proposalContent := `---
requires:
  - id: missing-dep
---

# Test Proposal
`
	if err := os.WriteFile(
		filepath.Join(changeDir, "proposal.md"),
		[]byte(proposalContent),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	err := ValidateDependenciesForAccept("test-change", tmpDir)
	if err == nil {
		t.Error("expected error for unmet dependency")
	}

	unmetErr, ok := err.(*UnmetDependenciesError)
	if !ok {
		t.Fatalf("expected UnmetDependenciesError, got %T", err)
	}

	if len(unmetErr.Dependencies) != 1 {
		t.Errorf("expected 1 unmet dependency, got %d", len(unmetErr.Dependencies))
	}

	if unmetErr.Dependencies[0] != "missing-dep" {
		t.Errorf("expected 'missing-dep', got %q", unmetErr.Dependencies[0])
	}
}

func TestValidateDependenciesForAccept_MetDependency(t *testing.T) {
	tmpDir := t.TempDir()

	// Create change directory
	changeDir := filepath.Join(tmpDir, "spectr", "changes", "test-change")
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create archive with the required dependency
	archiveDir := filepath.Join(tmpDir, "spectr", "changes", "archive", "2024-01-15-required-dep")
	if err := os.MkdirAll(archiveDir, 0o755); err != nil {
		t.Fatal(err)
	}

	proposalContent := `---
requires:
  - id: required-dep
---

# Test Proposal
`
	if err := os.WriteFile(
		filepath.Join(changeDir, "proposal.md"),
		[]byte(proposalContent),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	err := ValidateDependenciesForAccept("test-change", tmpDir)
	if err != nil {
		t.Errorf("expected no error when dependency is archived, got: %v", err)
	}
}

func TestUnmetDependenciesError_SingleDep(t *testing.T) {
	err := &UnmetDependenciesError{
		ChangeID:     "my-change",
		Dependencies: []string{"dep-a"},
	}

	expected := "cannot accept 'my-change': required dependency 'dep-a' is not archived"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestUnmetDependenciesError_MultipleDeps(t *testing.T) {
	err := &UnmetDependenciesError{
		ChangeID:     "my-change",
		Dependencies: []string{"dep-a", "dep-b"},
	}

	expected := "cannot accept 'my-change': required dependencies not archived: dep-a, dep-b"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestDependencyValidationResult_HasErrors(t *testing.T) {
	tests := []struct {
		name     string
		result   DependencyValidationResult
		expected bool
	}{
		{
			name:     "no errors",
			result:   DependencyValidationResult{},
			expected: false,
		},
		{
			name: "has cycles",
			result: DependencyValidationResult{
				Cycles: [][]string{{"A", "B", "A"}},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.result.HasErrors() != tt.expected {
				t.Errorf("HasErrors() = %v, want %v", tt.result.HasErrors(), tt.expected)
			}
		})
	}
}

func TestDependencyValidationResult_HasWarnings(t *testing.T) {
	tests := []struct {
		name     string
		result   DependencyValidationResult
		expected bool
	}{
		{
			name:     "no warnings",
			result:   DependencyValidationResult{},
			expected: false,
		},
		{
			name: "has unmet dependencies",
			result: DependencyValidationResult{
				UnmetDependencies: map[string][]UnmetDependency{
					"change": {{ID: "dep"}},
				},
			},
			expected: true,
		},
		{
			name: "empty unmet dependencies map",
			result: DependencyValidationResult{
				UnmetDependencies: map[string][]UnmetDependency{
					"change": {},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.result.HasWarnings() != tt.expected {
				t.Errorf("HasWarnings() = %v, want %v", tt.result.HasWarnings(), tt.expected)
			}
		})
	}
}

func TestBuildDependencyGraph(t *testing.T) {
	tmpDir := t.TempDir()
	changesDir := filepath.Join(tmpDir, "spectr", "changes")

	// Create changes with dependencies
	changes := []struct {
		id      string
		content string
	}{
		{
			id: "change-a",
			content: `---
requires:
  - id: change-b
---

# Change A
`,
		},
		{
			id: "change-b",
			content: `---
requires:
  - id: change-c
---

# Change B
`,
		},
		{
			id:      "change-c",
			content: "# Change C (no dependencies)",
		},
	}

	for _, c := range changes {
		changeDir := filepath.Join(changesDir, c.id)
		if err := os.MkdirAll(changeDir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(
			filepath.Join(changeDir, "proposal.md"),
			[]byte(c.content),
			0o644,
		); err != nil {
			t.Fatal(err)
		}
	}

	graph, err := BuildDependencyGraph(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(graph.Nodes) != 3 {
		t.Errorf("expected 3 nodes, got %d", len(graph.Nodes))
	}

	// Check edges
	if len(graph.Edges["change-a"]) != 1 || graph.Edges["change-a"][0] != "change-b" {
		t.Errorf("expected change-a to require change-b, got %v", graph.Edges["change-a"])
	}

	if len(graph.Edges["change-b"]) != 1 || graph.Edges["change-b"][0] != "change-c" {
		t.Errorf("expected change-b to require change-c, got %v", graph.Edges["change-b"])
	}

	if len(graph.Edges["change-c"]) != 0 {
		t.Errorf("expected change-c to have no dependencies, got %v", graph.Edges["change-c"])
	}
}
