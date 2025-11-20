## ADDED Requirements

### Requirement: Archive Command PR Flag

The `spectr archive` command SHALL accept a `--pr` flag that triggers automated pull request creation after successful archive completion, integrating git branch creation, commit, push, and platform-specific PR CLI invocation into the archive workflow.

#### Scenario: User archives with PR flag

- **WHEN** user runs `spectr archive my-feature --pr`
- **AND** the change exists and archive operation succeeds
- **AND** git repository is configured with origin remote
- **THEN** the archive workflow completes (validation, spec merging, directory move)
- **AND** a git branch `archive-my-feature` is created
- **AND** archived files and updated specs are committed
- **AND** the branch is pushed to origin
- **AND** a PR is created using the detected platform CLI tool
- **AND** the PR URL is displayed

#### Scenario: PR flag combines with other archive flags

- **WHEN** user runs `spectr archive my-feature --pr --yes`
- **THEN** confirmation prompts are skipped
- **AND** PR creation proceeds automatically after archive

- **WHEN** user runs `spectr archive my-feature --pr --skip-specs`
- **THEN** spec updates are skipped
- **AND** PR is created with only archived directory
- **AND** PR body notes specs were skipped

- **WHEN** user runs `spectr archive my-feature --pr --interactive`
- **THEN** interactive change selection is used
- **AND** after selection, PR workflow proceeds

#### Scenario: PR flag without change ID uses interactive mode

- **WHEN** user runs `spectr archive --pr` with no change ID argument
- **THEN** interactive change selection is displayed
- **AND** after user selects a change, archive proceeds with PR creation
- **AND** the selected change is archived and PR is created

#### Scenario: PR flag error when archive fails

- **WHEN** user runs `spectr archive invalid --pr`
- **AND** the change does not exist or validation fails
- **THEN** an error is displayed about the archive failure
- **AND** no git operations are performed
- **AND** no branch or PR is created
- **AND** the command exits with error code 1

#### Scenario: PR flag error when git not available

- **WHEN** user runs `spectr archive my-feature --pr`
- **AND** archive succeeds but git is not available or no git repository exists
- **THEN** an error is displayed about git requirements
- **AND** the archive operation is complete
- **AND** changes remain uncommitted
- **AND** user is guided to manually commit or configure git
- **AND** the command exits with error code 1

#### Scenario: Help text documents PR flag

- **WHEN** user runs `spectr archive --help`
- **THEN** the `--pr` flag is listed in the available flags
- **AND** the description explains: "Create pull request after successful archive"
- **AND** the help text notes the dependency on git and PR CLI tools (gh/glab/tea)
