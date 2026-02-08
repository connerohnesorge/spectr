package domain

// hookTimeout is the default timeout in seconds for hook commands.
const hookTimeout = 600

// BuildHooksFrontmatter constructs the hooks map for a slash command's
// YAML frontmatter. Each hook type gets an entry with a command that
// invokes spectr hooks.
func BuildHooksFrontmatter(cmd SlashCommand) map[string]any {
	hooks := make(map[string]any, len(AllHookTypes()))

	for _, ht := range AllHookTypes() {
		hooks[ht.String()] = []any{
			map[string]any{
				"matcher": "",
				"hooks": []any{
					map[string]any{
						"type":    "command",
						"command": "spectr hooks " + ht.String() + " --command " + cmd.String(),
						"timeout": hookTimeout,
					},
				},
			},
		}
	}

	return hooks
}
