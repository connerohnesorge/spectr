package providers

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/spf13/afero"
)

// mockSkillTemplateManager implements TemplateManager with skill support for testing
type mockSkillTemplateManager struct {
	skillFS fs.FS
	err     error
}

func (*mockSkillTemplateManager) InstructionPointer() domain.TemplateRef {
	return domain.TemplateRef{Name: "instruction-pointer.md.tmpl"}
}
func (*mockSkillTemplateManager) Agents() domain.TemplateRef {
	return domain.TemplateRef{Name: "AGENTS.md.tmpl"}
}
func (*mockSkillTemplateManager) SlashCommand(cmd domain.SlashCommand) domain.TemplateRef {
	return domain.TemplateRef{Name: fmt.Sprintf("slash-%s.md.tmpl", cmd.String())}
}
func (*mockSkillTemplateManager) TOMLSlashCommand(cmd domain.SlashCommand) domain.TemplateRef {
	return domain.TemplateRef{Name: fmt.Sprintf("slash-%s.toml.tmpl", cmd.String())}
}

func (m *mockSkillTemplateManager) SkillFS(_ string) (fs.FS, error) {
	if m.err != nil {
		return nil, m.err
	}

	return m.skillFS, nil
}

// TestAgentSkillsInitializer_Construction tests that NewAgentSkillsInitializer
// correctly creates an initializer with skill name, target dir, and template manager.
func TestAgentSkillsInitializer_Construction(t *testing.T) {
	tm := &mockSkillTemplateManager{}
	skillName := "test-skill"
	targetDir := ".claude/skills/test-skill"

	init := NewAgentSkillsInitializer(skillName, targetDir, tm)

	if init == nil {
		t.Fatal("NewAgentSkillsInitializer returned nil")
	}

	if init.skillName != skillName {
		t.Errorf("skillName = %v, want %v", init.skillName, skillName)
	}

	if init.targetDir != targetDir {
		t.Errorf("targetDir = %v, want %v", init.targetDir, targetDir)
	}

	// Template manager is stored as an interface, we just verify it was set
	if init.tm == nil {
		t.Error("template manager not set")
	}
}

// TestAgentSkillsInitializer_Init_CopySkillDirectory tests that Init()
// recursively copies all files from embedded skill directory, preserving
// directory structure and file permissions.
func TestAgentSkillsInitializer_Init_CopySkillDirectory(t *testing.T) {
	// Setup: Create a mock skill filesystem
	skillFS := fstest.MapFS{
		"SKILL.md": {
			Data: []byte("# Test Skill\n\nThis is a test skill."),
			Mode: 0o644,
		},
		"scripts/accept.sh": {
			Data: []byte("#!/bin/bash\necho 'test script'"),
			Mode: 0o755, // Executable
		},
		"scripts/helper.sh": {
			Data: []byte("#!/bin/bash\necho 'helper'"),
			Mode: 0o755,
		},
		"references/doc.md": {
			Data: []byte("# Documentation"),
			Mode: 0o644,
		},
	}

	tm := &mockSkillTemplateManager{skillFS: skillFS}
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	init := NewAgentSkillsInitializer("test-skill", ".claude/skills/test-skill", tm)

	// Execute
	result, err := init.Init(context.Background(), projectFs, homeFs, cfg, tm)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Check result - should have created files and directories
	expectedCreated := []string{
		".claude/skills/test-skill",
		".claude/skills/test-skill/SKILL.md",
		".claude/skills/test-skill/scripts",
		".claude/skills/test-skill/scripts/accept.sh",
		".claude/skills/test-skill/scripts/helper.sh",
		".claude/skills/test-skill/references",
		".claude/skills/test-skill/references/doc.md",
	}

	if !containsAll(result.CreatedFiles, expectedCreated) {
		t.Errorf(
			"CreatedFiles missing expected files.\nGot: %v\nWant to contain: %v",
			result.CreatedFiles,
			expectedCreated,
		)
	}

	if len(result.UpdatedFiles) != 0 {
		t.Errorf("UpdatedFiles = %v, want []", result.UpdatedFiles)
	}

	// Verify files exist with correct content
	tests := []struct {
		path            string
		expectedContent string
		shouldBeExec    bool
	}{
		{
			".claude/skills/test-skill/SKILL.md",
			"# Test Skill\n\nThis is a test skill.",
			false,
		},
		{
			".claude/skills/test-skill/scripts/accept.sh",
			"#!/bin/bash\necho 'test script'",
			true,
		},
		{
			".claude/skills/test-skill/scripts/helper.sh",
			"#!/bin/bash\necho 'helper'",
			true,
		},
		{
			".claude/skills/test-skill/references/doc.md",
			"# Documentation",
			false,
		},
	}

	for _, tt := range tests {
		// Check file exists
		exists, err := afero.Exists(projectFs, tt.path)
		if err != nil {
			t.Errorf("failed to check file %s: %v", tt.path, err)

			continue
		}
		if !exists {
			t.Errorf("file %s does not exist after Init()", tt.path)

			continue
		}

		// Check content
		content, err := afero.ReadFile(projectFs, tt.path)
		if err != nil {
			t.Errorf("failed to read file %s: %v", tt.path, err)

			continue
		}
		if string(content) != tt.expectedContent {
			t.Errorf(
				"file %s content = %q, want %q",
				tt.path,
				string(content),
				tt.expectedContent,
			)
		}

		// Check permissions for executable files
		info, err := projectFs.Stat(tt.path)
		if err != nil {
			t.Errorf("failed to stat file %s: %v", tt.path, err)

			continue
		}

		isExec := info.Mode()&0o111 != 0
		if tt.shouldBeExec && !isExec {
			t.Errorf(
				"file %s should be executable, mode = %v",
				tt.path,
				info.Mode(),
			)
		}
	}

	// Verify directory structure is preserved
	dirTests := []string{
		".claude/skills/test-skill",
		".claude/skills/test-skill/scripts",
		".claude/skills/test-skill/references",
	}

	for _, dir := range dirTests {
		exists, err := afero.DirExists(projectFs, dir)
		if err != nil {
			t.Errorf("failed to check directory %s: %v", dir, err)
		}
		if !exists {
			t.Errorf("directory %s does not exist after Init()", dir)
		}
	}
}

// TestAgentSkillsInitializer_Init_SkillNotFound tests that Init() returns
// an error when the skill name doesn't exist.
func TestAgentSkillsInitializer_Init_SkillNotFound(t *testing.T) {
	tm := &mockSkillTemplateManager{
		err: fs.ErrNotExist,
	}
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	init := NewAgentSkillsInitializer(
		"nonexistent-skill",
		".claude/skills/nonexistent-skill",
		tm,
	)

	// Execute
	_, err := init.Init(context.Background(), projectFs, homeFs, cfg, tm)

	// Should return error
	if err == nil {
		t.Fatal("Init() should fail when skill not found")
	}

	if !contains(err.Error(), "nonexistent-skill") {
		t.Errorf(
			"error message should contain skill name, got: %v",
			err.Error(),
		)
	}
}

// TestAgentSkillsInitializer_Init_IdempotentExecution tests that calling
// Init() multiple times overwrites existing files with embedded content.
func TestAgentSkillsInitializer_Init_IdempotentExecution(t *testing.T) {
	// Setup: Create a mock skill filesystem
	skillFS := fstest.MapFS{
		"SKILL.md": {
			Data: []byte("# Updated Skill\n\nThis is the new content."),
			Mode: 0o644,
		},
		"scripts/accept.sh": {
			Data: []byte("#!/bin/bash\necho 'updated script'"),
			Mode: 0o755,
		},
	}

	tm := &mockSkillTemplateManager{skillFS: skillFS}
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	init := NewAgentSkillsInitializer("test-skill", ".claude/skills/test-skill", tm)

	// First execution - creates files
	result1, err := init.Init(context.Background(), projectFs, homeFs, cfg, tm)
	if err != nil {
		t.Fatalf("First Init() failed: %v", err)
	}

	if len(result1.CreatedFiles) == 0 {
		t.Error("First Init() should create files")
	}
	if len(result1.UpdatedFiles) != 0 {
		t.Error("First Init() should not update files")
	}

	// Modify the files to have different content
	err = afero.WriteFile(
		projectFs,
		".claude/skills/test-skill/SKILL.md",
		[]byte("# Old Content\n\nThis should be overwritten."),
		0o644,
	)
	if err != nil {
		t.Fatalf("failed to modify file: %v", err)
	}

	// Second execution - updates files
	result2, err := init.Init(context.Background(), projectFs, homeFs, cfg, tm)
	if err != nil {
		t.Fatalf("Second Init() failed: %v", err)
	}

	// Should have updated files, not created them
	expectedUpdated := []string{
		".claude/skills/test-skill/SKILL.md",
		".claude/skills/test-skill/scripts/accept.sh",
	}

	if !containsAll(result2.UpdatedFiles, expectedUpdated) {
		t.Errorf(
			"Second Init() UpdatedFiles missing expected files.\nGot: %v\nWant to contain: %v",
			result2.UpdatedFiles,
			expectedUpdated,
		)
	}

	// Verify content was overwritten
	content, err := afero.ReadFile(
		projectFs,
		".claude/skills/test-skill/SKILL.md",
	)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	expected := "# Updated Skill\n\nThis is the new content."
	if string(content) != expected {
		t.Errorf(
			"file content after second Init() = %q, want %q",
			string(content),
			expected,
		)
	}
}

// TestAgentSkillsInitializer_Init_CreatesTargetDirectory tests that Init()
// creates the target directory if it doesn't exist.
func TestAgentSkillsInitializer_Init_CreatesTargetDirectory(t *testing.T) {
	skillFS := fstest.MapFS{
		"SKILL.md": {
			Data: []byte("# Test Skill"),
			Mode: 0o644,
		},
	}

	tm := &mockSkillTemplateManager{skillFS: skillFS}
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	// Verify target directory doesn't exist
	exists, _ := afero.DirExists(projectFs, ".claude/skills/test-skill")
	if exists {
		t.Fatal("target directory should not exist before Init()")
	}

	init := NewAgentSkillsInitializer("test-skill", ".claude/skills/test-skill", tm)

	// Execute
	_, err := init.Init(context.Background(), projectFs, homeFs, cfg, tm)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Verify target directory was created
	exists, err = afero.DirExists(projectFs, ".claude/skills/test-skill")
	if err != nil {
		t.Errorf("failed to check directory: %v", err)
	}
	if !exists {
		t.Error("target directory should exist after Init()")
	}

	// Verify parent directories were created too
	exists, err = afero.DirExists(projectFs, ".claude/skills")
	if err != nil {
		t.Errorf("failed to check parent directory: %v", err)
	}
	if !exists {
		t.Error("parent directory should exist after Init()")
	}
}

// TestAgentSkillsInitializer_IsSetup tests that IsSetup() returns true
// if SKILL.md exists and false if it's missing.
func TestAgentSkillsInitializer_IsSetup(t *testing.T) {
	tests := []struct {
		name       string
		setupFiles map[string]string
		want       bool
	}{
		{
			name: "returns true when SKILL.md exists",
			setupFiles: map[string]string{
				".claude/skills/test-skill/SKILL.md": "# Test Skill",
			},
			want: true,
		},
		{
			name:       "returns false when SKILL.md is missing",
			setupFiles: make(map[string]string),
			want:       false,
		},
		{
			name: "returns false when only other files exist",
			setupFiles: map[string]string{
				".claude/skills/test-skill/scripts/accept.sh": "#!/bin/bash",
				".claude/skills/test-skill/README.md":         "# README",
			},
			want: false,
		},
		{
			name: "returns true when SKILL.md exists with other files",
			setupFiles: map[string]string{
				".claude/skills/test-skill/SKILL.md":          "# Test Skill",
				".claude/skills/test-skill/scripts/accept.sh": "#!/bin/bash",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			projectFs := afero.NewMemMapFs()
			homeFs := afero.NewMemMapFs()
			cfg := &Config{SpectrDir: "spectr"}

			// Create setup files
			for path, content := range tt.setupFiles {
				dir := filepath.Dir(path)
				if err := projectFs.MkdirAll(dir, 0o755); err != nil {
					t.Fatalf("failed to create directory %s: %v", dir, err)
				}
				if err := afero.WriteFile(projectFs, path, []byte(content), 0o644); err != nil {
					t.Fatalf("failed to create file %s: %v", path, err)
				}
			}

			tm := &mockSkillTemplateManager{}
			init := NewAgentSkillsInitializer(
				"test-skill",
				".claude/skills/test-skill",
				tm,
			)

			// Execute
			got := init.IsSetup(projectFs, homeFs, cfg)

			// Check result
			if got != tt.want {
				t.Errorf("IsSetup() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestAgentSkillsInitializer_dedupeKey tests that dedupeKey() returns
// "AgentSkillsInitializer:<normalized-target-dir>".
func TestAgentSkillsInitializer_dedupeKey(t *testing.T) {
	tests := []struct {
		name      string
		targetDir string
		want      string
	}{
		{
			name:      "simple path",
			targetDir: ".claude/skills/test-skill",
			want:      "AgentSkillsInitializer:.claude/skills/test-skill",
		},
		{
			name:      "path with trailing slash",
			targetDir: ".claude/skills/test-skill/",
			want:      "AgentSkillsInitializer:.claude/skills/test-skill",
		},
		{
			name:      "path with dots",
			targetDir: ".claude/skills/./test-skill",
			want:      "AgentSkillsInitializer:.claude/skills/test-skill",
		},
		{
			name:      "path with double slashes",
			targetDir: ".claude//skills//test-skill",
			want:      "AgentSkillsInitializer:.claude/skills/test-skill",
		},
		{
			name:      "complex path",
			targetDir: ".claude/skills/../skills/test-skill",
			want:      "AgentSkillsInitializer:.claude/skills/test-skill",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := &mockSkillTemplateManager{}
			init := &AgentSkillsInitializer{
				skillName: "test-skill",
				targetDir: tt.targetDir,
				tm:        tm,
			}
			got := init.dedupeKey()
			if got != tt.want {
				t.Errorf("dedupeKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestAgentSkillsInitializer_PreservesExecutablePermissions tests that
// executable scripts maintain their executable permissions after copying.
func TestAgentSkillsInitializer_PreservesExecutablePermissions(t *testing.T) {
	// Setup: Create a skill with various file permissions
	skillFS := fstest.MapFS{
		"SKILL.md": {
			Data: []byte("# Test Skill"),
			Mode: 0o644, // Regular file
		},
		"scripts/executable.sh": {
			Data: []byte("#!/bin/bash\necho 'executable'"),
			Mode: 0o755, // Executable
		},
		"scripts/readonly.sh": {
			Data: []byte("#!/bin/bash\necho 'readonly'"),
			Mode: 0o444, // Read-only
		},
		"config/settings.json": {
			Data: []byte("{}"),
			Mode: 0o644, // Regular file
		},
	}

	tm := &mockSkillTemplateManager{skillFS: skillFS}
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	init := NewAgentSkillsInitializer("test-skill", ".claude/skills/test-skill", tm)

	// Execute
	_, err := init.Init(context.Background(), projectFs, homeFs, cfg, tm)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Verify permissions
	tests := []struct {
		path         string
		expectedMode os.FileMode
	}{
		{".claude/skills/test-skill/SKILL.md", 0o644},
		{".claude/skills/test-skill/scripts/executable.sh", 0o755},
		{".claude/skills/test-skill/scripts/readonly.sh", 0o444},
		{".claude/skills/test-skill/config/settings.json", 0o644},
	}

	for _, tt := range tests {
		info, err := projectFs.Stat(tt.path)
		if err != nil {
			t.Errorf("failed to stat file %s: %v", tt.path, err)

			continue
		}

		if info.Mode() != tt.expectedMode {
			t.Errorf(
				"file %s mode = %v, want %v",
				tt.path,
				info.Mode(),
				tt.expectedMode,
			)
		}
	}
}

// TestAgentSkillsInitializer_EmptySkillDirectory tests handling of
// an empty skill directory (only root, no files).
func TestAgentSkillsInitializer_EmptySkillDirectory(t *testing.T) {
	// Create an empty skill filesystem
	skillFS := fstest.MapFS{}

	tm := &mockSkillTemplateManager{skillFS: skillFS}
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	init := NewAgentSkillsInitializer("empty-skill", ".claude/skills/empty-skill", tm)

	// Execute
	result, err := init.Init(context.Background(), projectFs, homeFs, cfg, tm)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Should create the root directory
	if !containsString(result.CreatedFiles, ".claude/skills/empty-skill") {
		t.Errorf(
			"CreatedFiles should contain root directory, got: %v",
			result.CreatedFiles,
		)
	}

	// Verify directory exists
	exists, err := afero.DirExists(projectFs, ".claude/skills/empty-skill")
	if err != nil {
		t.Errorf("failed to check directory: %v", err)
	}
	if !exists {
		t.Error("target directory should exist after Init()")
	}
}

// Helper functions

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}

	return false
}

// containsString checks if a string slice contains a specific string
func containsString(list []string, item string) bool {
	for _, s := range list {
		if s == item {
			return true
		}
	}

	return false
}

// containsAll checks if actualList contains all items from expectedList
func containsAll(actualList, expectedList []string) bool {
	actualMap := make(map[string]bool)
	for _, item := range actualList {
		actualMap[item] = true
	}

	for _, expected := range expectedList {
		if !actualMap[expected] {
			return false
		}
	}

	return true
}
