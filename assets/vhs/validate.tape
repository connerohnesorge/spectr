# Demo: spectr validate - Validation with Error Correction
#
# This demo showcases:
# - Running validation on a change proposal with errors
# - Viewing and understanding validation errors
# - Fixing errors in the spec file
# - Re-running validation successfully
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

# Copy base project and broken change
Type "# Set up demo with a broken change proposal"
Enter
Type "cp -r examples/validate/spectr _demo"
Enter
Type "mkdir -p _demo/changes/test-validation/specs/test-feature"
Enter
Type "cp examples/validate/broken/proposal.md _demo/changes/test-validation/"
Enter
Type "cp examples/validate/broken/specs/test-feature/spec.md _demo/changes/test-validation/specs/test-feature/"
Enter
Sleep 500ms

Type "echo ''"
Enter
Sleep 500ms

# Run validation - it will fail
Type "echo '=== Running validation (will show errors) ==='"
Enter
Sleep 500ms
Type "cd _demo && spectr validate test-validation --strict"
Enter
Sleep 2s

Type "echo ''"
Enter
Sleep 1s

# Fix the error by copying the fixed version
Type "echo '=== Fixing the spec by using corrected version ==='"
Enter
Sleep 500ms
Type "cp examples/validate/fixed/specs/test-feature/spec.md _demo/changes/test-validation/specs/test-feature/"
Enter
Sleep 500ms

Type "echo ''"
Enter
Sleep 500ms

# Run validation again - success
Type "echo '=== Re-running validation (should pass) ==='"
Enter
Sleep 500ms
Type "cd _demo && spectr validate test-validation --strict"
Enter
Sleep 2s

# Cleanup
Type "echo ''"
Enter
Sleep 500ms
Type "rm -rf _demo"
Enter
Sleep 500ms

# Final pause
Sleep 1s
