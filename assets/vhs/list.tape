Output assets/gifs/list.gif

# VHS Configuration
Set FontSize 14
Set Width 1200
Set Height 600
Set Padding 10
Set Theme "Catppuccin Mocha"
Set TypingSpeed 50ms
Set Shell "zsh"

# Setup: Clear screen
Type "clear"
Enter
Sleep 500ms

# List active changes
Type "cd examples/list/"
Enter
Sleep 1s

Type "spectr list"
Sleep 1s
Enter
Sleep 2s

# List specifications
Type "spectr list --specs"
Sleep 1s
Enter
Sleep 2s

# Show long format for more details
Type "spectr list --specs --long"
Sleep 1s
Enter
Sleep 2.5s

# Show interative mode
Type "spectr list -I"
Sleep 1s
Enter
Sleep 2s

# Final pause
Sleep 1s
