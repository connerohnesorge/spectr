package providers

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/spf13/afero"
)

// File extension constants.
const (
	extMD   = ".md"
	extTOML = ".toml"
	// fileMode is the permission for created slash command files.
	fileMode = 0644
)

// SlashCommandsInitializer creates Markdown slash command files in the project filesystem. //nolint:lll
// Uses early binding - templates are resolved at construction time.
// Always overwrites existing files (idempotent behavior).
type SlashCommandsInitializer struct {
	dir      string                                     // Directory for slash commands (e.g., ".claude/commands/spectr") //nolint:lll
	commands map[domain.SlashCommand]domain.TemplateRef // Map of command to template //nolint:lll
}

// NewSlashCommandsInitializer creates an initializer for Markdown slash commands in project filesystem. //nolint:lll
//
//	Example: NewSlashCommandsInitializer(".claude/commands/spectr", map[domain.SlashCommand]domain.TemplateRef{ //nolint:lll
//	    domain.SlashProposal: tm.SlashCommand(domain.SlashProposal),
//	    domain.SlashApply: tm.SlashCommand(domain.SlashApply),
//	})
func NewSlashCommandsInitializer(
	dir string,
	commands map[domain.SlashCommand]domain.TemplateRef,
) Initializer {
	return &SlashCommandsInitializer{
		dir:      dir,
		commands: commands,
	}
}

// Init creates slash command files in the project filesystem.

// Always overwrites existing files (idempotent).
//
//nolint:revive // Init signature is defined by Initializer interface
func (s *SlashCommandsInitializer) Init(
	_ context.Context,
	projectFs, _ afero.Fs,
	cfg *Config,
	_ TemplateManager,
) (InitResult, error) {
	return createSlashCommands(projectFs, s.dir, s.commands, cfg, extMD)
}

// IsSetup checks if all slash command files exist in the project filesystem.
func (s *SlashCommandsInitializer) IsSetup(projectFs, _ afero.Fs, _ *Config) bool {
	return checkSlashCommandsExist(projectFs, s.dir, s.commands, extMD)
}

// dedupeKey returns a unique key for deduplication.
func (s *SlashCommandsInitializer) dedupeKey() string {
	return fmt.Sprintf("SlashCommandsInitializer:%s", filepath.Clean(s.dir))
}

// HomeSlashCommandsInitializer creates Markdown slash command files in the home filesystem.
// Uses early binding - templates are resolved at construction time.
// Always overwrites existing files (idempotent behavior).
type HomeSlashCommandsInitializer struct {
	dir      string                                     // Directory for slash commands (relative to home) //nolint:lll
	commands map[domain.SlashCommand]domain.TemplateRef // Map of command to template //nolint:lll
}

// NewHomeSlashCommandsInitializer creates an initializer for Markdown slash commands in home filesystem. //nolint:lll
// Example: NewHomeSlashCommandsInitializer(".config/mytool/commands", map[domain.SlashCommand]domain.TemplateRef{...}) //nolint:lll
func NewHomeSlashCommandsInitializer( //nolint:lll
	dir string,
	commands map[domain.SlashCommand]domain.TemplateRef,
) Initializer {
	return &HomeSlashCommandsInitializer{
		dir:      dir,
		commands: commands,
	}
}

// Init creates slash command files in the home filesystem.
// Always overwrites existing files (idempotent).
//
//nolint:revive // Init signature is defined by Initializer interface
func (h *HomeSlashCommandsInitializer) Init(
	_ context.Context,
	_, homeFs afero.Fs,
	cfg *Config,
	_ TemplateManager,
) (InitResult, error) {
	return createSlashCommands(homeFs, h.dir, h.commands, cfg, extMD)
}

// IsSetup checks if all slash command files exist in the home filesystem.
func (h *HomeSlashCommandsInitializer) IsSetup(_, homeFs afero.Fs, _ *Config) bool {
	return checkSlashCommandsExist(homeFs, h.dir, h.commands, extMD)
}

// dedupeKey returns a unique key for deduplication.
func (h *HomeSlashCommandsInitializer) dedupeKey() string {
	return fmt.Sprintf("HomeSlashCommandsInitializer:%s", filepath.Clean(h.dir))
}

// PrefixedSlashCommandsInitializer creates Markdown slash command files with custom prefix in project filesystem.
// Used by providers like Antigravity that need custom file naming.
// Output: {dir}/{prefix}{command}.md (e.g., .agent/workflows/spectr-proposal.md)
type PrefixedSlashCommandsInitializer struct {
	dir      string                                     // Directory for slash commands
	prefix   string                                     // Prefix for file names (e.g., "spectr-")
	commands map[domain.SlashCommand]domain.TemplateRef // Map of command to template //nolint:lll
}

// NewPrefixedSlashCommandsInitializer creates an initializer for prefixed Markdown slash commands in project filesystem.
// Example: NewPrefixedSlashCommandsInitializer(".agent/workflows", "spectr-", map[domain.SlashCommand]domain.TemplateRef{...})
func NewPrefixedSlashCommandsInitializer( //nolint:lll
	dir, prefix string,
	commands map[domain.SlashCommand]domain.TemplateRef,
) Initializer {
	return &PrefixedSlashCommandsInitializer{ //nolint:lll
		dir:      dir,
		prefix:   prefix,
		commands: commands,
	}
}

// Init creates prefixed slash command files in the project filesystem.
// Always overwrites existing files (idempotent).
//
//nolint:revive // Init signature is defined by Initializer interface
func (p *PrefixedSlashCommandsInitializer) Init(
	_ context.Context,
	projectFs, _ afero.Fs,
	cfg *Config,
	_ TemplateManager,
) (InitResult, error) {
	return createPrefixedSlashCommands(projectFs, p.dir, p.prefix, p.commands, cfg, extMD)
}

// IsSetup checks if all prefixed slash command files exist in the project filesystem.
func (p *PrefixedSlashCommandsInitializer) IsSetup(projectFs, _ afero.Fs, _ *Config) bool {
	return checkPrefixedSlashCommandsExist(projectFs, p.dir, p.prefix, p.commands, extMD)
}

// dedupeKey returns a unique key for deduplication.
func (p *PrefixedSlashCommandsInitializer) dedupeKey() string {
	return fmt.Sprintf("PrefixedSlashCommandsInitializer:%s:%s", filepath.Clean(p.dir), p.prefix)
}

// HomePrefixedSlashCommandsInitializer creates Markdown slash command files with custom prefix in home filesystem.
// Used by providers like Codex that need custom file naming in home directory.
// Output: {dir}/{prefix}{command}.md (e.g., ~/.codex/prompts/spectr-proposal.md)
type HomePrefixedSlashCommandsInitializer struct {
	dir      string                                     // Directory for slash commands (relative to home)
	prefix   string                                     // Prefix for file names (e.g., "spectr-") //nolint:lll
	commands map[domain.SlashCommand]domain.TemplateRef // Map of command to template
}

// NewHomePrefixedSlashCommandsInitializer creates an initializer for prefixed Markdown slash commands in home filesystem.
// Example: NewHomePrefixedSlashCommandsInitializer(".codex/prompts", "spectr-", map[domain.SlashCommand]domain.TemplateRef{...}) //nolint:lll
func NewHomePrefixedSlashCommandsInitializer(
	dir, prefix string,
	commands map[domain.SlashCommand]domain.TemplateRef,
) Initializer { //nolint:lll
	return &HomePrefixedSlashCommandsInitializer{
		dir:      dir,
		prefix:   prefix,
		commands: commands,
	}
}

// Init creates prefixed slash command files in the home filesystem.
// Always overwrites existing files (idempotent).
//
//nolint:revive // Init signature is defined by Initializer interface
func (h *HomePrefixedSlashCommandsInitializer) Init(
	_ context.Context,
	_, homeFs afero.Fs,
	cfg *Config,
	_ TemplateManager,
) (InitResult, error) {
	return createPrefixedSlashCommands(homeFs, h.dir, h.prefix, h.commands, cfg, extMD)
}

// IsSetup checks if all prefixed slash command files exist in the home filesystem.
func (h *HomePrefixedSlashCommandsInitializer) IsSetup(
	_, homeFs afero.Fs,
	_ *Config,
) bool {
	return checkPrefixedSlashCommandsExist(homeFs, h.dir, h.prefix, h.commands, extMD)
}

// dedupeKey returns a unique key for deduplication.
func (h *HomePrefixedSlashCommandsInitializer) dedupeKey() string {
	return fmt.Sprintf(
		"HomePrefixedSlashCommandsInitializer:%s:%s",
		filepath.Clean(h.dir),
		h.prefix,
	)
} //nolint:lll

// TOMLSlashCommandsInitializer creates TOML slash command files in the project filesystem.
// Used by Gemini provider for TOML format.
// Uses early binding - templates are resolved at construction time.
// Always overwrites existing files (idempotent behavior). //nolint:lll
type TOMLSlashCommandsInitializer struct {
	dir      string                                     // Directory for slash commands (e.g., ".gemini/commands/spectr")
	commands map[domain.SlashCommand]domain.TemplateRef // Map of command to template
}

// NewTOMLSlashCommandsInitializer creates an initializer for TOML slash commands in project filesystem.
//
//	Example: NewTOMLSlashCommandsInitializer(".gemini/commands/spectr", map[domain.SlashCommand]domain.TemplateRef{
//	    domain.SlashProposal: tm.TOMLSlashCommand(domain.SlashProposal),
//	    domain.SlashApply: tm.TOMLSlashCommand(domain.SlashApply),
//	})
func NewTOMLSlashCommandsInitializer(
	dir string,
	commands map[domain.SlashCommand]domain.TemplateRef,
) Initializer {
	return &TOMLSlashCommandsInitializer{
		dir:      dir,
		commands: commands,

		//nolint:revive // Init signature is defined by Initializer interface
	}
}

// Init creates TOML slash command files in the project filesystem.
// Always overwrites existing files (idempotent).
//
//nolint:revive // Init signature is defined by Initializer interface
func (t *TOMLSlashCommandsInitializer) Init(
	_ context.Context,
	projectFs, _ afero.Fs,
	cfg *Config,
	_ TemplateManager,
) (InitResult, error) {
	return createSlashCommands(projectFs, t.dir, t.commands, cfg, extTOML)
}

// IsSetup checks if all TOML slash command files exist in the project filesystem.
func (t *TOMLSlashCommandsInitializer) IsSetup(projectFs, _ afero.Fs, _ *Config) bool {
	return checkSlashCommandsExist(projectFs, t.dir, t.commands, extTOML)
}

// dedupeKey returns a unique key for deduplication.
func (t *TOMLSlashCommandsInitializer) dedupeKey() string {
	return fmt.Sprintf("TOMLSlashCommandsInitializer:%s", filepath.Clean(t.dir)) //nolint:lll
}

// createSlashCommands is a helper that creates slash command files with the given extension.
// Used by SlashCommandsInitializer, HomeSlashCommandsInitializer, and TOMLSlashCommandsInitializer.
//
//nolint:revive // Helper function needs these params for clarity
func createSlashCommands(
	fs afero.Fs,
	dir string,
	commands map[domain.SlashCommand]domain.TemplateRef,
	cfg *Config,
	ext string,
) (InitResult, error) {
	// Derive template context from config
	tmplCtx := domain.TemplateContext{
		BaseDir:     cfg.SpectrDir,
		SpecsDir:    cfg.SpecsDir(),
		ChangesDir:  cfg.ChangesDir(),
		ProjectFile: cfg.ProjectFile(),
		AgentsFile:  cfg.AgentsFile(),
	}

	// Create parent directory if it doesn't exist
	if err := fs.MkdirAll(dir, 0755); err != nil {
		return InitResult{}, fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	var created []string
	var updated []string

	for cmd, template := range commands {
		// Render template
		content, err := template.Render(tmplCtx)
		if err != nil {
			return InitResult{}, fmt.Errorf(
				"failed to render template for %s: %w", //nolint:lll
				cmd.String(),
				err,
			)
		}
		//nolint:lll
		// Create file path: {dir}/{command}{ext}
		filePath := filepath.Join(dir, cmd.String()+ext)

		// Check if file already exists //nolint:lll
		exists, err := afero.Exists(fs, filePath)
		if err != nil {
			return InitResult{}, fmt.Errorf("failed to check file %s: %w", filePath, err)
		}

		// Write file (always overwrite for idempotent behavior)
		if err := afero.WriteFile(fs, filePath, []byte(content), fileMode); err != nil {
			return InitResult{}, fmt.Errorf("failed to write file %s: %w", filePath, err)
		}

		if exists {
			updated = append(updated, filePath)
		} else {
			created = append(created, filePath)
		}
	}

	return InitResult{
		CreatedFiles: created,
		UpdatedFiles: updated,
	}, nil
}

// createPrefixedSlashCommands is a helper that creates prefixed slash command files.
// Used by PrefixedSlashCommandsInitializer and HomePrefixedSlashCommandsInitializer.
//
//nolint:revive // Helper function needs multiple params for clarity
func createPrefixedSlashCommands(
	fs afero.Fs,
	dir, prefix string,
	commands map[domain.SlashCommand]domain.TemplateRef,
	cfg *Config,
	ext string,
) (InitResult, error) {
	// Derive template context from config
	tmplCtx := domain.TemplateContext{
		BaseDir:     cfg.SpectrDir,
		SpecsDir:    cfg.SpecsDir(),
		ChangesDir:  cfg.ChangesDir(),
		ProjectFile: cfg.ProjectFile(),
		AgentsFile:  cfg.AgentsFile(),
	}

	// Create parent directory if it doesn't exist
	if err := fs.MkdirAll(dir, 0755); err != nil {
		return InitResult{}, fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	var created []string
	var updated []string

	for cmd, template := range commands {
		// Render template
		content, err := template.Render(tmplCtx)
		if err != nil {
			return InitResult{}, fmt.Errorf(
				"failed to render template for %s: %w",
				cmd.String(),
				err,
			)
		}

		// Create file path: {dir}/{prefix}{command}{ext}
		filePath := filepath.Join(dir, prefix+cmd.String()+ext)

		// Check if file already exists
		exists, err := afero.Exists(fs, filePath)
		if err != nil {
			return InitResult{}, fmt.Errorf("failed to check file %s: %w", filePath, err)
		}

		// Write file (always overwrite for idempotent behavior)
		if err := afero.WriteFile(fs, filePath, []byte(content), fileMode); err != nil {
			return InitResult{}, fmt.Errorf("failed to write file %s: %w", filePath, err)
		}

		if exists {
			updated = append(updated, filePath)
		} else {
			created = append(created, filePath)
		}
	}

	return InitResult{
		CreatedFiles: created,
		UpdatedFiles: updated,
	}, nil
}

// checkSlashCommandsExist checks if all slash command files exist.
func checkSlashCommandsExist(
	fs afero.Fs,
	dir string,
	commands map[domain.SlashCommand]domain.TemplateRef,
	ext string,
) bool {
	for cmd := range commands {
		filePath := filepath.Join(dir, cmd.String()+ext)
		exists, err := afero.Exists(fs, filePath)
		if err != nil || !exists {
			return false
		}
	}

	return true
}

// checkPrefixedSlashCommandsExist checks if all prefixed slash command files exist.
//
//nolint:revive // Helper function needs these params for clarity
func checkPrefixedSlashCommandsExist(
	fs afero.Fs,
	dir, prefix string,
	commands map[domain.SlashCommand]domain.TemplateRef,
	ext string,
) bool {
	for cmd := range commands {
		filePath := filepath.Join(dir, prefix+cmd.String()+ext)
		exists, err := afero.Exists(fs, filePath)
		if err != nil || !exists {
			return false
		}
	}

	return true
}
