package discovery

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindSpectrRoots_SingleRoot(t *testing.T) {
	// Create a temp directory structure with a single spectr/ dir
	tmpDir := t.TempDir()

	// Create spectr/ directory
	spectrDir := filepath.Join(tmpDir, "spectr")
	if err := os.MkdirAll(spectrDir, 0o755); err != nil {
		t.Fatalf("failed to create spectr dir: %v", err)
	}

	// Create .git directory to establish boundary
	gitDir := filepath.Join(tmpDir, ".git")
	if err := os.MkdirAll(gitDir, 0o755); err != nil {
		t.Fatalf("failed to create .git dir: %v", err)
	}

	roots, err := FindSpectrRoots(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(roots) != 1 {
		t.Fatalf("expected 1 root, got %d", len(roots))
	}

	if roots[0].Path != tmpDir {
		t.Errorf("expected Path %s, got %s", tmpDir, roots[0].Path)
	}

	if roots[0].RelativeTo != "." {
		t.Errorf("expected RelativeTo '.', got %s", roots[0].RelativeTo)
	}

	if roots[0].GitRoot != tmpDir {
		t.Errorf("expected GitRoot %s, got %s", tmpDir, roots[0].GitRoot)
	}
}

func TestFindSpectrRoots_MultipleRoots(t *testing.T) {
	// Create a nested directory structure:
	// tmpDir/
	//   .git/
	//   spectr/          <- root 2 (parent)
	//   project/
	//     spectr/        <- root 1 (child, closer to cwd)
	//     src/           <- cwd
	tmpDir := t.TempDir()

	// Create .git directory at root
	gitDir := filepath.Join(tmpDir, ".git")
	if err := os.MkdirAll(gitDir, 0o755); err != nil {
		t.Fatalf("failed to create .git dir: %v", err)
	}

	// Create parent spectr/
	parentSpectr := filepath.Join(tmpDir, "spectr")
	if err := os.MkdirAll(parentSpectr, 0o755); err != nil {
		t.Fatalf("failed to create parent spectr dir: %v", err)
	}

	// Create project/spectr/
	projectSpectr := filepath.Join(tmpDir, "project", "spectr")
	if err := os.MkdirAll(projectSpectr, 0o755); err != nil {
		t.Fatalf("failed to create project spectr dir: %v", err)
	}

	// Create src/ as cwd
	srcDir := filepath.Join(tmpDir, "project", "src")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatalf("failed to create src dir: %v", err)
	}

	roots, err := FindSpectrRoots(srcDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(roots) != 2 {
		t.Fatalf("expected 2 roots, got %d", len(roots))
	}

	// First root should be the closest one (project/)
	projectDir := filepath.Join(tmpDir, "project")
	if roots[0].Path != projectDir {
		t.Errorf("expected first root Path %s, got %s", projectDir, roots[0].Path)
	}
	if roots[0].RelativeTo != ".." {
		t.Errorf("expected first root RelativeTo '..', got %s", roots[0].RelativeTo)
	}

	// Second root should be the parent (tmpDir)
	if roots[1].Path != tmpDir {
		t.Errorf("expected second root Path %s, got %s", tmpDir, roots[1].Path)
	}
}

func TestFindSpectrRoots_GitBoundary(t *testing.T) {
	// Create a structure where spectr/ exists outside the git boundary
	// tmpDir/
	//   spectr/              <- should NOT be found (outside git boundary)
	//   repo/
	//     .git/
	//     spectr/            <- should be found
	tmpDir := t.TempDir()

	// Create outer spectr/ (outside git boundary)
	outerSpectr := filepath.Join(tmpDir, "spectr")
	if err := os.MkdirAll(outerSpectr, 0o755); err != nil {
		t.Fatalf("failed to create outer spectr dir: %v", err)
	}

	// Create repo/.git/
	repoGit := filepath.Join(tmpDir, "repo", ".git")
	if err := os.MkdirAll(repoGit, 0o755); err != nil {
		t.Fatalf("failed to create repo .git dir: %v", err)
	}

	// Create repo/spectr/
	repoSpectr := filepath.Join(tmpDir, "repo", "spectr")
	if err := os.MkdirAll(repoSpectr, 0o755); err != nil {
		t.Fatalf("failed to create repo spectr dir: %v", err)
	}

	repoDir := filepath.Join(tmpDir, "repo")
	roots, err := FindSpectrRoots(repoDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(roots) != 1 {
		t.Fatalf("expected 1 root (git boundary should stop discovery), got %d", len(roots))
	}

	if roots[0].Path != repoDir {
		t.Errorf("expected Path %s, got %s", repoDir, roots[0].Path)
	}
}

func TestFindSpectrRoots_EnvVarOverride(t *testing.T) {
	tmpDir := t.TempDir()

	// Create spectr/ directory
	spectrDir := filepath.Join(tmpDir, "spectr")
	if err := os.MkdirAll(spectrDir, 0o755); err != nil {
		t.Fatalf("failed to create spectr dir: %v", err)
	}

	// Create another directory to use as cwd
	otherDir := filepath.Join(tmpDir, "other")
	if err := os.MkdirAll(otherDir, 0o755); err != nil {
		t.Fatalf("failed to create other dir: %v", err)
	}

	// Set SPECTR_ROOT env var
	t.Setenv("SPECTR_ROOT", tmpDir)

	roots, err := FindSpectrRoots(otherDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(roots) != 1 {
		t.Fatalf("expected 1 root from env var, got %d", len(roots))
	}

	if roots[0].Path != tmpDir {
		t.Errorf("expected Path %s, got %s", tmpDir, roots[0].Path)
	}
}

func TestFindSpectrRoots_EnvVarInvalidPath(t *testing.T) {
	tmpDir := t.TempDir()

	// Create directory without spectr/
	noSpectrDir := filepath.Join(tmpDir, "nope")
	if err := os.MkdirAll(noSpectrDir, 0o755); err != nil {
		t.Fatalf("failed to create dir: %v", err)
	}

	// Set SPECTR_ROOT to invalid path
	t.Setenv("SPECTR_ROOT", noSpectrDir)

	_, err := FindSpectrRoots(tmpDir)
	if err == nil {
		t.Fatal("expected error for invalid SPECTR_ROOT, got nil")
	}

	expectedErrMsg := "SPECTR_ROOT path does not contain spectr/ directory"
	if !contains(err.Error(), expectedErrMsg) {
		t.Errorf("expected error containing %q, got %q", expectedErrMsg, err.Error())
	}
}

func TestFindSpectrRoots_NoRootsFound(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .git to establish boundary (no spectr/)
	gitDir := filepath.Join(tmpDir, ".git")
	if err := os.MkdirAll(gitDir, 0o755); err != nil {
		t.Fatalf("failed to create .git dir: %v", err)
	}

	roots, err := FindSpectrRoots(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(roots) != 0 {
		t.Errorf("expected 0 roots, got %d", len(roots))
	}
}

func TestFindSpectrRoots_FromSubdirectory(t *testing.T) {
	// Test running from deep within a project
	// tmpDir/
	//   .git/
	//   spectr/
	//   src/
	//     pkg/
	//       internal/    <- cwd
	tmpDir := t.TempDir()

	// Create .git
	gitDir := filepath.Join(tmpDir, ".git")
	if err := os.MkdirAll(gitDir, 0o755); err != nil {
		t.Fatalf("failed to create .git dir: %v", err)
	}

	// Create spectr/
	spectrDir := filepath.Join(tmpDir, "spectr")
	if err := os.MkdirAll(spectrDir, 0o755); err != nil {
		t.Fatalf("failed to create spectr dir: %v", err)
	}

	// Create nested cwd
	deepDir := filepath.Join(tmpDir, "src", "pkg", "internal")
	if err := os.MkdirAll(deepDir, 0o755); err != nil {
		t.Fatalf("failed to create deep dir: %v", err)
	}

	roots, err := FindSpectrRoots(deepDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(roots) != 1 {
		t.Fatalf("expected 1 root, got %d", len(roots))
	}

	if roots[0].Path != tmpDir {
		t.Errorf("expected Path %s, got %s", tmpDir, roots[0].Path)
	}

	if roots[0].RelativeTo != "../../.." {
		t.Errorf("expected RelativeTo '../../..', got %s", roots[0].RelativeTo)
	}
}

func TestSpectrRoot_SpectrDir(t *testing.T) {
	root := SpectrRoot{Path: "/home/user/project"}
	expected := "/home/user/project/spectr"
	if got := root.SpectrDir(); got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func TestSpectrRoot_ChangesDir(t *testing.T) {
	root := SpectrRoot{Path: "/home/user/project"}
	expected := "/home/user/project/spectr/changes"
	if got := root.ChangesDir(); got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func TestSpectrRoot_SpecsDir(t *testing.T) {
	root := SpectrRoot{Path: "/home/user/project"}
	expected := "/home/user/project/spectr/specs"
	if got := root.SpecsDir(); got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func TestSpectrRoot_DisplayName(t *testing.T) {
	tests := []struct {
		name     string
		root     SpectrRoot
		expected string
	}{
		{
			name: "current directory",
			root: SpectrRoot{
				Path:       "/home/user/project",
				RelativeTo: ".",
			},
			expected: "project",
		},
		{
			name: "relative path",
			root: SpectrRoot{
				Path:       "/home/user/mono/project",
				RelativeTo: "../project",
			},
			expected: "../project",
		},
		{
			name: "parent relative",
			root: SpectrRoot{
				Path:       "/home/user",
				RelativeTo: "../..",
			},
			expected: "../..",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.root.DisplayName(); got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}

// setupNestedGitRepoFixture creates a mono-repo structure with nested git repositories
// Returns the root temp directory path.
func setupNestedGitRepoFixture(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()

	// Create main repo .git and spectr/
	mustMkdirAll(t, filepath.Join(tmpDir, ".git"))
	mustMkdirAll(t, filepath.Join(tmpDir, "spectr"))

	// Create auth package with its own .git and spectr
	mustMkdirAll(t, filepath.Join(tmpDir, "packages", "auth", ".git"))
	mustMkdirAll(t, filepath.Join(tmpDir, "packages", "auth", "spectr"))

	// Create api package with its own .git and spectr
	mustMkdirAll(t, filepath.Join(tmpDir, "packages", "api", ".git"))
	mustMkdirAll(t, filepath.Join(tmpDir, "packages", "api", "spectr"))

	return tmpDir
}

// mustMkdirAll creates a directory, failing the test if it cannot.
func mustMkdirAll(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("failed to create directory %s: %v", path, err)
	}
}

// assertSingleRoot verifies exactly one root is found with the expected path.
func assertSingleRoot(t *testing.T, roots []SpectrRoot, expectedPath string) {
	t.Helper()

	if len(roots) != 1 {
		t.Errorf("expected 1 root, got %d", len(roots))
		for i, r := range roots {
			t.Logf("  root[%d]: %s", i, r.Path)
		}

		return
	}

	if roots[0].Path != expectedPath {
		t.Errorf("expected root to be %s, got %s", expectedPath, roots[0].Path)
	}
}

func TestFindSpectrRoots_NestedGitRepos(t *testing.T) {
	// Integration test: mono-repo with nested git repositories
	// Structure:
	// mono-repo/
	//   .git/                    <- main repo boundary
	//   spectr/                  <- main repo spectr
	//   packages/
	//     auth/
	//       .git/               <- nested git repo
	//       spectr/             <- auth's spectr
	//     api/
	//       .git/               <- nested git repo
	//       spectr/             <- api's spectr
	tmpDir := setupNestedGitRepoFixture(t)

	t.Run("from mono-repo root finds only main spectr", func(t *testing.T) {
		roots, err := FindSpectrRoots(tmpDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		assertSingleRoot(t, roots, tmpDir)
	})

	t.Run("from auth package finds only auth spectr", func(t *testing.T) {
		authDir := filepath.Join(tmpDir, "packages", "auth")
		roots, err := FindSpectrRoots(authDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		assertSingleRoot(t, roots, authDir)
	})

	t.Run("from api subdirectory finds only api spectr", func(t *testing.T) {
		apiSrc := filepath.Join(tmpDir, "packages", "api", "src")
		mustMkdirAll(t, apiSrc)

		roots, err := FindSpectrRoots(apiSrc)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		apiDir := filepath.Join(tmpDir, "packages", "api")
		assertSingleRoot(t, roots, apiDir)
	})

	t.Run("from packages directory finds main spectr", func(t *testing.T) {
		packagesDir := filepath.Join(tmpDir, "packages")

		roots, err := FindSpectrRoots(packagesDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// packages/ is within the main repo's git boundary
		assertSingleRoot(t, roots, tmpDir)
	})
}

func TestFindSpectrRoots_RelativePathCalculation(t *testing.T) {
	// Test that relative paths are correctly calculated for multi-root scenarios
	tmpDir := t.TempDir()

	// Create nested structure within same git repo:
	// tmpDir/
	//   .git/
	//   spectr/
	//   services/
	//     backend/
	//       spectr/
	//       src/
	//         handlers/   <- cwd

	gitDir := filepath.Join(tmpDir, ".git")
	if err := os.MkdirAll(gitDir, 0o755); err != nil {
		t.Fatalf("failed to create .git dir: %v", err)
	}

	rootSpectr := filepath.Join(tmpDir, "spectr")
	if err := os.MkdirAll(rootSpectr, 0o755); err != nil {
		t.Fatalf("failed to create root spectr dir: %v", err)
	}

	backendSpectr := filepath.Join(tmpDir, "services", "backend", "spectr")
	if err := os.MkdirAll(backendSpectr, 0o755); err != nil {
		t.Fatalf("failed to create backend spectr dir: %v", err)
	}

	handlersDir := filepath.Join(tmpDir, "services", "backend", "src", "handlers")
	if err := os.MkdirAll(handlersDir, 0o755); err != nil {
		t.Fatalf("failed to create handlers dir: %v", err)
	}

	roots, err := FindSpectrRoots(handlersDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(roots) != 2 {
		t.Fatalf("expected 2 roots, got %d", len(roots))
	}

	// First root should be backend (closest)
	backendDir := filepath.Join(tmpDir, "services", "backend")
	if roots[0].Path != backendDir {
		t.Errorf("expected first root Path %s, got %s", backendDir, roots[0].Path)
	}
	// nolint:goconst // Test-specific relative path, not worth extracting as constant
	if roots[0].RelativeTo != "../.." {
		t.Errorf("expected first root RelativeTo '../..', got %s", roots[0].RelativeTo)
	}

	// Second root should be the root tmpDir
	if roots[1].Path != tmpDir {
		t.Errorf("expected second root Path %s, got %s", tmpDir, roots[1].Path)
	}
	if roots[1].RelativeTo != "../../../.." {
		t.Errorf("expected second root RelativeTo '../../../..', got %s", roots[1].RelativeTo)
	}
}

// TestFindSpectrRootsDownward tests the downward directory discovery.
// nolint:revive // Complex test function with multiple comprehensive scenarios
func TestFindSpectrRootsDownward(t *testing.T) {
	t.Run("finds nested spectr directories", func(t *testing.T) {
		// Create structure:
		// tmpDir/
		//   project1/
		//     .git/
		//     spectr/
		//   project2/
		//     .git/
		//     spectr/
		//   project3/
		//     spectr/  (no .git)
		tmpDir := t.TempDir()

		mustMkdirAll(t, filepath.Join(tmpDir, "project1", ".git"))
		mustMkdirAll(t, filepath.Join(tmpDir, "project1", "spectr"))

		mustMkdirAll(t, filepath.Join(tmpDir, "project2", ".git"))
		mustMkdirAll(t, filepath.Join(tmpDir, "project2", "spectr"))

		mustMkdirAll(t, filepath.Join(tmpDir, "project3", "spectr"))

		roots, err := findSpectrRootsDownward(tmpDir, tmpDir, maxDiscoveryDepth)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(roots) != 3 {
			t.Fatalf("expected 3 roots, got %d", len(roots))
		}

		// Verify all projects were found
		foundProjects := make(map[string]bool)
		for _, root := range roots {
			baseName := filepath.Base(root.Path)
			foundProjects[baseName] = true

			// Verify git root is set correctly for projects with .git
			if baseName != "project1" && baseName != "project2" {
				continue
			}

			if root.GitRoot != root.Path {
				t.Errorf("expected GitRoot %s for %s, got %s",
					root.Path, baseName, root.GitRoot)
			}
		}

		expectedProjects := []string{"project1", "project2", "project3"}
		for _, project := range expectedProjects {
			if !foundProjects[project] {
				t.Errorf("expected to find %s", project)
			}
		}
	})

	t.Run("respects max depth limit", func(t *testing.T) {
		// Create deep nesting:
		// tmpDir/
		//   a/b/c/d/e/f/g/h/i/j/k/
		//     spectr/  <- at depth 11, should not be found with maxDepth=10
		tmpDir := t.TempDir()
		deepPath := filepath.Join(tmpDir, "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k")
		mustMkdirAll(t, filepath.Join(deepPath, "spectr"))

		roots, err := findSpectrRootsDownward(tmpDir, tmpDir, 10)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(roots) != 0 {
			t.Errorf("expected 0 roots (too deep), got %d", len(roots))
		}
	})

	t.Run("skips ignored directories", func(t *testing.T) {
		// Create structure with directories that should be skipped:
		// tmpDir/
		//   .git/spectr/         <- should skip
		//   node_modules/spectr/ <- should skip
		//   vendor/spectr/       <- should skip
		//   project/spectr/      <- should find
		tmpDir := t.TempDir()

		mustMkdirAll(t, filepath.Join(tmpDir, ".git", "spectr"))
		mustMkdirAll(t, filepath.Join(tmpDir, "node_modules", "spectr"))
		mustMkdirAll(t, filepath.Join(tmpDir, "vendor", "spectr"))
		mustMkdirAll(t, filepath.Join(tmpDir, "project", "spectr"))

		roots, err := findSpectrRootsDownward(tmpDir, tmpDir, maxDiscoveryDepth)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(roots) != 1 {
			t.Fatalf("expected 1 root (only project), got %d", len(roots))
		}

		if filepath.Base(roots[0].Path) != "project" {
			t.Errorf("expected to find project, got %s", roots[0].Path)
		}
	})

	t.Run("finds nested repos in subdirectories", func(t *testing.T) {
		// Create mono-repo structure:
		// tmpDir/
		//   packages/
		//     auth/
		//       .git/
		//       spectr/
		//     api/
		//       .git/
		//       spectr/
		tmpDir := t.TempDir()

		mustMkdirAll(t, filepath.Join(tmpDir, "packages", "auth", ".git"))
		mustMkdirAll(t, filepath.Join(tmpDir, "packages", "auth", "spectr"))

		mustMkdirAll(t, filepath.Join(tmpDir, "packages", "api", ".git"))
		mustMkdirAll(t, filepath.Join(tmpDir, "packages", "api", "spectr"))

		roots, err := findSpectrRootsDownward(tmpDir, tmpDir, maxDiscoveryDepth)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(roots) != 2 {
			t.Fatalf("expected 2 roots, got %d", len(roots))
		}

		// Verify both auth and api were found
		foundRepos := make(map[string]bool)
		for _, root := range roots {
			baseName := filepath.Base(root.Path)
			foundRepos[baseName] = true

			// Both should have .git, so GitRoot should be set
			if root.GitRoot == "" {
				t.Errorf("expected GitRoot to be set for %s", baseName)
			}
		}

		if !foundRepos["auth"] || !foundRepos["api"] {
			t.Error("expected to find both auth and api repos")
		}
	})

	t.Run("continues after finding spectr", func(t *testing.T) {
		// Verify it doesn't stop at first spectr/ found
		// tmpDir/
		//   outer/
		//     spectr/
		//     inner/
		//       spectr/
		tmpDir := t.TempDir()

		mustMkdirAll(t, filepath.Join(tmpDir, "outer", "spectr"))
		mustMkdirAll(t, filepath.Join(tmpDir, "outer", "inner", "spectr"))

		roots, err := findSpectrRootsDownward(tmpDir, tmpDir, maxDiscoveryDepth)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(roots) != 2 {
			t.Fatalf("expected 2 roots, got %d", len(roots))
		}
	})

	t.Run("handles permission errors gracefully", func(t *testing.T) {
		// This test verifies the function continues when it encounters
		// directories it can't read (permission errors)
		tmpDir := t.TempDir()

		// Create a readable project
		mustMkdirAll(t, filepath.Join(tmpDir, "readable", "spectr"))

		// Note: We can't reliably test permission errors in all environments
		// (CI, different OSes, etc.), so we just verify the function doesn't
		// crash and finds the readable directory
		roots, err := findSpectrRootsDownward(tmpDir, tmpDir, maxDiscoveryDepth)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(roots) < 1 {
			t.Error("expected to find at least the readable directory")
		}
	})

	t.Run("calculates relative paths correctly", func(t *testing.T) {
		// Test from different cwd
		tmpDir := t.TempDir()

		mustMkdirAll(t, filepath.Join(tmpDir, "projects", "myapp", "spectr"))

		// Pretend cwd is somewhere else
		fakeCwd := filepath.Join(tmpDir, "somedir")
		mustMkdirAll(t, fakeCwd)

		roots, err := findSpectrRootsDownward(tmpDir, fakeCwd, maxDiscoveryDepth)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(roots) != 1 {
			t.Fatalf("expected 1 root, got %d", len(roots))
		}

		// RelativeTo should be relative to fakeCwd
		expectedRelPath := "../projects/myapp"
		if roots[0].RelativeTo != expectedRelPath {
			t.Errorf("expected RelativeTo %s, got %s",
				expectedRelPath, roots[0].RelativeTo)
		}
	})
}

// TestDeduplicateRoots tests the deduplicateRoots helper function.
func TestDeduplicateRoots(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		var roots []SpectrRoot
		result := deduplicateRoots(roots)

		if len(result) != 0 {
			t.Errorf("expected empty slice, got %d roots", len(result))
		}
	})

	t.Run("no duplicates", func(t *testing.T) {
		roots := []SpectrRoot{
			{Path: "/home/user/project1", RelativeTo: ".", GitRoot: "/home/user/project1"},
			{
				Path:       "/home/user/project2",
				RelativeTo: "../project2",
				GitRoot:    "/home/user/project2",
			},
			{
				Path:       "/home/user/project3",
				RelativeTo: "../project3",
				GitRoot:    "/home/user/project3",
			},
		}

		result := deduplicateRoots(roots)

		if len(result) != 3 {
			t.Errorf("expected 3 roots, got %d", len(result))
		}

		for i, root := range result {
			if root.Path != roots[i].Path {
				t.Errorf("order changed: expected %s at index %d, got %s",
					roots[i].Path, i, root.Path)
			}
		}
	})

	t.Run("with duplicates preserves first occurrence", func(t *testing.T) {
		roots := []SpectrRoot{
			{Path: "/home/user/project1", RelativeTo: ".", GitRoot: "/home/user/project1"},
			{
				Path:       "/home/user/project2",
				RelativeTo: "../project2",
				GitRoot:    "/home/user/project2",
			},
			{
				Path:       "/home/user/project1",
				RelativeTo: "different",
				GitRoot:    "different",
			}, // duplicate
			{
				Path:       "/home/user/project3",
				RelativeTo: "../project3",
				GitRoot:    "/home/user/project3",
			},
			{
				Path:       "/home/user/project2",
				RelativeTo: "also-different",
				GitRoot:    "also-different",
			}, // duplicate
		}

		result := deduplicateRoots(roots)

		if len(result) != 3 {
			t.Errorf("expected 3 unique roots, got %d", len(result))
		}

		// Verify first occurrences were preserved
		if result[0].Path != "/home/user/project1" || result[0].RelativeTo != "." {
			t.Error("first occurrence of project1 not preserved")
		}

		if result[1].Path != "/home/user/project2" || result[1].RelativeTo != "../project2" {
			t.Error("first occurrence of project2 not preserved")
		}

		if result[2].Path != "/home/user/project3" {
			t.Error("project3 not in result")
		}
	})

	t.Run("all duplicates", func(t *testing.T) {
		roots := []SpectrRoot{
			{Path: "/home/user/project", RelativeTo: ".", GitRoot: "/home/user/project"},
			{Path: "/home/user/project", RelativeTo: "different1", GitRoot: "git1"},
			{Path: "/home/user/project", RelativeTo: "different2", GitRoot: "git2"},
		}

		result := deduplicateRoots(roots)

		if len(result) != 1 {
			t.Errorf("expected 1 unique root, got %d", len(result))
		}

		if result[0].RelativeTo != "." {
			t.Error("first occurrence not preserved")
		}
	})
}

// TestSortRootsByDistance tests the sortRootsByDistance helper function.
func TestSortRootsByDistance(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		var roots []SpectrRoot
		cwd := "/home/user/project"

		result := sortRootsByDistance(roots, cwd)

		if len(result) != 0 {
			t.Errorf("expected empty slice, got %d roots", len(result))
		}
	})

	t.Run("single root", func(t *testing.T) {
		roots := []SpectrRoot{
			{Path: "/home/user/project", RelativeTo: ".", GitRoot: "/home/user/project"},
		}
		cwd := "/home/user/project"

		result := sortRootsByDistance(roots, cwd)

		if len(result) != 1 {
			t.Errorf("expected 1 root, got %d", len(result))
		}

		if result[0].Path != roots[0].Path {
			t.Error("single root was modified")
		}
	})

	t.Run("sorts by distance - closest first", func(t *testing.T) {
		// Setup: cwd is at /home/user/mono/project/src
		cwd := "/home/user/mono/project/src"

		roots := []SpectrRoot{
			{
				Path:       "/home/user/mono",
				RelativeTo: "../..",
				GitRoot:    "/home/user/mono",
			}, // distance 2 (../..)
			{
				Path:       "/home/user/mono/project",
				RelativeTo: "..",
				GitRoot:    "/home/user/mono",
			}, // distance 1 (..)
			{
				Path:       "/home/user/mono/project/src",
				RelativeTo: ".",
				GitRoot:    "/home/user/mono",
			}, // distance 0 (.)
			{
				Path:       "/home/user/mono/other",
				RelativeTo: "../../other",
				GitRoot:    "/home/user/mono",
			}, // distance 2
			{
				Path:       "/home/user/mono/project/lib",
				RelativeTo: "../lib",
				GitRoot:    "/home/user/mono",
			}, // distance 1
		}

		result := sortRootsByDistance(roots, cwd)

		if len(result) != 5 {
			t.Fatalf("expected 5 roots, got %d", len(result))
		}

		// Check order: closest (.) should be first
		if result[0].Path != "/home/user/mono/project/src" {
			t.Errorf("expected closest root first, got %s", result[0].Path)
		}

		// Next should be distance 1 roots (.. and ../lib)
		// They should be sorted alphabetically among themselves
		if result[1].RelativeTo != ".." && result[1].RelativeTo != "../lib" {
			t.Errorf("expected distance-1 root at index 1, got %s", result[1].RelativeTo)
		}

		// Last should be distance 2 roots (../.. and ../../other)
		if result[3].RelativeTo != "../.." && result[3].RelativeTo != "../../other" {
			t.Errorf("expected distance-2 root at index 3, got %s", result[3].RelativeTo)
		}
	})

	t.Run("preserves original slice - creates copy", func(t *testing.T) {
		cwd := "/home/user/mono/project"
		roots := []SpectrRoot{
			{Path: "/home/user/mono", RelativeTo: "..", GitRoot: "/home/user/mono"},
			{Path: "/home/user/mono/project", RelativeTo: ".", GitRoot: "/home/user/mono"},
		}

		originalFirstPath := roots[0].Path

		result := sortRootsByDistance(roots, cwd)

		// Original slice should be unchanged
		if roots[0].Path != originalFirstPath {
			t.Error("original slice was modified")
		}

		// Result should be sorted differently
		if result[0].Path != originalFirstPath {
			return // Different, as expected
		}

		// This might happen if they were already sorted, check second element
		if len(result) > 1 && result[1].Path == roots[1].Path {
			t.Error("result appears to be the same as input (should be sorted)")
		}
	})

	t.Run("handles equal distances alphabetically", func(t *testing.T) {
		cwd := "/home/user/mono"

		roots := []SpectrRoot{
			{Path: "/home/user/mono/zeta", RelativeTo: "zeta", GitRoot: "/home/user/mono"},
			{Path: "/home/user/mono/alpha", RelativeTo: "alpha", GitRoot: "/home/user/mono"},
			{Path: "/home/user/mono/beta", RelativeTo: "beta", GitRoot: "/home/user/mono"},
		}

		result := sortRootsByDistance(roots, cwd)

		if len(result) != 3 {
			t.Fatalf("expected 3 roots, got %d", len(result))
		}

		// All at same distance, should be alphabetically sorted
		if result[0].RelativeTo != "alpha" {
			t.Errorf("expected 'alpha' first, got %s", result[0].RelativeTo)
		}

		if result[1].RelativeTo != "beta" {
			t.Errorf("expected 'beta' second, got %s", result[1].RelativeTo)
		}

		if result[2].RelativeTo != "zeta" {
			t.Errorf("expected 'zeta' third, got %s", result[2].RelativeTo)
		}
	})
}

// TestFindSpectrRoots_SubdirectoryDiscovery tests basic subdirectory discovery.
// This verifies that FindSpectrRoots can find spectr/ directories in subdirectories
// when running from a directory without a git repository.
func TestFindSpectrRoots_SubdirectoryDiscovery(t *testing.T) {
	// Create structure (no git boundary, so downward discovery should happen):
	// tmpDir/
	//   project1/
	//     spectr/
	//   project2/
	//     spectr/
	//   deeply/
	//     nested/
	//       project3/
	//         spectr/
	tmpDir := t.TempDir()

	mustMkdirAll(t, filepath.Join(tmpDir, "project1", "spectr"))
	mustMkdirAll(t, filepath.Join(tmpDir, "project2", "spectr"))
	mustMkdirAll(t, filepath.Join(tmpDir, "deeply", "nested", "project3", "spectr"))

	// Run from tmpDir (no .git, so downward discovery occurs)
	roots, err := FindSpectrRoots(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(roots) != 3 {
		for i, r := range roots {
			t.Logf("  root[%d]: %s", i, r.Path)
		}
		t.Fatalf("expected 3 roots, got %d", len(roots))
	}

	// Verify all projects were found
	foundProjects := make(map[string]bool)
	for _, root := range roots {
		baseName := filepath.Base(root.Path)
		foundProjects[baseName] = true

		// Verify GitRoot is empty (no git boundaries)
		if root.GitRoot != "" {
			t.Errorf("expected empty GitRoot for %s (no .git), got %s",
				baseName, root.GitRoot)
		}
	}

	expectedProjects := []string{"project1", "project2", "project3"}
	for _, project := range expectedProjects {
		if !foundProjects[project] {
			t.Errorf("expected to find %s", project)
		}
	}
}

// TestFindSpectrRoots_MonorepoWithSubprojects tests the exact GitHub issue #363 scenario:
// a mono-repo with nested git repositories, each with their own spectr/ directory.
func TestFindSpectrRoots_MonorepoWithSubprojects(t *testing.T) {
	// Create the exact structure from GitHub issue #363:
	// mono-repo/
	//   .git/                    <- main repo boundary
	//   spectr/                  <- main repo spectr
	//   packages/
	//     auth/
	//       .git/                <- nested git repo
	//       spectr/              <- auth's spectr
	//       src/
	//         lib/               <- test from here
	//     api/
	//       .git/                <- nested git repo
	//       spectr/              <- api's spectr
	tmpDir := t.TempDir()

	// Create main mono-repo with .git and spectr/
	mustMkdirAll(t, filepath.Join(tmpDir, ".git"))
	mustMkdirAll(t, filepath.Join(tmpDir, "spectr"))

	// Create auth package with its own .git and spectr/
	authDir := filepath.Join(tmpDir, "packages", "auth")
	mustMkdirAll(t, filepath.Join(authDir, ".git"))
	mustMkdirAll(t, filepath.Join(authDir, "spectr"))
	authSrcLib := filepath.Join(authDir, "src", "lib")
	mustMkdirAll(t, authSrcLib)

	// Create api package with its own .git and spectr/
	apiDir := filepath.Join(tmpDir, "packages", "api")
	mustMkdirAll(t, filepath.Join(apiDir, ".git"))
	mustMkdirAll(t, filepath.Join(apiDir, "spectr"))

	t.Run("from mono-repo root finds only main spectr", func(t *testing.T) {
		roots, err := FindSpectrRoots(tmpDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		assertSingleRoot(t, roots, tmpDir)
		if roots[0].GitRoot != tmpDir {
			t.Errorf("expected GitRoot %s, got %s", tmpDir, roots[0].GitRoot)
		}
	})

	t.Run("from auth/src/lib finds only auth spectr (upward to git boundary)", func(t *testing.T) {
		roots, err := FindSpectrRoots(authSrcLib)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		assertSingleRoot(t, roots, authDir)
		if roots[0].GitRoot != authDir {
			t.Errorf("expected GitRoot %s, got %s", authDir, roots[0].GitRoot)
		}
		// Verify RelativeTo is correct
		expectedRel := "../.."
		if roots[0].RelativeTo != expectedRel {
			t.Errorf("expected RelativeTo %s, got %s", expectedRel, roots[0].RelativeTo)
		}
	})

	t.Run("from api package finds only api spectr", func(t *testing.T) {
		roots, err := FindSpectrRoots(apiDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		assertSingleRoot(t, roots, apiDir)
		if roots[0].GitRoot != apiDir {
			t.Errorf("expected GitRoot %s, got %s", apiDir, roots[0].GitRoot)
		}
	})

	t.Run("from packages directory finds only main spectr", func(t *testing.T) {
		packagesDir := filepath.Join(tmpDir, "packages")
		roots, err := FindSpectrRoots(packagesDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// packages/ is within the main repo's git boundary
		assertSingleRoot(t, roots, tmpDir)
		if roots[0].GitRoot != tmpDir {
			t.Errorf("expected GitRoot %s, got %s", tmpDir, roots[0].GitRoot)
		}
	})
}

// TestFindSpectrRoots_DepthLimit verifies that the 10-level depth limit is enforced
// during downward discovery, preventing excessive directory traversal.
func TestFindSpectrRoots_DepthLimit(t *testing.T) {
	// Note: This test directly tests findSpectrRootsDownward since FindSpectrRoots
	// only does downward discovery when NOT in a git repo. Testing the internal
	// function ensures the depth limit logic is verified.
	//
	// Due to how WalkDir and calculateDepth work together, the effective depths are:
	// - tmpDir (start path): depth 0 in depthMap initialization, but calculated as depth 1
	// - tmpDir/a: depth 2
	// - tmpDir/a/b/c/d/e/f/g/h/i: depth 10 (at limit)
	// - tmpDir/a/b/c/d/e/f/g/h/i/j: depth 11 (exceeds limit)
	tmpDir := t.TempDir()

	// Create shallow spectr (depth 2 with current implementation)
	mustMkdirAll(t, filepath.Join(tmpDir, "shallow", "spectr"))

	// Create at-limit spectr (depth 10)
	// Using 9 path segments: a/b/c/d/e/f/g/h/i
	atLimitPath := filepath.Join(tmpDir, "a", "b", "c", "d", "e", "f", "g", "h", "i")
	mustMkdirAll(t, filepath.Join(atLimitPath, "spectr"))

	// Create too-deep spectr (depth 11, should not be found)
	// Using 10 path segments: a/b/c/d/e/f/g/h/i/j
	tooDeepPath := filepath.Join(tmpDir, "a", "b", "c", "d", "e", "f", "g", "h", "i", "j")
	mustMkdirAll(t, filepath.Join(tooDeepPath, "spectr"))

	// Test the internal downward discovery function directly
	roots, err := findSpectrRootsDownward(tmpDir, tmpDir, maxDiscoveryDepth)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should find 2 roots: shallow and at-limit, but NOT too-deep
	if len(roots) != 2 {
		for i, r := range roots {
			relPath, _ := filepath.Rel(tmpDir, r.Path)
			t.Logf("  root[%d]: %s", i, relPath)
		}
		t.Fatalf("expected 2 roots (shallow and at-limit), got %d", len(roots))
	}

	// Verify shallow was found
	foundShallow := false
	foundAtLimit := false
	foundTooDeep := false

	for _, root := range roots {
		if filepath.Base(root.Path) == "shallow" {
			foundShallow = true
		}
		if root.Path == atLimitPath {
			foundAtLimit = true
		}
		if root.Path == tooDeepPath {
			foundTooDeep = true
		}
	}

	if !foundShallow {
		t.Error("expected to find shallow spectr/")
	}
	if !foundAtLimit {
		t.Error("expected to find spectr/ at depth limit")
	}
	if foundTooDeep {
		t.Error("should NOT find spectr/ beyond depth limit")
	}
}

// TestFindSpectrRoots_SkipsIgnoredDirs verifies that discovery skips common
// directories that should not contain spectr/ directories (.git, node_modules, etc).
func TestFindSpectrRoots_SkipsIgnoredDirs(t *testing.T) {
	// Note: This test directly tests findSpectrRootsDownward since FindSpectrRoots
	// only does downward discovery when NOT in a git repo. Testing the internal
	// function ensures the skip logic is verified.
	tmpDir := t.TempDir()

	// Create spectr/ inside ignored directories
	mustMkdirAll(t, filepath.Join(tmpDir, ".git", "spectr"))
	mustMkdirAll(t, filepath.Join(tmpDir, "node_modules", "package", "spectr"))
	mustMkdirAll(t, filepath.Join(tmpDir, "vendor", "lib", "spectr"))
	mustMkdirAll(t, filepath.Join(tmpDir, "target", "spectr"))
	mustMkdirAll(t, filepath.Join(tmpDir, "dist", "spectr"))
	mustMkdirAll(t, filepath.Join(tmpDir, "build", "spectr"))

	// Create valid spectr/ in regular project directory
	mustMkdirAll(t, filepath.Join(tmpDir, "project", "spectr"))

	// Test the internal downward discovery function directly
	roots, err := findSpectrRootsDownward(tmpDir, tmpDir, maxDiscoveryDepth)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should find only 1 root: project (all others should be skipped)
	if len(roots) != 1 {
		for i, r := range roots {
			t.Logf("  root[%d]: %s", i, r.Path)
		}
		t.Fatalf("expected 1 root (only project), got %d", len(roots))
	}

	if filepath.Base(roots[0].Path) != "project" {
		t.Errorf("expected to find only project, got %s", roots[0].Path)
	}
}

// TestFindSpectrRoots_Deduplication verifies that when both upward and downward
// discovery find the same spectr/ directories, duplicates are properly removed.
func TestFindSpectrRoots_Deduplication(t *testing.T) {
	// Create structure where upward and downward discovery might find duplicates:
	// tmpDir/
	//   .git/
	//   spectr/                  <- found by both upward (from cwd) and would be found by downward
	//   subproject/
	//     spectr/                <- found by downward discovery
	//     subdir/                <- cwd (starting point)
	tmpDir := t.TempDir()

	// Create .git at root
	mustMkdirAll(t, filepath.Join(tmpDir, ".git"))

	// Create main spectr/
	mustMkdirAll(t, filepath.Join(tmpDir, "spectr"))

	// Create subproject with spectr/
	subprojectDir := filepath.Join(tmpDir, "subproject")
	mustMkdirAll(t, filepath.Join(subprojectDir, "spectr"))

	// Create subdirectory as cwd
	subdirPath := filepath.Join(subprojectDir, "subdir")
	mustMkdirAll(t, subdirPath)

	// Run discovery from subdir
	roots, err := FindSpectrRoots(subdirPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should find 2 unique roots: subproject (upward, closest) and tmpDir (upward)
	// No duplicates should exist
	if len(roots) != 2 {
		for i, r := range roots {
			t.Logf("  root[%d]: %s", i, r.Path)
		}
		t.Fatalf("expected 2 roots (subproject and tmpDir), got %d", len(roots))
	}

	// Verify no duplicates by checking each path appears only once
	seen := make(map[string]int)
	for _, root := range roots {
		seen[root.Path]++
	}

	for path, count := range seen {
		if count > 1 {
			t.Errorf("duplicate path found %d times: %s", count, path)
		}
	}

	// Verify correct roots were found (order matters: closest first)
	if roots[0].Path != subprojectDir {
		t.Errorf("expected first root to be subproject (closest), got %s", roots[0].Path)
	}
	if roots[1].Path != tmpDir {
		t.Errorf("expected second root to be tmpDir, got %s", roots[1].Path)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || s != "" && containsAt(s, substr, 0))
}

func containsAt(s, substr string, start int) bool {
	for i := start; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}

	return false
}
