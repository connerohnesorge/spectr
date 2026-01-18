package providers

import (
	"context"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/connerohnesorge/spectr/internal/domain"
	"github.com/spf13/afero"
)

const (
	testProposalContent = "proposal content"
	testApplyContent    = "apply content"
)

func TestSlashCommandsInitializer_Init(
	t *testing.T,
) {
	tests := []struct {
		name         string
		dir          string
		existingFile string
		wantCreated  []string
		wantUpdated  []string
	}{
		{
			name: "creates both slash command files",
			dir:  ".claude/commands/spectr",
			wantCreated: []string{
				".claude/commands/spectr/proposal.md",
				".claude/commands/spectr/apply.md",
			},
			wantUpdated: nil,
		},
		{
			name:         "overwrites existing file",
			dir:          ".claude/commands/spectr",
			existingFile: ".claude/commands/spectr/proposal.md",
			wantCreated: []string{
				".claude/commands/spectr/apply.md",
			},
			wantUpdated: []string{
				".claude/commands/spectr/proposal.md",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			projectFs := afero.NewMemMapFs()
			homeFs := afero.NewMemMapFs()
			cfg := &Config{SpectrDir: "spectr"}

			// Create directory
			_ = projectFs.MkdirAll(tt.dir, 0o755)

			// Create existing file if specified
			if tt.existingFile != "" {
				_ = afero.WriteFile(
					projectFs,
					tt.existingFile,
					[]byte("old content"),
					0o644,
				)
			}

			// Create templates
			commands := map[domain.SlashCommand]domain.TemplateRef{
				domain.SlashProposal: createTestTemplate(
					t,
					"proposal content",
				),
				domain.SlashApply: createTestTemplate(
					t,
					"apply content",
				),
			}

			// Create initializer
			init := NewSlashCommandsInitializer(
				tt.dir,
				commands,
			)

			// Execute
			result, err := init.Init(
				context.Background(),
				projectFs,
				homeFs,
				cfg,
				nil,
			)
			if err != nil {
				t.Fatalf("Init() failed: %v", err)
			}

			// Check result
			if !stringSliceEqualUnordered(
				result.CreatedFiles,
				tt.wantCreated,
			) {
				t.Errorf(
					"CreatedFiles = %v, want %v",
					result.CreatedFiles,
					tt.wantCreated,
				)
			}
			if !stringSliceEqualUnordered(
				result.UpdatedFiles,
				tt.wantUpdated,
			) {
				t.Errorf(
					"UpdatedFiles = %v, want %v",
					result.UpdatedFiles,
					tt.wantUpdated,
				)
			}

			// Verify files exist and have correct content
			proposalContent, err := afero.ReadFile(
				projectFs,
				filepath.Join(
					tt.dir,
					"proposal.md",
				),
			)
			if err != nil {
				t.Fatalf(
					"failed to read proposal.md: %v",
					err,
				)
			}
			if string(
				proposalContent,
			) != testProposalContent {
				t.Errorf(
					"proposal.md content = %q, want %q",
					string(proposalContent),
					"proposal content",
				)
			}

			applyContent, err := afero.ReadFile(
				projectFs,
				filepath.Join(tt.dir, "apply.md"),
			)
			if err != nil {
				t.Fatalf(
					"failed to read apply.md: %v",
					err,
				)
			}
			if string(
				applyContent,
			) != testApplyContent {
				t.Errorf(
					"apply.md content = %q, want %q",
					string(applyContent),
					"apply content",
				)
			}
		})
	}
}

func TestSlashCommandsInitializer_IsSetup(
	t *testing.T,
) {
	tests := []struct {
		name          string
		dir           string
		existingFiles []string
		want          bool
	}{
		{
			name: "returns true when all files exist",
			dir:  ".claude/commands/spectr",
			existingFiles: []string{
				".claude/commands/spectr/proposal.md",
				".claude/commands/spectr/apply.md",
			},
			want: true,
		},
		{
			name:          "returns false when no files exist",
			dir:           ".claude/commands/spectr",
			existingFiles: nil,
			want:          false,
		},
		{
			name: "returns false when some files missing",
			dir:  ".claude/commands/spectr",
			existingFiles: []string{
				".claude/commands/spectr/proposal.md",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			projectFs := afero.NewMemMapFs()
			homeFs := afero.NewMemMapFs()
			cfg := &Config{SpectrDir: "spectr"}

			// Create directory and files
			_ = projectFs.MkdirAll(tt.dir, 0o755)
			for _, file := range tt.existingFiles {
				_ = afero.WriteFile(
					projectFs,
					file,
					[]byte("content"),
					0o644,
				)
			}

			// Create templates
			commands := map[domain.SlashCommand]domain.TemplateRef{
				domain.SlashProposal: createTestTemplate(
					t,
					"proposal",
				),
				domain.SlashApply: createTestTemplate(
					t,
					"apply",
				),
			}

			// Create initializer
			init := NewSlashCommandsInitializer(
				tt.dir,
				commands,
			)

			// Execute
			got := init.IsSetup(
				projectFs,
				homeFs,
				cfg,
			)

			// Check result
			if got != tt.want {
				t.Errorf(
					"IsSetup() = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestSlashCommandsInitializer_dedupeKey(
	t *testing.T,
) {
	commands := map[domain.SlashCommand]domain.TemplateRef{
		domain.SlashProposal: createTestTemplate(
			t,
			"proposal",
		),
	}

	tests := []struct {
		name string
		dir  string
		want string
	}{
		{
			name: "simple path",
			dir:  ".claude/commands/spectr",
			want: "SlashCommandsInitializer:.claude/commands/spectr",
		},
		{
			name: "path with trailing slash",
			dir:  ".claude/commands/spectr/",
			want: "SlashCommandsInitializer:.claude/commands/spectr",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			init := &SlashCommandsInitializer{
				dir:      tt.dir,
				commands: commands,
			}
			got := init.dedupeKey()
			if got != tt.want {
				t.Errorf(
					"dedupeKey() = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestHomeSlashCommandsInitializer_Init(
	t *testing.T,
) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	dir := ".config/mytool/commands"
	_ = homeFs.MkdirAll(dir, 0o755)

	// Create templates
	commands := map[domain.SlashCommand]domain.TemplateRef{
		domain.SlashProposal: createTestTemplate(
			t,
			testProposalContent,
		),
		domain.SlashApply: createTestTemplate(
			t,
			testApplyContent,
		),
	}

	// Create initializer
	init := NewHomeSlashCommandsInitializer(
		dir,
		commands,
	)

	// Execute
	result, err := init.Init(
		context.Background(),
		projectFs,
		homeFs,
		cfg,
		nil,
	)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Check result
	if len(result.CreatedFiles) != 2 {
		t.Errorf(
			"CreatedFiles = %v, want 2 files",
			result.CreatedFiles,
		)
	}

	// Verify files exist in home filesystem
	proposalPath := filepath.Join(
		dir,
		"proposal.md",
	)
	applyPath := filepath.Join(dir, "apply.md")

	exists, _ := afero.Exists(
		homeFs,
		proposalPath,
	)
	if !exists {
		t.Error(
			"proposal.md should exist in home filesystem",
		)
	}

	exists, _ = afero.Exists(homeFs, applyPath)
	if !exists {
		t.Error(
			"apply.md should exist in home filesystem",
		)
	}

	// Verify files DO NOT exist in project filesystem
	exists, _ = afero.Exists(
		projectFs,
		proposalPath,
	)
	if exists {
		t.Error(
			"proposal.md should NOT exist in project filesystem",
		)
	}
}

func TestHomeSlashCommandsInitializer_dedupeKey(
	t *testing.T,
) {
	commands := map[domain.SlashCommand]domain.TemplateRef{
		domain.SlashProposal: createTestTemplate(
			t,
			"proposal",
		),
	}

	init := &HomeSlashCommandsInitializer{
		dir:      ".config/mytool",
		commands: commands,
	}
	got := init.dedupeKey()
	want := "HomeSlashCommandsInitializer:.config/mytool"
	if got != want {
		t.Errorf(
			"dedupeKey() = %v, want %v",
			got,
			want,
		)
	}
}

func TestPrefixedSlashCommandsInitializer_Init(
	t *testing.T,
) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	dir := ".agent/workflows"
	prefix := testSpectrPrefix
	_ = projectFs.MkdirAll(dir, 0o755)

	// Create templates
	commands := map[domain.SlashCommand]domain.TemplateRef{
		domain.SlashProposal: createTestTemplate(
			t,
			testProposalContent,
		),
		domain.SlashApply: createTestTemplate(
			t,
			testApplyContent,
		),
	}

	// Create initializer
	init := NewPrefixedSlashCommandsInitializer(
		dir,
		prefix,
		commands,
	)

	// Execute
	result, err := init.Init(
		context.Background(),
		projectFs,
		homeFs,
		cfg,
		nil,
	)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Check result
	if len(result.CreatedFiles) != 2 {
		t.Errorf(
			"CreatedFiles = %v, want 2 files",
			result.CreatedFiles,
		)
	}

	// Verify files exist with correct prefix
	proposalPath := filepath.Join(
		dir,
		"spectr-proposal.md",
	)
	applyPath := filepath.Join(
		dir,
		"spectr-apply.md",
	)

	proposalContent, err := afero.ReadFile(
		projectFs,
		proposalPath,
	)
	if err != nil {
		t.Fatalf(
			"failed to read spectr-proposal.md: %v",
			err,
		)
	}
	if string(
		proposalContent,
	) != "proposal content" {
		t.Errorf(
			"spectr-proposal.md content = %q, want %q",
			string(proposalContent),
			"proposal content",
		)
	}

	applyContent, err := afero.ReadFile(
		projectFs,
		applyPath,
	)
	if err != nil {
		t.Fatalf(
			"failed to read spectr-apply.md: %v",
			err,
		)
	}
	if string(applyContent) != "apply content" {
		t.Errorf(
			"spectr-apply.md content = %q, want %q",
			string(applyContent),
			"apply content",
		)
	}
}

func TestPrefixedSlashCommandsInitializer_IsSetup(
	t *testing.T,
) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	dir := ".agent/workflows"
	prefix := testSpectrPrefix
	_ = projectFs.MkdirAll(dir, 0o755)

	// Create templates
	commands := map[domain.SlashCommand]domain.TemplateRef{
		domain.SlashProposal: createTestTemplate(
			t,
			"proposal",
		),
		domain.SlashApply: createTestTemplate(
			t,
			"apply",
		),
	}

	// Create initializer
	init := NewPrefixedSlashCommandsInitializer(
		dir,
		prefix,
		commands,
	)

	// Should return false when files don't exist
	if init.IsSetup(projectFs, homeFs, cfg) {
		t.Error(
			"IsSetup() = true, want false when files don't exist",
		)
	}

	// Create files
	_ = afero.WriteFile(
		projectFs,
		filepath.Join(dir, "spectr-proposal.md"),
		[]byte("content"),
		0o644,
	)
	_ = afero.WriteFile(
		projectFs,
		filepath.Join(dir, "spectr-apply.md"),
		[]byte("content"),
		0o644,
	)

	// Should return true when files exist
	if !init.IsSetup(projectFs, homeFs, cfg) {
		t.Error(
			"IsSetup() = false, want true when files exist",
		)
	}
}

func TestPrefixedSlashCommandsInitializer_dedupeKey(
	t *testing.T,
) {
	commands := map[domain.SlashCommand]domain.TemplateRef{
		domain.SlashProposal: createTestTemplate(
			t,
			"proposal",
		),
	}

	init := &PrefixedSlashCommandsInitializer{
		dir:      ".agent/workflows",
		prefix:   "spectr-",
		commands: commands,
	}
	got := init.dedupeKey()
	want := "PrefixedSlashCommandsInitializer:.agent/workflows:spectr-"
	if got != want {
		t.Errorf(
			"dedupeKey() = %v, want %v",
			got,
			want,
		)
	}
}

func TestHomePrefixedSlashCommandsInitializer_Init(
	t *testing.T,
) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	dir := ".codex/prompts"
	prefix := testSpectrPrefix
	_ = homeFs.MkdirAll(dir, 0o755)

	// Create templates
	commands := map[domain.SlashCommand]domain.TemplateRef{
		domain.SlashProposal: createTestTemplate(
			t,
			testProposalContent,
		),
		domain.SlashApply: createTestTemplate(
			t,
			testApplyContent,
		),
	}

	// Create initializer
	init := NewHomePrefixedSlashCommandsInitializer(
		dir,
		prefix,
		commands,
	)

	// Execute
	result, err := init.Init(
		context.Background(),
		projectFs,
		homeFs,
		cfg,
		nil,
	)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Check result
	if len(result.CreatedFiles) != 2 {
		t.Errorf(
			"CreatedFiles = %v, want 2 files",
			result.CreatedFiles,
		)
	}

	// Verify files exist in home filesystem with correct prefix
	proposalPath := filepath.Join(
		dir,
		"spectr-proposal.md",
	)
	applyPath := filepath.Join(
		dir,
		"spectr-apply.md",
	)

	proposalContent, err := afero.ReadFile(
		homeFs,
		proposalPath,
	)
	if err != nil {
		t.Fatalf(
			"failed to read spectr-proposal.md from home fs: %v",
			err,
		)
	}
	if string(
		proposalContent,
	) != "proposal content" {
		t.Errorf(
			"spectr-proposal.md content = %q, want %q",
			string(proposalContent),
			"proposal content",
		)
	}

	applyContent, err := afero.ReadFile(
		homeFs,
		applyPath,
	)
	if err != nil {
		t.Fatalf(
			"failed to read spectr-apply.md from home fs: %v",
			err,
		)
	}
	if string(applyContent) != "apply content" {
		t.Errorf(
			"spectr-apply.md content = %q, want %q",
			string(applyContent),
			"apply content",
		)
	}

	// Verify files DO NOT exist in project filesystem
	exists, _ := afero.Exists(
		projectFs,
		proposalPath,
	)
	if exists {
		t.Error(
			"spectr-proposal.md should NOT exist in project filesystem",
		)
	}
}

func TestHomePrefixedSlashCommandsInitializer_dedupeKey(
	t *testing.T,
) {
	commands := map[domain.SlashCommand]domain.TemplateRef{
		domain.SlashProposal: createTestTemplate(
			t,
			"proposal",
		),
	}

	init := &HomePrefixedSlashCommandsInitializer{
		dir:      ".codex/prompts",
		prefix:   "spectr-",
		commands: commands,
	}
	got := init.dedupeKey()
	want := "HomePrefixedSlashCommandsInitializer:.codex/prompts:spectr-"
	if got != want {
		t.Errorf(
			"dedupeKey() = %v, want %v",
			got,
			want,
		)
	}
}

func TestTOMLSlashCommandsInitializer_Init(
	t *testing.T,
) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	dir := ".gemini/commands/spectr"
	_ = projectFs.MkdirAll(dir, 0o755)

	// Create templates
	commands := map[domain.SlashCommand]domain.TemplateRef{
		domain.SlashProposal: createTestTemplate(
			t,
			"proposal toml content",
		),
		domain.SlashApply: createTestTemplate(
			t,
			"apply toml content",
		),
	}

	// Create initializer
	init := NewTOMLSlashCommandsInitializer(
		dir,
		commands,
	)

	// Execute
	result, err := init.Init(
		context.Background(),
		projectFs,
		homeFs,
		cfg,
		nil,
	)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Check result
	if len(result.CreatedFiles) != 2 {
		t.Errorf(
			"CreatedFiles = %v, want 2 files",
			result.CreatedFiles,
		)
	}

	// Verify files exist with .toml extension
	proposalPath := filepath.Join(
		dir,
		"proposal.toml",
	)
	applyPath := filepath.Join(dir, "apply.toml")

	proposalContent, err := afero.ReadFile(
		projectFs,
		proposalPath,
	)
	if err != nil {
		t.Fatalf(
			"failed to read proposal.toml: %v",
			err,
		)
	}
	if string(
		proposalContent,
	) != "proposal toml content" {
		t.Errorf(
			"proposal.toml content = %q, want %q",
			string(proposalContent),
			"proposal toml content",
		)
	}

	applyContent, err := afero.ReadFile(
		projectFs,
		applyPath,
	)
	if err != nil {
		t.Fatalf(
			"failed to read apply.toml: %v",
			err,
		)
	}
	if string(
		applyContent,
	) != "apply toml content" {
		t.Errorf(
			"apply.toml content = %q, want %q",
			string(applyContent),
			"apply toml content",
		)
	}
}

func TestTOMLSlashCommandsInitializer_IsSetup(
	t *testing.T,
) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "spectr"}

	dir := ".gemini/commands/spectr"
	_ = projectFs.MkdirAll(dir, 0o755)

	// Create templates
	commands := map[domain.SlashCommand]domain.TemplateRef{
		domain.SlashProposal: createTestTemplate(
			t,
			"proposal",
		),
		domain.SlashApply: createTestTemplate(
			t,
			"apply",
		),
	}

	// Create initializer
	init := NewTOMLSlashCommandsInitializer(
		dir,
		commands,
	)

	// Should return false when files don't exist
	if init.IsSetup(projectFs, homeFs, cfg) {
		t.Error(
			"IsSetup() = true, want false when files don't exist",
		)
	}

	// Create files with .toml extension
	_ = afero.WriteFile(
		projectFs,
		filepath.Join(dir, "proposal.toml"),
		[]byte("content"),
		0o644,
	)
	_ = afero.WriteFile(
		projectFs,
		filepath.Join(dir, "apply.toml"),
		[]byte("content"),
		0o644,
	)

	// Should return true when files exist
	if !init.IsSetup(projectFs, homeFs, cfg) {
		t.Error(
			"IsSetup() = false, want true when files exist",
		)
	}
}

func TestTOMLSlashCommandsInitializer_dedupeKey(
	t *testing.T,
) {
	commands := map[domain.SlashCommand]domain.TemplateRef{
		domain.SlashProposal: createTestTemplate(
			t,
			"proposal",
		),
	}

	init := &TOMLSlashCommandsInitializer{
		dir:      ".gemini/commands/spectr",
		commands: commands,
	}
	got := init.dedupeKey()
	want := "TOMLSlashCommandsInitializer:.gemini/commands/spectr"
	if got != want {
		t.Errorf(
			"dedupeKey() = %v, want %v",
			got,
			want,
		)
	}
}

func TestSlashCommands_TemplateContextUsage(
	t *testing.T,
) {
	// Setup
	projectFs := afero.NewMemMapFs()
	homeFs := afero.NewMemMapFs()
	cfg := &Config{SpectrDir: "myspectr"}

	dir := ".claude/commands/spectr"
	_ = projectFs.MkdirAll(dir, 0o755)

	// Create template that uses context variables
	tmplText := "Specs: {{.SpecsDir}}"
	tmpl, err := template.New("test").
		Parse(tmplText)
	if err != nil {
		t.Fatalf(
			"failed to create template: %v",
			err,
		)
	}

	commands := map[domain.SlashCommand]domain.TemplateRef{
		domain.SlashProposal: {
			Name:     "test",
			Template: tmpl,
		},
	}

	// Create initializer
	init := NewSlashCommandsInitializer(
		dir,
		commands,
	)

	// Execute
	_, err = init.Init(
		context.Background(),
		projectFs,
		homeFs,
		cfg,
		nil,
	)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Verify template was rendered with correct context
	content, err := afero.ReadFile(
		projectFs,
		filepath.Join(dir, "proposal.md"),
	)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	expected := "Specs: myspectr/specs"
	if string(content) != expected {
		t.Errorf(
			"file content = %q, want %q",
			string(content),
			expected,
		)
	}
}

func TestSlashCommands_DifferentTypesHaveDifferentDedupeKeys(
	t *testing.T,
) {
	// Verify that different initializer types with same directory have different dedupe keys
	commands := map[domain.SlashCommand]domain.TemplateRef{
		domain.SlashProposal: createTestTemplate(
			t,
			"content",
		),
	}

	dir := ".claude/commands/spectr"

	slashInit := &SlashCommandsInitializer{
		dir:      dir,
		commands: commands,
	}
	homeSlashInit := &HomeSlashCommandsInitializer{
		dir:      dir,
		commands: commands,
	}
	tomlInit := &TOMLSlashCommandsInitializer{
		dir:      dir,
		commands: commands,
	}
	prefixedInit := &PrefixedSlashCommandsInitializer{
		dir:      dir,
		prefix:   "spectr-",
		commands: commands,
	}
	homePrefixedInit := &HomePrefixedSlashCommandsInitializer{
		dir:      dir,
		prefix:   "spectr-",
		commands: commands,
	}

	keys := []string{
		slashInit.dedupeKey(),
		homeSlashInit.dedupeKey(),
		tomlInit.dedupeKey(),
		prefixedInit.dedupeKey(),
		homePrefixedInit.dedupeKey(),
	}

	// Verify all keys are unique
	seen := make(map[string]bool)
	for i, key := range keys {
		if seen[key] {
			t.Errorf(
				"duplicate dedupe key at index %d: %s",
				i,
				key,
			)
		}
		seen[key] = true
	}

	// Verify expected key formats
	expectedKeys := []string{
		"SlashCommandsInitializer:.claude/commands/spectr",
		"HomeSlashCommandsInitializer:.claude/commands/spectr",
		"TOMLSlashCommandsInitializer:.claude/commands/spectr",
		"PrefixedSlashCommandsInitializer:.claude/commands/spectr:spectr-",
		"HomePrefixedSlashCommandsInitializer:.claude/commands/spectr:spectr-",
	}

	for i, expected := range expectedKeys {
		if keys[i] != expected {
			t.Errorf(
				"keys[%d] = %q, want %q",
				i,
				keys[i],
				expected,
			)
		}
	}
}

// Helper function to compare string slices without order
func stringSliceEqualUnordered(
	a, b []string,
) bool {
	if len(a) != len(b) {
		return false
	}
	if len(a) == 0 && len(b) == 0 {
		return true
	}

	// Create frequency maps
	countA := make(map[string]int)
	countB := make(map[string]int)

	for _, s := range a {
		countA[s]++
	}
	for _, s := range b {
		countB[s]++
	}

	// Compare maps
	if len(countA) != len(countB) {
		return false
	}
	for k, v := range countA {
		if countB[k] != v {
			return false
		}
	}

	return true
}
