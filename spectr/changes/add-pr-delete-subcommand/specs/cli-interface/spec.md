## ADDED Requirements

### Requirement: PR Delete Subcommand

The `spectr pr delete` command (alias: `spectr pr d`) SHALL create a pull request that removes an entire spec directory from the codebase using git worktree isolation.

#### Scenario: User deletes a spec via PR

- **WHEN** user runs `spectr pr delete my-spec`
- **AND** `spectr/specs/my-spec/` exists
- **THEN** a worktree is created on branch `spectr/delete-my-spec`
- **AND** the directory `spectr/specs/my-spec/` is removed in the worktree
- **AND** a commit is created with message "spectr: delete my-spec spec"
- **AND** the branch is pushed to origin
- **AND** a PR is created via platform CLI (gh, glab, tea)
- **AND** the worktree is cleaned up
- **AND** the PR URL is displayed

#### Scenario: User uses shorthand alias

- **WHEN** user runs `spectr pr d my-spec`
- **THEN** the behavior is identical to `spectr pr delete my-spec`

#### Scenario: Partial spec ID resolution

- **WHEN** user runs `spectr pr delete cli`
- **AND** only one spec ID starts with or contains `cli` (e.g., `cli-interface`)
- **THEN** a message is displayed: "Resolved 'cli' -> 'cli-interface'"
- **AND** the delete proceeds with the resolved ID

#### Scenario: Ambiguous partial spec ID

- **WHEN** user runs `spectr pr delete val`
- **AND** multiple spec IDs match `val` (e.g., `validation`, `value-objects`)
- **THEN** an error is displayed: "Ambiguous ID 'val' matches multiple specs: validation, value-objects"
- **AND** the command exits with error code 1

#### Scenario: Spec not found

- **WHEN** user runs `spectr pr delete nonexistent`
- **AND** no spec matches `nonexistent`
- **THEN** an error is displayed: "No spec found matching 'nonexistent'"
- **AND** the command exits with error code 1

#### Scenario: Dry run mode

- **WHEN** user runs `spectr pr delete my-spec --dry-run`
- **THEN** the command displays what would happen without executing
- **AND** no worktree is created
- **AND** no files are deleted
- **AND** no PR is created

#### Scenario: Draft PR mode

- **WHEN** user runs `spectr pr delete my-spec --draft`
- **THEN** the PR is created as a draft PR

#### Scenario: Force delete existing branch

- **WHEN** user runs `spectr pr delete my-spec --force`
- **AND** branch `spectr/delete-my-spec` already exists on remote
- **THEN** the existing remote branch is deleted first
- **AND** the workflow proceeds normally

#### Scenario: Branch already exists without force

- **WHEN** user runs `spectr pr delete my-spec`
- **AND** branch `spectr/delete-my-spec` already exists on remote
- **AND** `--force` flag is not provided
- **THEN** an error is displayed: "branch 'spectr/delete-my-spec' already exists on remote; use --force to delete"
- **AND** the command exits with error code 1

#### Scenario: Custom base branch

- **WHEN** user runs `spectr pr delete my-spec --base develop`
- **THEN** the worktree is based on `develop` branch
- **AND** the PR targets `develop` instead of the default branch

### Requirement: Delete PR Commit Message Format

The commit message for spec deletion PRs SHALL follow a structured format that clearly identifies the deletion operation.

#### Scenario: Delete commit message format

- **WHEN** a delete PR commit is created
- **THEN** the commit message SHALL be: "spectr: delete <spec-id> spec"
- **AND** the commit body SHALL include the spec ID being deleted

### Requirement: Delete PR Body Format

The PR body for spec deletion PRs SHALL clearly explain the deletion operation.

#### Scenario: Delete PR body content

- **WHEN** a delete PR is created
- **THEN** the PR title SHALL be: "spectr: Delete <spec-id> spec"
- **AND** the PR body SHALL include a summary section explaining the deletion
- **AND** the PR body SHALL list the spec being removed
