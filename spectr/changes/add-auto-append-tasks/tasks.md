# Tasks

## 1. Configuration Infrastructure

- [ ] 1.1 Create `internal/config/config.go` with YAML config struct
- [ ] 1.2 Implement `LoadConfig()` function to find and parse `spectr.yaml`
- [ ] 1.3 Add config file discovery (walk up from cwd to find spectr.yaml)
- [ ] 1.4 Write unit tests for config loading (valid, missing, malformed)

## 2. Accept Command Integration

- [ ] 2.1 Load config in `AcceptCmd.Run()` before processing
- [ ] 2.2 Pass append_tasks config to task processing pipeline
- [ ] 2.3 Modify `writeTasksJSONC()` to append configured tasks with section
- [ ] 2.4 Generate sequential task IDs for appended tasks (continue from last)
- [ ] 2.5 Write integration tests for accept with config

## 3. Validation and Documentation

- [ ] 3.1 Run existing tests to ensure no regressions
- [ ] 3.2 Test behavior when spectr.yaml is missing (should work unchanged)
- [ ] 3.3 Update CLI help text if needed
