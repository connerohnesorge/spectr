Output assets/gifs/archive.gif

# VHS Configuration
Set FontSize 14
Set Width 1200
Set Height 600
Set Padding 10
Set Theme "Catppuccin Mocha"
Set TypingSpeed 50ms
Set Shell "zsh"

# Setup: Clean leftover artifacts from previous runs
Hide
Type "rm -rf examples/archive/spectr/changes/archive/"
Enter
Show

# Setup: Clear screen
Type "clear"
Enter
Sleep 500ms

# Show the change in the list
Type "cd examples/archive/"
Enter
Type "spectr list"
Sleep 1s
Enter
Sleep 1s

# Run archive command
Type "spectr archive test-archive --yes"
Sleep 1s
Enter
Sleep 2s

# Verify it moved to archive
Type "ls -la spectr/changes/archive/ | tail -5"
Sleep 1s
Enter
Sleep 1.5s

# Show it's no longer in active changes
Type "spectr list"
Sleep 1s
Enter
Sleep 1.5s

# Restore the example directory to original state (hidden from recording)
Hide
Type "git checkout -- . && rm -rf spectr/changes/archive/"
Enter
Show
