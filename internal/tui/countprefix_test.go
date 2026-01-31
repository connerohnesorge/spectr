package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestCountPrefixState_HandleKey(t *testing.T) {
	tests := []struct {
		name         string
		sequence     []string // sequence of key presses
		wantCount    int
		wantIsNavKey bool
		wantHandled  bool
		wantIsActive bool   // state after all keys
		wantString   string // String() output after all keys
	}{
		{
			name:         "basic count: 9j",
			sequence:     []string{"9", "j"},
			wantCount:    9,
			wantIsNavKey: true,
			wantHandled:  true,
			wantIsActive: false,
			wantString:   "",
		},
		{
			name:         "multi-digit: 42k",
			sequence:     []string{"4", "2", "k"},
			wantCount:    42,
			wantIsNavKey: true,
			wantHandled:  true,
			wantIsActive: false,
			wantString:   "",
		},
		{
			name:         "three-digit: 123j",
			sequence:     []string{"1", "2", "3", "j"},
			wantCount:    123,
			wantIsNavKey: true,
			wantHandled:  true,
			wantIsActive: false,
			wantString:   "",
		},
		{
			name:         "four-digit max: 9999j",
			sequence:     []string{"9", "9", "9", "9", "j"},
			wantCount:    9999,
			wantIsNavKey: true,
			wantHandled:  true,
			wantIsActive: false,
			wantString:   "",
		},
		{
			name:         "capping: 99999j (5 digits input, capped at 4)",
			sequence:     []string{"9", "9", "9", "9", "9", "j"},
			wantCount:    9999,
			wantIsNavKey: true,
			wantHandled:  true,
			wantIsActive: false,
			wantString:   "",
		},
		{
			name:         "zero count: 0j",
			sequence:     []string{"0", "j"},
			wantCount:    0,
			wantIsNavKey: true,
			wantHandled:  true,
			wantIsActive: false,
			wantString:   "",
		},
		{
			name:         "navigation key: j (no prefix)",
			sequence:     []string{"j"},
			wantCount:    1,
			wantIsNavKey: true,
			wantHandled:  true,
			wantIsActive: false,
			wantString:   "",
		},
		{
			name:         "navigation key: k (no prefix)",
			sequence:     []string{"k"},
			wantCount:    1,
			wantIsNavKey: true,
			wantHandled:  true,
			wantIsActive: false,
			wantString:   "",
		},
		{
			name:         "navigation key: up (no prefix)",
			sequence:     []string{"up"},
			wantCount:    1,
			wantIsNavKey: true,
			wantHandled:  true,
			wantIsActive: false,
			wantString:   "",
		},
		{
			name:         "navigation key: down (no prefix)",
			sequence:     []string{"down"},
			wantCount:    1,
			wantIsNavKey: true,
			wantHandled:  true,
			wantIsActive: false,
			wantString:   "",
		},
		{
			name:         "cancellation: 9<ESC>j",
			sequence:     []string{"9", "esc", "j"},
			wantCount:    1,
			wantIsNavKey: true,
			wantHandled:  true,
			wantIsActive: false,
			wantString:   "",
		},
		{
			name:         "cancellation: 42<ESC>k",
			sequence:     []string{"4", "2", "esc", "k"},
			wantCount:    1,
			wantIsNavKey: true,
			wantHandled:  true,
			wantIsActive: false,
			wantString:   "",
		},
		{
			name:         "other key resets: 5x (x is not nav key)",
			sequence:     []string{"5", "x"},
			wantCount:    1,
			wantIsNavKey: false,
			wantHandled:  true,
			wantIsActive: false,
			wantString:   "",
		},
		{
			name:         "case insensitive: 5J (uppercase J)",
			sequence:     []string{"5", "J"},
			wantCount:    5,
			wantIsNavKey: true,
			wantHandled:  true,
			wantIsActive: false,
			wantString:   "",
		},
		{
			name:         "case insensitive: 5K (uppercase K)",
			sequence:     []string{"5", "K"},
			wantCount:    5,
			wantIsNavKey: true,
			wantHandled:  true,
			wantIsActive: false,
			wantString:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := &CountPrefixState{}
			var count int
			var isNavKey, handled bool

			// Process all keys in sequence
			for i, key := range tt.sequence {
				msg := tea.KeyMsg{Type: tea.KeyRunes}
				// Handle special keys
				switch key {
				case keyEsc:
					msg.Type = tea.KeyEscape
				case keyUp:
					msg.Type = tea.KeyUp
				case keyDown:
					msg.Type = tea.KeyDown
				default:
					msg.Type = tea.KeyRunes
					msg.Runes = []rune(key)
				}

				count, isNavKey, handled = state.HandleKey(msg)

				// Only check final result on last key
				if i != len(tt.sequence)-1 {
					continue
				}

				if count != tt.wantCount {
					t.Errorf(
						"HandleKey() count = %d, want %d",
						count,
						tt.wantCount,
					)
				}
				if isNavKey != tt.wantIsNavKey {
					t.Errorf(
						"HandleKey() isNavKey = %v, want %v",
						isNavKey,
						tt.wantIsNavKey,
					)
				}
				if handled != tt.wantHandled {
					t.Errorf(
						"HandleKey() handled = %v, want %v",
						handled,
						tt.wantHandled,
					)
				}
			}

			// Check final state
			if state.IsActive() != tt.wantIsActive {
				t.Errorf(
					"IsActive() = %v, want %v",
					state.IsActive(),
					tt.wantIsActive,
				)
			}
			if state.String() != tt.wantString {
				t.Errorf(
					"String() = %q, want %q",
					state.String(),
					tt.wantString,
				)
			}
		})
	}
}

func TestCountPrefixState_IsActive(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
		want   bool
	}{
		{
			name:   "empty prefix",
			prefix: "",
			want:   false,
		},
		{
			name:   "single digit",
			prefix: "5",
			want:   true,
		},
		{
			name:   "multiple digits",
			prefix: "42",
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := &CountPrefixState{
				prefix: tt.prefix,
			}
			if got := state.IsActive(); got != tt.want {
				t.Errorf(
					"IsActive() = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestCountPrefixState_Reset(t *testing.T) {
	state := &CountPrefixState{prefix: "42"}

	if !state.IsActive() {
		t.Error("Expected state to be active before Reset()")
	}

	state.Reset()

	if state.IsActive() {
		t.Error("Expected state to be inactive after Reset()")
	}
	if state.String() != "" {
		t.Errorf(
			"Expected empty string after Reset(), got %q",
			state.String(),
		)
	}
}

func TestCountPrefixState_String(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
		want   string
	}{
		{
			name:   "empty",
			prefix: "",
			want:   "",
		},
		{
			name:   "single digit",
			prefix: "5",
			want:   "5",
		},
		{
			name:   "multiple digits",
			prefix: "123",
			want:   "123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := &CountPrefixState{
				prefix: tt.prefix,
			}
			if got := state.String(); got != tt.want {
				t.Errorf(
					"String() = %q, want %q",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestCountPrefixState_AccumulateDigits(t *testing.T) {
	state := &CountPrefixState{}

	// Accumulate digits one by one
	digits := []string{"1", "2", "3"}
	for _, digit := range digits {
		msg := tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune(digit),
		}
		count, isNavKey, handled := state.HandleKey(msg)
		if !handled {
			t.Error("Expected digit to be handled")
		}
		if isNavKey {
			t.Error("Expected digit to not be a nav key")
		}
		if count != 1 {
			t.Errorf("Expected count=1 for digit, got %d", count)
		}
	}

	if !state.IsActive() {
		t.Error("Expected state to be active after accumulating digits")
	}
	if state.String() != "123" {
		t.Errorf("Expected prefix '123', got %q", state.String())
	}
}

func TestIsNavigationKey(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want bool
	}{
		{name: "j lowercase", key: "j", want: true},
		{name: "J uppercase", key: "J", want: true},
		{name: "k lowercase", key: "k", want: true},
		{name: "K uppercase", key: "K", want: true},
		{name: "up", key: "up", want: true},
		{name: "UP uppercase", key: "UP", want: true},
		{name: "down", key: "down", want: true},
		{name: "DOWN uppercase", key: "DOWN", want: true},
		{name: "h is not nav", key: "h", want: false},
		{name: "l is not nav", key: "l", want: false},
		{name: "x is not nav", key: "x", want: false},
		{name: "1 is not nav", key: "1", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNavigationKey(tt.key); got != tt.want {
				t.Errorf(
					"isNavigationKey(%q) = %v, want %v",
					tt.key,
					got,
					tt.want,
				)
			}
		})
	}
}
