// Package discovery provides test helpers for change and spec discovery.
package discovery

import (
	"os"
	"path/filepath"
	"testing"
)

const (
	testDirPerm  = 0o755
	testFilePerm = 0o644
)

// createChangeDir creates a change directory with a proposal.md file.
// It sets up the directory structure and writes the proposal content.
func createChangeDir(
	t *testing.T,
	changesDir, name, content string,
) {
	t.Helper()
	changeDir := filepath.Join(changesDir, name)
	if err := os.MkdirAll(changeDir, testDirPerm); err != nil {
		t.Fatal(err)
	}
	proposalPath := filepath.Join(
		changeDir,
		"proposal.md",
	)
	if err := os.WriteFile(
		proposalPath,
		[]byte(content),
		testFilePerm,
	); err != nil {
		t.Fatal(err)
	}
}

// verifyChangesExcluded checks that specified changes are not present.
// It verifies excluded items don't appear in the changes list.
func verifyChangesExcluded(
	t *testing.T,
	changes, excluded []string,
) {
	t.Helper()
	changeMap := make(map[string]bool)
	for _, c := range changes {
		changeMap[c] = true
	}
	for _, exc := range excluded {
		if changeMap[exc] {
			t.Errorf(
				"Change %s should be excluded but was found",
				exc,
			)
		}
	}
}

// verifyOrdering checks that items match the expected order.
// It compares actual vs expected sequences element by element.
func verifyOrdering(
	t *testing.T,
	actual, expected []string,
	itemType string,
) {
	t.Helper()
	for i, exp := range expected {
		if i >= len(actual) {
			t.Errorf(
				"Not enough %ss returned",
				itemType,
			)

			break
		}
		if actual[i] != exp {
			t.Errorf(
				"Expected %s[%d] to be %s, got %s",
				itemType,
				i,
				exp,
				actual[i],
			)
		}
	}
}
