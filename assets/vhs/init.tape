# Demo: spectr init - Project Initialization Wizard
#
# This demo showcases:
# - Running the spectr init wizard
# - Creating a new Spectr project structure
# - Viewing the generated files
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

# Run spectr init with _demo directory
Type "spectr init _demo"
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
Type "echo ''"
Enter
Type "echo 'Generated project structure:'"
Enter
Type "tree _demo -L 2"
Enter
Sleep 2s

# Show the project.md file was created
Type "echo ''"
Enter
Type "echo 'Project configuration:'"
Enter
Type "head -n 10 _demo/project.md"
Enter
Sleep 2s

# Cleanup
Type "echo ''"
Enter
Sleep 500ms
Type "rm -rf _demo"
Enter
Sleep 500ms

# Final pause to let viewers see the result
Sleep 1s
