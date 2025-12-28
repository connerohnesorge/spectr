package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_DefaultWhenMissing(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	cfg, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Dir != DefaultDir {
		t.Errorf("expected Dir=%q, got %q", DefaultDir, cfg.Dir)
	}
}

func TestLoad_ReadsConfigFile(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ConfigFileName)

	content := []byte("dir: my-custom-dir\n")
	if err := os.WriteFile(configPath, content, 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Dir != "my-custom-dir" {
		t.Errorf("expected Dir=%q, got %q", "my-custom-dir", cfg.Dir)
	}
}

func TestLoad_DefaultsEmptyDir(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ConfigFileName)

	content := []byte("# empty config\n")
	if err := os.WriteFile(configPath, content, 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Dir != DefaultDir {
		t.Errorf("expected Dir=%q, got %q", DefaultDir, cfg.Dir)
	}
}

func TestLoad_InvalidAbsolutePath(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ConfigFileName)

	content := []byte("dir: /absolute/path\n")
	if err := os.WriteFile(configPath, content, 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	_, err := Load(tmpDir)
	if err == nil {
		t.Fatal("expected error for absolute path, got nil")
	}
}

func TestLoad_InvalidPathTraversal(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ConfigFileName)

	content := []byte("dir: ../outside\n")
	if err := os.WriteFile(configPath, content, 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	_, err := Load(tmpDir)
	if err == nil {
		t.Fatal("expected error for path traversal, got nil")
	}
}

func TestConfig_DerivedPaths(t *testing.T) {
	t.Parallel()

	cfg := &Config{Dir: "openspec"}

	tests := []struct {
		name     string
		got      string
		expected string
	}{
		{"SpecsDir", cfg.SpecsDir(), "openspec/specs"},
		{"ChangesDir", cfg.ChangesDir(), "openspec/changes"},
		{"ProjectFile", cfg.ProjectFile(), "openspec/project.md"},
		{"AgentsFile", cfg.AgentsFile(), "openspec/AGENTS.md"},
	}

	for _, tc := range tests {
		if tc.got != tc.expected {
			t.Errorf("%s: expected %q, got %q", tc.name, tc.expected, tc.got)
		}
	}
}

func TestDefault(t *testing.T) {
	t.Parallel()

	cfg := Default()
	if cfg.Dir != DefaultDir {
		t.Errorf("expected Dir=%q, got %q", DefaultDir, cfg.Dir)
	}
}
