package initializers

import (
	"context"
	"testing"

	"github.com/connerohnesorge/spectr/internal/initialize/providers"
	"github.com/spf13/afero"
)

// mockTemplateManager is a mock implementation of providers.TemplateManager
// used for testing initializers that don't actually need template rendering.
type mockTemplateManager struct {
	instructionPointer string
	slashCommands      map[string]string
	err                error
}

func newMockTemplateManager() *mockTemplateManager {
	return &mockTemplateManager{
		instructionPointer: "# Spectr Instructions\n\nDefault mock content.",
		slashCommands: map[string]string{
			"proposal": "Proposal command content",
			"apply":    "Apply command content",
		},
	}
}

func (m *mockTemplateManager) RenderInstructionPointer(_ providers.TemplateContext) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.instructionPointer, nil
}

func (m *mockTemplateManager) RenderSlashCommand(commandType string, _ providers.TemplateContext) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	if content, ok := m.slashCommands[commandType]; ok {
		return content, nil
	}
	return "", nil
}

func TestDirectoryInitializer_Init(t *testing.T) {
	tests := []struct {
		name    string
		paths   []string
		global  bool
		wantErr bool
	}{
		{
			name:    "creates single directory",
			paths:   []string{".claude/commands/spectr"},
			global:  false,
			wantErr: false,
		},
		{
			name:    "creates multiple directories",
			paths:   []string{".claude/commands/spectr", ".claude/settings"},
			global:  false,
			wantErr: false,
		},
		{
			name:    "creates nested directories",
			paths:   []string{"deeply/nested/path/structure"},
			global:  false,
			wantErr: false,
		},
		{
			name:    "creates global directory",
			paths:   []string{".config/aider/commands/spectr"},
			global:  true,
			wantErr: false,
		},
		{
			name:    "handles empty paths",
			paths:   []string{},
			global:  false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			cfg := &providers.Config{SpectrDir: "spectr"}
			tm := newMockTemplateManager()
			ctx := context.Background()

			var init *DirectoryInitializer
			if tt.global {
				init = NewGlobalDirectoryInitializer(tt.paths...)
			} else {
				init = NewDirectoryInitializer(tt.paths...)
			}

			err := init.Init(ctx, fs, cfg, tm)
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verify directories were created
			for _, p := range tt.paths {
				exists, err := afero.DirExists(fs, p)
				if err != nil {
					t.Errorf("DirExists() error = %v", err)
					continue
				}
				if !exists {
					t.Errorf("Directory %s was not created", p)
				}
			}
		})
	}
}

func TestDirectoryInitializer_Init_Idempotent(t *testing.T) {
	fs := afero.NewMemMapFs()
	cfg := &providers.Config{SpectrDir: "spectr"}
	tm := newMockTemplateManager()
	ctx := context.Background()

	init := NewDirectoryInitializer(".claude/commands/spectr")

	// Run Init multiple times
	for i := 0; i < 3; i++ {
		err := init.Init(ctx, fs, cfg, tm)
		if err != nil {
			t.Errorf("Init() run %d: error = %v", i+1, err)
		}
	}

	// Verify directory exists
	exists, err := afero.DirExists(fs, ".claude/commands/spectr")
	if err != nil {
		t.Errorf("DirExists() error = %v", err)
	}
	if !exists {
		t.Error("Directory should exist after multiple Init() calls")
	}
}

func TestDirectoryInitializer_Init_CreatesParentDirectories(t *testing.T) {
	fs := afero.NewMemMapFs()
	cfg := &providers.Config{SpectrDir: "spectr"}
	tm := newMockTemplateManager()
	ctx := context.Background()

	init := NewDirectoryInitializer("a/b/c/d/e")

	err := init.Init(ctx, fs, cfg, tm)
	if err != nil {
		t.Errorf("Init() error = %v", err)
		return
	}

	// Verify all parent directories were created
	dirs := []string{"a", "a/b", "a/b/c", "a/b/c/d", "a/b/c/d/e"}
	for _, dir := range dirs {
		exists, err := afero.DirExists(fs, dir)
		if err != nil {
			t.Errorf("DirExists(%s) error = %v", dir, err)
			continue
		}
		if !exists {
			t.Errorf("Parent directory %s was not created", dir)
		}
	}
}

func TestDirectoryInitializer_IsSetup(t *testing.T) {
	tests := []struct {
		name       string
		paths      []string
		setupPaths []string // Paths to pre-create
		want       bool
	}{
		{
			name:       "returns false when directories do not exist",
			paths:      []string{".claude/commands/spectr"},
			setupPaths: []string{},
			want:       false,
		},
		{
			name:       "returns true when all directories exist",
			paths:      []string{".claude/commands/spectr", ".claude/settings"},
			setupPaths: []string{".claude/commands/spectr", ".claude/settings"},
			want:       true,
		},
		{
			name:       "returns false when only some directories exist",
			paths:      []string{".claude/commands/spectr", ".claude/settings"},
			setupPaths: []string{".claude/commands/spectr"},
			want:       false,
		},
		{
			name:       "returns true for empty paths",
			paths:      []string{},
			setupPaths: []string{},
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			cfg := &providers.Config{SpectrDir: "spectr"}

			// Pre-create directories
			for _, p := range tt.setupPaths {
				if err := fs.MkdirAll(p, 0755); err != nil {
					t.Fatalf("Failed to pre-create directory %s: %v", p, err)
				}
			}

			init := NewDirectoryInitializer(tt.paths...)

			got := init.IsSetup(fs, cfg)
			if got != tt.want {
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
			name:  "returns first path",
			paths: []string{".claude/commands/spectr", ".claude/settings"},
			want:  ".claude/commands/spectr",
		},
		{
			name:  "returns single path",
			paths: []string{".claude/commands/spectr"},
			want:  ".claude/commands/spectr",
		},
		{
			name:  "returns empty string for empty paths",
			paths: []string{},
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			init := NewDirectoryInitializer(tt.paths...)

			got := init.Path()
			if got != tt.want {
				t.Errorf("Path() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirectoryInitializer_IsGlobal(t *testing.T) {
	tests := []struct {
		name   string
		global bool
		want   bool
	}{
		{
			name:   "project initializer returns false",
			global: false,
			want:   false,
		},
		{
			name:   "global initializer returns true",
			global: true,
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var init *DirectoryInitializer
			if tt.global {
				init = NewGlobalDirectoryInitializer(".config/aider/commands")
			} else {
				init = NewDirectoryInitializer(".claude/commands")
			}

			got := init.IsGlobal()
			if got != tt.want {
				t.Errorf("IsGlobal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirectoryInitializer_Paths(t *testing.T) {
	tests := []struct {
		name  string
		paths []string
		want  []string
	}{
		{
			name:  "returns all paths",
			paths: []string{".claude/commands/spectr", ".claude/settings"},
			want:  []string{".claude/commands/spectr", ".claude/settings"},
		},
		{
			name:  "returns single path",
			paths: []string{".claude/commands/spectr"},
			want:  []string{".claude/commands/spectr"},
		},
		{
			name:  "returns empty slice for no paths",
			paths: []string{},
			want:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			init := NewDirectoryInitializer(tt.paths...)

			got := init.Paths()
			if len(got) != len(tt.want) {
				t.Errorf("Paths() returned %d paths, want %d", len(got), len(tt.want))
				return
			}
			for i, p := range got {
				if p != tt.want[i] {
					t.Errorf("Paths()[%d] = %v, want %v", i, p, tt.want[i])
				}
			}
		})
	}
}
