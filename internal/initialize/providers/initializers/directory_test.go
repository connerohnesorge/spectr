package initializers

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
		wantErr      bool
	}{
		{
			name:        "creates single directory",
			paths:       []string{"test/dir"},
			wantCreated: []string{"test/dir"},
			wantErr:     false,
		},
		{
			name: "creates multiple directories",
			paths: []string{
				"dir1",
				"dir2",
				"dir3",
			},
			wantCreated: []string{
				"dir1",
				"dir2",
				"dir3",
			},
			wantErr: false,
		},
		{
			name:         "skips existing directory",
			paths:        []string{"test/dir"},
			existingDirs: []string{"test/dir"},
			wantCreated:  nil,
			wantErr:      false,
		},
		{
			name: "creates only missing directories",
			paths: []string{
				"dir1",
				"dir2",
				"dir3",
			},
			existingDirs: []string{"dir2"},
			wantCreated: []string{
				"dir1",
				"dir3",
			},
			wantErr: false,
		},
		{
			name:        "creates nested directories",
			paths:       []string{"a/b/c/d"},
			wantCreated: []string{"a/b/c/d"},
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup in-memory filesystem
			fs := afero.NewMemMapFs()
			cfg := DefaultConfig()

			// Create existing directories
			for _, dir := range tt.existingDirs {
				if err := fs.MkdirAll(dir, 0755); err != nil {
					t.Fatalf(
						"failed to create existing dir: %v",
						err,
					)
				}
			}

			// Run initializer
			init := NewDirectoryInitializer(
				tt.paths...)
			result, err := init.Init(
				context.Background(),
				fs,
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

			// Check created files
			if len(
				result.CreatedFiles,
			) != len(
				tt.wantCreated,
			) {
				t.Errorf(
					"Init() created %d files, want %d",
					len(result.CreatedFiles),
					len(tt.wantCreated),
				)
			}

			for _, path := range tt.wantCreated {
				found := false
				for _, created := range result.CreatedFiles {
					if created == path {
						found = true

						break
					}
				}
				if !found {
					t.Errorf(
						"Init() did not create expected path: %s",
						path,
					)
				}
			}

			// Verify directories exist on filesystem
			for _, path := range tt.paths {
				exists, err := afero.DirExists(
					fs,
					path,
				)
				if err != nil {
					t.Errorf(
						"error checking directory %s: %v",
						path,
						err,
					)
				}
				if !exists {
					t.Errorf(
						"directory %s should exist but doesn't",
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
				"dir1",
				"dir2",
			},
			existingDirs: []string{
				"dir1",
				"dir2",
			},
			want: true,
		},
		{
			name: "returns false when some directories missing",
			paths: []string{
				"dir1",
				"dir2",
				"dir3",
			},
			existingDirs: []string{"dir1"},
			want:         false,
		},
		{
			name: "returns false when no directories exist",
			paths: []string{
				"dir1",
				"dir2",
			},
			existingDirs: nil,
			want:         false,
		},
		{
			name:         "returns true for single existing directory",
			paths:        []string{"test/dir"},
			existingDirs: []string{"test/dir"},
			want:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup in-memory filesystem
			fs := afero.NewMemMapFs()
			cfg := DefaultConfig()

			// Create existing directories
			for _, dir := range tt.existingDirs {
				if err := fs.MkdirAll(dir, 0755); err != nil {
					t.Fatalf(
						"failed to create existing dir: %v",
						err,
					)
				}
			}

			// Check IsSetup
			init := NewDirectoryInitializer(
				tt.paths...)
			got := init.IsSetup(fs, cfg)

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

func TestDirectoryInitializer_Path(t *testing.T) {
	tests := []struct {
		name  string
		paths []string
		want  string
	}{
		{
			name: "returns first path",
			paths: []string{
				"dir1",
				"dir2",
				"dir3",
			},
			want: "dir1",
		},
		{
			name:  "returns single path",
			paths: []string{"test/dir"},
			want:  "test/dir",
		},
		{
			name:  "returns empty string for empty paths",
			paths: nil,
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			init := NewDirectoryInitializer(
				tt.paths...)
			got := init.Path()

			if got != tt.want {
				t.Errorf(
					"Path() = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestDirectoryInitializer_IsGlobal(
	t *testing.T,
) {
	init := NewDirectoryInitializer("test/dir")
	if init.IsGlobal() {
		t.Error(
			"IsGlobal() should return false for DirectoryInitializer",
		)
	}
}

func TestDirectoryInitializer_Idempotent(
	t *testing.T,
) {
	fs := afero.NewMemMapFs()
	cfg := DefaultConfig()
	ctx := context.Background()

	init := NewDirectoryInitializer(
		"test/dir1",
		"test/dir2",
	)

	// First run
	result1, err := init.Init(ctx, fs, cfg, nil)
	if err != nil {
		t.Fatalf("first Init() failed: %v", err)
	}

	if len(result1.CreatedFiles) != 2 {
		t.Errorf(
			"first Init() created %d files, want 2",
			len(result1.CreatedFiles),
		)
	}

	// Second run - should be idempotent
	result2, err := init.Init(ctx, fs, cfg, nil)
	if err != nil {
		t.Fatalf("second Init() failed: %v", err)
	}

	if len(result2.CreatedFiles) != 0 {
		t.Errorf(
			"second Init() created %d files, want 0 (idempotent)",
			len(result2.CreatedFiles),
		)
	}

	// Verify directories still exist
	if !init.IsSetup(fs, cfg) {
		t.Error(
			"directories should still exist after second Init()",
		)
	}
}
