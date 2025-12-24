package providers

import (
	"context"
	"testing"

	"github.com/spf13/afero"
)

func TestDirectoryInitializer_Init(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		setupFs        func(afero.Fs)
		wantCreated    []string
		wantUpdated    []string
		wantErr        bool
		checkDirExists bool
	}{
		{
			name:    "creates new directory",
			path:    ".claude/commands/spectr",
			setupFs: func(_ afero.Fs) {},
			wantCreated: []string{
				".claude/commands/spectr",
			},
			wantUpdated:    nil,
			wantErr:        false,
			checkDirExists: true,
		},
		{
			name: "directory already exists",
			path: ".claude/commands/spectr",
			setupFs: func(fs afero.Fs) {
				_ = fs.MkdirAll(
					".claude/commands/spectr",
					0o755,
				)
			},
			wantCreated:    nil,
			wantUpdated:    nil,
			wantErr:        false,
			checkDirExists: true,
		},
		{
			name:    "creates nested directory",
			path:    "deep/nested/directory/structure",
			setupFs: func(_ afero.Fs) {},
			wantCreated: []string{
				"deep/nested/directory/structure",
			},
			wantUpdated:    nil,
			wantErr:        false,
			checkDirExists: true,
		},
		{
			name:           "creates single level directory",
			path:           "simple",
			setupFs:        func(_ afero.Fs) {},
			wantCreated:    []string{"simple"},
			wantUpdated:    nil,
			wantErr:        false,
			checkDirExists: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create in-memory filesystem
			fs := afero.NewMemMapFs()

			// Setup filesystem
			if tt.setupFs != nil {
				tt.setupFs(fs)
			}

			// Create initializer
			init := NewDirectoryInitializer(
				tt.path,
			)

			// Run Init
			cfg := NewDefaultConfig()
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

			// Check result
			if len(
				result.CreatedFiles,
			) != len(
				tt.wantCreated,
			) {
				t.Errorf(
					"Init() CreatedFiles = %v, want %v",
					result.CreatedFiles,
					tt.wantCreated,
				)
			} else {
				for i, path := range result.CreatedFiles {
					if path != tt.wantCreated[i] {
						t.Errorf("Init() CreatedFiles[%d] = %v, want %v", i, path, tt.wantCreated[i])
					}
				}
			}

			if len(
				result.UpdatedFiles,
			) != len(
				tt.wantUpdated,
			) {
				t.Errorf(
					"Init() UpdatedFiles = %v, want %v",
					result.UpdatedFiles,
					tt.wantUpdated,
				)
			}

			// Check directory exists
			if !tt.checkDirExists {
				return
			}

			exists, err := afero.DirExists(
				fs,
				tt.path,
			)
			if err != nil {
				t.Errorf(
					"DirExists() error = %v",
					err,
				)
			}
			if !exists {
				t.Errorf(
					"Directory %s should exist after Init()",
					tt.path,
				)
			}
		})
	}
}

func TestDirectoryInitializer_IsSetup(
	t *testing.T,
) {
	tests := []struct {
		name    string
		path    string
		setupFs func(afero.Fs)
		want    bool
	}{
		{
			name: "directory exists",
			path: ".claude/commands/spectr",
			setupFs: func(fs afero.Fs) {
				_ = fs.MkdirAll(
					".claude/commands/spectr",
					0o755,
				)
			},
			want: true,
		},
		{
			name:    "directory does not exist",
			path:    ".claude/commands/spectr",
			setupFs: func(_ afero.Fs) {},
			want:    false,
		},
		{
			name: "parent exists but not child",
			path: ".claude/commands/spectr",
			setupFs: func(fs afero.Fs) {
				_ = fs.MkdirAll(
					".claude/commands",
					0o755,
				)
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			if tt.setupFs != nil {
				tt.setupFs(fs)
			}

			init := NewDirectoryInitializer(
				tt.path,
			)
			cfg := NewDefaultConfig()
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
	path := ".claude/commands/spectr"
	init := NewDirectoryInitializer(path)

	if got := init.Path(); got != path {
		t.Errorf(
			"Path() = %v, want %v",
			got,
			path,
		)
	}
}

func TestDirectoryInitializer_IsGlobal(
	t *testing.T,
) {
	tests := []struct {
		name string
		init *DirectoryInitializer
		want bool
	}{
		{
			name: "project-relative directory",
			init: NewDirectoryInitializer(
				".claude/commands/spectr",
			),
			want: false,
		},
		{
			name: "global directory",
			init: NewGlobalDirectoryInitializer(
				".config/aider/commands",
			),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.init.IsGlobal(); got != tt.want {
				t.Errorf(
					"IsGlobal() = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestDirectoryInitializer_Idempotent(
	t *testing.T,
) {
	// Test that running Init multiple times is safe and idempotent
	fs := afero.NewMemMapFs()
	init := NewDirectoryInitializer(
		".claude/commands/spectr",
	)
	cfg := NewDefaultConfig()

	// First run - should create directory
	result1, err := init.Init(
		context.Background(),
		fs,
		cfg,
		nil,
	)
	if err != nil {
		t.Fatalf("First Init() error = %v", err)
	}
	if len(result1.CreatedFiles) != 1 {
		t.Errorf(
			"First Init() should create directory, got %v",
			result1.CreatedFiles,
		)
	}

	// Second run - should do nothing (directory already exists)
	result2, err := init.Init(
		context.Background(),
		fs,
		cfg,
		nil,
	)
	if err != nil {
		t.Fatalf("Second Init() error = %v", err)
	}
	if len(result2.CreatedFiles) != 0 {
		t.Errorf(
			"Second Init() should not create directory again, got %v",
			result2.CreatedFiles,
		)
	}

	// Third run - still safe
	result3, err := init.Init(
		context.Background(),
		fs,
		cfg,
		nil,
	)
	if err != nil {
		t.Fatalf("Third Init() error = %v", err)
	}
	if len(result3.CreatedFiles) != 0 {
		t.Errorf(
			"Third Init() should not create directory again, got %v",
			result3.CreatedFiles,
		)
	}

	// Verify directory still exists
	exists, err := afero.DirExists(
		fs,
		".claude/commands/spectr",
	)
	if err != nil {
		t.Fatalf("DirExists() error = %v", err)
	}
	if !exists {
		t.Error(
			"Directory should still exist after multiple Init() calls",
		)
	}
}
