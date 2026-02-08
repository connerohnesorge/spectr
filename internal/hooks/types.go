// Package hooks implements hook handling for Claude Code slash commands.
package hooks

import "encoding/json"

// HookInput represents the JSON payload received from Claude Code via stdin.
type HookInput struct {
	SessionID     string          `json:"session_id"`
	HookEventName string          `json:"hook_event_name"`
	ToolName      string          `json:"tool_name"`
	ToolInput     json.RawMessage `json:"tool_input"`
}

// HookOutput represents the JSON response sent to Claude Code via stdout.
type HookOutput struct {
	Blocked      bool    `json:"blocked"`
	Message      *string `json:"message"`
	SystemPrompt *string `json:"system_prompt"`
}

// toolInput represents parsed tool_input fields relevant to file operations.
type toolInput struct {
	FilePath string `json:"file_path"`
}
