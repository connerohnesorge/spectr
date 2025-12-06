# Change: Add CI Workflow Setup Option to Init Wizard

## Why
Users currently need to manually set up the `.github/workflows/spectr-ci.yml` file to enable automated Spectr validation in their CI/CD pipeline. This is an extra step that many users may not know about or may forget. Adding a checkbox in the init wizard's TUI to optionally create this workflow file provides a seamless onboarding experience and encourages best practices for CI integration.

## What Changes
- Add a new step or option in the init wizard TUI to ask users if they want to set up GitHub Actions CI workflow
- When selected, create `.github/workflows/spectr-ci.yml` with the standard Spectr validation job configuration
- The option should be presented clearly with a brief explanation of what it does
- The workflow file should follow the existing ci-integration spec requirements (full history checkout, concurrency management, etc.)
- Skip this option if `.github/workflows/spectr-ci.yml` already exists (treat as already configured)

## Impact
- Affected specs: `cli-interface` (MODIFIED - init wizard functionality)
- Affected code:
  - `internal/initialize/wizard.go` - Add CI workflow step or option to the TUI
  - `internal/initialize/executor.go` - Add CI workflow file creation logic
  - `internal/initialize/templates.go` - Add CI workflow template
- Breaking changes: None (additive feature)
- Benefits:
  - Reduces onboarding friction for users who want CI validation
  - Promotes best practices by making CI setup visible and easy
  - Maintains consistency with existing init wizard patterns
