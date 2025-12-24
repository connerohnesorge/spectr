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

func TestSlashCommand_TemplateRef(t *testing.T) {
	tests := []struct {
		name         string
		cmd          SlashCommand
		expectedName string
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
			ref, err := tt.cmd.TemplateRef()
			if err != nil {
				t.Errorf(
					"SlashCommand(%d).TemplateRef() unexpected error: %v",
					tt.cmd,
					err,
				)
			}
			if ref.Name != tt.expectedName {
				t.Errorf(
					"SlashCommand(%d).TemplateRef().Name = %q, want %q",
					tt.cmd,
					ref.Name,
					tt.expectedName,
				)
			}
			if ref.Template == nil {
				t.Error(
					"SlashCommand.TemplateRef().Template should not be nil",
				)
			}
		})
	}
}

func TestSlashCommand_TemplateRefUnknown(
	t *testing.T,
) {
	// Test that unknown command returns an error
	cmd := SlashCommand(999)
	ref, err := cmd.TemplateRef()
	if err == nil {
		t.Error(
			"SlashCommand(999).TemplateRef() expected error, got nil",
		)
	}
	if ref.Template != nil {
		t.Error(
			"SlashCommand(999).TemplateRef().Template should be nil on error",
		)
	}
}
