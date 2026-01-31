package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestLoadConfig_ValidConfig(t *testing.T) {
	// Create temp directory with valid config
	tmpDir := t.TempDir()
	configPath := filepath.Join(
		tmpDir,
		"spectr.yaml",
	)

	configContent := `append_tasks:
  section: "Project Workflow"
  tasks:
    - "Run linter and tests"
    - "Update changelog"
    - "Notify stakeholders"
`
	err := os.WriteFile(
		configPath,
		[]byte(configContent),
		0o644,
	)
	assert.NoError(t, err)

	cfg, err := LoadConfig(tmpDir)
	assert.NoError(t, err)
	assert.NotEqual(t, nil, cfg)
	assert.NotEqual(t, nil, cfg.AppendTasks)
	assert.Equal(
		t,
		"Project Workflow",
		cfg.AppendTasks.Section,
	)
	assert.Equal(t, 3, len(cfg.AppendTasks.Tasks))
	assert.Equal(
		t,
		"Run linter and tests",
		cfg.AppendTasks.Tasks[0],
	)
	assert.Equal(
		t,
		"Update changelog",
		cfg.AppendTasks.Tasks[1],
	)
	assert.Equal(
		t,
		"Notify stakeholders",
		cfg.AppendTasks.Tasks[2],
	)
}

func TestLoadConfig_MissingConfig(t *testing.T) {
	// Create temp directory without config
	tmpDir := t.TempDir()

	cfg, err := LoadConfig(tmpDir)
	assert.NoError(t, err)
	assert.Equal(t, (*Config)(nil), cfg)
}

func TestLoadConfig_MalformedConfig(
	t *testing.T,
) {
	// Create temp directory with invalid YAML
	tmpDir := t.TempDir()
	configPath := filepath.Join(
		tmpDir,
		"spectr.yaml",
	)

	malformedContent := `append_tasks:
  section: [invalid yaml
  tasks: {not: a list}
`
	err := os.WriteFile(
		configPath,
		[]byte(malformedContent),
		0o644,
	)
	assert.NoError(t, err)

	cfg, err := LoadConfig(tmpDir)
	assert.Error(t, err)
	assert.Equal(t, (*Config)(nil), cfg)
	assert.Contains(
		t,
		err.Error(),
		"config file is malformed",
	)
}

func TestLoadConfig_MissingSectionName(
	t *testing.T,
) {
	// Create temp directory with config missing section name
	tmpDir := t.TempDir()
	configPath := filepath.Join(
		tmpDir,
		"spectr.yaml",
	)

	configContent := `append_tasks:
  tasks:
    - "Task without section"
`
	err := os.WriteFile(
		configPath,
		[]byte(configContent),
		0o644,
	)
	assert.NoError(t, err)

	cfg, err := LoadConfig(tmpDir)
	assert.NoError(t, err)
	assert.NotEqual(t, nil, cfg)
	assert.NotEqual(t, nil, cfg.AppendTasks)
	assert.Equal(t, "", cfg.AppendTasks.Section)
	assert.Equal(
		t,
		DefaultAppendTasksSection,
		cfg.AppendTasks.GetSection(),
	)
}

func TestLoadConfig_EmptyTasksList(t *testing.T) {
	// Create temp directory with empty tasks list
	tmpDir := t.TempDir()
	configPath := filepath.Join(
		tmpDir,
		"spectr.yaml",
	)

	configContent := `append_tasks:
  section: "Empty Section"
  tasks: []
`
	err := os.WriteFile(
		configPath,
		[]byte(configContent),
		0o644,
	)
	assert.NoError(t, err)

	cfg, err := LoadConfig(tmpDir)
	assert.NoError(t, err)
	assert.NotEqual(t, nil, cfg)
	assert.NotEqual(t, nil, cfg.AppendTasks)
	assert.Equal(
		t,
		false,
		cfg.AppendTasks.HasTasks(),
	)
}

func TestLoadConfig_ConfigInParentDir(
	t *testing.T,
) {
	// Create nested directories with config in parent
	tmpDir := t.TempDir()
	subDir := filepath.Join(
		tmpDir,
		"sub",
		"nested",
	)
	err := os.MkdirAll(subDir, 0o755)
	assert.NoError(t, err)

	configPath := filepath.Join(
		tmpDir,
		"spectr.yaml",
	)
	configContent := `append_tasks:
  section: "Parent Config"
  tasks:
    - "Found in parent"
`
	err = os.WriteFile(
		configPath,
		[]byte(configContent),
		0o644,
	)
	assert.NoError(t, err)

	cfg, err := LoadConfig(subDir)
	assert.NoError(t, err)
	assert.NotEqual(t, nil, cfg)
	assert.Equal(
		t,
		"Parent Config",
		cfg.AppendTasks.Section,
	)
}

func TestAppendTasksConfig_GetSection(
	t *testing.T,
) {
	tests := []struct {
		name     string
		section  string
		expected string
	}{
		{
			name:     "custom section",
			section:  "Custom Section",
			expected: "Custom Section",
		},
		{
			name:     "empty section uses default",
			section:  "",
			expected: DefaultAppendTasksSection,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &AppendTasksConfig{
				Section: tt.section,
			}
			assert.Equal(
				t,
				tt.expected,
				cfg.GetSection(),
			)
		})
	}
}

func TestAppendTasksConfig_HasTasks(
	t *testing.T,
) {
	tests := []struct {
		name     string
		tasks    []string
		expected bool
	}{
		{
			name:     "with tasks",
			tasks:    []string{"task1", "task2"},
			expected: true,
		},
		{
			name:     "empty tasks",
			tasks:    make([]string, 0),
			expected: false,
		},
		{
			name:     "nil tasks",
			tasks:    nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &AppendTasksConfig{
				Tasks: tt.tasks,
			}
			assert.Equal(
				t,
				tt.expected,
				cfg.HasTasks(),
			)
		})
	}
}

func TestAppendTasksConfig_NilReceiver(
	t *testing.T,
) {
	var cfg *AppendTasksConfig

	t.Run(
		"GetSection on nil receiver",
		func(t *testing.T) {
			assert.Equal(
				t,
				DefaultAppendTasksSection,
				cfg.GetSection(),
			)
		},
	)

	t.Run(
		"HasTasks on nil receiver",
		func(t *testing.T) {
			assert.Equal(t, false, cfg.HasTasks())
		},
	)
}

func TestLoadConfig_RefsAlwaysPrepend(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "spectr.yaml")

	configContent := `refs_always_prepend:
  tasks:
    - "Review section requirements before starting"
    - "Verify prerequisites are met"
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	assert.NoError(t, err)

	cfg, err := LoadConfig(tmpDir)
	assert.NoError(t, err)
	assert.NotEqual(t, nil, cfg)
	assert.NotEqual(t, nil, cfg.RefsAlwaysPrepend)
	assert.True(t, cfg.RefsAlwaysPrepend.HasTasks())
	assert.Equal(t, 2, len(cfg.RefsAlwaysPrepend.Tasks))
	assert.Equal(
		t,
		"Review section requirements before starting",
		cfg.RefsAlwaysPrepend.Tasks[0],
	)
	assert.Equal(
		t,
		"Verify prerequisites are met",
		cfg.RefsAlwaysPrepend.Tasks[1],
	)
}

func TestLoadConfig_RefsAlwaysAppend(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "spectr.yaml")

	configContent := `refs_always_append:
  tasks:
    - "Verify all tasks in this section are complete"
    - "Run tests for this section"
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	assert.NoError(t, err)

	cfg, err := LoadConfig(tmpDir)
	assert.NoError(t, err)
	assert.NotEqual(t, nil, cfg)
	assert.NotEqual(t, nil, cfg.RefsAlwaysAppend)
	assert.True(t, cfg.RefsAlwaysAppend.HasTasks())
	assert.Equal(t, 2, len(cfg.RefsAlwaysAppend.Tasks))
	assert.Equal(
		t,
		"Verify all tasks in this section are complete",
		cfg.RefsAlwaysAppend.Tasks[0],
	)
	assert.Equal(
		t,
		"Run tests for this section",
		cfg.RefsAlwaysAppend.Tasks[1],
	)
}

func TestLoadConfig_BothRefsConfigs(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "spectr.yaml")

	configContent := `refs_always_prepend:
  tasks:
    - "Check prerequisites"

refs_always_append:
  tasks:
    - "Verify completion"

append_tasks:
  section: "Final Tasks"
  tasks:
    - "Output: COMPLETE"
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	assert.NoError(t, err)

	cfg, err := LoadConfig(tmpDir)
	assert.NoError(t, err)
	assert.NotEqual(t, nil, cfg)
	assert.NotEqual(t, nil, cfg.RefsAlwaysPrepend)
	assert.NotEqual(t, nil, cfg.RefsAlwaysAppend)
	assert.NotEqual(t, nil, cfg.AppendTasks)
	assert.True(t, cfg.RefsAlwaysPrepend.HasTasks())
	assert.True(t, cfg.RefsAlwaysAppend.HasTasks())
	assert.True(t, cfg.AppendTasks.HasTasks())
}

func TestRefsTasksConfig_HasTasks(t *testing.T) {
	tests := []struct {
		name     string
		tasks    []string
		expected bool
	}{
		{
			name:     "with tasks",
			tasks:    []string{"task1", "task2"},
			expected: true,
		},
		{
			name:     "empty tasks",
			tasks:    make([]string, 0),
			expected: false,
		},
		{
			name:     "nil tasks",
			tasks:    nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &RefsTasksConfig{
				Tasks: tt.tasks,
			}
			assert.Equal(t, tt.expected, cfg.HasTasks())
		})
	}
}

func TestRefsTasksConfig_NilReceiver(t *testing.T) {
	var cfg *RefsTasksConfig
	assert.False(t, cfg.HasTasks())
}
