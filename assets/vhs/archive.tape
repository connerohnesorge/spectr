# Demo: spectr archive - Archiving Completed Changes
#
# This demo showcases:
# - Running spectr archive on a completed change
# - Merging deltas into base specs
# - Moving change to archive with timestamp
# - Verifying the archived state
#
# Output: assets/gifs/archive.gif

Output assets/gifs/archive.gif

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

# Create a test change to archive
Type "# Create a completed change ready to archive"
Enter
Type "mkdir -p spectr/changes/test-archive/specs/sample-feature"
Enter
Sleep 500ms

# Create proposal
Type "cat > spectr/changes/test-archive/proposal.md << 'EOF'"
Enter
Type "# Change: Test Archive"
Enter
Type ""
Enter
Type "## Why"
Enter
Type "Demonstrating archive workflow."
Enter
Type ""
Enter
Type "## What Changes"
Enter
Type "- Add sample feature capability"
Enter
Type "EOF"
Enter
Sleep 500ms

# Create spec delta
Type "cat > spectr/changes/test-archive/specs/sample-feature/spec.md << 'EOF'"
Enter
Type "## ADDED Requirements"
Enter
Type "### Requirement: Sample Feature"
Enter
Type "The system SHALL provide sample functionality."
Enter
Type ""
Enter
Type "#### Scenario: Basic usage"
Enter
Type "- **WHEN** user invokes feature"
Enter
Type "- **THEN** system responds correctly"
Enter
Type "EOF"
Enter
Sleep 500ms

# Create base spec directory
Type "mkdir -p spectr/specs/sample-feature"
Enter
Type "touch spectr/specs/sample-feature/spec.md"
Enter
Sleep 500ms

Type "echo ''"
Enter
Sleep 500ms

# Show the change in the list
Type "echo '=== Current active changes ==='"
Enter
Sleep 500ms
Type "spectr list"
Enter
Sleep 1.5s

Type "echo ''"
Enter
Sleep 500ms

# Run archive command
Type "echo '=== Archiving the completed change ==='"
Enter
Sleep 500ms
Type "spectr archive test-archive --yes"
Enter
Sleep 2s

Type "echo ''"
Enter
Sleep 500ms

# Verify it moved to archive
Type "echo '=== Verifying archive ==='"
Enter
Sleep 500ms
Type "ls -la spectr/changes/archive/ | tail -5"
Enter
Sleep 1.5s

Type "echo ''"
Enter
Sleep 500ms

# Show it's no longer in active changes
Type "echo '=== Active changes (should not include test-archive) ==='"
Enter
Sleep 500ms
Type "spectr list"
Enter
Sleep 1.5s

# Cleanup
Type "echo ''"
Enter
Sleep 500ms
Type "rm -rf spectr/changes/archive/*/test-archive"
Enter
Type "rm -rf spectr/specs/sample-feature"
Enter
Sleep 500ms

# Final pause
Sleep 1s
