package domain

import "testing"

func TestSlashCommand_String(t *testing.T) {
	tests := []struct {
		name     string
		cmd      SlashCommand
		expected string
	}{
		{"Proposal command", SlashProposal, "proposal"},
		{"Apply command", SlashApply, "apply"},
		{"Unknown command", SlashCommand(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cmd.String()
			if got != tt.expected {
				t.Errorf("SlashCommand(%d).String() = %q, want %q", tt.cmd, got, tt.expected)
			}
		})
	}
}

func TestSlashCommand_TemplateName(t *testing.T) {
	tests := []struct {
		name     string
		cmd      SlashCommand
		expected string
	}{
		{"Proposal template", SlashProposal, "slash-proposal.md.tmpl"},
		{"Apply template", SlashApply, "slash-apply.md.tmpl"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cmd.TemplateName()
			if got != tt.expected {
				t.Errorf("SlashCommand(%d).TemplateName() = %q, want %q", tt.cmd, got, tt.expected)
			}
		})
	}
}

func TestSlashCommand_TemplateNameUnknown(t *testing.T) {
	// Test that unknown command returns empty string (zero value of map lookup)
	cmd := SlashCommand(999)
	got := cmd.TemplateName()
	if got != "" {
		t.Errorf("SlashCommand(999).TemplateName() = %q, want empty string", got)
	}
}
