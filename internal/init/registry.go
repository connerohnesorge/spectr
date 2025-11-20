package init

import "fmt"

// ToolRegistry manages the collection of available AI tool definitions
type ToolRegistry struct {
	tools map[ToolID]*ToolDefinition
}

// NewRegistry creates and initializes a new ToolRegistry with all
// 7 AI tool definitions (slash commands auto-installed)
func NewRegistry() *ToolRegistry {
	registry := &ToolRegistry{
		tools: make(map[ToolID]*ToolDefinition),
	}

	// Register all config-based tools from tool_definitions.go
	registry.registerTool(&ToolDefinition{
		ID:         ToolClaudeCode,
		Name:       "Claude Code",
		Type:       ToolTypeConfig,
		Priority:   1,
		Configured: false,
	})

	registry.registerTool(&ToolDefinition{
		ID:         ToolCline,
		Name:       "Cline",
		Type:       ToolTypeConfig,
		Priority:   2,
		Configured: false,
	})

	registry.registerTool(&ToolDefinition{
		ID:         ToolCostrictConfig,
		Name:       "Costrict",
		Type:       ToolTypeConfig,
		Priority:   3,
		Configured: false,
	})

	registry.registerTool(&ToolDefinition{
		ID:         ToolQoderConfig,
		Name:       "Qoder",
		Type:       ToolTypeConfig,
		Priority:   4,
		Configured: false,
	})

	registry.registerTool(&ToolDefinition{
		ID:         ToolCodeBuddy,
		Name:       "CodeBuddy",
		Type:       ToolTypeConfig,
		Priority:   5,
		Configured: false,
	})

	registry.registerTool(&ToolDefinition{
		ID:         ToolQwen,
		Name:       "Qwen",
		Type:       ToolTypeConfig,
		Priority:   6,
		Configured: false,
	})

	registry.registerTool(&ToolDefinition{
		ID:         ToolAntigravity,
		Name:       "Antigravity",
		Type:       ToolTypeConfig,
		Priority:   7,
		Configured: false,
	})

	return registry
}

// registerTool adds a tool to the registry
func (r *ToolRegistry) registerTool(tool *ToolDefinition) {
	r.tools[tool.ID] = tool
}

// GetTool retrieves a tool by its ID
// Returns an error if the tool ID is not found
func (r *ToolRegistry) GetTool(id ToolID) (*ToolDefinition, error) {
	tool, exists := r.tools[id]
	if !exists {
		return nil, fmt.Errorf("tool with ID '%s' not found", id)
	}

	return tool, nil
}

// GetToolByString retrieves a tool by its string ID
// (for backward compatibility). Returns an error if not found.
func (r *ToolRegistry) GetToolByString(id string) (*ToolDefinition, error) {
	return r.GetTool(ToolID(id))
}

// GetAllTools returns all registered tools as a slice
func (r *ToolRegistry) GetAllTools() []*ToolDefinition {
	tools := make([]*ToolDefinition, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}

	return tools
}

// GetToolsByType returns all tools of a specific type
func (r *ToolRegistry) GetToolsByType(toolType ToolType) []*ToolDefinition {
	tools := make([]*ToolDefinition, 0)
	for _, tool := range r.tools {
		if tool.Type == toolType {
			tools = append(tools, tool)
		}
	}

	return tools
}

// ListTools returns a list of all tool IDs
func (r *ToolRegistry) ListTools() []ToolID {
	ids := make([]ToolID, 0, len(r.tools))
	for id := range r.tools {
		ids = append(ids, id)
	}

	return ids
}

// GetSlashToolMapping returns the slash command tool ID for a
// config-based tool. Returns the slash tool ID and true if a mapping
// exists, zero value and false otherwise.
// This delegates to the GetSlashToolForConfig function in tool_definitions.go
func GetSlashToolMapping(configToolID ToolID) (ToolID, bool) {
	return GetSlashToolForConfig(configToolID)
}

// GetSlashToolMappingString is a backward-compatible wrapper
// that accepts strings instead of ToolID.
func GetSlashToolMappingString(configToolID string) (string, bool) {
	slashID, ok := GetSlashToolForConfig(ToolID(configToolID))

	return string(slashID), ok
}
