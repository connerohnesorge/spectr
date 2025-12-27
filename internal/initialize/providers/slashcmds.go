//nolint:revive // file-length-limit: initializers need 4 types + 4 helpers
package providers

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/connerohnesorge/spectr/internal/templates"
	"github.com/spf13/afero"
)

const (
	mdExtension = ".md"
)

// SlashCommandsInitializer creates Markdown slash command files in project
// filesystem. Uses early binding with map[SlashCommand]TemplateRef.
type SlashCommandsInitializer struct {
	dir      string                                     // directory for slash commands
	commands map[domain.SlashCommand]domain.TemplateRef // command templates
}

// NewSlashCommandsInitializer creates a new SlashCommandsInitializer for
// project filesystem.
func NewSlashCommandsInitializer(
	dir string,
	commands map[domain.SlashCommand]domain.TemplateRef,
) *SlashCommandsInitializer {
	return &SlashCommandsInitializer{
		dir:      dir,
		commands: commands,
	}
}

// Init creates Markdown slash command files in project filesystem.
//
//nolint:revive // argument-limit: interface requires 5 args for dual-fs support
func (s *SlashCommandsInitializer) Init(
	_ context.Context,
	projectFs, _ afero.Fs,
	cfg *Config,
	_ *templates.TemplateManager,
) (InitResult, error) {
	return createSlashCommands(projectFs, s.dir, s.commands, cfg, mdExtension)
}

// IsSetup returns true if all slash command files exist in project
// filesystem.
func (s *SlashCommandsInitializer) IsSetup(
	projectFs, _ afero.Fs,
	_ *Config,
) bool {
	return checkSlashCommandsExist(
		projectFs, s.dir, s.commands, mdExtension,
	)
}

// dedupeKey returns a unique key for deduplication.
func (s *SlashCommandsInitializer) dedupeKey() string {
	return fmt.Sprintf("SlashCommandsInitializer:%s", filepath.Clean(s.dir))
}

// HomeSlashCommandsInitializer creates Markdown slash command files in home
// filesystem. Uses early binding with map[SlashCommand]TemplateRef.
type HomeSlashCommandsInitializer struct {
	dir      string                                     // directory for slash commands
	commands map[domain.SlashCommand]domain.TemplateRef // command templates
}

// NewHomeSlashCommandsInitializer creates a new
// HomeSlashCommandsInitializer for home filesystem.
func NewHomeSlashCommandsInitializer(
	dir string,
	commands map[domain.SlashCommand]domain.TemplateRef,
) *HomeSlashCommandsInitializer {
	return &HomeSlashCommandsInitializer{
		dir:      dir,
		commands: commands,
	}
}

// Init creates Markdown slash command files in home filesystem.
//
//nolint:revive // argument-limit: interface requires 5 args for dual-fs support
func (h *HomeSlashCommandsInitializer) Init(
	_ context.Context,
	_, homeFs afero.Fs,
	cfg *Config,
	_ *templates.TemplateManager,
) (InitResult, error) {
	return createSlashCommands(homeFs, h.dir, h.commands, cfg, mdExtension)
}

// IsSetup returns true if all slash command files exist in home
// filesystem.
func (h *HomeSlashCommandsInitializer) IsSetup(
	_, homeFs afero.Fs,
	_ *Config,
) bool {
	return checkSlashCommandsExist(homeFs, h.dir, h.commands, mdExtension)
}

// dedupeKey returns a unique key for deduplication.
func (h *HomeSlashCommandsInitializer) dedupeKey() string {
	return fmt.Sprintf(
		"HomeSlashCommandsInitializer:%s",
		filepath.Clean(h.dir),
	)
}

// PrefixedSlashCommandsInitializer creates Markdown slash command files with
// custom prefix in project filesystem.
// Output: {dir}/{prefix}{command}.md
// (e.g., .agent/workflows/spectr-proposal.md)
type PrefixedSlashCommandsInitializer struct {
	dir    string // directory for slash commands
	prefix string // prefix for filenames (e.g., "spectr-")
	// command templates
	commands map[domain.SlashCommand]domain.TemplateRef
}

// NewPrefixedSlashCommandsInitializer creates a new
// PrefixedSlashCommandsInitializer.
func NewPrefixedSlashCommandsInitializer(
	dir, prefix string,
	commands map[domain.SlashCommand]domain.TemplateRef,
) *PrefixedSlashCommandsInitializer {
	return &PrefixedSlashCommandsInitializer{
		dir:      dir,
		prefix:   prefix,
		commands: commands,
	}
}

// Init creates prefixed Markdown slash command files in project filesystem.
//
//nolint:revive // argument-limit: interface requires 5 args for dual-fs support
func (p *PrefixedSlashCommandsInitializer) Init(
	_ context.Context,
	projectFs, _ afero.Fs,
	cfg *Config,
	_ *templates.TemplateManager,
) (InitResult, error) {
	return createPrefixedSlashCommands(
		projectFs, p.dir, p.prefix, p.commands, cfg, mdExtension,
	)
}

// IsSetup returns true if all prefixed slash command files exist in
// project filesystem.
func (p *PrefixedSlashCommandsInitializer) IsSetup(
	projectFs, _ afero.Fs,
	_ *Config,
) bool {
	return checkPrefixedSlashCommandsExist(
		projectFs, p.dir, p.prefix, p.commands, mdExtension,
	)
}

// dedupeKey returns a unique key for deduplication.
func (p *PrefixedSlashCommandsInitializer) dedupeKey() string {
	return fmt.Sprintf(
		"PrefixedSlashCommandsInitializer:%s:%s",
		filepath.Clean(p.dir),
		p.prefix,
	)
}

// HomePrefixedSlashCommandsInitializer creates Markdown slash command files
// with custom prefix in home filesystem.
// Output: {dir}/{prefix}{command}.md
// (e.g., ~/.codex/prompts/spectr-proposal.md)
type HomePrefixedSlashCommandsInitializer struct {
	dir    string // directory for slash commands
	prefix string // prefix for filenames (e.g., "spectr-")
	// command templates
	commands map[domain.SlashCommand]domain.TemplateRef
}

// NewHomePrefixedSlashCommandsInitializer creates a new
// HomePrefixedSlashCommandsInitializer.
func NewHomePrefixedSlashCommandsInitializer(
	dir, prefix string,
	commands map[domain.SlashCommand]domain.TemplateRef,
) *HomePrefixedSlashCommandsInitializer {
	return &HomePrefixedSlashCommandsInitializer{
		dir:      dir,
		prefix:   prefix,
		commands: commands,
	}
}

// Init creates prefixed Markdown slash command files in home filesystem.
//
//nolint:revive // argument-limit: interface requires 5 args for dual-fs support
func (h *HomePrefixedSlashCommandsInitializer) Init(
	_ context.Context,
	_, homeFs afero.Fs,
	cfg *Config,
	_ *templates.TemplateManager,
) (InitResult, error) {
	return createPrefixedSlashCommands(
		homeFs, h.dir, h.prefix, h.commands, cfg, mdExtension,
	)
}

// IsSetup returns true if all prefixed slash command files exist in home
// filesystem.
func (h *HomePrefixedSlashCommandsInitializer) IsSetup(
	_, homeFs afero.Fs,
	_ *Config,
) bool {
	return checkPrefixedSlashCommandsExist(
		homeFs, h.dir, h.prefix, h.commands, mdExtension,
	)
}

// dedupeKey returns a unique key for deduplication.
func (h *HomePrefixedSlashCommandsInitializer) dedupeKey() string {
	return fmt.Sprintf(
		"HomePrefixedSlashCommandsInitializer:%s:%s",
		filepath.Clean(h.dir),
		h.prefix,
	)
}

// TOMLSlashCommandsInitializer creates TOML slash command files in project filesystem.
// Uses early binding with map[SlashCommand]TemplateRef.
type TOMLSlashCommandsInitializer struct {
	dir      string                                     // directory for slash commands
	commands map[domain.SlashCommand]domain.TemplateRef // command templates
}

// NewTOMLSlashCommandsInitializer creates a new
// TOMLSlashCommandsInitializer.
func NewTOMLSlashCommandsInitializer(
	dir string,
	commands map[domain.SlashCommand]domain.TemplateRef,
) *TOMLSlashCommandsInitializer {
	return &TOMLSlashCommandsInitializer{
		dir:      dir,
		commands: commands,
	}
}

// Init creates TOML slash command files in project filesystem.
//
//nolint:revive // argument-limit: interface requires 5 args for dual-fs support
func (t *TOMLSlashCommandsInitializer) Init(
	_ context.Context,
	projectFs, _ afero.Fs,
	cfg *Config,
	_ *templates.TemplateManager,
) (InitResult, error) {
	return createSlashCommands(projectFs, t.dir, t.commands, cfg, ".toml")
}

// IsSetup returns true if all TOML slash command files exist in project
// filesystem.
func (t *TOMLSlashCommandsInitializer) IsSetup(
	projectFs, _ afero.Fs,
	_ *Config,
) bool {
	return checkSlashCommandsExist(projectFs, t.dir, t.commands, ".toml")
}

// dedupeKey returns a unique key for deduplication.
func (t *TOMLSlashCommandsInitializer) dedupeKey() string {
	return fmt.Sprintf(
		"TOMLSlashCommandsInitializer:%s",
		filepath.Clean(t.dir),
	)
}

// createSlashCommands is a helper function to create slash command files.
// Always overwrites existing files (idempotent behavior).
func createSlashCommands(
	fs afero.Fs,
	dir string,
	commands map[domain.SlashCommand]domain.TemplateRef,
	cfg *Config,
	ext string,
) (InitResult, error) {
	var created, updated []string

	// Create template context
	templateCtx := domain.TemplateContext{
		BaseDir:     cfg.SpectrDir,
		SpecsDir:    cfg.SpecsDir(),
		ChangesDir:  cfg.ChangesDir(),
		ProjectFile: cfg.ProjectFile(),
		AgentsFile:  cfg.AgentsFile(),
	}

	for cmd, tmpl := range commands {
		filename := cmd.String() + ext
		path := filepath.Join(dir, filename)

		// Render template
		content, err := tmpl.Render(templateCtx)
		if err != nil {
			return InitResult{}, fmt.Errorf(
				"failed to render template for %s: %w",
				filename, err,
			)
		}

		// Check if file exists
		exists, err := afero.Exists(fs, path)
		if err != nil {
			return InitResult{}, fmt.Errorf(
				"failed to check file existence for %s: %w",
				path, err,
			)
		}

		// Always write file (idempotent - overwrite if exists)
		//nolint:revive // octal-literal: file permissions use octal
		if err := afero.WriteFile(fs, path, []byte(content), 0644); err != nil {
			return InitResult{}, fmt.Errorf(
				"failed to write file %s: %w",
				path, err,
			)
		}

		if exists {
			updated = append(updated, path)
		} else {
			created = append(created, path)
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
		filename := cmd.String() + ext
		path := filepath.Join(dir, filename)
		exists, err := afero.Exists(fs, path)
		if err != nil || !exists {
			return false
		}
	}

	return true
}

// createPrefixedSlashCommands is a helper function to create prefixed
// slash command files. Always overwrites existing files (idempotent behavior).
//
//nolint:revive // argument-limit: helper func needs fs, dir, prefix, cmds, cfg, ext
func createPrefixedSlashCommands(
	fs afero.Fs,
	dir, prefix string,
	commands map[domain.SlashCommand]domain.TemplateRef,
	cfg *Config,
	ext string,
) (InitResult, error) {
	var created, updated []string

	// Create template context
	templateCtx := domain.TemplateContext{
		BaseDir:     cfg.SpectrDir,
		SpecsDir:    cfg.SpecsDir(),
		ChangesDir:  cfg.ChangesDir(),
		ProjectFile: cfg.ProjectFile(),
		AgentsFile:  cfg.AgentsFile(),
	}

	for cmd, tmpl := range commands {
		filename := prefix + cmd.String() + ext
		path := filepath.Join(dir, filename)

		// Render template
		content, err := tmpl.Render(templateCtx)
		if err != nil {
			return InitResult{}, fmt.Errorf(
				"failed to render template for %s: %w",
				filename, err,
			)
		}

		// Check if file exists
		exists, err := afero.Exists(fs, path)
		if err != nil {
			return InitResult{}, fmt.Errorf(
				"failed to check file existence for %s: %w",
				path, err,
			)
		}

		// Always write file (idempotent - overwrite if exists)
		//nolint:revive // octal-literal: file permissions use octal
		if err := afero.WriteFile(fs, path, []byte(content), 0644); err != nil {
			return InitResult{}, fmt.Errorf(
				"failed to write file %s: %w",
				path, err,
			)
		}

		if exists {
			updated = append(updated, path)
		} else {
			created = append(created, path)
		}
	}

	return InitResult{
		CreatedFiles: created,
		UpdatedFiles: updated,
	}, nil
}

// checkPrefixedSlashCommandsExist checks if all prefixed slash command
// files exist.
func checkPrefixedSlashCommandsExist(
	fs afero.Fs,
	dir, prefix string,
	commands map[domain.SlashCommand]domain.TemplateRef,
	ext string,
) bool {
	for cmd := range commands {
		filename := prefix + cmd.String() + ext
		path := filepath.Join(dir, filename)
		exists, err := afero.Exists(fs, path)
		if err != nil || !exists {
			return false
		}
	}

	return true
}
