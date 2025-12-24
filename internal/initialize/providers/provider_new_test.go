package providers

import (
	"context"
	"testing"
	"text/template"

	"github.com/connerohnesorge/spectr/internal/initialize/templates"
	"github.com/spf13/afero"
)

// TestNewProviderInterface verifies the new Provider interface works correctly.
// Task 8.2: Create provider_new_test.go with tests for new Provider interface
func TestNewProviderInterface(t *testing.T) {
	ctx := context.Background()

	// Test that ClaudeProvider implements the Provider interface
	var p Provider = &ClaudeProvider{}
	initializers := p.Initializers(ctx)

	if len(initializers) == 0 {
		t.Error(
			"ClaudeProvider.Initializers() returned empty list",
		)
	}

	// Verify initializers have expected methods
	for i, init := range initializers {
		if init == nil {
			t.Errorf("Initializer %d is nil", i)

			continue
		}

		// Verify Path() returns non-empty string
		path := init.Path()
		if path == "" {
			t.Errorf(
				"Initializer %d has empty Path()",
				i,
			)
		}

		// Verify IsGlobal() returns a boolean (no panic)
		_ = init.IsGlobal()

		// Verify IsSetup() doesn't panic with nil fs (it should handle gracefully or return false)
		// We'll test with a real MemMapFs below
	}
}

// TestAllProvidersReturnInitializers verifies all 15+ providers return expected initializers.
// Task 8.3: Add tests verifying all 15+ providers return expected initializers
func TestAllProvidersReturnInitializers(
	t *testing.T,
) {
	ctx := context.Background()
	allProviders := AllProviders()

	if len(allProviders) < 15 {
		t.Errorf(
			"Expected at least 15 providers, got %d",
			len(allProviders),
		)
	}

	for _, reg := range allProviders {
		t.Run(reg.ID, func(t *testing.T) {
			initializers := reg.Provider.Initializers(
				ctx,
			)

			// Every provider should return at least one initializer
			if len(initializers) == 0 {
				t.Errorf(
					"Provider %s returned no initializers",
					reg.ID,
				)

				return
			}

			// Count initializer types
			var hasDirectory, hasConfig, hasSlashCommands bool

			for _, init := range initializers {
				switch init.(type) {
				case *DirectoryInitializer:
					hasDirectory = true
				case *ConfigFileInitializer:
					hasConfig = true
				case *SlashCommandsInitializer:
					hasSlashCommands = true
				}
			}

			// Every provider should have at least:
			// - DirectoryInitializer (to create command directory)
			// - SlashCommandsInitializer (for proposal/apply commands)
			if !hasDirectory {
				t.Errorf(
					"Provider %s missing DirectoryInitializer",
					reg.ID,
				)
			}
			if !hasSlashCommands {
				t.Errorf(
					"Provider %s missing SlashCommandsInitializer",
					reg.ID,
				)
			}

			// Some providers have config files (Claude, Codex), others don't (Gemini, Cursor)
			// So we don't enforce hasConfig for all providers
			_ = hasConfig

			// Verify all initializers have valid paths
			for i, init := range initializers {
				path := init.Path()
				if path == "" {
					t.Errorf(
						"Provider %s initializer %d has empty path",
						reg.ID,
						i,
					)
				}
			}
		})
	}
}

// TestProviderRegistrationMetadata verifies provider registration metadata (ID, Name, Priority).
// Task 8.4: Add tests verifying provider registration metadata (ID, Name, Priority)
func TestProviderRegistrationMetadata(
	t *testing.T,
) {
	allProviders := AllProviders()

	// Track used IDs and priorities to detect duplicates
	seenIDs := make(map[string]bool)
	seenPriorities := make(map[int]string)

	for _, reg := range allProviders {
		// Test ID is non-empty and kebab-case
		if reg.ID == "" {
			t.Error(
				"Found provider with empty ID",
			)

			continue
		}

		// Verify ID is kebab-case (lowercase, numbers, hyphens only)
		for _, char := range reg.ID {
			if (char < 'a' || char > 'z') &&
				(char < '0' || char > '9') &&
				char != '-' {
				t.Errorf(
					"Provider ID %q is not kebab-case (invalid char: %c)",
					reg.ID,
					char,
				)
			}
		}

		// Test ID is unique
		if seenIDs[reg.ID] {
			t.Errorf(
				"Duplicate provider ID: %s",
				reg.ID,
			)
		}
		seenIDs[reg.ID] = true

		// Test Name is non-empty
		if reg.Name == "" {
			t.Errorf(
				"Provider %s has empty Name",
				reg.ID,
			)
		}

		// Test Priority is positive
		if reg.Priority < 1 {
			t.Errorf(
				"Provider %s has invalid priority: %d (must be >= 1)",
				reg.ID,
				reg.Priority,
			)
		}

		// Test Priority is unique
		if existingID, exists := seenPriorities[reg.Priority]; exists {
			t.Errorf(
				"Duplicate priority %d for providers %s and %s",
				reg.Priority,
				existingID,
				reg.ID,
			)
		}
		seenPriorities[reg.Priority] = reg.ID

		// Test Provider implementation is not nil
		if reg.Provider == nil {
			t.Errorf(
				"Provider %s has nil Provider implementation",
				reg.ID,
			)
		}
	}
}

// TestPrioritySorting verifies AllProviders() returns providers sorted by priority.
// Task 8.4: Verify priority-based sorting
func TestPrioritySorting(t *testing.T) {
	allProviders := AllProviders()

	// Verify they're sorted by priority (ascending)
	for i := 1; i < len(allProviders); i++ {
		prev := allProviders[i-1]
		curr := allProviders[i]

		if prev.Priority > curr.Priority {
			t.Errorf(
				"Providers not sorted by priority: %s (priority %d) comes before %s (priority %d)",
				prev.ID,
				prev.Priority,
				curr.ID,
				curr.Priority,
			)
		}
	}
}

// TestFullInitializationFlow tests the complete initialization flow using afero.MemMapFs.
// Task 8.7: Add integration test for full initialization flow using afero.MemMapFs
func TestFullInitializationFlow(t *testing.T) {
	ctx := context.Background()

	// Create in-memory filesystem for testing
	fs := afero.NewMemMapFs()

	// Create a mock template manager
	tm := &testMockTemplateManager{}

	// Create config
	cfg := &Config{
		SpectrDir: "spectr",
	}

	// Get Claude provider as test subject
	reg, ok := GetProvider("claude-code")
	if !ok {
		t.Fatal(
			"Failed to get claude-code provider",
		)
	}

	// Get initializers from provider
	initializers := reg.Provider.Initializers(ctx)
	if len(initializers) == 0 {
		t.Fatal(
			"Claude provider returned no initializers",
		)
	}

	// Execute initializers in order (simulating what executor.go does)
	var aggregatedResult InitResult

	for _, init := range initializers {
		// Check if already setup
		isSetup := init.IsSetup(fs, cfg)
		if isSetup {
			t.Logf(
				"Initializer for %s already setup (unexpected on first run)",
				init.Path(),
			)
		}

		// Run initialization
		result, err := init.Init(ctx, fs, cfg, tm)
		if err != nil {
			t.Fatalf(
				"Initializer for %s failed: %v",
				init.Path(),
				err,
			)
		}

		// Aggregate results
		aggregatedResult = aggregatedResult.Merge(
			result,
		)

		// Verify IsSetup returns true after initialization
		if !init.IsSetup(fs, cfg) {
			t.Errorf(
				"IsSetup() returned false after Init() for %s",
				init.Path(),
			)
		}
	}

	// Verify we created some files
	if aggregatedResult.IsEmpty() {
		t.Error(
			"Full initialization created no files",
		)
	}

	if len(aggregatedResult.CreatedFiles) == 0 {
		t.Error(
			"No files were created during initialization",
		)
	}

	// Verify expected files exist on filesystem
	expectedPaths := []string{
		".claude/commands/spectr",             // directory
		"CLAUDE.md",                           // config file
		".claude/commands/spectr/proposal.md", // slash command
		".claude/commands/spectr/apply.md",    // slash command
	}

	for _, path := range expectedPaths {
		exists, err := afero.Exists(fs, path)
		if err != nil {
			t.Errorf(
				"Error checking if %s exists: %v",
				path,
				err,
			)

			continue
		}
		if !exists {
			t.Errorf(
				"Expected path %s does not exist after initialization",
				path,
			)
		}
	}

	// Verify files are non-empty
	for _, path := range expectedPaths {
		// Skip directory check
		isDir, _ := afero.IsDir(fs, path)
		if isDir {
			continue
		}

		content, err := afero.ReadFile(fs, path)
		if err != nil {
			t.Errorf(
				"Failed to read %s: %v",
				path,
				err,
			)

			continue
		}

		if len(content) == 0 {
			t.Errorf("File %s is empty", path)
		}
	}
}

// TestInitResultAccumulation verifies InitResult accumulation across multiple initializers.
// Task 8.8: Add integration test verifying InitResult accumulation
func TestInitResultAccumulation(t *testing.T) {
	ctx := context.Background()
	fs := afero.NewMemMapFs()
	tm := &testMockTemplateManager{}
	cfg := &Config{
		SpectrDir: "spectr",
	}

	// Get a provider with multiple initializers (Claude has 3)
	reg, ok := GetProvider("claude-code")
	if !ok {
		t.Fatal(
			"Failed to get claude-code provider",
		)
	}

	initializers := reg.Provider.Initializers(ctx)
	if len(initializers) < 2 {
		t.Fatal(
			"Need provider with at least 2 initializers for accumulation test",
		)
	}

	// Execute each initializer and accumulate results
	var accumulated InitResult
	individualResults := make(
		[]InitResult,
		0,
		len(initializers),
	)

	for _, init := range initializers {
		result, err := init.Init(ctx, fs, cfg, tm)
		if err != nil {
			t.Fatalf(
				"Initializer failed: %v",
				err,
			)
		}

		individualResults = append(
			individualResults,
			result,
		)
		accumulated = accumulated.Merge(result)
	}

	// Verify accumulated result contains all individual results
	totalCreated := 0
	totalUpdated := 0
	for _, result := range individualResults {
		totalCreated += len(result.CreatedFiles)
		totalUpdated += len(result.UpdatedFiles)
	}

	if len(
		accumulated.CreatedFiles,
	) != totalCreated {
		t.Errorf(
			"Accumulated CreatedFiles count %d != sum of individual results %d",
			len(
				accumulated.CreatedFiles,
			),
			totalCreated,
		)
	}

	if len(
		accumulated.UpdatedFiles,
	) != totalUpdated {
		t.Errorf(
			"Accumulated UpdatedFiles count %d != sum of individual results %d",
			len(
				accumulated.UpdatedFiles,
			),
			totalUpdated,
		)
	}

	// Verify TotalFiles() method
	expectedTotal := totalCreated + totalUpdated
	if accumulated.TotalFiles() != expectedTotal {
		t.Errorf(
			"TotalFiles() returned %d, expected %d",
			accumulated.TotalFiles(),
			expectedTotal,
		)
	}

	// Verify IsEmpty() returns false when we have results
	if expectedTotal > 0 &&
		accumulated.IsEmpty() {
		t.Error(
			"IsEmpty() returned true when we have results",
		)
	}

	// Test empty result
	emptyResult := InitResult{}
	if !emptyResult.IsEmpty() {
		t.Error(
			"IsEmpty() should return true for empty result",
		)
	}
}

// TestIdempotentInitialization verifies running initializers twice doesn't cause errors.
func TestIdempotentInitialization(t *testing.T) {
	ctx := context.Background()
	fs := afero.NewMemMapFs()
	tm := &testMockTemplateManager{}
	cfg := &Config{
		SpectrDir: "spectr",
	}

	reg, ok := GetProvider("claude-code")
	if !ok {
		t.Fatal(
			"Failed to get claude-code provider",
		)
	}

	initializers := reg.Provider.Initializers(ctx)

	// Run initialization first time
	var firstRunResult InitResult
	for _, init := range initializers {
		result, err := init.Init(ctx, fs, cfg, tm)
		if err != nil {
			t.Fatalf(
				"First initialization failed: %v",
				err,
			)
		}
		firstRunResult = firstRunResult.Merge(
			result,
		)
	}

	// Verify first run created files
	if len(firstRunResult.CreatedFiles) == 0 {
		t.Error(
			"First run should have created files",
		)
	}

	// Run initialization second time (should be idempotent)
	var secondRunResult InitResult
	for _, init := range initializers {
		result, err := init.Init(ctx, fs, cfg, tm)
		if err != nil {
			t.Fatalf(
				"Second initialization failed: %v",
				err,
			)
		}
		secondRunResult = secondRunResult.Merge(
			result,
		)
	}

	// Second run should not create new files (already exist)
	// ConfigFileInitializer may update files with markers, but shouldn't create new ones
	if len(secondRunResult.CreatedFiles) > 0 {
		t.Logf(
			"Warning: Second run created %d files (expected 0 for full idempotency): %v",
			len(
				secondRunResult.CreatedFiles,
			),
			secondRunResult.CreatedFiles,
		)
	}
}

// TestGeminiProviderTOMLFormat verifies Gemini provider uses TOML format.
func TestGeminiProviderTOMLFormat(t *testing.T) {
	ctx := context.Background()
	fs := afero.NewMemMapFs()
	tm := &testMockTemplateManager{}
	cfg := &Config{
		SpectrDir: "spectr",
	}

	reg, ok := GetProvider("gemini")
	if !ok {
		t.Fatal("Failed to get gemini provider")
	}

	initializers := reg.Provider.Initializers(ctx)

	// Find SlashCommandsInitializer
	var slashInit *SlashCommandsInitializer
	for _, init := range initializers {
		if si, ok := init.(*SlashCommandsInitializer); ok {
			slashInit = si

			break
		}
	}

	if slashInit == nil {
		t.Fatal(
			"Gemini provider has no SlashCommandsInitializer",
		)
	}

	// Execute initialization
	_, err := slashInit.Init(ctx, fs, cfg, tm)
	if err != nil {
		t.Fatalf(
			"SlashCommandsInitializer failed: %v",
			err,
		)
	}

	// Verify TOML files were created
	tomlFiles := []string{
		".gemini/commands/spectr/proposal.toml",
		".gemini/commands/spectr/apply.toml",
	}

	for _, path := range tomlFiles {
		exists, err := afero.Exists(fs, path)
		if err != nil {
			t.Errorf(
				"Error checking %s: %v",
				path,
				err,
			)

			continue
		}
		if !exists {
			t.Errorf(
				"Expected TOML file %s does not exist",
				path,
			)
		}
	}
}

// testMockTemplateManager implements TemplateManager for testing
// Named with 'test' prefix to avoid collision with mockTemplateManager in configfile_test.go
type testMockTemplateManager struct{}

func (*testMockTemplateManager) RenderAgents(
	_ TemplateContext,
) (string, error) {
	return "# Mock AGENTS content", nil
}

func (*testMockTemplateManager) RenderInstructionPointer(
	_ TemplateContext,
) (string, error) {
	return "# Mock Spectr Instructions\nRead spectr/AGENTS.md", nil
}

func (*testMockTemplateManager) RenderSlashCommand(
	commandType string,
	_ TemplateContext,
) (string, error) {
	return "Mock slash command content for " + commandType, nil
}

func (*testMockTemplateManager) InstructionPointer() any {
	tmpl := template.Must(
		template.New("instruction-pointer.md.tmpl").
			Parse("# Mock Spectr Instructions\nRead spectr/AGENTS.md"),
	)

	return templates.NewTemplateRef(
		"instruction-pointer.md.tmpl",
		tmpl,
	)
}

func (*testMockTemplateManager) Agents() any {
	tmpl := template.Must(
		template.New("AGENTS.md.tmpl").
			Parse("# Mock AGENTS content"),
	)

	return templates.NewTemplateRef(
		"AGENTS.md.tmpl",
		tmpl,
	)
}

func (*testMockTemplateManager) Project() any {
	tmpl := template.Must(
		template.New("project.md.tmpl").
			Parse("# Mock PROJECT content"),
	)

	return templates.NewTemplateRef(
		"project.md.tmpl",
		tmpl,
	)
}

func (*testMockTemplateManager) CIWorkflow() any {
	tmpl := template.Must(
		template.New("spectr-ci.yml.tmpl").
			Parse("# Mock CI Workflow"),
	)

	return templates.NewTemplateRef(
		"spectr-ci.yml.tmpl",
		tmpl,
	)
}

func (*testMockTemplateManager) SlashCommand(
	_ any,
) any {
	tmpl := template.Must(
		template.New("slash-command.md.tmpl").
			Parse("Mock slash command content"),
	)

	return templates.NewTemplateRef(
		"slash-command.md.tmpl",
		tmpl,
	)
}
