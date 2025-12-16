// Package initialize provides utilities for initializing Spectr
// in a project directory.
//
// This file implements the new executor using the redesigned provider
// architecture with afero-based filesystem operations and initializer
// deduplication.
package initialize

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"

	"github.com/connerohnesorge/spectr/internal/initialize/providers"
	initpkg "github.com/connerohnesorge/spectr/internal/initialize/providers/initializers" //nolint:revive
)

// Key prefix length constants for parsing initializer keys.
const (
	keyPrefixDir       = 4  // len("dir:")
	keyPrefixConfig    = 7  // len("config:")
	keyPrefixSlashCmds = 10 // len("slashcmds:")
)

// Directory permission constant.
const dirPerm = 0755

// InitExecutorNew handles initialization using the new provider architecture.
// It uses afero.Fs rooted at the project directory for cleaner path handling,
// and collects/deduplicates initializers from the new provider system.
type InitExecutorNew struct {
	// fs is the project-rooted filesystem (afero.BasePathFs).
	// All paths used with this filesystem are relative to the project root.
	fs afero.Fs

	// projectPath is the absolute path to the project directory.
	projectPath string

	// registry is the provider registry containing all registered providers.
	registry *providers.ProviderRegistry

	// cfg is the configuration passed to initializers.
	cfg *providers.Config

	// tm is the template renderer for generating content.
	tm providers.TemplateRenderer
}

// NewInitExecutorNew creates a new initialization executor using the new
// provider architecture. It creates an afero filesystem rooted at the project
// path, registers all providers, and prepares for initialization.
func NewInitExecutorNew(
	projectPath string,
	tm providers.TemplateRenderer,
) (*InitExecutorNew, error) {
	if projectPath == "" {
		return nil, errors.New("project path is required")
	}

	osFs := afero.NewOsFs()
	exists, err := afero.DirExists(osFs, projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to check project path: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf(
			"project path does not exist: %s",
			projectPath,
		)
	}

	fs := afero.NewBasePathFs(osFs, projectPath)
	registry := providers.CreateRegistry()
	factory := initpkg.NewFactory()

	err = providers.RegisterAllProvidersIncludingBase(registry, tm, factory)
	if err != nil {
		return nil, fmt.Errorf("failed to register providers: %w", err)
	}

	cfg := providers.NewConfig()

	return &InitExecutorNew{
		fs:          fs,
		projectPath: projectPath,
		registry:    registry,
		cfg:         cfg,
		tm:          tm,
	}, nil
}

// Execute runs the initialization process for the selected providers.
// It creates the spectr directory structure, project files, and configures
// the selected providers using the new initializer-based architecture.
func (e *InitExecutorNew) Execute(
	ctx context.Context,
	selectedProviderIDs []string,
	opts ExecuteOptions,
) (*ExecutionResult, error) {
	result := &ExecutionResult{
		CreatedFiles: make([]string, 0),
		UpdatedFiles: make([]string, 0),
		Errors:       make([]string, 0),
	}

	if e.isSpectrInitialized() {
		result.Errors = append(
			result.Errors,
			"Spectr already initialized in this project",
		)
	}

	if err := e.createDirectoryStructure(result); err != nil {
		return result, fmt.Errorf(
			"failed to create directory structure: %w",
			err,
		)
	}

	if err := e.createProjectMd(result); err != nil {
		errMsg := fmt.Sprintf("failed to create project.md: %v", err)
		result.Errors = append(result.Errors, errMsg)
	}

	if err := e.createAgentsMd(result); err != nil {
		errMsg := fmt.Sprintf("failed to create AGENTS.md: %v", err)
		result.Errors = append(result.Errors, errMsg)
	}

	allInits, err := e.collectInitializers(ctx, selectedProviderIDs, result)
	if err != nil {
		errMsg := fmt.Sprintf("failed to collect initializers: %v", err)
		result.Errors = append(result.Errors, errMsg)
	}

	dedupedInitializers := providers.DedupeInitializers(allInits)

	if err := e.runInitializers(ctx, dedupedInitializers, result); err != nil {
		errMsg := fmt.Sprintf("failed to run initializers: %v", err)
		result.Errors = append(result.Errors, errMsg)
	}

	if opts.CIWorkflowEnabled {
		if err := e.createCIWorkflow(result); err != nil {
			errMsg := fmt.Sprintf("failed to create CI workflow: %v", err)
			result.Errors = append(result.Errors, errMsg)
		}
	}

	return result, nil
}

// isSpectrInitialized checks if Spectr is already initialized.
func (e *InitExecutorNew) isSpectrInitialized() bool {
	projectFile := filepath.Join(e.cfg.SpectrDir, "project.md")
	exists, err := afero.Exists(e.fs, projectFile)

	return err == nil && exists
}

// createDirectoryStructure creates the spectr/ directory and subdirectories.
func (e *InitExecutorNew) createDirectoryStructure(
	result *ExecutionResult,
) error {
	dirs := []string{
		e.cfg.SpectrDir,
		filepath.Join(e.cfg.SpectrDir, "specs"),
		filepath.Join(e.cfg.SpectrDir, "changes"),
	}

	for _, dir := range dirs {
		exists, err := afero.DirExists(e.fs, dir)
		if err != nil {
			return fmt.Errorf(
				"failed to check directory %s: %w",
				dir,
				err,
			)
		}
		if exists {
			continue
		}
		if err := e.fs.MkdirAll(dir, dirPerm); err != nil {
			return fmt.Errorf(
				"failed to create directory %s: %w",
				dir,
				err,
			)
		}
		result.CreatedFiles = append(result.CreatedFiles, dir+"/")
	}

	return nil
}

// createProjectMd creates the project.md file.
func (e *InitExecutorNew) createProjectMd(result *ExecutionResult) error {
	projectFile := filepath.Join(e.cfg.SpectrDir, "project.md")

	exists, err := afero.Exists(e.fs, projectFile)
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

	projectName := filepath.Base(e.projectPath)
	tm, err := NewTemplateManager()
	if err != nil {
		return fmt.Errorf("failed to create template manager: %w", err)
	}

	content, err := tm.RenderProject(ProjectContext{
		ProjectName: projectName,
		Description: "Add your project description here",
		TechStack:   []string{"Add", "Your", "Technologies", "Here"},
		Conventions: "",
	})
	if err != nil {
		return fmt.Errorf("failed to render project template: %w", err)
	}

	err = afero.WriteFile(e.fs, projectFile, []byte(content), filePerm)
	if err != nil {
		return fmt.Errorf("failed to write project.md: %w", err)
	}

	projectPath := filepath.Join(e.cfg.SpectrDir, "project.md")
	result.CreatedFiles = append(result.CreatedFiles, projectPath)

	return nil
}

// createAgentsMd creates the AGENTS.md file.
func (e *InitExecutorNew) createAgentsMd(result *ExecutionResult) error {
	agentsFile := filepath.Join(e.cfg.SpectrDir, "AGENTS.md")

	exists, err := afero.Exists(e.fs, agentsFile)
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

	content, err := e.tm.RenderAgents(providers.DefaultTemplateContext())
	if err != nil {
		return fmt.Errorf("failed to render agents template: %w", err)
	}

	err = afero.WriteFile(e.fs, agentsFile, []byte(content), filePerm)
	if err != nil {
		return fmt.Errorf("failed to write AGENTS.md: %w", err)
	}

	agentsPath := filepath.Join(e.cfg.SpectrDir, "AGENTS.md")
	result.CreatedFiles = append(result.CreatedFiles, agentsPath)

	return nil
}

// collectInitializers collects initializers from all selected providers.
func (e *InitExecutorNew) collectInitializers(
	ctx context.Context,
	selectedProviderIDs []string,
	result *ExecutionResult,
) ([]providers.Initializer, error) {
	if len(selectedProviderIDs) == 0 {
		return nil, nil
	}

	var allInits []providers.Initializer
	for _, providerID := range selectedProviderIDs {
		reg := e.registry.Get(providerID)
		if reg == nil {
			errMsg := fmt.Sprintf("provider %s not found", providerID)
			result.Errors = append(result.Errors, errMsg)

			continue
		}
		inits := reg.Provider.Initializers(ctx)
		allInits = append(allInits, inits...)
	}

	return allInits, nil
}

// runInitializers runs all initializers with the new Init signature.
func (e *InitExecutorNew) runInitializers(
	ctx context.Context,
	inits []providers.Initializer,
	result *ExecutionResult,
) error {
	for _, init := range inits {
		tracking := FileCreated
		if init.IsSetup(e.fs, e.cfg) {
			tracking = FileUpdated
		}
		if err := init.Init(ctx, e.fs, e.cfg); err != nil {
			return fmt.Errorf("initializer failed: %w", err)
		}
		trackInitializerFiles(init, tracking, result)
	}

	return nil
}

// createCIWorkflow creates the .github/workflows/spectr-ci.yml file.
func (e *InitExecutorNew) createCIWorkflow(result *ExecutionResult) error {
	workflowDir := filepath.Join(".github", "workflows")
	if err := e.fs.MkdirAll(workflowDir, dirPerm); err != nil {
		return fmt.Errorf("failed to create workflows directory: %w", err)
	}

	workflowFile := filepath.Join(workflowDir, "spectr-ci.yml")
	wasConfigured, _ := afero.Exists(e.fs, workflowFile)

	tm, err := NewTemplateManager()
	if err != nil {
		return fmt.Errorf("failed to create template manager: %w", err)
	}

	content, err := tm.RenderCIWorkflow()
	if err != nil {
		return fmt.Errorf("failed to render CI workflow: %w", err)
	}

	err = afero.WriteFile(e.fs, workflowFile, []byte(content), filePerm)
	if err != nil {
		return fmt.Errorf("failed to write CI workflow file: %w", err)
	}

	ciPath := ".github/workflows/spectr-ci.yml"
	ciTracking := FileCreated
	if wasConfigured {
		ciTracking = FileUpdated
	}
	trackFile(result, ciPath, ciTracking)

	return nil
}

// GetRegistry returns the provider registry for external access.
func (e *InitExecutorNew) GetRegistry() *providers.ProviderRegistry {
	return e.registry
}
