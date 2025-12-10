## 1. Implementation

- [ ] 1.1 Update `internal/archive/archiver.go` to remove auto-accept logic
- [ ] 1.2 Add check for `tasks.json` existence before archiving
- [ ] 1.3 Display actionable error message: "No tasks.json found. Run `spectr accept <change-id>` to accept the change first."
- [ ] 1.4 Block archive when no task files exist at all

## 2. Validation

- [ ] 2.1 Test archiving a change with tasks.json (should succeed)
- [ ] 2.2 Test archiving a change with only tasks.md (should fail with helpful error)
- [ ] 2.3 Test archiving a change with no task files (should fail with helpful error)
