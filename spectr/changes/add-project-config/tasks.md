## 1. Core Configuration Infrastructure

- [ ] 1.1 Create `internal/config/` package with `Config` struct and `Load()` function
- [ ] 1.2 Implement config file discovery (walk up directories for `spectr.yaml`)
- [ ] 1.3 Add YAML parsing with gopkg.in/yaml.v3 dependency
- [ ] 1.4 Implement default values when config file is absent
- [ ] 1.5 Write unit tests for config loading and defaults

## 2. Migrate Discovery Module

- [ ] 2.1 Update `discovery/changes.go` to use config for root directory
- [ ] 2.2 Update `discovery/specs.go` to use config for root directory
- [ ] 2.3 Add config parameter to discovery functions or use context pattern
- [ ] 2.4 Update discovery tests to cover custom root directories

## 3. Migrate Validation Module

- [ ] 3.1 Remove hardcoded `SpectrDir` constant from `validation/helpers.go`
- [ ] 3.2 Update `ValidateItemByType` to use config
- [ ] 3.3 Update validation tests for custom root directories

## 4. Migrate Archive Module

- [ ] 4.1 Update `archive/archiver.go` path construction to use config
- [ ] 4.2 Update archive tests for custom root directories

## 5. Migrate List Module

- [ ] 5.1 Update `list/lister.go` path construction to use config
- [ ] 5.2 Update list tests for custom root directories

## 6. Migrate View Module

- [ ] 6.1 Update `view/dashboard.go` path construction to use config
- [ ] 6.2 Update view tests for custom root directories

## 7. Update Init Command

- [ ] 7.1 Update `init/executor.go` to support custom root directory
- [ ] 7.2 Update `init/filesystem.go` `IsSpectrInitialized` to use config
- [ ] 7.3 Optionally create `spectr.yaml` during init (when non-default root requested)
- [ ] 7.4 Update init tests

## 8. Integration Testing

- [ ] 8.1 Add integration tests for projects with custom root directories
- [ ] 8.2 Test backward compatibility with projects lacking `spectr.yaml`
- [ ] 8.3 Test config file discovery from nested directories

## 9. Documentation

- [ ] 9.1 Update README with configuration options
- [ ] 9.2 Update AGENTS.md template with config information
