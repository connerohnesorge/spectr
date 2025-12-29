## MODIFIED Requirements

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

### Requirement: Automated GIF Commit Attribution

The system SHALL attribute auto-committed GIF changes to a dedicated bot identity to distinguish automated updates from developer contributions.

#### Scenario: Bot commits GIF updates

- **WHEN** the VHS workflow generates new GIFs
- **THEN** the commit author SHALL be identifiable as automated (e.g., `vhs-action` bot)
- **AND** the commit message SHALL indicate it is an automated GIF update
- **AND** the commit SHALL only include `*.gif` files from `docs/src/assets/gifs/`
