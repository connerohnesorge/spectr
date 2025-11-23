# Medium Test Corpus

## Purpose
This file represents a typical capability spec with multiple requirements and scenarios.

## Requirements

### Requirement: User Registration
The system SHALL provide user registration functionality with email verification.

#### Scenario: Valid registration
- **WHEN** user submits valid registration form
- **THEN** system creates new user account
- **AND** sends verification email
- **AND** returns 201 Created status

#### Scenario: Duplicate email
- **WHEN** user attempts to register with existing email
- **THEN** system returns 409 Conflict
- **AND** displays error message

#### Scenario: Invalid email format
- **WHEN** user submits invalid email format
- **THEN** system returns 400 Bad Request
- **AND** provides validation error details

### Requirement: Email Verification
The system SHALL verify user email addresses before allowing full account access.

#### Scenario: Successful verification
- **WHEN** user clicks verification link
- **THEN** account is activated
- **AND** user is redirected to login page

#### Scenario: Expired verification link
- **WHEN** user clicks expired verification link
- **THEN** system displays error message
- **AND** provides option to resend verification

### Requirement: Profile Management
Users SHALL be able to view and update their profile information.

#### Scenario: View profile
- **WHEN** authenticated user accesses profile page
- **THEN** system displays current profile information
- **AND** shows editable fields

#### Scenario: Update profile
- **WHEN** user submits updated profile information
- **THEN** system validates changes
- **AND** saves updated information
- **AND** returns success confirmation

#### Scenario: Invalid profile data
- **WHEN** user submits invalid profile data
- **THEN** system returns validation errors
- **AND** original data remains unchanged

### Requirement: Password Management
Users SHALL be able to change their password when authenticated.

#### Scenario: Successful password change
- **WHEN** user provides current password and new password
- **THEN** system verifies current password
- **AND** updates to new password
- **AND** invalidates existing sessions

#### Scenario: Incorrect current password
- **WHEN** user provides incorrect current password
- **THEN** system returns 401 Unauthorized
- **AND** password remains unchanged

### Requirement: Account Deletion
Users SHALL be able to permanently delete their account.

#### Scenario: Delete account
- **WHEN** user confirms account deletion
- **THEN** system marks account as deleted
- **AND** schedules data removal
- **AND** logs user out immediately

#### Scenario: Cancel deletion
- **WHEN** user initiates but cancels deletion
- **THEN** account remains active
- **AND** no data is removed

### Requirement: Two-Factor Authentication
The system SHALL support optional two-factor authentication for enhanced security.

#### Scenario: Enable 2FA
- **WHEN** user enables 2FA
- **THEN** system generates QR code
- **AND** requires verification code
- **AND** stores 2FA secret

#### Scenario: Login with 2FA
- **WHEN** 2FA-enabled user logs in
- **THEN** system prompts for verification code
- **AND** validates code before granting access

#### Scenario: Disable 2FA
- **WHEN** user disables 2FA
- **THEN** system requires password confirmation
- **AND** removes 2FA requirement

### Requirement: Login Rate Limiting
The system SHALL implement rate limiting to prevent brute force attacks.

#### Scenario: Successful login attempts
- **WHEN** user makes multiple successful login attempts
- **THEN** no rate limiting is applied
- **AND** all attempts succeed

#### Scenario: Failed login attempts
- **WHEN** user makes 5 failed login attempts
- **THEN** account is temporarily locked
- **AND** cooldown period is enforced

#### Scenario: Rate limit cooldown
- **WHEN** cooldown period expires
- **THEN** user can attempt login again
- **AND** counter is reset

### Requirement: Session Expiration
User sessions SHALL automatically expire after a period of inactivity.

#### Scenario: Active session
- **WHEN** user makes request within timeout period
- **THEN** session remains active
- **AND** timeout is reset

#### Scenario: Expired session
- **WHEN** user makes request after timeout
- **THEN** system returns 401 Unauthorized
- **AND** redirects to login page

### Requirement: Remember Me Functionality
Users SHALL have option to stay logged in across browser sessions.

#### Scenario: Remember me enabled
- **WHEN** user selects remember me option
- **THEN** system issues long-lived token
- **AND** token is stored securely

#### Scenario: Remember me disabled
- **WHEN** remember me is not selected
- **THEN** session expires when browser closes
- **AND** user must login again

### Requirement: Audit Logging
The system SHALL log all authentication events for security monitoring.

#### Scenario: Login event logged
- **WHEN** user successfully logs in
- **THEN** system records login event
- **AND** includes timestamp and IP address

#### Scenario: Failed login logged
- **WHEN** login attempt fails
- **THEN** system records failure event
- **AND** includes attempted username

#### Scenario: Logout event logged
- **WHEN** user logs out
- **THEN** system records logout event
- **AND** includes session duration

## Code Example

Here's an example of how authentication works:

```go
func Authenticate(username, password string) (*Token, error) {
    user, err := db.FindUser(username)
    if err != nil {
        return nil, err
    }

    if !user.VerifyPassword(password) {
        return nil, ErrInvalidCredentials
    }

    token := jwt.Generate(user.ID)
    return token, nil
}
```

This code shows the basic authentication flow.
