package providers

import (
	"context"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/spf13/afero"
)

func TestSlashCommandsInitializer_Init(t *testing.T) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	// Create templates
	proposalTmpl := template.Must(template.New("proposal.md.tmpl").Parse("Proposal content"))
	applyTmpl := template.Must(template.New("apply.md.tmpl").Parse("Apply content"))

	commands := map[domain.SlashCommand]domain.TemplateRef{
		domain.SlashProposal: {Name: "proposal.md.tmpl", Template: proposalTmpl},
		domain.SlashApply:    {Name: "apply.md.tmpl", Template: applyTmpl},
	}

	// Test
	init := NewSlashCommandsInitializer(".claude/commands/spectr", commands)
	result, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)

	// Verify
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	if len(result.CreatedFiles) != 2 {
		t.Errorf("CreatedFiles count = %d, want 2", len(result.CreatedFiles))
	}

	if len(result.UpdatedFiles) != 0 {
		t.Errorf("UpdatedFiles count = %d, want 0", len(result.UpdatedFiles))
	}

	// Verify files exist with correct extension
	proposalPath := filepath.Join(".claude/commands/spectr", "proposal.md")
	applyPath := filepath.Join(".claude/commands/spectr", "apply.md")

	exists, _ := afero.Exists(projectFs, proposalPath)
	if !exists {
		t.Errorf("file %s does not exist", proposalPath)
	}

	exists, _ = afero.Exists(projectFs, applyPath)
	if !exists {
		t.Errorf("file %s does not exist", applyPath)
	}

	// Verify content
	content, _ := afero.ReadFile(projectFs, proposalPath)
	if string(content) != "Proposal content" {
		t.Errorf("proposal content = %q, want %q", string(content), "Proposal content")
	}

	content, _ = afero.ReadFile(projectFs, applyPath)
	if string(content) != "Apply content" {
		t.Errorf("apply content = %q, want %q", string(content), "Apply content")
	}
}

func TestSlashCommandsInitializer_Init_Idempotent(t *testing.T) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	// Create templates
	proposalTmpl := template.Must(template.New("proposal.md.tmpl").Parse("New content"))

	commands := map[domain.SlashCommand]domain.TemplateRef{
		domain.SlashProposal: {Name: "proposal.md.tmpl", Template: proposalTmpl},
	}

	// Create existing file with old content
	proposalPath := filepath.Join(".claude/commands/spectr", "proposal.md")
	if err := projectFs.MkdirAll(".claude/commands/spectr", 0755); err != nil {
		t.Fatalf("failed to create dir: %v", err)
	}
	if err := afero.WriteFile(projectFs, proposalPath, []byte("Old content"), 0644); err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}

	// Test
	init := NewSlashCommandsInitializer(".claude/commands/spectr", commands)
	result, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)

	// Verify
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Should be reported as updated, not created
	if len(result.CreatedFiles) != 0 {
		t.Errorf("CreatedFiles count = %d, want 0", len(result.CreatedFiles))
	}

	if len(result.UpdatedFiles) != 1 {
		t.Errorf("UpdatedFiles count = %d, want 1", len(result.UpdatedFiles))
	}

	// Verify content was overwritten
	content, _ := afero.ReadFile(projectFs, proposalPath)
	if string(content) != "New content" {
		t.Errorf("content = %q, want %q (should be overwritten)", string(content), "New content")
	}
}

func TestSlashCommandsInitializer_IsSetup(t *testing.T) {
	tests := []struct {
		name          string
		existingFiles []string
		want          bool
	}{
		{
			name:          "returns true if all files exist",
			existingFiles: []string{"proposal.md", "apply.md"},
			want:          true,
		},
		{
			name:          "returns false if any file missing",
			existingFiles: []string{"proposal.md"},
			want:          false,
		},
		{
			name:          "returns false if no files exist",
			existingFiles: nil,
			want:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			projectFs := afero.NewMemMapFs()
			homeFs := afero.NewMemMapFs()
			cfg := &Config{SpectrDir: "spectr"}

			// Create existing files
			if err := projectFs.MkdirAll(".claude/commands/spectr", 0755); err != nil {
				t.Fatalf("failed to create dir: %v", err)
			}
			for _, file := range tt.existingFiles {
				path := filepath.Join(".claude/commands/spectr", file)
				if err := afero.WriteFile(projectFs, path, []byte("content"), 0644); err != nil {
					t.Fatalf("failed to create file: %v", err)
				}
			}

			// Create templates
			proposalTmpl := template.Must(template.New("proposal.md.tmpl").Parse("Proposal"))
			applyTmpl := template.Must(template.New("apply.md.tmpl").Parse("Apply"))

			commands := map[domain.SlashCommand]domain.TemplateRef{
				domain.SlashProposal: {Name: "proposal.md.tmpl", Template: proposalTmpl},
				domain.SlashApply:    {Name: "apply.md.tmpl", Template: applyTmpl},
			}

			// Test
			init := NewSlashCommandsInitializer(".claude/commands/spectr", commands)
			got := init.IsSetup(projectFs, homeFs, cfg)

			if got != tt.want {
				t.Errorf("IsSetup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlashCommandsInitializer_DedupeKey(t *testing.T) {
	// Setup
	proposalTmpl := template.Must(template.New("proposal.md.tmpl").Parse("Proposal"))
	commands := map[domain.SlashCommand]domain.TemplateRef{
		domain.SlashProposal: {Name: "proposal.md.tmpl", Template: proposalTmpl},
	}

	init := NewSlashCommandsInitializer(".claude/commands/spectr", commands)
	got := init.dedupeKey()
	want := "SlashCommandsInitializer:.claude/commands/spectr"

	if got != want {
		t.Errorf("dedupeKey() = %v, want %v", got, want)
	}
}

func TestHomeSlashCommandsInitializer_Init(t *testing.T) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	// Create templates
	proposalTmpl := template.Must(template.New("proposal.md.tmpl").Parse("Proposal content"))

	commands := map[domain.SlashCommand]domain.TemplateRef{
		domain.SlashProposal: {Name: "proposal.md.tmpl", Template: proposalTmpl},
	}

	// Test
	init := NewHomeSlashCommandsInitializer(".codex/prompts", commands)
	result, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)

	// Verify
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	if len(result.CreatedFiles) != 1 {
		t.Errorf("CreatedFiles count = %d, want 1", len(result.CreatedFiles))
	}

	// Verify file exists in homeFs (not projectFs)
	proposalPath := filepath.Join(".codex/prompts", "proposal.md")
	exists, _ := afero.Exists(homeFs, proposalPath)
	if !exists {
		t.Errorf("file %s does not exist in homeFs", proposalPath)
	}

	// Verify it doesn't exist in projectFs
	exists, _ = afero.Exists(projectFs, proposalPath)
	if exists {
		t.Errorf("file %s should not exist in projectFs", proposalPath)
	}
}

func TestHomeSlashCommandsInitializer_DedupeKey(t *testing.T) {
	// Setup
	proposalTmpl := template.Must(template.New("proposal.md.tmpl").Parse("Proposal"))
	commands := map[domain.SlashCommand]domain.TemplateRef{
		domain.SlashProposal: {Name: "proposal.md.tmpl", Template: proposalTmpl},
	}

	init := NewHomeSlashCommandsInitializer(".codex/prompts", commands)
	got := init.dedupeKey()
	want := "HomeSlashCommandsInitializer:.codex/prompts"

	if got != want {
		t.Errorf("dedupeKey() = %v, want %v", got, want)
	}
}

func TestPrefixedSlashCommandsInitializer_Init(t *testing.T) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	// Create templates
	proposalTmpl := template.Must(template.New("proposal.md.tmpl").Parse("Proposal content"))
	applyTmpl := template.Must(template.New("apply.md.tmpl").Parse("Apply content"))

	commands := map[domain.SlashCommand]domain.TemplateRef{
		domain.SlashProposal: {Name: "proposal.md.tmpl", Template: proposalTmpl},
		domain.SlashApply:    {Name: "apply.md.tmpl", Template: applyTmpl},
	}

	// Test
	init := NewPrefixedSlashCommandsInitializer(".agent/workflows", "spectr-", commands)
	result, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)

	// Verify
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	if len(result.CreatedFiles) != 2 {
		t.Errorf("CreatedFiles count = %d, want 2", len(result.CreatedFiles))
	}

	// Verify files exist with prefix
	proposalPath := filepath.Join(".agent/workflows", "spectr-proposal.md")
	applyPath := filepath.Join(".agent/workflows", "spectr-apply.md")

	exists, _ := afero.Exists(projectFs, proposalPath)
	if !exists {
		t.Errorf("file %s does not exist", proposalPath)
	}

	exists, _ = afero.Exists(projectFs, applyPath)
	if !exists {
		t.Errorf("file %s does not exist", applyPath)
	}

	// Verify content
	content, _ := afero.ReadFile(projectFs, proposalPath)
	if string(content) != "Proposal content" {
		t.Errorf("proposal content = %q, want %q", string(content), "Proposal content")
	}
}

func TestPrefixedSlashCommandsInitializer_DedupeKey(t *testing.T) {
	// Setup
	proposalTmpl := template.Must(template.New("proposal.md.tmpl").Parse("Proposal"))
	commands := map[domain.SlashCommand]domain.TemplateRef{
		domain.SlashProposal: {Name: "proposal.md.tmpl", Template: proposalTmpl},
	}

	init := NewPrefixedSlashCommandsInitializer(".agent/workflows", "spectr-", commands)
	got := init.dedupeKey()
	want := "PrefixedSlashCommandsInitializer:.agent/workflows:spectr-"

	if got != want {
		t.Errorf("dedupeKey() = %v, want %v", got, want)
	}
}

func TestHomePrefixedSlashCommandsInitializer_Init(t *testing.T) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	// Create templates
	proposalTmpl := template.Must(template.New("proposal.md.tmpl").Parse("Proposal content"))

	commands := map[domain.SlashCommand]domain.TemplateRef{
		domain.SlashProposal: {Name: "proposal.md.tmpl", Template: proposalTmpl},
	}

	// Test
	init := NewHomePrefixedSlashCommandsInitializer(".codex/prompts", "spectr-", commands)
	result, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)

	// Verify
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	if len(result.CreatedFiles) != 1 {
		t.Errorf("CreatedFiles count = %d, want 1", len(result.CreatedFiles))
	}

	// Verify file exists in homeFs with prefix
	proposalPath := filepath.Join(".codex/prompts", "spectr-proposal.md")
	exists, _ := afero.Exists(homeFs, proposalPath)
	if !exists {
		t.Errorf("file %s does not exist in homeFs", proposalPath)
	}

	// Verify it doesn't exist in projectFs
	exists, _ = afero.Exists(projectFs, proposalPath)
	if exists {
		t.Errorf("file %s should not exist in projectFs", proposalPath)
	}
}

func TestHomePrefixedSlashCommandsInitializer_DedupeKey(t *testing.T) {
	// Setup
	proposalTmpl := template.Must(template.New("proposal.md.tmpl").Parse("Proposal"))
	commands := map[domain.SlashCommand]domain.TemplateRef{
		domain.SlashProposal: {Name: "proposal.md.tmpl", Template: proposalTmpl},
	}

	init := NewHomePrefixedSlashCommandsInitializer(".codex/prompts", "spectr-", commands)
	got := init.dedupeKey()
	want := "HomePrefixedSlashCommandsInitializer:.codex/prompts:spectr-"

	if got != want {
		t.Errorf("dedupeKey() = %v, want %v", got, want)
	}
}

func TestTOMLSlashCommandsInitializer_Init(t *testing.T) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	// Create templates
	proposalTmpl := template.Must(
		template.New("proposal.toml.tmpl").Parse(`description = "Create proposal"
prompt = """
Proposal content
"""`),
	)
	applyTmpl := template.Must(template.New("apply.toml.tmpl").Parse(`description = "Apply changes"
prompt = """
Apply content
"""`))

	commands := map[domain.SlashCommand]domain.TemplateRef{
		domain.SlashProposal: {Name: "proposal.toml.tmpl", Template: proposalTmpl},
		domain.SlashApply:    {Name: "apply.toml.tmpl", Template: applyTmpl},
	}

	// Test
	init := NewTOMLSlashCommandsInitializer(".gemini/commands/spectr", commands)
	result, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)

	// Verify
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	if len(result.CreatedFiles) != 2 {
		t.Errorf("CreatedFiles count = %d, want 2", len(result.CreatedFiles))
	}

	// Verify files exist with .toml extension
	proposalPath := filepath.Join(".gemini/commands/spectr", "proposal.toml")
	applyPath := filepath.Join(".gemini/commands/spectr", "apply.toml")

	exists, _ := afero.Exists(projectFs, proposalPath)
	if !exists {
		t.Errorf("file %s does not exist", proposalPath)
	}

	exists, _ = afero.Exists(projectFs, applyPath)
	if !exists {
		t.Errorf("file %s does not exist", applyPath)
	}

	// Verify content
	content, _ := afero.ReadFile(projectFs, proposalPath)
	expectedContent := `description = "Create proposal"
prompt = """
Proposal content
"""`
	if string(content) != expectedContent {
		t.Errorf("proposal content = %q, want %q", string(content), expectedContent)
	}
}

func TestTOMLSlashCommandsInitializer_DedupeKey(t *testing.T) {
	// Setup
	proposalTmpl := template.Must(template.New("proposal.toml.tmpl").Parse("content"))
	commands := map[domain.SlashCommand]domain.TemplateRef{
		domain.SlashProposal: {Name: "proposal.toml.tmpl", Template: proposalTmpl},
	}

	init := NewTOMLSlashCommandsInitializer(".gemini/commands/spectr", commands)
	got := init.dedupeKey()
	want := "TOMLSlashCommandsInitializer:.gemini/commands/spectr"

	if got != want {
		t.Errorf("dedupeKey() = %v, want %v", got, want)
	}
}
