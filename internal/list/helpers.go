package list

import (
	"github.com/connerohnesorge/spectr/internal/tui"
)

// truncateString truncates a string and adds ellipsis if needed.
// This is a thin wrapper around tui.TruncateString for backward compatibility.
func truncateString(s string, maxLen int) string {
	return tui.TruncateString(s, maxLen)
}
