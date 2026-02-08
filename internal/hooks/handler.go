package hooks

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/connerohnesorge/spectr/internal/domain"
)

// Handle reads hook input from stdin, dispatches to the appropriate handler,
// and writes the hook output to stdout.
func Handle(
	hookType domain.HookType,
	command string,
	stdin io.Reader,
	stdout io.Writer,
) error {
	var input HookInput
	if err := json.NewDecoder(stdin).Decode(&input); err != nil {
		return fmt.Errorf("failed to decode hook input: %w", err)
	}

	output := dispatch(hookType, command, &input)

	if err := json.NewEncoder(stdout).Encode(output); err != nil {
		return fmt.Errorf("failed to encode hook output: %w", err)
	}

	return nil
}

// dispatch routes hook events to their handlers.
// Most hook types are no-ops that return an empty (non-blocking) response.
func dispatch(
	hookType domain.HookType,
	command string,
	input *HookInput,
) *HookOutput {
	switch hookType {
	case domain.HookPreToolUse:
		return handlePreToolUse(command, input)
	case domain.HookPostToolUse,
		domain.HookUserPromptSubmit,
		domain.HookStop,
		domain.HookSubagentStart,
		domain.HookSubagentStop,
		domain.HookPreCompact,
		domain.HookSessionStart,
		domain.HookSessionEnd,
		domain.HookNotification,
		domain.HookPermissionRequest:
		return &HookOutput{}
	}

	return &HookOutput{}
}
