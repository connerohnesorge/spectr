Output assets/gifs/pr-hotkey.gif

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

# Navigate to example directory
Type "cd examples/list/"
Enter
Sleep 1s

# Open interactive list mode
Type "spectr list -I"
Sleep 1s
Enter
Sleep 2s

# Show help menu to display available hotkeys
Type "?"
Sleep 2s

# Hide help menu
Type "?"
Sleep 1.5s

# Press Shift+P to trigger PR mode
Type "P"
Sleep 2s

# Exit the TUI
Type "q"
Sleep 500ms
