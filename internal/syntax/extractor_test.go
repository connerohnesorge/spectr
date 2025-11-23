package syntax

import (
	"testing"
)

func TestExtractRequirements(t *testing.T) {
	input := `
### Requirement: Req1
Some description.
#### Scenario: Scen1
Scenario details.
#### Scenario: Scen2
More details.

### Requirement: Req2
Another req.
`
	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	reqs := ExtractRequirements(doc)
	if len(reqs) != 2 {
		t.Fatalf("expected 2 requirements, got %d", len(reqs))
	}

	if reqs[0].Name != "Req1" {
		t.Errorf("expected Req1, got %q", reqs[0].Name)
	}
	if len(reqs[0].Scenarios) != 2 {
		t.Errorf("expected 2 scenarios in Req1, got %d", len(reqs[0].Scenarios))
	}
	if reqs[0].Scenarios[0] != "Scen1" {
		t.Errorf("expected Scen1, got %q", reqs[0].Scenarios[0])
	}

	if reqs[1].Name != "Req2" {
		t.Errorf("expected Req2, got %q", reqs[1].Name)
	}
}

func TestExtractDelta(t *testing.T) {
	input := `
## ADDED Requirements
### Requirement: AddedReq
#### Scenario: AddedScen

## MODIFIED Requirements
### Requirement: ModReq

## REMOVED Requirements
### Requirement: RemReq

## RENAMED Requirements
- FROM: ` + "`### Requirement: OldName`" + `
- TO: ` + "`### Requirement: NewName`" + `
`
	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	plan := ExtractDelta(doc)

	if len(plan.Added) != 1 {
		t.Errorf("expected 1 added req, got %d", len(plan.Added))
	} else if plan.Added[0].Name != "AddedReq" {
		t.Errorf("expected AddedReq, got %q", plan.Added[0].Name)
	}

	if len(plan.Modified) != 1 {
		t.Errorf("expected 1 modified req, got %d", len(plan.Modified))
	}

	if len(plan.Removed) != 1 {
		t.Errorf("expected 1 removed req, got %d", len(plan.Removed))
	} else if plan.Removed[0] != "RemReq" {
		t.Errorf("expected RemReq, got %q", plan.Removed[0])
	}

	if len(plan.Renamed) != 1 {
		t.Errorf("expected 1 renamed req, got %d", len(plan.Renamed))
	} else {
		if plan.Renamed[0].From != "OldName" {
			t.Errorf("expected OldName, got %q", plan.Renamed[0].From)
		}
		if plan.Renamed[0].To != "NewName" {
			t.Errorf("expected NewName, got %q", plan.Renamed[0].To)
		}
	}
}
