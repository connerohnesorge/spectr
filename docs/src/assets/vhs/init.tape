Output docs/src/assets/gifs/init.gif

# VHS Configuration
Set FontSize 14
Set Width 1200
Set Height 600
Set Padding 10
Set Theme "Catppuccin Mocha"
Set TypingSpeed 50ms
Set Shell "zsh"

Hide
Type "rm -rf /tmp/spectr-demo"
Enter
Type "mkdir /tmp/spectr-demo"
Enter
Type "clear"
Enter
Show

# Run spectr init with /tmp/spectr-demo directory
Type "spectr init /tmp/spectr-demo"
Enter
Sleep 1s

# Interact with wizard - Project name
Enter
Sleep 500ms
Type "a"
Enter
Sleep 500ms
Enter
Sleep 500ms
Enter

# Show the created structure
Type "tree /tmp/spectr-demo -L 2"
Enter
Sleep 2s

# Show the project.md file was created
Type "cat /tmp/spectr-demo/project.md"
Enter
Sleep 2s

# Final pause to let viewers see the result
Sleep 1s
