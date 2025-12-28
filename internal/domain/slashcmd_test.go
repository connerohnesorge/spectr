package domain

import "testing"

func TestSlashCommand_String(t *testing.T) {
	tests := []struct {
		cmd      SlashCommand
		expected string
	}{
		{SlashProposal, "proposal"},
		{SlashApply, "apply"},
	}

	for _, tt := range tests {
		got := tt.cmd.String()
		if got != tt.expected {
			t.Errorf("SlashCommand(%d).String() = %q, want %q", tt.cmd, got, tt.expected)
		}
	}
}

func TestSlashCommand_StringUnknown(t *testing.T) {
	// Test an out-of-range value
	unknownCmd := SlashCommand(999)
	got := unknownCmd.String()
	if got != "unknown" {
		t.Errorf("SlashCommand(999).String() = %q, want %q", got, "unknown")
	}
}

func TestSlashCommand_Iota(t *testing.T) {
	// Verify iota ordering
	if SlashProposal != 0 {
		t.Errorf("SlashProposal = %d, want 0", SlashProposal)
	}
	if SlashApply != 1 {
		t.Errorf("SlashApply = %d, want 1", SlashApply)
	}
}
