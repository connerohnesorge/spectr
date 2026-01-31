package cmd

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/connerohnesorge/spectr/internal/config"
	"github.com/connerohnesorge/spectr/internal/parsers"
)

func TestCreateChildPrependedTasks(t *testing.T) {
	tests := []struct {
		name       string
		sectionNum string
		cfg        *config.RefsTasksConfig
		wantIDs    []string
		wantLen    int
	}{
		{
			name:       "single task",
			sectionNum: "1",
			cfg: &config.RefsTasksConfig{
				Tasks: []string{"Check prerequisites"},
			},
			wantIDs: []string{"1.0.1"},
			wantLen: 1,
		},
		{
			name:       "multiple tasks",
			sectionNum: "2",
			cfg: &config.RefsTasksConfig{
				Tasks: []string{
					"Review requirements",
					"Verify prerequisites",
					"Check dependencies",
				},
			},
			wantIDs: []string{"2.0.1", "2.0.2", "2.0.3"},
			wantLen: 3,
		},
		{
			name:       "nil config",
			sectionNum: "1",
			cfg:        nil,
			wantIDs:    nil,
			wantLen:    0,
		},
		{
			name:       "empty tasks",
			sectionNum: "1",
			cfg: &config.RefsTasksConfig{
				Tasks: make([]string, 0),
			},
			wantIDs: nil,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks := createChildPrependedTasks(tt.sectionNum, tt.cfg)
			assert.Equal(t, tt.wantLen, len(tasks))

			for i, task := range tasks {
				assert.Equal(t, tt.wantIDs[i], task.ID)
				assert.Equal(t, "Section Prerequisites", task.Section)
				assert.Equal(t, parsers.TaskStatusPending, task.Status)
			}
		})
	}
}

func TestCreateChildAppendedTasks(t *testing.T) {
	tests := []struct {
		name       string
		sectionNum string
		cfg        *config.RefsTasksConfig
		wantIDs    []string
		wantLen    int
	}{
		{
			name:       "single task",
			sectionNum: "1",
			cfg: &config.RefsTasksConfig{
				Tasks: []string{"Verify completion"},
			},
			wantIDs: []string{"1.99.1"},
			wantLen: 1,
		},
		{
			name:       "multiple tasks",
			sectionNum: "3",
			cfg: &config.RefsTasksConfig{
				Tasks: []string{
					"Verify all tasks complete",
					"Run tests",
					"Update documentation",
				},
			},
			wantIDs: []string{"3.99.1", "3.99.2", "3.99.3"},
			wantLen: 3,
		},
		{
			name:       "nil config",
			sectionNum: "1",
			cfg:        nil,
			wantIDs:    nil,
			wantLen:    0,
		},
		{
			name:       "empty tasks",
			sectionNum: "1",
			cfg: &config.RefsTasksConfig{
				Tasks: make([]string, 0),
			},
			wantIDs: nil,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks := createChildAppendedTasks(tt.sectionNum, tt.cfg)
			assert.Equal(t, tt.wantLen, len(tasks))

			for i, task := range tasks {
				assert.Equal(t, tt.wantIDs[i], task.ID)
				assert.Equal(t, "Section Verification", task.Section)
				assert.Equal(t, parsers.TaskStatusPending, task.Status)
			}
		})
	}
}

func TestCreateChildPrependedTasks_TaskDescriptions(t *testing.T) {
	cfg := &config.RefsTasksConfig{
		Tasks: []string{
			"First task description",
			"Second task description",
		},
	}

	tasks := createChildPrependedTasks("1", cfg)
	assert.Equal(t, 2, len(tasks))
	assert.Equal(t, "First task description", tasks[0].Description)
	assert.Equal(t, "Second task description", tasks[1].Description)
}

func TestCreateChildAppendedTasks_TaskDescriptions(t *testing.T) {
	cfg := &config.RefsTasksConfig{
		Tasks: []string{
			"Verify completion",
			"Run final tests",
		},
	}

	tasks := createChildAppendedTasks("2", cfg)
	assert.Equal(t, 2, len(tasks))
	assert.Equal(t, "Verify completion", tasks[0].Description)
	assert.Equal(t, "Run final tests", tasks[1].Description)
}

func TestStatusPreservationForInjectedTasks(t *testing.T) {
	// Create status map with some injected task IDs
	statusMap := map[string]parsers.TaskStatusValue{
		"1.0.1":  parsers.TaskStatusCompleted,
		"1.0.2":  parsers.TaskStatusInProgress,
		"1.1":    parsers.TaskStatusCompleted,
		"1.99.1": parsers.TaskStatusCompleted,
	}

	// Create prepended tasks
	prependCfg := &config.RefsTasksConfig{
		Tasks: []string{"Task 1", "Task 2"},
	}
	prependedTasks := createChildPrependedTasks("1", prependCfg)

	// Create appended tasks
	appendCfg := &config.RefsTasksConfig{
		Tasks: []string{"Verify"},
	}
	appendedTasks := createChildAppendedTasks("1", appendCfg)

	// Apply status preservation
	applyStatusPreservation(prependedTasks, statusMap)
	applyStatusPreservation(appendedTasks, statusMap)

	// Verify statuses were preserved
	assert.Equal(t, parsers.TaskStatusCompleted, prependedTasks[0].Status)
	assert.Equal(t, parsers.TaskStatusInProgress, prependedTasks[1].Status)
	assert.Equal(t, parsers.TaskStatusCompleted, appendedTasks[0].Status)
}

func TestAggregateSectionStatusWithInjectedTasks(t *testing.T) {
	statusMap := map[string]parsers.TaskStatusValue{
		"1.0.1":  parsers.TaskStatusCompleted,
		"1.1":    parsers.TaskStatusCompleted,
		"1.2":    parsers.TaskStatusCompleted,
		"1.99.1": parsers.TaskStatusCompleted,
	}

	// Build combined task list like what would be in a child file
	tasks := []parsers.Task{
		{ID: "1.0.1", Status: parsers.TaskStatusPending},
		{ID: "1.1", Status: parsers.TaskStatusPending},
		{ID: "1.2", Status: parsers.TaskStatusPending},
		{ID: "1.99.1", Status: parsers.TaskStatusPending},
	}

	// All completed -> should return completed
	status := aggregateSectionStatus(tasks, statusMap)
	assert.Equal(t, parsers.TaskStatusCompleted, status)
}

func TestAggregateSectionStatusWithMixedInjectedTasks(t *testing.T) {
	statusMap := map[string]parsers.TaskStatusValue{
		"1.0.1":  parsers.TaskStatusCompleted,
		"1.1":    parsers.TaskStatusCompleted,
		"1.99.1": parsers.TaskStatusInProgress, // One injected task in progress
	}

	tasks := []parsers.Task{
		{ID: "1.0.1", Status: parsers.TaskStatusPending},
		{ID: "1.1", Status: parsers.TaskStatusPending},
		{ID: "1.99.1", Status: parsers.TaskStatusPending},
	}

	// One in progress -> should return in_progress
	status := aggregateSectionStatus(tasks, statusMap)
	assert.Equal(t, parsers.TaskStatusInProgress, status)
}
