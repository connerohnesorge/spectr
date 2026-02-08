package hooks

import (
	"encoding/json"
	"regexp"
	"strings"
)

// spectrChangesPattern matches file paths under spectr/changes/<id>/.
var spectrChangesPattern = regexp.MustCompile(
	`(?:^|/)spectr/changes/[^/]+/`,
)

// fileToolNames lists the tool names that perform file write operations.
var fileToolNames = map[string]bool{
	"Edit":  true,
	"Write": true,
}

// handlePreToolUse checks if a file write operation targets protected
// files under spectr/changes/. For the "apply" command, only tasks.jsonc
// is allowed to be modified.
func handlePreToolUse(
	command string,
	input *HookInput,
) *HookOutput {
	// Only guard file writes during apply
	if command != "apply" {
		return &HookOutput{}
	}

	// Only check file-writing tools
	if !fileToolNames[input.ToolName] {
		return &HookOutput{}
	}

	// Parse tool_input to get file_path
	var ti toolInput
	if err := json.Unmarshal(input.ToolInput, &ti); err != nil {
		return &HookOutput{}
	}

	if ti.FilePath == "" {
		return &HookOutput{}
	}

	// Check if the file is under spectr/changes/<id>/
	if !spectrChangesPattern.MatchString(ti.FilePath) {
		return &HookOutput{}
	}

	// Allow tasks.jsonc modifications
	if strings.HasSuffix(ti.FilePath, "/tasks.jsonc") {
		return &HookOutput{}
	}

	// Block all other modifications under spectr/changes/
	msg := "Blocked: cannot modify " + ti.FilePath +
		" during apply. Only tasks.jsonc may be modified under spectr/changes/."

	return &HookOutput{
		Blocked: true,
		Message: &msg,
	}
}
