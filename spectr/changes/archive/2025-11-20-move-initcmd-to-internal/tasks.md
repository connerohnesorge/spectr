# Implementation Tasks

## 1. Move InitCmd struct definition
- [ ] 1.1 Copy InitCmd struct (including doc comment) from cmd/init.go to internal/init/models.go
- [ ] 1.2 Verify InitCmd is exported and properly documented
- [ ] 1.3 Remove InitCmd struct from cmd/init.go

## 2. Update imports in cmd layer
- [ ] 2.1 Update cmd/init.go to import InitCmd from internal/init
- [ ] 2.2 Update cmd/root.go CLI struct to use initpkg.InitCmd
- [ ] 2.3 Verify all cmd layer files compile without errors

## 3. Refactor internal/init functions
- [ ] 3.1 Update internal/init/executor.go or related functions to accept InitCmd struct
- [ ] 3.2 Update all internal/init function signatures to use InitCmd instead of individual parameters
- [ ] 3.3 Update tests in cmd/init_test.go to reflect new import
- [ ] 3.4 Update tests in internal/init/*_test.go to use new function signatures

## 4. Validation and testing
- [ ] 4.1 Run full test suite: `go test ./...`
- [ ] 4.2 Verify no lint errors: `golangci-lint run`
- [ ] 4.3 Build the project: `go build -o spectr ./`
- [ ] 4.4 Manual test: `./spectr init --help` displays correct flags

## 5. Documentation updates
- [ ] 5.1 Verify doc comments on InitCmd in internal/init/models.go are clear
- [ ] 5.2 Update any code comments referencing InitCmd location
