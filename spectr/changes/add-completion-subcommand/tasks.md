## 1. Dependencies

- [ ] 1.1 Add `github.com/jotaen/kong-completion` dependency to `go.mod`
- [ ] 1.2 Run `go mod tidy` to resolve dependency tree

## 2. Core Implementation

- [ ] 2.1 Create `cmd/completion.go` with predictor registrations
- [ ] 2.2 Define custom predictor for change IDs (scans `spectr/changes/` directory)
- [ ] 2.3 Define custom predictor for spec IDs (scans `spectr/specs/` directory)
- [ ] 2.4 Define custom predictor for item types (`change`, `spec`)
- [ ] 2.5 Add `Completion` field to CLI struct in `cmd/root.go`
- [ ] 2.6 Modify `main.go` to use kong-completion registration pattern

## 3. Testing

- [ ] 3.1 Test completion subcommand outputs valid shell scripts for bash/zsh/fish
- [ ] 3.2 Test custom predictors return correct values
- [ ] 3.3 Verify existing commands work unchanged after modification

## 4. Documentation

- [ ] 4.1 Update CLI help text to describe completion command usage
