package domain

// HookType represents a type-safe hook event type identifier.
type HookType int

// hookEntry pairs a HookType constant with its PascalCase name.
type hookEntry struct {
	Type HookType
	Name string
}

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

// hookTable is the single source of truth for every hook type.
// Order matches the iota declaration above.
var hookTable = []hookEntry{
	{HookPreToolUse, "PreToolUse"},
	{HookPostToolUse, "PostToolUse"},
	{HookUserPromptSubmit, "UserPromptSubmit"},
	{HookStop, "Stop"},
	{HookSubagentStart, "SubagentStart"},
	{HookSubagentStop, "SubagentStop"},
	{HookPreCompact, "PreCompact"},
	{HookSessionStart, "SessionStart"},
	{HookSessionEnd, "SessionEnd"},
	{HookNotification, "Notification"},
	{HookPermissionRequest, "PermissionRequest"},
}

const unknownHookType = "unknown"

// String returns the PascalCase name for the hook type.
func (h HookType) String() string {
	for _, e := range hookTable {
		if e.Type == h {
			return e.Name
		}
	}

	return unknownHookType
}

// AllHookTypes returns all defined hook types.
func AllHookTypes() []HookType {
	out := make([]HookType, len(hookTable))
	for i, e := range hookTable {
		out[i] = e.Type
	}

	return out
}

// ParseHookType parses a string into a HookType.
// Returns the HookType and true if found, or zero value and false if not.
func ParseHookType(s string) (HookType, bool) {
	for _, e := range hookTable {
		if e.Name == s {
			return e.Type, true
		}
	}

	return 0, false
}
