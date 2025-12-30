# Delta Specification

## ADDED Requirements

### Requirement: CI Workflow Setup Option in Init Wizard Review Step

The initialization wizard's Review step SHALL include an optional checkbox to
create a GitHub Actions workflow file (`.github/workflows/spectr-ci.yml`) for
automated Spectr validation during CI/CD. This option is presented alongside the
tool selection summary, keeping the wizard flow quick without adding a separate
step.

#### Scenario: CI option displayed in Review step

- **WHEN** user completes tool selection and proceeds to the Review step
- **THEN** a "Spectr CI Validation" checkbox option is displayed after the tool
  summary
- **AND** the option appears before the creation plan section
- **AND** a description explains: "Validate specs automatically on push and pull
  requests"

#### Scenario: CI option detects existing workflow

- **WHEN** user runs `spectr init` on a project that already has
  `.github/workflows/spectr-ci.yml`
- **AND** user reaches the Review step
- **THEN** the "Spectr CI Validation" option shows a "(configured)" indicator
- **AND** the option is pre-selected by default
- **AND** selecting it will update the existing workflow file

#### Scenario: CI option not pre-selected on fresh projects

- **WHEN** user runs `spectr init` on a project without
  `.github/workflows/spectr-ci.yml`
- **AND** user reaches the Review step
- **THEN** the "Spectr CI Validation" option is NOT pre-selected by default
- **AND** the user must explicitly select it to enable CI workflow creation

#### Scenario: User toggles CI option in Review step

- **WHEN** user is on the Review step
- **AND** user presses Space
- **THEN** the CI workflow checkbox toggles between selected and unselected
- **AND** the creation plan updates to reflect the change
- **AND** the visual state updates immediately

#### Scenario: CI workflow created when selected

- **WHEN** user selects the "Spectr CI Validation" option in Review
- **AND** user presses Enter to proceed with initialization
- **THEN** the system creates `.github/workflows/` directory if it doesn't exist
- **AND** the system creates `.github/workflows/spectr-ci.yml` with the Spectr
  validation workflow
- **AND** the workflow file is tracked in the execution result as created or
  updated

#### Scenario: CI workflow not created when unselected

- **WHEN** user does NOT select the "Spectr CI Validation" option
- **AND** user proceeds with initialization
- **THEN** no `.github/workflows/spectr-ci.yml` file is created
- **AND** any existing `.github/workflows/spectr-ci.yml` file is left unchanged

#### Scenario: CI workflow content uses pinned action version

- **WHEN** the CI workflow file is created
- **THEN** the workflow contains a single `spectr-validate` job
- **AND** the workflow uses `connerohnesorge/spectr-action@v0.0.2` (pinned
  version)
- **AND** the workflow triggers on push to `main` branch only
- **AND** the workflow triggers on pull requests to all branches
- **AND** the workflow uses `fetch-depth: 0` for full git history
- **AND** the workflow includes concurrency management to cancel in-progress
  runs
- **AND** the workflow runs on `ubuntu-latest`

#### Scenario: Creation plan shows CI workflow when enabled

- **WHEN** user has selected the "Spectr CI Validation" option
- **THEN** the creation plan section shows `.github/workflows/spectr-ci.yml`
- **AND** the file is listed with the tool configurations

#### Scenario: Creation plan hides CI workflow when disabled

- **WHEN** user has NOT selected the "Spectr CI Validation" option
- **THEN** the creation plan does NOT mention `.github/workflows/spectr-ci.yml`

#### Scenario: Completion screen shows CI workflow file

- **WHEN** the CI workflow file is successfully created
- **THEN** the completion screen lists `.github/workflows/spectr-ci.yml` in
  created or updated files
- **AND** the file path is displayed with the appropriate icon

#### Scenario: Non-interactive mode supports CI workflow flag

- **WHEN** user runs `spectr init --non-interactive --ci-workflow`
- **THEN** the CI workflow file is created without TUI interaction
- **AND** the workflow file content matches the interactive mode output

#### Scenario: Non-interactive mode without CI flag skips workflow

- **WHEN** user runs `spectr init --non-interactive` without `--ci-workflow`
- **THEN** no CI workflow file is created
- **AND** existing workflow files are not modified

#### Scenario: Review step help text includes Space for toggle

- **WHEN** user is on the Review step
- **THEN** the help text shows: "Space: Toggle CI Enter: Initialize Backspace:
  Back q: Quit"
