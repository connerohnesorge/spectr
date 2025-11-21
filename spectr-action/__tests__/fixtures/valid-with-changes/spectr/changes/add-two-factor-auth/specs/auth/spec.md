## ADDED Requirements

### Requirement: Two-Factor Authentication
The system SHALL support TOTP-based two-factor authentication for enhanced account security.

#### Scenario: Successful 2FA enrollment
- **WHEN** a user initiates 2FA enrollment
- **THEN** the system generates a TOTP secret
- **AND** displays a QR code for authenticator app setup
- **AND** generates recovery codes for backup access

#### Scenario: Successful login with 2FA enabled
- **WHEN** a user with 2FA enabled provides valid credentials
- **THEN** the system prompts for a TOTP code
- **AND** validates the provided TOTP code
- **AND** issues an authentication token upon successful verification

#### Scenario: Failed login with invalid TOTP code
- **WHEN** a user provides valid credentials but an invalid TOTP code
- **THEN** the system rejects the login attempt
- **AND** does not issue an authentication token
- **AND** logs the failed 2FA verification

### Requirement: Recovery Code Support
The system MUST provide recovery codes as a backup authentication method when 2FA is enabled.

#### Scenario: Login using recovery code
- **WHEN** a user cannot access their authenticator app
- **AND** provides a valid recovery code
- **THEN** the system authenticates the user
- **AND** marks the recovery code as used
- **AND** issues an authentication token

#### Scenario: Recovery code exhaustion warning
- **WHEN** a user has used all but one recovery code
- **THEN** the system displays a warning to generate new codes
- **AND** prompts the user to regenerate recovery codes after login
