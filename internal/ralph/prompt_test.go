package ralph

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGeneratePrompt(t *testing.T) {
	tests := []struct {
		name        string
		task        *Task
		changeDir   string
		setupFiles  func(t *testing.T, dir string)
		wantErr     bool
		wantContain []string
	}{
		{
			name:      "nil task returns error",
			task:      nil,
			changeDir: "/some/path",
			wantErr:   true,
		},
		{
			name:      "empty changeDir returns error",
			task:      &Task{ID: "1.1", Section: "Test", Description: "Test task"},
			changeDir: "",
			wantErr:   true,
		},
		{
			name:      "missing proposal.md returns error",
			task:      &Task{ID: "1.1", Section: "Test", Description: "Test task"},
			changeDir: "/nonexistent",
			wantErr:   true,
		},
		{
			name: "basic prompt with proposal only",
			task: &Task{
				ID:          "1.1",
				Section:     "Core Infrastructure",
				Description: "Create package structure",
				Status:      "pending",
			},
			setupFiles: func(t *testing.T, dir string) {
				content := "# Change: Test Change\n\n## Why\n\nBecause reasons."
				if err := os.WriteFile(filepath.Join(dir, "proposal.md"), []byte(content), 0o644); err != nil {
					t.Fatal(err)
				}
			},
			wantErr: false,
			wantContain: []string{
				"# Task: 1.1 - Core Infrastructure",
				"## Task Description",
				"Create package structure",
				"## Change Context",
				"### Proposal",
				"# Change: Test Change",
				"## Instructions",
				"update the task status in tasks.jsonc",
			},
		},
		{
			name: "prompt with proposal and design",
			task: &Task{
				ID:          "2.1",
				Section:     "Implementation",
				Description: "Implement feature X",
				Status:      "pending",
			},
			setupFiles: func(t *testing.T, dir string) {
				proposal := "# Proposal\n\nAdd feature X"
				design := "# Design\n\n## Architecture\n\nUse pattern Y"

				if err := os.WriteFile(filepath.Join(dir, "proposal.md"), []byte(proposal), 0o644); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(dir, "design.md"), []byte(design), 0o644); err != nil {
					t.Fatal(err)
				}
			},
			wantErr: false,
			wantContain: []string{
				"# Task: 2.1 - Implementation",
				"Implement feature X",
				"### Proposal",
				"Add feature X",
				"### Design",
				"## Architecture",
				"Use pattern Y",
			},
		},
		{
			name: "prompt with delta specs",
			task: &Task{
				ID:          "3.1",
				Section:     "Specs",
				Description: "Update specs",
				Status:      "pending",
			},
			setupFiles: func(t *testing.T, dir string) {
				proposal := "# Proposal\n\nUpdate specs"

				if err := os.WriteFile(filepath.Join(dir, "proposal.md"), []byte(proposal), 0o644); err != nil {
					t.Fatal(err)
				}

				// Create specs directory structure
				specsDir := filepath.Join(dir, "specs")
				if err := os.MkdirAll(filepath.Join(specsDir, "feature-a"), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.MkdirAll(filepath.Join(specsDir, "feature-b"), 0o755); err != nil {
					t.Fatal(err)
				}

				spec1 := "## ADDED Requirements\n\n### Requirement: Feature A\n\nThe system SHALL..."
				spec2 := "## MODIFIED Requirements\n\n### Requirement: Feature B\n\nThe system SHALL..."

				if err := os.WriteFile(filepath.Join(specsDir, "feature-a", "spec.md"), []byte(spec1), 0o644); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(specsDir, "feature-b", "spec.md"), []byte(spec2), 0o644); err != nil {
					t.Fatal(err)
				}
			},
			wantErr: false,
			wantContain: []string{
				"### Relevant Specs",
				"#### feature-a",
				"## ADDED Requirements",
				"Feature A",
				"#### feature-b",
				"## MODIFIED Requirements",
				"Feature B",
			},
		},
		{
			name: "prompt with tasks.jsonc files",
			task: &Task{
				ID:          "4.1",
				Section:     "Tasks",
				Description: "Test task",
				Status:      "pending",
			},
			setupFiles: func(t *testing.T, dir string) {
				proposal := "# Proposal\n\nTest"
				if err := os.WriteFile(filepath.Join(dir, "proposal.md"), []byte(proposal), 0o644); err != nil {
					t.Fatal(err)
				}

				// Create multiple tasks.jsonc files
				tasksContent := `{"version": 1, "tasks": []}`
				if err := os.WriteFile(filepath.Join(dir, "tasks.jsonc"), []byte(tasksContent), 0o644); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(dir, "tasks-2.jsonc"), []byte(tasksContent), 0o644); err != nil {
					t.Fatal(err)
				}
			},
			wantErr: false,
			wantContain: []string{
				"Task status files:",
				"tasks.jsonc",
				"tasks-2.jsonc",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dir string
			if tt.setupFiles != nil {
				// Create temporary directory for test
				dir = t.TempDir()
				tt.setupFiles(t, dir)
			} else if tt.changeDir != "" {
				dir = tt.changeDir
			}

			got, err := GeneratePrompt(tt.task, dir)

			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratePrompt() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if tt.wantErr {
				return
			}

			// Check that all expected strings are present in the output
			for _, want := range tt.wantContain {
				if !strings.Contains(got, want) {
					t.Errorf(
						"GeneratePrompt() output missing expected content:\nwant substring: %q\ngot:\n%s",
						want,
						got,
					)
				}
			}
		})
	}
}

func TestLoadDeltaSpecs(t *testing.T) {
	tests := []struct {
		name       string
		setupFiles func(t *testing.T, dir string)
		wantSpecs  int
		wantErr    bool
	}{
		{
			name:       "missing specs directory returns no error",
			setupFiles: func(_ *testing.T, _ string) {},
			wantSpecs:  0,
			wantErr:    false,
		},
		{
			name: "empty specs directory",
			setupFiles: func(t *testing.T, dir string) {
				if err := os.MkdirAll(filepath.Join(dir, "specs"), 0o755); err != nil {
					t.Fatal(err)
				}
			},
			wantSpecs: 0,
			wantErr:   false,
		},
		{
			name: "single spec file",
			setupFiles: func(t *testing.T, dir string) {
				specsDir := filepath.Join(dir, "specs", "feature-x")
				if err := os.MkdirAll(specsDir, 0o755); err != nil {
					t.Fatal(err)
				}
				content := "# Spec content"
				if err := os.WriteFile(filepath.Join(specsDir, "spec.md"), []byte(content), 0o644); err != nil {
					t.Fatal(err)
				}
			},
			wantSpecs: 1,
			wantErr:   false,
		},
		{
			name: "multiple spec files",
			setupFiles: func(t *testing.T, dir string) {
				specsDir := filepath.Join(dir, "specs")
				features := []string{"feature-a", "feature-b", "feature-c"}
				for _, feature := range features {
					featureDir := filepath.Join(specsDir, feature)
					if err := os.MkdirAll(featureDir, 0o755); err != nil {
						t.Fatal(err)
					}
					content := "# " + feature
					if err := os.WriteFile(filepath.Join(featureDir, "spec.md"), []byte(content), 0o644); err != nil {
						t.Fatal(err)
					}
				}
			},
			wantSpecs: 3,
			wantErr:   false,
		},
		{
			name: "ignores non-spec.md files",
			setupFiles: func(t *testing.T, dir string) {
				specsDir := filepath.Join(dir, "specs", "feature-x")
				if err := os.MkdirAll(specsDir, 0o755); err != nil {
					t.Fatal(err)
				}
				// Create spec.md (should be loaded)
				if err := os.WriteFile(filepath.Join(specsDir, "spec.md"), []byte("valid"), 0o644); err != nil {
					t.Fatal(err)
				}
				// Create other files (should be ignored)
				if err := os.WriteFile(filepath.Join(specsDir, "README.md"), []byte("ignore"), 0o644); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(specsDir, "notes.txt"), []byte("ignore"), 0o644); err != nil {
					t.Fatal(err)
				}
			},
			wantSpecs: 1,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			tt.setupFiles(t, dir)

			template := &PromptTemplate{
				DeltaSpecs: make(map[string]string),
			}

			specsDir := filepath.Join(dir, "specs")
			err := loadDeltaSpecs(specsDir, template)

			if (err != nil) != tt.wantErr {
				t.Errorf("loadDeltaSpecs() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if len(template.DeltaSpecs) != tt.wantSpecs {
				t.Errorf(
					"loadDeltaSpecs() loaded %d specs, want %d",
					len(template.DeltaSpecs),
					tt.wantSpecs,
				)
			}
		})
	}
}

func TestAssemblePrompt(t *testing.T) {
	tests := []struct {
		name           string
		template       *PromptTemplate
		wantContain    []string
		wantNotContain []string
	}{
		{
			name: "minimal template",
			template: &PromptTemplate{
				Task: &Task{
					ID:          "1.1",
					Section:     "Test",
					Description: "Test task",
				},
				Proposal: "# Proposal\n\nTest proposal",
			},
			wantContain: []string{
				"# Task: 1.1 - Test",
				"## Task Description",
				"Test task",
				"### Proposal",
				"Test proposal",
				"## Instructions",
			},
			wantNotContain: []string{
				"### Design",
				"### Relevant Specs",
				"Task status files:",
			},
		},
		{
			name: "full template with all fields",
			template: &PromptTemplate{
				Task: &Task{
					ID:          "2.1",
					Section:     "Full Test",
					Description: "Complete test",
				},
				Proposal: "# Proposal content",
				Design:   "# Design content",
				DeltaSpecs: map[string]string{
					"feature-a": "Spec A content",
					"feature-b": "Spec B content",
				},
				TasksJSONCs: []string{
					"/path/to/tasks.jsonc",
					"/path/to/tasks-2.jsonc",
				},
			},
			wantContain: []string{
				"# Task: 2.1 - Full Test",
				"Complete test",
				"### Proposal",
				"Proposal content",
				"### Design",
				"Design content",
				"### Relevant Specs",
				"#### feature-a",
				"Spec A content",
				"#### feature-b",
				"Spec B content",
				"Task status files:",
				"/path/to/tasks.jsonc",
				"/path/to/tasks-2.jsonc",
			},
		},
		{
			name: "template with design but no specs",
			template: &PromptTemplate{
				Task: &Task{
					ID:          "3.1",
					Section:     "Partial",
					Description: "Partial test",
				},
				Proposal: "Proposal",
				Design:   "Design",
			},
			wantContain: []string{
				"### Design",
			},
			wantNotContain: []string{
				"### Relevant Specs",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := assemblePrompt(tt.template)

			for _, want := range tt.wantContain {
				if !strings.Contains(got, want) {
					t.Errorf("assemblePrompt() missing expected content: %q", want)
				}
			}

			for _, notWant := range tt.wantNotContain {
				if strings.Contains(got, notWant) {
					t.Errorf("assemblePrompt() contains unexpected content: %q", notWant)
				}
			}
		})
	}
}
