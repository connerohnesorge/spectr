# Small Test Corpus

## Purpose
This is a small test file for benchmarking parser performance on typical change deltas.

## ADDED Requirements

### Requirement: User Authentication
The system SHALL provide secure user authentication using JWT tokens.

#### Scenario: Successful login
- **WHEN** user provides valid credentials
- **THEN** system returns JWT token
- **AND** token expires after 24 hours

#### Scenario: Invalid credentials
- **WHEN** user provides invalid credentials
- **THEN** system returns 401 Unauthorized
- **AND** error message is displayed

### Requirement: Session Management
The system SHALL manage user sessions securely.

#### Scenario: Token validation
- **WHEN** authenticated request is made
- **THEN** system validates JWT token
- **AND** grants access if valid

## MODIFIED Requirements

### Requirement: Password Reset
The system SHALL allow users to reset their password via email.

#### Scenario: Reset email sent
- **WHEN** user requests password reset
- **THEN** system sends email with reset link
- **AND** link expires after 1 hour

## REMOVED Requirements

### Requirement: Legacy Login Method
**Reason**: Deprecated in favor of JWT authentication
**Migration**: Update all clients to use JWT tokens
