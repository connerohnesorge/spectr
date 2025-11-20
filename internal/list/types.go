package list

import "github.com/connerohnesorge/spectr/internal/parsers"

// ChangeInfo represents information about a change
type ChangeInfo struct {
	ID         string             `json:"id"`
	Title      string             `json:"title"`
	DeltaCount int                `json:"deltaCount"`
	TaskStatus parsers.TaskStatus `json:"taskStatus"`
}

// SpecInfo represents information about a spec
type SpecInfo struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	RequirementCount int    `json:"requirementCount"`
}

// ItemType represents the type of an item (change or spec)
type ItemType int

const (
	// ItemTypeChange represents a change item
	ItemTypeChange ItemType = iota
	// ItemTypeSpec represents a spec item
	ItemTypeSpec
)

// String returns the string representation of an ItemType
func (it ItemType) String() string {
	switch it {
	case ItemTypeChange:
		return "change"
	case ItemTypeSpec:
		return "spec"
	default:
		return "unknown"
	}
}

// Item represents a unified wrapper around either a ChangeInfo or SpecInfo.
// Exactly one of Change or Spec will be non-nil, determined by Type.
type Item struct {
	// Type indicates whether this item is a change or spec
	Type ItemType `json:"type"`
	// Change contains the change information if Type == ItemTypeChange
	Change *ChangeInfo `json:"change,omitempty"`
	// Spec contains the spec information if Type == ItemTypeSpec
	Spec *SpecInfo `json:"spec,omitempty"`
}

// ID returns the identifier for this item (change ID or spec ID)
func (i *Item) ID() string {
	switch i.Type {
	case ItemTypeChange:
		if i.Change != nil {
			return i.Change.ID
		}
	case ItemTypeSpec:
		if i.Spec != nil {
			return i.Spec.ID
		}
	}

	return ""
}

// Title returns the title for this item
func (i *Item) Title() string {
	switch i.Type {
	case ItemTypeChange:
		if i.Change != nil {
			return i.Change.Title
		}
	case ItemTypeSpec:
		if i.Spec != nil {
			return i.Spec.Title
		}
	}

	return ""
}

// ItemList represents a collection of mixed changes and specs
type ItemList []Item

// NewChangeItem creates a new Item wrapping a ChangeInfo
func NewChangeItem(change ChangeInfo) Item {
	return Item{
		Type:   ItemTypeChange,
		Change: &change,
	}
}

// NewSpecItem creates a new Item wrapping a SpecInfo
func NewSpecItem(spec SpecInfo) Item {
	return Item{
		Type: ItemTypeSpec,
		Spec: &spec,
	}
}

// FilterByType returns a new ItemList containing only items of the specified
// type.
func (il ItemList) FilterByType(itemType ItemType) ItemList {
	var filtered ItemList
	for _, item := range il {
		if item.Type == itemType {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// Changes returns all ChangeInfo items from the list
func (il ItemList) Changes() []ChangeInfo {
	var changes []ChangeInfo
	for _, item := range il {
		if item.Type == ItemTypeChange && item.Change != nil {
			changes = append(changes, *item.Change)
		}
	}

	return changes
}

// Specs returns all SpecInfo items from the list
func (il ItemList) Specs() []SpecInfo {
	var specs []SpecInfo
	for _, item := range il {
		if item.Type == ItemTypeSpec && item.Spec != nil {
			specs = append(specs, *item.Spec)
		}
	}

	return specs
}
