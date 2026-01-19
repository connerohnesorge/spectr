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
			name: "creates single directory",
			paths: []string{
				".claude/commands/spectr",
			},
			wantCreated: []string{
				".claude/commands/spectr",
			},
			wantUpdated: nil,
		},
		{
			name: "creates multiple directories",
			paths: []string{
				".claude/commands",
				".claude/contexts",
			},
			wantCreated: []string{
				".claude/commands",
				".claude/contexts",
			},
			wantUpdated: nil,
		},
		{
			name: "silent success if directory exists",
			paths: []string{
				".claude/commands/spectr",
			},
			existingDirs: []string{
				".claude/commands/spectr",
			},
			wantCreated: nil,
			wantUpdated: nil,
		},
		{
			name: "mixed existing and new directories",
			paths: []string{
				".claude/commands",
				".claude/contexts",
			},
			existingDirs: []string{
				".claude/commands",
			},
			wantCreated: []string{
				".claude/contexts",
			},
			wantUpdated: nil,
		},
		{
			name: "creates nested directories (MkdirAll semantics)",
			paths: []string{
				".claude/commands/spectr/nested",
			},
			wantCreated: []string{
				".claude/commands/spectr/nested",
			},
			wantUpdated: nil,
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
				if err := projectFs.MkdirAll(dir, 0o755); err != nil {
					t.Fatalf(
						"failed to setup existing directory: %v",
						err,
					)
				}
			}

			// Create initializer
			init := NewDirectoryInitializer(
				tt.paths...)

			// Execute
			result, err := init.Init(
				context.Background(),
				projectFs,
				homeFs,
				cfg,
				nil,
			)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Init() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)

				return
			}

			// Check result
			if !stringSliceEqual(
				result.CreatedFiles,
				tt.wantCreated,
			) {
				t.Errorf(
					"Init() CreatedFiles = %v, want %v",
					result.CreatedFiles,
					tt.wantCreated,
				)
			}
			if !stringSliceEqual(
				result.UpdatedFiles,
				tt.wantUpdated,
			) {
				t.Errorf(
					"Init() UpdatedFiles = %v, want %v",
					result.UpdatedFiles,
					tt.wantUpdated,
				)
			}

			// Verify directories exist
			for _, path := range tt.paths {
				exists, err := afero.DirExists(
					projectFs,
					path,
				)
				if err != nil {
					t.Errorf(
						"failed to check directory %s: %v",
						path,
						err,
					)
				}
				if !exists {
					t.Errorf(
						"directory %s does not exist after Init()",
						path,
					)
				}
			}
		})
	}
}

func TestDirectoryInitializer_IsSetup(
	t *testing.T,
) {
	tests := []struct {
		name         string
		paths        []string
		existingDirs []string
		want         bool
	}{
		{
			name: "returns true when all directories exist",
			paths: []string{
				".claude/commands",
			},
			existingDirs: []string{
				".claude/commands",
			},
			want: true,
		},
		{
			name: "returns false when directory does not exist",
			paths: []string{
				".claude/commands",
			},
			existingDirs: nil,
			want:         false,
		},
		{
			name: "returns true when all multiple directories exist",
			paths: []string{
				".claude/commands",
				".claude/contexts",
			},
			existingDirs: []string{
				".claude/commands",
				".claude/contexts",
			},
			want: true,
		},
		{
			name: "returns false when some directories missing",
			paths: []string{
				".claude/commands",
				".claude/contexts",
			},
			existingDirs: []string{
				".claude/commands",
			},
			want: false,
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
				if err := projectFs.MkdirAll(dir, 0o755); err != nil {
					t.Fatalf(
						"failed to setup existing directory: %v",
						err,
					)
				}
			}

			// Create initializer
			init := NewDirectoryInitializer(
				tt.paths...)

			// Execute
			got := init.IsSetup(
				projectFs,
				homeFs,
				cfg,
			)

			// Check result
			if got != tt.want {
				t.Errorf(
					"IsSetup() = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestDirectoryInitializer_dedupeKey(
	t *testing.T,
) {
	tests := []struct {
		name  string
		paths []string
		want  string
	}{
		{
			name:  "simple path",
			paths: []string{".claude/commands"},
			want:  "DirectoryInitializer:.claude/commands",
		},
		{
			name:  "path with trailing slash",
			paths: []string{".claude/commands/"},
			want:  "DirectoryInitializer:.claude/commands",
		},
		{
			name:  "path with dots",
			paths: []string{".claude/./commands"},
			want:  "DirectoryInitializer:.claude/commands",
		},
		{
			name: "multiple paths uses first",
			paths: []string{
				".claude/commands",
				".claude/contexts",
			},
			want: "DirectoryInitializer:.claude/commands",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			init := &DirectoryInitializer{
				paths: tt.paths,
			}
			got := init.dedupeKey()
			if got != tt.want {
				t.Errorf(
					"dedupeKey() = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestHomeDirectoryInitializer_Init(
	t *testing.T,
) {
	tests := []struct {
		name         string
		paths        []string
		existingDirs []string
		wantCreated  []string
		wantUpdated  []string
		wantErr      bool
	}{
		{
			name: "creates single directory in home",
			paths: []string{
				".codex/prompts",
			},
			wantCreated: []string{
				".codex/prompts",
			},
			wantUpdated: nil,
		},
		{
			name: "creates multiple directories in home",
			paths: []string{
				".codex/prompts",
				".codex/configs",
			},
			wantCreated: []string{
				".codex/prompts",
				".codex/configs",
			},
			wantUpdated: nil,
		},
		{
			name: "silent success if directory exists in home",
			paths: []string{
				".codex/prompts",
			},
			existingDirs: []string{
				".codex/prompts",
			},
			wantCreated: nil,
			wantUpdated: nil,
		},
		{
			name: "creates nested directories in home (MkdirAll semantics)",
			paths: []string{
				".codex/prompts/spectr",
			},
			wantCreated: []string{
				".codex/prompts/spectr",
			},
			wantUpdated: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			projectFs := afero.NewMemMapFs()
			homeFs := afero.NewMemMapFs()
			cfg := &Config{SpectrDir: "spectr"}

			// Create existing directories in home filesystem
			for _, dir := range tt.existingDirs {
				if err := homeFs.MkdirAll(dir, 0o755); err != nil {
					t.Fatalf(
						"failed to setup existing directory: %v",
						err,
					)
				}
			}

			// Create initializer
			init := NewHomeDirectoryInitializer(
				tt.paths...)

			// Execute
			result, err := init.Init(
				context.Background(),
				projectFs,
				homeFs,
				cfg,
				nil,
			)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Init() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)

				return
			}

			// Check result
			if !stringSliceEqual(
				result.CreatedFiles,
				tt.wantCreated,
			) {
				t.Errorf(
					"Init() CreatedFiles = %v, want %v",
					result.CreatedFiles,
					tt.wantCreated,
				)
			}
			if !stringSliceEqual(
				result.UpdatedFiles,
				tt.wantUpdated,
			) {
				t.Errorf(
					"Init() UpdatedFiles = %v, want %v",
					result.UpdatedFiles,
					tt.wantUpdated,
				)
			}

			// Verify directories exist in home filesystem
			for _, path := range tt.paths {
				exists, err := afero.DirExists(
					homeFs,
					path,
				)
				if err != nil {
					t.Errorf(
						"failed to check directory %s: %v",
						path,
						err,
					)
				}
				if !exists {
					t.Errorf(
						"directory %s does not exist in home filesystem after Init()",
						path,
					)
				}
			}

			// Verify directories DO NOT exist in project filesystem
			for _, path := range tt.paths {
				exists, _ := afero.DirExists(
					projectFs,
					path,
				)
				if exists {
					t.Errorf(
						"directory %s should NOT exist in project filesystem",
						path,
					)
				}
			}
		})
	}
}

func TestHomeDirectoryInitializer_IsSetup(
	t *testing.T,
) {
	tests := []struct {
		name         string
		paths        []string
		existingDirs []string
		want         bool
	}{
		{
			name: "returns true when directory exists in home",
			paths: []string{
				".codex/prompts",
			},
			existingDirs: []string{
				".codex/prompts",
			},
			want: true,
		},
		{
			name: "returns false when directory does not exist in home",
			paths: []string{
				".codex/prompts",
			},
			existingDirs: nil,
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			projectFs := afero.NewMemMapFs()
			homeFs := afero.NewMemMapFs()
			cfg := &Config{SpectrDir: "spectr"}

			// Create existing directories in home filesystem
			for _, dir := range tt.existingDirs {
				if err := homeFs.MkdirAll(dir, 0o755); err != nil {
					t.Fatalf(
						"failed to setup existing directory: %v",
						err,
					)
				}
			}

			// Create initializer
			init := NewHomeDirectoryInitializer(
				tt.paths...)

			// Execute
			got := init.IsSetup(
				projectFs,
				homeFs,
				cfg,
			)

			// Check result
			if got != tt.want {
				t.Errorf(
					"IsSetup() = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestHomeDirectoryInitializer_dedupeKey(
	t *testing.T,
) {
	tests := []struct {
		name  string
		paths []string
		want  string
	}{
		{
			name:  "simple path",
			paths: []string{".codex/prompts"},
			want:  "HomeDirectoryInitializer:.codex/prompts",
		},
		{
			name:  "path with trailing slash",
			paths: []string{".codex/prompts/"},
			want:  "HomeDirectoryInitializer:.codex/prompts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			init := &HomeDirectoryInitializer{
				paths: tt.paths,
			}
			got := init.dedupeKey()
			if got != tt.want {
				t.Errorf(
					"dedupeKey() = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestDirectoryInitializer_SeparateTypeFromHome(
	t *testing.T,
) {
	// Verify that DirectoryInitializer and HomeDirectoryInitializer have different dedupeKeys
	// even for the same path. This ensures they are not deduplicated against each other.

	projectInit := &DirectoryInitializer{
		paths: []string{".config/tool"},
	}
	homeInit := &HomeDirectoryInitializer{
		paths: []string{".config/tool"},
	}

	projectKey := projectInit.dedupeKey()
	homeKey := homeInit.dedupeKey()

	if projectKey == homeKey {
		t.Errorf(
			"DirectoryInitializer and HomeDirectoryInitializer should have different dedupeKeys, both got: %s",
			projectKey,
		)
	}

	expectedProjectKey := "DirectoryInitializer:.config/tool"
	expectedHomeKey := "HomeDirectoryInitializer:.config/tool"

	if projectKey != expectedProjectKey {
		t.Errorf(
			"DirectoryInitializer dedupeKey() = %v, want %v",
			projectKey,
			expectedProjectKey,
		)
	}
	if homeKey != expectedHomeKey {
		t.Errorf(
			"HomeDirectoryInitializer dedupeKey() = %v, want %v",
			homeKey,
			expectedHomeKey,
		)
	}
}

// Helper function to compare string slices
func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
