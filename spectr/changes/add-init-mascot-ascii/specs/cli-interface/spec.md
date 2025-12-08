## ADDED Requirements

### Requirement: ASCII Mascot in Initialization Banner

The `spectr init` wizard intro screen SHALL display an ASCII art representation of the Spectr ghost gopher mascot alongside the "SPECTR" text logo, creating a branded visual identity during project initialization.

The mascot design depicts a ghost with Go gopher characteristics (round eyes, buck teeth, small ears), matching the official Spectr logo in `assets/logo.png`.

#### Scenario: Init wizard displays mascot and text logo

- **WHEN** user runs `spectr init` and enters the interactive wizard
- **THEN** the intro screen displays both the ASCII ghost mascot and "SPECTR" text
- **AND** both elements are rendered with the purple-to-pink gradient styling
- **AND** the combined banner fits within 80 columns for standard terminal widths

#### Scenario: Mascot renders with gradient colors

- **WHEN** the init wizard intro screen is rendered
- **THEN** the mascot ASCII art receives the same gradient color treatment as the text logo
- **AND** colors transition smoothly from left to right (or top to bottom) across the mascot

#### Scenario: Banner maintains readability

- **WHEN** the combined mascot and text banner is displayed
- **THEN** the mascot is clearly recognizable as a ghost figure
- **AND** the "SPECTR" text remains fully legible
- **AND** there is appropriate spacing between mascot and text elements
