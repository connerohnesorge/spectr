## MODIFIED Requirements

### Requirement: Pre-made Example Projects for VHS Demos

The system SHALL provide pre-made spectr project examples in the `examples/` directory that VHS tape files use for demonstrations by executing commands directly within the example directories, keeping demos focused on spectr commands rather than file creation boilerplate.

#### Scenario: Developer creates clean demo

- **WHEN** a VHS tape file needs a spectr project for demonstration
- **THEN** it SHALL execute spectr commands directly in the `examples/<demo-name>/` directory
- **AND** the demo output SHALL focus on spectr commands, not temporary directory copying

#### Scenario: Demo runs in example directory

- **WHEN** a VHS tape demonstrates a spectr command
- **THEN** the tape SHALL use `Hide`/`Show` to conceal the directory change
- **AND** spectr commands SHALL run directly in the example directory without copying to `_demo`

#### Scenario: Developer maintains example project

- **WHEN** a change to demo content is needed
- **THEN** the developer SHALL edit the pre-made example in `examples/` directory
- **AND** the change will automatically apply to any tape using that example

#### Scenario: Demo cleanup is minimal

- **WHEN** a tape completes execution
- **THEN** no `rm -rf _demo` cleanup commands SHALL be needed
- **AND** the example directory remains unchanged

## ADDED Requirements

### Requirement: VHS Tape Output Clarity

VHS tape files SHALL NOT contain empty echo statements that only produce blank lines without meaningful output.

#### Scenario: No useless echo statements

- **WHEN** a VHS tape file is reviewed
- **THEN** it SHALL contain no `Type "echo ''"` commands
- **AND** visual spacing SHALL be achieved through `Sleep` commands only
