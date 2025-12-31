## 1. Implementation

- [ ] 1.1 Add TTY check at the beginning of `Archive()` function in
  `internal/archive/archiver.go`
- [ ] 1.2 Use `isatty.IsTerminal(os.Stdout.Fd())` to detect TTY (consistent with
  existing pattern)
- [ ] 1.3 Return clear error message when TTY is not detected
- [ ] 1.4 Ensure TTY check happens before any file operations or validation

## 2. Testing

- [ ] 2.1 Write unit test for TTY detection in archive command
- [ ] 2.2 Test archive succeeds when TTY is present (manual testing)
- [ ] 2.3 Test archive fails gracefully when TTY is not present (e.g., `echo |
  spectr archive`)
- [ ] 2.4 Verify error message clarity and helpfulness

## 3. Documentation

- [ ] 3.1 Update archive command documentation to mention TTY requirement
- [ ] 3.2 Add note about human-only execution in relevant specs
