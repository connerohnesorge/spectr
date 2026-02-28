package hooks

import (
	"encoding/json"
	"regexp"
	"strings"
)

// spectrChangesPattern matches references to spectr/changes/<id>/ in both file
// paths and shell command strings. It accepts spectr/changes/ when preceded by
// start-of-string or any non-alphanumeric/underscore character (e.g. /, space,
// >, &) so it catches redirect targets in Bash commands as well as plain paths.
var spectrChangesPattern = regexp.MustCompile(
	`(?:^|[^a-zA-Z0-9_])spectr/changes/[^/]+/`,
)

// fileToolNames lists the tool names that perform file write operations.
var fileToolNames = map[string]bool{
	"Edit":  true,
	"Write": true,
}

// TODO: Currently only the "apply" command has pre-tool-use guards.
// Future commands may need their own guards (e.g., "proposal" preventing
// writes to spectr/specs/, "next" restricting edits to task-related files).
// Consider a command→rules registry when adding more commands.

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

	// Best-effort guard for Bash commands that may write to protected paths.
	// Shell commands are inherently hard to analyze statically; this catches
	// common patterns like redirects and explicit paths in command strings.
	if input.ToolName == "Bash" {
		var ti toolInput
		if err := json.Unmarshal(input.ToolInput, &ti); err != nil {
			return &HookOutput{} // Can't parse Bash input; allow it
		}
		if ti.Command != "" && spectrChangesPattern.MatchString(ti.Command) {
			msg := "Blocked: Bash command references spectr/changes/ path during apply. " +
				"Use spectr CLI tools to modify change proposals."

			return &HookOutput{Blocked: true, Message: &msg}
		}

		return &HookOutput{}
	}

	// Only check file-writing tools
	if !fileToolNames[input.ToolName] {
		return &HookOutput{}
	}

	// Parse tool_input to get file_path
	var ti toolInput
	if err := json.Unmarshal(input.ToolInput, &ti); err != nil {
		msg := "Blocked: could not parse tool input for " + input.ToolName +
			" during apply"

		return &HookOutput{Blocked: true, Message: &msg}
	}

	if ti.FilePath == "" {
		msg := "Blocked: empty file path for " + input.ToolName +
			" during apply"

		return &HookOutput{Blocked: true, Message: &msg}
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
