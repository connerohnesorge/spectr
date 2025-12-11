## 1. Implementation

- [ ] 1.1 Add `NormalizeItemPath` function in `internal/discovery/normalize.go` that extracts item ID and infers type from paths
- [ ] 1.2 Add unit tests for `NormalizeItemPath` in `internal/discovery/normalize_test.go`
- [ ] 1.3 Integrate `NormalizeItemPath` into `cmd/validate.go` before calling `DetermineItemType`
- [ ] 1.4 Integrate `NormalizeItemPath` into `cmd/archive.go` before resolving change ID
- [ ] 1.5 Integrate `NormalizeItemPath` into `cmd/accept.go` before resolving change ID

## 2. Validation

- [ ] 2.1 Run `spectr validate add-path-normalization-validate --strict` to verify proposal
- [ ] 2.2 Run `go test ./internal/discovery/...` to ensure tests pass
- [ ] 2.3 Manually test `spectr validate spectr/changes/<id>` works correctly
- [ ] 2.4 Manually test `spectr archive spectr/changes/<id>` works correctly
