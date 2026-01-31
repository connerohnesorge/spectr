package providers

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/spf13/afero"
)

// SkillFileInitializer creates a single SKILL.md file from a template.
// Unlike AgentSkillsInitializer which copies entire embedded skill directories,
// this initializer renders a template to generate a single SKILL.md file.
// Used for simple skills that only need a SKILL.md without additional resources.
type SkillFileInitializer struct {
	targetPath string             // Full path to SKILL.md file (e.g., ".agents/skills/spectr-proposal/SKILL.md")
	template   domain.TemplateRef // Template to render for SKILL.md content
}

// NewSkillFileInitializer creates an initializer for a single SKILL.md file.
//
// Parameters:
//   - targetPath: Full path to the SKILL.md file (e.g., ".agents/skills/spectr-proposal/SKILL.md")
//   - template: TemplateRef to render for the SKILL.md content
//
// Example:
//
//	init := NewSkillFileInitializer(
//		".agents/skills/spectr-proposal/SKILL.md",
//		tm.SkillProposal(),
//	)
func NewSkillFileInitializer(
	targetPath string,
	template domain.TemplateRef,
) *SkillFileInitializer {
	return &SkillFileInitializer{
		targetPath: targetPath,
		template:   template,
	}
}

// Init creates the SKILL.md file in the project filesystem.
//
// Behavior:
//   - Creates parent directory if it doesn't exist
//   - Renders the template using TemplateContext from config
//   - Overwrites existing file (idempotent operation)
//   - Returns InitResult with CreatedFiles or UpdatedFiles
//
//nolint:revive // Init signature is defined by Initializer interface
func (s *SkillFileInitializer) Init(
	ctx context.Context,
	projectFs, homeFs afero.Fs,
	cfg *Config,
	tm TemplateManager,
) (InitResult, error) {
	// Derive template context from config
	tmplCtx := &domain.TemplateContext{
		BaseDir:     cfg.SpectrDir,
		SpecsDir:    cfg.SpecsDir(),
		ChangesDir:  cfg.ChangesDir(),
		ProjectFile: cfg.ProjectFile(),
		AgentsFile:  cfg.AgentsFile(),
	}

	// Render template content
	content, err := s.template.Render(tmplCtx)
	if err != nil {
		return InitResult{}, fmt.Errorf(
			"failed to render template for %s: %w",
			s.targetPath,
			err,
		)
	}

	// Create parent directory if needed
	parentDir := filepath.Dir(s.targetPath)
	if err := projectFs.MkdirAll(parentDir, 0o755); err != nil {
		return InitResult{}, fmt.Errorf(
			"failed to create directory %s: %w",
			parentDir,
			err,
		)
	}

	// Check if file already exists
	exists, err := afero.Exists(projectFs, s.targetPath)
	if err != nil {
		return InitResult{}, fmt.Errorf(
			"failed to check file %s: %w",
			s.targetPath,
			err,
		)
	}

	// Write file (always overwrite for idempotent behavior)
	if err := afero.WriteFile(projectFs, s.targetPath, []byte(content), fileMode); err != nil {
		return InitResult{}, fmt.Errorf(
			"failed to write file %s: %w",
			s.targetPath,
			err,
		)
	}

	if exists {
		return InitResult{
			CreatedFiles: nil,
			UpdatedFiles: []string{s.targetPath},
		}, nil
	}

	return InitResult{
		CreatedFiles: []string{s.targetPath},
		UpdatedFiles: nil,
	}, nil
}

// IsSetup checks if the SKILL.md file exists in the project filesystem.
func (s *SkillFileInitializer) IsSetup(
	projectFs, _ afero.Fs,
	_ *Config,
) bool {
	exists, err := afero.Exists(projectFs, s.targetPath)

	return err == nil && exists
}

// dedupeKey returns a unique key for deduplication.
// Uses type name + normalized target path to prevent duplicate skill file creation.
//
// The key format is: SkillFileInitializer:<normalized-target-path>
//
// Example: SkillFileInitializer:.agents/skills/spectr-proposal/SKILL.md
//
//nolint:unused // Called through deduplicator interface in executor
func (s *SkillFileInitializer) dedupeKey() string {
	normalizedPath := filepath.Clean(s.targetPath)

	return fmt.Sprintf(
		"SkillFileInitializer:%s",
		normalizedPath,
	)
}
