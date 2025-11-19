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

# Create a test change with intentional errors
Type "# Create a test change proposal with errors"
Enter
Type "mkdir -p spectr/changes/test-validation/specs/test-feature"
Enter
Sleep 500ms

# Create a spec with missing scenarios
Type "cat > spectr/changes/test-validation/specs/test-feature/spec.md << 'EOF'"
Enter
Type "## ADDED Requirements"
Enter
Type "### Requirement: User Authentication"
Enter
Type "The system SHALL provide secure authentication."
Enter
Type ""
Enter
Type "### Requirement: Missing Scenario Example"
Enter
Type "This requirement is missing a scenario."
Enter
Type "EOF"
Enter
Sleep 500ms

# Create minimal proposal
Type "cat > spectr/changes/test-validation/proposal.md << 'EOF'"
Enter
Type "# Change: Test Validation"
Enter
Type ""
Enter
Type "## Why"
Enter
Type "Testing validation workflow."
Enter
Type "EOF"
Enter
Sleep 500ms

Type "echo ''"
Enter
Sleep 500ms

# Run validation - it will fail
Type "echo '=== Running validation (will show errors) ==='"
Enter
Sleep 500ms
Type "spectr validate test-validation --strict"
Enter
Sleep 2s

Type "echo ''"
Enter
Sleep 1s

# Fix the error by adding scenarios
Type "echo '=== Fixing the spec by adding missing scenario ==='"
Enter
Sleep 500ms
Type "cat > spectr/changes/test-validation/specs/test-feature/spec.md << 'EOF'"
Enter
Type "## ADDED Requirements"
Enter
Type "### Requirement: User Authentication"
Enter
Type "The system SHALL provide secure authentication."
Enter
Type ""
Enter
Type "#### Scenario: Successful login"
Enter
Type "- **WHEN** user provides valid credentials"
Enter
Type "- **THEN** system grants access token"
Enter
Type ""
Enter
Type "### Requirement: Missing Scenario Example"
Enter
Type "This requirement now has a scenario."
Enter
Type ""
Enter
Type "#### Scenario: Example scenario"
Enter
Type "- **WHEN** condition occurs"
Enter
Type "- **THEN** expected result happens"
Enter
Type "EOF"
Enter
Sleep 500ms

Type "echo ''"
Enter
Sleep 500ms

# Run validation again - success
Type "echo '=== Re-running validation (should pass) ==='"
Enter
Sleep 500ms
Type "spectr validate test-validation --strict"
Enter
Sleep 2s

# Cleanup
Type "echo ''"
Enter
Sleep 500ms
Type "rm -rf spectr/changes/test-validation"
Enter
Sleep 500ms

# Final pause
Sleep 1s
