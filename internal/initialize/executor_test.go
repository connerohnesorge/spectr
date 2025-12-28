package initialize

import (
	"context"
	"sort"
	"strings"
	"testing"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/connerohnesorge/spectr/internal/initialize/providers"
	"github.com/connerohnesorge/spectr/internal/initialize/providers/initializers"
	"github.com/spf13/afero"
)

// TestSortInitializersByType tests that initializers are sorted by type priority.
func TestSortInitializersByType(t *testing.T) {
	// Create mock initializers in wrong order
	mockInits := []domain.Initializer{
		initializers.NewSlashCommandsInitializer(".claude/commands/spectr", nil),
		initializers.NewConfigFileInitializer("CLAUDE.md", domain.TemplateRef{}),
		initializers.NewDirectoryInitializer(".claude/commands/spectr"),
	}

	sortInitializersByType(mockInits)

	// Verify order: Directory (1), ConfigFile (2), SlashCommands (3)
	expectedOrder := []string{
		"*initializers.DirectoryInitializer",
		"*initializers.ConfigFileInitializer",
		"*initializers.SlashCommandsInitializer",
	}

	for i, init := range mockInits {
		typeName := getTypeName(init)
		if typeName != expectedOrder[i] {
			t.Errorf("Position %d: got %s, want %s", i, typeName, expectedOrder[i])
		}
	}
}

// TestSortInitializersByType_AllTypes tests sorting with all initializer types.
func TestSortInitializersByType_AllTypes(t *testing.T) {
	mockInits := []domain.Initializer{
		initializers.NewTOMLSlashCommandsInitializer(".gemini/commands/spectr", nil),
		initializers.NewConfigFileInitializer("CLAUDE.md", domain.TemplateRef{}),
		initializers.NewHomeDirectoryInitializer(".codex/prompts"),
		initializers.NewPrefixedSlashCommandsInitializer(".agent/workflows", "spectr-", nil),
		initializers.NewDirectoryInitializer(".claude/commands/spectr"),
		initializers.NewHomePrefixedSlashCommandsInitializer(".codex/prompts", "spectr-", nil),
	}

	sortInitializersByType(mockInits)

	// Verify directories come first (priority 1)
	for i := range 2 {
		priority := initializerPriority(mockInits[i])
		if priority != 1 {
			t.Errorf(
				"Position %d should have priority 1, got %d (type: %T)",
				i,
				priority,
				mockInits[i],
			)
		}
	}

	// ConfigFile should be at position 2 (priority 2)
	if initializerPriority(mockInits[2]) != 2 {
		t.Errorf(
			"Position 2 should have priority 2, got %d (type: %T)",
			initializerPriority(mockInits[2]),
			mockInits[2],
		)
	}

	// SlashCommands should be at positions 3-5 (priority 3)
	for i := 3; i < 6; i++ {
		priority := initializerPriority(mockInits[i])
		if priority != 3 {
			t.Errorf(
				"Position %d should have priority 3, got %d (type: %T)",
				i,
				priority,
				mockInits[i],
			)
		}
	}
}

// TestDedupeInitializers tests deduplication of initializers.
func TestDedupeInitializers(t *testing.T) {
	// Create duplicate initializers
	inits := []domain.Initializer{
		initializers.NewDirectoryInitializer(".claude/commands/spectr"),
		initializers.NewDirectoryInitializer(".claude/commands/spectr"), // Duplicate
		initializers.NewDirectoryInitializer(".gemini/commands/spectr"),
		initializers.NewConfigFileInitializer("CLAUDE.md", domain.TemplateRef{}),
		initializers.NewConfigFileInitializer("CLAUDE.md", domain.TemplateRef{}), // Duplicate
	}

	result := dedupeInitializers(inits)

	// Should have 3 unique initializers
	if len(result) != 3 {
		t.Errorf("Expected 3 unique initializers, got %d", len(result))
	}
}

// TestDedupeInitializers_PreservesFirst tests that first occurrence is kept.
func TestDedupeInitializers_PreservesFirst(t *testing.T) {
	// The first DirectoryInitializer should be kept
	inits := []domain.Initializer{
		initializers.NewDirectoryInitializer(".claude/commands/spectr"),
		initializers.NewDirectoryInitializer(".claude/commands/spectr"),
	}

	result := dedupeInitializers(inits)

	if len(result) != 1 {
		t.Fatalf("Expected 1 initializer, got %d", len(result))
	}

	// Verify it's the first one
	d, ok := result[0].(*initializers.DirectoryInitializer)
	if !ok {
		t.Fatal("Expected DirectoryInitializer")
	}
	expectedKey := "DirectoryInitializer:.claude/commands/spectr"
	if d.DedupeKey() != expectedKey {
		t.Errorf("Kept wrong initializer: got key %q", d.DedupeKey())
	}
}

// TestTemplateContextFromConfig tests that TemplateContext is derived correctly.
func TestTemplateContextFromConfig(t *testing.T) {
	cfg := &domain.Config{SpectrDir: "spectr"}

	ctx := templateContextFromConfig(cfg)

	if ctx.BaseDir != "spectr" {
		t.Errorf("BaseDir = %q, want %q", ctx.BaseDir, "spectr")
	}
	if ctx.SpecsDir != "spectr/specs" {
		t.Errorf("SpecsDir = %q, want %q", ctx.SpecsDir, "spectr/specs")
	}
	if ctx.ChangesDir != "spectr/changes" {
		t.Errorf("ChangesDir = %q, want %q", ctx.ChangesDir, "spectr/changes")
	}
	if ctx.ProjectFile != "spectr/project.md" {
		t.Errorf("ProjectFile = %q, want %q", ctx.ProjectFile, "spectr/project.md")
	}
	if ctx.AgentsFile != "spectr/AGENTS.md" {
		t.Errorf("AgentsFile = %q, want %q", ctx.AgentsFile, "spectr/AGENTS.md")
	}
}

// TestTemplateContextFromConfig_CustomDir tests with custom spectr directory.
func TestTemplateContextFromConfig_CustomDir(t *testing.T) {
	cfg := &domain.Config{SpectrDir: "custom-dir"}

	ctx := templateContextFromConfig(cfg)

	if ctx.BaseDir != "custom-dir" {
		t.Errorf("BaseDir = %q, want %q", ctx.BaseDir, "custom-dir")
	}
	if ctx.SpecsDir != "custom-dir/specs" {
		t.Errorf("SpecsDir = %q, want %q", ctx.SpecsDir, "custom-dir/specs")
	}
}

// TestInitializerPriority tests the priority function for all initializer types.
func TestInitializerPriority(t *testing.T) {
	tests := []struct {
		name     string
		init     domain.Initializer
		expected int
	}{
		{"DirectoryInitializer", initializers.NewDirectoryInitializer("test"), 1},
		{"HomeDirectoryInitializer", initializers.NewHomeDirectoryInitializer("test"), 1},
		{
			"ConfigFileInitializer",
			initializers.NewConfigFileInitializer("test", domain.TemplateRef{}),
			2,
		},
		{"SlashCommandsInitializer", initializers.NewSlashCommandsInitializer("test", nil), 3},
		{
			"HomeSlashCommandsInitializer",
			initializers.NewHomeSlashCommandsInitializer("test", nil),
			3,
		},
		{
			"PrefixedSlashCommandsInitializer",
			initializers.NewPrefixedSlashCommandsInitializer("test", "prefix", nil),
			3,
		},
		{
			"HomePrefixedSlashCommandsInitializer",
			initializers.NewHomePrefixedSlashCommandsInitializer("test", "prefix", nil),
			3,
		},
		{
			"TOMLSlashCommandsInitializer",
			initializers.NewTOMLSlashCommandsInitializer("test", nil),
			3,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			priority := initializerPriority(tc.init)
			if priority != tc.expected {
				t.Errorf("%s priority = %d, want %d", tc.name, priority, tc.expected)
			}
		})
	}
}

// TestIntegration_DirectoryInitializer tests DirectoryInitializer with afero.MemMapFs.
func TestIntegration_DirectoryInitializer(t *testing.T) {
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &domain.Config{SpectrDir: "spectr"}

	init := initializers.NewDirectoryInitializer(".claude/commands/spectr")

	// Test Init creates directory
	result, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	if len(result.CreatedFiles) != 1 {
		t.Errorf("CreatedFiles count = %d, want 1", len(result.CreatedFiles))
	}
	if result.CreatedFiles[0] != ".claude/commands/spectr" {
		t.Errorf("CreatedFiles[0] = %q, want %q", result.CreatedFiles[0], ".claude/commands/spectr")
	}

	// Verify directory exists
	exists, err := afero.DirExists(projectFs, ".claude/commands/spectr")
	if err != nil {
		t.Fatalf("DirExists() error: %v", err)
	}
	if !exists {
		t.Error("Directory was not created")
	}

	// Test IsSetup returns true after creation
	if !init.IsSetup(projectFs, homeFs, cfg) {
		t.Error("IsSetup() should return true after Init()")
	}

	// Test Init is idempotent (second run doesn't report creation)
	result2, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err != nil {
		t.Fatalf("Second Init() error: %v", err)
	}
	if len(result2.CreatedFiles) != 0 {
		t.Errorf(
			"Second Init CreatedFiles count = %d, want 0 (idempotent)",
			len(result2.CreatedFiles),
		)
	}
}

// TestIntegration_HomeDirectoryInitializer tests HomeDirectoryInitializer with afero.MemMapFs.
func TestIntegration_HomeDirectoryInitializer(t *testing.T) {
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &domain.Config{SpectrDir: "spectr"}

	init := initializers.NewHomeDirectoryInitializer(".codex/prompts")

	// Test Init creates directory in home fs
	result, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
	if err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	if len(result.CreatedFiles) != 1 {
		t.Errorf("CreatedFiles count = %d, want 1", len(result.CreatedFiles))
	}

	// Verify directory exists in home fs, not project fs
	existsHome, _ := afero.DirExists(homeFs, ".codex/prompts")
	existsProject, _ := afero.DirExists(projectFs, ".codex/prompts")

	if !existsHome {
		t.Error("Directory was not created in home fs")
	}
	if existsProject {
		t.Error("Directory should not exist in project fs")
	}
}

// TestIntegration_ExecutionResultAccumulation tests that ExecutionResults are correctly merged.
func TestIntegration_ExecutionResultAccumulation(t *testing.T) {
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &domain.Config{SpectrDir: "spectr"}

	// Create multiple initializers
	inits := []domain.Initializer{
		initializers.NewDirectoryInitializer(".claude/commands/spectr"),
		initializers.NewDirectoryInitializer(".gemini/commands/spectr"),
		initializers.NewHomeDirectoryInitializer(".codex/prompts"),
	}

	// Simulate what executor does: run all and accumulate results
	var allCreated []string
	for _, init := range inits {
		result, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
		if err != nil {
			t.Fatalf("Init() error: %v", err)
		}
		allCreated = append(allCreated, result.CreatedFiles...)
	}

	// Verify accumulation
	if len(allCreated) != 3 {
		t.Errorf("Total CreatedFiles = %d, want 3", len(allCreated))
	}

	expectedCreated := []string{
		".claude/commands/spectr",
		".gemini/commands/spectr",
		".codex/prompts",
	}

	// Sort for comparison
	sort.Strings(allCreated)
	sort.Strings(expectedCreated)

	for i, expected := range expectedCreated {
		if i >= len(allCreated) || allCreated[i] != expected {
			t.Errorf("CreatedFiles mismatch at %d: got %v, want %v", i, allCreated, expectedCreated)

			break
		}
	}
}

// TestIntegration_FailFastOnError tests that execution stops on first error.
func TestIntegration_FailFastOnError(t *testing.T) {
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &domain.Config{SpectrDir: "spectr"}

	// Create an initializer that will fail (requires template but none provided)
	// SlashCommandsInitializer with nil commands map will work, but ConfigFileInitializer
	// with nil template will fail during Init
	nilTemplateInit := initializers.NewConfigFileInitializer("TEST.md", domain.TemplateRef{
		Name:     "nonexistent.tmpl",
		Template: nil, // This will cause a nil pointer panic which we need to recover from
	})

	// Create a list with failing initializer in the middle
	inits := []domain.Initializer{
		initializers.NewDirectoryInitializer(".first/dir"),
		nilTemplateInit,
		initializers.NewDirectoryInitializer(".second/dir"),
	}

	var allCreated []string
	var sawError bool

	for _, init := range inits {
		// Recover from panic in case nil template causes it
		func() {
			defer func() {
				if r := recover(); r != nil {
					sawError = true
				}
			}()
			result, err := init.Init(context.Background(), projectFs, homeFs, cfg, nil)
			if err != nil {
				sawError = true

				return
			}
			allCreated = append(allCreated, result.CreatedFiles...)
		}()

		if sawError {
			break // Fail-fast behavior
		}
	}

	// First directory should have been created
	if len(allCreated) != 1 {
		t.Errorf("Should have created 1 directory before failing, got %d", len(allCreated))
	}

	// Second directory should NOT exist (stopped before)
	exists, _ := afero.DirExists(projectFs, ".second/dir")
	if exists {
		t.Error(".second/dir should not exist due to fail-fast")
	}

	// Verify error was detected
	if !sawError {
		t.Error("Expected error during execution")
	}
}

// TestIntegration_ProviderInitializerCollection tests collecting initializers from providers.
func TestIntegration_ProviderInitializerCollection(t *testing.T) {
	providers.ResetRegistry()
	defer func() {
		providers.ResetRegistry()
		_ = providers.RegisterAllProviders()
	}()

	if err := providers.RegisterAllProviders(); err != nil {
		t.Fatalf("RegisterAllProviders() failed: %v", err)
	}

	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() failed: %v", err)
	}

	// Get Claude provider
	reg, found := providers.Get("claude-code")
	if !found {
		t.Fatal("claude-code provider not found")
	}

	// Get initializers
	inits := reg.Provider.Initializers(context.Background(), tm)

	// Verify count (DirectoryInitializer + ConfigFileInitializer + SlashCommandsInitializer)
	if len(inits) != 3 {
		t.Errorf("Claude provider returned %d initializers, want 3", len(inits))
	}

	// Verify types
	var hasDir, hasConfig, hasSlash bool
	for _, init := range inits {
		switch init.(type) {
		case *initializers.DirectoryInitializer:
			hasDir = true
		case *initializers.ConfigFileInitializer:
			hasConfig = true
		case *initializers.SlashCommandsInitializer:
			hasSlash = true
		}
	}

	if !hasDir {
		t.Error("Missing DirectoryInitializer")
	}
	if !hasConfig {
		t.Error("Missing ConfigFileInitializer")
	}
	if !hasSlash {
		t.Error("Missing SlashCommandsInitializer")
	}
}

// TestIntegration_DeduplicationAcrossProviders tests deduplication when providers share resources.
func TestIntegration_DeduplicationAcrossProviders(t *testing.T) {
	providers.ResetRegistry()
	defer func() {
		providers.ResetRegistry()
		_ = providers.RegisterAllProviders()
	}()

	if err := providers.RegisterAllProviders(); err != nil {
		t.Fatalf("RegisterAllProviders() failed: %v", err)
	}

	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() failed: %v", err)
	}

	// Antigravity and Codex both use AGENTS.md
	antigravityReg, _ := providers.Get("antigravity")
	codexReg, _ := providers.Get("codex")

	// Collect initializers from both
	var allInits []domain.Initializer
	allInits = append(allInits, antigravityReg.Provider.Initializers(context.Background(), tm)...)
	allInits = append(allInits, codexReg.Provider.Initializers(context.Background(), tm)...)

	// Sort by type (as executor does)
	sortInitializersByType(allInits)

	// Dedupe
	deduped := dedupeInitializers(allInits)

	// Count ConfigFileInitializers for AGENTS.md
	agentsCount := 0
	for _, init := range deduped {
		cfg, ok := init.(*initializers.ConfigFileInitializer)
		if !ok {
			continue
		}
		if strings.Contains(cfg.DedupeKey(), "AGENTS.md") {
			agentsCount++
		}
	}

	// Should only have one AGENTS.md initializer
	if agentsCount != 1 {
		t.Errorf("Expected 1 AGENTS.md initializer after dedup, got %d", agentsCount)
	}
}

// TestIntegration_FullInitializationFlow tests the complete initialization flow.
func TestIntegration_FullInitializationFlow(t *testing.T) {
	providers.ResetRegistry()
	defer func() {
		providers.ResetRegistry()
		_ = providers.RegisterAllProviders()
	}()

	if err := providers.RegisterAllProviders(); err != nil {
		t.Fatalf("RegisterAllProviders() failed: %v", err)
	}

	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &domain.Config{SpectrDir: "spectr"}

	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("NewTemplateManager() failed: %v", err)
	}

	// Select a few providers to test
	selectedIDs := []string{"claude-code", "cursor", "gemini"}

	// Collect initializers from selected providers (as executor does)
	allRegistrations := providers.RegisteredProviders()
	var allInits []domain.Initializer

	for _, reg := range allRegistrations {
		isSelected := false
		for _, id := range selectedIDs {
			if id == reg.ID {
				isSelected = true

				break
			}
		}
		if !isSelected {
			continue
		}

		inits := reg.Provider.Initializers(context.Background(), tm)
		allInits = append(allInits, inits...)
	}

	// Sort and dedupe
	sortInitializersByType(allInits)
	deduped := dedupeInitializers(allInits)

	// Execute all
	var allCreated []string
	for _, init := range deduped {
		result, err := init.Init(context.Background(), projectFs, homeFs, cfg, tm)
		if err != nil {
			t.Fatalf("Init() error: %v", err)
		}
		allCreated = append(allCreated, result.CreatedFiles...)
	}

	// Verify directories were created
	expectedDirs := []string{
		".claude/commands/spectr",
		".cursorrules/commands/spectr",
		".gemini/commands/spectr",
	}

	for _, dir := range expectedDirs {
		exists, err := afero.DirExists(projectFs, dir)
		if err != nil {
			t.Errorf("Error checking %s: %v", dir, err)

			continue
		}
		if !exists {
			t.Errorf("Directory %s was not created", dir)
		}
	}

	// Verify config file was created (CLAUDE.md)
	exists, err := afero.Exists(projectFs, "CLAUDE.md")
	if err != nil {
		t.Errorf("Error checking CLAUDE.md: %v", err)
	}
	if !exists {
		t.Error("CLAUDE.md was not created")
	}

	// Verify slash commands were created
	expectedSlashCommands := []string{
		".claude/commands/spectr/proposal.md",
		".claude/commands/spectr/apply.md",
		".cursorrules/commands/spectr/proposal.md",
		".cursorrules/commands/spectr/apply.md",
		".gemini/commands/spectr/proposal.toml",
		".gemini/commands/spectr/apply.toml",
	}

	for _, path := range expectedSlashCommands {
		exists, err := afero.Exists(projectFs, path)
		if err != nil {
			t.Errorf("Error checking %s: %v", path, err)

			continue
		}
		if !exists {
			t.Errorf("Slash command %s was not created", path)
		}
	}

	// Verify results contain created files
	if len(allCreated) == 0 {
		t.Error("No files reported as created")
	}
}

// Helper function to get type name.
func getTypeName(init domain.Initializer) string {
	switch init.(type) {
	case *initializers.DirectoryInitializer:
		return "*initializers.DirectoryInitializer"
	case *initializers.HomeDirectoryInitializer:
		return "*initializers.HomeDirectoryInitializer"
	case *initializers.ConfigFileInitializer:
		return "*initializers.ConfigFileInitializer"
	case *initializers.SlashCommandsInitializer:
		return "*initializers.SlashCommandsInitializer"
	case *initializers.HomeSlashCommandsInitializer:
		return "*initializers.HomeSlashCommandsInitializer"
	case *initializers.PrefixedSlashCommandsInitializer:
		return "*initializers.PrefixedSlashCommandsInitializer"
	case *initializers.HomePrefixedSlashCommandsInitializer:
		return "*initializers.HomePrefixedSlashCommandsInitializer"
	case *initializers.TOMLSlashCommandsInitializer:
		return "*initializers.TOMLSlashCommandsInitializer"
	default:
		return "unknown"
	}
}
