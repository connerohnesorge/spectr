# Tasks

## 1. Skill Definition and Discovery

- [ ] 1.1 Create `.agents/skills/spectr-timeline/SKILL.md` with clear instructions
- [ ] 1.2 Include command description, examples, and expected output format
- [ ] 1.3 Document the timeline.json schema and all fields
- [ ] 1.4 Add guardrails about minimal dependencies and straightforward analysis

## 2. Timeline Analysis and JSON Generation

- [ ] 2.1 Discover all active changes in `spectr/changes/` (excluding archive/)
- [ ] 2.2 Parse proposal.md frontmatter for each change (requires/enables metadata)
- [ ] 2.3 Build dependency graph from parsed metadata
- [ ] 2.4 Detect circular dependencies and report errors
- [ ] 2.5 Calculate implementation phases (parallel vs sequential batches)
- [ ] 2.6 Generate timeline.json with:
  - [ ] 2.6a Change metadata (ID, title, description)
  - [ ] 2.6b Task counts and completion status
  - [ ] 2.6c Dependency relationships with reasons
  - [ ] 2.6d Implementation phase assignment
  - [ ] 2.6e Notes on blockers and parallelization opportunities
- [ ] 2.7 Ensure human-readable formatting (well-indented JSON, explanatory fields)

## 3. Testing and Validation

- [ ] 3.1 Create sample timeline.json output manually to verify structure
- [ ] 3.2 Test with existing chained-proposals (add-chained-proposals change)
- [ ] 3.3 Verify circular dependency detection works
- [ ] 3.4 Verify implementation phases are correctly ordered
- [ ] 3.5 Validate JSON output against expected schema

## 4. Documentation

- [ ] 4.1 Update `/spectr/AGENTS.md` to document /spectr:timeline command
- [ ] 4.2 Include example timeline.json output in documentation
- [ ] 4.3 Explain phase-based implementation strategy
- [ ] 4.4 Document how to interpret timeline for planning purposes

## 5. Integration

- [ ] 5.1 Test that timeline.json is generated correctly in project root
- [ ] 5.2 Verify skill is discoverable by claude-code CLI
- [ ] 5.3 Confirm timeline output includes all active changes
- [ ] 5.4 Add timeline.json to .gitignore if it's generated (or track it if it's reference docs)
