package domain

import "testing"

func TestSlashCommand_String(t *testing.T) {
	tests := []struct {
		name string
		cmd  SlashCommand
		want string
	}{
		{
			name: "SlashProposal",
			cmd:  SlashProposal,
			want: "proposal",
		},
		{
			name: "SlashApply",
			cmd:  SlashApply,
			want: "apply",
		},
		{
			name: "invalid command",
			cmd:  SlashCommand(99),
			want: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cmd.String(); got != tt.want {
				t.Errorf("SlashCommand.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
