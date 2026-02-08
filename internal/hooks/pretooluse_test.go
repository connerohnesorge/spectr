package hooks

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/connerohnesorge/spectr/internal/domain"
)

func TestHandlePreToolUse_BlocksProposalModification(
	t *testing.T,
) {
	input := HookInput{
		SessionID:     "test-session",
		HookEventName: "PreToolUse",
		ToolName:      "Edit",
		ToolInput: json.RawMessage(
			`{"file_path": "spectr/changes/foo/proposal.md", "old_string": "a", "new_string": "b"}`,
		),
	}

	output := handlePreToolUse("apply", &input)

	if !output.Blocked {
		t.Error("expected Edit on proposal.md to be blocked")
	}
	if output.Message == nil {
		t.Fatal("expected message when blocked")
	}
	if !strings.Contains(*output.Message, "proposal.md") {
		t.Errorf(
			"message should mention the file, got: %s",
			*output.Message,
		)
	}
}

func TestHandlePreToolUse_AllowsTasksJsonc(
	t *testing.T,
) {
	input := HookInput{
		SessionID:     "test-session",
		HookEventName: "PreToolUse",
		ToolName:      "Write",
		ToolInput: json.RawMessage(
			`{"file_path": "spectr/changes/foo/tasks.jsonc", "content": "{}"}`,
		),
	}

	output := handlePreToolUse("apply", &input)

	if output.Blocked {
		t.Error("expected Write on tasks.jsonc to be allowed")
	}
}

func TestHandlePreToolUse_AllowsNonApplyCommand(
	t *testing.T,
) {
	input := HookInput{
		SessionID:     "test-session",
		HookEventName: "PreToolUse",
		ToolName:      "Edit",
		ToolInput: json.RawMessage(
			`{"file_path": "spectr/changes/foo/proposal.md", "old_string": "a", "new_string": "b"}`,
		),
	}

	output := handlePreToolUse("proposal", &input)

	if output.Blocked {
		t.Error("expected non-apply command to not block")
	}
}

func TestHandlePreToolUse_AllowsNonFileTools(
	t *testing.T,
) {
	input := HookInput{
		SessionID:     "test-session",
		HookEventName: "PreToolUse",
		ToolName:      "Read",
		ToolInput: json.RawMessage(
			`{"file_path": "spectr/changes/foo/proposal.md"}`,
		),
	}

	output := handlePreToolUse("apply", &input)

	if output.Blocked {
		t.Error("expected Read tool to not be blocked")
	}
}

func TestHandlePreToolUse_AllowsFilesOutsideChanges(
	t *testing.T,
) {
	input := HookInput{
		SessionID:     "test-session",
		HookEventName: "PreToolUse",
		ToolName:      "Edit",
		ToolInput: json.RawMessage(
			`{"file_path": "src/main.go", "old_string": "a", "new_string": "b"}`,
		),
	}

	output := handlePreToolUse("apply", &input)

	if output.Blocked {
		t.Error("expected files outside spectr/changes/ to not be blocked")
	}
}

func TestHandlePreToolUse_BlocksWriteTool(
	t *testing.T,
) {
	input := HookInput{
		SessionID:     "test-session",
		HookEventName: "PreToolUse",
		ToolName:      "Write",
		ToolInput: json.RawMessage(
			`{"file_path": "spectr/changes/bar/specs/auth/spec.md", "content": "new content"}`,
		),
	}

	output := handlePreToolUse("apply", &input)

	if !output.Blocked {
		t.Error("expected Write on spec.md under changes to be blocked")
	}
}

func TestHandlePreToolUse_AbsolutePath(
	t *testing.T,
) {
	input := HookInput{
		SessionID:     "test-session",
		HookEventName: "PreToolUse",
		ToolName:      "Edit",
		ToolInput: json.RawMessage(
			`{"file_path": "/home/user/project/spectr/changes/foo/proposal.md", "old_string": "a", "new_string": "b"}`,
		),
	}

	output := handlePreToolUse("apply", &input)

	if !output.Blocked {
		t.Error("expected absolute path under spectr/changes/ to be blocked")
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
