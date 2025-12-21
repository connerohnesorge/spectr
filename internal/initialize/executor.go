// Package initialize provides utilities for initializing Spectr
// in a project directory.
//
//nolint:revive // file-length-limit - logically cohesive, no benefit to split
package initialize

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/spf13/afero"
	"github.com/connerohnesorge/spectr/internal/initialize/git"
	"github.com/connerohnesorge/spectr/internal/initialize/providers"
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
	"github.com/connerohnesorge/spectr/internal/initialize/templates"
	"github.com/connerohnesorge/spectr/internal/initialize/types"
)

// InitExecutor handles the actual initialization process
type InitExecutor struct {
	projectPath string
	tm          *templates.TemplateManager
	projectFs   afero.Fs
	globalFs    afero.Fs
	detector    *git.ChangeDetector
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
		projectFs:   afero.NewBasePathFs(afero.NewOsFs(), projectPath),
		globalFs:    afero.NewOsFs(), // Ideally this would be limited or specialized
		detector:    git.NewChangeDetector(projectPath),
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

	// 0. Require git repository
	if !git.IsGitRepo(e.projectPath) {
		return result, fmt.Errorf("spectr init requires a git repository. Run 'git init' first")
	}

	// 0.1 Capture initial state
	snapshotBefore, err := e.detector.Snapshot()
	if err != nil {
		return result, fmt.Errorf("failed to capture git snapshot: %w", err)
	}

	// 1. Check if Spectr is already initialized
	if IsSpectrInitialized(e.projectPath) {
		result.Errors = append(
			result.Errors,
			"Spectr already initialized in this project",
		)
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

	// 5. Configure selected providers
	err = e.configureProviders(
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

	// 7. Calculate changed files using git
	changedFiles, err := e.detector.ChangedFiles(snapshotBefore)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("failed to detect changed files: %v", err))
	} else {
        // Merge git detected files into CreatedFiles for display, avoiding duplicates
        seen := make(map[string]bool)
        for _, f := range result.CreatedFiles { seen[f] = true }
        for _, f := range result.UpdatedFiles { seen[f] = true }
        
        for _, f := range changedFiles {
            if !seen[f] {
                result.CreatedFiles = append(result.CreatedFiles, f)
                seen[f] = true
            }
        }
	}

	return result, nil
}

// createDirectoryStructure creates the spectr/ directory
// and subdirectories
func (e *InitExecutor) createDirectoryStructure(
	spectrDir string,
	result *ExecutionResult,
) error {
	dirs := []string{
		"spectr",
		filepath.Join("spectr", "specs"),
		filepath.Join("spectr", "changes"),
	}

	for _, dir := range dirs {
		exists, _ := afero.Exists(e.projectFs, dir)
		if !exists {
			if err := e.projectFs.MkdirAll(dir, 0755); err != nil {
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
		"spectr",
		"project.md",
	)

	// Check if it already exists
	exists, _ := afero.Exists(e.projectFs, projectFile)
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
		types.ProjectContext{
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
	if err := afero.WriteFile(e.projectFs, projectFile, []byte(content), 0644); err != nil {
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
		"spectr",
		"AGENTS.md",
	)

	// Check if it already exists
	exists, _ := afero.Exists(e.projectFs, agentsFile)
	if exists {
		result.Errors = append(
			result.Errors,
			"AGENTS.md already exists, skipping",
		)

		return nil
	}

	// Render template
	content, err := e.tm.RenderAgents(
		types.DefaultTemplateContext(),
	)
	if err != nil {
		return fmt.Errorf(
			"failed to render agents template: %w",
			err,
		)
	}

	// Write file
	if err := afero.WriteFile(e.projectFs, agentsFile, []byte(content), 0644); err != nil {
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

// configureProviders configures the selected providers using the new composable architecture.
func (e *InitExecutor) configureProviders(
	selectedProviderIDs []string,
	result *ExecutionResult,
) error {
	if len(selectedProviderIDs) == 0 {
		return nil
	}

	var allInitializers []types.Initializer
	seenPaths := make(map[string]bool)

	// 1. Collect initializers from all selected providers
	for _, id := range selectedProviderIDs {
		reg, exists := providers.Get(id)
		if !exists {
			result.Errors = append(result.Errors, fmt.Sprintf("provider %s not found", id))
			continue
		}

		inits := reg.Provider.Initializers()
		for _, ini := range inits {
			path := ini.Path()
			if path != "" {
				if seenPaths[path] {
					continue // Deduplicate by path
				}
				seenPaths[path] = true
			}
			allInitializers = append(allInitializers, ini)
		}
	}

	// 2. Sort initializers by type (Directories first)
	sort.SliceStable(allInitializers, func(i, j int) bool {
		_, iIsDir := allInitializers[i].(*initializers.DirectoryInitializer)
		_, jIsDir := allInitializers[j].(*initializers.DirectoryInitializer)
		if iIsDir && !jIsDir {
			return true
		}
		return false
	})

	// 3. Run initializers
	cfg := &types.Config{SpectrDir: "spectr"}
	ctx := context.Background()
	for _, ini := range allInitializers {
		setup, err := ini.IsSetup(e.projectFs, e.globalFs, cfg)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("failed to check status of %s: %v", ini.Path(), err))
			continue
		}

		if err := ini.Init(ctx, e.projectFs, e.globalFs, cfg, e.tm); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("failed to initialize %s: %v", ini.Path(), err))
			if !setup {
				// If it wasn't setup and failed, we might want to continue or stop.
				// Proposal says: "No rollback, report failures. Keep simple; users can re-run init"
			}
		}
	}

	return nil
}

// createCIWorkflow creates the .github/workflows/spectr-ci.yml file
func (e *InitExecutor) createCIWorkflow(
	result *ExecutionResult,
) error {
	workflowDir := filepath.Join(
		e.projectPath,
		".github",
		"workflows",
	)
	if err := e.projectFs.MkdirAll(workflowDir, 0755); err != nil {
		return fmt.Errorf(
			"failed to create workflows directory: %w",
			err,
		)
	}

	workflowFile := filepath.Join(
		workflowDir,
		"spectr-ci.yml",
	)

	// Render the CI workflow template
	content, err := e.tm.RenderCIWorkflow()
	if err != nil {
		return fmt.Errorf(
			"failed to render CI workflow template: %w",
			err,
		)
	}

	// Write the workflow file
	if err := afero.WriteFile(e.projectFs, workflowFile, []byte(content), 0644); err != nil {
		return fmt.Errorf(
			"failed to write CI workflow file: %w",
			err,
		)
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
