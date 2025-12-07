# Change: Add CI Workflow Setup Option to Init Wizard

## Why
Users currently need to manually set up the `.github/workflows/spectr-ci.yml` file to enable automated Spectr validation in their CI/CD pipeline. This is an extra step that many users may not know about or may forget. Adding an option in the init wizard to optionally create this workflow file provides a seamless onboarding experience and encourages best practices for CI integration.

## What Changes
- Add a "Spectr CI Validation" checkbox option to the Review step of the init wizard
- The option appears after the tool selection summary, before the creation plan
- When enabled, create `.github/workflows/spectr-ci.yml` with Spectr validation job only
- The workflow uses `spectr-action@v0.0.2` and triggers on push to main and all PRs
- Pre-select if `.github/workflows/spectr-ci.yml` already exists (treat as configured)
- Add `--ci-workflow` flag for non-interactive mode

## Design Decisions
- **Combined with Review step**: Avoids adding another step, keeps wizard flow quick
- **Workflow name**: `spectr-ci.yml` - dedicated Spectr workflow file that won't conflict with existing CI
- **Workflow scope**: Spectr validation only (single `spectr-validate` job) - minimal and focused
- **Action version**: `@v0.0.2` - pinned version for reproducibility
- **Push triggers**: main only (PRs trigger on all branches)
- **Label**: "Spectr CI Validation" - focuses on function rather than platform
- **Directory handling**: Create file normally even if `.github/workflows/` already exists

## Impact
- Affected specs: `cli-interface` (MODIFIED - init wizard functionality)
- Affected code:
  - `internal/initialize/wizard.go` - Add CI option to Review step rendering and state
  - `internal/initialize/executor.go` - Add CI workflow file creation logic
  - `internal/initialize/templates.go` - Add CI workflow template
  - `cmd/init.go` - Add `--ci-workflow` flag
- Breaking changes: None (additive feature)
- Benefits:
  - Reduces onboarding friction for users who want CI validation
  - Promotes best practices by making CI setup visible and easy
  - No additional wizard step keeps the flow quick
  - Maintains consistency with existing init wizard patterns
