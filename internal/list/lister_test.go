package list

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/connerohnesorge/spectr/internal/parsers"
)

func TestListChanges(t *testing.T) {
	tmpDir := t.TempDir()
	changesDir := filepath.Join(
		tmpDir,
		"spectr",
		"changes",
	)
	if err := os.MkdirAll(changesDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create test change with all components
	changeDir := filepath.Join(
		changesDir,
		"add-feature",
	)
	if err := os.MkdirAll(filepath.Join(changeDir, "specs", "test-spec"), 0o755); err != nil {
		t.Fatal(err)
	}

	// Write proposal.md
	proposalContent := `# Change: Add Amazing Feature

More details here.`
	if err := os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte(proposalContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Write tasks.md
	tasksContent := `## Tasks
- [x] Task 1
- [ ] Task 2
- [x] Task 3`
	if err := os.WriteFile(filepath.Join(changeDir, "tasks.md"), []byte(tasksContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Write spec delta
	specContent := `## ADDED Requirements
### Requirement: New Feature

## MODIFIED Requirements
### Requirement: Updated Feature`
	if err := os.WriteFile(filepath.Join(changeDir, "specs", "test-spec", "spec.md"), []byte(specContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Test listing
	lister := NewLister(tmpDir)
	changes, err := lister.ListChanges()
	if err != nil {
		t.Fatalf("ListChanges failed: %v", err)
	}

	if len(changes) != 1 {
		t.Fatalf(
			"Expected 1 change, got %d",
			len(changes),
		)
	}

	change := changes[0]
	if change.ID != "add-feature" {
		t.Errorf(
			"Expected ID 'add-feature', got %q",
			change.ID,
		)
	}
	if change.Title != "Add Amazing Feature" {
		t.Errorf(
			"Expected title 'Add Amazing Feature', got %q",
			change.Title,
		)
	}
	if change.DeltaCount != 2 {
		t.Errorf(
			"Expected delta count 2, got %d",
			change.DeltaCount,
		)
	}
	if change.TaskStatus.Total != 3 {
		t.Errorf(
			"Expected 3 total tasks, got %d",
			change.TaskStatus.Total,
		)
	}
	if change.TaskStatus.Completed != 2 {
		t.Errorf(
			"Expected 2 completed tasks, got %d",
			change.TaskStatus.Completed,
		)
	}
}

func TestListChanges_NoChanges(t *testing.T) {
	tmpDir := t.TempDir()
	lister := NewLister(tmpDir)
	changes, err := lister.ListChanges()
	if err != nil {
		t.Fatalf("ListChanges failed: %v", err)
	}
	if len(changes) != 0 {
		t.Errorf(
			"Expected empty list, got %d changes",
			len(changes),
		)
	}
}

func TestListChanges_FallbackTitle(t *testing.T) {
	tmpDir := t.TempDir()
	changesDir := filepath.Join(
		tmpDir,
		"spectr",
		"changes",
	)
	changeDir := filepath.Join(
		changesDir,
		"test-change",
	)
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Write proposal.md without H1 heading
	proposalContent := `Some content without heading`
	if err := os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte(proposalContent), 0o644); err != nil {
		t.Fatal(err)
	}

	lister := NewLister(tmpDir)
	changes, err := lister.ListChanges()
	if err != nil {
		t.Fatalf("ListChanges failed: %v", err)
	}

	if len(changes) != 1 {
		t.Fatalf(
			"Expected 1 change, got %d",
			len(changes),
		)
	}

	// Should fall back to ID as title
	if changes[0].Title != "test-change" {
		t.Errorf(
			"Expected fallback title 'test-change', got %q",
			changes[0].Title,
		)
	}
}

func TestListSpecs(t *testing.T) {
	tmpDir := t.TempDir()
	specsDir := filepath.Join(
		tmpDir,
		"spectr",
		"specs",
	)
	specDir := filepath.Join(
		specsDir,
		"authentication",
	)
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Write spec.md
	specContent := `# Authentication

### Requirement: User Login
Login feature

### Requirement: Password Reset
Reset feature

### Requirement: Two-Factor Auth
2FA feature`
	if err := os.WriteFile(filepath.Join(specDir, "spec.md"), []byte(specContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Test listing
	lister := NewLister(tmpDir)
	specs, err := lister.ListSpecs()
	if err != nil {
		t.Fatalf("ListSpecs failed: %v", err)
	}

	if len(specs) != 1 {
		t.Fatalf(
			"Expected 1 spec, got %d",
			len(specs),
		)
	}

	spec := specs[0]
	if spec.ID != "authentication" {
		t.Errorf(
			"Expected ID 'authentication', got %q",
			spec.ID,
		)
	}
	if spec.Title != "Authentication" {
		t.Errorf(
			"Expected title 'Authentication', got %q",
			spec.Title,
		)
	}
	if spec.RequirementCount != 3 {
		t.Errorf(
			"Expected 3 requirements, got %d",
			spec.RequirementCount,
		)
	}
}

func TestListSpecs_NoSpecs(t *testing.T) {
	tmpDir := t.TempDir()
	lister := NewLister(tmpDir)
	specs, err := lister.ListSpecs()
	if err != nil {
		t.Fatalf("ListSpecs failed: %v", err)
	}
	if len(specs) != 0 {
		t.Errorf(
			"Expected empty list, got %d specs",
			len(specs),
		)
	}
}

func TestListSpecs_FallbackTitle(t *testing.T) {
	tmpDir := t.TempDir()
	specsDir := filepath.Join(
		tmpDir,
		"spectr",
		"specs",
	)
	specDir := filepath.Join(
		specsDir,
		"test-spec",
	)
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Write spec.md without H1 heading
	specContent := `Some content without heading

### Requirement: Feature`
	if err := os.WriteFile(filepath.Join(specDir, "spec.md"), []byte(specContent), 0o644); err != nil {
		t.Fatal(err)
	}

	lister := NewLister(tmpDir)
	specs, err := lister.ListSpecs()
	if err != nil {
		t.Fatalf("ListSpecs failed: %v", err)
	}

	if len(specs) != 1 {
		t.Fatalf(
			"Expected 1 spec, got %d",
			len(specs),
		)
	}

	// Should fall back to ID as title
	if specs[0].Title != "test-spec" {
		t.Errorf(
			"Expected fallback title 'test-spec', got %q",
			specs[0].Title,
		)
	}
}

func TestListAll(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a change
	changesDir := filepath.Join(
		tmpDir,
		"spectr",
		"changes",
	)
	changeDir := filepath.Join(
		changesDir,
		"add-feature",
	)
	if err := os.MkdirAll(filepath.Join(changeDir, "specs", "test-spec"), 0o755); err != nil {
		t.Fatal(err)
	}

	proposalContent := `# Change: Add Feature`
	if err := os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte(proposalContent), 0o644); err != nil {
		t.Fatal(err)
	}

	tasksContent := `## Tasks
- [x] Task 1`
	if err := os.WriteFile(filepath.Join(changeDir, "tasks.md"), []byte(tasksContent), 0o644); err != nil {
		t.Fatal(err)
	}

	specDeltaContent := `## ADDED Requirements
### Requirement: New Feature`
	if err := os.WriteFile(filepath.Join(changeDir, "specs", "test-spec", "spec.md"), []byte(specDeltaContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create a spec
	specsDir := filepath.Join(
		tmpDir,
		"spectr",
		"specs",
	)
	specDir := filepath.Join(
		specsDir,
		"authentication",
	)
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatal(err)
	}

	specContent := `# Authentication

### Requirement: User Login`
	if err := os.WriteFile(filepath.Join(specDir, "spec.md"), []byte(specContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Test listing all items
	lister := NewLister(tmpDir)
	items, err := lister.ListAll(nil)
	if err != nil {
		t.Fatalf("ListAll failed: %v", err)
	}

	if len(items) != 2 {
		t.Fatalf(
			"Expected 2 items, got %d",
			len(items),
		)
	}

	// Verify sorting by ID (add-feature comes before authentication)
	if items[0].ID() != "add-feature" {
		t.Errorf(
			"Expected first item to be 'add-feature', got %q",
			items[0].ID(),
		)
	}
	if items[0].Type != ItemTypeChange {
		t.Errorf(
			"Expected first item to be a change, got %v",
			items[0].Type,
		)
	}

	if items[1].ID() != "authentication" {
		t.Errorf(
			"Expected second item to be 'authentication', got %q",
			items[1].ID(),
		)
	}
	if items[1].Type != ItemTypeSpec {
		t.Errorf(
			"Expected second item to be a spec, got %v",
			items[1].Type,
		)
	}
}

func TestListAll_FilterByType(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a change
	changesDir := filepath.Join(
		tmpDir,
		"spectr",
		"changes",
	)
	changeDir := filepath.Join(
		changesDir,
		"add-feature",
	)
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	proposalContent := `# Change: Add Feature`
	if err := os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte(proposalContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create a spec
	specsDir := filepath.Join(
		tmpDir,
		"spectr",
		"specs",
	)
	specDir := filepath.Join(
		specsDir,
		"authentication",
	)
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatal(err)
	}

	specContent := `# Authentication

### Requirement: User Login`
	if err := os.WriteFile(filepath.Join(specDir, "spec.md"), []byte(specContent), 0o644); err != nil {
		t.Fatal(err)
	}

	lister := NewLister(tmpDir)

	// Test filtering for changes only
	changeType := ItemTypeChange
	items, err := lister.ListAll(&ListAllOptions{
		FilterType: &changeType,
		SortByID:   true,
	})
	if err != nil {
		t.Fatalf(
			"ListAll with change filter failed: %v",
			err,
		)
	}

	if len(items) != 1 {
		t.Fatalf(
			"Expected 1 item (change), got %d",
			len(items),
		)
	}

	if items[0].Type != ItemTypeChange {
		t.Errorf(
			"Expected change item, got %v",
			items[0].Type,
		)
	}

	// Test filtering for specs only
	specType := ItemTypeSpec
	items, err = lister.ListAll(&ListAllOptions{
		FilterType: &specType,
		SortByID:   true,
	})
	if err != nil {
		t.Fatalf(
			"ListAll with spec filter failed: %v",
			err,
		)
	}

	if len(items) != 1 {
		t.Fatalf(
			"Expected 1 item (spec), got %d",
			len(items),
		)
	}

	if items[0].Type != ItemTypeSpec {
		t.Errorf(
			"Expected spec item, got %v",
			items[0].Type,
		)
	}
}

func TestListAll_NoSorting(t *testing.T) {
	tmpDir := t.TempDir()

	// Create multiple changes with IDs that sort differently
	changesDir := filepath.Join(
		tmpDir,
		"spectr",
		"changes",
	)
	for _, id := range []string{"zebra-change", "alpha-change"} {
		changeDir := filepath.Join(changesDir, id)
		if err := os.MkdirAll(changeDir, 0o755); err != nil {
			t.Fatal(err)
		}

		content := `# Change: Test`
		if err := os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	lister := NewLister(tmpDir)

	// Test with sorting disabled
	items, err := lister.ListAll(&ListAllOptions{
		SortByID: false,
	})
	if err != nil {
		t.Fatalf(
			"ListAll without sorting failed: %v",
			err,
		)
	}

	if len(items) != 2 {
		t.Fatalf(
			"Expected 2 items, got %d",
			len(items),
		)
	}

	// Without sorting, order depends on filesystem readdir order
	// We just verify all items are present
	ids := make(map[string]bool)
	for _, item := range items {
		ids[item.ID()] = true
	}

	if !ids["zebra-change"] ||
		!ids["alpha-change"] {
		t.Error(
			"Expected both zebra-change and alpha-change to be present",
		)
	}
}

func TestListAll_Empty(t *testing.T) {
	tmpDir := t.TempDir()
	lister := NewLister(tmpDir)

	items, err := lister.ListAll(nil)
	if err != nil {
		t.Fatalf(
			"ListAll on empty directory failed: %v",
			err,
		)
	}

	if len(items) != 0 {
		t.Errorf(
			"Expected empty list, got %d items",
			len(items),
		)
	}
}

// isGitAvailable checks if git command is available.
func isGitAvailable() bool {
	_, err := exec.LookPath("git")

	return err == nil
}

// isInGitRepo checks if we're currently in a git repository.
func isInGitRepo() bool {
	cmd := exec.Command(
		"git",
		"rev-parse",
		"--git-dir",
	)

	return cmd.Run() == nil
}

// hasOriginRemote checks if the origin remote is configured.
func hasOriginRemote() bool {
	cmd := exec.Command(
		"git",
		"remote",
		"get-url",
		"origin",
	)

	return cmd.Run() == nil
}

func TestFilterChangesNotOnRef_EmptySlice(
	t *testing.T,
) {
	// Test with empty changes slice - should return empty slice without error
	var changes []ChangeInfo
	result, err := FilterChangesNotOnRef(
		changes,
		"origin/main",
	)
	if err != nil {
		t.Fatalf(
			"FilterChangesNotOnRef with empty slice failed: %v",
			err,
		)
	}
	if len(result) != 0 {
		t.Errorf(
			"Expected empty result, got %d changes",
			len(result),
		)
	}
}

func TestFilterChangesNotOnRef_Integration(
	t *testing.T,
) {
	if !isGitAvailable() {
		t.Skip("git is not available")
	}
	if !isInGitRepo() {
		t.Skip("not in a git repository")
	}
	if !hasOriginRemote() {
		t.Skip("origin remote not configured")
	}

	// Change to the repo root for consistent path resolution
	cmd := exec.Command(
		"git",
		"rev-parse",
		"--show-toplevel",
	)
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf(
			"failed to get repo root: %v",
			err,
		)
	}
	repoRoot := string(
		output[:len(output)-1],
	) // trim newline

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf(
			"failed to get working directory: %v",
			err,
		)
	}
	defer func() { _ = os.Chdir(oldWd) }()
	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf(
			"failed to change to repo root: %v",
			err,
		)
	}

	// Create test changes - some that exist on origin/main and some that don't
	testCases := []struct {
		name             string
		changes          []ChangeInfo
		ref              string
		expectFiltered   int // number of changes expected in result
		expectContains   []string
		expectNotContain []string
	}{
		{
			name: "filter changes not on origin/main",
			changes: []ChangeInfo{
				// This is a made-up change that should NOT exist on origin/main
				{
					ID:         "nonexistent-test-change-xyz-12345",
					Title:      "Test Change",
					DeltaCount: 0,
					TaskStatus: parsers.TaskStatus{},
				},
				// These are known paths that should exist on origin/main (from archive)
				// We use "archive" since it's in the archive directory
			},
			ref:            "origin/main",
			expectFiltered: 1,
			expectContains: []string{
				"nonexistent-test-change-xyz-12345",
			},
		},
		{
			name:           "all changes already merged",
			changes:        nil, // Empty slice - no changes to filter
			ref:            "origin/main",
			expectFiltered: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := FilterChangesNotOnRef(
				tc.changes,
				tc.ref,
			)
			if err != nil {
				t.Fatalf(
					"FilterChangesNotOnRef failed: %v",
					err,
				)
			}

			if len(result) != tc.expectFiltered {
				t.Errorf(
					"Expected %d filtered changes, got %d",
					tc.expectFiltered,
					len(result),
				)
			}

			// Check expected IDs are in the result
			for _, expectedID := range tc.expectContains {
				found := false
				for _, change := range result {
					if change.ID == expectedID {
						found = true

						break
					}
				}
				if !found {
					t.Errorf(
						"Expected change %q to be in result, but it wasn't",
						expectedID,
					)
				}
			}
		})
	}
}

func TestFilterChangesNotOnRef_MixedChanges(
	t *testing.T,
) {
	if !isGitAvailable() {
		t.Skip("git is not available")
	}

	tmpDir := t.TempDir()

	repoRoot := filepath.Join(tmpDir, "repo")
	if err := os.MkdirAll(repoRoot, 0o755); err != nil {
		t.Fatalf("failed to create repo root: %v", err)
	}

	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf("failed to change to repo root: %v", err)
	}
	defer func() { _ = os.Chdir("..") }()

	initGitRepo := func() {
		cmd := exec.Command("git", "init")
		cmd.Dir = repoRoot
		if err := cmd.Run(); err != nil {
			t.Fatalf("failed to init git repo: %v", err)
		}

		cmd = exec.Command("git", "config", "user.email", "test@example.com")
		cmd.Dir = repoRoot
		if err := cmd.Run(); err != nil {
			t.Fatalf("failed to set git user email: %v", err)
		}

		cmd = exec.Command("git", "config", "user.name", "Test User")
		cmd.Dir = repoRoot
		if err := cmd.Run(); err != nil {
			t.Fatalf("failed to set git user name: %v", err)
		}
	}

	createChange := func(changeID string) {
		changeDir := filepath.Join(repoRoot, "spectr", "changes", changeID)
		if err := os.MkdirAll(changeDir, 0o755); err != nil {
			t.Fatalf("failed to create change dir %s: %v", changeID, err)
		}

		proposalContent := "# Change: " + changeID
		if err := os.WriteFile(filepath.Join(changeDir, "proposal.md"), []byte(proposalContent), 0o644); err != nil {
			t.Fatalf("failed to write proposal.md for %s: %v", changeID, err)
		}
	}

	commitChange := func(changeID string) {
		cmd := exec.Command("git", "add", ".")
		cmd.Dir = repoRoot
		if err := cmd.Run(); err != nil {
			t.Fatalf("failed to git add: %v", err)
		}

		cmd = exec.Command("git", "commit", "-m", "Add "+changeID)
		cmd.Dir = repoRoot
		if err := cmd.Run(); err != nil {
			t.Fatalf("failed to git commit: %v", err)
		}
	}

	initGitRepo()

	createChange("existing-change-on-main")
	commitChange("existing-change-on-main")

	changes := []ChangeInfo{
		{
			ID:         "existing-change-on-main",
			Title:      "Existing Change",
			DeltaCount: 1,
			TaskStatus: parsers.TaskStatus{},
		},
		{
			ID:         "nonexistent-change-abc",
			Title:      "Non-existent 1",
			DeltaCount: 1,
			TaskStatus: parsers.TaskStatus{
				Total:     2,
				Completed: 1,
			},
		},
		{
			ID:         "another-fake-change-xyz",
			Title:      "Non-existent 2",
			DeltaCount: 2,
			TaskStatus: parsers.TaskStatus{
				Total:     3,
				Completed: 0,
			},
		},
	}

	result, err := FilterChangesNotOnRef(
		changes,
		"HEAD",
	)
	if err != nil {
		t.Fatalf("FilterChangesNotOnRef failed: %v", err)
	}

	if len(result) != 2 {
		t.Errorf(
			"Expected 2 filtered changes (non-existent ones), got %d",
			len(result),
		)
		for _, c := range result {
			t.Logf("  - %s", c.ID)
		}
	}

	for _, change := range result {
		if change.ID == "existing-change-on-main" {
			t.Error(
				"existing-change-on-main should have been filtered out (it exists on HEAD)",
			)
		}
	}

	foundAbc := false
	foundXyz := false
	for _, change := range result {
		if change.ID == "nonexistent-change-abc" {
			foundAbc = true
		}
		if change.ID == "another-fake-change-xyz" {
			foundXyz = true
		}
	}
	if !foundAbc {
		t.Error(
			"nonexistent-change-abc should be in the result",
		)
	}
	if !foundXyz {
		t.Error(
			"another-fake-change-xyz should be in the result",
		)
	}
}

func TestFilterChangesNotOnRef_PreservesChangeInfo(
	t *testing.T,
) {
	if !isGitAvailable() {
		t.Skip("git is not available")
	}

	tmpDir := t.TempDir()

	repoRoot := filepath.Join(tmpDir, "repo")
	if err := os.MkdirAll(repoRoot, 0o755); err != nil {
		t.Fatalf("failed to create repo root: %v", err)
	}

	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf("failed to change to repo root: %v", err)
	}
	defer func() { _ = os.Chdir("..") }()

	initGitRepo := func() {
		cmd := exec.Command("git", "init")
		cmd.Dir = repoRoot
		if err := cmd.Run(); err != nil {
			t.Fatalf("failed to init git repo: %v", err)
		}

		cmd = exec.Command("git", "config", "user.email", "test@example.com")
		cmd.Dir = repoRoot
		if err := cmd.Run(); err != nil {
			t.Fatalf("failed to set git user email: %v", err)
		}

		cmd = exec.Command("git", "config", "user.name", "Test User")
		cmd.Dir = repoRoot
		if err := cmd.Run(); err != nil {
			t.Fatalf("failed to set git user name: %v", err)
		}
	}

	initGitRepo()

	cmd := exec.Command("git", "commit", "--allow-empty", "-m", "Initial commit")
	cmd.Dir = repoRoot
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to create initial commit: %v", err)
	}

	originalChange := ChangeInfo{
		ID:         "test-change-preserve-info-xyz",
		Title:      "Test Change with Full Info",
		DeltaCount: 5,
		TaskStatus: parsers.TaskStatus{
			Total:     10,
			Completed: 7,
		},
	}

	changes := []ChangeInfo{originalChange}

	result, err := FilterChangesNotOnRef(
		changes,
		"HEAD",
	)
	if err != nil {
		t.Fatalf(
			"FilterChangesNotOnRef failed: %v",
			err,
		)
	}

	if len(result) != 1 {
		t.Fatalf(
			"Expected 1 filtered change, got %d",
			len(result),
		)
	}

	filteredChange := result[0]
	if filteredChange.ID != originalChange.ID {
		t.Errorf(
			"ID mismatch: got %q, want %q",
			filteredChange.ID,
			originalChange.ID,
		)
	}
	if filteredChange.Title != originalChange.Title {
		t.Errorf(
			"Title mismatch: got %q, want %q",
			filteredChange.Title,
			originalChange.Title,
		)
	}
	if filteredChange.DeltaCount != originalChange.DeltaCount {
		t.Errorf(
			"DeltaCount mismatch: got %d, want %d",
			filteredChange.DeltaCount,
			originalChange.DeltaCount,
		)
	}
	if filteredChange.TaskStatus.Total != originalChange.TaskStatus.Total {
		t.Errorf(
			"TaskStatus.Total mismatch: got %d, want %d",
			filteredChange.TaskStatus.Total,
			originalChange.TaskStatus.Total,
		)
	}
	if filteredChange.TaskStatus.Completed != originalChange.TaskStatus.Completed {
		t.Errorf(
			"TaskStatus.Completed mismatch: got %d, want %d",
			filteredChange.TaskStatus.Completed,
			originalChange.TaskStatus.Completed,
		)
	}
}
