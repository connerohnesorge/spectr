# Demo: spectr workflow - Complete End-to-End Workflow
#
# This demo showcases the complete Spectr workflow:
# 1. Create a new change proposal with spec deltas
# 2. Validate the proposal
# 3. View it in the changes list
# 4. Archive it after completion
#
# Output: assets/gifs/workflow.gif

Output assets/gifs/workflow.gif

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

Type "# Complete Spectr Workflow Demo"
Enter
Type "echo ''"
Enter
Sleep 500ms

# STAGE 1: Set up the change proposal (from pre-made example)
Type "echo '=== STAGE 1: Setting Up Change Proposal ==='"
Enter
Sleep 500ms

Type "# Using pre-made example project"
Enter
Type "cp -r examples/workflow/spectr _demo"
Enter
Sleep 500ms

Type "echo ''"
Enter
Sleep 500ms

# STAGE 2: Validate
Type "echo '=== STAGE 2: Validating Proposal ==='"
Enter
Sleep 500ms
Type "cd _demo && spectr validate add-user-notifications --strict"
Enter
Sleep 2s

Type "echo ''"
Enter
Sleep 500ms

# STAGE 3: View in list
Type "echo '=== STAGE 3: Viewing Active Changes ==='"
Enter
Sleep 500ms
Type "cd _demo && spectr list"
Enter
Sleep 1.5s

Type "echo ''"
Enter
Sleep 500ms

# Show details
Type "echo '=== Change Details ==='"
Enter
Sleep 500ms
Type "cd _demo && spectr show add-user-notifications"
Enter
Sleep 2s

Type "echo ''"
Enter
Sleep 500ms

# STAGE 4: Archive (simulating completion)
Type "echo '=== STAGE 4: Archiving Completed Change ==='"
Enter
Sleep 500ms
Type "echo '(After implementation is complete and deployed)'"
Enter
Sleep 1s
Type "cd _demo && spectr archive add-user-notifications --yes"
Enter
Sleep 2s

Type "echo ''"
Enter
Sleep 500ms

# Verify archive
Type "echo '=== Verification: Change Archived ==='"
Enter
Sleep 500ms
Type "cd _demo && spectr list | grep -q add-user-notifications && echo 'Still active' || echo 'Successfully archived'"
Enter
Sleep 1.5s

# Cleanup
Type "echo ''"
Enter
Sleep 500ms
Type "rm -rf _demo"
Enter
Sleep 500ms

Type "echo ''"
Enter
Type "echo 'Complete workflow demonstrated!'"
Enter
Sleep 1s

# Final pause
Sleep 1s
