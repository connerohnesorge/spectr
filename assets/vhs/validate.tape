# Demo: spectr validate - Validation with Error Correction
#
# This demo showcases:
# - Running validation on a change proposal with errors
# - Viewing and understanding validation errors
# - Fixing errors in the spec file
# - Re-running validation successfully
#
# Execution Model:
# - Runs commands directly in examples/validate/spectr/ directory
# - Uses Hide/Show to hide setup commands that create the broken change
# - Shows the validation errors and the fix process
# - Uses Hide/Show to hide cleanup commands
#
# Output: assets/gifs/validate.gif

Output assets/gifs/validate.gif

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

# HIDDEN: Set up the broken change proposal
Hide
Type "mkdir -p examples/validate/spectr/changes/test-validation/specs/test-feature"
Enter
Type "cp examples/validate/broken/proposal.md examples/validate/spectr/changes/test-validation/"
Enter
Type "cp examples/validate/broken/specs/test-feature/spec.md examples/validate/spectr/changes/test-validation/specs/test-feature/"
Enter
Sleep 200ms
Show

# Run validation - it will fail
Type "cd examples/validate/spectr && spectr validate test-validation --strict"
Enter
Sleep 2s

Sleep 1s

# Fix the error by copying the fixed version
Type "cp examples/validate/fixed/specs/test-feature/spec.md examples/validate/spectr/changes/test-validation/specs/test-feature/"
Enter
Sleep 500ms

Sleep 500ms

# Run validation again - success
Type "cd examples/validate/spectr && spectr validate test-validation --strict"
Enter
Sleep 2s

# HIDDEN: Cleanup - restore the repository
Hide
Type "cd examples/validate/spectr && git restore ."
Enter
Sleep 200ms
Show

# Final pause
Sleep 1s
