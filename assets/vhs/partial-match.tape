# Demo: spectr archive - Partial ID Matching
#
# This demo showcases:
# - Partial ID matching with prefix (e.g., "refactor" -> "refactor-unified-interactive-tui")
# - Partial ID matching with substring (e.g., "unified" -> "refactor-unified-interactive-tui")
# - Case-insensitive matching
#
# Execution Model:
# - All commands run directly in examples/partial-match/spectr/ directory
# - Git restore commands are hidden using VHS Hide/Show directives
# - No temporary _demo directory needed - uses the real example in place
# - After each archive, git restore resets the state for the next demo
#
# Output: assets/gifs/partial-match.gif

Output assets/gifs/partial-match.gif

# VHS Configuration
Set FontSize 14
Set Width 1200
Set Height 600
Set Padding 10
Set Theme "Catppuccin Mocha"
Set TypingSpeed 50ms

# Setup: Clear screen
Type "clear"
Enter
Sleep 500ms

# Show the change in the list (full name visible)
Type "echo '=== Active change with long ID ==='"
Enter
Sleep 500ms
Type "cd examples/partial-match/spectr && spectr list"
Enter
Sleep 1.5s

# Demonstrate prefix matching
Type "echo '=== Archive using prefix: refactor ==='"
Enter
Sleep 500ms
Type "cd examples/partial-match/spectr && spectr archive refactor --yes"
Enter
Sleep 2s

# Restore for next demo (hidden)
Hide
Type "cd examples/partial-match/spectr && git checkout -- ."
Enter
Sleep 500ms
Show

# Demonstrate substring matching
Type "echo '=== Archive using substring: unified ==='"
Enter
Sleep 500ms
Type "cd examples/partial-match/spectr && spectr archive unified --yes"
Enter
Sleep 2s

# Restore for next demo (hidden)
Hide
Type "cd examples/partial-match/spectr && git checkout -- ."
Enter
Sleep 500ms
Show

# Demonstrate case-insensitive matching
Type "echo '=== Case-insensitive: REFACTOR ==='"
Enter
Sleep 500ms
Type "cd examples/partial-match/spectr && spectr archive REFACTOR --yes"
Enter
Sleep 2s

# Final restore (hidden)
Hide
Type "cd examples/partial-match/spectr && git checkout -- ."
Enter
Sleep 500ms
Show

# Final pause
Sleep 1s
