package initialize

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"strings"
	"text/template"

	"github.com/connerohnesorge/spectr/internal/domain"
)

//go:embed templates/**/*.tmpl templates/providers/**/*.tmpl
var templateFS embed.FS

//go:embed templates/skills
var skillFS embed.FS

// TemplateManager manages embedded templates for initialization
type TemplateManager struct {
	templates         *template.Template
	slashTemplates    map[string]*template.Template
	providerTemplates map[string]map[string]*template.Template
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

	slashTemplates, err := parseSlashTemplates()
	if err != nil {
		return nil, err
	}

	providerTmpls, err := loadProviderTemplates(slashTemplates, templateFS)
	if err != nil {
		return nil, err
	}

	return &TemplateManager{
		templates:         mainTmpl,
		slashTemplates:    slashTemplates,
		providerTemplates: providerTmpls,
	}, nil
}

func validateProviderTemplate(_, provider *template.Template) error {
	knownSections := map[string]struct{}{
		"guardrails":      {},
		"steps":           {},
		"reference":       {},
		"main":            {},
		"base_guardrails": {},
		"base_steps":      {},
		"base_reference":  {},
	}

	for _, tmpl := range provider.Templates() {
		name := tmpl.Name()
		if strings.HasSuffix(name, ".tmpl") {
			continue
		}
		if _, ok := knownSections[name]; !ok {
			return fmt.Errorf(
				"unknown section %q (allowed: guardrails, steps, reference, main, base_guardrails, base_steps, base_reference)",
				name,
			)
		}
	}

	return nil
}

func loadProviderTemplates(
	baseTemplates map[string]*template.Template,
	providerFS fs.FS,
) (map[string]map[string]*template.Template, error) {
	providerTmpls := make(map[string]map[string]*template.Template)
	entries, err := fs.ReadDir(providerFS, "templates/providers")
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return providerTmpls, nil
		}

		return nil, fmt.Errorf(
			"failed to read provider template directory: %w",
			err,
		)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		providerID := entry.Name()
		providerPath := fmt.Sprintf(
			"templates/providers/%s",
			providerID,
		)
		providerTemplates, err := loadProviderTemplatesForProvider(
			providerID,
			providerPath,
			baseTemplates,
			providerFS,
		)
		if err != nil {
			return nil, err
		}
		if len(providerTemplates) == 0 {
			continue
		}

		providerTmpls[providerID] = providerTemplates
	}

	return providerTmpls, nil
}

func loadProviderTemplatesForProvider(
	providerID, providerPath string,
	baseTemplates map[string]*template.Template,
	providerFS fs.FS,
) (map[string]*template.Template, error) {
	providerEntries, err := fs.ReadDir(providerFS, providerPath)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read provider %s directory: %w",
			providerID,
			err,
		)
	}

	providerTemplates := make(map[string]*template.Template)
	for _, providerEntry := range providerEntries {
		if providerEntry.IsDir() {
			continue
		}
		if !strings.HasSuffix(providerEntry.Name(), ".tmpl") {
			continue
		}

		templateName := providerEntry.Name()
		baseTemplate, ok := baseTemplates[templateName]
		if !ok {
			return nil, fmt.Errorf(
				"unknown base template %s for provider %s",
				templateName,
				providerID,
			)
		}

		pattern := fmt.Sprintf("%s/%s", providerPath, templateName)
		providerTmpl, err := template.ParseFS(
			providerFS,
			pattern,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to parse provider %s template %s: %w",
				providerID,
				templateName,
				err,
			)
		}

		if err := validateProviderTemplate(baseTemplate, providerTmpl); err != nil {
			return nil, fmt.Errorf(
				"provider %s template validation failed: %w",
				providerID,
				err,
			)
		}

		providerTemplates[templateName] = providerTmpl
	}

	return providerTemplates, nil
}

func parseSlashTemplates() (map[string]*template.Template, error) {
	slashTemplateNames := []string{
		"slash-proposal.md.tmpl",
		"slash-apply.md.tmpl",
		"slash-proposal.toml.tmpl",
		"slash-apply.toml.tmpl",
	}

	slashTemplates := make(map[string]*template.Template, len(slashTemplateNames))
	for _, name := range slashTemplateNames {
		pattern := fmt.Sprintf("templates/%s", name)
		tmpl, err := template.ParseFS(domain.TemplateFS, pattern)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to parse domain template %s: %w",
				name,
				err,
			)
		}

		slashTemplates[name] = tmpl
	}

	return slashTemplates, nil
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
	tmpl, ok := tm.slashTemplates[templateName]
	if !ok {
		return "", fmt.Errorf(
			"unknown slash command template %s",
			commandType,
		)
	}
	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(
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
func (tm *TemplateManager) SlashCommand(cmd domain.SlashCommand) domain.TemplateRef {
	return tm.ProviderSlashCommand("", cmd)
}

// ProviderSlashCommand returns a provider-aware template reference for the given slash command.
// Providers without overrides will receive a nil ProviderTemplate and use the base template.
func (tm *TemplateManager) ProviderSlashCommand(
	providerID string,
	cmd domain.SlashCommand,
) domain.TemplateRef {
	names := map[domain.SlashCommand]string{
		domain.SlashProposal: "slash-proposal.md.tmpl",
		domain.SlashApply:    "slash-apply.md.tmpl",
	}

	return domain.TemplateRef{
		Name:             names[cmd],
		Template:         tm.slashTemplates[names[cmd]],
		ProviderTemplate: tm.providerTemplates[providerID][names[cmd]],
	}
}

// TOMLSlashCommand returns a TOML template reference for the given slash command type.
// Used by TOMLSlashCommandsInitializer (Gemini only).
func (tm *TemplateManager) TOMLSlashCommand(cmd domain.SlashCommand) domain.TemplateRef {
	return tm.ProviderTOMLSlashCommand("", cmd)
}

// ProviderTOMLSlashCommand returns a provider-aware TOML template reference for the given slash command.
// Providers without overrides will receive a nil ProviderTemplate and use the base template.
func (tm *TemplateManager) ProviderTOMLSlashCommand(
	providerID string,
	cmd domain.SlashCommand,
) domain.TemplateRef {
	names := map[domain.SlashCommand]string{
		domain.SlashProposal: "slash-proposal.toml.tmpl",
		domain.SlashApply:    "slash-apply.toml.tmpl",
	}

	return domain.TemplateRef{
		Name:             names[cmd],
		Template:         tm.slashTemplates[names[cmd]],
		ProviderTemplate: tm.providerTemplates[providerID][names[cmd]],
	}
}

// SkillFS returns an fs.FS rooted at the skill directory for the given skill name.
// Returns an error if the skill does not exist.
// The filesystem contains all files under templates/skills/<skillName>/ with paths
// relative to the skill root (e.g., SKILL.md, scripts/accept.sh).
//
//nolint:revive // receiver not used but required by TemplateManager interface
func (*TemplateManager) SkillFS(skillName string) (fs.FS, error) {
	// Create a sub-filesystem rooted at templates/skills/<skillName>
	skillPath := fmt.Sprintf("templates/skills/%s", skillName)
	subFS, err := fs.Sub(skillFS, skillPath)
	if err != nil {
		return nil, fmt.Errorf("skill %s not found: %w", skillName, err)
	}

	return subFS, nil
}
