// Package initialize provides utilities for initializing Spectr
// in a project directory.
//
//nolint:revive // file-length-limit - logically cohesive, no benefit to split
package initialize

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/connerohnesorge/spectr/internal/initialize/providers"
	"github.com/connerohnesorge/spectr/internal/initialize/templates"
	"github.com/spf13/afero"
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

// configureProviders configures the selected providers using the new initializer-based architecture.
// Collects initializers from all providers, deduplicates them, sorts by type, and executes them.
func (e *InitExecutor) configureProviders(
	selectedProviderIDs []string,
	spectrDir string,
	result *ExecutionResult,
) error {
	if len(selectedProviderIDs) == 0 {
		return nil // No providers to configure
	}

	// Task 6.1: Create dual filesystems
	// Project-relative filesystem (for CLAUDE.md, .claude/commands/, etc.)
	projectFs := afero.NewBasePathFs(
		afero.NewOsFs(),
		e.projectPath,
	)

	// Global filesystem (for ~/.config/tool/commands/, etc.)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf(
			"failed to get user home directory: %w",
			err,
		)
	}
	globalFs := afero.NewBasePathFs(
		afero.NewOsFs(),
		homeDir,
	)

	// Create configuration
	cfg := &providers.Config{
		SpectrDir: "spectr",
	}

	// Task 6.2 & 6.3: Collect initializers from selected providers
	var allInitializers []providers.Initializer
	for _, providerID := range selectedProviderIDs {
		reg, ok := providers.GetProvider(
			providerID,
		)
		if !ok {
			result.Errors = append(
				result.Errors,
				fmt.Sprintf(
					"provider %s not found",
					providerID,
				),
			)

			continue
		}

		// Collect initializers from this provider
		ctx := context.Background()
		inits := reg.Provider.Initializers(ctx)
		allInitializers = append(
			allInitializers,
			inits...)
	}

	// Task 6.4: Deduplicate initializers by path
	dedupedInitializers := dedupeInitializers(
		allInitializers,
	)

	// Task 6.5: Sort initializers by type
	sortedInitializers := sortInitializers(
		dedupedInitializers,
	)

	// Create a template manager adapter that satisfies providers.TemplateManager interface
	tmAdapter := &templateManagerAdapter{tm: e.tm}

	// Task 6.6: Execute initializers
	for _, init := range sortedInitializers {
		// Select filesystem based on initializer scope
		fs := projectFs
		if init.IsGlobal() {
			fs = globalFs
		}

		// Execute initializer
		initResult, err := init.Init(
			context.Background(),
			fs,
			cfg,
			tmAdapter,
		)
		// Task 6.8: Handle partial failures - continue with rest
		if err != nil {
			result.Errors = append(
				result.Errors,
				fmt.Sprintf(
					"failed to initialize %s: %v",
					init.Path(),
					err,
				),
			)

			continue
		}

		// Task 6.7: Aggregate InitResult into ExecutionResult
		result.CreatedFiles = append(
			result.CreatedFiles,
			initResult.CreatedFiles...)
		result.UpdatedFiles = append(
			result.UpdatedFiles,
			initResult.UpdatedFiles...)
	}

	return nil
}

// dedupeInitializers removes duplicate initializers by path.
// Same path = run once, even if multiple providers return it.
func dedupeInitializers(
	all []providers.Initializer,
) []providers.Initializer {
	seen := make(map[string]bool)
	var result []providers.Initializer
	for _, init := range all {
		key := init.Path()
		if !seen[key] {
			seen[key] = true
			result = append(result, init)
		}
	}

	return result
}

// sortInitializers sorts initializers by type to guarantee execution order:
// 1. DirectoryInitializer   - Create directories first
// 2. ConfigFileInitializer  - Then config files
// 3. SlashCommandsInitializer - Then slash commands
func sortInitializers(
	all []providers.Initializer,
) []providers.Initializer {
	sorted := make(
		[]providers.Initializer,
		len(all),
	)
	copy(sorted, all)

	sort.SliceStable(sorted, func(i, j int) bool {
		return initializerPriority(
			sorted[i],
		) < initializerPriority(
			sorted[j],
		)
	})

	return sorted
}

// initializerPriority returns the execution priority for an initializer type.
// Lower numbers execute first.
func initializerPriority(
	init providers.Initializer,
) int {
	switch init.(type) {
	case *providers.DirectoryInitializer:
		return 1
	case *providers.ConfigFileInitializer:
		return 2
	case *providers.SlashCommandsInitializer:
		return 3
	default:
		return 999 // Unknown types execute last
	}
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

// templateManagerAdapter adapts *TemplateManager to satisfy providers.TemplateManager interface.
// This is needed because the concrete TemplateManager returns templates.TemplateRef,
// but the interface requires interface{} to avoid import cycles.
type templateManagerAdapter struct {
	tm *TemplateManager
}

// RenderAgents delegates to the concrete TemplateManager
func (a *templateManagerAdapter) RenderAgents(
	ctx providers.TemplateContext,
) (string, error) {
	return a.tm.RenderAgents(ctx)
}

// RenderInstructionPointer delegates to the concrete TemplateManager
func (a *templateManagerAdapter) RenderInstructionPointer(
	ctx providers.TemplateContext,
) (string, error) {
	return a.tm.RenderInstructionPointer(ctx)
}

// RenderSlashCommand delegates to the concrete TemplateManager
func (a *templateManagerAdapter) RenderSlashCommand(
	commandType string,
	ctx providers.TemplateContext,
) (string, error) {
	return a.tm.RenderSlashCommand(
		commandType,
		ctx,
	)
}

// InstructionPointer returns the template reference as interface{}
func (a *templateManagerAdapter) InstructionPointer() interface{} {
	return a.tm.InstructionPointer()
}

// Agents returns the template reference as interface{}
func (a *templateManagerAdapter) Agents() interface{} {
	return a.tm.Agents()
}

// Project returns the template reference as interface{}
func (a *templateManagerAdapter) Project() interface{} {
	return a.tm.Project()
}

// CIWorkflow returns the template reference as interface{}
func (a *templateManagerAdapter) CIWorkflow() interface{} {
	return a.tm.CIWorkflow()
}

// SlashCommand returns the template reference as interface{}
func (a *templateManagerAdapter) SlashCommand(
	cmd interface{},
) interface{} {
	// Type assert cmd to templates.SlashCommand
	slashCmd, ok := cmd.(templates.SlashCommand)
	if !ok {
		// Return nil if conversion fails - will be caught by initializer
		return nil
	}

	return a.tm.SlashCommand(slashCmd)
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
