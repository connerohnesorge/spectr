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

# STAGE 1: Create change proposal
Type "echo '=== STAGE 1: Creating Change Proposal ==='"
Enter
Sleep 500ms

Type "mkdir -p spectr/changes/add-user-notifications/specs/notifications"
Enter
Sleep 300ms

# Create proposal
Type "cat > spectr/changes/add-user-notifications/proposal.md << 'EOF'"
Enter
Type "# Change: Add User Notifications"
Enter
Type ""
Enter
Type "## Why"
Enter
Type "Users need to be notified of important system events."
Enter
Type ""
Enter
Type "## What Changes"
Enter
Type "- Add notification capability"
Enter
Type "- Support email and in-app notifications"
Enter
Type ""
Enter
Type "## Impact"
Enter
Type "- Affected specs: notifications"
Enter
Type "- New feature, no breaking changes"
Enter
Type "EOF"
Enter
Sleep 500ms

# Create spec delta
Type "cat > spectr/changes/add-user-notifications/specs/notifications/spec.md << 'EOF'"
Enter
Type "## ADDED Requirements"
Enter
Type "### Requirement: Email Notifications"
Enter
Type "The system SHALL send email notifications for critical events."
Enter
Type ""
Enter
Type "#### Scenario: Critical event notification"
Enter
Type "- **WHEN** critical event occurs"
Enter
Type "- **THEN** email is sent to user"
Enter
Type "- **AND** notification is logged"
Enter
Type ""
Enter
Type "### Requirement: In-App Notifications"
Enter
Type "The system SHALL display in-app notifications."
Enter
Type ""
Enter
Type "#### Scenario: Real-time notification"
Enter
Type "- **WHEN** user is online"
Enter
Type "- **THEN** notification appears in real-time"
Enter
Type "EOF"
Enter
Sleep 500ms

# Create tasks file
Type "cat > spectr/changes/add-user-notifications/tasks.md << 'EOF'"
Enter
Type "## 1. Implementation"
Enter
Type "- [ ] 1.1 Create notification service"
Enter
Type "- [ ] 1.2 Implement email delivery"
Enter
Type "- [ ] 1.3 Add in-app notification UI"
Enter
Type "- [ ] 1.4 Write tests"
Enter
Type "EOF"
Enter
Sleep 500ms

# Create base spec
Type "mkdir -p spectr/specs/notifications"
Enter
Type "touch spectr/specs/notifications/spec.md"
Enter
Sleep 500ms

Type "echo ''"
Enter
Sleep 500ms

# STAGE 2: Validate
Type "echo '=== STAGE 2: Validating Proposal ==='"
Enter
Sleep 500ms
Type "spectr validate add-user-notifications --strict"
Enter
Sleep 2s

Type "echo ''"
Enter
Sleep 500ms

# STAGE 3: View in list
Type "echo '=== STAGE 3: Viewing Active Changes ==='"
Enter
Sleep 500ms
Type "spectr list"
Enter
Sleep 1.5s

Type "echo ''"
Enter
Sleep 500ms

# Show details
Type "echo '=== Change Details ==='"
Enter
Sleep 500ms
Type "spectr show add-user-notifications"
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
Type "spectr archive add-user-notifications --yes"
Enter
Sleep 2s

Type "echo ''"
Enter
Sleep 500ms

# Verify archive
Type "echo '=== Verification: Change Archived ==='"
Enter
Sleep 500ms
Type "spectr list | grep -q add-user-notifications && echo 'Still active' || echo 'Successfully archived'"
Enter
Sleep 1.5s

# Cleanup
Type "echo ''"
Enter
Sleep 500ms
Type "rm -rf spectr/changes/archive/*/add-user-notifications"
Enter
Type "rm -rf spectr/specs/notifications"
Enter
Sleep 500ms

Type "echo ''"
Enter
Type "echo 'Complete workflow demonstrated!'"
Enter
Sleep 1s

# Final pause
Sleep 1s
