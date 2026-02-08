package domain

// HookType represents a type-safe hook event type identifier.
type HookType int

const (
	HookPreToolUse HookType = iota
	HookPostToolUse
	HookUserPromptSubmit
	HookStop
	HookSubagentStart
	HookSubagentStop
	HookPreCompact
	HookSessionStart
	HookSessionEnd
	HookNotification
	HookPermissionRequest
)

const unknownHookType = "unknown"

// String returns the PascalCase name for the hook type.
func (h HookType) String() string {
	names := []string{
		"PreToolUse",
		"PostToolUse",
		"UserPromptSubmit",
		"Stop",
		"SubagentStart",
		"SubagentStop",
		"PreCompact",
		"SessionStart",
		"SessionEnd",
		"Notification",
		"PermissionRequest",
	}
	if int(h) < len(names) {
		return names[h]
	}

	return unknownHookType
}

// AllHookTypes returns all defined hook types.
func AllHookTypes() []HookType {
	return []HookType{
		HookPreToolUse,
		HookPostToolUse,
		HookUserPromptSubmit,
		HookStop,
		HookSubagentStart,
		HookSubagentStop,
		HookPreCompact,
		HookSessionStart,
		HookSessionEnd,
		HookNotification,
		HookPermissionRequest,
	}
}

// ParseHookType parses a string into a HookType.
// Returns the HookType and true if found, or zero value and false if not.
func ParseHookType(s string) (HookType, bool) {
	for _, ht := range AllHookTypes() {
		if ht.String() == s {
			return ht, true
		}
	}

	return 0, false
}
