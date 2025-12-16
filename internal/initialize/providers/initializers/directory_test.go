package initializers

import (
	"context"
	"testing"

	"github.com/connerohnesorge/spectr/internal/initialize/providers"
	"github.com/spf13/afero"
)

func TestNewDirectoryInitializer(t *testing.T) {
	tests := []struct {
		name    string
		paths   []string
		wantLen int
	}{
		{
			name:    "single path",
			paths:   []string{".claude/commands"},
			wantLen: 1,
		},
		{
			name: "multiple paths",
			paths: []string{
				".claude/commands",
				".gemini/commands",
			},
			wantLen: 2,
		},
		{
			name:    "no paths",
			paths:   nil,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDirectoryInitializer(
				tt.paths...)

			if d == nil {
				t.Fatal(
					"NewDirectoryInitializer() returned nil",
				)
			}

			if len(d.Paths) != tt.wantLen {
				t.Errorf(
					"len(Paths) = %d, want %d",
					len(d.Paths),
					tt.wantLen,
				)
			}

			for i, path := range tt.paths {
				if d.Paths[i] != path {
					t.Errorf(
						"Paths[%d] = %s, want %s",
						i,
						d.Paths[i],
						path,
					)
				}
			}
		})
	}
}

func TestDirectoryInitializer_Init_SingleDirectory(
	t *testing.T,
) {
	fs := afero.NewMemMapFs()
	ctx := context.Background()
	cfg := providers.NewConfig()

	d := NewDirectoryInitializer(
		".claude/commands/spectr",
	)

	err := d.Init(ctx, fs, cfg)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Verify directory was created
	info, err := fs.Stat(
		".claude/commands/spectr",
	)
	if err != nil {
		t.Fatalf(
			"Directory was not created: %v",
			err,
		)
	}

	if !info.IsDir() {
		t.Error("Created path is not a directory")
	}
}

func TestDirectoryInitializer_Init_MultipleDirectories(
	t *testing.T,
) {
	fs := afero.NewMemMapFs()
	ctx := context.Background()
	cfg := providers.NewConfig()

	paths := []string{
		".claude/commands/spectr",
		".gemini/commands/spectr",
		".cursor/rules",
	}
	d := NewDirectoryInitializer(paths...)

	err := d.Init(ctx, fs, cfg)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Verify all directories were created
	for _, path := range paths {
		info, err := fs.Stat(path)
		if err != nil {
			t.Errorf(
				"Directory %s was not created: %v",
				path,
				err,
			)

			continue
		}

		if !info.IsDir() {
			t.Errorf(
				"Path %s is not a directory",
				path,
			)
		}
	}
}

func TestDirectoryInitializer_Init_NestedDirectories(
	t *testing.T,
) {
	fs := afero.NewMemMapFs()
	ctx := context.Background()
	cfg := providers.NewConfig()

	// Test deeply nested directory
	d := NewDirectoryInitializer("a/b/c/d/e/f")

	err := d.Init(ctx, fs, cfg)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Verify the full path was created
	info, err := fs.Stat("a/b/c/d/e/f")
	if err != nil {
		t.Fatalf(
			"Nested directory was not created: %v",
			err,
		)
	}

	if !info.IsDir() {
		t.Error("Created path is not a directory")
	}

	// Verify intermediate directories exist
	intermediates := []string{
		"a",
		"a/b",
		"a/b/c",
		"a/b/c/d",
		"a/b/c/d/e",
	}
	for _, path := range intermediates {
		info, err := fs.Stat(path)
		if err != nil {
			t.Errorf(
				"Intermediate directory %s was not created: %v",
				path,
				err,
			)

			continue
		}

		if !info.IsDir() {
			t.Errorf(
				"Intermediate path %s is not a directory",
				path,
			)
		}
	}
}

func TestDirectoryInitializer_Init_Idempotent(
	t *testing.T,
) {
	fs := afero.NewMemMapFs()
	ctx := context.Background()
	cfg := providers.NewConfig()

	d := NewDirectoryInitializer(
		".claude/commands/spectr",
	)

	// Call Init multiple times
	for i := range 3 {
		err := d.Init(ctx, fs, cfg)
		if err != nil {
			t.Fatalf(
				"Init() call %d failed: %v",
				i+1,
				err,
			)
		}
	}

	// Verify directory still exists
	info, err := fs.Stat(
		".claude/commands/spectr",
	)
	if err != nil {
		t.Fatalf(
			"Directory does not exist after multiple Init calls: %v",
			err,
		)
	}

	if !info.IsDir() {
		t.Error(
			"Path is not a directory after multiple Init calls",
		)
	}
}

func TestDirectoryInitializer_Init_EmptyPaths(
	t *testing.T,
) {
	fs := afero.NewMemMapFs()
	ctx := context.Background()
	cfg := providers.NewConfig()

	d := NewDirectoryInitializer()

	// Should succeed with no paths
	err := d.Init(ctx, fs, cfg)
	if err != nil {
		t.Fatalf(
			"Init() with empty paths should succeed, got: %v",
			err,
		)
	}
}

func TestDirectoryInitializer_IsSetup_NotExists(
	t *testing.T,
) {
	fs := afero.NewMemMapFs()
	cfg := providers.NewConfig()

	d := NewDirectoryInitializer(
		".claude/commands/spectr",
	)

	// Directory doesn't exist yet
	if d.IsSetup(fs, cfg) {
		t.Error(
			"IsSetup() should return false when directory doesn't exist",
		)
	}
}

func TestDirectoryInitializer_IsSetup_Exists(
	t *testing.T,
) {
	fs := afero.NewMemMapFs()
	ctx := context.Background()
	cfg := providers.NewConfig()

	d := NewDirectoryInitializer(
		".claude/commands/spectr",
	)

	// Create the directory
	err := d.Init(ctx, fs, cfg)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Now IsSetup should return true
	if !d.IsSetup(fs, cfg) {
		t.Error(
			"IsSetup() should return true when directory exists",
		)
	}
}

func TestDirectoryInitializer_IsSetup_AllDirectoriesExist(
	t *testing.T,
) {
	fs := afero.NewMemMapFs()
	ctx := context.Background()
	cfg := providers.NewConfig()

	paths := []string{
		".claude/commands",
		".gemini/commands",
	}
	d := NewDirectoryInitializer(paths...)

	// Create all directories
	err := d.Init(ctx, fs, cfg)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// IsSetup should return true when all directories exist
	if !d.IsSetup(fs, cfg) {
		t.Error(
			"IsSetup() should return true when all directories exist",
		)
	}
}

func TestDirectoryInitializer_IsSetup_PartialSetup(
	t *testing.T,
) {
	fs := afero.NewMemMapFs()
	cfg := providers.NewConfig()

	paths := []string{
		".claude/commands",
		".gemini/commands",
	}
	d := NewDirectoryInitializer(paths...)

	// Create only one directory
	err := fs.MkdirAll(".claude/commands", 0755)
	if err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}

	// IsSetup should return false when not all directories exist
	if d.IsSetup(fs, cfg) {
		t.Error(
			"IsSetup() should return false when only some directories exist",
		)
	}
}

func TestDirectoryInitializer_IsSetup_FileNotDirectory(
	t *testing.T,
) {
	fs := afero.NewMemMapFs()
	cfg := providers.NewConfig()

	d := NewDirectoryInitializer(
		".claude/commands",
	)

	// Create a file instead of a directory
	err := afero.WriteFile(
		fs,
		".claude/commands",
		[]byte("file content"),
		0644,
	)
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	// IsSetup should return false when path is a file, not a directory
	if d.IsSetup(fs, cfg) {
		t.Error(
			"IsSetup() should return false when path is a file, not a directory",
		)
	}
}

func TestDirectoryInitializer_IsSetup_EmptyPaths(
	t *testing.T,
) {
	fs := afero.NewMemMapFs()
	cfg := providers.NewConfig()

	d := NewDirectoryInitializer()

	// Empty paths should return true (nothing to check)
	if !d.IsSetup(fs, cfg) {
		t.Error(
			"IsSetup() with empty paths should return true",
		)
	}
}

func TestDirectoryInitializer_Key_SinglePath(
	t *testing.T,
) {
	d := NewDirectoryInitializer(
		".claude/commands",
	)

	key := d.Key()
	expected := "dir:.claude/commands"

	if key != expected {
		t.Errorf(
			"Key() = %s, want %s",
			key,
			expected,
		)
	}
}

func TestDirectoryInitializer_Key_MultiplePaths(
	t *testing.T,
) {
	// Note: paths are sorted in Key()
	d := NewDirectoryInitializer(
		".gemini/commands",
		".claude/commands",
	)

	key := d.Key()
	expected := "dir:.claude/commands,.gemini/commands"

	if key != expected {
		t.Errorf(
			"Key() = %s, want %s",
			key,
			expected,
		)
	}
}

func TestDirectoryInitializer_Key_EmptyPaths(
	t *testing.T,
) {
	d := NewDirectoryInitializer()

	key := d.Key()
	expected := "dir:"

	if key != expected {
		t.Errorf(
			"Key() = %s, want %s",
			key,
			expected,
		)
	}
}

func TestDirectoryInitializer_Key_Consistent(
	t *testing.T,
) {
	d := NewDirectoryInitializer(
		".claude/commands",
		".gemini/commands",
	)

	// Key should be consistent across multiple calls
	key1 := d.Key()
	key2 := d.Key()
	key3 := d.Key()

	if key1 != key2 || key2 != key3 {
		t.Errorf(
			"Key() is not consistent: %s, %s, %s",
			key1,
			key2,
			key3,
		)
	}
}

func TestDirectoryInitializer_Key_OrderIndependent(
	t *testing.T,
) {
	d1 := NewDirectoryInitializer(
		".claude/commands",
		".gemini/commands",
	)
	d2 := NewDirectoryInitializer(
		".gemini/commands",
		".claude/commands",
	)

	// Keys should be the same regardless of path order
	if d1.Key() != d2.Key() {
		t.Errorf(
			"Keys differ for same paths in different order: %s vs %s",
			d1.Key(),
			d2.Key(),
		)
	}
}

func TestDirectoryInitializer_ImplementsInterface(
	_ *testing.T,
) {
	// Compile-time check is in directory.go, but this is a runtime verification
	var _ providers.Initializer = (*DirectoryInitializer)(nil)
}
