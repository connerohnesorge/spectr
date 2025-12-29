# Change: Show Already-Configured Providers in Init Wizard TUI

## Why

When running `spectr init` on a project that already has some providers configured, users have no visual indication of which tools are already set up. This forces users to either:

1. Manually check file existence before running init
2. Accidentally re-select already-configured providers (harmless but confusing)
3. Miss the opportunity to understand their current configuration state

The `IsConfigured()` method exists on all providers but isn't used in the wizard TUI.

## What Changes

- **Detect configured providers**: Call `IsConfigured(projectPath)` on each provider during wizard initialization
- **Visual indicator**: Show a distinct marker (e.g., `[✓ configured]` or different color) for already-configured providers
- **Pre-selection behavior**: Already-configured providers are pre-selected by default, allowing users to keep or deselect them
- **Legend/help text**: Add explanation of the configured indicator in the selection screen

## Impact

- Affected specs: `cli-interface` (adds requirement for configured provider display)
- Affected code:
  - `internal/init/wizard.go` - WizardModel initialization and rendering
  - `internal/init/wizard.go` - `renderProviderGroup()` and `renderSelect()` functions
