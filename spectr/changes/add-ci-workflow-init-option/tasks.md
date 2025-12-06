## 1. Template and Configuration

- [ ] 1.1 Add CI workflow template constant to `internal/initialize/templates.go` with embedded `spectr-ci.yml` content (spectr-validate job only, using `spectr-action@v0.0.2`)
- [ ] 1.2 Add `RenderCIWorkflow()` method to `TemplateManager` to generate the workflow file content

## 2. Wizard TUI Updates

- [ ] 2.1 Add `ciWorkflowEnabled` and `ciWorkflowConfigured` fields to `WizardModel`
- [ ] 2.2 Add detection logic in `NewWizardModel` to check if `.github/workflows/spectr-ci.yml` exists (set `ciWorkflowConfigured` and pre-select if true)
- [ ] 2.3 Update `renderReview()` to display the "Spectr CI Validation" checkbox option after tool summary
- [ ] 2.4 Update `handleReviewKeys()` to support Space key for toggling CI option (in addition to Enter/Backspace)
- [ ] 2.5 Update `renderCreationPlan()` to conditionally show `.github/workflows/spectr-ci.yml` when CI is enabled

## 3. Executor Integration

- [ ] 3.1 Add `createCIWorkflow()` method to `InitExecutor` to create the `.github/workflows/` directory and `spectr-ci.yml` file
- [ ] 3.2 Update `Execute()` signature to accept CI workflow enablement flag
- [ ] 3.3 Call `createCIWorkflow()` when CI workflow is enabled, track created/updated files in result

## 4. Non-Interactive Mode

- [ ] 4.1 Add `CIWorkflow bool` field with `--ci-workflow` flag to `InitCmd` in `cmd/init.go`
- [ ] 4.2 Pass `--ci-workflow` flag value to executor in non-interactive mode

## 5. Testing

- [ ] 5.1 Add unit tests for CI workflow template rendering
- [ ] 5.2 Add unit tests for CI workflow file existence detection
- [ ] 5.3 Add unit tests for Review step with CI option toggling

## 6. Validation

- [ ] 6.1 Run `spectr validate add-ci-workflow-init-option --strict` and fix any issues
- [ ] 6.2 Manually test the init wizard flow with CI option enabled and disabled
- [ ] 6.3 Test non-interactive mode with and without `--ci-workflow` flag
