# Delta Specification

## ADDED Requirements

### Requirement: VHS GIF Auto-Generation Workflow

The system SHALL provide automated GIF generation through a GitHub Actions
workflow that triggers when VHS tape files are modified and commits updated GIFs
back to the repository.

#### Scenario: GIF regeneration on tape file change

- **WHEN** a developer pushes changes to any `.tape` file in `assets/vhs/`
- **THEN** the VHS workflow executes automatically
- **AND** VHS processes all tape files to generate updated GIFs
- **AND** the generated GIFs are committed back to the `assets/gifs/` directory

#### Scenario: Workflow ignores non-tape changes

- **WHEN** a developer pushes changes that do not modify any `.tape` file
- **THEN** the VHS workflow does NOT trigger
- **AND** existing GIFs remain unchanged

#### Scenario: Multiple tape files processed

- **WHEN** multiple tape files are modified in a single commit
- **THEN** all modified tape files are processed in a single workflow run
- **AND** all corresponding GIFs are generated and committed together

### Requirement: VHS Action Version Pinning

The system SHALL use a specific major version tag of the vhs-action (e.g.,
`@v2`) to ensure reproducible builds while still receiving compatible updates.

#### Scenario: Version-pinned VHS action reference

- **WHEN** the VHS workflow is defined
- **THEN** the vhs-action uses a major version tag (e.g., `@v2`)
- **AND** minor and patch updates are received automatically
- **AND** breaking changes require explicit version updates

### Requirement: Automated GIF Commit Attribution

The system SHALL attribute auto-committed GIF changes to a dedicated bot
identity to distinguish automated updates from developer contributions.

#### Scenario: Bot commits GIF updates

- **WHEN** the VHS workflow generates new GIFs
- **THEN** the commit author SHALL be identifiable as automated (e.g.,
  `vhs-action` bot)
- **AND** the commit message SHALL indicate it is an automated GIF update
- **AND** the commit SHALL only include `*.gif` files from `assets/gifs/`

### Requirement: Workflow Permissions Configuration

The system SHALL configure the VHS workflow with minimal required permissions to
commit changes back to the repository securely.

#### Scenario: Workflow has write access

- **WHEN** the VHS workflow needs to commit generated GIFs
- **THEN** it SHALL have `contents: write` permission
- **AND** it SHALL use the default `GITHUB_TOKEN` for authentication
- **AND** no additional secrets SHALL be required

### Requirement: Concurrency Management for VHS Workflow

The system SHALL cancel in-progress VHS workflow runs when new tape file changes
are pushed to conserve CI resources.

#### Scenario: Stale VHS run cancellation

- **WHEN** a developer pushes new tape file changes while a VHS workflow is
  running
- **THEN** the previous run is automatically cancelled
- **AND** a new workflow run starts for the latest changes
