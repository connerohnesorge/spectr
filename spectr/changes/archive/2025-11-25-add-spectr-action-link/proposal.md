# Change: Add Spectr Action Documentation Link

## Why

Users who want to integrate Spectr validation into their CI/CD pipelines need an easy way to discover and use the spectr-action GitHub Action. While the action is already used in this project's CI workflow, there is no user-facing documentation pointing users to the action repository or explaining how to add it to their own projects.

## What Changes

- Add reference to `connerohnesorge/spectr-action` in README.md Links & Resources section
- Add brief CI Integration section in README.md explaining how to use the action
- Update documentation specification to include CI integration documentation requirements

## Impact

- Affected specs: `documentation`
- Affected code: `README.md` (Links & Resources section and new CI Integration section)
- Breaking changes: None
- Benefits:
  - Users can easily discover the spectr-action for their own projects
  - Clear guidance on how to add automated validation to CI pipelines
  - Improved discoverability of the ecosystem around Spectr
