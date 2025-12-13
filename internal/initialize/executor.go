// Package initialize provides utilities for initializing Spectr
// in a project directory.
//
//nolint:revive // file-length-limit - logically cohesive, no benefit to split
package initialize

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/connerohnesorge/spectr/internal/initialize/providers"
)

// InitExecutor handles the actual initialization process
type InitExecutor struct {
	projectPath string
	tm          *TemplateManager
}

// NewInitExecutor creates a new initialization executor
func NewInitExecutor(
	cmd *InitCmd,
) (*InitExecutor, error) {
	// Use the resolved path from InitCmd
	projectPath := cmd.Path
	if projectPath == "" {
		return nil, fmt.Errorf(
			"project path is required",
		)
	}

	// Check if path exists
	if !FileExists(projectPath) {
		return nil,
			fmt.Errorf(
				"project path does not exist: %s",
				projectPath,
			)
	}

	// Initialize template manager
	tm, err := NewTemplateManager()
	if err != nil {
		return nil,
			fmt.Errorf(
				"failed to initialize template manager: %w",
				err,
			)
	}

	return &InitExecutor{
		projectPath: projectPath,
		tm:          tm,
	}, nil
}

// Execute runs the initialization process
func (e *InitExecutor) Execute(
	selectedProviderIDs []string,
	ciWorkflowEnabled bool,
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
	spectrDir := filepath.Join(
		e.projectPath,
		"spectr",
	)
	err := e.createDirectoryStructure(
		spectrDir,
		result,
	)
	if err != nil {
		return result, fmt.Errorf(
			"failed to create directory structure: %w",
			err,
		)
	}

	// 3. Create project.md
	err = e.createProjectMd(spectrDir, result)
	if err != nil {
		result.Errors = append(
			result.Errors,
			fmt.Sprintf(
				"failed to create project.md: %v",
				err,
			),
		)
	}

	// 4. Create AGENTS.md
	err = e.createAgentsMd(spectrDir, result)
	if err != nil {
		result.Errors = append(
			result.Errors,
			fmt.Sprintf(
				"failed to create AGENTS.md: %v",
				err,
			),
		)
	}

	// 5. Configure selected providers
	err = e.configureProviders(
		selectedProviderIDs,
		spectrDir,
		result,
	)
	if err != nil {
		result.Errors = append(
			result.Errors,
			fmt.Sprintf(
				"failed to configure tools: %v",
				err,
			),
		)
	}

	// 6. Create CI workflow if enabled
	if ciWorkflowEnabled {
		err = e.createCIWorkflow(result)
		if err != nil {
			result.Errors = append(
				result.Errors,
				fmt.Sprintf(
					"failed to create CI workflow: %v",
					err,
				),
			)
		}
	}

	return result, nil
}

// createDirectoryStructure creates the spectr/ directory
// and subdirectories
func (*InitExecutor) createDirectoryStructure(
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
			result.CreatedFiles = append(
				result.CreatedFiles,
				dir+"/",
			)
		}
	}

	return nil
}

// createProjectMd creates the project.md file
func (e *InitExecutor) createProjectMd(
	spectrDir string,
	result *ExecutionResult,
) error {
	projectFile := filepath.Join(
		spectrDir,
		"project.md",
	)

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
	content, err := e.tm.RenderProject(
		ProjectContext{
			ProjectName: projectName,
			Description: "Add your project description here",
			TechStack: []string{
				"Add",
				"Your",
				"Technologies",
				"Here",
			},
			Conventions: "",
		},
	)
	if err != nil {
		return fmt.Errorf(
			"failed to render project template: %w",
			err,
		)
	}

	// Write file
	if err := os.WriteFile(projectFile, []byte(content), filePerm); err != nil {
		return fmt.Errorf(
			"failed to write project.md: %w",
			err,
		)
	}

	result.CreatedFiles = append(
		result.CreatedFiles,
		"spectr/project.md",
	)

	return nil
}

// createAgentsMd creates the AGENTS.md file
func (e *InitExecutor) createAgentsMd(
	spectrDir string,
	result *ExecutionResult,
) error {
	agentsFile := filepath.Join(
		spectrDir,
		"AGENTS.md",
	)

	// Check if it already exists
	if FileExists(agentsFile) {
		result.Errors = append(
			result.Errors,
			"AGENTS.md already exists, skipping",
		)

		return nil
	}

	// Render template
	content, err := e.tm.RenderAgents(
		providers.DefaultTemplateContext(),
	)
	if err != nil {
		return fmt.Errorf(
			"failed to render agents template: %w",
			err,
		)
	}

	// Write file
	if err := os.WriteFile(agentsFile, []byte(content), filePerm); err != nil {
		return fmt.Errorf(
			"failed to write AGENTS.md: %w",
			err,
		)
	}

	result.CreatedFiles = append(
		result.CreatedFiles,
		"spectr/AGENTS.md",
	)

	return nil
}

// configureProviders configures the selected providers using the new interface-driven architecture.
// Each provider handles both its instruction file AND slash commands in a single Configure() call.
func (e *InitExecutor) configureProviders(
	selectedProviderIDs []string,
	spectrDir string,
	result *ExecutionResult,
) error {
	if len(selectedProviderIDs) == 0 {
		return nil // No providers to configure
	}

	for _, providerID := range selectedProviderIDs {
		provider := providers.Get(providerID)
		if provider == nil {
			result.Errors = append(
				result.Errors,
				fmt.Sprintf(
					"provider %s not found",
					providerID,
				),
			)

			continue
		}

		// Check if already configured
		wasConfigured := provider.IsConfigured(
			e.projectPath,
		)

		// Configure the provider using its initializers
		if err := providers.ConfigureInitializers(provider.Initializers(), e.projectPath, e.tm); err != nil {
			result.Errors = append(
				result.Errors,
				fmt.Sprintf(
					"failed to configure %s: %v",
					provider.Name(),
					err,
				),
			)

			continue
		}

		// Track created/updated files
		filePaths := provider.GetFilePaths()
		if wasConfigured {
			result.UpdatedFiles = append(
				result.UpdatedFiles,
				filePaths...)
		} else {
			result.CreatedFiles = append(result.CreatedFiles, filePaths...)
		}
	}

	return nil
}

// createCIWorkflow creates the .github/workflows/spectr-ci.yml file
func (e *InitExecutor) createCIWorkflow(
	result *ExecutionResult,
) error {
	// Ensure .github/workflows directory exists
	workflowDir := filepath.Join(
		e.projectPath,
		".github",
		"workflows",
	)
	if err := EnsureDir(workflowDir); err != nil {
		return fmt.Errorf(
			"failed to create workflows directory: %w",
			err,
		)
	}

	workflowFile := filepath.Join(
		workflowDir,
		"spectr-ci.yml",
	)
	wasConfigured := FileExists(workflowFile)

	// Render the CI workflow template
	content, err := e.tm.RenderCIWorkflow()
	if err != nil {
		return fmt.Errorf(
			"failed to render CI workflow template: %w",
			err,
		)
	}

	// Write the workflow file
	if err := os.WriteFile(workflowFile, []byte(content), filePerm); err != nil {
		return fmt.Errorf(
			"failed to write CI workflow file: %w",
			err,
		)
	}

	// Track in results
	if wasConfigured {
		result.UpdatedFiles = append(
			result.UpdatedFiles,
			".github/workflows/spectr-ci.yml",
		)
	} else {
		result.CreatedFiles = append(result.CreatedFiles, ".github/workflows/spectr-ci.yml")
	}

	return nil
}

// FormatNextStepsMessage returns a formatted next steps message for display after initialization
func FormatNextStepsMessage() string {
	return `
────────────────────────────────────────────────────────────────

Next steps:

1. Populate your project context by telling your AI assistant:

   "` + PopulateContextPrompt + `"

2. Create your first change proposal by saying:

   "Help me create a change proposal for [YOUR FEATURE HERE]. Walk me through
   the process and ask questions to understand the requirements."

3. Learn the Spectr workflow:

   "Review spectr/AGENTS.md and explain how Spectr's change workflow works."

────────────────────────────────────────────────────────────────
`
}
