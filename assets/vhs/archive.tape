# Demo: spectr archive - Archiving Completed Changes
#
# This demo showcases:
# - Running spectr archive on a completed change
# - Merging deltas into base specs
# - Moving change to archive with timestamp
# - Verifying the archived state
#
# Execution Model:
# - Runs directly in examples/archive/spectr/ directory (no _demo copy)
# - The archive command WILL modify the directory (moves test-archive to changes/archive/)
# - Git restore commands are hidden at the end to reset the example directory
#
# Output: assets/gifs/archive.gif

Output assets/gifs/archive.gif

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

# Show the change in the list
Type "echo '=== Current active changes ==='"
Enter
Sleep 500ms
Type "cd examples/archive/spectr && spectr list"
Enter
Sleep 1.5s

# Run archive command
Type "echo '=== Archiving the completed change ==='"
Enter
Sleep 500ms
Type "cd examples/archive/spectr && spectr archive test-archive --yes"
Enter
Sleep 2s

# Verify it moved to archive
Type "echo '=== Verifying archive ==='"
Enter
Sleep 500ms
Type "ls -la examples/archive/spectr/changes/archive/ | tail -5"
Enter
Sleep 1.5s

# Show it's no longer in active changes
Type "echo '=== Active changes (should not include test-archive) ==='"
Enter
Sleep 500ms
Type "cd examples/archive/spectr && spectr list"
Enter
Sleep 1.5s

# Restore the example directory to original state (hidden from recording)
Hide
Type "cd examples/archive/spectr && git checkout -- ."
Enter
Sleep 500ms
Show

# Final pause
Sleep 1s
