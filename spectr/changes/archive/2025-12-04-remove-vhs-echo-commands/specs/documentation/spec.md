# Documentation Specification (Delta)

## MODIFIED Requirements

### Requirement: VHS Tape Output Clarity

VHS tape files SHALL NOT contain typed echo statements that display section
headers or commentary. Demos SHALL let the spectr commands and their output
speak for themselves. Comments within the tape file (lines starting with `#`)
SHOULD be used to document sections for maintainers, but these are not displayed
in the recording.

#### Scenario: No typed echo section headers

- **WHEN** a VHS tape file is reviewed
- **THEN** it SHALL contain no `Type "echo ..."` commands for section headers
- **AND** visual context SHALL be provided through VHS comments (starting with
  `#`) which are not recorded

#### Scenario: No useless echo statements

- **WHEN** a VHS tape file is reviewed
- **THEN** it SHALL contain no `Type "echo ''"` commands
- **AND** visual spacing SHALL be achieved through `Sleep` commands only

#### Scenario: Commands are self-documenting

- **WHEN** a user views a demo GIF
- **THEN** they SHALL see only the actual spectr commands being typed
- **AND** they SHALL NOT see preparatory echo statements being typed out
