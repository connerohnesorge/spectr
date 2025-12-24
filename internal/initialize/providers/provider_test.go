package providers

import (
	"context"
	"testing"
)

// TestClaudeProvider tests that ClaudeProvider returns correct initializers
func TestClaudeProvider(t *testing.T) {
	p := &ClaudeProvider{}
	ctx := context.Background()

	inits := p.Initializers(ctx)

	if len(inits) != 3 {
		t.Errorf(
			"ClaudeProvider should return 3 initializers, got %d",
			len(inits),
		)
	}

	// Verify paths
	expectedPaths := map[string]bool{
		".claude/commands/spectr": false,
		"CLAUDE.md":               false,
	}

	for _, init := range inits {
		path := init.Path()
		if _, ok := expectedPaths[path]; ok {
			expectedPaths[path] = true
		}
	}

	for path, found := range expectedPaths {
		if !found {
			t.Errorf(
				"ClaudeProvider should return initializer for path %s",
				path,
			)
		}
	}

	// Verify none are global
	for _, init := range inits {
		if init.IsGlobal() {
			t.Errorf(
				"ClaudeProvider initializers should not be global, but %s is global",
				init.Path(),
			)
		}
	}
}

// TestGeminiProvider tests that GeminiProvider returns correct initializers
func TestGeminiProvider(t *testing.T) {
	p := &GeminiProvider{}
	ctx := context.Background()

	inits := p.Initializers(ctx)

	if len(inits) != 2 {
		t.Errorf(
			"GeminiProvider should return 2 initializers, got %d",
			len(inits),
		)
	}

	// Verify paths
	expectedPaths := map[string]bool{
		".gemini/commands/spectr": false,
	}

	for _, init := range inits {
		path := init.Path()
		if _, ok := expectedPaths[path]; ok {
			expectedPaths[path] = true
		}
	}

	for path, found := range expectedPaths {
		if !found {
			t.Errorf(
				"GeminiProvider should return initializer for path %s",
				path,
			)
		}
	}
}

// TestCursorProvider tests that CursorProvider returns correct initializers
func TestCursorProvider(t *testing.T) {
	p := &CursorProvider{}
	ctx := context.Background()

	inits := p.Initializers(ctx)

	if len(inits) != 2 {
		t.Errorf(
			"CursorProvider should return 2 initializers, got %d",
			len(inits),
		)
	}

	// Verify paths
	expectedPaths := map[string]bool{
		".cursorrules/commands/spectr": false,
	}

	for _, init := range inits {
		path := init.Path()
		if _, ok := expectedPaths[path]; ok {
			expectedPaths[path] = true
		}
	}

	for path, found := range expectedPaths {
		if !found {
			t.Errorf(
				"CursorProvider should return initializer for path %s",
				path,
			)
		}
	}
}

// TestClineProvider tests that ClineProvider returns correct initializers
func TestClineProvider(t *testing.T) {
	p := &ClineProvider{}
	ctx := context.Background()

	inits := p.Initializers(ctx)

	if len(inits) != 3 {
		t.Errorf(
			"ClineProvider should return 3 initializers, got %d",
			len(inits),
		)
	}

	// Verify paths
	foundClineDir := false
	foundClineMd := false
	for _, init := range inits {
		path := init.Path()
		if path == ".clinerules/commands/spectr" {
			foundClineDir = true
		}
		if path == "CLINE.md" {
			foundClineMd = true
		}
	}

	if !foundClineDir {
		t.Error(
			"ClineProvider should return initializer for .clinerules/commands/spectr",
		)
	}
	if !foundClineMd {
		t.Error(
			"ClineProvider should return initializer for CLINE.md",
		)
	}
}

// TestContinueProvider tests that ContinueProvider returns correct initializers
func TestContinueProvider(t *testing.T) {
	p := &ContinueProvider{}
	ctx := context.Background()

	inits := p.Initializers(ctx)

	if len(inits) != 2 {
		t.Errorf(
			"ContinueProvider should return 2 initializers, got %d",
			len(inits),
		)
	}

	// Build count of found paths
	foundPaths := make(map[string]int)
	for _, init := range inits {
		foundPaths[init.Path()]++
	}

	// Verify all expected paths are present with correct counts
	// Both DirectoryInitializer and SlashCommandsInitializer use the same path
	expectedPaths := map[string]int{
		".continue/commands/spectr": 2,
	}
	for expected, expectedCount := range expectedPaths {
		if foundPaths[expected] != expectedCount {
			t.Errorf(
				"ContinueProvider path %q: got count %d, want %d",
				expected,
				foundPaths[expected],
				expectedCount,
			)
		}
	}

	// Verify no unexpected paths
	for path := range foundPaths {
		if _, ok := expectedPaths[path]; !ok {
			t.Errorf(
				"ContinueProvider has unexpected path: %s",
				path,
			)
		}
	}
}

// TestWindsurfProvider tests that WindsurfProvider returns correct initializers
func TestWindsurfProvider(t *testing.T) {
	p := &WindsurfProvider{}
	ctx := context.Background()

	inits := p.Initializers(ctx)

	if len(inits) != 2 {
		t.Errorf(
			"WindsurfProvider should return 2 initializers, got %d",
			len(inits),
		)
	}

	// Build count of found paths
	foundPaths := make(map[string]int)
	for _, init := range inits {
		foundPaths[init.Path()]++
	}

	// Verify all expected paths are present with correct counts
	// Both DirectoryInitializer and SlashCommandsInitializer use the same path
	expectedPaths := map[string]int{
		".windsurf/commands/spectr": 2,
	}
	for expected, expectedCount := range expectedPaths {
		if foundPaths[expected] != expectedCount {
			t.Errorf(
				"WindsurfProvider path %q: got count %d, want %d",
				expected,
				foundPaths[expected],
				expectedCount,
			)
		}
	}

	// Verify no unexpected paths
	for path := range foundPaths {
		if _, ok := expectedPaths[path]; !ok {
			t.Errorf(
				"WindsurfProvider has unexpected path: %s",
				path,
			)
		}
	}
}

// TestAiderProvider tests that AiderProvider returns correct initializers
func TestAiderProvider(t *testing.T) {
	p := &AiderProvider{}
	ctx := context.Background()

	inits := p.Initializers(ctx)

	if len(inits) != 2 {
		t.Errorf(
			"AiderProvider should return 2 initializers, got %d",
			len(inits),
		)
	}

	// Build count of found paths
	foundPaths := make(map[string]int)
	for _, init := range inits {
		foundPaths[init.Path()]++
	}

	// Verify all expected paths are present with correct counts
	// Both DirectoryInitializer and SlashCommandsInitializer use the same path
	expectedPaths := map[string]int{
		".aider/commands/spectr": 2,
	}
	for expected, expectedCount := range expectedPaths {
		if foundPaths[expected] != expectedCount {
			t.Errorf(
				"AiderProvider path %q: got count %d, want %d",
				expected,
				foundPaths[expected],
				expectedCount,
			)
		}
	}

	// Verify no unexpected paths
	for path := range foundPaths {
		if _, ok := expectedPaths[path]; !ok {
			t.Errorf(
				"AiderProvider has unexpected path: %s",
				path,
			)
		}
	}
}

// TestCostrictProvider tests that CostrictProvider returns correct initializers
func TestCostrictProvider(t *testing.T) {
	p := &CostrictProvider{}
	ctx := context.Background()

	inits := p.Initializers(ctx)

	if len(inits) != 3 {
		t.Errorf(
			"CostrictProvider should return 3 initializers, got %d",
			len(inits),
		)
	}

	// Build count of found paths
	foundPaths := make(map[string]int)
	for _, init := range inits {
		foundPaths[init.Path()]++
	}

	// Verify all expected paths are present with correct counts
	// DirectoryInitializer and SlashCommandsInitializer share the same path (count 2)
	// ConfigFileInitializer uses COSTRICT.md (count 1)
	expectedPaths := map[string]int{
		".costrict/commands/spectr": 2,
		"COSTRICT.md":               1,
	}
	for expected, expectedCount := range expectedPaths {
		if foundPaths[expected] != expectedCount {
			t.Errorf(
				"CostrictProvider path %q: got count %d, want %d",
				expected,
				foundPaths[expected],
				expectedCount,
			)
		}
	}

	// Verify no unexpected paths
	for path := range foundPaths {
		if _, ok := expectedPaths[path]; !ok {
			t.Errorf(
				"CostrictProvider has unexpected path: %s",
				path,
			)
		}
	}
}

// TestKilocodeProvider tests that KilocodeProvider returns correct initializers
func TestKilocodeProvider(t *testing.T) {
	p := &KilocodeProvider{}
	ctx := context.Background()

	inits := p.Initializers(ctx)

	if len(inits) != 2 {
		t.Errorf(
			"KilocodeProvider should return 2 initializers, got %d",
			len(inits),
		)
	}

	// Build count of found paths
	foundPaths := make(map[string]int)
	for _, init := range inits {
		foundPaths[init.Path()]++
	}

	// Verify all expected paths are present with correct counts
	// Both DirectoryInitializer and SlashCommandsInitializer use the same path
	expectedPaths := map[string]int{
		".kilocode/commands/spectr": 2,
	}
	for expected, expectedCount := range expectedPaths {
		if foundPaths[expected] != expectedCount {
			t.Errorf(
				"KilocodeProvider path %q: got count %d, want %d",
				expected,
				foundPaths[expected],
				expectedCount,
			)
		}
	}

	// Verify no unexpected paths
	for path := range foundPaths {
		if _, ok := expectedPaths[path]; !ok {
			t.Errorf(
				"KilocodeProvider has unexpected path: %s",
				path,
			)
		}
	}
}

// TestAntigravityProvider tests that AntigravityProvider returns correct initializers
func TestAntigravityProvider(t *testing.T) {
	p := &AntigravityProvider{}
	ctx := context.Background()

	inits := p.Initializers(ctx)

	if len(inits) != 3 {
		t.Errorf(
			"AntigravityProvider should return 3 initializers, got %d",
			len(inits),
		)
	}

	// Verify paths
	foundAntigravity := false
	foundAgents := false
	for _, init := range inits {
		path := init.Path()
		if path == ".agent/workflows/spectr" {
			foundAntigravity = true
		}
		if path == "AGENTS.md" {
			foundAgents = true
		}
	}

	if !foundAntigravity {
		t.Error(
			"AntigravityProvider should return initializer for .agent/workflows/spectr",
		)
	}
	if !foundAgents {
		t.Error(
			"AntigravityProvider should return initializer for AGENTS.md",
		)
	}
}

// TestCodexProvider tests that CodexProvider returns correct initializers
func TestCodexProvider(t *testing.T) {
	p := &CodexProvider{}
	ctx := context.Background()

	inits := p.Initializers(ctx)

	// CodexProvider currently returns 1 initializer (AGENTS.md)
	// TODO: Will return more when global slash commands are implemented
	if len(inits) != 1 {
		t.Errorf(
			"CodexProvider should return 1 initializer, got %d",
			len(inits),
		)
	}

	// Verify paths
	foundAgents := false
	for _, init := range inits {
		path := init.Path()
		if path == "AGENTS.md" {
			foundAgents = true
		}
	}

	if !foundAgents {
		t.Error(
			"CodexProvider should return initializer for AGENTS.md",
		)
	}
}

// TestAllProvidersReturnInitializers verifies all providers return at least one initializer
func TestAllProvidersReturnInitializers(
	t *testing.T,
) {
	ctx := context.Background()

	providers := []Provider{
		&ClaudeProvider{},
		&GeminiProvider{},
		&CursorProvider{},
		&ClineProvider{},
		&ContinueProvider{},
		&WindsurfProvider{},
		&AiderProvider{},
		&CostrictProvider{},
		&KilocodeProvider{},
		&AntigravityProvider{},
		&CodexProvider{},
	}

	for _, p := range providers {
		inits := p.Initializers(ctx)
		if len(inits) == 0 {
			t.Errorf(
				"Provider %T should return at least one initializer",
				p,
			)
		}

		// Verify all initializers have non-empty paths
		for i, init := range inits {
			if init.Path() == "" {
				t.Errorf(
					"Provider %T initializer %d has empty path",
					p,
					i,
				)
			}
		}
	}
}

// TestInitializersAreIdempotent verifies that calling Initializers multiple times
// returns consistent results (same number of initializers with same paths)
func TestInitializersAreIdempotent(t *testing.T) {
	ctx := context.Background()
	p := &ClaudeProvider{}

	inits1 := p.Initializers(ctx)
	inits2 := p.Initializers(ctx)

	if len(inits1) != len(inits2) {
		t.Errorf(
			"ClaudeProvider should return same number of initializers: got %d and %d",
			len(inits1),
			len(inits2),
		)
	}

	// Verify paths match
	paths1 := make(map[string]bool)
	for _, init := range inits1 {
		paths1[init.Path()] = true
	}

	paths2 := make(map[string]bool)
	for _, init := range inits2 {
		paths2[init.Path()] = true
	}

	for path := range paths1 {
		if !paths2[path] {
			t.Errorf(
				"Path %s found in first call but not second",
				path,
			)
		}
	}

	for path := range paths2 {
		if !paths1[path] {
			t.Errorf(
				"Path %s found in second call but not first",
				path,
			)
		}
	}
}
