package parsers

import (
	"encoding/json"
	"testing"
)

func TestTask_SerializationWithChildren(
	t *testing.T,
) {
	tests := []struct {
		name     string
		task     Task
		wantJSON string
	}{
		{
			name: "task without children field",
			task: Task{
				ID:          "1",
				Section:     "Implementation",
				Description: "Implement feature",
				Status:      TaskStatusPending,
			},
			wantJSON: `{"id":"1","section":"Implementation","description":"Implement feature","status":"pending"}`,
		},
		{
			name: "task with children field",
			task: Task{
				ID:          "1",
				Section:     "Implementation",
				Description: "Implement feature",
				Status:      TaskStatusPending,
				Children:    "$ref:tasks-1.jsonc",
			},
			wantJSON: `{"id":"1","section":"Implementation","description":"Implement feature","status":"pending","children":"$ref:tasks-1.jsonc"}`,
		},
		{
			name: "task with empty children field omits it",
			task: Task{
				ID:          "1",
				Section:     "Implementation",
				Description: "Implement feature",
				Status:      TaskStatusPending,
				Children:    "",
			},
			wantJSON: `{"id":"1","section":"Implementation","description":"Implement feature","status":"pending"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.task)
			if err != nil {
				t.Fatalf(
					"json.Marshal() error = %v",
					err,
				)
			}
			if string(got) != tt.wantJSON {
				t.Errorf(
					"json.Marshal() = %s, want %s",
					got,
					tt.wantJSON,
				)
			}
		})
	}
}

func TestTask_DeserializationWithChildren(
	t *testing.T,
) {
	tests := []struct {
		name     string
		jsonData string
		want     Task
	}{
		{
			name:     "task without children field",
			jsonData: `{"id":"1","section":"Implementation","description":"Implement feature","status":"pending"}`,
			want: Task{
				ID:          "1",
				Section:     "Implementation",
				Description: "Implement feature",
				Status:      TaskStatusPending,
			},
		},
		{
			name:     "task with children field",
			jsonData: `{"id":"1","section":"Implementation","description":"Implement feature","status":"pending","children":"$ref:tasks-1.jsonc"}`,
			want: Task{
				ID:          "1",
				Section:     "Implementation",
				Description: "Implement feature",
				Status:      TaskStatusPending,
				Children:    "$ref:tasks-1.jsonc",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Task
			err := json.Unmarshal(
				[]byte(tt.jsonData),
				&got,
			)
			if err != nil {
				t.Fatalf(
					"json.Unmarshal() error = %v",
					err,
				)
			}
			if got != tt.want {
				t.Errorf(
					"json.Unmarshal() = %+v, want %+v",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestTasksFile_SerializationWithParentAndIncludes(
	t *testing.T,
) {
	tests := []struct {
		name     string
		file     TasksFile
		wantJSON string
	}{
		{
			name: "version 1 flat file",
			file: TasksFile{
				Version: 1,
				Tasks: []Task{
					{
						ID:          "1",
						Section:     "Implementation",
						Description: "Implement feature",
						Status:      TaskStatusPending,
					},
				},
			},
			wantJSON: `{"version":1,"tasks":[{"id":"1","section":"Implementation","description":"Implement feature","status":"pending"}]}`,
		},
		{
			name: "version 2 root file with includes",
			file: TasksFile{
				Version: 2,
				Tasks: []Task{
					{
						ID:          "1",
						Section:     "Implementation",
						Description: "Implement feature",
						Status:      TaskStatusPending,
						Children:    "$ref:tasks-1.jsonc",
					},
				},
				Includes: []string{
					"tasks-*.jsonc",
				},
			},
			wantJSON: `{"version":2,"tasks":[{"id":"1","section":"Implementation","description":"Implement feature","status":"pending","children":"$ref:tasks-1.jsonc"}],"includes":["tasks-*.jsonc"]}`,
		},
		{
			name: "version 2 child file with parent",
			file: TasksFile{
				Version: 2,
				Tasks: []Task{
					{
						ID:          "1.1",
						Section:     "Implementation",
						Description: "Create database schema",
						Status:      TaskStatusPending,
					},
				},
				Parent: "1",
			},
			wantJSON: `{"version":2,"tasks":[{"id":"1.1","section":"Implementation","description":"Create database schema","status":"pending"}],"parent":"1"}`,
		},
		{
			name: "empty parent and includes omitted",
			file: TasksFile{
				Version: 2,
				Tasks: []Task{
					{
						ID:          "1",
						Section:     "Implementation",
						Description: "Implement feature",
						Status:      TaskStatusPending,
					},
				},
				Parent:   "",
				Includes: nil,
			},
			wantJSON: `{"version":2,"tasks":[{"id":"1","section":"Implementation","description":"Implement feature","status":"pending"}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.file)
			if err != nil {
				t.Fatalf(
					"json.Marshal() error = %v",
					err,
				)
			}
			if string(got) != tt.wantJSON {
				t.Errorf(
					"json.Marshal() = %s, want %s",
					got,
					tt.wantJSON,
				)
			}
		})
	}
}

func TestTasksFile_DeserializationWithParentAndIncludes(
	t *testing.T,
) {
	tests := []struct {
		name     string
		jsonData string
		want     TasksFile
	}{
		{
			name:     "version 1 flat file",
			jsonData: `{"version":1,"tasks":[{"id":"1","section":"Implementation","description":"Implement feature","status":"pending"}]}`,
			want: TasksFile{
				Version: 1,
				Tasks: []Task{
					{
						ID:          "1",
						Section:     "Implementation",
						Description: "Implement feature",
						Status:      TaskStatusPending,
					},
				},
			},
		},
		{
			name:     "version 2 root file with includes",
			jsonData: `{"version":2,"tasks":[{"id":"1","section":"Implementation","description":"Implement feature","status":"pending","children":"$ref:tasks-1.jsonc"}],"includes":["tasks-*.jsonc"]}`,
			want: TasksFile{
				Version: 2,
				Tasks: []Task{
					{
						ID:          "1",
						Section:     "Implementation",
						Description: "Implement feature",
						Status:      TaskStatusPending,
						Children:    "$ref:tasks-1.jsonc",
					},
				},
				Includes: []string{
					"tasks-*.jsonc",
				},
			},
		},
		{
			name:     "version 2 child file with parent",
			jsonData: `{"version":2,"tasks":[{"id":"1.1","section":"Implementation","description":"Create database schema","status":"pending"}],"parent":"1"}`,
			want: TasksFile{
				Version: 2,
				Tasks: []Task{
					{
						ID:          "1.1",
						Section:     "Implementation",
						Description: "Create database schema",
						Status:      TaskStatusPending,
					},
				},
				Parent: "1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got TasksFile
			err := json.Unmarshal(
				[]byte(tt.jsonData),
				&got,
			)
			if err != nil {
				t.Fatalf(
					"json.Unmarshal() error = %v",
					err,
				)
			}

			// Compare fields individually for better error messages
			if got.Version != tt.want.Version {
				t.Errorf(
					"Version = %d, want %d",
					got.Version,
					tt.want.Version,
				)
			}
			if got.Parent != tt.want.Parent {
				t.Errorf(
					"Parent = %s, want %s",
					got.Parent,
					tt.want.Parent,
				)
			}
			if len(
				got.Tasks,
			) != len(
				tt.want.Tasks,
			) {
				t.Errorf(
					"Tasks length = %d, want %d",
					len(got.Tasks),
					len(tt.want.Tasks),
				)
			} else {
				for i := range got.Tasks {
					if got.Tasks[i] != tt.want.Tasks[i] {
						t.Errorf("Tasks[%d] = %+v, want %+v", i, got.Tasks[i], tt.want.Tasks[i])
					}
				}
			}
			if len(
				got.Includes,
			) != len(
				tt.want.Includes,
			) {
				t.Errorf(
					"Includes length = %d, want %d",
					len(got.Includes),
					len(tt.want.Includes),
				)
			} else {
				for i := range got.Includes {
					if got.Includes[i] != tt.want.Includes[i] {
						t.Errorf("Includes[%d] = %s, want %s", i, got.Includes[i], tt.want.Includes[i])
					}
				}
			}
		})
	}
}
