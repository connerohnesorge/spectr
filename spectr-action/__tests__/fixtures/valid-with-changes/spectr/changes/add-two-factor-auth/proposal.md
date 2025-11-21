# Change: Add Two-Factor Authentication

## Why

Current authentication relies solely on username/password, which is vulnerable to credential theft. Adding two-factor authentication (2FA) provides an additional security layer to protect user accounts from unauthorized access.

## What Changes

- Add support for TOTP-based two-factor authentication
- Require 2FA verification after successful password authentication
- Provide enrollment flow for users to set up 2FA
- Add recovery codes for account recovery

## Impact

- Affected specs: `auth`
- Affected code: Authentication service, user management API, login UI
- Breaking change: **BREAKING** - Login flow now requires additional verification step for users with 2FA enabled
