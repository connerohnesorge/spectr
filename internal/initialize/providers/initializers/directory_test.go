package initializers

import (
	"context"
	"testing"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/spf13/afero"
)

func TestDirectoryInitializer_Init(t *testing.T) {
	tests := []struct {
		name        string
		paths       []string
		setup       func(afero.Fs)
		wantCreated []string
		wantErr     bool
	}{
		{
			name:        "creates single directory",
			paths:       []string{".claude/commands/spectr"},
			setup:       func(_ afero.Fs) {},
			wantCreated: []string{".claude/commands/spectr"},
		},
		{
			name:        "creates multiple directories",
			paths:       []string{"dir1", "dir2"},
			setup:       func(_ afero.Fs) {},
			wantCreated: []string{"dir1", "dir2"},
		},
		{
			name:  "silent success if directory exists",
			paths: []string{"existing"},
			setup: func(fs afero.Fs) {
				_ = fs.MkdirAll("existing", 0o755)
			},
			wantCreated: nil, // Empty because dir already existed
		},
		{
			name:  "creates only non-existing directories",
			paths: []string{"existing", "new"},
			setup: func(fs afero.Fs) {
				_ = fs.MkdirAll("existing", 0o755)
			},
			wantCreated: []string{"new"},
		},
		{
			name:        "creates nested directories with MkdirAll",
			paths:       []string{"a/b/c/d"},
			setup:       func(_ afero.Fs) {},
			wantCreated: []string{"a/b/c/d"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectFs := afero.NewMemMapFs()
			homeFs := afero.NewMemMapFs()
			tt.setup(projectFs)

			cfg := &domain.Config{SpectrDir: "spectr"}
			init := NewDirectoryInitializer(tt.paths...)

			result, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if len(result.CreatedFiles) != len(tt.wantCreated) {
				t.Errorf(
					"Init() created %d files, want %d",
					len(result.CreatedFiles),
					len(tt.wantCreated),
				)

				return
			}

			for i, path := range tt.wantCreated {
				if result.CreatedFiles[i] != path {
					t.Errorf("Init() created[%d] = %s, want %s", i, result.CreatedFiles[i], path)
				}
			}

			// Verify directories actually exist
			for _, path := range tt.paths {
				exists, err := afero.DirExists(projectFs, path)
				if err != nil || !exists {
					t.Errorf("Directory %s should exist after Init()", path)
				}
			}
		})
	}
}

func TestDirectoryInitializer_IsSetup(t *testing.T) {
	tests := []struct {
		name  string
		paths []string
		setup func(afero.Fs)
		want  bool
	}{
		{
			name:  "returns true when all directories exist",
			paths: []string{"dir1", "dir2"},
			setup: func(fs afero.Fs) {
				_ = fs.MkdirAll("dir1", 0o755)
				_ = fs.MkdirAll("dir2", 0o755)
			},
			want: true,
		},
		{
			name:  "returns false when some directories missing",
			paths: []string{"dir1", "dir2"},
			setup: func(fs afero.Fs) {
				_ = fs.MkdirAll("dir1", 0o755)
			},
			want: false,
		},
		{
			name:  "returns false when no directories exist",
			paths: []string{"dir1"},
			setup: func(_ afero.Fs) {},
			want:  false,
		},
		{
			name:  "returns true for empty paths slice",
			paths: nil,
			setup: func(_ afero.Fs) {},
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectFs := afero.NewMemMapFs()
			homeFs := afero.NewMemMapFs()
			tt.setup(projectFs)

			cfg := &domain.Config{SpectrDir: "spectr"}
			init := NewDirectoryInitializer(tt.paths...)

			if got := init.IsSetup(projectFs, homeFs, cfg); got != tt.want {
				t.Errorf("IsSetup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirectoryInitializer_dedupeKey(t *testing.T) {
	tests := []struct {
		name  string
		paths []string
		want  string
	}{
		{
			name:  "uses first path with Clean",
			paths: []string{".claude/commands/spectr"},
			want:  "DirectoryInitializer:.claude/commands/spectr",
		},
		{
			name:  "normalizes path with trailing slash",
			paths: []string{".claude/commands/spectr/"},
			want:  "DirectoryInitializer:.claude/commands/spectr",
		},
		{
			name:  "handles empty paths",
			paths: nil,
			want:  "DirectoryInitializer:",
		},
		{
			name:  "uses only first path when multiple",
			paths: []string{"first", "second"},
			want:  "DirectoryInitializer:first",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			init, ok := NewDirectoryInitializer(tt.paths...).(*DirectoryInitializer)
			if !ok {
				t.Fatal("NewDirectoryInitializer did not return *DirectoryInitializer")
			}
			if got := init.DedupeKey(); got != tt.want {
				t.Errorf("dedupeKey() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHomeDirectoryInitializer_Init(t *testing.T) {
	tests := []struct {
		name        string
		paths       []string
		setup       func(afero.Fs)
		wantCreated []string
		wantErr     bool
	}{
		{
			name:        "creates directory in home filesystem",
			paths:       []string{".codex/prompts"},
			setup:       func(_ afero.Fs) {},
			wantCreated: []string{".codex/prompts"},
		},
		{
			name:  "silent success if directory exists",
			paths: []string{"existing"},
			setup: func(fs afero.Fs) {
				_ = fs.MkdirAll("existing", 0o755)
			},
			wantCreated: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectFs := afero.NewMemMapFs()
			homeFs := afero.NewMemMapFs()
			tt.setup(homeFs) // Setup on homeFs, not projectFs

			cfg := &domain.Config{SpectrDir: "spectr"}
			init := NewHomeDirectoryInitializer(tt.paths...)

			result, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if len(result.CreatedFiles) != len(tt.wantCreated) {
				t.Errorf(
					"Init() created %d files, want %d",
					len(result.CreatedFiles),
					len(tt.wantCreated),
				)

				return
			}

			// Verify directories exist in homeFs, not projectFs
			for _, path := range tt.paths {
				exists, err := afero.DirExists(homeFs, path)
				if err != nil || !exists {
					t.Errorf("Directory %s should exist in homeFs after Init()", path)
				}

				// Should NOT exist in projectFs
				exists, _ = afero.DirExists(projectFs, path)
				if exists {
					t.Errorf("Directory %s should NOT exist in projectFs", path)
				}
			}
		})
	}
}

func TestHomeDirectoryInitializer_IsSetup(t *testing.T) {
	tests := []struct {
		name  string
		paths []string
		setup func(afero.Fs)
		want  bool
	}{
		{
			name:  "checks home filesystem",
			paths: []string{".codex/prompts"},
			setup: func(fs afero.Fs) {
				_ = fs.MkdirAll(".codex/prompts", 0o755)
			},
			want: true,
		},
		{
			name:  "returns false if not in home filesystem",
			paths: []string{".codex/prompts"},
			setup: func(_ afero.Fs) {},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectFs := afero.NewMemMapFs()
			homeFs := afero.NewMemMapFs()
			tt.setup(homeFs)

			cfg := &domain.Config{SpectrDir: "spectr"}
			init := NewHomeDirectoryInitializer(tt.paths...)

			if got := init.IsSetup(projectFs, homeFs, cfg); got != tt.want {
				t.Errorf("IsSetup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHomeDirectoryInitializer_dedupeKey(t *testing.T) {
	tests := []struct {
		name  string
		paths []string
		want  string
	}{
		{
			name:  "uses HomeDirectoryInitializer prefix",
			paths: []string{".codex/prompts"},
			want:  "HomeDirectoryInitializer:.codex/prompts",
		},
		{
			name:  "handles empty paths",
			paths: nil,
			want:  "HomeDirectoryInitializer:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			init, ok := NewHomeDirectoryInitializer(tt.paths...).(*HomeDirectoryInitializer)
			if !ok {
				t.Fatal("NewHomeDirectoryInitializer did not return *HomeDirectoryInitializer")
			}
			if got := init.DedupeKey(); got != tt.want {
				t.Errorf("dedupeKey() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDirectoryInitializer_UsesProjectFs(t *testing.T) {
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()

	cfg := &domain.Config{SpectrDir: "spectr"}
	init := NewDirectoryInitializer("testdir")

	_, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Should exist in projectFs
	exists, _ := afero.DirExists(projectFs, "testdir")
	if !exists {
		t.Error("Directory should exist in projectFs")
	}

	// Should NOT exist in homeFs
	exists, _ = afero.DirExists(homeFs, "testdir")
	if exists {
		t.Error("Directory should NOT exist in homeFs")
	}
}

func TestHomeDirectoryInitializer_UsesHomeFs(t *testing.T) {
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()

	cfg := &domain.Config{SpectrDir: "spectr"}
	init := NewHomeDirectoryInitializer("testdir")

	_, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Should exist in homeFs
	exists, _ := afero.DirExists(homeFs, "testdir")
	if !exists {
		t.Error("Directory should exist in homeFs")
	}

	// Should NOT exist in projectFs
	exists, _ = afero.DirExists(projectFs, "testdir")
	if exists {
		t.Error("Directory should NOT exist in projectFs")
	}
}
