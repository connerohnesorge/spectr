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
