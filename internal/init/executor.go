// Package init provides utilities for initializing Spectr
// in a project directory.
//
//nolint:revive // file-length-limit - logically cohesive, no benefit to split
package init

import (
	"fmt"
	"os"
	"path/filepath"
)

// InitExecutor handles the actual initialization process
type InitExecutor struct {
	projectPath string
	registry    *ToolRegistry
	tm          *TemplateManager
}

// NewInitExecutor creates a new initialization executor
func NewInitExecutor(projectPath string) (*InitExecutor, error) {
	// Expand and validate path
	expandedPath, err := ExpandPath(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to expand project path: %w", err)
	}

	// Check if path exists
	if !FileExists(expandedPath) {
		return nil, fmt.Errorf("project path does not exist: %s", expandedPath)
	}

	// Initialize template manager
	tm, err := NewTemplateManager()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize template manager: %w", err)
	}

	return &InitExecutor{
		projectPath: expandedPath,
		registry:    NewRegistry(),
		tm:          tm,
	}, nil
}

// Execute runs the initialization process
func (e *InitExecutor) Execute(
	selectedToolIDs []string,
) (*ExecutionResult, error) {
	result := &ExecutionResult{
		CreatedFiles: make([]string, 0),
		UpdatedFiles: make([]string, 0),
		Errors:       make([]string, 0),
	}

	// 1. Check if Spectr is already initialized
	if IsSpectrInitialized(e.projectPath) {
		result.Errors = append(
			result.Errors,
			"Spectr already initialized in this project",
		)
		// Don't return error - allow updating tool configurations
	}

	// 2. Create spectr/ directory structure
	spectrDir := filepath.Join(e.projectPath, "spectr")
	if err := e.createDirectoryStructure(spectrDir, result); err != nil {
		return result, fmt.Errorf(
			"failed to create directory structure: %w",
			err,
		)
	}

	// 3. Create project.md
	if err := e.createProjectMd(spectrDir, result); err != nil {
		result.Errors = append(
			result.Errors,
			fmt.Sprintf("failed to create project.md: %v", err),
		)
	}

	// 4. Create AGENTS.md
	if err := e.createAgentsMd(spectrDir, result); err != nil {
		result.Errors = append(
			result.Errors,
			fmt.Sprintf("failed to create AGENTS.md: %v", err),
		)
	}

	// 5. Configure selected tools
	if err := e.configureTools(selectedToolIDs, spectrDir, result); err != nil {
		result.Errors = append(
			result.Errors,
			fmt.Sprintf("failed to configure tools: %v", err),
		)
	}

	// 6. Create README if it doesn't exist
	if err := e.createReadmeIfMissing(result); err != nil {
		result.Errors = append(
			result.Errors,
			fmt.Sprintf("failed to create README: %v", err),
		)
	}

	return result, nil
}

// createDirectoryStructure creates the spectr/ directory
// and subdirectories
func (_e *InitExecutor) createDirectoryStructure(
	spectrDir string,
	result *ExecutionResult,
) error {
	dirs := []string{
		spectrDir,
		filepath.Join(spectrDir, "specs"),
		filepath.Join(spectrDir, "changes"),
	}

	for _, dir := range dirs {
		if !FileExists(dir) {
			if err := EnsureDir(dir); err != nil {
				return fmt.Errorf(
					"failed to create directory %s: %w",
					dir,
					err,
				)
			}
			result.CreatedFiles = append(result.CreatedFiles, dir+"/")
		}
	}

	return nil
}

// createProjectMd creates the project.md file
func (e *InitExecutor) createProjectMd(
	spectrDir string,
	result *ExecutionResult,
) error {
	projectFile := filepath.Join(spectrDir, "project.md")

	// Check if it already exists
	if FileExists(projectFile) {
		result.Errors = append(
			result.Errors,
			"project.md already exists, skipping",
		)

		return nil
	}

	// Get project name from directory
	projectName := filepath.Base(e.projectPath)

	// Render template
	content, err := e.tm.RenderProject(ProjectContext{
		ProjectName: projectName,
		Description: "Add your project description here",
		TechStack:   []string{"Add", "Your", "Technologies", "Here"},
		Conventions: "",
	})
	if err != nil {
		return fmt.Errorf("failed to render project template: %w", err)
	}

	// Write file
	if err := os.WriteFile(projectFile, []byte(content), filePerm); err != nil {
		return fmt.Errorf("failed to write project.md: %w", err)
	}

	result.CreatedFiles = append(result.CreatedFiles, "spectr/project.md")

	return nil
}

// createAgentsMd creates the AGENTS.md file
func (e *InitExecutor) createAgentsMd(
	spectrDir string,
	result *ExecutionResult,
) error {
	agentsFile := filepath.Join(spectrDir, "AGENTS.md")

	// Check if it already exists
	if FileExists(agentsFile) {
		result.Errors = append(
			result.Errors,
			"AGENTS.md already exists, skipping",
		)

		return nil
	}

	// Render template
	content, err := e.tm.RenderAgents()
	if err != nil {
		return fmt.Errorf("failed to render agents template: %w", err)
	}

	// Write file
	if err := os.WriteFile(agentsFile, []byte(content), filePerm); err != nil {
		return fmt.Errorf("failed to write AGENTS.md: %w", err)
	}

	result.CreatedFiles = append(result.CreatedFiles, "spectr/AGENTS.md")

	return nil
}

// configureTools configures the selected tools
func (e *InitExecutor) configureTools(
	selectedToolIDs []string,
	spectrDir string,
	result *ExecutionResult,
) error {
	if len(selectedToolIDs) == 0 {
		return nil // No tools to configure
	}

	for _, toolID := range selectedToolIDs {
		tool, err := e.registry.GetTool(toolID)
		if err != nil {
			result.Errors = append(
				result.Errors,
				fmt.Sprintf("tool %s not found: %v", toolID, err),
			)

			continue
		}

		configurator := e.getConfigurator(toolID)
		if configurator == nil {
			result.Errors = append(
				result.Errors,
				fmt.Sprintf(
					"no configurator found for tool: %s",
					toolID,
				),
			)

			continue
		}

		// Check if already configured
		wasConfigured := configurator.IsConfigured(e.projectPath)

		// Configure the tool
		if err := configurator.Configure(
			e.projectPath,
			spectrDir,
		); err != nil {
			result.Errors = append(
				result.Errors,
				fmt.Sprintf(
					"failed to configure %s: %v",
					tool.Name,
					err,
				),
			)

			continue
		}

		// Track created/updated files
		fileInfo := e.getToolFileInfo(tool)
		if wasConfigured {
			result.UpdatedFiles = append(result.UpdatedFiles, fileInfo...)
		} else {
			result.CreatedFiles = append(result.CreatedFiles, fileInfo...)
		}

		// Auto-install slash commands if this tool has a mapping
		if slashToolID, hasMapping := GetSlashToolMapping(toolID); hasMapping {
			slashConfigurator := e.getConfigurator(slashToolID)
			if slashConfigurator == nil {
				result.Errors = append(
					result.Errors,
					fmt.Sprintf(
						"no slash configurator found for: %s",
						slashToolID,
					),
				)

				continue
			}

			// Check if slash commands already configured
			slashWasConfigured := slashConfigurator.IsConfigured(e.projectPath)

			// Configure slash commands
			if err := slashConfigurator.Configure(
				e.projectPath,
				spectrDir,
			); err != nil {
				result.Errors = append(
					result.Errors,
					fmt.Sprintf(
						"failed to configure slash commands for %s: %v",
						tool.Name,
						err,
					),
				)

				continue
			}

			// Track slash command files
			slashFileInfo := e.getSlashCommandFileInfo(slashToolID)
			if slashWasConfigured {
				result.UpdatedFiles = append(result.UpdatedFiles, slashFileInfo...)
			} else {
				result.CreatedFiles = append(result.CreatedFiles, slashFileInfo...)
			}
		}
	}

	return nil
}

// getConfigurator returns the configurator for a tool ID
func (_e *InitExecutor) getConfigurator(toolID string) Configurator {
	switch toolID {
	// Config-based tools
	case "claude-code":
		return &ClaudeCodeConfigurator{}
	case "cline":
		return &ClineConfigurator{}
	case "costrict-config":
		return &CostrictConfigurator{}
	case "qoder-config":
		return &QoderConfigurator{}
	case "codebuddy":
		return &CodeBuddyConfigurator{}
	case "qwen":
		return &QwenConfigurator{}
	case "antigravity":
		return &AntigravityConfigurator{}

	// Slash command tools
	case "claude":
		return NewClaudeSlashConfigurator()
	case "cline-slash":
		return NewClineSlashConfigurator()
	case "kilocode":
		return NewKilocodeSlashConfigurator()
	case "qoder-slash":
		return NewQoderSlashConfigurator()
	case "cursor":
		return NewCursorSlashConfigurator()
	case "aider":
		return NewAiderSlashConfigurator()
	case "continue":
		return NewContinueSlashConfigurator()
	case "copilot":
		return NewCopilotSlashConfigurator()
	case "mentat":
		return NewMentatSlashConfigurator()
	case "tabnine":
		return NewTabnineSlashConfigurator()
	case "smol":
		return NewSmolSlashConfigurator()
	case "costrict-slash":
		return NewCostrictSlashConfigurator()
	case "windsurf":
		return NewWindsurfSlashConfigurator()
	case "codebuddy-slash":
		return NewCodeBuddySlashConfigurator()
	case "qwen-slash":
		return NewQwenSlashConfigurator()
	case "antigravity-slash":
		return NewAntigravitySlashConfigurator()

	default:
		return nil
	}
}

// getToolFileInfo returns the files that would be created/updated for a tool
func (e *InitExecutor) getToolFileInfo(tool *ToolDefinition) []string {
	switch tool.Type {
	case ToolTypeConfig:
		// Config-based tools create a single instruction file
		// Get the actual file path from the tool ID mapping
		var filePath string
		switch tool.ID {
		case "claude-code":
			filePath = "CLAUDE.md"
		case "cline":
			filePath = "CLINE.md"
		case "costrict-config":
			filePath = "COSTRICT.md"
		case "qoder-config":
			filePath = "QODER.md"
		case "codebuddy":
			filePath = "CODEBUDDY.md"
		case "qwen":
			filePath = "QWEN.md"
		case "antigravity":
			filePath = "AGENTS.md"
		default:
			return make([]string, 0)
		}

		return []string{filePath}

	case ToolTypeSlash:
		// Slash command tools create 3 files
		configurator := e.getConfigurator(tool.ID)
		slashConfig, ok := configurator.(*SlashCommandConfigurator)
		if !ok {
			return make([]string, 0)
		}
		files := make([]string, 0)
		for _, path := range slashConfig.config.FilePaths {
			files = append(files, path)
		}

		return files
	}

	return make([]string, 0)
}

// getSlashCommandFileInfo returns the files for a slash command tool ID
func (e *InitExecutor) getSlashCommandFileInfo(slashToolID string) []string {
	configurator := e.getConfigurator(slashToolID)
	slashConfig, ok := configurator.(*SlashCommandConfigurator)
	if !ok {
		return make([]string, 0)
	}

	files := make([]string, 0)
	for _, path := range slashConfig.config.FilePaths {
		files = append(files, path)
	}

	return files
}

// FormatNextStepsMessage returns a formatted next steps message for display after initialization
func FormatNextStepsMessage() string {
	return `
────────────────────────────────────────────────────────────────

Next steps:

1. Populate your project context by telling your AI assistant:

   "Review spectr/project.md and help me fill in our project's tech stack,
   conventions, and description. Ask me questions to understand the codebase."

2. Create your first change proposal by saying:

   "Help me create a change proposal for [YOUR FEATURE HERE]. Walk me through
   the process and ask questions to understand the requirements."

3. Learn the Spectr workflow:

   "Review spectr/AGENTS.md and explain how Spectr's change workflow works."

────────────────────────────────────────────────────────────────
`
}

// createReadmeIfMissing creates a basic README.md if it doesn't exist
func (e *InitExecutor) createReadmeIfMissing(result *ExecutionResult) error {
	readmePath := filepath.Join(e.projectPath, "README.md")

	// Only create if it doesn't exist
	if FileExists(readmePath) {
		return nil
	}

	// Get project name
	projectName := filepath.Base(e.projectPath)

	content := fmt.Sprintf(`# %s

This project uses [Spectr](https://spectr.dev) for structured development and change management.

## Getting Started

1. Review the project documentation in `+"`spectr/project.md`"+`
2. Explore the Spectr documentation: https://spectr.dev
3. Create your first change proposal: `+"`spectr proposal <change-name>`"+`

## Spectr Commands

- `+"`spectr proposal <name>`"+` - Create a new change proposal
- `+"`spectr apply <change-id>`"+` - Apply an approved change
- `+"`spectr archive <change-id>`"+` - Archive a deployed change

## Documentation

- [Project Overview](spectr/project.md)
- [AI Agent Instructions](spectr/AGENTS.md)
- [Specifications](spectr/specs/)
- [Change Proposals](spectr/changes/)
`, projectName)

	if err := os.WriteFile(readmePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create README.md: %w", err)
	}

	result.CreatedFiles = append(result.CreatedFiles, "README.md")

	return nil
}
