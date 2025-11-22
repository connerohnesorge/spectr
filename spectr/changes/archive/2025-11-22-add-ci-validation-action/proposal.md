# Change: Add GitHub Action for CI Validation

## Why
Spectr changes need to be validated automatically in CI/CD pipelines to catch validation errors before code is merged. Manual validation is error-prone and slows down development. A GitHub Action provides automated validation on every push and pull request, ensuring all Spectr changes meet quality standards before they reach the main branch.

## What Changes
- Added `spectr-validate` job to `.github/workflows/ci.yml`
- Integrated `connerohnesorge/spectr-action@v0.0.1` for automated validation
- Configured job to run on all branches for pushes and pull requests
- Job runs with full git history (`fetch-depth: 0`) to support change detection
- Positioned as first job in CI pipeline to fail fast on spec violations

## Impact
- Affected specs: `ci-integration` (new capability)
- Affected code: `.github/workflows/ci.yml` (lines 15-24)
- Breaking changes: None
- Benefits:
  - Automated validation prevents invalid specs from being merged
  - Developers get immediate feedback on proposal quality
  - Reduces manual review burden for spec correctness
  - Enforces Spectr conventions across all contributions
