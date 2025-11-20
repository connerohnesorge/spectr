package init

import (
	"testing"
)

func TestNewRegistryFromProviders(t *testing.T) {
	registry := NewRegistryFromProviders()

	// Test that the registry is not nil
	if registry == nil {
		t.Fatal("NewRegistryFromProviders() returned nil")
	}

	// Test that all 7 tools are registered (slash commands auto-installed)
	allTools := registry.GetAllTools()
	if len(allTools) != 7 {
		t.Errorf("Expected 7 tools, got %d", len(allTools))
	}

	// Test that the tools map is not nil
	if registry.tools == nil {
		t.Error("registry.tools map is nil")
	}
}

func TestGetTool(t *testing.T) {
	registry := NewRegistryFromProviders()

	tests := []struct {
		name    string
		toolID  string
		wantErr bool
	}{
		{"Get Claude Code", "claude-code", false},
		{"Get Cline", "cline", false},
		{"Get Costrict Config", "costrict-config", false},
		{"Get Qoder Config", "qoder-config", false},
		{"Get CodeBuddy", "codebuddy", false},
		{"Get Qwen", "qwen", false},
		{"Get Antigravity", "antigravity", false},
		{"Get Invalid Tool", "nonexistent", true},
		{"Get Slash Tool (removed)", "claude", true},
		{"Get Slash Tool (removed)", "cursor", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tool, err := registry.GetTool(tt.toolID)
			if tt.wantErr {
				verifyErrorCase(t, tt.toolID, tool, err)
			} else {
				verifySuccessCase(t, tt.toolID, tool, err)
			}
		})
	}
}

// verifyErrorCase checks that GetTool properly returns an error
func verifyErrorCase(
	t *testing.T,
	toolID string,
	tool *ToolDefinition,
	err error,
) {
	t.Helper()
	if err == nil {
		t.Errorf("GetTool(%s) expected error, got nil", toolID)
	}
	if tool != nil {
		t.Errorf(
			"GetTool(%s) expected nil tool, got %v",
			toolID,
			tool,
		)
	}
}

// verifySuccessCase checks that GetTool returns a valid tool
func verifySuccessCase(
	t *testing.T,
	toolID string,
	tool *ToolDefinition,
	err error,
) {
	t.Helper()
	if err != nil {
		t.Errorf("GetTool(%s) unexpected error: %v", toolID, err)
	}
	if tool == nil {
		t.Errorf("GetTool(%s) returned nil tool", toolID)

		return
	}
	if tool.ID != toolID {
		t.Errorf(
			"GetTool(%s) returned tool with ID %s",
			toolID,
			tool.ID,
		)
	}
}

func TestGetToolsByType(t *testing.T) {
	registry := NewRegistryFromProviders()

	// Test config-based tools
	configTools := registry.GetToolsByType(ToolTypeConfig)
	if len(configTools) != 7 {
		t.Errorf("Expected 7 config tools, got %d", len(configTools))
	}

	// Verify all config tools have correct type
	for _, tool := range configTools {
		if tool.Type != ToolTypeConfig {
			t.Errorf("Config tool %s has wrong type: %s", tool.ID, tool.Type)
		}
	}

	// Test slash command tools (should be 0 - auto-installed now)
	slashTools := registry.GetToolsByType(ToolTypeSlash)
	if len(slashTools) != 0 {
		t.Errorf("Expected 0 slash tools (auto-installed), got %d", len(slashTools))
	}
}

func TestSlashToolsNotInRegistry(t *testing.T) {
	registry := NewRegistryFromProviders()

	// Slash-only tools should no longer be in registry (auto-installed)
	removedSlashTools := []string{
		"claude",
		"kilocode",
		"qoder-slash",
		"cursor",
		"aider",
		"continue",
		"copilot",
		"mentat",
		"tabnine",
		"smol",
		"costrict-slash",
	}

	for _, id := range removedSlashTools {
		tool, err := registry.GetTool(id)
		if err == nil {
			t.Errorf("Slash-only tool %s should not be in registry (auto-installed now)", id)
		}
		if tool != nil {
			t.Errorf("GetTool(%s) should return nil, got %v", id, tool)
		}
	}
}

func TestListTools(t *testing.T) {
	registry := NewRegistryFromProviders()

	toolIDs := registry.ListTools()

	// Test that we get 7 tool IDs (slash commands auto-installed)
	if len(toolIDs) != 7 {
		t.Errorf("Expected 7 tool IDs, got %d", len(toolIDs))
	}

	// Test that all tool IDs are unique
	seen := make(map[string]bool)
	for _, id := range toolIDs {
		if seen[id] {
			t.Errorf("Duplicate tool ID found: %s", id)
		}
		seen[id] = true
	}
}

func TestToolIDsAreKebabCase(t *testing.T) {
	registry := NewRegistryFromProviders()

	allTools := registry.GetAllTools()

	for _, tool := range allTools {
		// Check that ID contains only lowercase letters, numbers, and hyphens
		for _, char := range tool.ID {
			if (char < 'a' || char > 'z') && (char < '0' || char > '9') && char != '-' {
				t.Errorf(
					"Tool ID %s is not in kebab-case (contains invalid character: %c)",
					tool.ID,
					char,
				)
			}
		}
	}
}

func TestAllToolsHaveRequiredFields(t *testing.T) {
	registry := NewRegistryFromProviders()

	allTools := registry.GetAllTools()

	for _, tool := range allTools {
		if tool.ID == "" {
			t.Error("Found tool with empty ID")
		}
		if tool.Name == "" {
			t.Errorf("Tool %s has empty Name", tool.ID)
		}
		if tool.Type != ToolTypeConfig && tool.Type != ToolTypeSlash {
			t.Errorf("Tool %s has invalid Type: %s", tool.ID, tool.Type)
		}
		if tool.Priority < 1 || tool.Priority > 7 {
			t.Errorf("Tool %s has invalid Priority: %d (should be 1-7)", tool.ID, tool.Priority)
		}
		if tool.Configured {
			t.Errorf("Tool %s should start with Configured=false", tool.ID)
		}
	}
}

func TestPrioritiesAreUnique(t *testing.T) {
	registry := NewRegistryFromProviders()

	allTools := registry.GetAllTools()
	priorities := make(map[int]string)

	for _, tool := range allTools {
		if existingTool, exists := priorities[tool.Priority]; exists {
			t.Errorf(
				"Duplicate priority %d found for tools %s and %s",
				tool.Priority,
				existingTool,
				tool.ID,
			)
		}
		priorities[tool.Priority] = tool.ID
	}
}
