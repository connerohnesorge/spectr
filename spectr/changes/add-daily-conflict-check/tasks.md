## 1. Implementation

- [ ] 1.1 Create conflict detection logic in `internal/conflicts/` package
  - [ ] 1.1.1 Implement function to scan all pending changes
  - [ ] 1.1.2 Implement function to extract affected capabilities from each change
  - [ ] 1.1.3 Implement function to detect requirement-level overlaps (same requirement in multiple changes)
  - [ ] 1.1.4 Generate conflict report with details (change IDs, affected specs, conflicting requirements)

- [ ] 1.2 Add CLI command `spectr conflicts` to run conflict detection
  - [ ] 1.2.1 Add command definition in cmd/
  - [ ] 1.2.2 Support `--json` output for CI consumption
  - [ ] 1.2.3 Return non-zero exit code when conflicts detected

- [ ] 1.3 Create GitHub Action workflow `.github/workflows/conflict-check.yml`
  - [ ] 1.3.1 Configure cron schedule for 5 AM UTC daily
  - [ ] 1.3.2 Checkout repository and run spectr conflicts --json
  - [ ] 1.3.3 Parse JSON output and create GitHub issue if conflicts found
  - [ ] 1.3.4 Deduplicate issues (don't create duplicates for same conflicts)

- [ ] 1.4 Write tests
  - [ ] 1.4.1 Unit tests for conflict detection logic
  - [ ] 1.4.2 Integration test with sample conflicting changes

- [ ] 1.5 Update documentation
  - [ ] 1.5.1 Add conflict detection to AGENTS.md workflow guidance
