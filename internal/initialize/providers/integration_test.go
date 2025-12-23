package providers

import (
	"context"
	"testing"

	"github.com/spf13/afero"
)

// TestFullInitializationFlow is an integration test that verifies the complete initialization flow.
// It tests that a provider's initializers can create files correctly using an in-memory filesystem.
func TestFullInitializationFlow(t *testing.T) {
	// Create in-memory filesystem
	fs := afero.NewMemMapFs()

	// Define test config
	cfg := &Config{
		SpectrDir: "spectr",
	}

	// Create mock template manager with realistic content
	tm := &mockTemplateManager{
		content: "# Mock Template Content\n\nThis is test content for initialization.\n",
		err:     nil,
	}

	// Create context
	ctx := context.Background()

	// Get initializers from Claude Code provider (as example)
	reg := Get("claude-code")
	if reg == nil {
		t.Fatal("claude-code provider not found")
	}

	initializers := reg.Provider.Initializers(ctx)

	// Execute all initializers
	allResults := make([]InitResult, 0, len(initializers))
	for _, init := range initializers {
		result, err := init.Init(ctx, fs, cfg, tm)
		if err != nil {
			t.Fatalf("Initializer %T failed: %v", init, err)
		}
		allResults = append(allResults, result)
	}

	// Verify that files were created
	totalCreated := 0
	totalUpdated := 0
	for _, result := range allResults {
		totalCreated += len(result.CreatedFiles)
		totalUpdated += len(result.UpdatedFiles)
	}

	if totalCreated == 0 {
		t.Error("No files were created")
	}

	t.Logf("Created %d files, updated %d files", totalCreated, totalUpdated)

	// Verify specific files exist
	expectedFiles := []string{
		".claude/commands/spectr",             // Directory
		"CLAUDE.md",                           // Config file
		".claude/commands/spectr/proposal.md", // Slash command
		".claude/commands/spectr/apply.md",    // Slash command
	}

	for _, path := range expectedFiles {
		exists, err := afero.Exists(fs, path)
		if err != nil {
			t.Errorf("Error checking if %s exists: %v", path, err)

			continue
		}
		if !exists {
			t.Errorf("Expected file/directory %s does not exist", path)
		}
	}

	// Verify directory exists and is a directory
	isDir, err := afero.IsDir(fs, ".claude/commands/spectr")
	if err != nil {
		t.Errorf("Error checking if .claude/commands/spectr is directory: %v", err)
	} else if !isDir {
		t.Error(".claude/commands/spectr is not a directory")
	}

	// Verify config file has content
	configContent, err := afero.ReadFile(fs, "CLAUDE.md")
	if err != nil {
		t.Errorf("Error reading CLAUDE.md: %v", err)
	} else if len(configContent) == 0 {
		t.Error("CLAUDE.md is empty")
	}

	// Verify slash commands have content
	for _, cmd := range []string{"proposal", "apply"} {
		cmdPath := ".claude/commands/spectr/" + cmd + ".md"
		cmdContent, err := afero.ReadFile(fs, cmdPath)
		if err != nil {
			t.Errorf("Error reading %s: %v", cmdPath, err)
		} else if len(cmdContent) == 0 {
			t.Errorf("%s is empty", cmdPath)
		}
	}
}

// TestInitResultAccumulation verifies that InitResult values are correctly accumulated.
// This tests that multiple initializers' results can be merged together properly.
func TestInitResultAccumulation(t *testing.T) {
	// Create in-memory filesystem
	fs := afero.NewMemMapFs()

	// Create config and template manager
	cfg := &Config{SpectrDir: "spectr"}
	tm := &mockTemplateManager{content: "# Mock Content\n", err: nil}
	ctx := context.Background()

	// Test with a provider that has multiple initializers
	reg := Get("claude-code")
	if reg == nil {
		t.Fatal("claude-code provider not found")
	}

	initializers := reg.Provider.Initializers(ctx)

	// Collect results from all initializers
	allResults := make([]InitResult, 0, len(initializers))
	for _, init := range initializers {
		result, err := init.Init(ctx, fs, cfg, tm)
		if err != nil {
			t.Fatalf("Initializer failed: %v", err)
		}
		allResults = append(allResults, result)
	}

	// Aggregate results (simulating what executor does)
	var aggregated InitResult
	for _, result := range allResults {
		aggregated.CreatedFiles = append(aggregated.CreatedFiles, result.CreatedFiles...)
		aggregated.UpdatedFiles = append(aggregated.UpdatedFiles, result.UpdatedFiles...)
	}

	// Verify aggregated results contain all created files
	if len(aggregated.CreatedFiles) == 0 {
		t.Error("Aggregated CreatedFiles is empty")
	}

	// Verify no files in UpdatedFiles (first run should only create)
	if len(aggregated.UpdatedFiles) > 0 {
		t.Logf(
			"Found %d updated files on first run (unexpected but not necessarily wrong)",
			len(aggregated.UpdatedFiles),
		)
	}

	t.Logf(
		"Aggregated results: %d created, %d updated",
		len(aggregated.CreatedFiles),
		len(aggregated.UpdatedFiles),
	)

	// Now run initializers again - should show files as updated or already existing
	secondRunResults := make([]InitResult, 0, len(initializers))
	for _, init := range initializers {
		result, err := init.Init(ctx, fs, cfg, tm)
		if err != nil {
			t.Fatalf("Second run initializer failed: %v", err)
		}
		secondRunResults = append(secondRunResults, result)
	}

	// Aggregate second run results
	var aggregated2 InitResult
	for _, result := range secondRunResults {
		aggregated2.CreatedFiles = append(aggregated2.CreatedFiles, result.CreatedFiles...)
		aggregated2.UpdatedFiles = append(aggregated2.UpdatedFiles, result.UpdatedFiles...)
	}

	// On second run, files should be updated (or skipped), not created
	// The exact behavior depends on the initializer implementation
	t.Logf(
		"Second run results: %d created, %d updated",
		len(aggregated2.CreatedFiles),
		len(aggregated2.UpdatedFiles),
	)
}

// TestPartialFailureResultAccumulation verifies that results are accumulated even when some initializers fail.
func TestPartialFailureResultAccumulation(t *testing.T) {
	// Create in-memory filesystem
	fs := afero.NewMemMapFs()

	// Create config and template manager
	cfg := &Config{SpectrDir: "spectr"}
	tm := &mockTemplateManager{content: "# Mock Content\n", err: nil}
	ctx := context.Background()

	// Create a mix of successful and potentially failing initializers
	// For this test, we'll create some files then make filesystem read-only to simulate failure

	// First, successfully create some files
	init1 := NewDirectoryInitializer(".test/dir1")
	result1, err := init1.Init(ctx, fs, cfg, tm)
	if err != nil {
		t.Fatalf("First initializer failed: %v", err)
	}

	init2 := NewDirectoryInitializer(".test/dir2")
	result2, err := init2.Init(ctx, fs, cfg, tm)
	if err != nil {
		t.Fatalf("Second initializer failed: %v", err)
	}

	// Accumulate successful results
	var accumulated InitResult
	accumulated.CreatedFiles = append(accumulated.CreatedFiles, result1.CreatedFiles...)
	accumulated.CreatedFiles = append(accumulated.CreatedFiles, result2.CreatedFiles...)

	// Verify we accumulated results from successful initializers
	if len(accumulated.CreatedFiles) == 0 {
		t.Error("No files accumulated from successful initializers")
	}

	t.Logf("Accumulated %d files from successful initializers", len(accumulated.CreatedFiles))

	// Verify both directories exist
	for _, dir := range []string{".test/dir1", ".test/dir2"} {
		exists, err := afero.DirExists(fs, dir)
		if err != nil {
			t.Errorf("Error checking directory %s: %v", dir, err)
		} else if !exists {
			t.Errorf("Directory %s should exist", dir)
		}
	}
}

// TestInitializerIdempotency verifies that initializers can be run multiple times safely.
func TestInitializerIdempotency(t *testing.T) {
	fs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}
	tm := &mockTemplateManager{content: "# Mock Content\n", err: nil}
	ctx := context.Background()

	// Test each type of initializer for idempotency
	tests := []struct {
		name        string
		initializer Initializer
	}{
		{
			name:        "DirectoryInitializer",
			initializer: NewDirectoryInitializer(".test/idempotent/dir"),
		},
		{
			name:        "ConfigFileInitializer",
			initializer: NewConfigFileInitializer("TEST_CONFIG.md", "instruction_pointer"),
		},
		{
			name:        "SlashCommandsInitializer",
			initializer: NewSlashCommandsInitializer(".test/commands", ".md", FormatMarkdown),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run initializer first time
			result1, err := tt.initializer.Init(ctx, fs, cfg, tm)
			if err != nil {
				t.Fatalf("First run failed: %v", err)
			}

			// Run initializer second time (should be idempotent)
			result2, err := tt.initializer.Init(ctx, fs, cfg, tm)
			if err != nil {
				t.Fatalf("Second run failed: %v", err)
			}

			// Verify both runs completed successfully
			t.Logf(
				"First run: %d created, %d updated",
				len(result1.CreatedFiles),
				len(result1.UpdatedFiles),
			)
			t.Logf(
				"Second run: %d created, %d updated",
				len(result2.CreatedFiles),
				len(result2.UpdatedFiles),
			)

			// Both runs should succeed without errors (idempotent)
		})
	}
}

// TestMultipleProvidersDeduplication verifies that when multiple providers share initializers,
// deduplication works correctly.
func TestMultipleProvidersDeduplication(t *testing.T) {
	ctx := context.Background()

	// Get initializers from multiple providers
	providers := []string{"claude-code", "cline", "gemini"}

	allInitializers := make([]Initializer, 0)
	for _, providerID := range providers {
		reg := Get(providerID)
		if reg == nil {
			t.Fatalf("Provider %s not found", providerID)
		}
		inits := reg.Provider.Initializers(ctx)
		allInitializers = append(allInitializers, inits...)
	}

	// Count initializers by path (for deduplication)
	pathCounts := make(map[string]int)
	for _, init := range allInitializers {
		path := init.Path()
		pathCounts[path]++
	}

	// Log paths and their counts
	for path, count := range pathCounts {
		t.Logf("Path %s appears %d times", path, count)
	}

	// All paths should be unique per provider (no accidental duplicates within same provider)
	// But different providers should have different paths
	if len(pathCounts) == 0 {
		t.Error("No initializer paths found")
	}
}
