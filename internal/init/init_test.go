package init

import (
	"testing"
)

// TestToolTypeConstants verifies that the ToolType constants are defined correctly
func TestToolTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		toolType ToolType
		expected string
	}{
		{
			name:     "Config tool type",
			toolType: ToolTypeConfig,
			expected: "config",
		},
		{
			name:     "Slash tool type",
			toolType: ToolTypeSlash,
			expected: "slash",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.toolType) != tt.expected {
				t.Errorf("ToolType = %v, want %v", tt.toolType, tt.expected)
			}
		})
	}
}

// TestToolDefinitionCreation verifies that ToolDefinition can be created properly
func TestToolDefinitionCreation(t *testing.T) {
	tool := ToolDefinition{
		ID:           "test-tool",
		Name:         "Test Tool",
		Type:         ToolTypeConfig,
		SlashCommand: "",
		Priority:     1,
		Configured:   false,
	}

	if tool.ID != "test-tool" {
		t.Errorf("ToolDefinition.ID = %v, want %v", tool.ID, "test-tool")
	}
	if tool.Name != "Test Tool" {
		t.Errorf("ToolDefinition.Name = %v, want %v", tool.Name, "Test Tool")
	}
	if tool.Type != ToolTypeConfig {
		t.Errorf("ToolDefinition.Type = %v, want %v", tool.Type, ToolTypeConfig)
	}
	if tool.Configured {
		t.Errorf("ToolDefinition.Configured = %v, want %v", tool.Configured, false)
	}
}

// TestProjectConfigCreation verifies that ProjectConfig can be created properly
func TestProjectConfigCreation(t *testing.T) {
	config := ProjectConfig{
		ProjectPath:   "/test/path",
		SelectedTools: []string{"tool1", "tool2"},
		SpectrEnabled: true,
	}

	if config.ProjectPath != "/test/path" {
		t.Errorf("ProjectConfig.ProjectPath = %v, want %v", config.ProjectPath, "/test/path")
	}
	if len(config.SelectedTools) != 2 {
		t.Errorf("ProjectConfig.SelectedTools length = %v, want %v", len(config.SelectedTools), 2)
	}
	if !config.SpectrEnabled {
		t.Errorf("ProjectConfig.SpectrEnabled = %v, want %v", config.SpectrEnabled, true)
	}
}

// TestInitStateConstants verifies that the InitState constants are defined correctly
func TestInitStateConstants(t *testing.T) {
	tests := []struct {
		name     string
		state    InitState
		expected int
	}{
		{
			name:     "Select tools state",
			state:    StateSelectTools,
			expected: 0,
		},
		{
			name:     "Configure tools state",
			state:    StateConfigureTools,
			expected: 1,
		},
		{
			name:     "Confirmation state",
			state:    StateConfirmation,
			expected: 2,
		},
		{
			name:     "Complete state",
			state:    StateComplete,
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.state) != tt.expected {
				t.Errorf("InitState = %v, want %v", tt.state, tt.expected)
			}
		})
	}
}
