## 1. Template and Configuration

- [ ] 1.1 Add CI workflow template to `internal/initialize/templates.go` with embedded `spectr-ci.yml` content
- [ ] 1.2 Add `RenderCIWorkflow()` method to `TemplateManager` to generate the workflow file content

## 2. Wizard TUI Updates

- [ ] 2.1 Add a new checkbox option in the tool selection step for "GitHub Actions CI Workflow"
- [ ] 2.2 Implement detection logic to check if `.github/workflows/spectr-ci.yml` already exists (mark as configured)
- [ ] 2.3 Add descriptive text explaining what the CI workflow does when the option is highlighted

## 3. Executor Integration

- [ ] 3.1 Add `createCIWorkflow()` method to `InitExecutor` to create the `.github/workflows/` directory and workflow file
- [ ] 3.2 Integrate CI workflow creation into `Execute()` flow based on selection state
- [ ] 3.3 Track created/updated files appropriately in the execution result

## 4. Testing

- [ ] 4.1 Add unit tests for CI workflow template rendering
- [ ] 4.2 Add unit tests for CI workflow file existence detection
- [ ] 4.3 Add integration test verifying CI workflow creation during init

## 5. Validation

- [ ] 5.1 Run `spectr validate add-ci-workflow-init-option --strict` and fix any issues
- [ ] 5.2 Manually test the init wizard with CI workflow option enabled and disabled
