package hooks

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/connerohnesorge/spectr/internal/domain"
)

func TestHandlePreToolUse(t *testing.T) {
	tests := []struct {
		name      string
		command   string
		input     HookInput
		wantBlock bool
		wantMsg   string // substring expected in message; empty = no check
	}{
		{
			name:    "blocks proposal modification",
			command: "apply",
			input: HookInput{
				ToolName: "Edit",
				ToolInput: json.RawMessage(
					`{"file_path": "spectr/changes/foo/proposal.md", "old_string": "a", "new_string": "b"}`,
				),
			},
			wantBlock: true,
			wantMsg:   "proposal.md",
		},
		{
			name:    "allows tasks.jsonc",
			command: "apply",
			input: HookInput{
				ToolName: "Write",
				ToolInput: json.RawMessage(
					`{"file_path": "spectr/changes/foo/tasks.jsonc", "content": "{}"}`,
				),
			},
			wantBlock: false,
		},
		{
			name:    "allows non-apply command",
			command: "proposal",
			input: HookInput{
				ToolName: "Edit",
				ToolInput: json.RawMessage(
					`{"file_path": "spectr/changes/foo/proposal.md", "old_string": "a", "new_string": "b"}`,
				),
			},
			wantBlock: false,
		},
		{
			name:    "allows non-file tools",
			command: "apply",
			input: HookInput{
				ToolName: "Read",
				ToolInput: json.RawMessage(
					`{"file_path": "spectr/changes/foo/proposal.md"}`,
				),
			},
			wantBlock: false,
		},
		{
			name:    "allows files outside changes",
			command: "apply",
			input: HookInput{
				ToolName: "Edit",
				ToolInput: json.RawMessage(
					`{"file_path": "src/main.go", "old_string": "a", "new_string": "b"}`,
				),
			},
			wantBlock: false,
		},
		{
			name:    "blocks Write on spec under changes",
			command: "apply",
			input: HookInput{
				ToolName: "Write",
				ToolInput: json.RawMessage(
					`{"file_path": "spectr/changes/bar/specs/auth/spec.md", "content": "new content"}`,
				),
			},
			wantBlock: true,
		},
		{
			name:    "blocks absolute path under changes",
			command: "apply",
			input: HookInput{
				ToolName: "Edit",
				ToolInput: json.RawMessage(
					`{"file_path": "/home/user/project/spectr/changes/foo/proposal.md", "old_string": "a", "new_string": "b"}`,
				),
			},
			wantBlock: true,
		},
		{
			name:    "blocks malformed tool input JSON",
			command: "apply",
			input: HookInput{
				ToolName:  "Edit",
				ToolInput: json.RawMessage(`{bad json`),
			},
			wantBlock: true,
			wantMsg:   "could not parse",
		},
		{
			name:    "blocks empty file path",
			command: "apply",
			input: HookInput{
				ToolName:  "Write",
				ToolInput: json.RawMessage(`{"file_path": ""}`),
			},
			wantBlock: true,
			wantMsg:   "empty file path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := handlePreToolUse(tt.command, &tt.input)

			if output.Blocked != tt.wantBlock {
				t.Errorf(
					"Blocked = %v, want %v",
					output.Blocked,
					tt.wantBlock,
				)
			}

			if tt.wantMsg == "" {
				return
			}

			if output.Message == nil {
				t.Fatal("expected message when blocked")
			}
			if !strings.Contains(*output.Message, tt.wantMsg) {
				t.Errorf(
					"message %q should contain %q",
					*output.Message,
					tt.wantMsg,
				)
			}
		})
	}
}

func TestHandle_FullRoundTrip(t *testing.T) {
	inputJSON := `{
		"session_id": "test",
		"hook_event_name": "PreToolUse",
		"tool_name": "Edit",
		"tool_input": {"file_path": "spectr/changes/foo/proposal.md", "old_string": "a", "new_string": "b"}
	}`

	var buf bytes.Buffer
	err := Handle(
		domain.HookPreToolUse,
		"apply",
		strings.NewReader(inputJSON),
		&buf,
	)
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}

	var output HookOutput
	if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
		t.Fatalf("failed to parse output: %v", err)
	}

	if !output.Blocked {
		t.Error("expected blocked output from full round trip")
	}
}

func TestHandle_NoopHookType(t *testing.T) {
	inputJSON := `{
		"session_id": "test",
		"hook_event_name": "Stop",
		"tool_name": "",
		"tool_input": {}
	}`

	var buf bytes.Buffer
	err := Handle(
		domain.HookStop,
		"apply",
		strings.NewReader(inputJSON),
		&buf,
	)
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}

	var output HookOutput
	if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
		t.Fatalf("failed to parse output: %v", err)
	}

	if output.Blocked {
		t.Error("expected non-blocking output for Stop hook type")
	}
}
