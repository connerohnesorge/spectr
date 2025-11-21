# Example Capability

## Purpose

This capability demonstrates a minimal valid Spectr specification for testing purposes. It includes basic authentication and data validation requirements.

## Requirements

### Requirement: User Authentication
The system SHALL authenticate users before granting access to protected resources.

#### Scenario: Successful authentication with valid credentials
- **WHEN** a user provides valid username and password
- **THEN** the system returns an authentication token
- **AND** the user gains access to protected resources

#### Scenario: Failed authentication with invalid credentials
- **WHEN** a user provides invalid credentials
- **THEN** the system returns an authentication error
- **AND** the user is denied access to protected resources

### Requirement: Data Validation
The system MUST validate all user input before processing.

#### Scenario: Valid input processing
- **WHEN** user input passes validation rules
- **THEN** the system processes the input
- **AND** returns a success response

#### Scenario: Invalid input rejection
- **WHEN** user input fails validation rules
- **THEN** the system rejects the input
- **AND** returns clear validation error messages
