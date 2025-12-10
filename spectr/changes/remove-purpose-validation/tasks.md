## 1. Spec Updates

- [ ] 1.1 Modify validation spec to remove Purpose-related scenarios

## 2. Code Implementation

- [ ] 2.1 Remove Purpose validation from `internal/validation/spec_rules.go`
- [ ] 2.2 Remove `MinPurposeLength` constant
- [ ] 2.3 Update `internal/archive/spec_merger.go` to remove Purpose from skeleton
- [ ] 2.4 Update test files that reference Purpose validation

## 3. Documentation

- [ ] 3.1 Update `spectr/project.md` to remove Purpose length constraint
- [ ] 3.2 Update `spectr/AGENTS.md` if it references Purpose requirements

## 4. Validation

- [ ] 4.1 Run `spectr validate --all --strict` to ensure no regressions
- [ ] 4.2 Run Go tests to verify validation logic still works
