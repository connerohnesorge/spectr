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

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/connerohnesorge/spectr/internal/initialize/providers"
	"github.com/connerohnesorge/spectr/internal/templates"
	"github.com/spf13/afero"
)

// InitExecutor handles the actual initialization process
type InitExecutor struct {
	projectPath string
	tm          *templates.TemplateManager
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
	tm, err := templates.NewTemplateManager()
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

	// 5. Configure selected providers using new architecture
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
		templates.ProjectContext{
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
		domain.DefaultTemplateContext(),
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
// This implementation uses the composable initializer pattern.
func (e *InitExecutor) configureProviders(
	selectedProviderIDs []string,
	spectrDir string,
	result *ExecutionResult,
) error {
	if len(selectedProviderIDs) == 0 {
		return nil // No providers to configure
	}

	// Task 8.2: Create dual filesystem
	// Project filesystem rooted at project directory
	projectFs := afero.NewBasePathFs(afero.NewOsFs(), e.projectPath)

	// Home filesystem rooted at home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	homeFs := afero.NewBasePathFs(afero.NewOsFs(), homeDir)

	// Create Config from spectrDir
	cfg := &providers.Config{
		SpectrDir: "spectr", // relative to project root
	}
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// Task 8.5: Create TemplateContext from Config
	tmplCtx := templateContextFromConfig(cfg)
	_ = tmplCtx // Will be used when rendering templates in initializers

	// Task 8.3: Use RegisteredProviders() to get sorted provider list
	allRegistrations := providers.RegisteredProviders()

	// Filter to only selected providers (maintain priority order)
	var selectedProviders []providers.Registration
	for _, reg := range allRegistrations {
		for _, id := range selectedProviderIDs {
			if reg.ID == id {
				selectedProviders = append(selectedProviders, reg)

				break
			}
		}
	}

	// Task 8.4: Collect initializers from selected providers
	ctx := context.Background()
	var allInitializers []providers.Initializer

	for _, reg := range selectedProviders {
		inits := reg.Provider.Initializers(ctx, e.tm)
		allInitializers = append(allInitializers, inits...)
	}

	// Task 8.7: Sort initializers by type priority
	sortedInitializers := sortInitializers(allInitializers)

	// Task 8.6: Deduplicate initializers (keep first occurrence)
	deduplicatedInitializers := deduplicateInitializers(sortedInitializers)

	// Task 8.8, 8.9, 8.10: Execute initializers with fail-fast behavior
	initResults := make([]providers.InitResult, 0, len(deduplicatedInitializers))

	for _, init := range deduplicatedInitializers {
		// Pass both filesystems and TemplateManager to initializer
		initResult, err := init.Init(ctx, projectFs, homeFs, cfg, e.tm)
		if err != nil {
			// Task 8.10: Fail-fast - stop on first error
			// Return partial results from successful initializers
			partialExecResult := providers.AggregateResults(initResults)
			result.CreatedFiles = append(result.CreatedFiles, partialExecResult.CreatedFiles...)
			result.UpdatedFiles = append(result.UpdatedFiles, partialExecResult.UpdatedFiles...)

			return fmt.Errorf("initializer failed: %w", err)
		}
		initResults = append(initResults, initResult)
	}

	// Task 8.9: Aggregate all results on success
	execResult := providers.AggregateResults(initResults)
	result.CreatedFiles = append(result.CreatedFiles, execResult.CreatedFiles...)
	result.UpdatedFiles = append(result.UpdatedFiles, execResult.UpdatedFiles...)

	return nil
}

// templateContextFromConfig derives TemplateContext from Config.SpectrDir (Task 8.5)
func templateContextFromConfig(cfg *providers.Config) domain.TemplateContext {
	return domain.TemplateContext{
		BaseDir:     cfg.SpectrDir,
		SpecsDir:    cfg.SpecsDir(),
		ChangesDir:  cfg.ChangesDir(),
		ProjectFile: cfg.ProjectFile(),
		AgentsFile:  cfg.AgentsFile(),
	}
}

// sortInitializers sorts initializers by type priority (Task 8.7)
// Order: Directory (1) -> ConfigFile (2) -> SlashCommands (3)
func sortInitializers(inits []providers.Initializer) []providers.Initializer {
	sorted := make([]providers.Initializer, len(inits))
	copy(sorted, inits)

	sort.SliceStable(sorted, func(i, j int) bool {
		return initializerPriority(sorted[i]) < initializerPriority(sorted[j])
	})

	return sorted
}

// initializerPriority returns the execution priority for an initializer type
func initializerPriority(init providers.Initializer) int {
	switch init.(type) {
	case *providers.DirectoryInitializer, *providers.HomeDirectoryInitializer:
		return 1
	case *providers.ConfigFileInitializer:
		return 2
	case *providers.SlashCommandsInitializer, *providers.HomeSlashCommandsInitializer,
		*providers.PrefixedSlashCommandsInitializer, *providers.HomePrefixedSlashCommandsInitializer,
		*providers.TOMLSlashCommandsInitializer:
		return 3
	default:
		return 99 // Unknown types go last
	}
}

// deduplicatable is the optional interface for initializers that support deduplication
type deduplicatable interface {
	dedupeKey() string
}

// deduplicateInitializers removes duplicate initializers, keeping first occurrence (Task 8.6)
func deduplicateInitializers(inits []providers.Initializer) []providers.Initializer {
	seen := make(map[string]bool)
	result := make([]providers.Initializer, 0, len(inits))

	for _, init := range inits {
		// Check if initializer implements deduplicatable interface
		if d, ok := init.(deduplicatable); ok {
			key := d.dedupeKey()
			if seen[key] {
				continue // Skip duplicate
			}
			seen[key] = true
		}
		result = append(result, init)
	}

	return result
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
