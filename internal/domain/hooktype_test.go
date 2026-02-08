package domain

import (
	"testing"
)

func TestHookTypeString(t *testing.T) {
	tests := []struct {
		hookType HookType
		want     string
	}{
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

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.hookType.String(); got != tt.want {
				t.Errorf(
					"HookType.String() = %q, want %q",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestHookTypeString_Unknown(t *testing.T) {
	unknown := HookType(999)
	if got := unknown.String(); got != unknownHookType {
		t.Errorf(
			"HookType(999).String() = %q, want %q",
			got,
			unknownHookType,
		)
	}
}

func TestAllHookTypes(t *testing.T) {
	types := AllHookTypes()
	if len(types) != 11 {
		t.Errorf(
			"AllHookTypes() returned %d types, want 11",
			len(types),
		)
	}
}

func TestParseHookType(t *testing.T) {
	tests := []struct {
		input string
		want  HookType
		ok    bool
	}{
		{"PreToolUse", HookPreToolUse, true},
		{"PostToolUse", HookPostToolUse, true},
		{"UserPromptSubmit", HookUserPromptSubmit, true},
		{"Stop", HookStop, true},
		{"SubagentStart", HookSubagentStart, true},
		{"SubagentStop", HookSubagentStop, true},
		{"PreCompact", HookPreCompact, true},
		{"SessionStart", HookSessionStart, true},
		{"SessionEnd", HookSessionEnd, true},
		{"Notification", HookNotification, true},
		{"PermissionRequest", HookPermissionRequest, true},
		{"invalid", 0, false},
		{"pretooluse", 0, false},
		{"", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, ok := ParseHookType(tt.input)
			if ok != tt.ok {
				t.Errorf(
					"ParseHookType(%q) ok = %v, want %v",
					tt.input,
					ok,
					tt.ok,
				)
			}
			if got != tt.want {
				t.Errorf(
					"ParseHookType(%q) = %v, want %v",
					tt.input,
					got,
					tt.want,
				)
			}
		})
	}
}
