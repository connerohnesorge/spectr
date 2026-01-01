# Validation Package

Enforces spec and change formatting rules. Critical quality gate.

## OVERVIEW
Validates specs (`spectr/specs/`) and changes (`spectr/changes/`). Ensures scenarios exist, headers follow format, MODIFIED requirements include complete content.

## STRUCTURE
```
internal/validation/
├── validator.go           # Main validation orchestration
├── spec_validators.go    # Spec-specific rules
├── delta_validators.go   # Delta spec rules
├── change_rules.go       # Change directory rules
├── formatters.go         # Error formatting
├── constants.go          # Markdown formatting constants
└── *_test.go            # Table-driven tests
```

## WHERE TO LOOK
| Task | Location | Notes |
|------|----------|-------|
| Validate specs | ValidateSpec() | Spec-level rules |
| Validate changes | ValidateChange() | Change + delta rules |
| Check scenarios | RequirementScenarios rule | Every requirement must have ≥1 scenario |
| Format headers | ScenarioFormatting rule | Must use `#### Scenario:` (4 hashtags) |

## CONVENTIONS
- **Strict validation**: All issues are errors (validation is strict, no warnings in strict mode)
- **Early return**: Return on first error in critical paths
- **Table tests**: All validators use t.Run() subtests

## RULES (ENFORCED)
| Rule | Severity | Description |
|------|----------|-------------|
| RequirementScenarios | Error | Every requirement MUST have ≥1 scenario |
| ScenarioFormatting | Error | Scenarios MUST use `#### Scenario:` (4 hashtags) |
| PurposeLength | Warning | Purpose sections MUST be ≥50 chars |
| ModifiedComplete | Error | MODIFIED requirements MUST include full updated content (no partial) |
| DeltaPresence | Error | Changes MUST have ≥1 delta spec |
| ScenarioStructure | Warning | Scenarios SHOULD have WHEN/THEN bullets |

## ANTI-PATTERNS
- **NEVER relax validation**: Quality gate intentional
- **DON'T skip scenarios**: Format requirement → error, fix source

## KEY FUNCTIONS
- `ValidateSpec(spec *Spec) []ValidationIssue` - Validates spec structure
- `ValidateChange(change *Change) []ValidationIssue` - Validates change + delta specs
- `checkScenarios(requirements []Requirement) []ValidationIssue` - Enforces scenario rule
- `checkScenarioFormatting(scenarios []Scenario) []ValidationIssue` - Enforces header format

## ERROR FORMATTING
- `ValidationIssue{File, Line, Rule, Message, Severity}` - Structured errors
- Formatters output human-readable or JSON
- Line numbers: 1-based, point to issue location
