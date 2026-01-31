// Package tui provides terminal UI components and utilities.
package tui

import (
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	// MaxCountPrefixDigits is the maximum number of digits allowed
	// in a count prefix (4 digits = max 9999).
	MaxCountPrefixDigits = 4
	// MaxCountPrefixValue is the maximum count value allowed.
	MaxCountPrefixValue = 9999
)

// CountPrefixState manages vim-style count prefix state (e.g., "42j").
// It accumulates digits and detects navigation keys to return the count.
type CountPrefixState struct {
	prefix string // accumulated digit string
}

// HandleKey processes a key press and updates count prefix state.
// Returns (count, isNavKey, handled):
//   - count: parsed count (1 if no prefix, capped at MaxCountPrefixValue)
//   - isNavKey: true if key is a navigation key (j/k/up/down)
//   - handled: true if this key was processed by count prefix logic
func (c *CountPrefixState) HandleKey(msg tea.KeyMsg) (
	count int,
	isNavKey, handled bool,
) {
	key := msg.String()

	// Check if it's a digit 0-9
	if len(key) == 1 && key[0] >= '0' && key[0] <= '9' {
		// Accumulate digit if we haven't reached max digits
		if len(c.prefix) < MaxCountPrefixDigits {
			c.prefix += key
		}

		return 1, false, true
	}

	// Check if it's a navigation key
	if isNavigationKey(key) {
		count := c.parseCount()
		c.Reset()

		return count, true, true
	}

	// Check if it's ESC key and we have an active prefix
	if key == keyEsc && c.IsActive() {
		c.Reset()

		return 1, false, true
	}

	// Any other key resets state if we have an active prefix
	// We return handled=true because we consumed the prefix
	if c.IsActive() {
		c.Reset()

		return 1, false, true
	}

	// No active prefix, key is not handled by count prefix
	return 1, false, false
}

// IsActive returns true if there are accumulated digits.
func (c *CountPrefixState) IsActive() bool {
	return c.prefix != ""
}

// Reset clears the accumulated prefix.
func (c *CountPrefixState) Reset() {
	c.prefix = ""
}

// String returns the current prefix string for display.
func (c *CountPrefixState) String() string {
	return c.prefix
}

// parseCount parses the accumulated prefix as an integer.
// Returns 1 if no prefix, or the parsed value capped at MaxCountPrefixValue.
func (c *CountPrefixState) parseCount() int {
	if c.prefix == "" {
		return 1
	}

	count, err := strconv.Atoi(c.prefix)
	if err != nil {
		return 1
	}

	// Cap at maximum value
	if count > MaxCountPrefixValue {
		return MaxCountPrefixValue
	}

	return count
}

// isNavigationKey checks if a key is a navigation key.
func isNavigationKey(key string) bool {
	lowerKey := strings.ToLower(key)

	return lowerKey == keyJ || lowerKey == keyK ||
		lowerKey == keyUp || lowerKey == keyDown
}
