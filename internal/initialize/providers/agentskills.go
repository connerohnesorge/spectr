package providers

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"

	"github.com/spf13/afero"
)

// AgentSkillsInitializer copies an embedded skill directory to a target path.
// This initializer handles copying entire skill directories from embedded
// templates to the project filesystem, preserving directory structure and
// file permissions (e.g., executable scripts).
//
// Skills follow the AgentSkills specification (https://agentskills.io) and
// must contain at minimum a SKILL.md file.
type AgentSkillsInitializer struct {
	skillName string          // name of the skill (matches embedded directory name)
	targetDir string          // target directory path (e.g., ".claude/skills/my-skill")
	tm        TemplateManager // template manager for accessing embedded skill files
}

// NewAgentSkillsInitializer creates an AgentSkillsInitializer for the given
// skill name and target directory.
//
// Parameters:
//   - skillName: Name of the embedded skill directory (e.g., "spectr-accept-wo-spectr-bin")
//   - targetDir: Target directory path relative to project root (e.g., ".claude/skills/spectr-accept-wo-spectr-bin")
//   - tm: TemplateManager providing access to embedded skill files via SkillFS
//
// Example:
//
//	init := NewAgentSkillsInitializer(
//		"spectr-accept-wo-spectr-bin",
//		".claude/skills/spectr-accept-wo-spectr-bin",
//		tm,
//	)
func NewAgentSkillsInitializer(
	skillName, targetDir string,
	tm TemplateManager,
) *AgentSkillsInitializer {
	return &AgentSkillsInitializer{
		skillName: skillName,
		targetDir: targetDir,
		tm:        tm,
	}
}

// Init recursively copies all files from the embedded skill directory to the
// target directory in the project filesystem.
//
// Behavior:
//   - Creates target directory if it doesn't exist
//   - Recursively copies all files and subdirectories
//   - Preserves directory structure (e.g., scripts/accept.sh)
//   - Preserves file permissions (executable scripts remain executable)
//   - Overwrites existing files (idempotent operation)
//   - Returns error if skill not found in embedded templates
//
// The initializer uses TemplateManager.SkillFS() to access the embedded
// skill filesystem rooted at the skill directory.
//
//nolint:revive // Init signature is defined by Initializer interface
func (a *AgentSkillsInitializer) Init(
	ctx context.Context,
	projectFs, homeFs afero.Fs,
	cfg *Config,
	tm TemplateManager,
) (InitResult, error) {
	// Get the embedded skill filesystem
	skillFS, err := tm.SkillFS(a.skillName)
	if err != nil {
		return InitResult{}, fmt.Errorf(
			"failed to get skill filesystem for %s: %w",
			a.skillName,
			err,
		)
	}

	var createdFiles []string
	var updatedFiles []string

	// Walk the skill filesystem and copy all files
	err = fs.WalkDir(skillFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Compute target path
		targetPath := filepath.Join(a.targetDir, path)

		// Create directory
		if d.IsDir() {
			exists, existsErr := afero.DirExists(projectFs, targetPath)
			if existsErr != nil {
				return fmt.Errorf(
					"failed to check directory %s: %w",
					targetPath,
					existsErr,
				)
			}

			if mkdirErr := projectFs.MkdirAll(targetPath, 0o755); mkdirErr != nil {
				return fmt.Errorf(
					"failed to create directory %s: %w",
					targetPath,
					mkdirErr,
				)
			}

			if !exists {
				createdFiles = append(createdFiles, targetPath)
			}

			return nil
		}

		// Check if file exists
		fileExists, existsErr := afero.Exists(projectFs, targetPath)
		if existsErr != nil {
			return fmt.Errorf(
				"failed to check file %s: %w",
				targetPath,
				existsErr,
			)
		}

		// Read source file
		sourceFile, openErr := skillFS.Open(path)
		if openErr != nil {
			return fmt.Errorf(
				"failed to open skill file %s: %w",
				path,
				openErr,
			)
		}
		defer func() {
			_ = sourceFile.Close()
		}()

		sourceData, readErr := io.ReadAll(sourceFile)
		if readErr != nil {
			return fmt.Errorf(
				"failed to read skill file %s: %w",
				path,
				readErr,
			)
		}

		// Get source file permissions to detect executables
		sourceInfo, statErr := d.Info()
		if statErr != nil {
			return fmt.Errorf(
				"failed to get file info for %s: %w",
				path,
				statErr,
			)
		}

		// Determine target file mode: use 0755 for executables, 0644 for regular files
		// Don't preserve readonly permissions from embed.FS
		targetMode := fs.FileMode(0o644)

		// Check if source has executable bit OR if it's a .sh file
		// (embed.FS doesn't preserve executable bits from git, so we check extension)
		isExecutable := sourceInfo.Mode()&0o111 != 0
		isShellScript := filepath.Ext(path) == ".sh"

		if isExecutable || isShellScript {
			targetMode = 0o755
		}

		// Write target file with normal write permissions
		writeErr := afero.WriteFile(
			projectFs,
			targetPath,
			sourceData,
			targetMode,
		)
		if writeErr != nil {
			return fmt.Errorf(
				"failed to write file %s: %w",
				targetPath,
				writeErr,
			)
		}

		// Track created/updated files
		if fileExists {
			updatedFiles = append(updatedFiles, targetPath)
		} else {
			createdFiles = append(createdFiles, targetPath)
		}

		return nil
	})

	if err != nil {
		return InitResult{}, fmt.Errorf(
			"failed to copy skill %s to %s: %w",
			a.skillName,
			a.targetDir,
			err,
		)
	}

	return InitResult{
		CreatedFiles: createdFiles,
		UpdatedFiles: updatedFiles,
	}, nil
}

// IsSetup returns true if the SKILL.md file exists in the target directory.
// This indicates that the skill has been initialized at least once.
//
// The check is simple - we only verify SKILL.md exists, not the complete
// skill structure. This matches the pattern used by other initializers
// (e.g., DirectoryInitializer checks directory existence).
func (a *AgentSkillsInitializer) IsSetup(
	projectFs, _ afero.Fs,
	_ *Config,
) bool { //nolint:lll // Function signature defined by Initializer interface
	skillMdPath := filepath.Join(a.targetDir, "SKILL.md")
	exists, err := afero.Exists(projectFs, skillMdPath)

	return err == nil && exists
}

// dedupeKey returns a unique key for deduplication.
// Uses type name + normalized target directory path to prevent duplicate
// skill initialization.
//
// The key format is: AgentSkillsInitializer:<normalized-target-dir>
//
// Example: AgentSkillsInitializer:.claude/skills/spectr-accept-wo-spectr-bin
//
//nolint:unused // Called through deduplicator interface in executor
func (a *AgentSkillsInitializer) dedupeKey() string {
	normalizedPath := filepath.Clean(a.targetDir)

	return fmt.Sprintf("AgentSkillsInitializer:%s", normalizedPath)
}
