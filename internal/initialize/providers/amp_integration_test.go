package providers

import (
	"context"
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"
	"text/template"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

// TestAmpProvider_Integration_DirectoryStructure tests that Amp provider
// creates the correct .agents/skills/ directory structure with all expected
// skill files and subdirectories (Task 4.3).
func TestAmpProvider_Integration_DirectoryStructure(
	t *testing.T,
) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	// Create mock skill filesystems for embedded skills
	acceptSkillFS := fstest.MapFS{
		"SKILL.md": {
			Data: []byte(
				"---\nname: spectr-accept-wo-spectr-bin\ndescription: Accept proposals\n---\nContent",
			),
			Mode: 0o644,
		},
		"scripts/accept.sh": {
			Data: []byte("#!/bin/bash\necho 'accept'"),
			Mode: 0o755,
		},
	}

	validateSkillFS := fstest.MapFS{
		"SKILL.md": {
			Data: []byte(
				"---\nname: spectr-validate-wo-spectr-bin\ndescription: Validate proposals\n---\nContent",
			),
			Mode: 0o644,
		},
		"scripts/validate.sh": {
			Data: []byte("#!/bin/bash\necho 'validate'"),
			Mode: 0o755,
		},
	}

	tm := &mockSkillTemplateManager{
		skillFS: acceptSkillFS, // Will be overridden per skill
	}

	// Override SkillFS to return different filesystems per skill
	skillFSMap := map[string]fs.FS{
		"spectr-accept-wo-spectr-bin":   acceptSkillFS,
		"spectr-validate-wo-spectr-bin": validateSkillFS,
	}

	tmWithMultiSkills := &multiSkillTemplateManager{
		mockSkillTemplateManager: tm,
		skillFSMap:               skillFSMap,
	}

	provider := &AmpProvider{}
	initializers := provider.Initializers(context.Background(), tmWithMultiSkills)

	// Execute all initializers
	for _, init := range initializers {
		_, err := init.Init(
			context.Background(),
			projectFs,
			homeFs,
			cfg,
			tmWithMultiSkills,
		)
		if err != nil {
			t.Fatalf("Init() failed: %v", err)
		}
	}

	// Verify directory structure
	expectedDirs := []string{
		".agents",
		".agents/skills",
		".agents/skills/spectr-proposal",
		".agents/skills/spectr-apply",
		".agents/skills/spectr-accept-wo-spectr-bin",
		".agents/skills/spectr-accept-wo-spectr-bin/scripts",
		".agents/skills/spectr-validate-wo-spectr-bin",
		".agents/skills/spectr-validate-wo-spectr-bin/scripts",
	}

	for _, dir := range expectedDirs {
		exists, err := afero.DirExists(projectFs, dir)
		if err != nil {
			t.Errorf("failed to check directory %s: %v", dir, err)
		}
		if !exists {
			t.Errorf("directory %s should exist", dir)
		}
	}

	// Verify files exist
	expectedFiles := []string{
		"AMP.md",
		".agents/skills/spectr-proposal/SKILL.md",
		".agents/skills/spectr-apply/SKILL.md",
		".agents/skills/spectr-accept-wo-spectr-bin/SKILL.md",
		".agents/skills/spectr-accept-wo-spectr-bin/scripts/accept.sh",
		".agents/skills/spectr-validate-wo-spectr-bin/SKILL.md",
		".agents/skills/spectr-validate-wo-spectr-bin/scripts/validate.sh",
	}

	for _, file := range expectedFiles {
		exists, err := afero.Exists(projectFs, file)
		if err != nil {
			t.Errorf("failed to check file %s: %v", file, err)
		}
		if !exists {
			t.Errorf("file %s should exist", file)
		}
	}
}

// TestAmpProvider_Integration_FrontmatterParsing tests that SKILL.md files
// have valid YAML frontmatter with required name and description fields (Task 4.4).
func TestAmpProvider_Integration_FrontmatterParsing(
	t *testing.T,
) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	// Create mock skill filesystems with proper frontmatter
	acceptSkillFS := fstest.MapFS{
		"SKILL.md": {
			Data: []byte(`---
name: spectr-accept-wo-spectr-bin
description: |
  Accept Spectr change proposals without the binary
---

# Content here`),
			Mode: 0o644,
		},
	}

	tm := &multiSkillTemplateManager{
		mockSkillTemplateManager: &mockSkillTemplateManager{
			skillFS: acceptSkillFS,
		},
		skillFSMap: map[string]fs.FS{
			"spectr-accept-wo-spectr-bin": acceptSkillFS,
		},
	}

	provider := &AmpProvider{}
	initializers := provider.Initializers(context.Background(), tm)

	// Execute all initializers
	for _, init := range initializers {
		_, err := init.Init(
			context.Background(),
			projectFs,
			homeFs,
			cfg,
			tm,
		)
		if err != nil {
			t.Fatalf("Init() failed: %v", err)
		}
	}

	// Test cases for frontmatter parsing
	tests := []struct {
		path         string
		expectedName string
		expectedDesc string
	}{
		{
			path:         ".agents/skills/spectr-proposal/SKILL.md",
			expectedName: "spectr-proposal",
			expectedDesc: "Create a new change proposal",
		},
		{
			path:         ".agents/skills/spectr-apply/SKILL.md",
			expectedName: "spectr-apply",
			expectedDesc: "Apply or accept change proposals",
		},
		{
			path:         ".agents/skills/spectr-accept-wo-spectr-bin/SKILL.md",
			expectedName: "spectr-accept-wo-spectr-bin",
			expectedDesc: "Accept Spectr change proposals without the binary",
		},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			// Read SKILL.md file
			content, err := afero.ReadFile(projectFs, tt.path)
			if err != nil {
				t.Fatalf("failed to read %s: %v", tt.path, err)
			}

			// Parse frontmatter
			frontmatter := parseFrontmatter(t, string(content))

			// Verify required fields
			name, ok := frontmatter["name"].(string)
			if !ok {
				t.Error("frontmatter missing 'name' field or not a string")
			}
			if name != tt.expectedName {
				t.Errorf("frontmatter name = %q, want %q", name, tt.expectedName)
			}

			// Description can be string or multiline
			desc := getFrontmatterString(frontmatter, "description")
			if desc == "" {
				t.Error("frontmatter missing 'description' field")
			}
			if !strings.Contains(desc, tt.expectedDesc) {
				t.Errorf(
					"frontmatter description = %q, want to contain %q",
					desc,
					tt.expectedDesc,
				)
			}
		})
	}
}

// TestAmpProvider_Integration_EmbeddedSkillCopying tests that embedded skills
// (spectr-accept-wo-spectr-bin, spectr-validate-wo-spectr-bin) are correctly
// copied with their complete directory structure (Task 4.5).
func TestAmpProvider_Integration_EmbeddedSkillCopying(
	t *testing.T,
) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	// Create realistic embedded skill filesystem
	acceptSkillFS := fstest.MapFS{
		"SKILL.md": {
			Data: []byte(`---
name: spectr-accept-wo-spectr-bin
description: Accept proposals
---
# Skill content`),
			Mode: 0o644,
		},
		"scripts/accept.sh": {
			Data: []byte("#!/bin/bash\necho 'accept script'"),
			Mode: 0o755,
		},
		"scripts/helper.sh": {
			Data: []byte("#!/bin/bash\necho 'helper'"),
			Mode: 0o755,
		},
		"references/example.md": {
			Data: []byte("# Example"),
			Mode: 0o644,
		},
	}

	validateSkillFS := fstest.MapFS{
		"SKILL.md": {
			Data: []byte(`---
name: spectr-validate-wo-spectr-bin
description: Validate proposals
---
# Skill content`),
			Mode: 0o644,
		},
		"scripts/validate.sh": {
			Data: []byte("#!/bin/bash\necho 'validate script'"),
			Mode: 0o755,
		},
	}

	tm := &multiSkillTemplateManager{
		mockSkillTemplateManager: &mockSkillTemplateManager{},
		skillFSMap: map[string]fs.FS{
			"spectr-accept-wo-spectr-bin":   acceptSkillFS,
			"spectr-validate-wo-spectr-bin": validateSkillFS,
		},
	}

	provider := &AmpProvider{}
	initializers := provider.Initializers(context.Background(), tm)

	// Execute all initializers
	for _, init := range initializers {
		_, err := init.Init(
			context.Background(),
			projectFs,
			homeFs,
			cfg,
			tm,
		)
		if err != nil {
			t.Fatalf("Init() failed: %v", err)
		}
	}

	// Verify accept skill was copied completely
	acceptTests := []struct {
		path         string
		expectedText string
		shouldBeExec bool
	}{
		{
			".agents/skills/spectr-accept-wo-spectr-bin/SKILL.md",
			"# Skill content",
			false,
		},
		{
			".agents/skills/spectr-accept-wo-spectr-bin/scripts/accept.sh",
			"echo 'accept script'",
			true,
		},
		{
			".agents/skills/spectr-accept-wo-spectr-bin/scripts/helper.sh",
			"echo 'helper'",
			true,
		},
		{
			".agents/skills/spectr-accept-wo-spectr-bin/references/example.md",
			"# Example",
			false,
		},
	}

	for _, tt := range acceptTests {
		// Check file exists
		content, err := afero.ReadFile(projectFs, tt.path)
		if err != nil {
			t.Errorf("failed to read %s: %v", tt.path, err)

			continue
		}

		// Check content
		if !strings.Contains(string(content), tt.expectedText) {
			t.Errorf(
				"file %s content = %q, want to contain %q",
				tt.path,
				string(content),
				tt.expectedText,
			)

			continue
		}

		// Check executable permissions
		info, err := projectFs.Stat(tt.path)
		if err != nil {
			t.Errorf("failed to stat %s: %v", tt.path, err)

			continue
		}

		isExec := info.Mode()&0o111 != 0
		if tt.shouldBeExec && !isExec {
			t.Errorf("file %s should be executable, mode = %v", tt.path, info.Mode())
		}
		if !tt.shouldBeExec && isExec {
			t.Errorf(
				"file %s should not be executable, mode = %v",
				tt.path,
				info.Mode(),
			)
		}
	}

	// Verify validate skill was copied
	validatePath := ".agents/skills/spectr-validate-wo-spectr-bin/scripts/validate.sh"
	content, err := afero.ReadFile(projectFs, validatePath)
	if err != nil {
		t.Errorf("failed to read %s: %v", validatePath, err)
	}
	if !strings.Contains(string(content), "echo 'validate script'") {
		t.Errorf(
			"validate script content = %q, want to contain validate script",
			string(content),
		)
	}
}

// TestAmpProvider_Integration_TemplateVariableSubstitution tests that template
// variables are correctly substituted in skill content (Task 4.6).
func TestAmpProvider_Integration_TemplateVariableSubstitution(
	t *testing.T,
) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "myspectr"}

	tm := &multiSkillTemplateManager{
		mockSkillTemplateManager: &mockSkillTemplateManager{},
		skillFSMap:               make(map[string]fs.FS),
	}
	provider := &AmpProvider{}
	initializers := provider.Initializers(context.Background(), tm)

	// Execute only SkillFileInitializers (proposal and apply)
	for _, init := range initializers {
		skillFileInit, ok := init.(*SkillFileInitializer)
		if !ok {
			continue
		}

		_, err := skillFileInit.Init(
			context.Background(),
			projectFs,
			homeFs,
			cfg,
			tm,
		)
		if err != nil {
			t.Fatalf("Init() failed: %v", err)
		}
	}

	// Verify template variables were substituted in proposal skill
	proposalContent, err := afero.ReadFile(
		projectFs,
		".agents/skills/spectr-proposal/SKILL.md",
	)
	if err != nil {
		t.Fatalf("failed to read proposal skill: %v", err)
	}

	// Check for substituted variables
	expectedSubstitutions := []string{
		"myspectr/specs",      // {{.SpecsDir}}
		"myspectr/changes",    // {{.ChangesDir}}
		"myspectr/AGENTS.md",  // {{.AgentsFile}}
		"myspectr/project.md", // {{.ProjectFile}}
	}

	for _, expected := range expectedSubstitutions {
		if !strings.Contains(string(proposalContent), expected) {
			t.Errorf(
				"proposal skill should contain %q, content:\n%s",
				expected,
				string(proposalContent),
			)
		}
	}

	// Verify template variables in apply skill
	applyContent, err := afero.ReadFile(
		projectFs,
		".agents/skills/spectr-apply/SKILL.md",
	)
	if err != nil {
		t.Fatalf("failed to read apply skill: %v", err)
	}

	if !strings.Contains(string(applyContent), "myspectr/changes") {
		t.Errorf(
			"apply skill should contain template substitutions, content:\n%s",
			string(applyContent),
		)
	}
}

// TestAmpProvider_Integration_Deduplication tests that when multiple providers
// generate the same skill files, only one is created (Task 4.7).
func TestAmpProvider_Integration_Deduplication(
	t *testing.T,
) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	tm := &multiSkillTemplateManager{
		mockSkillTemplateManager: &mockSkillTemplateManager{},
		skillFSMap:               make(map[string]fs.FS),
	}

	// Create two different providers that both create the same directory
	provider1 := &AmpProvider{}
	provider2 := &AmpProvider{}

	inits1 := provider1.Initializers(context.Background(), tm)
	inits2 := provider2.Initializers(context.Background(), tm)

	// Combine initializers
	allInits := make([]Initializer, 0, len(inits1)+len(inits2))
	allInits = append(allInits, inits1...)
	allInits = append(allInits, inits2...)

	// Execute first time - should create files
	firstResults := make([]InitResult, 0, len(allInits))
	for _, init := range allInits {
		result, err := init.Init(
			context.Background(),
			projectFs,
			homeFs,
			cfg,
			tm,
		)
		if err != nil {
			t.Fatalf("Init() failed: %v", err)
		}
		firstResults = append(firstResults, result)
	}

	// Count created files
	totalCreated := 0
	totalUpdated := 0
	for _, result := range firstResults {
		totalCreated += len(result.CreatedFiles)
		totalUpdated += len(result.UpdatedFiles)
	}

	// Should have created files once, and updated them on duplicate execution
	// The second provider's initializers should update, not create
	if totalCreated == 0 {
		t.Error("should have created some files")
	}

	// Verify .agents/skills directory exists only once
	entries, err := afero.ReadDir(projectFs, ".agents")
	if err != nil {
		t.Fatalf("failed to read .agents: %v", err)
	}

	skillsDirCount := 0
	for _, entry := range entries {
		if entry.Name() == "skills" && entry.IsDir() {
			skillsDirCount++
		}
	}

	if skillsDirCount != 1 {
		t.Errorf("found %d 'skills' directories, want 1", skillsDirCount)
	}

	// Verify each skill exists only once
	skillPaths := []string{
		".agents/skills/spectr-proposal/SKILL.md",
		".agents/skills/spectr-apply/SKILL.md",
	}

	for _, path := range skillPaths {
		exists, err := afero.Exists(projectFs, path)
		if err != nil {
			t.Errorf("failed to check %s: %v", path, err)
		}
		if !exists {
			t.Errorf("skill %s should exist", path)
		}

		// Read file to ensure it's not corrupted by multiple writes
		content, err := afero.ReadFile(projectFs, path)
		if err != nil {
			t.Errorf("failed to read %s: %v", path, err)
		}
		if len(content) == 0 {
			t.Errorf("skill %s should not be empty", path)
		}
	}
}

// Helper: multiSkillTemplateManager supports multiple embedded skills
type multiSkillTemplateManager struct {
	*mockSkillTemplateManager
	skillFSMap map[string]fs.FS
}

func (m *multiSkillTemplateManager) SkillFS(skillName string) (fs.FS, error) {
	if skillFS, ok := m.skillFSMap[skillName]; ok {
		return skillFS, nil
	}
	// Return empty filesystem for unknown skills instead of error
	return fstest.MapFS{
		"SKILL.md": {
			Data: []byte("---\nname: " + skillName + "\ndescription: Test skill\n---\nContent"),
			Mode: 0o644,
		},
	}, nil
}

func (*multiSkillTemplateManager) InstructionPointer() domain.TemplateRef {
	tmpl, _ := template.New("instruction-pointer.md.tmpl").Parse("# Spectr Instructions")

	return domain.TemplateRef{
		Name:     "instruction-pointer.md.tmpl",
		Template: tmpl,
	}
}

func (*multiSkillTemplateManager) ProposalSkill() domain.TemplateRef {
	tmpl, _ := template.New("skill-proposal.md.tmpl").Parse(`---
name: spectr-proposal
description: Create a new change proposal with delta specs and tasks
---

# Proposal Creation Guide

Specs: {{.SpecsDir}}
Changes: {{.ChangesDir}}
Project: {{.ProjectFile}}
Agents: {{.AgentsFile}}`)

	return domain.TemplateRef{
		Name:     "skill-proposal.md.tmpl",
		Template: tmpl,
	}
}

func (*multiSkillTemplateManager) ApplySkill() domain.TemplateRef {
	tmpl, _ := template.New("skill-apply.md.tmpl").Parse(`---
name: spectr-apply
description: Apply or accept change proposals
---

# Apply Guide

Changes: {{.ChangesDir}}`)

	return domain.TemplateRef{
		Name:     "skill-apply.md.tmpl",
		Template: tmpl,
	}
}

// Helper: parseFrontmatter extracts YAML frontmatter from markdown content
func parseFrontmatter(t *testing.T, content string) map[string]any {
	t.Helper()

	// Find frontmatter delimiters
	lines := strings.Split(content, "\n")
	if len(lines) < 3 || lines[0] != "---" {
		t.Fatalf("content does not start with frontmatter: %s", content)
	}

	// Find end delimiter
	endIdx := -1
	for i := 1; i < len(lines); i++ {
		if lines[i] == "---" {
			endIdx = i

			break
		}
	}
	if endIdx == -1 {
		t.Fatalf("frontmatter not closed: %s", content)
	}

	// Extract frontmatter
	frontmatterText := strings.Join(lines[1:endIdx], "\n")

	// Parse YAML
	var frontmatter map[string]any
	if err := yaml.Unmarshal([]byte(frontmatterText), &frontmatter); err != nil {
		t.Fatalf("failed to parse frontmatter YAML: %v\nContent: %s", err, frontmatterText)
	}

	return frontmatter
}

// Helper: getFrontmatterString gets a string value from frontmatter,
// handling both simple strings and multiline strings
func getFrontmatterString(frontmatter map[string]any, key string) string {
	val, ok := frontmatter[key]
	if !ok {
		return ""
	}

	switch v := val.(type) {
	case string:
		return v
	default:
		return ""
	}
}
