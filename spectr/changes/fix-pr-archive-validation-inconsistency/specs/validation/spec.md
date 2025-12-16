## ADDED Requirements

### Requirement: Base Branch Validation Mode
The validation system SHALL support validating change deltas against base specs from a specified git branch instead of local files.

#### Scenario: Validate with --base-branch flag
- **WHEN** the validate command is invoked with `--base-branch <ref>` flag
- **THEN** base specs SHALL be read from the specified git ref using `git show <ref>:spectr/specs/<capability>/spec.md`
- **AND** validation SHALL use the branch's spec content instead of local filesystem
- **AND** if the spec file doesn't exist on the branch, it SHALL be treated as non-existent (new capability)

#### Scenario: Git ref not found
- **WHEN** the validate command is invoked with `--base-branch` pointing to a non-existent ref
- **THEN** validation SHALL fail with a clear error message
- **AND** the error SHALL suggest checking if the branch exists and has been fetched

#### Scenario: Spec file read error from git
- **WHEN** git show fails to read a spec file for reasons other than "not found"
- **THEN** validation SHALL report the git error with the file path and ref
- **AND** validation SHALL NOT fall back to local specs silently

### Requirement: PR Archive Pre-flight Validation
The PR archive workflow SHALL perform pre-flight validation against the target branch before creating a worktree, providing faster feedback with better error context.

#### Scenario: Pre-flight validation catches ADDED conflict
- **WHEN** a change adds a requirement that already exists in the target branch's base spec
- **THEN** pre-flight validation SHALL fail before worktree creation
- **AND** the error message SHALL explain that the requirement exists on the target branch
- **AND** the error SHALL suggest running `spectr validate <change> --base-branch <target>` for debugging

#### Scenario: Pre-flight validation passes
- **WHEN** a change's deltas are valid against the target branch's specs
- **THEN** pre-flight validation SHALL pass
- **AND** the PR workflow SHALL proceed to create the worktree
- **AND** no duplicate validation error SHALL occur in the worktree

#### Scenario: Skip pre-flight when --no-validate specified
- **WHEN** pr archive is invoked with `--no-validate` flag
- **THEN** pre-flight validation SHALL be skipped
- **AND** the worktree SHALL be created directly

### Requirement: Improved Validation Error Messages for Branch Discrepancy
The validation system SHALL provide clear error messages when validation fails due to differences between local and remote specs.

#### Scenario: Error message for ADDED requirement conflict
- **WHEN** validation fails because an ADDED requirement exists in the base spec
- **AND** validation is running in base-branch mode
- **THEN** the error message SHALL indicate: "ADDED requirement '{name}' already exists in base spec on branch '{branch}'"
- **AND** the error SHALL suggest: "This may indicate the requirement was merged in another PR. Consider using MODIFIED instead of ADDED."

#### Scenario: Diff suggestion on conflict
- **WHEN** validation fails due to local vs branch spec differences
- **THEN** the error message SHALL suggest: "Run `spectr validate <change> --base-branch <branch>` to validate against the target branch"
