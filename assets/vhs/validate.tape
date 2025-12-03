Output assets/gifs/validate.gif

# VHS Configuration
Set FontSize 14
Set Width 1200
Set Height 600
Set Padding 10
Set Theme "Catppuccin Mocha"
Set TypingSpeed 50ms
Set Shell "zsh"

Hide
Type "mkdir -p examples/validate/spectr/changes/test-validation/specs/test-feature"
Enter
Type "cp examples/validate/broken/proposal.md examples/validate/spectr/changes/test-validation/"
Enter
Type "cp examples/validate/broken/specs/test-feature/spec.md examples/validate/spectr/changes/test-validation/specs/test-feature/"
Enter
Type "rm -rf examples/validate/spectr/changes/"
Enter
Type "clear"
Enter
Show

Type "cd examples/validate"
Enter
Sleep 1s

Type "spectr validate test-validation --strict"
Enter
Sleep 1s

Type "cp fixed/specs/test-feature/spec.md spectr/changes/test-validation/specs/test-feature/"
Enter
Sleep 500ms

Sleep 500ms

# Run validation again - success
Type "spectr validate test-validation --strict"
Enter
Sleep 2s

# HIDDEN: Cleanup - restore the repository
Hide
Type "git restore ."
Enter
Sleep 200ms
Show

# Final pause
Sleep 1s
