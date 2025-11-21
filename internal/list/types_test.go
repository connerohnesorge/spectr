package list

import (
	"encoding/json"
	"testing"

	"github.com/connerohnesorge/spectr/internal/parsers"
)

func TestItemType_String(t *testing.T) {
	tests := []struct {
		name     string
		itemType ItemType
		want     string
	}{
		{"Change type", ItemTypeChange, "change"},
		{"Spec type", ItemTypeSpec, "spec"},
		{"Unknown type", ItemType(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.itemType.String(); got != tt.want {
				t.Errorf("ItemType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewChangeItem(t *testing.T) {
	change := ChangeInfo{
		ID:         "test-change",
		Title:      "Test Change",
		DeltaCount: 3,
		TaskStatus: parsers.TaskStatus{Total: 5, Completed: 2},
	}

	item := NewChangeItem(change)

	if item.Type != ItemTypeChange {
		t.Errorf("NewChangeItem() Type = %v, want %v", item.Type, ItemTypeChange)
	}
	if item.Change == nil {
		t.Fatal("NewChangeItem() Change is nil")
	}
	if item.Spec != nil {
		t.Error("NewChangeItem() Spec should be nil")
	}
	if item.ID() != "test-change" {
		t.Errorf("Item.ID() = %v, want %v", item.ID(), "test-change")
	}
	if item.Title() != "Test Change" {
		t.Errorf("Item.Title() = %v, want %v", item.Title(), "Test Change")
	}
}

func TestNewSpecItem(t *testing.T) {
	spec := SpecInfo{
		ID:               "test-spec",
		Title:            "Test Spec",
		RequirementCount: 10,
	}

	item := NewSpecItem(spec)

	if item.Type != ItemTypeSpec {
		t.Errorf("NewSpecItem() Type = %v, want %v", item.Type, ItemTypeSpec)
	}
	if item.Spec == nil {
		t.Fatal("NewSpecItem() Spec is nil")
	}
	if item.Change != nil {
		t.Error("NewSpecItem() Change should be nil")
	}
	if item.ID() != "test-spec" {
		t.Errorf("Item.ID() = %v, want %v", item.ID(), "test-spec")
	}
	if item.Title() != "Test Spec" {
		t.Errorf("Item.Title() = %v, want %v", item.Title(), "Test Spec")
	}
}

func TestItem_ID_EmptyWhenNil(t *testing.T) {
	item := Item{Type: ItemTypeChange}
	if got := item.ID(); got != "" {
		t.Errorf("Item.ID() with nil Change = %v, want empty string", got)
	}

	item = Item{Type: ItemTypeSpec}
	if got := item.ID(); got != "" {
		t.Errorf("Item.ID() with nil Spec = %v, want empty string", got)
	}
}

func TestItem_Title_EmptyWhenNil(t *testing.T) {
	item := Item{Type: ItemTypeChange}
	if got := item.Title(); got != "" {
		t.Errorf("Item.Title() with nil Change = %v, want empty string", got)
	}

	item = Item{Type: ItemTypeSpec}
	if got := item.Title(); got != "" {
		t.Errorf("Item.Title() with nil Spec = %v, want empty string", got)
	}
}

func TestItemList_FilterByType(t *testing.T) {
	change1 := NewChangeItem(ChangeInfo{ID: "change1", Title: "Change 1"})
	change2 := NewChangeItem(ChangeInfo{ID: "change2", Title: "Change 2"})
	spec1 := NewSpecItem(SpecInfo{ID: "spec1", Title: "Spec 1"})
	spec2 := NewSpecItem(SpecInfo{ID: "spec2", Title: "Spec 2"})

	items := ItemList{change1, spec1, change2, spec2}

	changes := items.FilterByType(ItemTypeChange)
	if len(changes) != 2 {
		t.Errorf("FilterByType(ItemTypeChange) length = %v, want 2", len(changes))
	}

	specs := items.FilterByType(ItemTypeSpec)
	if len(specs) != 2 {
		t.Errorf("FilterByType(ItemTypeSpec) length = %v, want 2", len(specs))
	}
}

func TestItemList_Changes(t *testing.T) {
	change1 := NewChangeItem(ChangeInfo{ID: "change1", Title: "Change 1"})
	change2 := NewChangeItem(ChangeInfo{ID: "change2", Title: "Change 2"})
	spec1 := NewSpecItem(SpecInfo{ID: "spec1", Title: "Spec 1"})

	items := ItemList{change1, spec1, change2}

	changes := items.Changes()
	if len(changes) != 2 {
		t.Errorf("Changes() length = %v, want 2", len(changes))
	}

	// Verify content
	if changes[0].ID != "change1" {
		t.Errorf("Changes()[0].ID = %v, want change1", changes[0].ID)
	}
	if changes[1].ID != "change2" {
		t.Errorf("Changes()[1].ID = %v, want change2", changes[1].ID)
	}
}

func TestItemList_Specs(t *testing.T) {
	change1 := NewChangeItem(ChangeInfo{ID: "change1", Title: "Change 1"})
	spec1 := NewSpecItem(SpecInfo{ID: "spec1", Title: "Spec 1"})
	spec2 := NewSpecItem(SpecInfo{ID: "spec2", Title: "Spec 2"})

	items := ItemList{change1, spec1, spec2}

	specs := items.Specs()
	if len(specs) != 2 {
		t.Errorf("Specs() length = %v, want 2", len(specs))
	}

	// Verify content
	if specs[0].ID != "spec1" {
		t.Errorf("Specs()[0].ID = %v, want spec1", specs[0].ID)
	}
	if specs[1].ID != "spec2" {
		t.Errorf("Specs()[1].ID = %v, want spec2", specs[1].ID)
	}
}

func TestItem_JSONMarshaling(t *testing.T) {
	change := ChangeInfo{
		ID:         "test-change",
		Title:      "Test Change",
		DeltaCount: 3,
		TaskStatus: parsers.TaskStatus{Total: 5, Completed: 2},
	}
	item := NewChangeItem(change)

	data, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var decoded Item
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if decoded.Type != ItemTypeChange {
		t.Errorf("Decoded Type = %v, want %v", decoded.Type, ItemTypeChange)
	}
	if decoded.Change == nil {
		t.Fatal("Decoded Change is nil")
	}
	if decoded.Change.ID != "test-change" {
		t.Errorf("Decoded Change.ID = %v, want test-change", decoded.Change.ID)
	}
}

func TestItemList_EmptyLists(t *testing.T) {
	items := ItemList{}

	changes := items.Changes()
	if len(changes) != 0 {
		t.Errorf("Empty ItemList.Changes() length = %v, want 0", len(changes))
	}

	specs := items.Specs()
	if len(specs) != 0 {
		t.Errorf("Empty ItemList.Specs() length = %v, want 0", len(specs))
	}

	filtered := items.FilterByType(ItemTypeChange)
	if len(filtered) != 0 {
		t.Errorf("Empty ItemList.FilterByType() length = %v, want 0", len(filtered))
	}
}
