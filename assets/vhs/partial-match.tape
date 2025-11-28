# Demo: spectr archive - Partial ID Matching
#
# This demo showcases:
# - Partial ID matching with prefix (e.g., "refactor" -> "refactor-unified-interactive-tui")
# - Partial ID matching with substring (e.g., "unified" -> "refactor-unified-interactive-tui")
# - Case-insensitive matching
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

# Copy example project to temp workspace
Type "# Set up demo environment from pre-made example"
Enter
Type "cp -r examples/partial-match/spectr _demo"
Enter
Sleep 500ms

Type "echo ''"
Enter
Sleep 500ms

# Show the change in the list (full name visible)
Type "echo '=== Active change with long ID ==='"
Enter
Sleep 500ms
Type "cd _demo && spectr list"
Enter
Sleep 1.5s

Type "echo ''"
Enter
Sleep 500ms

# Demonstrate prefix matching
Type "echo '=== Archive using prefix: refactor ==='"
Enter
Sleep 500ms
Type "cd _demo && spectr archive refactor --yes"
Enter
Sleep 2s

Type "echo ''"
Enter
Sleep 500ms

# Restore for next demo
Type "# Reset demo for substring example"
Enter
Type "rm -rf _demo && cp -r examples/partial-match/spectr _demo"
Enter
Sleep 500ms

Type "echo ''"
Enter
Sleep 500ms

# Demonstrate substring matching
Type "echo '=== Archive using substring: unified ==='"
Enter
Sleep 500ms
Type "cd _demo && spectr archive unified --yes"
Enter
Sleep 2s

Type "echo ''"
Enter
Sleep 500ms

# Restore for next demo
Type "# Reset demo for case-insensitive example"
Enter
Type "rm -rf _demo && cp -r examples/partial-match/spectr _demo"
Enter
Sleep 500ms

Type "echo ''"
Enter
Sleep 500ms

# Demonstrate case-insensitive matching
Type "echo '=== Case-insensitive: REFACTOR ==='"
Enter
Sleep 500ms
Type "cd _demo && spectr archive REFACTOR --yes"
Enter
Sleep 2s

# Cleanup
Type "echo ''"
Enter
Sleep 500ms
Type "rm -rf _demo"
Enter
Sleep 500ms

# Final pause
Sleep 1s
