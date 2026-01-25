package domain

import (
	"strings"
	"testing"
)

func TestSlashNextFrontmatter(t *testing.T) {
	// Get the base frontmatter for SlashNext
	fm := GetBaseFrontmatter(SlashNext)

	// Verify expected field values
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

	// Verify no unexpected fields
	if len(fm) != len(expectedValues) {
		t.Errorf(
			"SlashNext frontmatter has %d fields, want %d",
			len(fm),
			len(expectedValues),
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
