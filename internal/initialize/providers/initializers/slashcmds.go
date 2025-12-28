package initializers

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/spf13/afero"
)

// SlashCommandsInitializer creates Markdown slash command files in the project filesystem.
// Files are named {command}.md in the specified directory.
type SlashCommandsInitializer struct {
	dir      string
	commands map[domain.SlashCommand]domain.TemplateRef
}

// NewSlashCommandsInitializer creates an initializer that creates Markdown slash command files
// in the project filesystem.
func NewSlashCommandsInitializer(
	dir string,
	commands map[domain.SlashCommand]domain.TemplateRef,
) domain.Initializer {
	return &SlashCommandsInitializer{
		dir:      dir,
		commands: commands,
	}
}

// Init creates slash command files. Always overwrites existing files (idempotent).
//
//nolint:revive // argument-limit - interface signature requires these parameters
func (s *SlashCommandsInitializer) Init(
	_ context.Context,
	projectFs, _ afero.Fs,
	cfg *domain.Config,
	_ any,
) (domain.ExecutionResult, error) {
	return createSlashCommands(projectFs, s.dir, s.commands, extMarkdown, cfg)
}

// IsSetup returns true if all slash command files exist.
func (s *SlashCommandsInitializer) IsSetup(projectFs, _ afero.Fs, _ *domain.Config) bool {
	return checkSlashCommandsExist(projectFs, s.dir, s.commands, extMarkdown)
}

// DedupeKey returns a unique key for deduplication.
// Exported to allow deduplication from the executor package.
func (s *SlashCommandsInitializer) DedupeKey() string {
	return "SlashCommandsInitializer:" + filepath.Clean(s.dir)
}

var _ Deduplicatable = (*SlashCommandsInitializer)(nil)

// HomeSlashCommandsInitializer creates Markdown slash command files in the home filesystem.
// Files are named {command}.md in the specified directory.
type HomeSlashCommandsInitializer struct {
	dir      string
	commands map[domain.SlashCommand]domain.TemplateRef
}

// NewHomeSlashCommandsInitializer creates an initializer that creates Markdown slash command files
// in the home filesystem.
func NewHomeSlashCommandsInitializer(
	dir string,
	commands map[domain.SlashCommand]domain.TemplateRef,
) domain.Initializer {
	return &HomeSlashCommandsInitializer{
		dir:      dir,
		commands: commands,
	}
}

// Init creates slash command files. Always overwrites existing files (idempotent).
//
//nolint:revive // argument-limit - interface signature requires these parameters
func (s *HomeSlashCommandsInitializer) Init(
	_ context.Context,
	_, homeFs afero.Fs,
	cfg *domain.Config,
	_ any,
) (domain.ExecutionResult, error) {
	return createSlashCommands(homeFs, s.dir, s.commands, extMarkdown, cfg)
}

// IsSetup returns true if all slash command files exist.
func (s *HomeSlashCommandsInitializer) IsSetup(_, homeFs afero.Fs, _ *domain.Config) bool {
	return checkSlashCommandsExist(homeFs, s.dir, s.commands, extMarkdown)
}

// DedupeKey returns a unique key for deduplication.
// Exported to allow deduplication from the executor package.
func (s *HomeSlashCommandsInitializer) DedupeKey() string {
	return "HomeSlashCommandsInitializer:" + filepath.Clean(s.dir)
}

var _ Deduplicatable = (*HomeSlashCommandsInitializer)(nil)

// PrefixedSlashCommandsInitializer creates Markdown slash command files with a prefix in the project filesystem.
// Files are named {prefix}{command}.md in the specified directory.
// Used by providers like Antigravity that use non-standard paths.
type PrefixedSlashCommandsInitializer struct {
	dir      string
	prefix   string
	commands map[domain.SlashCommand]domain.TemplateRef
}

// NewPrefixedSlashCommandsInitializer creates an initializer that creates Markdown slash command files
// with a custom prefix in the project filesystem.
func NewPrefixedSlashCommandsInitializer(
	dir, prefix string,
	commands map[domain.SlashCommand]domain.TemplateRef,
) domain.Initializer {
	return &PrefixedSlashCommandsInitializer{
		dir:      dir,
		prefix:   prefix,
		commands: commands,
	}
}

// Init creates prefixed slash command files. Always overwrites existing files (idempotent).
//
//nolint:revive // argument-limit - interface signature requires these parameters
func (s *PrefixedSlashCommandsInitializer) Init(
	_ context.Context,
	projectFs, _ afero.Fs,
	cfg *domain.Config,
	_ any,
) (domain.ExecutionResult, error) {
	return createPrefixedSlashCommands(projectFs, s.dir, s.prefix, s.commands, extMarkdown, cfg)
}

// IsSetup returns true if all slash command files exist.
func (s *PrefixedSlashCommandsInitializer) IsSetup(projectFs, _ afero.Fs, _ *domain.Config) bool {
	return checkPrefixedSlashCommandsExist(projectFs, s.dir, s.prefix, s.commands, extMarkdown)
}

// DedupeKey returns a unique key for deduplication.
// Exported to allow deduplication from the executor package.
func (s *PrefixedSlashCommandsInitializer) DedupeKey() string {
	return "PrefixedSlashCommandsInitializer:" + filepath.Clean(s.dir) + ":" + s.prefix
}

var _ Deduplicatable = (*PrefixedSlashCommandsInitializer)(nil)

// HomePrefixedSlashCommandsInitializer creates Markdown slash command files with a prefix in the home filesystem.
// Files are named {prefix}{command}.md in the specified directory.
// Used by providers like Codex that use home paths with non-standard naming.
type HomePrefixedSlashCommandsInitializer struct {
	dir      string
	prefix   string
	commands map[domain.SlashCommand]domain.TemplateRef
}

// NewHomePrefixedSlashCommandsInitializer creates an initializer that creates Markdown slash command files
// with a custom prefix in the home filesystem.
func NewHomePrefixedSlashCommandsInitializer(
	dir, prefix string,
	commands map[domain.SlashCommand]domain.TemplateRef,
) domain.Initializer {
	return &HomePrefixedSlashCommandsInitializer{
		dir:      dir,
		prefix:   prefix,
		commands: commands,
	}
}

// Init creates prefixed slash command files. Always overwrites existing files (idempotent).
//
//nolint:revive // argument-limit - interface signature requires these parameters
func (s *HomePrefixedSlashCommandsInitializer) Init(
	_ context.Context,
	_, homeFs afero.Fs,
	cfg *domain.Config,
	_ any,
) (domain.ExecutionResult, error) {
	return createPrefixedSlashCommands(homeFs, s.dir, s.prefix, s.commands, extMarkdown, cfg)
}

// IsSetup returns true if all slash command files exist.
func (s *HomePrefixedSlashCommandsInitializer) IsSetup(_, homeFs afero.Fs, _ *domain.Config) bool {
	return checkPrefixedSlashCommandsExist(homeFs, s.dir, s.prefix, s.commands, extMarkdown)
}

// DedupeKey returns a unique key for deduplication.
// Exported to allow deduplication from the executor package.
func (s *HomePrefixedSlashCommandsInitializer) DedupeKey() string {
	return "HomePrefixedSlashCommandsInitializer:" + filepath.Clean(s.dir) + ":" + s.prefix
}

var _ Deduplicatable = (*HomePrefixedSlashCommandsInitializer)(nil)

// TOMLSlashCommandsInitializer creates TOML slash command files in the project filesystem.
// Files are named {command}.toml in the specified directory.
// Used by providers like Gemini that use TOML format.
type TOMLSlashCommandsInitializer struct {
	dir      string
	commands map[domain.SlashCommand]domain.TemplateRef
}

// NewTOMLSlashCommandsInitializer creates an initializer that creates TOML slash command files
// in the project filesystem.
func NewTOMLSlashCommandsInitializer(
	dir string,
	commands map[domain.SlashCommand]domain.TemplateRef,
) domain.Initializer {
	return &TOMLSlashCommandsInitializer{
		dir:      dir,
		commands: commands,
	}
}

// Init creates TOML slash command files. Always overwrites existing files (idempotent).
//
//nolint:revive // argument-limit - interface signature requires these parameters
func (s *TOMLSlashCommandsInitializer) Init(
	_ context.Context,
	projectFs, _ afero.Fs,
	cfg *domain.Config,
	_ any,
) (domain.ExecutionResult, error) {
	return createSlashCommands(projectFs, s.dir, s.commands, extTOML, cfg)
}

// IsSetup returns true if all TOML slash command files exist.
func (s *TOMLSlashCommandsInitializer) IsSetup(projectFs, _ afero.Fs, _ *domain.Config) bool {
	return checkSlashCommandsExist(projectFs, s.dir, s.commands, extTOML)
}

// DedupeKey returns a unique key for deduplication.
// Exported to allow deduplication from the executor package.
func (s *TOMLSlashCommandsInitializer) DedupeKey() string {
	return "TOMLSlashCommandsInitializer:" + filepath.Clean(s.dir)
}

var _ Deduplicatable = (*TOMLSlashCommandsInitializer)(nil)

// File extension constants.
const (
	extMarkdown = ".md"
	extTOML     = ".toml"
)

// createSlashCommands creates slash command files with the given extension.
// Helper function shared by SlashCommandsInitializer, HomeSlashCommandsInitializer, and TOMLSlashCommandsInitializer.
//
//nolint:revive // argument-limit - helper function needs these parameters for flexibility
func createSlashCommands(
	fs afero.Fs,
	dir string,
	commands map[domain.SlashCommand]domain.TemplateRef,
	ext string,
	cfg *domain.Config,
) (domain.ExecutionResult, error) {
	created := make([]string, 0, len(commands))

	// Create template context from config
	templateCtx := domain.TemplateContext{
		BaseDir:     cfg.SpectrDir,
		SpecsDir:    cfg.SpecsDir(),
		ChangesDir:  cfg.ChangesDir(),
		ProjectFile: cfg.ProjectFile(),
		AgentsFile:  cfg.AgentsFile(),
	}

	for cmd, tmplRef := range commands {
		filename := cmd.String() + ext
		path := filepath.Join(dir, filename)

		// Render template
		var buf bytes.Buffer
		if err := tmplRef.Template.ExecuteTemplate(&buf, tmplRef.Name, templateCtx); err != nil {
			return domain.ExecutionResult{
					CreatedFiles: created,
				}, fmt.Errorf(
					"failed to render template %s: %w",
					tmplRef.Name,
					err,
				)
		}

		// Always overwrite (idempotent)
		if err := afero.WriteFile(fs, path, buf.Bytes(), filePerm); err != nil {
			return domain.ExecutionResult{CreatedFiles: created}, err
		}
		created = append(created, path)
	}

	return domain.ExecutionResult{CreatedFiles: created}, nil
}

// createPrefixedSlashCommands creates slash command files with a prefix and extension.
// Helper function shared by PrefixedSlashCommandsInitializer and HomePrefixedSlashCommandsInitializer.
//
//nolint:revive // argument-limit - helper function needs these parameters for flexibility
func createPrefixedSlashCommands(
	fs afero.Fs,
	dir, prefix string,
	commands map[domain.SlashCommand]domain.TemplateRef,
	ext string,
	cfg *domain.Config,
) (domain.ExecutionResult, error) {
	created := make([]string, 0, len(commands))

	// Create template context from config
	templateCtx := domain.TemplateContext{
		BaseDir:     cfg.SpectrDir,
		SpecsDir:    cfg.SpecsDir(),
		ChangesDir:  cfg.ChangesDir(),
		ProjectFile: cfg.ProjectFile(),
		AgentsFile:  cfg.AgentsFile(),
	}

	for cmd, tmplRef := range commands {
		filename := prefix + cmd.String() + ext
		path := filepath.Join(dir, filename)

		// Render template
		var buf bytes.Buffer
		if err := tmplRef.Template.ExecuteTemplate(&buf, tmplRef.Name, templateCtx); err != nil {
			return domain.ExecutionResult{
					CreatedFiles: created,
				}, fmt.Errorf(
					"failed to render template %s: %w",
					tmplRef.Name,
					err,
				)
		}

		// Always overwrite (idempotent)
		if err := afero.WriteFile(fs, path, buf.Bytes(), filePerm); err != nil {
			return domain.ExecutionResult{CreatedFiles: created}, err
		}
		created = append(created, path)
	}

	return domain.ExecutionResult{CreatedFiles: created}, nil
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

// checkPrefixedSlashCommandsExist checks if all prefixed slash command files exist.
//
//nolint:revive // argument-limit - helper function needs these parameters for flexibility
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
