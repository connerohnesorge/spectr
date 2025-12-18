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

	"github.com/connerohnesorge/spectr/internal/initialize/git"
	"github.com/connerohnesorge/spectr/internal/initialize/providers"
	"github.com/spf13/afero"
)

// InitExecutor handles the actual initialization process
type InitExecutor struct {
	projectPath string
	tm          *TemplateManager
	// Dual filesystem for project and global paths (Task 6.1)
	projectFs afero.Fs
	globalFs  afero.Fs
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

	// Create dual filesystem (Task 6.1)
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

	return &InitExecutor{
		projectPath: projectPath,
		tm:          tm,
		projectFs:   projectFs,
		globalFs:    globalFs,
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

	// Task 6.2: Check git repo at start - fail fast
	if !git.IsGitRepo(e.projectPath) {
		return nil, git.ErrNotGitRepo
	}

	// Task 6.8: Take snapshot before initialization for change detection
	detector := git.NewChangeDetector(e.projectPath)
	snapshot, err := detector.Snapshot()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to take git snapshot: %w",
			err,
		)
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

	// 5. Configure selected providers using new architecture (Task 6.7)
	ctx := context.Background()
	err = e.configureProviders(
		ctx,
		selectedProviderIDs,
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

	// Task 6.9: Update result with git-detected changed files
	changedFiles, err := detector.ChangedFiles(snapshot)
	if err != nil {
		result.Errors = append(
			result.Errors,
			fmt.Sprintf(
				"failed to detect changed files: %v",
				err,
			),
		)
	} else {
		// Replace declared paths with git-detected changes
		// This gives us the actual files that were modified
		result.CreatedFiles = changedFiles
		result.UpdatedFiles = nil // Git status doesn't distinguish, so clear this
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
// Task 6.7: Each provider returns initializers that are collected, deduplicated, sorted, and executed.
func (e *InitExecutor) configureProviders(
	ctx context.Context,
	selectedProviderIDs []string,
	result *ExecutionResult,
) error {
	if len(selectedProviderIDs) == 0 {
		return nil // No providers to configure
	}

	// Task 6.4: Collect initializers from all selected providers
	allInitializers := e.collectInitializers(ctx, selectedProviderIDs, result)

	// Task 6.5: Deduplicate by Path()
	dedupedInitializers := dedupeInitializers(allInitializers)

	// Task 6.6: Sort by type (guaranteed order)
	sortedInitializers := sortInitializers(dedupedInitializers)

	// Create config for initializers
	cfg := &providers.Config{SpectrDir: "spectr"}

	// Task 6.10: Execute each initializer, handling partial failures
	var initErrors []string
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
			// Task 6.10: Report failure but continue with rest
			errMsg := fmt.Sprintf(
				"failed to initialize %s: %v",
				init.Path(),
				err,
			)
			initErrors = append(initErrors, errMsg)
			continue
		}
	}

	// Add all initialization errors to result
	result.Errors = append(result.Errors, initErrors...)

	return nil
}

// collectInitializers collects initializers from all selected providers.
// Task 6.3 & 6.4: Uses new registry API (Registration-based retrieval)
func (e *InitExecutor) collectInitializers(
	ctx context.Context,
	providerIDs []string,
	result *ExecutionResult,
) []providers.Initializer {
	var all []providers.Initializer

	for _, id := range providerIDs {
		// Task 6.3: Use GetV2 instead of Get
		reg := providers.GetV2(id)
		if reg == nil {
			result.Errors = append(
				result.Errors,
				fmt.Sprintf(
					"provider %s not found",
					id,
				),
			)
			continue
		}

		// Get initializers from provider
		inits := reg.Provider.Initializers(ctx)
		all = append(all, inits...)
	}

	return all
}

// dedupeInitializers removes duplicate initializers based on their Path() and type.
// Task 6.5: Same path + same type = run once
// Different initializer types with the same path are NOT duplicates (e.g., DirectoryInitializer
// and SlashCommandsInitializer can both target the same directory path).
func dedupeInitializers(all []providers.Initializer) []providers.Initializer {
	seen := make(map[string]bool)
	var result []providers.Initializer

	for _, init := range all {
		// Create key from both type and path to allow different initializer types
		// to operate on the same path
		key := fmt.Sprintf("%T:%s", init, init.Path())
		if !seen[key] {
			seen[key] = true
			result = append(result, init)
		}
	}

	return result
}

// sortInitializers sorts initializers by type to ensure correct execution order.
// Task 6.6: Guaranteed order - directories first, then config files, then slash commands
func sortInitializers(all []providers.Initializer) []providers.Initializer {
	// Create a copy to avoid modifying the input
	sorted := make([]providers.Initializer, len(all))
	copy(sorted, all)

	sort.SliceStable(sorted, func(i, j int) bool {
		return initializerPriority(sorted[i]) < initializerPriority(sorted[j])
	})

	return sorted
}

// initializerPriority returns the execution priority for an initializer based on its type.
// Lower values execute first.
// Uses type name matching since initializers can come from either providers or initializers package.
func initializerPriority(init providers.Initializer) int {
	typeName := fmt.Sprintf("%T", init)

	// Directory initializers run first
	if strings.Contains(typeName, "directoryInitializer") || strings.Contains(typeName, "DirectoryInitializer") {
		return 1
	}
	// Config file initializers run second
	if strings.Contains(typeName, "configFileInitializer") || strings.Contains(typeName, "ConfigFileInitializer") {
		return 2
	}
	// Slash command initializers run last
	if strings.Contains(typeName, "slashCommandsInitializer") || strings.Contains(typeName, "SlashCommandsInitializer") {
		return 3
	}
	// Unknown types run at the end
	return 99
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
--------------------------------------------------------------------

Next steps:

1. Populate your project context by telling your AI assistant:

   "` + PopulateContextPrompt + `"

2. Create your first change proposal by saying:

   "Help me create a change proposal for [YOUR FEATURE HERE]. Walk me through
   the process and ask questions to understand the requirements."

3. Learn the Spectr workflow:

   "Review spectr/AGENTS.md and explain how Spectr's change workflow works."

--------------------------------------------------------------------
`
}
