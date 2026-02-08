package domain

import (
	"strings"
	"testing"
)

func TestSlashNextFrontmatter(t *testing.T) {
	// Get the base frontmatter for SlashNext
	fm := GetBaseFrontmatter(SlashNext)

	// Verify expected scalar field values
	expectedValues := map[string]any{
		"description":   "Spectr: Next Task Execution",
		"allowed-tools": "Read, Glob, Grep, Write, Edit, Bash(spectr:*)",
		"subtask":       false,
	}

	for key, expectedValue := range expectedValues {
		actualValue, exists := fm[key]
		if !exists {
			t.Errorf(
				"SlashNext frontmatter missing field %q",
				key,
			)

			continue
		}

		if actualValue != expectedValue {
			t.Errorf(
				"SlashNext frontmatter field %q = %v, want %v",
				key,
				actualValue,
				expectedValue,
			)
		}
	}

	// Verify hooks field exists and is a map
	hooksVal, hasHooks := fm["hooks"]
	if !hasHooks {
		t.Error("SlashNext frontmatter missing field \"hooks\"")
	} else if hooksMap, ok := hooksVal.(map[string]any); !ok {
		t.Error("SlashNext frontmatter \"hooks\" is not map[string]any")
	} else if len(hooksMap) != len(AllHookTypes()) {
		t.Errorf(
			"SlashNext frontmatter hooks has %d entries, want %d",
			len(hooksMap),
			len(AllHookTypes()),
		)
	}

	// Verify total field count (3 scalar + 1 hooks = 4)
	if len(fm) != 4 {
		t.Errorf(
			"SlashNext frontmatter has %d fields, want 4",
			len(fm),
		)
	}
}

func TestSlashNextRenderedFrontmatter(t *testing.T) {
	// Get the base frontmatter for SlashNext
	fm := GetBaseFrontmatter(SlashNext)

	// Render it with sample body content
	body := "# Spectr: Next Task Execution\n\nThis is the command body.\n"
	rendered, err := RenderFrontmatter(fm, body)
	if err != nil {
		t.Fatalf(
			"RenderFrontmatter failed: %v",
			err,
		)
	}

	// Verify the rendered output contains expected frontmatter
	// Note: YAML adds quotes around values with colons
	expectedStrings := []string{
		"---",
		"allowed-tools: Read, Glob, Grep, Write, Edit, Bash(spectr:*)",
		"description: 'Spectr: Next Task Execution'",
		"hooks:",
		"subtask: false",
		"# Spectr: Next Task Execution",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(rendered, expected) {
			t.Errorf(
				"Rendered frontmatter missing expected string %q\nGot:\n%s",
				expected,
				rendered,
			)
		}
	}

	// Verify the frontmatter block structure
	if !strings.HasPrefix(rendered, "---\n") {
		t.Error("Rendered output should start with '---\\n'")
	}

	parts := strings.Split(rendered, "---\n")
	if len(parts) < 3 {
		t.Errorf(
			"Rendered output should have 3+ parts separated by '---\\n', got %d",
			len(parts),
		)
	}
}
