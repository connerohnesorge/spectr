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
		wantErr      bool
	}{
		{
			name:        "creates single directory",
			paths:       []string{".claude/commands"},
			wantCreated: []string{".claude/commands"},
		},
		{
			name:        "creates multiple directories",
			paths:       []string{".claude/commands", ".claude/contexts"},
			wantCreated: []string{".claude/commands", ".claude/contexts"},
		},
		{
			name:         "skips existing directory",
			paths:        []string{".claude/commands"},
			existingDirs: []string{".claude/commands"},
			wantCreated:  nil,
		},
		{
			name:         "creates only missing directories",
			paths:        []string{".claude/commands", ".claude/contexts"},
			existingDirs: []string{".claude/commands"},
			wantCreated:  []string{".claude/contexts"},
		},
		{
			name:        "creates nested directories",
			paths:       []string{".claude/commands/spectr/deep"},
			wantCreated: []string{".claude/commands/spectr/deep"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create in-memory filesystem
			fs := afero.NewMemMapFs()

			// Create existing directories
			for _, dir := range tt.existingDirs {
				if err := fs.MkdirAll(dir, 0755); err != nil {
					t.Fatalf("failed to create existing directory: %v", err)
				}
			}

			// Create initializer and run Init
			init := NewDirectoryInitializer(tt.paths...)
			cfg := &Config{SpectrDir: "spectr"}
			ctx := context.Background()

			result, err := init.Init(ctx, fs, cfg, nil)

			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if err != nil {
				return
			}

			// Verify created files list
			if len(result.CreatedFiles) != len(tt.wantCreated) {
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
					t.Errorf("Init() did not create %s", path)
				}
			}

			// Verify all directories exist
			for _, path := range tt.paths {
				exists, err := afero.DirExists(fs, path)
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
			name:         "single directory exists",
			paths:        []string{".claude/commands"},
			existingDirs: []string{".claude/commands"},
			want:         true,
		},
		{
			name:         "single directory missing",
			paths:        []string{".claude/commands"},
			existingDirs: nil,
			want:         false,
		},
		{
			name:         "all directories exist",
			paths:        []string{".claude/commands", ".claude/contexts"},
			existingDirs: []string{".claude/commands", ".claude/contexts"},
			want:         true,
		},
		{
			name:         "some directories missing",
			paths:        []string{".claude/commands", ".claude/contexts"},
			existingDirs: []string{".claude/commands"},
			want:         false,
		},
		{
			name:         "no directories exist",
			paths:        []string{".claude/commands", ".claude/contexts"},
			existingDirs: nil,
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()

			// Create existing directories
			for _, dir := range tt.existingDirs {
				if err := fs.MkdirAll(dir, 0755); err != nil {
					t.Fatalf("failed to create existing directory: %v", err)
				}
			}

			init := NewDirectoryInitializer(tt.paths...)
			cfg := &Config{SpectrDir: "spectr"}

			if got := init.IsSetup(fs, cfg); got != tt.want {
				t.Errorf("IsSetup() = %v, want %v", got, tt.want)
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
			name:  "single path",
			paths: []string{".claude/commands"},
			want:  ".claude/commands",
		},
		{
			name:  "multiple paths returns first",
			paths: []string{".claude/commands", ".claude/contexts"},
			want:  ".claude/commands",
		},
		{
			name:  "no paths",
			paths: nil,
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			init := NewDirectoryInitializer(tt.paths...)
			if got := init.Path(); got != tt.want {
				t.Errorf("Path() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirectoryInitializer_IsGlobal(t *testing.T) {
	tests := []struct {
		name     string
		isGlobal bool
		want     bool
	}{
		{
			name:     "default is not global",
			isGlobal: false,
			want:     false,
		},
		{
			name:     "configured as global",
			isGlobal: true,
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			init := NewDirectoryInitializer(".claude/commands")
			if tt.isGlobal {
				init = init.WithGlobal(true)
			}

			if got := init.IsGlobal(); got != tt.want {
				t.Errorf("IsGlobal() = %v, want %v", got, tt.want)
			}
		})
	}
}
