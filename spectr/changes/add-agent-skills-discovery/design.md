# Design: Agent Skills Discovery

## Architecture Overview

This design extends Spectr's provider system with skill discovery capabilities by adding a new `internal/discovery` package for skill metadata parsing and integrating with the existing `internal/list` and `cmd/list` infrastructure.

### Design Principles

1. **Simplicity First**: <500 lines total, single-file implementations
2. **Follow Existing Patterns**: Mirror `internal/list` architecture for consistency
3. **Embedded-Only**: Scan only embedded skills, defer installed skill discovery
4. **Read-Only**: No state changes, purely informational queries
5. **Standard Formats**: Use Hugo/Jekyll-style YAML frontmatter

## Package Structure

```
internal/
├── discovery/
│   ├── skills.go         # NEW: Skill discovery functions (100 lines)
│   ├── skills_test.go    # NEW: Discovery tests (150 lines)
│   └── doc.go            # NEW: Package documentation
├── list/
│   ├── types.go          # MODIFIED: Add SkillInfo struct (15 lines)
│   └── lister.go         # MODIFIED: Add ListSkills() method (30 lines)
└── cmd/
    └── list.go           # MODIFIED: Add --skills flag and handler (60 lines)
```

## Data Structures

### SkillMetadata (internal/discovery/skills.go)

Represents raw SKILL.md frontmatter before enrichment:

```go
// SkillMetadata represents the frontmatter in SKILL.md files.
// Fields map directly to YAML frontmatter structure following agentskills.io spec.
type SkillMetadata struct {
    Name         string `yaml:"name"`         // Required: skill identifier
    Description  string `yaml:"description"`  // Required: short description
    Compatibility struct {
        Requirements []string `yaml:"requirements"` // Optional: required dependencies
        Optional     []string `yaml:"optional"`     // Optional: optional dependencies
        Platforms    []string `yaml:"platforms"`    // Optional: supported platforms
    } `yaml:"compatibility"`
}
```

**Why separate from SkillInfo?**
- `SkillMetadata`: Pure YAML parsing, no side effects
- `SkillInfo`: Enriched with installation status from filesystem
- Separation of concerns: parsing vs enrichment

### SkillInfo (internal/list/types.go)

Public-facing skill information with installation status:

```go
// SkillInfo represents information about an Agent Skill.
// Extends SkillMetadata with installation status for user display.
type SkillInfo struct {
    Name         string   `json:"name"`         // Skill name from frontmatter
    Description  string   `json:"description"`  // Skill description
    Installed    bool     `json:"installed"`    // Is it installed in project?
    Requirements []string `json:"requirements"` // Required dependencies (jq, bash, etc.)
    Optional     []string `json:"optional"`     // Optional dependencies
    Platforms    []string `json:"platforms"`    // Supported platforms (Linux, macOS, etc.)
}
```

**JSON tags** enable `--json` output format for scripting/automation.

## Core Functions

### 1. ParseSkillMetadata (internal/discovery/skills.go)

**Purpose:** Extract YAML frontmatter from SKILL.md content

**Signature:**
```go
func ParseSkillMetadata(content []byte) (*SkillMetadata, error)
```

**Algorithm:**
```
1. Split content into lines
2. Find first "---" marker (start of frontmatter)
3. Find second "---" marker (end of frontmatter)
4. Extract lines between markers
5. Parse YAML using gopkg.in/yaml.v3
6. Validate required fields (name, description)
7. Return SkillMetadata or error
```

**Error Handling:**
- No frontmatter markers → `ErrNoFrontmatter`
- Invalid YAML → `ErrInvalidYAML` with line context
- Missing required fields → `ErrMissingRequiredField`

**Example Input:**
```yaml
---
name: spectr-accept-wo-spectr-bin
description: Accept Spectr change proposals without spectr binary
compatibility:
  requirements:
    - jq
  platforms:
    - Linux
    - macOS
---

# Skill documentation...
```

### 2. ListEmbeddedSkills (internal/discovery/skills.go)

**Purpose:** Discover all embedded skills from skillFS

**Signature:**
```go
func ListEmbeddedSkills(skillFS fs.FS) ([]SkillMetadata, error)
```

**Algorithm:**
```
1. Read root directory entries from skillFS
2. For each entry:
   a. Skip if not a directory
   b. Check if SKILL.md exists in directory
   c. If missing, log warning and skip
   d. Read SKILL.md content
   e. Parse frontmatter via ParseSkillMetadata()
   f. Append to results
3. Sort results alphabetically by name
4. Return slice of SkillMetadata
```

**Why skillFS parameter?**
- Testable: Can inject mock filesystem
- Flexible: Works with any fs.FS implementation
- Consistent: Matches TemplateManager.SkillFS() pattern

### 3. IsSkillInstalled (internal/discovery/skills.go)

**Purpose:** Check if a skill is installed in project

**Signature:**
```go
func IsSkillInstalled(projectPath, skillName string) (bool, error)
```

**Algorithm:**
```
1. Construct path: <projectPath>/.claude/skills/<skillName>/SKILL.md
2. Check file existence using os.Stat
3. Return true if exists, false if not
4. Return error only for unexpected filesystem errors
```

**Design Decision: Why check SKILL.md instead of directory?**
- More reliable: Directory could exist but be empty
- Matches AgentSkillsInitializer.IsSetup() pattern
- Simpler: Single file check vs directory walking

## List Integration

### ListSkills Method (internal/list/lister.go)

**Purpose:** Combine skill discovery with installation status

**Signature:**
```go
func (l *Lister) ListSkills(skillFS fs.FS) ([]SkillInfo, error)
```

**Algorithm:**
```
1. Call discovery.ListEmbeddedSkills(skillFS) → []SkillMetadata
2. For each metadata:
   a. Call discovery.IsSkillInstalled(l.projectPath, meta.Name) → installed bool
   b. Convert SkillMetadata → SkillInfo with installed field
   c. Append to results
3. Return []SkillInfo (already sorted by ListEmbeddedSkills)
```

**Why method on Lister?**
- Consistent with ListChanges(), ListSpecs(), ListAll()
- Access to `projectPath` field for installation checking
- Maintains existing architecture patterns

## CLI Command Design

### Flag Addition (cmd/list.go)

Extend `ListCmd` struct:

```go
type ListCmd struct {
    Specs       bool `name:"specs" help:"List specifications instead of changes"`
    All         bool `name:"all"   help:"List both changes and specs in unified mode"`
    Skills      bool `name:"skills" help:"List embedded Agent Skills"` // NEW

    Long        bool `name:"long" help:"Show detailed output with titles and counts"`
    JSON        bool `name:"json" help:"Output as JSON"`
    Interactive bool `name:"interactive" help:"Interactive mode" short:"I"`
    Stdout      bool `name:"stdout" help:"Print ID to stdout (requires -I)"`
}
```

### Mutual Exclusivity Validation

```go
// In Run() method, add:
if c.Skills && (c.Specs || c.All || c.Interactive) {
    return &specterrs.IncompatibleFlagsError{
        Flag1: "--skills",
        Flag2: determineConflictingFlag(c), // --specs, --all, or --interactive
    }
}
```

**Rationale:** Skills are fundamentally different from changes/specs:
- Changes/specs are project-specific, skills are embedded/global
- Interactive mode designed for clipboard selection, skills are informational
- Prevents user confusion about what's being listed

### listSkills Method (cmd/list.go)

```go
func (c *ListCmd) listSkills(lister *list.Lister) error {
    // 1. Get skillFS from TemplateManager
    tm, err := initialize.NewTemplateManager()
    if err != nil {
        return fmt.Errorf("failed to create template manager: %w", err)
    }

    // 2. Get root skills filesystem
    skillFS, err := tm.SkillFS("")
    if err != nil {
        return fmt.Errorf("failed to access skills: %w", err)
    }

    // 3. List skills with installation status
    skills, err := lister.ListSkills(skillFS)
    if err != nil {
        return fmt.Errorf("failed to list skills: %w", err)
    }

    // 4. Format and output
    output, err := c.formatSkillsOutput(skills)
    if err != nil {
        return fmt.Errorf("failed to format output: %w", err)
    }

    fmt.Println(output)
    return nil
}
```

## Output Formats

### Text Format (Default)

**Format:** One skill name per line, sorted alphabetically

**Example:**
```
spectr-accept-wo-spectr-bin
spectr-validate-wo-spectr-bin
```

**Implementation:**
```go
func formatSkillsText(skills []SkillInfo) string {
    names := make([]string, len(skills))
    for i, skill := range skills {
        names[i] = skill.Name
    }
    return strings.Join(names, "\n")
}
```

### Long Format (--long)

**Format:** Multi-line per skill with description, requirements, platforms, status

**Example:**
```
spectr-accept-wo-spectr-bin
  Accept Spectr change proposals without spectr binary
  Requirements: jq
  Platforms: Linux, macOS, Unix-like systems with bash
  Status: Not Installed

spectr-validate-wo-spectr-bin
  Validate Spectr specs without spectr binary
  Requirements: bash 4.0+, grep, sed, find
  Optional: jq (for JSON output)
  Platforms: Linux, macOS, Unix-like systems with bash
  Status: Installed
```

**Implementation:**
```go
func formatSkillsLong(skills []SkillInfo) string {
    var buf strings.Builder
    for i, skill := range skills {
        if i > 0 {
            buf.WriteString("\n\n")
        }
        buf.WriteString(skill.Name)
        buf.WriteString("\n  ")
        buf.WriteString(skill.Description)

        if len(skill.Requirements) > 0 {
            buf.WriteString("\n  Requirements: ")
            buf.WriteString(strings.Join(skill.Requirements, ", "))
        }

        if len(skill.Optional) > 0 {
            buf.WriteString("\n  Optional: ")
            buf.WriteString(strings.Join(skill.Optional, ", "))
        }

        if len(skill.Platforms) > 0 {
            buf.WriteString("\n  Platforms: ")
            buf.WriteString(strings.Join(skill.Platforms, ", "))
        }

        buf.WriteString("\n  Status: ")
        if skill.Installed {
            buf.WriteString("Installed")
        } else {
            buf.WriteString("Not Installed")
        }
    }
    return buf.String()
}
```

### JSON Format (--json)

**Format:** JSON array of SkillInfo objects

**Example:**
```json
[
  {
    "name": "spectr-accept-wo-spectr-bin",
    "description": "Accept Spectr change proposals without spectr binary",
    "installed": false,
    "requirements": ["jq"],
    "optional": [],
    "platforms": ["Linux", "macOS", "Unix-like systems with bash"]
  },
  {
    "name": "spectr-validate-wo-spectr-bin",
    "description": "Validate Spectr specs without spectr binary",
    "installed": true,
    "requirements": ["bash 4.0+", "grep", "sed", "find"],
    "optional": ["jq"],
    "platforms": ["Linux", "macOS", "Unix-like systems with bash"]
  }
]
```

**Implementation:**
```go
func formatSkillsJSON(skills []SkillInfo) (string, error) {
    data, err := json.MarshalIndent(skills, "", "  ")
    if err != nil {
        return "", fmt.Errorf("failed to marshal JSON: %w", err)
    }
    return string(data), nil
}
```

## Testing Strategy

### Unit Tests (internal/discovery/skills_test.go)

**Test Cases:**

1. **TestParseSkillMetadata_ValidFrontmatter**
   - Input: Complete SKILL.md with all fields
   - Assert: All fields correctly parsed, no errors

2. **TestParseSkillMetadata_MinimalFrontmatter**
   - Input: Only required fields (name, description)
   - Assert: Optional fields are empty slices, no errors

3. **TestParseSkillMetadata_MissingRequiredField**
   - Input: Missing `name` field
   - Assert: Error with clear message about missing field

4. **TestParseSkillMetadata_InvalidYAML**
   - Input: Malformed YAML (invalid indentation, syntax errors)
   - Assert: Error with YAML parser context

5. **TestParseSkillMetadata_NoFrontmatter**
   - Input: SKILL.md without `---` markers
   - Assert: ErrNoFrontmatter with helpful message

6. **TestListEmbeddedSkills_MultipleSkills**
   - Setup: Mock fs.FS with 2+ skill directories
   - Assert: Returns all skills, sorted alphabetically

7. **TestListEmbeddedSkills_MissingSkillMd**
   - Setup: Directory without SKILL.md
   - Assert: Skips directory, logs warning, continues

8. **TestListEmbeddedSkills_EmptyDirectory**
   - Setup: skillFS with no skill directories
   - Assert: Returns empty slice, no error

9. **TestIsSkillInstalled_Installed**
   - Setup: Create .claude/skills/<name>/SKILL.md
   - Assert: Returns true, no error

10. **TestIsSkillInstalled_NotInstalled**
    - Setup: No .claude/skills directory
    - Assert: Returns false, no error

**Test Data Structure:**
```
internal/discovery/testdata/
├── valid_skill/
│   └── SKILL.md                    # Complete frontmatter
├── minimal_skill/
│   └── SKILL.md                    # Only required fields
├── invalid_yaml/
│   └── SKILL.md                    # Malformed YAML
├── no_frontmatter/
│   └── SKILL.md                    # No --- markers
└── missing_name/
    └── SKILL.md                    # Missing required field
```

### Integration Tests (cmd/list_test.go)

**Test Cases:**

1. **TestListCmd_Skills_TextFormat**
   - Execute: `spectr list --skills`
   - Assert: Outputs skill names, one per line

2. **TestListCmd_Skills_LongFormat**
   - Execute: `spectr list --skills --long`
   - Assert: Multi-line output with descriptions, requirements, status

3. **TestListCmd_Skills_JSONFormat**
   - Execute: `spectr list --skills --json`
   - Assert: Valid JSON array, all fields present

4. **TestListCmd_Skills_IncompatibleWithSpecs**
   - Execute: `spectr list --skills --specs`
   - Assert: Returns IncompatibleFlagsError

5. **TestListCmd_Skills_IncompatibleWithAll**
   - Execute: `spectr list --skills --all`
   - Assert: Returns IncompatibleFlagsError

6. **TestListCmd_Skills_IncompatibleWithInteractive**
   - Execute: `spectr list --skills --interactive`
   - Assert: Returns IncompatibleFlagsError

## Implementation Sequence

### Phase 1: Core Discovery (Day 1)

**Files:**
- `internal/discovery/skills.go`
- `internal/discovery/skills_test.go`
- `internal/discovery/doc.go`

**Tasks:**
1. Create `SkillMetadata` struct
2. Implement `ParseSkillMetadata()` with error handling
3. Implement `ListEmbeddedSkills()` with directory walking
4. Implement `IsSkillInstalled()` with filesystem checks
5. Write unit tests for all functions (10 test cases)
6. Ensure test coverage >80%

**Validation:** All discovery tests pass

### Phase 2: List Integration (Day 1-2)

**Files:**
- `internal/list/types.go`
- `internal/list/lister.go`

**Tasks:**
1. Add `SkillInfo` struct to types.go
2. Add `ListSkills()` method to lister.go
3. Test skill listing with mock skillFS
4. Test installation status enrichment

**Validation:** List integration tests pass

### Phase 3: CLI Command (Day 2)

**Files:**
- `cmd/list.go`
- `cmd/list_test.go`

**Tasks:**
1. Add `Skills` field to ListCmd
2. Add mutual exclusivity validation
3. Implement `listSkills()` method
4. Implement output formatters (text, long, JSON)
5. Wire up in Run() method
6. Write CLI integration tests (6 test cases)

**Validation:** CLI tests pass, manual smoke testing

### Phase 4: Specification (Day 2)

**Files:**
- `spectr/changes/add-agent-skills-discovery/specs/provider-system/spec.md`

**Tasks:**
1. Draft delta spec with ADDED requirements
2. Write scenarios for each requirement (8 scenarios total)
3. Cross-reference existing AgentSkillsInitializer requirements
4. Run `spectr validate add-agent-skills-discovery`
5. Fix any validation errors

**Validation:** `spectr validate` passes

### Phase 5: Final Testing (Day 3)

**Tasks:**
1. Run full test suite: `nix develop -c tests`
2. Run linter: `nix develop -c lint`
3. Manual testing with real embedded skills
4. Test all output formats with edge cases
5. Verify code size <500 lines

**Validation:** All tests pass, no lint errors

## Error Handling

### Discovery Errors

```go
var (
    ErrNoFrontmatter = errors.New("SKILL.md missing frontmatter markers (---)")
    ErrInvalidYAML = errors.New("invalid YAML in SKILL.md frontmatter")
    ErrMissingRequiredField = errors.New("missing required frontmatter field")
    ErrSkillNotFound = errors.New("skill not found in embedded filesystem")
)
```

**Error Wrapping:**
```go
if err := yaml.Unmarshal(yamlContent, &meta); err != nil {
    return nil, fmt.Errorf("%w: %v", ErrInvalidYAML, err)
}
```

### CLI Errors

- Follow existing error patterns from `cmd/list.go`
- Use `specterrs.IncompatibleFlagsError` for flag conflicts
- Return wrapped errors with context: `fmt.Errorf("failed to list skills: %w", err)`

## Design Trade-offs

### Embedded-Only vs All Sources

**Decision:** Embedded-only for initial implementation

**Rationale:**
- Simpler: No need to scan multiple directories or handle conflicts
- Sufficient: Current use case is only embedded skills (spectr-accept, spectr-validate)
- Extensible: Can add installed skill discovery in future without breaking changes

**Future Work:** Add `ListInstalledSkills()` and merge with embedded results

### Read-Only vs Write Operations

**Decision:** Read-only discovery, no install/remove commands

**Rationale:**
- Lower complexity: No error recovery, rollback, or validation needed
- Aligns with `list` command philosophy (informational, no side effects)
- Reduces testing surface: No need for filesystem modification tests

**Future Work:** Add `spectr install skill <name>` and `spectr remove skill <name>` commands

### YAML vs Other Formats

**Decision:** YAML frontmatter (Hugo/Jekyll style)

**Rationale:**
- Industry standard: Used by Hugo, Jekyll, OpenAI Codex, GitHub Actions
- Already embedded: `gopkg.in/yaml.v3` already in go.mod for other purposes
- Human-readable: Easier to write and edit than JSON or TOML
- Separator clarity: `---` markers unambiguous

**Alternative Considered:** JSON frontmatter rejected due to lack of standard marker syntax

## Dependencies

All dependencies already in `go.mod`:

- `gopkg.in/yaml.v3` v3.0.1 - YAML parsing (already used for other features)
- `github.com/spf13/afero` v1.11.0 - Filesystem abstraction (already used)
- `io/fs` - Standard library, no version

**No new dependencies required** ✅

## Code Size Estimates

| Component | File | Lines | Complexity |
|-----------|------|-------|------------|
| SkillMetadata struct | skills.go | 15 | Low |
| ParseSkillMetadata() | skills.go | 40 | Medium (YAML) |
| ListEmbeddedSkills() | skills.go | 30 | Low |
| IsSkillInstalled() | skills.go | 15 | Low |
| **Discovery Package** | | **100** | **Medium** |
| SkillInfo struct | types.go | 15 | Low |
| ListSkills() method | lister.go | 30 | Low |
| **List Integration** | | **45** | **Low** |
| Skills field + validation | list.go | 15 | Low |
| listSkills() method | list.go | 25 | Low |
| Output formatters | list.go | 20 | Low |
| **CLI Command** | | **60** | **Low** |
| Discovery unit tests | skills_test.go | 150 | Medium |
| CLI integration tests | list_test.go | 50 | Low |
| **Tests** | | **200** | **Medium** |
| Delta spec scenarios | spec.md | 50 | Low |
| **Total** | | **455** | **<500 ✅** |

## Validation Checklist

- [ ] All requirements have WHEN/THEN scenarios
- [ ] Tests verify actual behavior (not mocked)
- [ ] CLI flags mutually exclusive (--skills vs --specs/--all)
- [ ] JSON output is valid and parseable
- [ ] Frontmatter parsing handles missing optional fields gracefully
- [ ] Installation check uses afero (filesystem-independent for testing)
- [ ] Code size <500 lines total
- [ ] No new external dependencies
- [ ] Follows existing architectural patterns (list, discovery)
- [ ] Error messages are clear and actionable
