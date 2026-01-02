## 1. Conclaude Schema Extension

- [ ] 1.1 Add `postToolUse` section to conclaude JSON schema
- [ ] 1.2 Define `commands` array structure with tool filter support
- [ ] 1.3 Add environment variable documentation to schema

## 2. Hook Command Configuration

- [ ] 2.1 Implement tool name matching (exact, glob patterns)
- [ ] 2.2 Implement command execution with environment variables
- [ ] 2.3 Add timeout and error handling configuration

## 3. Environment Variable Interface

- [ ] 3.1 Define CONCLAUDE_TOOL_NAME variable
- [ ] 3.2 Define CONCLAUDE_TOOL_INPUT variable (JSON)
- [ ] 3.3 Define CONCLAUDE_TOOL_OUTPUT variable (JSON)
- [ ] 3.4 Define CONCLAUDE_TOOL_TIMESTAMP variable

## 4. Integration Examples

- [ ] 4.1 Create AskUserQuestion Q&A logger example script
- [ ] 4.2 Document single append file logging pattern
- [ ] 4.3 Add example to .conclaude.yaml template

## 5. Validation

- [ ] 5.1 Update conclaude validate to check postToolUse config
- [ ] 5.2 Add tests for hook execution flow
- [ ] 5.3 Test environment variable escaping edge cases
