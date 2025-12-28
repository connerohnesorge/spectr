package initializers

import (
	"context"
	"html/template"
	"path/filepath"
	"strings"
	"testing"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/spf13/afero"
)

func newTestSlashTemplates(t *testing.T) map[domain.SlashCommand]domain.TemplateRef {
	t.Helper()

	proposalTmpl, err := template.New("proposal.md.tmpl").Parse("# Proposal\nBaseDir: {{.BaseDir}}")
	if err != nil {
		t.Fatalf("Failed to parse proposal template: %v", err)
	}

	applyTmpl, err := template.New("apply.md.tmpl").Parse("# Apply\nSpecsDir: {{.SpecsDir}}")
	if err != nil {
		t.Fatalf("Failed to parse apply template: %v", err)
	}

	return map[domain.SlashCommand]domain.TemplateRef{
		domain.SlashProposal: {Name: "proposal.md.tmpl", Template: proposalTmpl},
		domain.SlashApply:    {Name: "apply.md.tmpl", Template: applyTmpl},
	}
}

func newTestTOMLTemplates(t *testing.T) map[domain.SlashCommand]domain.TemplateRef {
	t.Helper()

	proposalTmpl, err := template.New("proposal.toml.tmpl").Parse(`description = "Proposal"
prompt = "BaseDir: {{.BaseDir}}"`)
	if err != nil {
		t.Fatalf("Failed to parse proposal toml template: %v", err)
	}

	applyTmpl, err := template.New("apply.toml.tmpl").Parse(`description = "Apply"
prompt = "SpecsDir: {{.SpecsDir}}"`)
	if err != nil {
		t.Fatalf("Failed to parse apply toml template: %v", err)
	}

	return map[domain.SlashCommand]domain.TemplateRef{
		domain.SlashProposal: {Name: "proposal.toml.tmpl", Template: proposalTmpl},
		domain.SlashApply:    {Name: "apply.toml.tmpl", Template: applyTmpl},
	}
}

func TestSlashCommandsInitializer_Init(t *testing.T) {
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &domain.Config{SpectrDir: "spectr"}

	// Create parent directory
	_ = projectFs.MkdirAll(".claude/commands/spectr", 0o755)

	commands := newTestSlashTemplates(t)
	init := NewSlashCommandsInitializer(".claude/commands/spectr", commands)

	result, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	if len(result.CreatedFiles) != 2 {
		t.Errorf("Init() created %d files, want 2", len(result.CreatedFiles))
	}

	// Check proposal.md
	proposalPath := ".claude/commands/spectr/proposal.md"
	content, err := afero.ReadFile(projectFs, proposalPath)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", proposalPath, err)
	}
	if !strings.Contains(string(content), "# Proposal") {
		t.Error("proposal.md should contain Proposal header")
	}
	if !strings.Contains(string(content), "BaseDir: spectr") {
		t.Error("proposal.md should have rendered BaseDir")
	}

	// Check apply.md
	applyPath := ".claude/commands/spectr/apply.md"
	content, err = afero.ReadFile(projectFs, applyPath)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", applyPath, err)
	}
	if !strings.Contains(string(content), "# Apply") {
		t.Error("apply.md should contain Apply header")
	}
	if !strings.Contains(string(content), "SpecsDir: spectr/specs") {
		t.Error("apply.md should have rendered SpecsDir")
	}
}

func TestSlashCommandsInitializer_Idempotent(t *testing.T) {
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &domain.Config{SpectrDir: "spectr"}

	_ = projectFs.MkdirAll(".claude/commands/spectr", 0o755)

	commands := newTestSlashTemplates(t)
	init := NewSlashCommandsInitializer(".claude/commands/spectr", commands)

	// Run twice
	_, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err != nil {
		t.Fatalf("First Init() error = %v", err)
	}

	_, err = init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err != nil {
		t.Fatalf("Second Init() error = %v", err)
	}

	// Should still have files
	proposalPath := ".claude/commands/spectr/proposal.md"
	exists, _ := afero.Exists(projectFs, proposalPath)
	if !exists {
		t.Error("proposal.md should exist after second Init")
	}
}

func TestSlashCommandsInitializer_IsSetup(t *testing.T) {
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &domain.Config{SpectrDir: "spectr"}

	commands := newTestSlashTemplates(t)
	init := NewSlashCommandsInitializer(".claude/commands/spectr", commands)

	// Initially not setup
	if init.IsSetup(projectFs, homeFs, cfg) {
		t.Error("IsSetup() should return false when files don't exist")
	}

	// Create the files
	_ = projectFs.MkdirAll(".claude/commands/spectr", 0o755)
	_ = afero.WriteFile(projectFs, ".claude/commands/spectr/proposal.md", []byte("test"), 0o644)
	_ = afero.WriteFile(projectFs, ".claude/commands/spectr/apply.md", []byte("test"), 0o644)

	// Now should be setup
	if !init.IsSetup(projectFs, homeFs, cfg) {
		t.Error("IsSetup() should return true when all files exist")
	}
}

func TestSlashCommandsInitializer_dedupeKey(t *testing.T) {
	commands := newTestSlashTemplates(t)
	init, ok := NewSlashCommandsInitializer(".claude/commands/spectr/", commands).(*SlashCommandsInitializer)
	if !ok {
		t.Fatal("NewSlashCommandsInitializer did not return *SlashCommandsInitializer")
	}

	want := "SlashCommandsInitializer:" + filepath.Clean(".claude/commands/spectr/")
	if got := init.DedupeKey(); got != want {
		t.Errorf("dedupeKey() = %q, want %q", got, want)
	}
}

func TestHomeSlashCommandsInitializer_Init(t *testing.T) {
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &domain.Config{SpectrDir: "spectr"}

	// Create parent directory in HOME fs
	_ = homeFs.MkdirAll(".codex/prompts", 0o755)

	commands := newTestSlashTemplates(t)
	init := NewHomeSlashCommandsInitializer(".codex/prompts", commands)

	result, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	if len(result.CreatedFiles) != 2 {
		t.Errorf("Init() created %d files, want 2", len(result.CreatedFiles))
	}

	// Check files exist in homeFs, not projectFs
	proposalPath := ".codex/prompts/proposal.md"
	exists, _ := afero.Exists(homeFs, proposalPath)
	if !exists {
		t.Error("proposal.md should exist in homeFs")
	}

	exists, _ = afero.Exists(projectFs, proposalPath)
	if exists {
		t.Error("proposal.md should NOT exist in projectFs")
	}
}

func TestHomeSlashCommandsInitializer_IsSetup(t *testing.T) {
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &domain.Config{SpectrDir: "spectr"}

	commands := newTestSlashTemplates(t)
	init := NewHomeSlashCommandsInitializer(".codex/prompts", commands)

	// Create files in projectFs (wrong fs)
	_ = projectFs.MkdirAll(".codex/prompts", 0o755)
	_ = afero.WriteFile(projectFs, ".codex/prompts/proposal.md", []byte("test"), 0o644)
	_ = afero.WriteFile(projectFs, ".codex/prompts/apply.md", []byte("test"), 0o644)

	// Should still return false because files are in wrong fs
	if init.IsSetup(projectFs, homeFs, cfg) {
		t.Error("IsSetup() should return false when files exist in wrong fs")
	}

	// Create files in homeFs (correct fs)
	_ = homeFs.MkdirAll(".codex/prompts", 0o755)
	_ = afero.WriteFile(homeFs, ".codex/prompts/proposal.md", []byte("test"), 0o644)
	_ = afero.WriteFile(homeFs, ".codex/prompts/apply.md", []byte("test"), 0o644)

	// Now should be setup
	if !init.IsSetup(projectFs, homeFs, cfg) {
		t.Error("IsSetup() should return true when files exist in homeFs")
	}
}

func TestHomeSlashCommandsInitializer_dedupeKey(t *testing.T) {
	commands := newTestSlashTemplates(t)
	init, ok := NewHomeSlashCommandsInitializer(".codex/prompts", commands).(*HomeSlashCommandsInitializer)
	if !ok {
		t.Fatal("NewHomeSlashCommandsInitializer did not return *HomeSlashCommandsInitializer")
	}

	want := "HomeSlashCommandsInitializer:.codex/prompts"
	if got := init.DedupeKey(); got != want {
		t.Errorf("dedupeKey() = %q, want %q", got, want)
	}
}

func TestPrefixedSlashCommandsInitializer_Init(t *testing.T) {
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &domain.Config{SpectrDir: "spectr"}

	// Create parent directory
	_ = projectFs.MkdirAll(".agent/workflows", 0o755)

	commands := newTestSlashTemplates(t)
	init := NewPrefixedSlashCommandsInitializer(".agent/workflows", "spectr-", commands)

	result, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	if len(result.CreatedFiles) != 2 {
		t.Errorf("Init() created %d files, want 2", len(result.CreatedFiles))
	}

	// Check prefixed filenames
	proposalPath := ".agent/workflows/spectr-proposal.md"
	exists, _ := afero.Exists(projectFs, proposalPath)
	if !exists {
		t.Errorf("Expected %s to exist", proposalPath)
	}

	applyPath := ".agent/workflows/spectr-apply.md"
	exists, _ = afero.Exists(projectFs, applyPath)
	if !exists {
		t.Errorf("Expected %s to exist", applyPath)
	}
}

func TestPrefixedSlashCommandsInitializer_IsSetup(t *testing.T) {
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &domain.Config{SpectrDir: "spectr"}

	commands := newTestSlashTemplates(t)
	init := NewPrefixedSlashCommandsInitializer(".agent/workflows", "spectr-", commands)

	// Initially not setup
	if init.IsSetup(projectFs, homeFs, cfg) {
		t.Error("IsSetup() should return false when files don't exist")
	}

	// Create prefixed files
	_ = projectFs.MkdirAll(".agent/workflows", 0o755)
	_ = afero.WriteFile(projectFs, ".agent/workflows/spectr-proposal.md", []byte("test"), 0o644)
	_ = afero.WriteFile(projectFs, ".agent/workflows/spectr-apply.md", []byte("test"), 0o644)

	// Now should be setup
	if !init.IsSetup(projectFs, homeFs, cfg) {
		t.Error("IsSetup() should return true when all prefixed files exist")
	}
}

func TestPrefixedSlashCommandsInitializer_dedupeKey(t *testing.T) {
	commands := newTestSlashTemplates(t)
	init, ok := NewPrefixedSlashCommandsInitializer(".agent/workflows", "spectr-", commands).(*PrefixedSlashCommandsInitializer)
	if !ok {
		t.Fatal(
			"NewPrefixedSlashCommandsInitializer did not return *PrefixedSlashCommandsInitializer",
		)
	}

	want := "PrefixedSlashCommandsInitializer:.agent/workflows:spectr-"
	if got := init.DedupeKey(); got != want {
		t.Errorf("dedupeKey() = %q, want %q", got, want)
	}
}

func TestHomePrefixedSlashCommandsInitializer_Init(t *testing.T) {
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &domain.Config{SpectrDir: "spectr"}

	// Create parent directory in HOME fs
	_ = homeFs.MkdirAll(".codex/prompts", 0o755)

	commands := newTestSlashTemplates(t)
	init := NewHomePrefixedSlashCommandsInitializer(".codex/prompts", "spectr-", commands)

	result, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	if len(result.CreatedFiles) != 2 {
		t.Errorf("Init() created %d files, want 2", len(result.CreatedFiles))
	}

	// Check files exist in homeFs with prefix
	proposalPath := ".codex/prompts/spectr-proposal.md"
	exists, _ := afero.Exists(homeFs, proposalPath)
	if !exists {
		t.Errorf("Expected %s to exist in homeFs", proposalPath)
	}

	// Should NOT exist in projectFs
	exists, _ = afero.Exists(projectFs, proposalPath)
	if exists {
		t.Error("File should NOT exist in projectFs")
	}
}

func TestHomePrefixedSlashCommandsInitializer_IsSetup(t *testing.T) {
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &domain.Config{SpectrDir: "spectr"}

	commands := newTestSlashTemplates(t)
	init := NewHomePrefixedSlashCommandsInitializer(".codex/prompts", "spectr-", commands)

	// Initially not setup
	if init.IsSetup(projectFs, homeFs, cfg) {
		t.Error("IsSetup() should return false when files don't exist")
	}

	// Create prefixed files in homeFs
	_ = homeFs.MkdirAll(".codex/prompts", 0o755)
	_ = afero.WriteFile(homeFs, ".codex/prompts/spectr-proposal.md", []byte("test"), 0o644)
	_ = afero.WriteFile(homeFs, ".codex/prompts/spectr-apply.md", []byte("test"), 0o644)

	// Now should be setup
	if !init.IsSetup(projectFs, homeFs, cfg) {
		t.Error("IsSetup() should return true when all prefixed files exist in homeFs")
	}
}

func TestHomePrefixedSlashCommandsInitializer_dedupeKey(t *testing.T) {
	commands := newTestSlashTemplates(t)
	init, ok := NewHomePrefixedSlashCommandsInitializer(".codex/prompts", "spectr-", commands).(*HomePrefixedSlashCommandsInitializer)
	if !ok {
		t.Fatal(
			"NewHomePrefixedSlashCommandsInitializer did not return *HomePrefixedSlashCommandsInitializer",
		)
	}

	want := "HomePrefixedSlashCommandsInitializer:.codex/prompts:spectr-"
	if got := init.DedupeKey(); got != want {
		t.Errorf("dedupeKey() = %q, want %q", got, want)
	}
}

func TestTOMLSlashCommandsInitializer_Init(t *testing.T) {
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &domain.Config{SpectrDir: "spectr"}

	// Create parent directory
	_ = projectFs.MkdirAll(".gemini/commands/spectr", 0o755)

	commands := newTestTOMLTemplates(t)
	init := NewTOMLSlashCommandsInitializer(".gemini/commands/spectr", commands)

	result, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	if len(result.CreatedFiles) != 2 {
		t.Errorf("Init() created %d files, want 2", len(result.CreatedFiles))
	}

	// Check proposal.toml (NOT proposal.md)
	proposalPath := ".gemini/commands/spectr/proposal.toml"
	content, err := afero.ReadFile(projectFs, proposalPath)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", proposalPath, err)
	}
	if !strings.Contains(string(content), `description = "Proposal"`) {
		t.Error("proposal.toml should contain TOML description")
	}
	if !strings.Contains(string(content), "BaseDir: spectr") {
		t.Error("proposal.toml should have rendered BaseDir in prompt")
	}

	// Check apply.toml
	applyPath := ".gemini/commands/spectr/apply.toml"
	content, err = afero.ReadFile(projectFs, applyPath)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", applyPath, err)
	}
	if !strings.Contains(string(content), `description = "Apply"`) {
		t.Error("apply.toml should contain TOML description")
	}
}

func TestTOMLSlashCommandsInitializer_IsSetup(t *testing.T) {
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &domain.Config{SpectrDir: "spectr"}

	commands := newTestTOMLTemplates(t)
	init := NewTOMLSlashCommandsInitializer(".gemini/commands/spectr", commands)

	// Initially not setup
	if init.IsSetup(projectFs, homeFs, cfg) {
		t.Error("IsSetup() should return false when files don't exist")
	}

	// Create .md files (wrong extension)
	_ = projectFs.MkdirAll(".gemini/commands/spectr", 0o755)
	_ = afero.WriteFile(projectFs, ".gemini/commands/spectr/proposal.md", []byte("test"), 0o644)
	_ = afero.WriteFile(projectFs, ".gemini/commands/spectr/apply.md", []byte("test"), 0o644)

	// Should still return false because extension is wrong
	if init.IsSetup(projectFs, homeFs, cfg) {
		t.Error("IsSetup() should return false when files have wrong extension")
	}

	// Create .toml files (correct extension)
	_ = afero.WriteFile(projectFs, ".gemini/commands/spectr/proposal.toml", []byte("test"), 0o644)
	_ = afero.WriteFile(projectFs, ".gemini/commands/spectr/apply.toml", []byte("test"), 0o644)

	// Now should be setup
	if !init.IsSetup(projectFs, homeFs, cfg) {
		t.Error("IsSetup() should return true when all .toml files exist")
	}
}

func TestTOMLSlashCommandsInitializer_dedupeKey(t *testing.T) {
	commands := newTestTOMLTemplates(t)
	init, ok := NewTOMLSlashCommandsInitializer(".gemini/commands/spectr", commands).(*TOMLSlashCommandsInitializer)
	if !ok {
		t.Fatal("NewTOMLSlashCommandsInitializer did not return *TOMLSlashCommandsInitializer")
	}

	want := "TOMLSlashCommandsInitializer:.gemini/commands/spectr"
	if got := init.DedupeKey(); got != want {
		t.Errorf("dedupeKey() = %q, want %q", got, want)
	}
}

func TestSlashCommandsInitializer_UsesProjectFs(t *testing.T) {
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &domain.Config{SpectrDir: "spectr"}

	_ = projectFs.MkdirAll("testdir", 0o755)

	commands := newTestSlashTemplates(t)
	init := NewSlashCommandsInitializer("testdir", commands)

	_, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Should exist in projectFs
	exists, _ := afero.Exists(projectFs, "testdir/proposal.md")
	if !exists {
		t.Error("File should exist in projectFs")
	}

	// Should NOT exist in homeFs
	exists, _ = afero.Exists(homeFs, "testdir/proposal.md")
	if exists {
		t.Error("File should NOT exist in homeFs")
	}
}

func TestHomeSlashCommandsInitializer_UsesHomeFs(t *testing.T) {
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &domain.Config{SpectrDir: "spectr"}

	_ = homeFs.MkdirAll("testdir", 0o755)

	commands := newTestSlashTemplates(t)
	init := NewHomeSlashCommandsInitializer("testdir", commands)

	_, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Should exist in homeFs
	exists, _ := afero.Exists(homeFs, "testdir/proposal.md")
	if !exists {
		t.Error("File should exist in homeFs")
	}

	// Should NOT exist in projectFs
	exists, _ = afero.Exists(projectFs, "testdir/proposal.md")
	if exists {
		t.Error("File should NOT exist in projectFs")
	}
}
