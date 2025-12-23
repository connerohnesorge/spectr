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

	"github.com/connerohnesorge/spectr/internal/initialize/git"
	"github.com/connerohnesorge/spectr/internal/initialize/providers"
	"github.com/spf13/afero"
)

// InitExecutor handles the actual initialization process
type InitExecutor struct {
	projectPath    string
	tm             *TemplateManager
	projectFs      afero.Fs            // project-relative filesystem
	globalFs       afero.Fs            // home directory filesystem
	changeDetector *git.ChangeDetector // git-based change detection
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

	// Task 6.2: Check if project is a git repo - fail early with clear error
	if !git.IsGitRepo(projectPath) {
		return nil, fmt.Errorf(
			"spectr init requires a git repository for change detection, " +
				"run 'git init' first, then retry 'spectr init'",
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

	// Task 6.1: Create dual filesystem instances
	// Project-relative filesystem (for CLAUDE.md, .claude/commands/, etc.)
	projectFs := afero.NewBasePathFs(afero.NewOsFs(), projectPath)

	// Global filesystem (for ~/.config/tool/commands/, etc.)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get user home directory: %w",
			err,
		)
	}
	globalFs := afero.NewBasePathFs(afero.NewOsFs(), homeDir)

	// Create change detector for git-based file change reporting
	changeDetector := git.NewChangeDetector(projectPath)

	return &InitExecutor{
		projectPath:    projectPath,
		tm:             tm,
		projectFs:      projectFs,
		globalFs:       globalFs,
		changeDetector: changeDetector,
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

	// Task 6.8: Take a git snapshot before initialization
	beforeSnapshot, err := e.changeDetector.Snapshot()
	if err != nil {
		// Non-fatal: continue without git tracking if snapshot fails
		result.Errors = append(result.Errors,
			fmt.Sprintf("git snapshot failed (continuing without change tracking): %v", err))
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
	err = e.createDirectoryStructure(
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

	// 5. Configure selected providers using new initializer-based architecture
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

	// Task 6.9: Get changed files from git after initialization
	if beforeSnapshot != "" {
		changedFiles, gitErr := e.changeDetector.ChangedFiles(beforeSnapshot)
		if gitErr != nil {
			result.Errors = append(result.Errors,
				fmt.Sprintf("git change detection failed: %v", gitErr))
		} else {
			// Update result with git-detected files
			e.updateResultWithGitChanges(result, changedFiles)
		}
	}

	return result, nil
}

// updateResultWithGitChanges updates the ExecutionResult with git-detected changed files.
// This replaces the declared paths from providers with actual changes detected by git.
func (e *InitExecutor) updateResultWithGitChanges(result *ExecutionResult, changedFiles []string) {
	// Track which files were in the original created/updated lists
	originalCreated := make(map[string]bool)
	for _, f := range result.CreatedFiles {
		originalCreated[f] = true
	}

	originalUpdated := make(map[string]bool)
	for _, f := range result.UpdatedFiles {
		originalUpdated[f] = true
	}

	// Clear and rebuild the lists based on git detection
	// Files that are new to git (untracked before) go to CreatedFiles
	// Files that existed and changed go to UpdatedFiles
	newCreated := make([]string, 0)
	newUpdated := make([]string, 0)

	for _, file := range changedFiles {
		// Check if this file was in the original created list (new file)
		// or was already being tracked
		switch {
		case originalCreated[file]:
			newCreated = append(newCreated, file)
		case originalUpdated[file]:
			newUpdated = append(newUpdated, file)
		default:
			// File was detected by git but not in our tracking
			// This could be a file created by an initializer that we didn't track
			// Add it based on whether the file existed before (heuristic: if file ends with /)
			// For simplicity, add all git-detected changes to created
			newCreated = append(newCreated, file)
		}
	}

	// Only update if we have git-detected changes
	if len(newCreated) > 0 || len(newUpdated) > 0 {
		result.CreatedFiles = newCreated
		result.UpdatedFiles = newUpdated
	}
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
// It uses the Registry API with Registration-based retrieval and runs initializers from each provider.
func (e *InitExecutor) configureProviders(
	selectedProviderIDs []string,
	spectrDir string,
	result *ExecutionResult,
) error {
	if len(selectedProviderIDs) == 0 {
		return nil // No providers to configure
	}

	// Create config for initializers
	cfg := providers.NewConfig(spectrDir)

	ctx := context.Background()

	// Task 6.4: Collect initializers from selected providers
	allInitializers := e.collectInitializers(ctx, selectedProviderIDs, result)

	// Task 6.5: Deduplicate initializers by Path()
	dedupedInitializers := e.dedupeInitializers(allInitializers)

	// Task 6.6: Sort initializers by type (Directory -> ConfigFile -> SlashCommands)
	sortedInitializers := e.sortInitializers(dedupedInitializers)

	// Task 6.7 & 6.10: Run each initializer, collecting errors (don't fail on first error)
	var initErrors []InitializerError
	for _, init := range sortedInitializers {
		// Select appropriate filesystem based on IsGlobal()
		var fs afero.Fs
		if init.IsGlobal() {
			fs = e.globalFs
		} else {
			fs = e.projectFs
		}

		// Run the initializer
		if err := init.Init(ctx, fs, cfg, e.tm); err != nil {
			initErrors = append(initErrors, InitializerError{
				Path:  init.Path(),
				Error: err,
			})
		} else {
			// Track the file path for this initializer
			path := init.Path()
			if path != "" {
				if init.IsSetup(fs, cfg) {
					result.UpdatedFiles = append(result.UpdatedFiles, path)
				} else {
					result.CreatedFiles = append(result.CreatedFiles, path)
				}
			}
		}
	}

	// Task 6.10: Report which initializers failed
	if len(initErrors) > 0 {
		for _, initErr := range initErrors {
			result.Errors = append(result.Errors,
				fmt.Sprintf("initializer %q failed: %v", initErr.Path, initErr.Error))
		}
	}

	return nil
}

// InitializerError tracks failures for individual initializers
type InitializerError struct {
	Path  string
	Error error
}

// collectInitializers gathers initializers from all selected providers using Registry.
func (e *InitExecutor) collectInitializers(
	ctx context.Context,
	providerIDs []string,
	result *ExecutionResult,
) []providers.Initializer {
	var all []providers.Initializer

	for _, providerID := range providerIDs {
		// Task 6.3: Use new registry API (Get with Registration-based retrieval)
		reg, found := providers.Get(providerID)
		if !found {
			result.Errors = append(result.Errors,
				fmt.Sprintf("provider %s not found in registry", providerID))

			continue
		}

		if reg.Provider == nil {
			result.Errors = append(result.Errors,
				fmt.Sprintf("provider %s has nil Provider implementation", providerID))

			continue
		}

		// Get initializers from the provider
		initializers := reg.Provider.Initializers(ctx)
		all = append(all, initializers...)
	}

	return all
}

// dedupeInitializers removes duplicate initializers based on Path().
// When multiple providers return initializers with the same Path(), only the first one is kept.
func (e *InitExecutor) dedupeInitializers(all []providers.Initializer) []providers.Initializer {
	seen := make(map[string]bool)
	var result []providers.Initializer

	for _, init := range all {
		if init == nil {
			continue
		}
		key := init.Path()
		if key == "" {
			// Include initializers without a path (shouldn't happen normally)
			result = append(result, init)

			continue
		}
		if !seen[key] {
			seen[key] = true
			result = append(result, init)
		}
	}

	return result
}

// sortInitializers sorts initializers by type to ensure correct execution order.
// Order: Directory (1) -> ConfigFile (2) -> SlashCommands (3) -> Others (99)
func (e *InitExecutor) sortInitializers(all []providers.Initializer) []providers.Initializer {
	// Use stable sort to maintain relative ordering within each type
	sort.SliceStable(all, func(i, j int) bool {
		return initializerPriority(all[i]) < initializerPriority(all[j])
	})

	return all
}

// initializerPriority returns the execution priority for an initializer.
// Lower numbers execute first.
func initializerPriority(init providers.Initializer) int {
	if init == nil {
		return 99
	}

	// Type switch to determine priority based on initializer type
	switch init.(type) {
	case *providers.DirectoryInitializerBuiltin:
		return 1
	case *providers.ConfigFileInitializerBuiltin:
		return 2
	case *providers.SlashCommandsInitializerBuiltin:
		return 3
	default:
		// Check for initializers from the initializers subpackage
		// We use the Path() method to infer type when concrete type matching fails
		path := init.Path()
		if path == "" {
			return 99
		}
		// Heuristic: directories often don't have extensions
		// Config files often have .md extension
		// This is a fallback for unknown types
		return 50
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
