## 1. Implementation

- [ ] 1.1 Add `PriorityCrush` constant to `internal/initialize/providers/constants.go`
- [ ] 1.2 Create `internal/initialize/providers/crush.go` with Crush provider implementation
- [ ] 1.3 Verify provider registers correctly with `go build ./...`

## 2. Testing

- [ ] 2.1 Run existing provider tests to ensure no regressions: `go test ./internal/initialize/providers/...`
- [ ] 2.2 Test `spectr init` with Crush provider manually selected
- [ ] 2.3 Verify generated files are correct (CRUSH.md, .crush/commands/spectr/)

## 3. Documentation

- [ ] 3.1 Update README.md to include Crush in the list of supported AI tools
