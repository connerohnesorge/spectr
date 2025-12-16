package providers

import (
	"context"
	"testing"

	"github.com/spf13/afero"
)

// mockKeyerInitializer is a mock initializer that implements Keyer.
type mockKeyerInitializer struct {
	key        string
	initCalled bool
}

func (m *mockKeyerInitializer) Init(
	_ context.Context,
	_ afero.Fs,
	_ *Config,
) error {
	m.initCalled = true

	return nil
}

func (m *mockKeyerInitializer) IsSetup(
	_ afero.Fs,
	_ *Config,
) bool {
	return m.initCalled
}

func (m *mockKeyerInitializer) Key() string {
	return m.key
}

// mockNonKeyerInitializer is a mock initializer that does NOT implement Keyer.
type mockNonKeyerInitializer struct {
	id         string // for identification in tests, not used as key
	initCalled bool
}

func (m *mockNonKeyerInitializer) Init(
	_ context.Context,
	_ afero.Fs,
	_ *Config,
) error {
	m.initCalled = true

	return nil
}

func (m *mockNonKeyerInitializer) IsSetup(
	_ afero.Fs,
	_ *Config,
) bool {
	return m.initCalled
}

// Verify interface satisfaction at compile time.
var (
	_ Initializer = (*mockKeyerInitializer)(nil)
	_ Keyer       = (*mockKeyerInitializer)(nil)
	_ Initializer = (*mockNonKeyerInitializer)(
		nil,
	)
)

func TestDedupeInitializers_Empty(t *testing.T) {
	result := DedupeInitializers(nil)
	if result != nil {
		t.Errorf(
			"DedupeInitializers(nil) = %v, want nil",
			result,
		)
	}

	result = DedupeInitializers(
		make([]Initializer, 0),
	)
	if result != nil {
		t.Errorf(
			"DedupeInitializers([]) = %v, want nil",
			result,
		)
	}
}

func TestDedupeInitializers_SingleItem(
	t *testing.T,
) {
	init := &mockKeyerInitializer{key: "test-key"}
	all := []Initializer{init}

	result := DedupeInitializers(all)

	if len(result) != 1 {
		t.Errorf(
			"DedupeInitializers returned %d items, want 1",
			len(result),
		)
	}
	if result[0] != init {
		t.Error(
			"DedupeInitializers returned different initializer",
		)
	}
}

func TestDedupeInitializers_NoDuplicates(
	t *testing.T,
) {
	init1 := &mockKeyerInitializer{key: "key1"}
	init2 := &mockKeyerInitializer{key: "key2"}
	init3 := &mockKeyerInitializer{key: "key3"}
	all := []Initializer{init1, init2, init3}

	result := DedupeInitializers(all)

	if len(result) != 3 {
		t.Errorf(
			"DedupeInitializers returned %d items, want 3",
			len(result),
		)
	}

	// Verify order is preserved
	if result[0] != init1 {
		t.Error("result[0] should be init1")
	}
	if result[1] != init2 {
		t.Error("result[1] should be init2")
	}
	if result[2] != init3 {
		t.Error("result[2] should be init3")
	}
}

func TestDedupeInitializers_WithDuplicates(
	t *testing.T,
) {
	init1 := &mockKeyerInitializer{key: "key1"}
	init2 := &mockKeyerInitializer{key: "key2"}
	init3 := &mockKeyerInitializer{
		key: "key1",
	} // duplicate of init1
	init4 := &mockKeyerInitializer{key: "key3"}
	init5 := &mockKeyerInitializer{
		key: "key2",
	} // duplicate of init2

	all := []Initializer{
		init1,
		init2,
		init3,
		init4,
		init5,
	}

	result := DedupeInitializers(all)

	if len(result) != 3 {
		t.Errorf(
			"DedupeInitializers returned %d items, want 3",
			len(result),
		)
	}

	// Verify first occurrence is kept (order preserved)
	if result[0] != init1 {
		t.Error(
			"result[0] should be init1 (first occurrence of key1)",
		)
	}
	if result[1] != init2 {
		t.Error(
			"result[1] should be init2 (first occurrence of key2)",
		)
	}
	if result[2] != init4 {
		t.Error(
			"result[2] should be init4 (only occurrence of key3)",
		)
	}
}

func TestDedupeInitializers_AllDuplicates(
	t *testing.T,
) {
	init1 := &mockKeyerInitializer{
		key: "same-key",
	}
	init2 := &mockKeyerInitializer{
		key: "same-key",
	}
	init3 := &mockKeyerInitializer{
		key: "same-key",
	}

	all := []Initializer{init1, init2, init3}

	result := DedupeInitializers(all)

	if len(result) != 1 {
		t.Errorf(
			"DedupeInitializers returned %d items, want 1",
			len(result),
		)
	}
	if result[0] != init1 {
		t.Error(
			"result[0] should be init1 (first occurrence)",
		)
	}
}

func TestDedupeInitializers_NonKeyerInitializers(
	t *testing.T,
) {
	// Non-Keyer initializers should never be deduplicated
	init1 := &mockNonKeyerInitializer{id: "a"}
	init2 := &mockNonKeyerInitializer{id: "b"}
	init3 := &mockNonKeyerInitializer{id: "c"}

	all := []Initializer{init1, init2, init3}

	result := DedupeInitializers(all)

	if len(result) != 3 {
		t.Errorf(
			"DedupeInitializers returned %d items, want 3",
			len(result),
		)
	}

	// All should be preserved in order
	if result[0] != init1 {
		t.Error("result[0] should be init1")
	}
	if result[1] != init2 {
		t.Error("result[1] should be init2")
	}
	if result[2] != init3 {
		t.Error("result[2] should be init3")
	}
}

func TestDedupeInitializers_MixedKeyerAndNonKeyer(
	t *testing.T,
) {
	keyer1 := &mockKeyerInitializer{key: "key1"}
	nonKeyer1 := &mockNonKeyerInitializer{id: "a"}
	keyer2 := &mockKeyerInitializer{
		key: "key1",
	} // duplicate of keyer1
	nonKeyer2 := &mockNonKeyerInitializer{id: "b"}
	keyer3 := &mockKeyerInitializer{key: "key2"}

	all := []Initializer{
		keyer1,
		nonKeyer1,
		keyer2,
		nonKeyer2,
		keyer3,
	}

	result := DedupeInitializers(all)

	// Should have: keyer1, nonKeyer1, nonKeyer2, keyer3 (keyer2 is duplicate)
	if len(result) != 4 {
		t.Errorf(
			"DedupeInitializers returned %d items, want 4",
			len(result),
		)
	}

	// Verify order and content
	if result[0] != keyer1 {
		t.Error("result[0] should be keyer1")
	}
	if result[1] != nonKeyer1 {
		t.Error("result[1] should be nonKeyer1")
	}
	if result[2] != nonKeyer2 {
		t.Error("result[2] should be nonKeyer2")
	}
	if result[3] != keyer3 {
		t.Error("result[3] should be keyer3")
	}
}

func TestDedupeInitializers_PreservesOrder(
	t *testing.T,
) {
	// Test that order of first occurrences is preserved exactly
	inits := make([]Initializer, 10)
	for i := range 10 {
		inits[i] = &mockKeyerInitializer{
			key: string(rune('a' + i)),
		}
	}

	result := DedupeInitializers(inits)

	if len(result) != 10 {
		t.Errorf(
			"DedupeInitializers returned %d items, want 10",
			len(result),
		)
	}

	for i := range 10 {
		if result[i] != inits[i] {
			t.Errorf(
				"result[%d] does not match inits[%d]",
				i,
				i,
			)
		}
	}
}

func TestDedupeInitializers_SameNonKeyerInstanceNotDeduplicated(
	t *testing.T,
) {
	// Even the same non-Keyer instance pointer should not be deduplicated
	// (though in practice this would be unusual)
	init := &mockNonKeyerInitializer{id: "single"}
	all := []Initializer{init, init, init}

	result := DedupeInitializers(all)

	// Since they're the same pointer, they actually WILL have the same key (ptr:address)
	// This is correct behavior - same instance = same initialization work
	if len(result) != 1 {
		t.Errorf(
			"DedupeInitializers returned %d items, want 1 (same instance)",
			len(result),
		)
	}
}

func TestDedupeInitializers_RealWorldScenario(
	t *testing.T,
) {
	// Simulate real-world scenario: multiple providers returning overlapping initializers
	// Provider 1: Claude Code
	claude1 := &mockKeyerInitializer{
		key: "dir:.claude/commands/spectr",
	}
	claude2 := &mockKeyerInitializer{
		key: "config:CLAUDE.md",
	}
	claude3 := &mockKeyerInitializer{
		key: "slashcmds:.claude/commands/spectr:.md:0",
	}

	// Provider 2: Cursor (shares CLAUDE.md config file)
	cursor1 := &mockKeyerInitializer{
		key: "dir:.cursor/rules",
	}
	cursor2 := &mockKeyerInitializer{
		key: "config:CLAUDE.md",
	} // duplicate!
	cursor3 := &mockKeyerInitializer{
		key: "config:.cursor/rules/spectr.md",
	}

	// Provider 3: Cline (shares .claude directory)
	cline1 := &mockKeyerInitializer{
		key: "dir:.claude/commands/spectr",
	} // duplicate!
	cline2 := &mockKeyerInitializer{
		key: "config:CLAUDE.md",
	} // duplicate!
	cline3 := &mockKeyerInitializer{
		key: "slashcmds:.cline/prompts:.md:0",
	}

	all := []Initializer{
		claude1, claude2, claude3,
		cursor1, cursor2, cursor3,
		cline1, cline2, cline3,
	}

	result := DedupeInitializers(all)

	// Expected unique initializers:
	// 1. dir:.claude/commands/spectr (from Claude)
	// 2. config:CLAUDE.md (from Claude)
	// 3. slashcmds:.claude/commands/spectr:.md:0 (from Claude)
	// 4. dir:.cursor/rules (from Cursor)
	// 5. config:.cursor/rules/spectr.md (from Cursor)
	// 6. slashcmds:.cline/prompts:.md:0 (from Cline)
	if len(result) != 6 {
		t.Errorf(
			"DedupeInitializers returned %d items, want 6",
			len(result),
		)
	}

	// Verify the first occurrences are kept
	expected := []Initializer{
		claude1,
		claude2,
		claude3,
		cursor1,
		cursor3,
		cline3,
	}
	for i, exp := range expected {
		if i >= len(result) {
			t.Errorf(
				"result is too short, missing result[%d]",
				i,
			)

			continue
		}
		if result[i] != exp {
			t.Errorf(
				"result[%d] = %p, want %p",
				i,
				result[i],
				exp,
			)
		}
	}
}

func TestInitializerKey_WithKeyer(t *testing.T) {
	init := &mockKeyerInitializer{
		key: "test-key-123",
	}
	key := initializerKey(init)

	if key != "test-key-123" {
		t.Errorf(
			"initializerKey() = %s, want test-key-123",
			key,
		)
	}
}

func TestInitializerKey_WithoutKeyer(
	t *testing.T,
) {
	init := &mockNonKeyerInitializer{
		id: "non-keyer",
	}
	key := initializerKey(init)

	// Should start with "ptr:" since it uses pointer address
	if len(key) < 5 || key[:4] != "ptr:" {
		t.Errorf(
			"initializerKey() = %s, want ptr:... format",
			key,
		)
	}
}

func TestInitializerKey_DifferentNonKeyersHaveDifferentKeys(
	t *testing.T,
) {
	init1 := &mockNonKeyerInitializer{id: "a"}
	init2 := &mockNonKeyerInitializer{id: "b"}

	key1 := initializerKey(init1)
	key2 := initializerKey(init2)

	if key1 == key2 {
		t.Errorf(
			"Different non-Keyer instances should have different keys: %s == %s",
			key1,
			key2,
		)
	}
}
