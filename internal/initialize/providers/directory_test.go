package providers

import (
	"context"
	"testing"

	"github.com/spf13/afero"
)

func TestDirectoryInitializer_Init(t *testing.T) {
	tests := []struct {
		name         string
		paths        []string
		existingDirs []string
		wantCreated  []string
		wantUpdated  []string
		wantErr      bool
	}{
		{
			name:        "creates new directory",
			paths:       []string{".claude/commands/spectr"},
			wantCreated: []string{".claude/commands/spectr"},
			wantUpdated: nil,
		},
		{
			name:        "creates multiple directories",
			paths:       []string{".claude/commands/spectr", ".claude/contexts"},
			wantCreated: []string{".claude/commands/spectr", ".claude/contexts"},
			wantUpdated: nil,
		},
		{
			name:         "silent success if directory exists",
			paths:        []string{".claude/commands/spectr"},
			existingDirs: []string{".claude/commands/spectr"},
			wantCreated:  nil,
			wantUpdated:  nil,
		},
		{
			name:         "creates some, skips existing",
			paths:        []string{".claude/commands/spectr", ".claude/contexts"},
			existingDirs: []string{".claude/commands/spectr"},
			wantCreated:  []string{".claude/contexts"},
			wantUpdated:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			projectFs := afero.NewMemMapFs()
			homeFs := afero.NewMemMapFs()
			cfg := &Config{SpectrDir: "spectr"}

			// Create existing directories
			for _, dir := range tt.existingDirs {
				if err := projectFs.MkdirAll(dir, 0755); err != nil {
					t.Fatalf("failed to create existing dir: %v", err)
				}
			}

			// Test
			init := NewDirectoryInitializer(tt.paths...)
			result, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)

			// Verify error
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			// Verify result
			if !equalStringSlices(result.CreatedFiles, tt.wantCreated) {
				t.Errorf("Init() CreatedFiles = %v, want %v", result.CreatedFiles, tt.wantCreated)
			}
			if !equalStringSlices(result.UpdatedFiles, tt.wantUpdated) {
				t.Errorf("Init() UpdatedFiles = %v, want %v", result.UpdatedFiles, tt.wantUpdated)
			}

			// Verify directories exist
			for _, path := range tt.paths {
				exists, err := afero.DirExists(projectFs, path)
				if err != nil {
					t.Errorf("failed to check directory %s: %v", path, err)
				}
				if !exists {
					t.Errorf("directory %s does not exist", path)
				}
			}
		})
	}
}

func TestDirectoryInitializer_IsSetup(t *testing.T) {
	tests := []struct {
		name         string
		paths        []string
		existingDirs []string
		want         bool
	}{
		{
			name:         "returns true if all directories exist",
			paths:        []string{".claude/commands/spectr"},
			existingDirs: []string{".claude/commands/spectr"},
			want:         true,
		},
		{
			name:         "returns false if any directory missing",
			paths:        []string{".claude/commands/spectr", ".claude/contexts"},
			existingDirs: []string{".claude/commands/spectr"},
			want:         false,
		},
		{
			name:  "returns false if no directories exist",
			paths: []string{".claude/commands/spectr"},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			projectFs := afero.NewMemMapFs()
			homeFs := afero.NewMemMapFs()
			cfg := &Config{SpectrDir: "spectr"}

			// Create existing directories
			for _, dir := range tt.existingDirs {
				if err := projectFs.MkdirAll(dir, 0755); err != nil {
					t.Fatalf("failed to create existing dir: %v", err)
				}
			}

			// Test
			init := NewDirectoryInitializer(tt.paths...)
			got := init.IsSetup(projectFs, homeFs, cfg)

			if got != tt.want {
				t.Errorf("IsSetup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirectoryInitializer_DedupeKey(t *testing.T) {
	tests := []struct {
		name  string
		paths []string
		want  string
	}{
		{
			name:  "single path",
			paths: []string{".claude/commands/spectr"},
			want:  "DirectoryInitializer:.claude/commands/spectr",
		},
		{
			name:  "multiple paths",
			paths: []string{".claude/commands/spectr", ".claude/contexts"},
			want:  "DirectoryInitializer:.claude/commands/spectr:.claude/contexts",
		},
		{
			name:  "normalizes path",
			paths: []string{".claude/commands/spectr/"},
			want:  "DirectoryInitializer:.claude/commands/spectr",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			init := NewDirectoryInitializer(tt.paths...)
			got := init.dedupeKey()
			if got != tt.want {
				t.Errorf("dedupeKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHomeDirectoryInitializer_Init(t *testing.T) {
	tests := []struct {
		name         string
		paths        []string
		existingDirs []string
		wantCreated  []string
		wantUpdated  []string
		wantErr      bool
	}{
		{
			name:        "creates new directory in home",
			paths:       []string{".codex/prompts"},
			wantCreated: []string{".codex/prompts"},
			wantUpdated: nil,
		},
		{
			name:        "creates multiple directories in home",
			paths:       []string{".codex/prompts", ".codex/config"},
			wantCreated: []string{".codex/prompts", ".codex/config"},
			wantUpdated: nil,
		},
		{
			name:         "silent success if directory exists in home",
			paths:        []string{".codex/prompts"},
			existingDirs: []string{".codex/prompts"},
			wantCreated:  nil,
			wantUpdated:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup - note we use homeFs, not projectFs
			projectFs := afero.NewMemMapFs()
			homeFs := afero.NewMemMapFs()
			cfg := &Config{SpectrDir: "spectr"}

			// Create existing directories in homeFs
			for _, dir := range tt.existingDirs {
				if err := homeFs.MkdirAll(dir, 0755); err != nil {
					t.Fatalf("failed to create existing dir: %v", err)
				}
			}

			// Test
			init := NewHomeDirectoryInitializer(tt.paths...)
			result, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)

			// Verify error
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			// Verify result
			if !equalStringSlices(result.CreatedFiles, tt.wantCreated) {
				t.Errorf("Init() CreatedFiles = %v, want %v", result.CreatedFiles, tt.wantCreated)
			}
			if !equalStringSlices(result.UpdatedFiles, tt.wantUpdated) {
				t.Errorf("Init() UpdatedFiles = %v, want %v", result.UpdatedFiles, tt.wantUpdated)
			}

			// Verify directories exist in homeFs (not projectFs)
			for _, path := range tt.paths {
				exists, err := afero.DirExists(homeFs, path)
				if err != nil {
					t.Errorf("failed to check directory %s: %v", path, err)
				}
				if !exists {
					t.Errorf("directory %s does not exist in homeFs", path)
				}
			}
		})
	}
}

func TestHomeDirectoryInitializer_IsSetup(t *testing.T) {
	tests := []struct {
		name         string
		paths        []string
		existingDirs []string
		want         bool
	}{
		{
			name:         "returns true if all directories exist in home",
			paths:        []string{".codex/prompts"},
			existingDirs: []string{".codex/prompts"},
			want:         true,
		},
		{
			name:  "returns false if directory missing in home",
			paths: []string{".codex/prompts"},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			projectFs := afero.NewMemMapFs()
			homeFs := afero.NewMemMapFs()
			cfg := &Config{SpectrDir: "spectr"}

			// Create existing directories in homeFs
			for _, dir := range tt.existingDirs {
				if err := homeFs.MkdirAll(dir, 0755); err != nil {
					t.Fatalf("failed to create existing dir: %v", err)
				}
			}

			// Test
			init := NewHomeDirectoryInitializer(tt.paths...)
			got := init.IsSetup(projectFs, homeFs, cfg)

			if got != tt.want {
				t.Errorf("IsSetup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHomeDirectoryInitializer_DedupeKey(t *testing.T) {
	tests := []struct {
		name  string
		paths []string
		want  string
	}{
		{
			name:  "single path",
			paths: []string{".codex/prompts"},
			want:  "HomeDirectoryInitializer:.codex/prompts",
		},
		{
			name:  "normalizes path",
			paths: []string{".codex/prompts/"},
			want:  "HomeDirectoryInitializer:.codex/prompts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			init := NewHomeDirectoryInitializer(tt.paths...)
			got := init.dedupeKey()
			if got != tt.want {
				t.Errorf("dedupeKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to compare string slices
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
