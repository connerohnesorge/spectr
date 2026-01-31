package initialize

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"text/template"

	"github.com/connerohnesorge/spectr/internal/domain"
)

//go:embed templates/**/*.tmpl
var templateFS embed.FS

//go:embed templates/skills
var skillFS embed.FS

// TemplateManager manages embedded templates for initialization
type TemplateManager struct {
	templates *template.Template
}

// NewTemplateManager creates a new template manager with all
// embedded templates loaded.
// It merges templates from:
// 1. internal/initialize/templates (main templates: AGENTS.md, instruction-pointer.md)
// 2. internal/domain (slash command templates: slash-proposal.md, slash-apply.md, TOML variants)
func NewTemplateManager() (*TemplateManager, error) {
	// Parse main templates
	mainTmpl, err := template.ParseFS(
		templateFS,
		"templates/**/*.tmpl",
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse main templates: %w",
			err,
		)
	}

	// Parse and merge domain templates (slash commands)
	domainTmpl, err := template.ParseFS(
		domain.TemplateFS,
		"templates/*.tmpl",
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse domain templates: %w",
			err,
		)
	}

	// Merge: add domain templates to main template set
	// If duplicate template names exist, last-wins precedence applies
	for _, t := range domainTmpl.Templates() {
		if _, err := mainTmpl.AddParseTree(t.Name(), t.Tree); err != nil {
			return nil, fmt.Errorf(
				"failed to merge template %s: %w",
				t.Name(),
				err,
			)
		}
	}

	return &TemplateManager{
		templates: mainTmpl,
	}, nil
}

// RenderProject renders the project.md template with the given context
func (tm *TemplateManager) RenderProject(
	ctx ProjectContext,
) (string, error) {
	var buf bytes.Buffer
	// Template names in ParseFS include the full path from the embed directive
	err := tm.templates.ExecuteTemplate(
		&buf,
		"project.md.tmpl",
		ctx,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to render project template: %w",
			err,
		)
	}

	return buf.String(), nil
}

// RenderAgents renders the AGENTS.md template with the given template context
// The context provides path variables for dynamic directory names
func (tm *TemplateManager) RenderAgents(
	ctx *domain.TemplateContext,
) (string, error) {
	var buf bytes.Buffer
	err := tm.templates.ExecuteTemplate(
		&buf,
		"AGENTS.md.tmpl",
		ctx,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to render agents template: %w",
			err,
		)
	}

	return buf.String(), nil
}

// RenderInstructionPointer renders the instruction-pointer.md template
// This is a short pointer that directs AI assistants to read the AGENTS.md file
// The context provides path variables for dynamic directory names
func (tm *TemplateManager) RenderInstructionPointer(
	ctx *domain.TemplateContext,
) (string, error) {
	var buf bytes.Buffer
	err := tm.templates.ExecuteTemplate(
		&buf,
		"instruction-pointer.md.tmpl",
		ctx,
	)
	if err != nil {
		return "",
			fmt.Errorf(
				"failed to render instruction pointer template: %w",
				err,
			)
	}

	return buf.String(), nil
}

// RenderSlashCommand renders a slash command template with the given context
// commandType must be one of: "proposal", "apply", "archive"
// The context provides path variables for dynamic directory names
func (tm *TemplateManager) RenderSlashCommand(
	commandType string,
	ctx *domain.TemplateContext,
) (string, error) {
	templateName := fmt.Sprintf(
		"slash-%s.md.tmpl",
		commandType,
	)
	var buf bytes.Buffer
	err := tm.templates.ExecuteTemplate(
		&buf,
		templateName,
		ctx,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to render slash command template %s: %w",
			commandType,
			err,
		)
	}

	return buf.String(), nil
}

// RenderCIWorkflow renders the spectr-ci.yml template for GitHub Actions
// This template has no variables and returns the CI workflow configuration
func (tm *TemplateManager) RenderCIWorkflow() (string, error) {
	var buf bytes.Buffer
	err := tm.templates.ExecuteTemplate(
		&buf,
		"spectr-ci.yml.tmpl",
		nil,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to render CI workflow template: %w",
			err,
		)
	}

	return buf.String(), nil
}

// InstructionPointer returns a type-safe reference to the instruction-pointer.md.tmpl template
func (tm *TemplateManager) InstructionPointer() domain.TemplateRef {
	return domain.TemplateRef{
		Name:     "instruction-pointer.md.tmpl",
		Template: tm.templates,
	}
}

// Agents returns a type-safe reference to the AGENTS.md.tmpl template
func (tm *TemplateManager) Agents() domain.TemplateRef {
	return domain.TemplateRef{
		Name:     "AGENTS.md.tmpl",
		Template: tm.templates,
	}
}

// SlashCommand returns a Markdown template reference for the given slash command type.
// Used by SlashCommandsInitializer, HomeSlashCommandsInitializer, and PrefixedSlashCommandsInitializer.
// The returned TemplateRef will assemble frontmatter from BaseSlashCommandFrontmatter.
func (tm *TemplateManager) SlashCommand(
	cmd domain.SlashCommand,
) domain.TemplateRef {
	names := map[domain.SlashCommand]string{
		domain.SlashProposal: "slash-proposal.md.tmpl",
		domain.SlashApply:    "slash-apply.md.tmpl",
		domain.SlashNext:     "slash-next.md.tmpl",
	}

	return domain.TemplateRef{
		Name:     names[cmd],
		Template: tm.templates,
		Command:  &cmd,
	}
}

// SlashCommandWithOverrides returns a Markdown template with frontmatter overrides.
// Used when providers need to customize slash command frontmatter.
// If overrides is nil, behaves identically to SlashCommand(cmd).
func (tm *TemplateManager) SlashCommandWithOverrides(
	cmd domain.SlashCommand,
	overrides *domain.FrontmatterOverride,
) domain.TemplateRef {
	ref := tm.SlashCommand(cmd)
	ref.Overrides = overrides

	return ref
}

// TOMLSlashCommand returns a TOML template reference for the given slash command type.
// Used by TOMLSlashCommandsInitializer (Gemini only).
// Note: TOML templates have frontmatter embedded in the template file, not assembled dynamically.
func (tm *TemplateManager) TOMLSlashCommand(
	cmd domain.SlashCommand,
) domain.TemplateRef {
	names := map[domain.SlashCommand]string{
		domain.SlashProposal: "slash-proposal.toml.tmpl",
		domain.SlashApply:    "slash-apply.toml.tmpl",
		domain.SlashNext:     "slash-next.toml.tmpl",
	}

	// TOML templates don't use dynamic frontmatter (Command is nil)
	return domain.TemplateRef{
		Name:     names[cmd],
		Template: tm.templates,
	}
}

// ProposalSkill returns a type-safe reference to the skill-proposal.md.tmpl template.
// This template is used for Amp agent skills that create change proposals.
func (tm *TemplateManager) ProposalSkill() domain.TemplateRef {
	return domain.TemplateRef{
		Name:     "skill-proposal.md.tmpl",
		Template: tm.templates,
	}
}

// ApplySkill returns a type-safe reference to the skill-apply.md.tmpl template.
// This template is used for Amp agent skills that apply/accept change proposals.
func (tm *TemplateManager) ApplySkill() domain.TemplateRef {
	return domain.TemplateRef{
		Name:     "skill-apply.md.tmpl",
		Template: tm.templates,
	}
}

// NextSkill returns a type-safe reference to the skill-next.md.tmpl template.
// This template is used for agent skills that execute the next pending task.
func (tm *TemplateManager) NextSkill() domain.TemplateRef {
	return domain.TemplateRef{
		Name:     "skill-next.md.tmpl",
		Template: tm.templates,
	}
}

// SkillFS returns an fs.FS rooted at the skill directory for the given skill name.
// Returns an error if the skill does not exist.
// The filesystem contains all files under templates/skills/<skillName>/ with paths
// relative to the skill root (e.g., SKILL.md, scripts/accept.sh).
//
//nolint:revive // receiver not used but required by TemplateManager interface
func (*TemplateManager) SkillFS(
	skillName string,
) (fs.FS, error) {
	// Create a sub-filesystem rooted at templates/skills/<skillName>
	skillPath := fmt.Sprintf(
		"templates/skills/%s",
		skillName,
	)
	subFS, err := fs.Sub(skillFS, skillPath)
	if err != nil {
		return nil, fmt.Errorf(
			"skill %s not found: %w",
			skillName,
			err,
		)
	}

	return subFS, nil
}
