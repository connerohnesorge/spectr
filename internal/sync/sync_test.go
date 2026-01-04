package sync

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/connerohnesorge/spectr/internal/parsers"
)

func TestBuildStatusMap(t *testing.T) {
	tests := []struct {
		name  string
		tasks []parsers.Task
		want  map[string]rune
	}{
		{
			name:  "empty tasks",
			tasks: nil,
			want:  make(map[string]rune),
		},
		{
			name: "pending maps to unchecked",
			tasks: []parsers.Task{
				{ID: "1.1", Status: parsers.TaskStatusPending},
			},
			want: map[string]rune{"1.1": ' '},
		},
		{
			name: "in_progress maps to unchecked",
			tasks: []parsers.Task{
				{ID: "1.1", Status: parsers.TaskStatusInProgress},
			},
			want: map[string]rune{"1.1": ' '},
		},
		{
			name: "completed maps to checked",
			tasks: []parsers.Task{
				{ID: "1.1", Status: parsers.TaskStatusCompleted},
			},
			want: map[string]rune{"1.1": 'x'},
		},
		{
			name: "mixed statuses",
			tasks: []parsers.Task{
				{ID: "1.1", Status: parsers.TaskStatusCompleted},
				{ID: "1.2", Status: parsers.TaskStatusInProgress},
				{ID: "2.1", Status: parsers.TaskStatusPending},
			},
			want: map[string]rune{
				"1.1": 'x',
				"1.2": ' ',
				"2.1": ' ',
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildStatusMap(tt.tasks)
			if len(got) != len(tt.want) {
				t.Errorf("buildStatusMap() returned %d items, want %d", len(got), len(tt.want))

				return
			}
			for id, wantStatus := range tt.want {
				if gotStatus, ok := got[id]; !ok {
					t.Errorf("buildStatusMap() missing key %q", id)
				} else if gotStatus != wantStatus {
					t.Errorf("buildStatusMap()[%q] = %q, want %q", id, string(gotStatus), string(wantStatus))
				}
			}
		})
	}
}

func TestUpdateTaskLine(t *testing.T) {
	tests := []struct {
		name       string
		line       string
		statusMap  map[string]rune
		wantLine   string
		wantChange bool
	}{
		{
			name:       "not a task line",
			line:       "## Section Header",
			statusMap:  map[string]rune{"1.1": 'x'},
			wantLine:   "## Section Header",
			wantChange: false,
		},
		{
			name:       "task without number",
			line:       "- [ ] Do something",
			statusMap:  map[string]rune{"1.1": 'x'},
			wantLine:   "- [ ] Do something",
			wantChange: false,
		},
		{
			name:       "task not in status map",
			line:       "- [ ] 1.1 Task one",
			statusMap:  map[string]rune{"2.1": 'x'},
			wantLine:   "- [ ] 1.1 Task one",
			wantChange: false,
		},
		{
			name:       "unchecked to checked",
			line:       "- [ ] 1.1 Task one",
			statusMap:  map[string]rune{"1.1": 'x'},
			wantLine:   "- [x] 1.1 Task one",
			wantChange: true,
		},
		{
			name:       "checked to unchecked",
			line:       "- [x] 1.1 Task one",
			statusMap:  map[string]rune{"1.1": ' '},
			wantLine:   "- [ ] 1.1 Task one",
			wantChange: true,
		},
		{
			name:       "already correct status - unchecked",
			line:       "- [ ] 1.1 Task one",
			statusMap:  map[string]rune{"1.1": ' '},
			wantLine:   "- [ ] 1.1 Task one",
			wantChange: false,
		},
		{
			name:       "already correct status - checked",
			line:       "- [x] 1.1 Task one",
			statusMap:  map[string]rune{"1.1": 'x'},
			wantLine:   "- [x] 1.1 Task one",
			wantChange: false,
		},
		{
			name:       "simple dot format",
			line:       "- [ ] 1. Simple task",
			statusMap:  map[string]rune{"1.": 'x'},
			wantLine:   "- [x] 1. Simple task",
			wantChange: true,
		},
		{
			name:       "uppercase X matches lowercase x - no change",
			line:       "- [X] 1.1 Task one",
			statusMap:  map[string]rune{"1.1": 'x'},
			wantLine:   "- [X] 1.1 Task one",
			wantChange: false, // X and x both count as checked, so no change needed
		},
		{
			name:       "preserves content with special chars",
			line:       "- [ ] 1.1 Task with `code` and **bold**",
			statusMap:  map[string]rune{"1.1": 'x'},
			wantLine:   "- [x] 1.1 Task with `code` and **bold**",
			wantChange: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLine, gotChange := updateTaskLine(tt.line, tt.statusMap)
			if gotLine != tt.wantLine {
				t.Errorf("updateTaskLine() line = %q, want %q", gotLine, tt.wantLine)
			}
			if gotChange != tt.wantChange {
				t.Errorf("updateTaskLine() changed = %v, want %v", gotChange, tt.wantChange)
			}
		})
	}
}

func TestSyncTasksToMarkdown(t *testing.T) {
	tests := []struct {
		name          string
		tasksJsonc    string
		tasksMd       string
		wantTasksMd   string
		wantSyncCount int
		skipJsonc     bool // Don't create tasks.jsonc
		skipMd        bool // Don't create tasks.md
	}{
		{
			name:          "no tasks.jsonc - skip",
			skipJsonc:     true,
			tasksMd:       "- [ ] 1.1 Task one\n",
			wantTasksMd:   "- [ ] 1.1 Task one\n",
			wantSyncCount: 0,
		},
		{
			name:          "no tasks.md - skip",
			tasksJsonc:    `{"version":1,"tasks":[{"id":"1.1","status":"completed"}]}`,
			skipMd:        true,
			wantSyncCount: 0,
		},
		{
			name:          "update single task from pending to completed",
			tasksJsonc:    `{"version":1,"tasks":[{"id":"1.1","section":"Test","description":"Task one","status":"completed"}]}`,
			tasksMd:       "- [ ] 1.1 Task one\n",
			wantTasksMd:   "- [x] 1.1 Task one\n",
			wantSyncCount: 1,
		},
		{
			name:          "update single task from completed to pending",
			tasksJsonc:    `{"version":1,"tasks":[{"id":"1.1","section":"Test","description":"Task one","status":"pending"}]}`,
			tasksMd:       "- [x] 1.1 Task one\n",
			wantTasksMd:   "- [ ] 1.1 Task one\n",
			wantSyncCount: 1,
		},
		{
			name: "preserve formatting and headers",
			tasksJsonc: `{"version":1,"tasks":[
				{"id":"1.1","section":"Setup","description":"First task","status":"completed"},
				{"id":"1.2","section":"Setup","description":"Second task","status":"pending"}
			]}`,
			tasksMd: `# Tasks

## 1. Setup

- [ ] 1.1 First task
- [x] 1.2 Second task

Some notes here.
`,
			wantTasksMd: `# Tasks

## 1. Setup

- [x] 1.1 First task
- [ ] 1.2 Second task

Some notes here.
`,
			wantSyncCount: 2,
		},
		{
			name: "no changes needed",
			tasksJsonc: `{"version":1,"tasks":[
				{"id":"1.1","section":"Test","description":"Task one","status":"completed"},
				{"id":"1.2","section":"Test","description":"Task two","status":"pending"}
			]}`,
			tasksMd: `- [x] 1.1 Task one
- [ ] 1.2 Task two
`,
			wantTasksMd: `- [x] 1.1 Task one
- [ ] 1.2 Task two
`,
			wantSyncCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tmpDir, err := os.MkdirTemp("", "sync-test-*")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer func() { _ = os.RemoveAll(tmpDir) }()

			// Write tasks.jsonc if not skipping
			if !tt.skipJsonc {
				jsonPath := filepath.Join(tmpDir, "tasks.jsonc")
				if err := os.WriteFile(jsonPath, []byte(tt.tasksJsonc), 0o644); err != nil {
					t.Fatalf("failed to write tasks.jsonc: %v", err)
				}
			}

			// Write tasks.md if not skipping
			if !tt.skipMd {
				mdPath := filepath.Join(tmpDir, "tasks.md")
				if err := os.WriteFile(mdPath, []byte(tt.tasksMd), 0o644); err != nil {
					t.Fatalf("failed to write tasks.md: %v", err)
				}
			}

			// Run sync
			syncCount, err := SyncTasksToMarkdown(tmpDir)
			if err != nil {
				t.Fatalf("SyncTasksToMarkdown() error = %v", err)
			}

			if syncCount != tt.wantSyncCount {
				t.Errorf(
					"SyncTasksToMarkdown() syncCount = %d, want %d",
					syncCount,
					tt.wantSyncCount,
				)
			}

			// Verify tasks.md content if it should exist
			if tt.skipMd {
				return
			}

			mdPath := filepath.Join(tmpDir, "tasks.md")
			gotMd, err := os.ReadFile(mdPath)
			if err != nil {
				t.Fatalf("failed to read tasks.md: %v", err)
			}
			if string(gotMd) != tt.wantTasksMd {
				t.Errorf("tasks.md content:\ngot:\n%s\nwant:\n%s", string(gotMd), tt.wantTasksMd)
			}
		})
	}
}

func TestSyncAllActiveChanges(t *testing.T) {
	// Create a temp project structure
	tmpDir, err := os.MkdirTemp("", "sync-all-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create spectr/changes structure
	changesDir := filepath.Join(tmpDir, "spectr", "changes")
	change1Dir := filepath.Join(changesDir, "change-one")
	change2Dir := filepath.Join(changesDir, "change-two")
	archiveDir := filepath.Join(changesDir, "archive", "old-change")

	for _, dir := range []string{change1Dir, change2Dir, archiveDir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("failed to create dir: %v", err)
		}
	}

	// Create proposal.md files (required for discovery)
	for _, dir := range []string{change1Dir, change2Dir, archiveDir} {
		if err := os.WriteFile(filepath.Join(dir, "proposal.md"), []byte("# Test"), 0o644); err != nil {
			t.Fatalf("failed to write proposal.md: %v", err)
		}
	}

	// Create tasks.jsonc and tasks.md for change-one
	jsonc1 := `{"version":1,"tasks":[{"id":"1.1","section":"Test","description":"Task","status":"completed"}]}`
	md1 := "- [ ] 1.1 Task\n"
	if err := os.WriteFile(filepath.Join(change1Dir, "tasks.jsonc"), []byte(jsonc1), 0o644); err != nil {
		t.Fatalf("failed to write tasks.jsonc: %v", err)
	}
	if err := os.WriteFile(filepath.Join(change1Dir, "tasks.md"), []byte(md1), 0o644); err != nil {
		t.Fatalf("failed to write tasks.md: %v", err)
	}

	// Create tasks.jsonc and tasks.md for change-two (no changes needed)
	jsonc2 := `{"version":1,"tasks":[{"id":"2.1","section":"Test","description":"Task","status":"pending"}]}`
	md2 := "- [ ] 2.1 Task\n"
	if err := os.WriteFile(filepath.Join(change2Dir, "tasks.jsonc"), []byte(jsonc2), 0o644); err != nil {
		t.Fatalf("failed to write tasks.jsonc: %v", err)
	}
	if err := os.WriteFile(filepath.Join(change2Dir, "tasks.md"), []byte(md2), 0o644); err != nil {
		t.Fatalf("failed to write tasks.md: %v", err)
	}

	// Create tasks.jsonc and tasks.md for archived change (should not be synced)
	jsonc3 := `{"version":1,"tasks":[{"id":"3.1","section":"Test","description":"Task","status":"completed"}]}`
	md3 := "- [ ] 3.1 Task\n"
	if err := os.WriteFile(filepath.Join(archiveDir, "tasks.jsonc"), []byte(jsonc3), 0o644); err != nil {
		t.Fatalf("failed to write tasks.jsonc: %v", err)
	}
	if err := os.WriteFile(filepath.Join(archiveDir, "tasks.md"), []byte(md3), 0o644); err != nil {
		t.Fatalf("failed to write tasks.md: %v", err)
	}

	// Run sync
	err = SyncAllActiveChanges(tmpDir, false)
	if err != nil {
		t.Fatalf("SyncAllActiveChanges() error = %v", err)
	}

	// Verify change-one was updated
	gotMd1, _ := os.ReadFile(filepath.Join(change1Dir, "tasks.md"))
	wantMd1 := "- [x] 1.1 Task\n"
	if string(gotMd1) != wantMd1 {
		t.Errorf("change-one/tasks.md = %q, want %q", string(gotMd1), wantMd1)
	}

	// Verify change-two was not changed (already correct)
	gotMd2, _ := os.ReadFile(filepath.Join(change2Dir, "tasks.md"))
	if string(gotMd2) != md2 {
		t.Errorf("change-two/tasks.md = %q, want %q", string(gotMd2), md2)
	}

	// Verify archived change was NOT synced (should remain unchanged)
	gotMd3, _ := os.ReadFile(filepath.Join(archiveDir, "tasks.md"))
	if string(gotMd3) != md3 {
		t.Errorf(
			"archived change was synced but shouldn't have been: got %q, want %q",
			string(gotMd3),
			md3,
		)
	}
}
