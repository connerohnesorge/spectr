# Delta Specification

## ADDED Requirements

### Requirement: Two-Factor Authentication

The system SHALL provide two-factor authentication support.

#### Scenario: OTP verification

- **WHEN** user enables 2FA
- **THEN** system sends OTP code

#### Scenario: OTP validation

- **WHEN** user enters correct OTP
- **THEN** system grants access

## MODIFIED Requirements

### Requirement: User Login

The system SHALL support user login with optional 2FA.

#### Scenario: Login with 2FA

- **WHEN** user with 2FA enabled logs in
- **THEN** system requires OTP verification

#### Scenario: Login without 2FA

- **WHEN** user without 2FA logs in
- **THEN** system grants immediate access
