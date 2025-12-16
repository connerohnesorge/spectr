# Change: Remove CodeBuddy Support

## Why
CodeBuddy is being removed because the tool has been discontinued/deprecated and is no longer actively maintained. Removing the integration reduces maintenance burden and simplifies the provider list for active, supported tools.

## What Changes
- **BREAKING**: Remove CodeBuddy provider implementation
- Delete `internal/initialize/providers/codebuddy.go` (33 lines)
- Remove `PriorityCodeBuddy` constant from `internal/initialize/providers/constants.go`
- Remove CodeBuddy from README.md supported tools table
- Remove `support-codebuddy` specification directory

## Impact
- Existing CodeBuddy users will need to reconfigure with a different provider
- Provider count reduces from 17 to 16 supported tools
- No migration path needed - users can run `spectr init` with a new provider selection
- The priority constant gap at 5 is acceptable as priorities are only used for sorting
