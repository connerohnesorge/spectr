# Implementation Tasks

## 1. Implementation

- [x] 1.1 Add `configuredProviders` map to `WizardModel` struct in
  `internal/init/wizard.go`
- [x] 1.2 Update `NewWizardModel()` to call `IsConfigured()` on each provider
  and populate the map
- [x] 1.3 Update `NewWizardModel()` to pre-select already-configured providers
  in `selectedProviders` map
- [x] 1.4 Update `renderProviderGroup()` to display a "configured" indicator for
  already-configured providers
- [x] 1.5 Update `renderSelect()` help text to explain the configured indicator

## 2. Testing

- [x] 2.1 Add test case for `NewWizardModel` with configured providers
- [x] 2.2 Add test case for `NewWizardModel` with no configured providers
- [x] 2.3 Verify visual rendering shows configured indicator correctly

## 3. Validation

- [x] 3.1 Run `spectr validate add-configured-provider-detection --strict`
- [ ] 3.2 Manual testing with fresh project (no configured providers)
- [ ] 3.3 Manual testing with project that has some providers configured
