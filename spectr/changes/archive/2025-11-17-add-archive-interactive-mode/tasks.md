# Implementation Tasks

## 1. Preparation

- [x] 1.1 Review internal/list/interactive.go to understand the table interface implementation
- [x] 1.2 Review cmd/list.go to understand how the -I flag is wired up
- [x] 1.3 Identify reusable components vs archive-specific customizations needed

## 2. Archive Command Updates

- [x] 2.1 Add `Interactive bool` field to ArchiveCmd struct in cmd/archive.go with appropriate tags
- [x] 2.2 Update ArchiveCmd.Run() to pass interactive flag to archiver
- [x] 2.3 Update archive.NewArchiver() signature to accept interactive flag
- [x] 2.4 Store interactive flag in Archiver struct

## 3. Interactive Selection Implementation

- [x] 3.1 Create RunInteractiveArchive function in internal/list/interactive.go (or internal/archive package)
- [x] 3.2 Adapt the table model to show changes with appropriate columns (ID, Title, Deltas, Tasks)
- [x] 3.3 Return selected change ID instead of copying to clipboard
- [x] 3.4 Handle empty change list gracefully
- [x] 3.5 Support keyboard navigation (arrow keys, j/k) and quit commands (q, Ctrl+C)

## 4. Archiver Integration

- [x] 4.1 Update Archiver.Archive() to check interactive flag when changeID is empty
- [x] 4.2 If interactive=true, call interactive table function instead of old selectChange
- [x] 4.3 If interactive=false and changeID empty, fall back to old text-based selection (or show error)
- [x] 4.4 Remove or deprecate old selectChange method once interactive is default

## 5. Testing

- [x] 5.1 Test `spectr archive` with no arguments (should show interactive table)
- [x] 5.2 Test `spectr archive -I` explicitly triggering interactive mode
- [x] 5.3 Test `spectr archive <change-id>` direct invocation still works
- [x] 5.4 Test keyboard navigation (arrow keys, j/k, enter, q, Ctrl+C)
- [x] 5.5 Test with no available changes (should show message and exit cleanly)
- [x] 5.6 Test that selected change ID properly flows into archive workflow
- [x] 5.7 Verify table columns display correct data (ID, title, delta count, task status)

## 6. Documentation

- [x] 6.1 Update help text for archive command to mention interactive mode
- [x] 6.2 Verify consistency with list command's -I flag behavior
- [x] 6.3 Update AGENTS.md if needed to document the interactive mode
