# Validation Spec Delta

## ADDED Requirements

### Requirement: Root-Level Prevention
The validation system SHALL prevent requirements or specifications from being added directly at the root level without proper capability organization.

#### Scenario: Detect root-level spec file
- **WHEN** a spec file exists directly under `spectr/specs/` without a capability subdirectory
- **THEN** validation SHALL fail with error "Spec files must be organized under capability directories"

#### Scenario: Detect root-level requirements
- **WHEN** requirements are defined outside of a capability directory structure
- **THEN** validation SHALL fail with error "Requirements must be defined within a capability directory"

#### Scenario: Valid capability structure
- **WHEN** spec files are properly organized under `spectr/specs/[capability]/spec.md`
- **THEN** validation SHALL pass

## MODIFIED Requirements

### Requirement: Directory Structure Validation
The validation system SHALL enforce proper capability directory structure and provide clear guidance when violations are detected.

#### Scenario: Invalid directory hierarchy in strict mode
- **WHEN** strict mode validation is enabled
- **THEN** the system SHALL check for proper capability directory hierarchy
- **AND** SHALL report specific structural violations with corrective guidance

#### Scenario: Helpful error messages
- **WHEN** a structural violation is detected
- **THEN** error messages SHALL include:
  - Clear description of the violation
  - Expected directory structure
  - Example of correct organization

#### Scenario: Backward compatibility
- **WHEN** validating existing specs with proper structure
- **THEN** validation SHALL continue to pass
- **AND** no breaking changes SHALL be introduced to valid structures
