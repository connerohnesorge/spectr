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
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
	"github.com/spf13/afero"
)

// InitExecutor handles the actual initialization process
type InitExecutor struct {
	projectPath string
	tm          *TemplateManager
	projectFs   afero.Fs // Filesystem rooted at project directory
	homeFs      afero.Fs // Filesystem rooted at home directory
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

	// Create dual filesystem: projectFs rooted at project, homeFs rooted at home directory
	// Task 8.2: Fail if os.UserHomeDir() errors
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	projectFs := afero.NewBasePathFs(afero.NewOsFs(), projectPath)
	homeFs := afero.NewBasePathFs(afero.NewOsFs(), homeDir)

	return &InitExecutor{
		projectPath: projectPath,
		tm:          tm,
		projectFs:   projectFs,
		homeFs:      homeFs,
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
		defaultTemplateContext(),
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
// Uses dual filesystems, collects initializers, deduplicates, sorts by type, and executes with fail-fast.
func (e *InitExecutor) configureProviders(
	selectedProviderIDs []string,
	spectrDir string,
	result *ExecutionResult,
) error {
	if len(selectedProviderIDs) == 0 {
		return nil // No providers to configure
	}

	// Create config with relative path (projectFs is already rooted at project directory)
	cfg := &domain.Config{SpectrDir: "spectr"}
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// Task 8.3: Get all registered providers sorted by priority
	allRegistrations := providers.RegisteredProviders()

	// Task 8.4: Collect initializers from selected providers
	// Process in priority order (allRegistrations is already sorted)
	var allInitializers []domain.Initializer
	for _, reg := range allRegistrations {
		// Only process selected providers
		isSelected := false
		for _, id := range selectedProviderIDs {
			if id == reg.ID {
				isSelected = true

				break
			}
		}
		if !isSelected {
			continue
		}

		// Get initializers from provider
		inits := reg.Provider.Initializers(context.Background(), e.tm)
		allInitializers = append(allInitializers, inits...)
	}

	// Task 8.7: Sort initializers by type priority (stable sort preserves provider order)
	sortInitializersByType(allInitializers)

	// Task 8.6: Deduplicate initializers (keep first occurrence)
	deduplicatedInits := dedupeInitializers(allInitializers)

	// Task 8.8, 8.9, 8.10: Execute initializers with fail-fast, merge results inline
	for _, init := range deduplicatedInits {
		initResult, err := init.Init(context.Background(), e.projectFs, e.homeFs, cfg, e.tm)
		if err != nil {
			// Task 8.10: Fail-fast - stop on first error, return partial result
			return fmt.Errorf("initializer failed: %w", err)
		}
		// Task 8.9: Merge ExecutionResults inline
		result.CreatedFiles = append(result.CreatedFiles, initResult.CreatedFiles...)
		result.UpdatedFiles = append(result.UpdatedFiles, initResult.UpdatedFiles...)
	}

	return nil
}

// Task 8.5: templateContextFromConfig derives TemplateContext from Config.SpectrDir
func templateContextFromConfig(cfg *domain.Config) domain.TemplateContext {
	return domain.TemplateContext{
		BaseDir:     cfg.SpectrDir,
		SpecsDir:    cfg.SpecsDir(),
		ChangesDir:  cfg.ChangesDir(),
		ProjectFile: cfg.ProjectFile(),
		AgentsFile:  cfg.AgentsFile(),
	}
}

// Task 8.7: sortInitializersByType sorts initializers by type priority using stable sort.
// Priority 1: DirectoryInitializer, HomeDirectoryInitializer
// Priority 2: ConfigFileInitializer
// Priority 3: SlashCommandsInitializer, HomeSlashCommandsInitializer, PrefixedSlashCommandsInitializer,
//
//	HomePrefixedSlashCommandsInitializer, TOMLSlashCommandsInitializer
func sortInitializersByType(inits []domain.Initializer) {
	sort.SliceStable(inits, func(i, j int) bool {
		return initializerPriority(inits[i]) < initializerPriority(inits[j])
	})
}

// initializerPriority returns the type-based priority for an initializer.
func initializerPriority(init domain.Initializer) int {
	switch init.(type) {
	case *initializers.DirectoryInitializer, *initializers.HomeDirectoryInitializer:
		return 1
	case *initializers.ConfigFileInitializer:
		return 2
	case *initializers.SlashCommandsInitializer, *initializers.HomeSlashCommandsInitializer,
		*initializers.PrefixedSlashCommandsInitializer, *initializers.HomePrefixedSlashCommandsInitializer,
		*initializers.TOMLSlashCommandsInitializer:
		return 3
	default:
		return 99
	}
}

// Task 8.6: dedupeInitializers deduplicates initializers using the optional Deduplicatable interface.
// Keeps first occurrence when duplicates are found.
func dedupeInitializers(inits []domain.Initializer) []domain.Initializer {
	seen := make(map[string]bool)
	result := make([]domain.Initializer, 0, len(inits))

	for _, init := range inits {
		// Check if initializer supports deduplication via the Deduplicatable interface
		if d, ok := init.(initializers.Deduplicatable); ok {
			key := d.DedupeKey()
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

// defaultTemplateContext returns a default TemplateContext with standard paths.
// This is used for rendering templates during initialization.
func defaultTemplateContext() *domain.TemplateContext {
	return &domain.TemplateContext{
		BaseDir:     "spectr",
		SpecsDir:    "spectr/specs",
		ChangesDir:  "spectr/changes",
		ProjectFile: "spectr/project.md",
		AgentsFile:  "spectr/AGENTS.md",
	}
}
