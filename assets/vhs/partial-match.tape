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
Type "cd examples/partial-match/"
Enter
Sleep 1.5s

# Demonstrate prefix matching
Type "spectr archive refactor --yes"
Enter
Sleep 2s

# Restore for next demo (hidden)
Hide
Type "git checkout -- ."
Enter
Sleep 500ms
Show

# Demonstrate substring matching
Type "spectr archive unified --yes"
Enter
Sleep 2s

# Restore for next demo (hidden)
Hide
Type "git checkout -- ."
Enter
Sleep 500ms
Show

# Demonstrate case-insensitive matching
Type "spectr archive REFACTOR --yes"
Enter
Sleep 2s

# Final restore (hidden)
Hide
Type "git checkout -- ."
Enter
Sleep 500ms
Show

# Final pause
Sleep 1s
