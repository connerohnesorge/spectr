package init

// ============================================================================
// Composite Tool Providers
// ============================================================================
// Composite providers combine memory file and slash command providers using
// Go embedded fields for composition. Each composite provider represents a
// complete tool integration with both memory files and slash commands.

// ClaudeCodeToolProvider combines Claude memory file and slash command
// providers
type ClaudeCodeToolProvider struct {
	*ClaudeMemoryFileProvider
	*ClaudeSlashCommandProvider
	*SpectrAgentsUpdater
}

// NewClaudeCodeToolProvider creates a new Claude Code tool provider
func NewClaudeCodeToolProvider() *ClaudeCodeToolProvider {
	return &ClaudeCodeToolProvider{
		ClaudeMemoryFileProvider:   &ClaudeMemoryFileProvider{},
		ClaudeSlashCommandProvider: NewClaudeSlashCommandProvider(),
		SpectrAgentsUpdater:        &SpectrAgentsUpdater{},
	}
}

// GetName returns the tool name
func (*ClaudeCodeToolProvider) GetName() string {
	return "Claude Code"
}

// GetMemoryFileProvider returns the Claude memory file provider
func (p *ClaudeCodeToolProvider) GetMemoryFileProvider() MemoryFileProvider {
	return p.ClaudeMemoryFileProvider
}

// GetSlashCommandProvider returns the Claude slash command provider
//
//nolint:revive // line-length-limit - method name length required by interface
func (p *ClaudeCodeToolProvider) GetSlashCommandProvider() SlashCommandProvider {
	return p.ClaudeSlashCommandProvider
}

// ClineToolProvider combines Cline memory file and slash command providers
type ClineToolProvider struct {
	*ClineMemoryFileProvider
	*ClineSlashCommandProvider
	*SpectrAgentsUpdater
}

// NewClineToolProvider creates a new Cline tool provider
func NewClineToolProvider() *ClineToolProvider {
	return &ClineToolProvider{
		ClineMemoryFileProvider:   &ClineMemoryFileProvider{},
		ClineSlashCommandProvider: NewClineSlashCommandProvider(),
		SpectrAgentsUpdater:       &SpectrAgentsUpdater{},
	}
}

// GetName returns the tool name
func (*ClineToolProvider) GetName() string {
	return "Cline"
}

// GetMemoryFileProvider returns the Cline memory file provider
func (p *ClineToolProvider) GetMemoryFileProvider() MemoryFileProvider {
	return p.ClineMemoryFileProvider
}

// GetSlashCommandProvider returns the Cline slash command provider
func (p *ClineToolProvider) GetSlashCommandProvider() SlashCommandProvider {
	return p.ClineSlashCommandProvider
}

// QoderToolProvider combines Qoder memory file and slash command providers
type QoderToolProvider struct {
	*QoderMemoryFileProvider
	*QoderSlashCommandProvider
	*SpectrAgentsUpdater
}

// NewQoderToolProvider creates a new Qoder tool provider
func NewQoderToolProvider() *QoderToolProvider {
	return &QoderToolProvider{
		QoderMemoryFileProvider:   &QoderMemoryFileProvider{},
		QoderSlashCommandProvider: NewQoderSlashCommandProvider(),
		SpectrAgentsUpdater:       &SpectrAgentsUpdater{},
	}
}

// GetName returns the tool name
func (*QoderToolProvider) GetName() string {
	return "Qoder"
}

// GetMemoryFileProvider returns the Qoder memory file provider
func (p *QoderToolProvider) GetMemoryFileProvider() MemoryFileProvider {
	return p.QoderMemoryFileProvider
}

// GetSlashCommandProvider returns the Qoder slash command provider
func (p *QoderToolProvider) GetSlashCommandProvider() SlashCommandProvider {
	return p.QoderSlashCommandProvider
}

// CodeBuddyToolProvider combines CodeBuddy memory file and slash command
// providers
type CodeBuddyToolProvider struct {
	*CodeBuddyMemoryFileProvider
	*CodeBuddySlashCommandProvider
	*SpectrAgentsUpdater
}

// NewCodeBuddyToolProvider creates a new CodeBuddy tool provider
func NewCodeBuddyToolProvider() *CodeBuddyToolProvider {
	return &CodeBuddyToolProvider{
		CodeBuddyMemoryFileProvider:   &CodeBuddyMemoryFileProvider{},
		CodeBuddySlashCommandProvider: NewCodeBuddySlashCommandProvider(),
		SpectrAgentsUpdater:           &SpectrAgentsUpdater{},
	}
}

// GetName returns the tool name
func (*CodeBuddyToolProvider) GetName() string {
	return "CodeBuddy"
}

// GetMemoryFileProvider returns the CodeBuddy memory file provider
func (p *CodeBuddyToolProvider) GetMemoryFileProvider() MemoryFileProvider {
	return p.CodeBuddyMemoryFileProvider
}

// GetSlashCommandProvider returns the CodeBuddy slash command provider
func (p *CodeBuddyToolProvider) GetSlashCommandProvider() SlashCommandProvider {
	return p.CodeBuddySlashCommandProvider
}

// QwenToolProvider combines Qwen memory file and slash command providers
type QwenToolProvider struct {
	*QwenMemoryFileProvider
	*QwenSlashCommandProvider
	*SpectrAgentsUpdater
}

// NewQwenToolProvider creates a new Qwen tool provider
func NewQwenToolProvider() *QwenToolProvider {
	return &QwenToolProvider{
		QwenMemoryFileProvider:   &QwenMemoryFileProvider{},
		QwenSlashCommandProvider: NewQwenSlashCommandProvider(),
		SpectrAgentsUpdater:      &SpectrAgentsUpdater{},
	}
}

// GetName returns the tool name
func (*QwenToolProvider) GetName() string {
	return "Qwen"
}

// GetMemoryFileProvider returns the Qwen memory file provider
func (p *QwenToolProvider) GetMemoryFileProvider() MemoryFileProvider {
	return p.QwenMemoryFileProvider
}

// GetSlashCommandProvider returns the Qwen slash command provider
func (p *QwenToolProvider) GetSlashCommandProvider() SlashCommandProvider {
	return p.QwenSlashCommandProvider
}

// CostrictToolProvider combines Costrict memory file and slash command
// providers
type CostrictToolProvider struct {
	*CostrictMemoryFileProvider
	*CostrictSlashCommandProvider
	*SpectrAgentsUpdater
}

// NewCostrictToolProvider creates a new Costrict tool provider
func NewCostrictToolProvider() *CostrictToolProvider {
	return &CostrictToolProvider{
		CostrictMemoryFileProvider:   &CostrictMemoryFileProvider{},
		CostrictSlashCommandProvider: NewCostrictSlashCommandProvider(),
		SpectrAgentsUpdater:          &SpectrAgentsUpdater{},
	}
}

// GetName returns the tool name
func (*CostrictToolProvider) GetName() string {
	return "CoStrict"
}

// GetMemoryFileProvider returns the Costrict memory file provider
func (p *CostrictToolProvider) GetMemoryFileProvider() MemoryFileProvider {
	return p.CostrictMemoryFileProvider
}

// GetSlashCommandProvider returns the Costrict slash command provider
func (p *CostrictToolProvider) GetSlashCommandProvider() SlashCommandProvider {
	return p.CostrictSlashCommandProvider
}

// AntigravityToolProvider combines Antigravity memory file and slash
// command providers
type AntigravityToolProvider struct {
	*AgentsFileProvider
	*AntigravitySlashCommandProvider
	*SpectrAgentsUpdater
}

// NewAntigravityToolProvider creates a new Antigravity tool provider
func NewAntigravityToolProvider() *AntigravityToolProvider {
	return &AntigravityToolProvider{
		AgentsFileProvider:              &AgentsFileProvider{},
		AntigravitySlashCommandProvider: NewAntigravitySlashCommandProvider(),
		SpectrAgentsUpdater:             &SpectrAgentsUpdater{},
	}
}

// GetName returns the tool name
func (*AntigravityToolProvider) GetName() string {
	return "Antigravity"
}

// GetMemoryFileProvider returns the Antigravity memory file provider
func (p *AntigravityToolProvider) GetMemoryFileProvider() MemoryFileProvider {
	return p.AgentsFileProvider
}

// GetSlashCommandProvider returns the Antigravity slash command provider
//
//nolint:revive // line-length-limit - method name length required by interface
func (p *AntigravityToolProvider) GetSlashCommandProvider() SlashCommandProvider {
	return p.AntigravitySlashCommandProvider
}
