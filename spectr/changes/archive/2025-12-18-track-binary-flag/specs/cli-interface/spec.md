# Delta Specification

## ADDED Requirements

### Requirement: Track Command Binary Filtering

The track command SHALL support binary file filtering to prevent unintentional
commits of binary files.

#### Scenario: Include-binaries flag enables binary file commits

- **WHEN** user provides the `--include-binaries` flag
- **THEN** the system includes binary files in automatic commits
- **AND** commits all modified files as before (excluding task files)

#### Scenario: Default binary file exclusion

- **WHEN** user runs `spectr track` without the `--include-binaries` flag
- **THEN** the system detects binary files using git diff --numstat
- **AND** excludes binary files from staging
- **AND** displays a warning listing skipped binary files
- **AND** continues to commit non-binary files normally

#### Scenario: Binary detection via git

- **WHEN** the system needs to determine if a file is binary
- **THEN** it SHALL use `git diff --numstat` to check for binary markers
- **AND** treat files with `-       -` output as binary

#### Scenario: Only binary files modified

- **WHEN** a task status changes and only binary files were modified (with no
  --include-binaries flag)
- **THEN** the system displays a warning: "No files to commit for task
  \<task-id\> (binary files excluded)"
- **AND** lists the skipped binary files
- **AND** continues watching for more task changes
