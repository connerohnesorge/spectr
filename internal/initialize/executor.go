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
	"strings"

	"github.com/spf13/afero"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/connerohnesorge/spectr/internal/initialize/providers"
)

// InitExecutor handles the actual initialization process using the new provider architecture
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

// Execute runs the initialization process using the new provider architecture
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

	// 2. Create dual filesystem (Task 8.2)
	projectFs := afero.NewBasePathFs(afero.NewOsFs(), e.projectPath)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return result, fmt.Errorf("failed to get home directory: %w", err)
	}
	homeFs := afero.NewBasePathFs(afero.NewOsFs(), homeDir)

	// 3. Create spectr/ directory structure
	spectrDir := "spectr" // Relative to projectFs root
	err = e.createDirectoryStructure(
		projectFs,
		spectrDir,
		result,
	)
	if err != nil {
		return result, fmt.Errorf(
			"failed to create directory structure: %w",
			err,
		)
	}

	// 4. Create project.md
	err = e.createProjectMd(projectFs, spectrDir, result)
	if err != nil {
		result.Errors = append(
			result.Errors,
			fmt.Sprintf(
				"failed to create project.md: %v",
				err,
			),
		)
	}

	// 5. Create AGENTS.md
	err = e.createAgentsMd(projectFs, spectrDir, result)
	if err != nil {
		result.Errors = append(
			result.Errors,
			fmt.Sprintf(
				"failed to create AGENTS.md: %v",
				err,
			),
		)
	}

	// 6. Configure selected providers using new architecture (Tasks 8.3-8.10)
	providerResult, err := e.configureProviders(
		selectedProviderIDs,
		projectFs,
		homeFs,
		spectrDir,
	)

	// Merge provider results into main result (even on error, to preserve partial results)
	result.CreatedFiles = append(result.CreatedFiles, providerResult.CreatedFiles...)
	result.UpdatedFiles = append(result.UpdatedFiles, providerResult.UpdatedFiles...)

	if err != nil {
		return result, fmt.Errorf("failed to configure providers: %w", err)
	}

	// 7. Create CI workflow if enabled
	if ciWorkflowEnabled {
		err = e.createCIWorkflow(projectFs, result)
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
	projectFs afero.Fs,
	spectrDir string,
	result *ExecutionResult,
) error {
	dirs := []string{
		spectrDir,
		filepath.Join(spectrDir, "specs"),
		filepath.Join(spectrDir, "changes"),
	}

	for _, dir := range dirs {
		exists, err := afero.DirExists(projectFs, dir)
		if err != nil {
			return fmt.Errorf(
				"failed to check directory %s: %w",
				dir,
				err,
			)
		}

		if !exists {
			if err := projectFs.MkdirAll(dir, 0755); err != nil {
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
	projectFs afero.Fs,
	spectrDir string,
	result *ExecutionResult,
) error {
	projectFile := filepath.Join(
		spectrDir,
		"project.md",
	)

	// Check if it already exists
	exists, err := afero.Exists(projectFs, projectFile)
	if err != nil {
		return fmt.Errorf("failed to check project.md: %w", err)
	}
	if exists {
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
	if err := afero.WriteFile(projectFs, projectFile, []byte(content), filePerm); err != nil {
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
	projectFs afero.Fs,
	spectrDir string,
	result *ExecutionResult,
) error {
	agentsFile := filepath.Join(
		spectrDir,
		"AGENTS.md",
	)

	// Check if it already exists
	exists, err := afero.Exists(projectFs, agentsFile)
	if err != nil {
		return fmt.Errorf("failed to check AGENTS.md: %w", err)
	}
	if exists {
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
	if err := afero.WriteFile(projectFs, agentsFile, []byte(content), filePerm); err != nil {
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

// deduplicatable is an optional interface for initializers that support deduplication (Task 8.6)
type deduplicatable interface {
	dedupeKey() string
}

// configureProviders configures the selected providers using the new interface-driven architecture (Tasks 8.3-8.10)
func (e *InitExecutor) configureProviders(
	selectedProviderIDs []string,
	projectFs, homeFs afero.Fs,
	spectrDir string,
) (providers.ExecutionResult, error) {
	if len(selectedProviderIDs) == 0 {
		return providers.ExecutionResult{}, nil
	}

	ctx := context.Background()
	cfg := &providers.Config{
		SpectrDir: spectrDir,
	}

	// Validate config
	if err := cfg.Validate(); err != nil {
		return providers.ExecutionResult{}, fmt.Errorf("invalid config: %w", err)
	}

	// Task 8.3: Use RegisteredProviders() for sorted provider list
	allRegistrations := providers.RegisteredProviders()

	// Filter to only selected providers (preserve priority order)
	var selectedRegistrations []providers.Registration
	for _, reg := range allRegistrations {
		for _, id := range selectedProviderIDs {
			if reg.ID == id {
				selectedRegistrations = append(selectedRegistrations, reg)

				break
			}
		}
	}

	// Task 8.4: Collect initializers from selected providers
	var allInitializers []providers.Initializer
	for _, reg := range selectedRegistrations {
		inits := reg.Provider.Initializers(ctx, e.tm)
		allInitializers = append(allInitializers, inits...)
	}

	// Task 8.7: Sort initializers by type (stable sort preserves provider priority order)
	allInitializers = sortInitializers(allInitializers)

	// Task 8.6: Deduplicate initializers (keep first occurrence = highest priority provider wins)
	allInitializers = dedupeInitializers(allInitializers)

	// Task 8.8, 8.9, 8.10: Execute initializers with fail-fast behavior
	allResults := make([]providers.InitResult, 0, len(allInitializers))

	for _, init := range allInitializers {
		result, err := init.Init(ctx, projectFs, homeFs, cfg, e.tm)
		if err != nil {
			// Task 8.10: Fail-fast - stop on first error, return partial results
			partialResult := aggregateResults(allResults)

			return partialResult, fmt.Errorf("initializer failed: %w", err)
		}
		allResults = append(allResults, result)
	}

	// Task 8.9: Aggregate all results on success
	return aggregateResults(allResults), nil
}

// sortInitializers sorts initializers by type priority (Task 8.7)
// Uses stable sort to preserve provider priority order within each type
func sortInitializers(all []providers.Initializer) []providers.Initializer {
	sorted := make([]providers.Initializer, len(all))
	copy(sorted, all)

	sort.SliceStable(sorted, func(i, j int) bool {
		return initializerPriority(sorted[i]) < initializerPriority(sorted[j])
	})

	return sorted
}

// initializerPriority returns the priority for initializer ordering (Task 8.7)
func initializerPriority(init providers.Initializer) int {
	// Use type assertion to determine initializer type
	// We need to check the concrete types from the providers package
	typeName := fmt.Sprintf("%T", init)

	switch {
	// Priority 1: Directory initializers
	case strings.Contains(typeName, "DirectoryInitializer") || strings.Contains(typeName, "HomeDirectoryInitializer"):
		return 1
	// Priority 2: Config file initializers
	case strings.Contains(typeName, "ConfigFileInitializer"):
		return 2
	// Priority 3: Slash command initializers
	case strings.Contains(typeName, "SlashCommandsInitializer") ||
		strings.Contains(typeName, "HomeSlashCommandsInitializer") ||
		strings.Contains(typeName, "PrefixedSlashCommandsInitializer") ||
		strings.Contains(typeName, "HomePrefixedSlashCommandsInitializer") ||
		strings.Contains(typeName, "TOMLSlashCommandsInitializer"):
		return 3
	default:
		return 99 // Unknown types go last
	}
}

// dedupeInitializers removes duplicate initializers (Task 8.6)
// Keeps first occurrence when duplicates are found
func dedupeInitializers(all []providers.Initializer) []providers.Initializer {
	seen := make(map[string]bool)
	result := make([]providers.Initializer, 0, len(all))

	for _, init := range all {
		// Check if initializer supports deduplication
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

// aggregateResults combines multiple InitResult values into a single ExecutionResult (Task 8.9)
func aggregateResults(results []providers.InitResult) providers.ExecutionResult {
	var created, updated []string

	for _, r := range results {
		created = append(created, r.CreatedFiles...)
		updated = append(updated, r.UpdatedFiles...)
	}

	return providers.ExecutionResult{
		CreatedFiles: created,
		UpdatedFiles: updated,
	}
}

// createCIWorkflow creates the .github/workflows/spectr-ci.yml file
func (e *InitExecutor) createCIWorkflow(
	projectFs afero.Fs,
	result *ExecutionResult,
) error {
	// Ensure .github/workflows directory exists
	workflowDir := filepath.Join(
		".github",
		"workflows",
	)
	if err := projectFs.MkdirAll(workflowDir, 0755); err != nil {
		return fmt.Errorf(
			"failed to create workflows directory: %w",
			err,
		)
	}

	workflowFile := filepath.Join(
		workflowDir,
		"spectr-ci.yml",
	)
	wasConfigured, err := afero.Exists(projectFs, workflowFile)
	if err != nil {
		return fmt.Errorf("failed to check workflow file: %w", err)
	}

	// Render the CI workflow template
	content, err := e.tm.RenderCIWorkflow()
	if err != nil {
		return fmt.Errorf(
			"failed to render CI workflow template: %w",
			err,
		)
	}

	// Write the workflow file
	if err := afero.WriteFile(projectFs, workflowFile, []byte(content), filePerm); err != nil {
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
