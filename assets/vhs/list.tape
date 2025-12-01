# Demo: spectr list - Listing Changes and Specifications
#
# This demo showcases:
# - Listing active changes with spectr list
# - Listing specifications with spectr list --specs
# - Understanding what's in progress vs what's built
#
# Execution model:
# - Runs commands directly in examples/list/spectr/ directory
# - Read-only operations (no setup or cleanup needed)
# - No git restore needed (list commands don't modify state)
#
# Output: assets/gifs/list.gif

Output assets/gifs/list.gif

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

# List active changes
Type "echo '=== Active Changes (What is being proposed) ==='"
Enter
Sleep 500ms
Type "cd examples/list/spectr && spectr list"
Enter
Sleep 2s

# List specifications
Type "echo '=== Specifications (What is built) ==='"
Enter
Sleep 500ms
Type "cd examples/list/spectr && spectr list --specs"
Enter
Sleep 2s

# Show long format for more details
Type "echo '=== Detailed Spec View ==='"
Enter
Sleep 500ms
Type "cd examples/list/spectr && spectr list --specs --long"
Enter
Sleep 2.5s

# Final pause
Sleep 1s
