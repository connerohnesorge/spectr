# Tasks

## Code Changes

- [x] Rename `ModeNew` constant to `ModeProposal` in `internal/pr/templates.go`
- [x] Rename `PRNewCmd` struct to `PRProposalCmd` in `cmd/pr.go`
- [x] Rename `New` field to `Proposal` in `PRCmd` struct with `cmd:"proposal"`
  annotation
- [x] Update error message in `cmd/pr.go` from "pr new failed" to "pr proposal
  failed"
- [x] Update all `ModeNew` references in `internal/pr/workflow.go`
- [x] Update all `ModeNew` references in `internal/pr/dryrun.go`
- [x] Update all `ModeNew` references in `internal/pr/templates.go` (switch
  cases, etc.)

## Test Updates

- [x] Update test files: `internal/pr/templates_test.go` (ModeNew ->
  ModeProposal)
- [x] Update test files: `internal/pr/workflow_test.go` (ModeNew ->
  ModeProposal)
- [x] Update test files: `internal/pr/integration_test.go` (ModeNew ->
  ModeProposal)
- [x] Update test files: `cmd/pr_test.go` (PRNewCmd -> PRProposalCmd)

## Pending Change Updates

- [x] Update `spectr/changes/add-pr-subcommand/specs/cli-interface/spec.md` to
  use `proposal` terminology
- [x] Update `spectr/changes/add-pr-subcommand/proposal.md` if it mentions `pr
  new`
- [x] Update `spectr/changes/add-pr-subcommand/tasks.md` if it mentions `pr new`

## Validation

- [x] Run `go build` to verify compilation
- [x] Run `go test ./...` to verify all tests pass
- [x] Run `spectr validate rename-pr-new-to-proposal --strict` to validate the
  change
