package init

import "fmt"

// ToolRegistry manages the collection of available AI tool definitions
type ToolRegistry struct {
	tools map[string]*ToolDefinition
}

// NewRegistry creates and initializes a new ToolRegistry with all
// 6 AI tool definitions (slash commands auto-installed)
func NewRegistry() *ToolRegistry {
	registry := &ToolRegistry{
		tools: make(map[string]*ToolDefinition),
	}

	// Config-based tools (6 tools)
	// Each tool auto-installs its corresponding slash commands
	registry.registerTool(&ToolDefinition{
		ID:         "claude-code",
		Name:       "Claude Code",
		Type:       ToolTypeConfig,
		ConfigPath: ".claude/claude.json",
		Priority:   1,
		Configured: false,
	})

	registry.registerTool(&ToolDefinition{
		ID:         "cline",
		Name:       "Cline",
		Type:       ToolTypeConfig,
		ConfigPath: ".cline/cline_mcp_settings.json",
		Priority:   2,
		Configured: false,
	})

	registry.registerTool(&ToolDefinition{
		ID:         "costrict-config",
		Name:       "Costrict",
		Type:       ToolTypeConfig,
		ConfigPath: ".costrict/config.json",
		Priority:   3,
		Configured: false,
	})

	registry.registerTool(&ToolDefinition{
		ID:         "qoder-config",
		Name:       "Qoder",
		Type:       ToolTypeConfig,
		ConfigPath: ".qoder/config.json",
		Priority:   4,
		Configured: false,
	})

	registry.registerTool(&ToolDefinition{
		ID:         "codebuddy",
		Name:       "CodeBuddy",
		Type:       ToolTypeConfig,
		ConfigPath: ".codebuddy/config.json",
		Priority:   5,
		Configured: false,
	})

	registry.registerTool(&ToolDefinition{
		ID:         "qwen",
		Name:       "Qwen",
		Type:       ToolTypeConfig,
		ConfigPath: ".qwen/config.json",
		Priority:   6,
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
func (r *ToolRegistry) GetTool(id string) (*ToolDefinition, error) {
	tool, exists := r.tools[id]
	if !exists {
		return nil, fmt.Errorf("tool with ID '%s' not found", id)
	}

	return tool, nil
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
func (r *ToolRegistry) ListTools() []string {
	ids := make([]string, 0, len(r.tools))
	for id := range r.tools {
		ids = append(ids, id)
	}

	return ids
}

// configToSlashMapping maps config-based tool IDs to their slash
// command equivalents
var configToSlashMapping = map[string]string{
	"claude-code":     "claude",
	"cline":           "cline-slash",
	"costrict-config": "costrict-slash",
	"qoder-config":    "qoder-slash",
	"codebuddy":       "codebuddy-slash",
	"qwen":            "qwen-slash",
}

// GetSlashToolMapping returns the slash command tool ID for a
// config-based tool. Returns the slash tool ID and true if a mapping
// exists, empty string and false otherwise
func GetSlashToolMapping(configToolID string) (string, bool) {
	slashID, exists := configToSlashMapping[configToolID]

	return slashID, exists
}
