package domain

import "testing"

func TestSlashCommand_String(t *testing.T) {
	tests := []struct {
		name     string
		cmd      SlashCommand
		expected string
	}{
		{
			"Proposal command",
			SlashProposal,
			"proposal",
		},
		{"Apply command", SlashApply, "apply"},
		{
			"Unknown command",
			SlashCommand(999),
			"unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cmd.String()
			if got != tt.expected {
				t.Errorf(
					"SlashCommand(%d).String() = %q, want %q",
					tt.cmd,
					got,
					tt.expected,
				)
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
		{
			"Proposal template",
			SlashProposal,
			"slash-proposal.md.tmpl",
		},
		{
			"Apply template",
			SlashApply,
			"slash-apply.md.tmpl",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cmd.TemplateName()
			if err != nil {
				t.Errorf(
					"SlashCommand(%d).TemplateName() unexpected error: %v",
					tt.cmd,
					err,
				)
			}
			if got != tt.expected {
				t.Errorf(
					"SlashCommand(%d).TemplateName() = %q, want %q",
					tt.cmd,
					got,
					tt.expected,
				)
			}
		})
	}
}

func TestSlashCommand_TemplateNameUnknown(
	t *testing.T,
) {
	// Test that unknown command returns an error
	cmd := SlashCommand(999)
	got, err := cmd.TemplateName()
	if err == nil {
		t.Error(
			"SlashCommand(999).TemplateName() expected error, got nil",
		)
	}
	if got != "" {
		t.Errorf(
			"SlashCommand(999).TemplateName() = %q, want empty string on error",
			got,
		)
	}
}
