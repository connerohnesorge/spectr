# Tasks

## Implementation

- [ ] 1. Update `cmd/validate.go` to remove `Strict` field from `ValidateCmd` struct
- [ ] 2. Update `cmd/validate.go` to always pass `true` to `NewValidator()`
- [ ] 3. Update `internal/validation/interactive.go` to remove strict parameter from `RunInteractiveValidation()`
- [ ] 4. Update `internal/validation/validator.go` to remove `strictMode` field and always behave strictly
- [ ] 5. Update `internal/validation/spec_rules.go` to always apply warning-to-error conversion
- [ ] 6. Update `internal/validation/change_rules.go` to always apply warning-to-error conversion

## Testing

- [ ] 7. Update `internal/validation/validator_test.go` to remove strict mode tests and verify always-strict behavior
- [ ] 8. Run `go test ./...` to verify all tests pass
- [ ] 9. Run `go build ./...` to verify build succeeds

## Documentation

- [ ] 10. Update `docs/src/content/docs/reference/cli-commands.md` to remove `--strict` flag documentation
- [ ] 11. Update `README.md` to remove any `--strict` flag examples

## Validation

- [ ] 12. Run `spectr validate remove-strict-flag` to validate this change (using current strict mode)
- [ ] 13. Test CLI to verify `--strict` flag is rejected as unknown
