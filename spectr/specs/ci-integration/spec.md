# Ci Integration Specification

## Purpose

Automates Spectr validation in GitHub Actions to ensure specification compliance for all commits and pull requests.

## Requirements

### Requirement: GitHub Action Workflow Integration
The system SHALL provide automated Spectr validation through a GitHub Action workflow job that executes on push and pull request events across all branches.

#### Scenario: Validation on push event
- **WHEN** a developer pushes commits to any branch
- **THEN** the `spectr-validate` job executes automatically
- **AND** the job runs the spectr-action to validate all changes

#### Scenario: Validation on pull request
- **WHEN** a pull request is opened or updated
- **THEN** the `spectr-validate` job executes as a required check
- **AND** validation results are visible in the PR status checks

### Requirement: Full Git History Access
The system SHALL configure the GitHub Action to checkout the repository with full git history to enable change detection and validation across commits.

#### Scenario: Full history checkout
- **WHEN** the spectr-validate job executes
- **THEN** the repository is checked out with `fetch-depth: 0`
- **AND** all git history is available for change tracking
- **AND** the spectr-action can detect changes across the full commit range

### Requirement: Fast Failure Pipeline Position
The system SHALL position the spectr-validate job as the first job in the CI pipeline to provide rapid feedback on specification violations before running longer-running tests.

#### Scenario: Job ordering for fast failure
- **WHEN** a CI pipeline is triggered
- **THEN** the `spectr-validate` job executes before lint, test, and format-check jobs
- **AND** developers receive validation feedback within seconds
- **AND** subsequent jobs do not run if spectr validation fails

### Requirement: Concurrency Management
The system SHALL cancel in-progress validation runs when new commits are pushed to the same branch to conserve CI resources and provide feedback on the latest changes.

#### Scenario: Stale run cancellation
- **WHEN** a developer pushes a new commit while a validation run is in progress
- **THEN** the previous run is automatically cancelled
- **AND** a new validation run starts for the latest commit
- **AND** CI resources are freed for the new run

### Requirement: Multi-Branch Support
The system SHALL execute spectr validation on all branches, not just main or specific feature branches, to ensure consistent quality across the development workflow.

#### Scenario: Validation on feature branch
- **WHEN** a developer pushes to a feature branch
- **THEN** the spectr-validate job executes with the same configuration as main
- **AND** validation rules are applied consistently

#### Scenario: Validation on main branch
- **WHEN** commits are merged to the main branch
- **THEN** the spectr-validate job executes to verify final state
- **AND** any validation failures block the merge

### Requirement: Action Version Pinning
The system SHALL use a specific version tag of the spectr-action (not `latest` or branch references) to ensure reproducible builds and prevent unexpected behavior from action updates.

#### Scenario: Version-pinned action reference
- **WHEN** the CI workflow is defined
- **THEN** the spectr-action uses a semantic version tag (e.g., `@v0.0.1`)
- **AND** the action version does not change unless explicitly updated
- **AND** builds are reproducible across time

### Requirement: VHS GIF Auto-Generation Workflow
The system SHALL provide automated GIF generation through a GitHub Actions workflow that triggers when VHS tape files are modified and commits updated GIFs back to the repository.

#### Scenario: GIF regeneration on tape file change
- **WHEN** a developer pushes changes to any `.tape` file in `docs/src/assets/vhs/`
- **THEN** the VHS workflow executes automatically
- **AND** VHS processes all tape files to generate updated GIFs
- **AND** the generated GIFs are committed back to the `docs/src/assets/gifs/` directory

#### Scenario: Workflow ignores non-tape changes
- **WHEN** a developer pushes changes that do not modify any `.tape` file
- **THEN** the VHS workflow does NOT trigger
- **AND** existing GIFs remain unchanged

#### Scenario: Multiple tape files processed
- **WHEN** multiple tape files are modified in a single commit
- **THEN** all modified tape files are processed in a single workflow run
- **AND** all corresponding GIFs are generated and committed together

### Requirement: VHS Action Version Pinning
The system SHALL use a specific major version tag of the vhs-action (e.g., `@v2`) to ensure reproducible builds while still receiving compatible updates.

#### Scenario: Version-pinned VHS action reference
- **WHEN** the VHS workflow is defined
- **THEN** the vhs-action uses a major version tag (e.g., `@v2`)
- **AND** minor and patch updates are received automatically
- **AND** breaking changes require explicit version updates

### Requirement: Automated GIF Commit Attribution
The system SHALL attribute auto-committed GIF changes to a dedicated bot identity to distinguish automated updates from developer contributions.

#### Scenario: Bot commits GIF updates
- **WHEN** the VHS workflow generates new GIFs
- **THEN** the commit author SHALL be identifiable as automated (e.g., `vhs-action` bot)
- **AND** the commit message SHALL indicate it is an automated GIF update
- **AND** the commit SHALL only include `*.gif` files from `docs/src/assets/gifs/`

### Requirement: Workflow Permissions Configuration
The system SHALL configure the VHS workflow with minimal required permissions to commit changes back to the repository securely.

#### Scenario: Workflow has write access
- **WHEN** the VHS workflow needs to commit generated GIFs
- **THEN** it SHALL have `contents: write` permission
- **AND** it SHALL use the default `GITHUB_TOKEN` for authentication
- **AND** no additional secrets SHALL be required

### Requirement: Concurrency Management for VHS Workflow
The system SHALL cancel in-progress VHS workflow runs when new tape file changes are pushed to conserve CI resources.

#### Scenario: Stale VHS run cancellation
- **WHEN** a developer pushes new tape file changes while a VHS workflow is running
- **THEN** the previous run is automatically cancelled
- **AND** a new workflow run starts for the latest changes
