# Demo: spectr archive - Archiving Completed Changes
#
# This demo showcases:
# - Running spectr archive on a completed change
# - Merging deltas into base specs
# - Moving change to archive with timestamp
# - Verifying the archived state
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

# Copy example project to temp workspace
Type "# Set up demo environment from pre-made example"
Enter
Type "cp -r examples/archive/spectr _demo"
Enter
Sleep 500ms

Type "echo ''"
Enter
Sleep 500ms

# Show the change in the list
Type "echo '=== Current active changes ==='"
Enter
Sleep 500ms
Type "cd _demo && spectr list"
Enter
Sleep 1.5s

Type "echo ''"
Enter
Sleep 500ms

# Run archive command
Type "echo '=== Archiving the completed change ==='"
Enter
Sleep 500ms
Type "cd _demo && spectr archive test-archive --yes"
Enter
Sleep 2s

Type "echo ''"
Enter
Sleep 500ms

# Verify it moved to archive
Type "echo '=== Verifying archive ==='"
Enter
Sleep 500ms
Type "ls -la _demo/changes/archive/ | tail -5"
Enter
Sleep 1.5s

Type "echo ''"
Enter
Sleep 500ms

# Show it's no longer in active changes
Type "echo '=== Active changes (should not include test-archive) ==='"
Enter
Sleep 500ms
Type "cd _demo && spectr list"
Enter
Sleep 1.5s

# Cleanup
Type "echo ''"
Enter
Sleep 500ms
Type "rm -rf _demo"
Enter
Sleep 500ms

# Final pause
Sleep 1s
