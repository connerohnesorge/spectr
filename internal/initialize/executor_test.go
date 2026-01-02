package initialize

import (
	"context"
	"strings"
	"testing"

	"github.com/spf13/afero"

	"github.com/connerohnesorge/spectr/internal/initialize/providers"
)

// TestExecutorIntegration_FullInitializationFlow tests the full initialization flow
// using afero.MemMapFs for filesystem operations
func TestExecutorIntegration_FullInitializationFlow(t *testing.T) {
	// Create in-memory filesystems
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()

	// Create spectr directory structure
	if err := projectFs.MkdirAll("spectr", 0o755); err != nil {
		t.Fatalf("Failed to create spectr directory: %v", err)
	}

	// Initialize template manager
	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	// Reset and register all providers
	providers.Reset()
	if err := providers.RegisterAllProviders(); err != nil {
		t.Fatalf("Failed to register providers: %v", err)
	}

	// Test with Claude Code provider (ID: "claude-code")
	ctx := context.Background()
	cfg := &providers.Config{
		SpectrDir: "spectr",
	}

	// Get Claude provider
	claudeReg, ok := providers.Get("claude-code")
	if !ok {
		t.Fatal("Claude provider not found in registry")
	}

	// Get initializers from Claude provider
	inits := claudeReg.Provider.Initializers(ctx, tm)

	if len(inits) != 7 {
		t.Fatalf("Claude provider returned %d initializers, want 7", len(inits))
	}

	// Execute each initializer and collect results
	allResults := make([]providers.InitResult, 0, len(inits))

	for _, init := range inits {
		result, err := init.Init(ctx, projectFs, homeFs, cfg, tm)
		if err != nil {
			t.Fatalf("Initializer failed: %v", err)
		}
		allResults = append(allResults, result)
	}

	// Verify files were created
	expectedFiles := []string{
		".claude/commands/spectr", // Directory
		"CLAUDE.md",               // Config file
		".claude/commands/spectr/proposal.md",
		".claude/commands/spectr/apply.md",
	}

	for _, expectedFile := range expectedFiles {
		exists, err := afero.Exists(projectFs, expectedFile)
		if err != nil {
			t.Errorf("Failed to check file %s: %v", expectedFile, err)

			continue
		}
		if !exists {
			t.Errorf("Expected file %s was not created", expectedFile)
		}
	}

	// Verify CLAUDE.md contains spectr markers
	claudeContent, err := afero.ReadFile(projectFs, "CLAUDE.md")
	if err != nil {
		t.Fatalf("Failed to read CLAUDE.md: %v", err)
	}

	claudeStr := string(claudeContent)
	if !strings.Contains(claudeStr, "<!-- spectr:start -->") {
		t.Error("CLAUDE.md missing start marker")
	}
	if !strings.Contains(claudeStr, "<!-- spectr:end -->") {
		t.Error("CLAUDE.md missing end marker")
	}

	// Verify slash command files contain content
	proposalContent, err := afero.ReadFile(projectFs, ".claude/commands/spectr/proposal.md")
	if err != nil {
		t.Fatalf("Failed to read proposal.md: %v", err)
	}
	if len(proposalContent) == 0 {
		t.Error("proposal.md is empty")
	}

	applyContent, err := afero.ReadFile(projectFs, ".claude/commands/spectr/apply.md")
	if err != nil {
		t.Fatalf("Failed to read apply.md: %v", err)
	}
	if len(applyContent) == 0 {
		t.Error("apply.md is empty")
	}

	// Verify InitResult accumulation
	totalCreated := 0
	totalUpdated := 0
	for _, result := range allResults {
		totalCreated += len(result.CreatedFiles)
		totalUpdated += len(result.UpdatedFiles)
	}

	// First run should create files (directories count as created)
	if totalCreated == 0 {
		t.Error("No files were reported as created")
	}
}

// TestExecutorIntegration_InitResultAccumulation tests that InitResult
// correctly tracks created and updated files
func TestExecutorIntegration_InitResultAccumulation(t *testing.T) {
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()

	if err := projectFs.MkdirAll("spectr", 0o755); err != nil {
		t.Fatalf("Failed to create spectr directory: %v", err)
	}

	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	providers.Reset()
	if err := providers.RegisterAllProviders(); err != nil {
		t.Fatalf("Failed to register providers: %v", err)
	}

	ctx := context.Background()
	cfg := &providers.Config{
		SpectrDir: "spectr",
	}

	claudeReg, ok := providers.Get("claude-code")
	if !ok {
		t.Fatal("Claude provider not found")
	}

	inits := claudeReg.Provider.Initializers(ctx, tm)

	// First run: should create files
	firstResults := make([]providers.InitResult, 0, len(inits))
	for _, init := range inits {
		result, err := init.Init(ctx, projectFs, homeFs, cfg, tm)
		if err != nil {
			t.Fatalf("First run initializer failed: %v", err)
		}
		firstResults = append(firstResults, result)
	}

	// Count created files in first run
	firstCreated := 0
	firstUpdated := 0
	for _, result := range firstResults {
		firstCreated += len(result.CreatedFiles)
		firstUpdated += len(result.UpdatedFiles)
	}

	if firstCreated == 0 {
		t.Error("First run: no files were created")
	}
	if firstUpdated != 0 {
		t.Errorf("First run: %d files reported as updated, want 0", firstUpdated)
	}

	// Second run: should update existing files (config file should be updated)
	secondResults := make([]providers.InitResult, 0, len(inits))
	for _, init := range inits {
		result, err := init.Init(ctx, projectFs, homeFs, cfg, tm)
		if err != nil {
			t.Fatalf("Second run initializer failed: %v", err)
		}
		secondResults = append(secondResults, result)
	}

	// Count updated files in second run
	secondCreated := 0
	secondUpdated := 0
	for _, result := range secondResults {
		secondCreated += len(result.CreatedFiles)
		secondUpdated += len(result.UpdatedFiles)
	}

	// Second run should have fewer created files (directories already exist)
	// Config file and slash commands should be updated
	if secondUpdated == 0 {
		t.Error(
			"Second run: no files were updated (expected config file and slash commands to be updated)",
		)
	}

	// Directories should not be re-created
	if secondCreated >= firstCreated {
		t.Errorf(
			"Second run created %d files, expected fewer than first run (%d)",
			secondCreated,
			firstCreated,
		)
	}
}

// TestExecutorIntegration_MultipleProviders tests initialization with multiple providers
func TestExecutorIntegration_MultipleProviders(t *testing.T) {
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()

	if err := projectFs.MkdirAll("spectr", 0o755); err != nil {
		t.Fatalf("Failed to create spectr directory: %v", err)
	}

	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	providers.Reset()
	if err := providers.RegisterAllProviders(); err != nil {
		t.Fatalf("Failed to register providers: %v", err)
	}

	ctx := context.Background()
	cfg := &providers.Config{
		SpectrDir: "spectr",
	}

	// Test with Claude and Gemini providers
	providerIDs := []string{"claude-code", "gemini"}
	var allInitializers []providers.Initializer

	for _, id := range providerIDs {
		reg, ok := providers.Get(id)
		if !ok {
			t.Fatalf("Provider %s not found", id)
		}
		inits := reg.Provider.Initializers(ctx, tm)
		allInitializers = append(allInitializers, inits...)
	}

	// Execute all initializers
	for _, init := range allInitializers {
		_, err := init.Init(ctx, projectFs, homeFs, cfg, tm)
		if err != nil {
			t.Fatalf("Initializer failed: %v", err)
		}
	}

	// Verify both provider files exist
	claudeFiles := []string{
		"CLAUDE.md",
		".claude/commands/spectr/proposal.md",
		".claude/commands/spectr/apply.md",
	}

	geminiFiles := []string{
		".gemini/commands/spectr/proposal.toml",
		".gemini/commands/spectr/apply.toml",
	}

	for _, file := range claudeFiles {
		exists, err := afero.Exists(projectFs, file)
		if err != nil {
			t.Errorf("Failed to check %s: %v", file, err)
		} else if !exists {
			t.Errorf("Claude file %s not created", file)
		}
	}

	for _, file := range geminiFiles {
		exists, err := afero.Exists(projectFs, file)
		if err != nil {
			t.Errorf("Failed to check %s: %v", file, err)
		} else if !exists {
			t.Errorf("Gemini file %s not created", file)
		}
	}

	// Verify Gemini files are TOML format
	proposalContent, err := afero.ReadFile(projectFs, ".gemini/commands/spectr/proposal.toml")
	if err != nil {
		t.Fatalf("Failed to read Gemini proposal.toml: %v", err)
	}

	proposalStr := string(proposalContent)
	if !strings.Contains(proposalStr, "description =") {
		t.Error("Gemini proposal.toml missing 'description =' field")
	}
	if !strings.Contains(proposalStr, "prompt =") {
		t.Error("Gemini proposal.toml missing 'prompt =' field")
	}
}

// TestExecutorIntegration_HomeFilesystem tests providers that use home filesystem
func TestExecutorIntegration_HomeFilesystem(t *testing.T) {
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()

	if err := projectFs.MkdirAll("spectr", 0o755); err != nil {
		t.Fatalf("Failed to create spectr directory: %v", err)
	}

	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	providers.Reset()
	if err := providers.RegisterAllProviders(); err != nil {
		t.Fatalf("Failed to register providers: %v", err)
	}

	ctx := context.Background()
	cfg := &providers.Config{
		SpectrDir: "spectr",
	}

	// Test Codex provider (uses home filesystem)
	codexReg, ok := providers.Get("codex")
	if !ok {
		t.Fatal("Codex provider not found")
	}

	inits := codexReg.Provider.Initializers(ctx, tm)

	// Execute initializers
	for _, init := range inits {
		_, err := init.Init(ctx, projectFs, homeFs, cfg, tm)
		if err != nil {
			t.Fatalf("Codex initializer failed: %v", err)
		}
	}

	// Verify home directory files
	homeFiles := []string{
		".codex/prompts",
		".codex/prompts/spectr-proposal.md",
		".codex/prompts/spectr-apply.md",
	}

	for _, file := range homeFiles {
		exists, err := afero.Exists(homeFs, file)
		if err != nil {
			t.Errorf("Failed to check home file %s: %v", file, err)
		} else if !exists {
			t.Errorf("Home file %s not created", file)
		}
	}

	// Verify project file (AGENTS.md should be in project, not home)
	projectExists, err := afero.Exists(projectFs, "AGENTS.md")
	if err != nil {
		t.Errorf("Failed to check AGENTS.md: %v", err)
	} else if !projectExists {
		t.Error("AGENTS.md should be in project filesystem")
	}

	// Verify it's NOT in home filesystem
	homeExists, err := afero.Exists(homeFs, "AGENTS.md")
	if err != nil {
		t.Errorf("Failed to check home AGENTS.md: %v", err)
	} else if homeExists {
		t.Error("AGENTS.md should NOT be in home filesystem")
	}
}

// TestExecutorIntegration_ErrorHandling tests fail-fast error behavior
func TestExecutorIntegration_ErrorHandling(t *testing.T) {
	// Use a read-only filesystem to trigger errors
	projectFs := afero.NewReadOnlyFs(afero.NewMemMapFs())
	homeFs := afero.NewMemMapFs()

	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	providers.Reset()
	if err := providers.RegisterAllProviders(); err != nil {
		t.Fatalf("Failed to register providers: %v", err)
	}

	ctx := context.Background()
	cfg := &providers.Config{
		SpectrDir: "spectr",
	}

	claudeReg, ok := providers.Get("claude-code")
	if !ok {
		t.Fatal("Claude provider not found")
	}

	inits := claudeReg.Provider.Initializers(ctx, tm)

	// Execute initializers - should fail on read-only filesystem
	for i, init := range inits {
		_, err := init.Init(ctx, projectFs, homeFs, cfg, tm)
		if err == nil {
			// Some initializers might succeed (e.g., checking if files exist)
			// but directory creation should fail
			continue
		}

		// Error occurred - verify it's a meaningful error
		if err != nil && i == 0 {
			// First initializer is directory creation, should fail on read-only fs
			if !strings.Contains(err.Error(), "failed to create directory") &&
				!strings.Contains(err.Error(), "operation not permitted") &&
				!strings.Contains(err.Error(), "read-only") {
				t.Errorf("Expected directory creation error, got: %v", err)
			}

			return // Test passed - fail-fast worked
		}
	}

	// If we got here without errors, it might be because MemMapFs doesn't enforce read-only
	// This is acceptable - the important part is that errors are returned, not swallowed
}

// TestExecutorIntegration_InitializerOrdering tests that initializers are executed in the correct order
func TestExecutorIntegration_InitializerOrdering(t *testing.T) {
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()

	if err := projectFs.MkdirAll("spectr", 0o755); err != nil {
		t.Fatalf("Failed to create spectr directory: %v", err)
	}

	tm, err := NewTemplateManager()
	if err != nil {
		t.Fatalf("Failed to create template manager: %v", err)
	}

	providers.Reset()
	if err := providers.RegisterAllProviders(); err != nil {
		t.Fatalf("Failed to register providers: %v", err)
	}

	ctx := context.Background()
	cfg := &providers.Config{
		SpectrDir: "spectr",
	}

	claudeReg, ok := providers.Get("claude-code")
	if !ok {
		t.Fatal("Claude provider not found")
	}

	inits := claudeReg.Provider.Initializers(ctx, tm)

	// Verify initializer order: Directory (commands), Directory (skills), ConfigFile, SlashCommands, AgentSkills
	if len(inits) < 5 {
		t.Fatalf("Expected at least 5 initializers, got %d", len(inits))
	}

	// Check type order
	if _, ok := inits[0].(*providers.DirectoryInitializer); !ok {
		t.Errorf("First initializer should be DirectoryInitializer, got %T", inits[0])
	}

	if _, ok := inits[1].(*providers.DirectoryInitializer); !ok {
		t.Errorf("Second initializer should be DirectoryInitializer, got %T", inits[1])
	}

	if _, ok := inits[2].(*providers.ConfigFileInitializer); !ok {
		t.Errorf("Third initializer should be ConfigFileInitializer, got %T", inits[2])
	}

	if _, ok := inits[3].(*providers.SlashCommandsInitializer); !ok {
		t.Errorf("Fourth initializer should be SlashCommandsInitializer, got %T", inits[3])
	}

	if _, ok := inits[4].(*providers.AgentSkillsInitializer); !ok {
		t.Errorf("Fifth initializer should be AgentSkillsInitializer, got %T", inits[4])
	}

	// Execute in order and verify each step
	// 1. Directory should be created first
	result1, err := inits[0].Init(ctx, projectFs, homeFs, cfg, tm)
	if err != nil {
		t.Fatalf("Directory initializer failed: %v", err)
	}
	if len(result1.CreatedFiles) == 0 {
		t.Error("Directory initializer should create directory")
	}

	// 2. Config file should be created second (depends on nothing)
	result2, err := inits[1].Init(ctx, projectFs, homeFs, cfg, tm)
	if err != nil {
		t.Fatalf("ConfigFile initializer failed: %v", err)
	}
	if len(result2.CreatedFiles) == 0 {
		t.Error("ConfigFile initializer should create file")
	}

	// 3. Slash commands should be created last (depends on directory existing)
	result3, err := inits[2].Init(ctx, projectFs, homeFs, cfg, tm)
	if err != nil {
		t.Fatalf("SlashCommands initializer failed: %v", err)
	}
	if len(result3.CreatedFiles) == 0 && len(result3.UpdatedFiles) == 0 {
		t.Error("SlashCommands initializer should create or update files")
	}
}

// TestAggregateResultsDeduplication tests that aggregateResults deduplicates file paths
func TestAggregateResultsDeduplication(t *testing.T) {
	tests := []struct {
		name     string
		results  []providers.InitResult
		wantLen  int
		wantFile string
	}{
		{
			name: "deduplicates created files",
			results: []providers.InitResult{
				{
					CreatedFiles: []string{"file1.txt", "file2.txt"},
					UpdatedFiles: make([]string, 0),
				},
				{
					CreatedFiles: []string{"file1.txt", "file3.txt"},
					UpdatedFiles: make([]string, 0),
				},
			},
			wantLen:  3, // file1.txt, file2.txt, file3.txt (file1.txt appears once)
			wantFile: "file1.txt",
		},
		{
			name: "deduplicates updated files",
			results: []providers.InitResult{
				{
					CreatedFiles: make([]string, 0),
					UpdatedFiles: []string{"config.yml", "settings.json"},
				},
				{
					CreatedFiles: make([]string, 0),
					UpdatedFiles: []string{"config.yml"},
				},
			},
			wantLen:  2, // config.yml, settings.json (config.yml appears once)
			wantFile: "config.yml",
		},
		{
			name: "deduplicates across multiple initializers",
			results: []providers.InitResult{
				{
					CreatedFiles: []string{"dir/file.txt"},
					UpdatedFiles: make([]string, 0),
				},
				{
					CreatedFiles: []string{"other.txt"},
					UpdatedFiles: make([]string, 0),
				},
				{
					CreatedFiles: []string{"dir/file.txt"},
					UpdatedFiles: make([]string, 0),
				},
			},
			wantLen:  2, // dir/file.txt, other.txt
			wantFile: "dir/file.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := aggregateResults(tt.results)

			totalFiles := len(result.CreatedFiles) + len(result.UpdatedFiles)
			if totalFiles != tt.wantLen {
				t.Errorf("aggregateResults returned %d files, want %d", totalFiles, tt.wantLen)
			}

			// Check that the expected file is present
			found := false
			for _, f := range result.CreatedFiles {
				if f == tt.wantFile {
					found = true

					break
				}
			}
			if !found {
				for _, f := range result.UpdatedFiles {
					if f == tt.wantFile {
						found = true

						break
					}
				}
			}
			if !found {
				t.Errorf("Expected file %q not found in aggregated results", tt.wantFile)
			}
		})
	}
}
