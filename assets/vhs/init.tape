# Demo: spectr init - Project Initialization Wizard
#
# This demo showcases:
# - Running the spectr init wizard
# - Creating a new Spectr project structure
# - Viewing the generated files
#
# Execution Model:
# - Creates a temporary directory (/tmp/spectr-demo) for the demo
# - Unlike other demos, we NEED to create a fresh directory since
#   we're demonstrating project initialization from scratch
# - Uses Hide/Show to hide cleanup commands at the end
#
# Output: assets/gifs/init.gif

Output assets/gifs/init.gif

# VHS Configuration
Set FontSize 14
Set Width 1200
Set Height 600
Set Padding 10
Set Theme "Catppuccin Mocha"
Set TypingSpeed 50ms

# Setup: Clear screen and show we're starting fresh
Type "clear"
Enter
Sleep 500ms

# Show we're in a fresh directory
Type "# Starting in a fresh directory"
Enter
Sleep 500ms

# Run spectr init with /tmp/spectr-demo directory
Type "spectr init /tmp/spectr-demo"
Enter
Sleep 1s

# Interact with wizard - Project name
Type "my-awesome-project"
Enter
Sleep 500ms

# Description
Type "A demonstration of spectr initialization"
Enter
Sleep 500ms

# Confirm creation
Type "y"
Enter
Sleep 1.5s

# Show the created structure
Type "tree /tmp/spectr-demo -L 2"
Enter
Sleep 2s

# Show the project.md file was created
Type "head -n 10 /tmp/spectr-demo/project.md"
Enter
Sleep 2s

# Hide cleanup commands from viewers
Hide

# Cleanup
Type "rm -rf /tmp/spectr-demo"
Enter
Sleep 500ms

# Show the screen again for final pause
Show

# Final pause to let viewers see the result
Sleep 1s
