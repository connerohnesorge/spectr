# Implementation Tasks

## 1. Template and Configuration

- [x] 1.1 Add CI workflow template constant to
  `internal/initialize/templates.go` with embedded `spectr-ci.yml` content
  (spectr-validate job only, using `spectr-action@v0.0.2`)
- [x] 1.2 Add `RenderCIWorkflow()` method to `TemplateManager` to generate the
  workflow file content

## 2. Wizard TUI Updates

- [x] 2.1 Add `ciWorkflowEnabled` and `ciWorkflowConfigured` fields to
  `WizardModel`
- [x] 2.2 Add detection logic in `NewWizardModel` to check if
  `.github/workflows/spectr-ci.yml` exists (set `ciWorkflowConfigured` and
  pre-select if true)
- [x] 2.3 Update `renderReview()` to display the "Spectr CI Validation" checkbox
  option after tool summary
- [x] 2.4 Update `handleReviewKeys()` to support Space key for toggling CI
  option (in addition to Enter/Backspace)
- [x] 2.5 Update `renderCreationPlan()` to conditionally show
  `.github/workflows/spectr-ci.yml` when CI is enabled

## 3. Executor Integration

- [x] 3.1 Add `createCIWorkflow()` method to `InitExecutor` to create the
  `.github/workflows/` directory and `spectr-ci.yml` file
- [x] 3.2 Update `Execute()` signature to accept CI workflow enablement flag
- [x] 3.3 Call `createCIWorkflow()` when CI workflow is enabled, track
  created/updated files in result

## 4. Non-Interactive Mode

- [x] 4.1 Add `CIWorkflow bool` field with `--ci-workflow` flag to `InitCmd` in
  `cmd/init.go`
- [x] 4.2 Pass `--ci-workflow` flag value to executor in non-interactive mode

## 5. Testing

- [x] 5.1 Add unit tests for CI workflow template rendering
- [x] 5.2 Add unit tests for CI workflow file existence detection
- [x] 5.3 Add unit tests for Review step with CI option toggling

## 6. Validation

- [x] 6.1 Run `spectr validate add-ci-workflow-init-option --strict` and fix any
  issues
- [x] 6.2 Manually test the init wizard flow with CI option enabled and disabled
- [x] 6.3 Test non-interactive mode with and without `--ci-workflow` flag
