## 1. Package Setup
- [ ] 1.1 Create `internal/specterrs/` directory
- [ ] 1.2 Create `doc.go` with package documentation
- [ ] 1.3 Create `git.go` with 5 error types
- [ ] 1.4 Create `archive.go` with 6 error types
- [ ] 1.5 Create `validation.go` with 3 error types
- [ ] 1.6 Create `initialize.go` with 3 error types
- [ ] 1.7 Create `list.go` with 1 error type
- [ ] 1.8 Create `environment.go` with 1 error type
- [ ] 1.9 Create `pr.go` with 2 error types
- [ ] 1.10 Run `go build ./...` to verify package compiles

## 2. Migration - Low Coupling
- [ ] 2.1 Migrate `internal/list/interactive.go` (EditorNotSetError)
- [ ] 2.2 Migrate `cmd/list.go` (IncompatibleFlagsError)
- [ ] 2.3 Migrate `internal/initialize/filesystem.go` (EmptyPathError)
- [ ] 2.4 Migrate `cmd/init.go` (WizardModelCastError, InitializationCompletedWithErrorsError)
- [ ] 2.5 Run tests: `go test ./internal/list/... ./internal/initialize/... ./cmd/...`

## 3. Migration - Medium Coupling
- [ ] 3.1 Migrate `internal/git/worktree.go` (BranchNameRequiredError, BaseBranchRequiredError, NotInGitRepositoryError)
- [ ] 3.2 Migrate `internal/git/platform.go` (EmptyRemoteURLError)
- [ ] 3.3 Migrate `internal/git/branch.go` (NotInGitRepositoryError, BaseBranchNotFoundError)
- [ ] 3.4 Migrate `internal/pr/platforms.go` (UnknownPlatformError)
- [ ] 3.5 Migrate `internal/pr/workflow.go` (PRPrerequisiteError)
- [ ] 3.6 Run tests: `go test ./internal/git/... ./internal/pr/...`

## 4. Migration - High Coupling
- [ ] 4.1 Migrate `cmd/validate.go` (ValidationFailedError)
- [ ] 4.2 Migrate `internal/validation/delta_validators.go` (DeltaSpecParseError)
- [ ] 4.3 Migrate `cmd/accept.go` (ValidationRequiredError)
- [ ] 4.4 Run tests: `go test ./cmd/... ./internal/validation/...`

## 5. Migration - Archive (Highest Coupling)
- [ ] 5.1 Migrate `internal/archive/archiver.go` (UserCancelledError, ArchiveCancelledError, ValidationRequiredError)
- [ ] 5.2 Migrate `internal/archive/validator.go` (DeltaConflictError - 6 instances)
- [ ] 5.3 Remove `ErrUserCancelled` from `internal/archive/types.go`
- [ ] 5.4 Remove `errEmptyPath` constant from `internal/initialize/filesystem.go`
- [ ] 5.5 Run tests: `go test ./internal/archive/...`

## 6. Verification
- [ ] 6.1 Run full test suite: `go test ./...`
- [ ] 6.2 Run linter: `golangci-lint run`
- [ ] 6.3 Verify no `errors.New()` with hardcoded strings remain (grep check)
- [ ] 6.4 Manual smoke test of CLI commands
