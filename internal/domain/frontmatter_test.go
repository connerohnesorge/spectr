package domain

import (
	"strings"
	"testing"
)

func TestGetBaseFrontmatter(t *testing.T) {
	tests := []struct {
		name       string
		cmd        SlashCommand
		wantFields []string
		wantAgent  bool // whether "agent" field should be present
	}{
		{
			name:       "proposal command has expected fields",
			cmd:        SlashProposal,
			wantFields: []string{"description", "allowed-tools", "agent", "subtask"},
			wantAgent:  true,
		},
		{
			name:       "apply command has expected fields",
			cmd:        SlashApply,
			wantFields: []string{"description", "allowed-tools", "subtask"},
			wantAgent:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm := GetBaseFrontmatter(tt.cmd)

			// Check expected fields exist
			for _, field := range tt.wantFields {
				if _, ok := fm[field]; !ok {
					t.Errorf("GetBaseFrontmatter(%v) missing field %q", tt.cmd, field)
				}
			}

			// Check agent field presence
			_, hasAgent := fm["agent"]
			if hasAgent != tt.wantAgent {
				t.Errorf(
					"GetBaseFrontmatter(%v) agent field presence = %v, want %v",
					tt.cmd,
					hasAgent,
					tt.wantAgent,
				)
			}
		})
	}
}

const testModifiedValue = "modified"

func TestGetBaseFrontmatter_ReturnsDeepCopy(t *testing.T) {
	fm1 := GetBaseFrontmatter(SlashProposal)
	fm2 := GetBaseFrontmatter(SlashProposal)

	// Modify fm1
	fm1["description"] = testModifiedValue

	// fm2 should be unaffected
	if fm2["description"] == testModifiedValue {
		t.Error("GetBaseFrontmatter did not return a deep copy - mutations affect other copies")
	}

	// Original BaseSlashCommandFrontmatter should be unaffected
	original := BaseSlashCommandFrontmatter[SlashProposal]
	if original["description"] == testModifiedValue {
		t.Error("GetBaseFrontmatter did not return a deep copy - mutations affect original")
	}
}

func TestGetBaseFrontmatter_UnknownCommand(t *testing.T) {
	// Using an invalid command value
	fm := GetBaseFrontmatter(SlashCommand(999))
	if fm == nil {
		t.Error("GetBaseFrontmatter should return empty map for unknown command, got nil")
	}
	if len(fm) != 0 {
		t.Errorf("GetBaseFrontmatter should return empty map for unknown command, got %v", fm)
	}
}

func TestCopyMap(t *testing.T) {
	tests := []struct {
		name string
		src  map[string]any
	}{
		{
			name: "nil map",
			src:  nil,
		},
		{
			name: "empty map",
			src:  make(map[string]any),
		},
		{
			name: "simple map",
			src: map[string]any{
				"string": "value",
				"bool":   true,
				"int":    42,
			},
		},
		{
			name: "nested map",
			src: map[string]any{
				"outer": map[string]any{
					"inner": "value",
				},
			},
		},
		{
			name: "with slice",
			src: map[string]any{
				"items": []any{"a", "b", "c"},
			},
		},
		{
			name: "with string slice",
			src: map[string]any{
				"tags": []string{"tag1", "tag2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := copyMap(tt.src)

			if tt.src == nil {
				if dst != nil {
					t.Error("copyMap(nil) should return nil")
				}

				return
			}

			// Verify same length
			if len(dst) != len(tt.src) {
				t.Errorf("copyMap length = %d, want %d", len(dst), len(tt.src))
			}
		})
	}
}

func TestCopyMap_Mutation(t *testing.T) {
	src := map[string]any{
		"key": "original",
		"nested": map[string]any{
			"inner": "value",
		},
	}

	dst := copyMap(src)

	// Modify dst
	dst["key"] = "modified"
	if nested, ok := dst["nested"].(map[string]any); ok {
		nested["inner"] = "modified"
	}

	// src should be unaffected
	if src["key"] != "original" {
		t.Error("copyMap: modifying copy affected original (top-level)")
	}

	nested, ok := src["nested"].(map[string]any)
	if !ok {
		return
	}
	if nested["inner"] != "value" {
		t.Error("copyMap: modifying copy affected original (nested)")
	}
}

func TestApplyFrontmatterOverrides(t *testing.T) {
	tests := []struct {
		name      string
		base      map[string]any
		overrides *FrontmatterOverride
		wantKeys  []string
		wantValue map[string]any // specific values to check
	}{
		{
			name: "nil overrides returns copy of base",
			base: map[string]any{
				"key": "value",
			},
			overrides: nil,
			wantKeys:  []string{"key"},
			wantValue: map[string]any{"key": "value"},
		},
		{
			name: "set adds new field",
			base: map[string]any{
				"existing": "value",
			},
			overrides: &FrontmatterOverride{
				Set: map[string]any{
					"new": "field",
				},
			},
			wantKeys: []string{"existing", "new"},
			wantValue: map[string]any{
				"existing": "value",
				"new":      "field",
			},
		},
		{
			name: "set overwrites existing field",
			base: map[string]any{
				"key": "original",
			},
			overrides: &FrontmatterOverride{
				Set: map[string]any{
					"key": "modified",
				},
			},
			wantKeys:  []string{"key"},
			wantValue: map[string]any{"key": "modified"},
		},
		{
			name: "remove deletes field",
			base: map[string]any{
				"keep":   "value",
				"remove": "me",
			},
			overrides: &FrontmatterOverride{
				Remove: []string{"remove"},
			},
			wantKeys:  []string{"keep"},
			wantValue: map[string]any{"keep": "value"},
		},
		{
			name: "remove ignores non-existent field",
			base: map[string]any{
				"key": "value",
			},
			overrides: &FrontmatterOverride{
				Remove: []string{"nonexistent"},
			},
			wantKeys:  []string{"key"},
			wantValue: map[string]any{"key": "value"},
		},
		{
			name: "set then remove same field removes it",
			base: map[string]any{
				"key": "original",
			},
			overrides: &FrontmatterOverride{
				Set:    map[string]any{"key": "set-value"},
				Remove: []string{"key"},
			},
			wantKeys: nil,
		},
		{
			name: "claude code proposal override scenario",
			base: map[string]any{
				"description":   "Test",
				"allowed-tools": "Read, Write",
				"agent":         "plan",
				"subtask":       false,
			},
			overrides: &FrontmatterOverride{
				Set:    map[string]any{"context": "fork"},
				Remove: []string{"agent"},
			},
			wantKeys: []string{"description", "allowed-tools", "context", "subtask"},
			wantValue: map[string]any{
				"context": "fork",
				"subtask": false,
			},
		},
		{
			name:      "nil base with overrides",
			base:      nil,
			overrides: &FrontmatterOverride{Set: map[string]any{"key": "value"}},
			wantKeys:  []string{"key"},
			wantValue: map[string]any{"key": "value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ApplyFrontmatterOverrides(tt.base, tt.overrides)

			// Check expected keys
			if len(result) != len(tt.wantKeys) {
				t.Errorf(
					"ApplyFrontmatterOverrides got %d keys, want %d",
					len(result),
					len(tt.wantKeys),
				)
			}

			for _, key := range tt.wantKeys {
				if _, ok := result[key]; !ok {
					t.Errorf("ApplyFrontmatterOverrides missing expected key %q", key)
				}
			}

			// Check specific values
			for key, wantVal := range tt.wantValue {
				if result[key] != wantVal {
					t.Errorf(
						"ApplyFrontmatterOverrides[%q] = %v, want %v",
						key,
						result[key],
						wantVal,
					)
				}
			}
		})
	}
}

func TestApplyFrontmatterOverrides_DoesNotMutateBase(t *testing.T) {
	base := map[string]any{
		"key": "original",
	}
	overrides := &FrontmatterOverride{
		Set: map[string]any{"key": "modified", "new": "value"},
	}

	ApplyFrontmatterOverrides(base, overrides)

	// base should be unaffected
	if base["key"] != "original" {
		t.Error("ApplyFrontmatterOverrides mutated base map")
	}
	if _, ok := base["new"]; ok {
		t.Error("ApplyFrontmatterOverrides added field to base map")
	}
}

func TestRenderFrontmatter(t *testing.T) {
	tests := []struct {
		name           string
		fm             map[string]any
		body           string
		wantPrefix     string // prefix the result should have
		wantSuffix     string // suffix the result should have
		wantContain    []string
		wantNotContain []string
	}{
		{
			name:       "empty frontmatter",
			fm:         make(map[string]any),
			body:       "# Body",
			wantPrefix: "---\n---\n",
			wantSuffix: "# Body",
		},
		{
			name: "simple frontmatter",
			fm: map[string]any{
				"description": "Test",
			},
			body:        "# Body",
			wantPrefix:  "---\n",
			wantSuffix:  "---\n# Body",
			wantContain: []string{"description: Test"},
		},
		{
			name: "multiple fields",
			fm: map[string]any{
				"description": "Test",
				"subtask":     false,
			},
			body:        "# Body",
			wantContain: []string{"description: Test", "subtask: false"},
		},
		{
			name: "claude code proposal frontmatter",
			fm: map[string]any{
				"allowed-tools": "Read, Write",
				"context":       "fork",
				"description":   "Test",
				"subtask":       false,
			},
			body: "# Body",
			wantContain: []string{
				"context: fork",
				"description: Test",
				"allowed-tools: Read, Write",
				"subtask: false",
			},
			wantNotContain: []string{"agent:"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RenderFrontmatter(tt.fm, tt.body)
			if err != nil {
				t.Fatalf("RenderFrontmatter error = %v", err)
			}

			if tt.wantPrefix != "" && !strings.HasPrefix(result, tt.wantPrefix) {
				t.Errorf("RenderFrontmatter result does not have expected prefix\ngot:\n%s", result)
			}

			if tt.wantSuffix != "" && !strings.HasSuffix(result, tt.wantSuffix) {
				t.Errorf(
					"RenderFrontmatter result does not have expected suffix\ngot:\n%s\nwant suffix:\n%s",
					result,
					tt.wantSuffix,
				)
			}

			for _, want := range tt.wantContain {
				if !strings.Contains(result, want) {
					t.Errorf("RenderFrontmatter result does not contain %q\ngot:\n%s", want, result)
				}
			}

			for _, notWant := range tt.wantNotContain {
				if strings.Contains(result, notWant) {
					t.Errorf(
						"RenderFrontmatter result should not contain %q\ngot:\n%s",
						notWant,
						result,
					)
				}
			}
		})
	}
}

func TestRenderFrontmatter_DeterministicOrder(t *testing.T) {
	fm := map[string]any{
		"zebra":  "last",
		"alpha":  "first",
		"middle": "mid",
	}
	body := "# Body"

	// Render multiple times and verify consistent output
	first, err := RenderFrontmatter(fm, body)
	if err != nil {
		t.Fatalf("RenderFrontmatter error = %v", err)
	}

	for i := range 10 {
		result, err := RenderFrontmatter(fm, body)
		if err != nil {
			t.Fatalf("RenderFrontmatter error on iteration %d = %v", i, err)
		}
		if result != first {
			t.Errorf(
				"RenderFrontmatter output is not deterministic\nfirst:\n%s\niteration %d:\n%s",
				first,
				i,
				result,
			)
		}
	}
}
