## ADDED Requirements

### Requirement: CI Workflow Setup Option in Init Wizard
The initialization wizard SHALL provide an optional checkbox to create a GitHub Actions workflow file (`.github/workflows/spectr-ci.yml`) for automated Spectr validation during CI/CD. This option integrates seamlessly with the existing tool selection step.

#### Scenario: CI workflow option displayed in tool selection
- **WHEN** user runs `spectr init` and reaches the tool selection screen
- **THEN** a "GitHub Actions CI Workflow" option is displayed alongside AI tool options
- **AND** the option includes a brief description: "Automated Spectr validation on push/PR"
- **AND** the option uses the same checkbox styling as AI tool entries

#### Scenario: CI workflow option detects existing configuration
- **WHEN** user runs `spectr init` on a project that already has `.github/workflows/spectr-ci.yml`
- **AND** user reaches the tool selection screen
- **THEN** the "GitHub Actions CI Workflow" option shows a "(configured)" indicator
- **AND** the option is pre-selected by default
- **AND** selecting it will update the existing workflow file

#### Scenario: CI workflow option not pre-selected on fresh projects
- **WHEN** user runs `spectr init` on a project without `.github/workflows/spectr-ci.yml`
- **AND** user reaches the tool selection screen
- **THEN** the "GitHub Actions CI Workflow" option is NOT pre-selected by default
- **AND** the user must explicitly select it to enable CI workflow creation

#### Scenario: CI workflow created when selected
- **WHEN** user selects the "GitHub Actions CI Workflow" option
- **AND** user confirms the selection and proceeds with initialization
- **THEN** the system creates `.github/workflows/` directory if it doesn't exist
- **AND** the system creates `.github/workflows/spectr-ci.yml` with the standard validation workflow
- **AND** the workflow file is tracked in the execution result as created or updated

#### Scenario: CI workflow not created when unselected
- **WHEN** user does NOT select the "GitHub Actions CI Workflow" option
- **AND** user confirms the selection and proceeds with initialization
- **THEN** no `.github/workflows/spectr-ci.yml` file is created
- **AND** any existing `.github/workflows/spectr-ci.yml` file is left unchanged

#### Scenario: CI workflow content follows spec requirements
- **WHEN** the CI workflow file is created
- **THEN** the workflow triggers on push and pull request events
- **AND** the workflow uses `fetch-depth: 0` for full git history
- **AND** the workflow uses the `spectr-action` for validation
- **AND** the workflow includes concurrency management to cancel in-progress runs
- **AND** the workflow runs on `ubuntu-latest`

#### Scenario: Review screen shows CI workflow selection
- **WHEN** user has selected the "GitHub Actions CI Workflow" option
- **AND** user proceeds to the review screen
- **THEN** the review shows "GitHub Actions CI Workflow" in the list of selected items
- **AND** the creation plan mentions the workflow file path

#### Scenario: Completion screen shows CI workflow file
- **WHEN** the CI workflow file is successfully created
- **THEN** the completion screen lists `.github/workflows/spectr-ci.yml` in created files
- **AND** the file path is displayed with the appropriate icon (created or updated)

#### Scenario: Non-interactive mode supports CI workflow flag
- **WHEN** user runs `spectr init --non-interactive --ci-workflow`
- **THEN** the CI workflow file is created without TUI interaction
- **AND** the workflow file content matches the interactive mode output
