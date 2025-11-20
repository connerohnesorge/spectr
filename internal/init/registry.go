package init

import (
	"fmt"

	"github.com/conneroisu/spectr/internal/providers"
)

// ToolRegistry manages the collection of available AI tool definitions
type ToolRegistry struct {
	tools map[string]*ToolDefinition
}

// NewRegistryFromProviders creates a ToolRegistry populated from the
// provider registry. This replaces the old hardcoded NewRegistry() function.
func NewRegistryFromProviders() *ToolRegistry {
	configProviders := providers.ListProvidersByType(providers.TypeConfig)
	registry := &ToolRegistry{tools: make(map[string]*ToolDefinition)}

	for _, p := range configProviders {
		registry.registerTool(&ToolDefinition{
			ID:         p.ID,
			Name:       p.Name,
			Type:       ToolTypeConfig,
			Priority:   p.Priority,
			Configured: false,
		})
	}

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
