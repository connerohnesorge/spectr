package discovery

import (
	"testing"
)

func TestNormalizeItemPath(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedID   string
		expectedType string
	}{
		// Test case 1: spectr/changes/<id> pattern
		{
			name:         "spectr/changes prefix extracts change ID",
			input:        "spectr/changes/my-change",
			expectedID:   "my-change",
			expectedType: "change",
		},
		// Test case 2: spectr/changes/<id>/specs/foo/spec.md pattern
		{
			name:         "spectr/changes with nested spec path extracts change ID",
			input:        "spectr/changes/my-change/specs/foo/spec.md",
			expectedID:   "my-change",
			expectedType: "change",
		},
		// Test case 3: spectr/specs/<id> pattern
		{
			name:         "spectr/specs prefix extracts spec ID",
			input:        "spectr/specs/my-spec",
			expectedID:   "my-spec",
			expectedType: "spec",
		},
		// Test case 4: spectr/specs/<id>/spec.md pattern
		{
			name:         "spectr/specs with spec.md extracts spec ID",
			input:        "spectr/specs/my-spec/spec.md",
			expectedID:   "my-spec",
			expectedType: "spec",
		},
		// Test case 5: simple ID without path prefix
		{
			name:         "simple ID returns as-is with empty type",
			input:        "my-change",
			expectedID:   "my-change",
			expectedType: "",
		},
		// Test case 6: absolute path with spectr/changes
		{
			name:         "absolute path extracts change ID",
			input:        "/home/user/project/spectr/changes/my-change",
			expectedID:   "my-change",
			expectedType: "change",
		},
		// Test case 7: empty string
		{
			name:         "empty string returns empty ID and type",
			input:        "",
			expectedID:   "",
			expectedType: "",
		},
		// Test case 8: spectr/changes/archive should NOT extract archive as change ID
		{
			name:         "spectr/changes/archive does not extract archive as change ID",
			input:        "spectr/changes/archive",
			expectedID:   "spectr/changes/archive",
			expectedType: "",
		},
		// Additional edge cases
		{
			name:         "absolute path with spectr/specs extracts spec ID",
			input:        "/home/user/project/spectr/specs/my-spec",
			expectedID:   "my-spec",
			expectedType: "spec",
		},
		{
			name:         "absolute path with spectr/specs/spec.md extracts spec ID",
			input:        "/home/user/project/spectr/specs/my-spec/spec.md",
			expectedID:   "my-spec",
			expectedType: "spec",
		},
		{
			name:         "spectr/changes with trailing slash",
			input:        "spectr/changes/my-change/",
			expectedID:   "my-change",
			expectedType: "change",
		},
		{
			name:         "spectr/specs with trailing slash",
			input:        "spectr/specs/my-spec/",
			expectedID:   "my-spec",
			expectedType: "spec",
		},
		{
			name:         "deeply nested absolute path with change",
			input:        "/very/deep/path/to/spectr/changes/deep-change/specs/capability/spec.md",
			expectedID:   "deep-change",
			expectedType: "change",
		},
		{
			name:         "simple ID with dashes",
			input:        "add-path-normalization-validate",
			expectedID:   "add-path-normalization-validate",
			expectedType: "",
		},
		{
			name:         "unrecognized path pattern returns original input",
			input:        "/some/random/path/to/file.md",
			expectedID:   "/some/random/path/to/file.md",
			expectedType: "",
		},
		{
			name:         "spectr/changes/archive/archived-change returns original",
			input:        "spectr/changes/archive/archived-change",
			expectedID:   "spectr/changes/archive/archived-change",
			expectedType: "",
		},
		{
			name:         "spectr prefix without changes or specs",
			input:        "spectr/something-else/my-id",
			expectedID:   "spectr/something-else/my-id",
			expectedType: "",
		},
		{
			name:         "proposal.md file path extracts change ID",
			input:        "spectr/changes/my-change/proposal.md",
			expectedID:   "my-change",
			expectedType: "change",
		},
		{
			name:         "tasks.json file path extracts change ID",
			input:        "spectr/changes/my-change/tasks.json",
			expectedID:   "my-change",
			expectedType: "change",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, inferredType := NormalizeItemPath(
				tt.input,
			)

			if id != tt.expectedID {
				t.Errorf(
					"NormalizeItemPath(%q) ID = %q, want %q",
					tt.input,
					id,
					tt.expectedID,
				)
			}

			if inferredType != tt.expectedType {
				t.Errorf(
					"NormalizeItemPath(%q) Type = %q, want %q",
					tt.input,
					inferredType,
					tt.expectedType,
				)
			}
		})
	}
}

func TestExtractFirstPathComponent(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "path with multiple components",
			input:    "my-change/specs/foo/spec.md",
			expected: "my-change",
		},
		{
			name:     "single component",
			input:    "my-change",
			expected: "my-change",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "leading slash",
			input:    "/my-change/specs",
			expected: "my-change",
		},
		{
			name:     "trailing slash",
			input:    "my-change/",
			expected: "my-change",
		},
		{
			name:     "only slashes",
			input:    "///",
			expected: "",
		},
		{
			name:     "component with dashes",
			input:    "add-path-normalization/tasks.json",
			expected: "add-path-normalization",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractFirstPathComponent(
				tt.input,
			)

			if result != tt.expected {
				t.Errorf(
					"extractFirstPathComponent(%q) = %q, want %q",
					tt.input,
					result,
					tt.expected,
				)
			}
		})
	}
}
