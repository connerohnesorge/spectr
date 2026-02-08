package domain

import (
	"strings"
	"testing"
)

func TestBuildHooksFrontmatter(t *testing.T) {
	tests := []struct {
		name string
		cmd  SlashCommand
	}{
		{"proposal", SlashProposal},
		{"apply", SlashApply},
		{"next", SlashNext},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hooks := BuildHooksFrontmatter(tt.cmd)

			// Should have one entry per hook type
			if len(hooks) != len(AllHookTypes()) {
				t.Errorf(
					"BuildHooksFrontmatter(%v) returned %d entries, want %d",
					tt.cmd,
					len(hooks),
					len(AllHookTypes()),
				)
			}

			// Verify each hook type is present with correct structure
			for _, ht := range AllHookTypes() {
				entry, ok := hooks[ht.String()]
				if !ok {
					t.Errorf(
						"BuildHooksFrontmatter(%v) missing hook type %q",
						tt.cmd,
						ht.String(),
					)

					continue
				}

				// Verify entry is []any with one element
				arr, ok := entry.([]any)
				if !ok || len(arr) != 1 {
					t.Errorf(
						"BuildHooksFrontmatter(%v)[%q] expected []any with 1 element",
						tt.cmd,
						ht.String(),
					)

					continue
				}

				// Verify the matcher/hooks structure
				hookEntry, ok := arr[0].(map[string]any)
				if !ok {
					t.Errorf(
						"BuildHooksFrontmatter(%v)[%q][0] expected map[string]any",
						tt.cmd,
						ht.String(),
					)

					continue
				}

				// Check matcher
				if matcher, ok := hookEntry["matcher"]; !ok || matcher != "" {
					t.Errorf(
						"BuildHooksFrontmatter(%v)[%q] matcher = %v, want empty string",
						tt.cmd,
						ht.String(),
						matcher,
					)
				}

				// Check hooks array
				innerHooks, ok := hookEntry["hooks"].([]any)
				if !ok || len(innerHooks) != 1 {
					t.Errorf(
						"BuildHooksFrontmatter(%v)[%q] hooks expected []any with 1 element",
						tt.cmd,
						ht.String(),
					)

					continue
				}

				// Check command structure
				hookCmd, ok := innerHooks[0].(map[string]any)
				if !ok {
					t.Errorf(
						"BuildHooksFrontmatter(%v)[%q] hook command expected map[string]any",
						tt.cmd,
						ht.String(),
					)

					continue
				}

				if hookCmd["type"] != "command" {
					t.Errorf(
						"BuildHooksFrontmatter(%v)[%q] type = %v, want \"command\"",
						tt.cmd,
						ht.String(),
						hookCmd["type"],
					)
				}

				// Verify command string contains hook type and command name
				cmdStr, ok := hookCmd["command"].(string)
				if !ok {
					t.Errorf(
						"BuildHooksFrontmatter(%v)[%q] command is not a string",
						tt.cmd,
						ht.String(),
					)

					continue
				}

				if !strings.Contains(cmdStr, ht.String()) {
					t.Errorf(
						"BuildHooksFrontmatter(%v)[%q] command %q does not contain hook type",
						tt.cmd,
						ht.String(),
						cmdStr,
					)
				}

				if !strings.Contains(cmdStr, tt.cmd.String()) {
					t.Errorf(
						"BuildHooksFrontmatter(%v)[%q] command %q does not contain command name",
						tt.cmd,
						ht.String(),
						cmdStr,
					)
				}

				expectedCmd := "spectr hooks " + ht.String() + " --command " + tt.cmd.String()
				if cmdStr != expectedCmd {
					t.Errorf(
						"BuildHooksFrontmatter(%v)[%q] command = %q, want %q",
						tt.cmd,
						ht.String(),
						cmdStr,
						expectedCmd,
					)
				}

				if hookCmd["timeout"] != 600 {
					t.Errorf(
						"BuildHooksFrontmatter(%v)[%q] timeout = %v, want 600",
						tt.cmd,
						ht.String(),
						hookCmd["timeout"],
					)
				}
			}
		})
	}
}
