package providers

import (
	"context"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/spf13/afero"
)

const (
	testSkillPath = ".agents/skills/test-skill/SKILL.md"
)

// TestSkillFileInitializer_Construction tests that NewSkillFileInitializer
// correctly creates an initializer with target path and template.
func TestSkillFileInitializer_Construction(
	t *testing.T,
) {
	targetPath := testProposalSkillPath
	tmpl := createTestTemplate(t, "test content")

	init := NewSkillFileInitializer(targetPath, tmpl)

	if init == nil {
		t.Fatal("NewSkillFileInitializer returned nil")
	}

	if init.targetPath != targetPath {
		t.Errorf(
			"targetPath = %v, want %v",
			init.targetPath,
			targetPath,
		)
	}

	if init.template.Name != tmpl.Name {
		t.Errorf(
			"template.Name = %v, want %v",
			init.template.Name,
			tmpl.Name,
		)
	}
}

// TestSkillFileInitializer_Init_NewFile tests that Init() creates a new
// SKILL.md file with rendered template content and proper frontmatter.
func TestSkillFileInitializer_Init_NewFile(
	t *testing.T,
) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	// Create template with YAML frontmatter
	tmplText := `---
name: spectr-proposal
description: Create a new change proposal
---

# Proposal Creation Guide

Specs: {{.SpecsDir}}
Changes: {{.ChangesDir}}`

	tmpl := createTestTemplate(t, tmplText)
	targetPath := testProposalSkillPath

	init := NewSkillFileInitializer(targetPath, tmpl)

	// Execute
	result, err := init.Init(
		context.Background(),
		projectFs,
		homeFs,
		cfg,
		nil,
	)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Check result
	if len(result.CreatedFiles) != 1 || result.CreatedFiles[0] != targetPath {
		t.Errorf(
			"CreatedFiles = %v, want [%s]",
			result.CreatedFiles,
			targetPath,
		)
	}
	if len(result.UpdatedFiles) != 0 {
		t.Errorf(
			"UpdatedFiles = %v, want []",
			result.UpdatedFiles,
		)
	}

	// Verify file exists
	exists, err := afero.Exists(projectFs, targetPath)
	if err != nil {
		t.Fatalf("failed to check file: %v", err)
	}
	if !exists {
		t.Error("SKILL.md file should exist after Init()")
	}

	// Verify content
	content, err := afero.ReadFile(projectFs, targetPath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	expected := `---
name: spectr-proposal
description: Create a new change proposal
---

# Proposal Creation Guide

Specs: spectr/specs
Changes: spectr/changes`

	if string(content) != expected {
		t.Errorf(
			"file content = %q, want %q",
			string(content),
			expected,
		)
	}
}

// TestSkillFileInitializer_Init_CreatesParentDirectory tests that Init()
// creates the parent directory if it doesn't exist.
func TestSkillFileInitializer_Init_CreatesParentDirectory(
	t *testing.T,
) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	tmpl := createTestTemplate(t, "test content")
	targetPath := testProposalSkillPath

	// Verify parent directory doesn't exist
	parentDir := filepath.Dir(targetPath)
	exists, _ := afero.DirExists(projectFs, parentDir)
	if exists {
		t.Fatal("parent directory should not exist before Init()")
	}

	init := NewSkillFileInitializer(targetPath, tmpl)

	// Execute
	_, err := init.Init(
		context.Background(),
		projectFs,
		homeFs,
		cfg,
		nil,
	)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Verify parent directory was created
	exists, err = afero.DirExists(projectFs, parentDir)
	if err != nil {
		t.Errorf("failed to check directory: %v", err)
	}
	if !exists {
		t.Error("parent directory should exist after Init()")
	}

	// Verify all ancestor directories were created
	exists, err = afero.DirExists(projectFs, ".agents/skills")
	if err != nil {
		t.Errorf("failed to check ancestor directory: %v", err)
	}
	if !exists {
		t.Error("ancestor directory .agents/skills should exist after Init()")
	}
}

// TestSkillFileInitializer_Init_UpdateExistingFile tests that Init()
// overwrites an existing SKILL.md file (idempotent operation).
func TestSkillFileInitializer_Init_UpdateExistingFile(
	t *testing.T,
) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	targetPath := testProposalSkillPath

	// Create existing file with old content
	parentDir := filepath.Dir(targetPath)
	_ = projectFs.MkdirAll(parentDir, 0o755)
	_ = afero.WriteFile(
		projectFs,
		targetPath,
		[]byte("old content"),
		0o644,
	)

	// Create template with new content
	tmpl := createTestTemplate(t, "new content")
	init := NewSkillFileInitializer(targetPath, tmpl)

	// Execute
	result, err := init.Init(
		context.Background(),
		projectFs,
		homeFs,
		cfg,
		nil,
	)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Check result - should be UpdatedFiles, not CreatedFiles
	if len(result.UpdatedFiles) != 1 || result.UpdatedFiles[0] != targetPath {
		t.Errorf(
			"UpdatedFiles = %v, want [%s]",
			result.UpdatedFiles,
			targetPath,
		)
	}
	if len(result.CreatedFiles) != 0 {
		t.Errorf(
			"CreatedFiles = %v, want []",
			result.CreatedFiles,
		)
	}

	// Verify content was updated
	content, err := afero.ReadFile(projectFs, targetPath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	if string(content) != "new content" {
		t.Errorf(
			"file content = %q, want %q",
			string(content),
			"new content",
		)
	}
}

// TestSkillFileInitializer_Init_TemplateContextUsage tests that Init()
// correctly renders templates with TemplateContext variables.
func TestSkillFileInitializer_Init_TemplateContextUsage(
	t *testing.T,
) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "myspectr"}

	// Create template that uses context variables
	tmplText := `Base: {{.BaseDir}}
Specs: {{.SpecsDir}}
Changes: {{.ChangesDir}}
Project: {{.ProjectFile}}
Agents: {{.AgentsFile}}`

	tmpl := createTestTemplate(t, tmplText)
	targetPath := testSkillPath

	init := NewSkillFileInitializer(targetPath, tmpl)

	// Execute
	_, err := init.Init(
		context.Background(),
		projectFs,
		homeFs,
		cfg,
		nil,
	)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Verify template was rendered with correct context
	content, err := afero.ReadFile(projectFs, targetPath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	expected := `Base: myspectr
Specs: myspectr/specs
Changes: myspectr/changes
Project: myspectr/project.md
Agents: myspectr/AGENTS.md`

	if string(content) != expected {
		t.Errorf(
			"file content = %q, want %q",
			string(content),
			expected,
		)
	}
}

// TestSkillFileInitializer_Init_ErrorTemplateRenderFailure tests that Init()
// returns an error when template rendering fails.
func TestSkillFileInitializer_Init_ErrorTemplateRenderFailure(
	t *testing.T,
) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	// Create template with invalid syntax
	tmplText := "{{.NonExistentField}}"
	tmpl, err := template.New("test").Parse(tmplText)
	if err != nil {
		t.Fatalf("failed to create test template: %v", err)
	}

	templateRef := domain.TemplateRef{
		Name:     "test",
		Template: tmpl,
	}

	targetPath := testSkillPath
	init := NewSkillFileInitializer(targetPath, templateRef)

	// Execute - should fail during template rendering
	_, err = init.Init(
		context.Background(),
		projectFs,
		homeFs,
		cfg,
		nil,
	)

	// Should return error
	if err == nil {
		t.Fatal("Init() should fail when template rendering fails")
	}

	if !contains(err.Error(), "failed to render template") {
		t.Errorf(
			"error message should contain 'failed to render template', got: %v",
			err.Error(),
		)
	}
}

// TestSkillFileInitializer_Init_IdempotentExecution tests that calling
// Init() multiple times produces consistent results.
func TestSkillFileInitializer_Init_IdempotentExecution(
	t *testing.T,
) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	tmpl := createTestTemplate(t, "skill content")
	targetPath := testSkillPath

	init := NewSkillFileInitializer(targetPath, tmpl)

	// First execution - creates file
	result1, err := init.Init(
		context.Background(),
		projectFs,
		homeFs,
		cfg,
		nil,
	)
	if err != nil {
		t.Fatalf("First Init() failed: %v", err)
	}

	if len(result1.CreatedFiles) != 1 {
		t.Error("First Init() should create file")
	}

	// Second execution - updates file
	result2, err := init.Init(
		context.Background(),
		projectFs,
		homeFs,
		cfg,
		nil,
	)
	if err != nil {
		t.Fatalf("Second Init() failed: %v", err)
	}

	if len(result2.UpdatedFiles) != 1 {
		t.Error("Second Init() should update file")
	}

	// Verify content is identical
	content, err := afero.ReadFile(projectFs, targetPath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	if string(content) != "skill content" {
		t.Errorf(
			"file content = %q, want %q",
			string(content),
			"skill content",
		)
	}
}

// TestSkillFileInitializer_IsSetup tests that IsSetup() returns true
// if SKILL.md exists and false if it's missing.
func TestSkillFileInitializer_IsSetup(
	t *testing.T,
) {
	tests := []struct {
		name       string
		targetPath string
		setupFile  bool
		want       bool
	}{
		{
			name:       "returns true when SKILL.md exists",
			targetPath: "testSkillPath",
			setupFile:  true,
			want:       true,
		},
		{
			name:       "returns false when SKILL.md is missing",
			targetPath: "testSkillPath",
			setupFile:  false,
			want:       false,
		},
		{
			name:       "returns true for different skill directory",
			targetPath: ".agents/skills/other-skill/SKILL.md",
			setupFile:  true,
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			projectFs := afero.NewMemMapFs()
			homeFs := afero.NewMemMapFs()
			cfg := &Config{SpectrDir: "spectr"}

			if tt.setupFile {
				parentDir := filepath.Dir(tt.targetPath)
				_ = projectFs.MkdirAll(parentDir, 0o755)
				_ = afero.WriteFile(
					projectFs,
					tt.targetPath,
					[]byte("test content"),
					0o644,
				)
			}

			tmpl := createTestTemplate(t, "content")
			init := NewSkillFileInitializer(tt.targetPath, tmpl)

			// Execute
			got := init.IsSetup(projectFs, homeFs, cfg)

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

// TestSkillFileInitializer_dedupeKey tests that dedupeKey() returns
// "SkillFileInitializer:<normalized-target-path>".
func TestSkillFileInitializer_dedupeKey(
	t *testing.T,
) {
	tests := []struct {
		name       string
		targetPath string
		want       string
	}{
		{
			name:       "simple path",
			targetPath: "testSkillPath",
			want:       "SkillFileInitializer:testSkillPath",
		},
		{
			name:       "path with trailing slash",
			targetPath: "testSkillPath/",
			want:       "SkillFileInitializer:testSkillPath",
		},
		{
			name:       "path with dots",
			targetPath: ".agents/skills/./test-skill/SKILL.md",
			want:       "SkillFileInitializer:.agents/skills/test-skill/SKILL.md",
		},
		{
			name:       "path with double slashes",
			targetPath: ".agents//skills//test-skill//SKILL.md",
			want:       "SkillFileInitializer:.agents/skills/test-skill/SKILL.md",
		},
		{
			name:       "complex path",
			targetPath: ".agents/skills/../skills/test-skill/SKILL.md",
			want:       "SkillFileInitializer:.agents/skills/test-skill/SKILL.md",
		},
		{
			name:       "different skill",
			targetPath: testProposalSkillPath,
			want:       "SkillFileInitializer:.agents/skills/spectr-proposal/SKILL.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl := createTestTemplate(t, "content")
			init := &SkillFileInitializer{
				targetPath: tt.targetPath,
				template:   tmpl,
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

// TestSkillFileInitializer_MultipleSkills tests that multiple
// SkillFileInitializers can coexist without conflicts.
func TestSkillFileInitializer_MultipleSkills(
	t *testing.T,
) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	// Create two different skill initializers
	tmpl1 := createTestTemplate(t, "proposal skill content")
	tmpl2 := createTestTemplate(t, "apply skill content")

	init1 := NewSkillFileInitializer(
		"testProposalSkillPath",
		tmpl1,
	)
	init2 := NewSkillFileInitializer(
		".agents/skills/spectr-apply/SKILL.md",
		tmpl2,
	)

	// Execute both
	_, err := init1.Init(
		context.Background(),
		projectFs,
		homeFs,
		cfg,
		nil,
	)
	if err != nil {
		t.Fatalf("init1.Init() failed: %v", err)
	}

	_, err = init2.Init(
		context.Background(),
		projectFs,
		homeFs,
		cfg,
		nil,
	)
	if err != nil {
		t.Fatalf("init2.Init() failed: %v", err)
	}

	// Verify both files exist with correct content
	content1, err := afero.ReadFile(
		projectFs,
		"testProposalSkillPath",
	)
	if err != nil {
		t.Fatalf("failed to read proposal skill: %v", err)
	}
	if string(content1) != "proposal skill content" {
		t.Errorf(
			"proposal skill content = %q, want %q",
			string(content1),
			"proposal skill content",
		)
	}

	content2, err := afero.ReadFile(
		projectFs,
		".agents/skills/spectr-apply/SKILL.md",
	)
	if err != nil {
		t.Fatalf("failed to read apply skill: %v", err)
	}
	if string(content2) != "apply skill content" {
		t.Errorf(
			"apply skill content = %q, want %q",
			string(content2),
			"apply skill content",
		)
	}

	// Verify both skills report as setup
	if !init1.IsSetup(projectFs, homeFs, cfg) {
		t.Error("init1 should report as setup")
	}
	if !init2.IsSetup(projectFs, homeFs, cfg) {
		t.Error("init2 should report as setup")
	}

	// Verify dedupe keys are different
	if init1.dedupeKey() == init2.dedupeKey() {
		t.Error("dedupe keys should be different for different skills")
	}
}
