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
	configPath := filepath.Join(tmpDir, "spectr.yaml")

	configContent := `append_tasks:
  section: "Project Workflow"
  tasks:
    - "Run linter and tests"
    - "Update changelog"
    - "Notify stakeholders"
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	assert.NoError(t, err)

	cfg, err := LoadConfig(tmpDir)
	assert.NoError(t, err)
	assert.NotEqual(t, nil, cfg)
	assert.NotEqual(t, nil, cfg.AppendTasks)
	assert.Equal(t, "Project Workflow", cfg.AppendTasks.Section)
	assert.Equal(t, 3, len(cfg.AppendTasks.Tasks))
	assert.Equal(t, "Run linter and tests", cfg.AppendTasks.Tasks[0])
	assert.Equal(t, "Update changelog", cfg.AppendTasks.Tasks[1])
	assert.Equal(t, "Notify stakeholders", cfg.AppendTasks.Tasks[2])
}

func TestLoadConfig_MissingConfig(t *testing.T) {
	// Create temp directory without config
	tmpDir := t.TempDir()

	cfg, err := LoadConfig(tmpDir)
	assert.NoError(t, err)
	assert.Equal(t, (*Config)(nil), cfg)
}

func TestLoadConfig_MalformedConfig(t *testing.T) {
	// Create temp directory with invalid YAML
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "spectr.yaml")

	malformedContent := `append_tasks:
  section: [invalid yaml
  tasks: {not: a list}
`
	err := os.WriteFile(configPath, []byte(malformedContent), 0o644)
	assert.NoError(t, err)

	cfg, err := LoadConfig(tmpDir)
	assert.Error(t, err)
	assert.Equal(t, (*Config)(nil), cfg)
	assert.Contains(t, err.Error(), "failed to parse config file")
}

func TestLoadConfig_MissingSectionName(t *testing.T) {
	// Create temp directory with config missing section name
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "spectr.yaml")

	configContent := `append_tasks:
  tasks:
    - "Task without section"
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	assert.NoError(t, err)

	cfg, err := LoadConfig(tmpDir)
	assert.NoError(t, err)
	assert.NotEqual(t, nil, cfg)
	assert.NotEqual(t, nil, cfg.AppendTasks)
	assert.Equal(t, "", cfg.AppendTasks.Section)
	assert.Equal(t, DefaultAppendTasksSection, cfg.AppendTasks.GetSection())
}

func TestLoadConfig_EmptyTasksList(t *testing.T) {
	// Create temp directory with empty tasks list
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "spectr.yaml")

	configContent := `append_tasks:
  section: "Empty Section"
  tasks: []
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	assert.NoError(t, err)

	cfg, err := LoadConfig(tmpDir)
	assert.NoError(t, err)
	assert.NotEqual(t, nil, cfg)
	assert.NotEqual(t, nil, cfg.AppendTasks)
	assert.Equal(t, false, cfg.AppendTasks.HasTasks())
}

func TestLoadConfig_ConfigInParentDir(t *testing.T) {
	// Create nested directories with config in parent
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "sub", "nested")
	err := os.MkdirAll(subDir, 0o755)
	assert.NoError(t, err)

	configPath := filepath.Join(tmpDir, "spectr.yaml")
	configContent := `append_tasks:
  section: "Parent Config"
  tasks:
    - "Found in parent"
`
	err = os.WriteFile(configPath, []byte(configContent), 0o644)
	assert.NoError(t, err)

	cfg, err := LoadConfig(subDir)
	assert.NoError(t, err)
	assert.NotEqual(t, nil, cfg)
	assert.Equal(t, "Parent Config", cfg.AppendTasks.Section)
}

func TestAppendTasksConfig_GetSection(t *testing.T) {
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
			cfg := &AppendTasksConfig{Section: tt.section}
			assert.Equal(t, tt.expected, cfg.GetSection())
		})
	}
}

func TestAppendTasksConfig_HasTasks(t *testing.T) {
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
			cfg := &AppendTasksConfig{Tasks: tt.tasks}
			assert.Equal(t, tt.expected, cfg.HasTasks())
		})
	}
}
