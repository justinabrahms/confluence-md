package markdown

import (
	"strings"
	"testing"
)

func TestPreprocessConfluenceTasks(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantChecked   []string
		wantUnchecked []string
	}{
		{
			name: "single complete task",
			input: `<ac:task-list ac:task-list-id="123">
<ac:task>
<ac:task-id>1</ac:task-id>
<ac:task-uuid>abc-123</ac:task-uuid>
<ac:task-status>complete</ac:task-status>
<ac:task-body><span class="placeholder-inline-tasks">Do the thing</span></ac:task-body>
</ac:task>
</ac:task-list>`,
			wantChecked:   []string{"Do the thing"},
			wantUnchecked: nil,
		},
		{
			name: "single incomplete task",
			input: `<ac:task-list ac:task-list-id="123">
<ac:task>
<ac:task-id>1</ac:task-id>
<ac:task-uuid>abc-123</ac:task-uuid>
<ac:task-status>incomplete</ac:task-status>
<ac:task-body><span class="placeholder-inline-tasks">Do the thing</span></ac:task-body>
</ac:task>
</ac:task-list>`,
			wantChecked:   nil,
			wantUnchecked: []string{"Do the thing"},
		},
		{
			name: "mixed complete and incomplete tasks",
			input: `<ac:task-list ac:task-list-id="48f431f0-7f5e-4131-8de0-f0f1b5ceb499">
<ac:task>
<ac:task-id>1</ac:task-id>
<ac:task-uuid>66c1cef3-a623-4249-94f2-77be088c9b9a</ac:task-uuid>
<ac:task-status>complete</ac:task-status>
<ac:task-body><span class="placeholder-inline-tasks">moa api thing should have error rate by route <a href="https://onenr.io/0Zw0eEPanRv">https://onenr.io/0Zw0eEPanRv</a> not just response time by route https://root.ly/49niw4</span></ac:task-body>
</ac:task>
<ac:task>
<ac:task-id>3</ac:task-id>
<ac:task-uuid>efe285e4-7b5d-4ee8-b342-449b3c788643</ac:task-uuid>
<ac:task-status>incomplete</ac:task-status>
<ac:task-body><span class="placeholder-inline-tasks">Swap "main services view" dashboard to be prod by default. https://root.ly/49niw4</span></ac:task-body>
</ac:task>
</ac:task-list>`,
			wantChecked:   []string{"moa api thing"},
			wantUnchecked: []string{"Swap"},
		},
		{
			name: "task with inline comment markers",
			input: `<ac:task-list ac:task-list-id="123">
<ac:task>
<ac:task-id>3</ac:task-id>
<ac:task-uuid>xyz</ac:task-uuid>
<ac:task-status>incomplete</ac:task-status>
<ac:task-body><span class="placeholder-inline-tasks">Swap<ac:inline-comment-marker ac:ref="229c7393-af3d-42d7-a5e8-0e72e571e602"> "main services view" dashboard to be prod by defa</ac:inline-comment-marker>ult.</span></ac:task-body>
</ac:task>
</ac:task-list>`,
			wantChecked:   nil,
			wantUnchecked: []string{"Swap", "main services view", "dashboard"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := preprocessConfluenceTasks(tt.input)

			// Check that completed tasks have [x]
			for _, want := range tt.wantChecked {
				if !strings.Contains(result, "[x]") {
					t.Errorf("expected [x] checkbox in result for %q", want)
				}
				if !strings.Contains(result, want) {
					t.Errorf("expected %q in result, got: %s", want, result)
				}
			}

			// Check that incomplete tasks have [ ]
			for _, want := range tt.wantUnchecked {
				if !strings.Contains(result, "[ ]") {
					t.Errorf("expected [ ] checkbox in result for %q", want)
				}
				if !strings.Contains(result, want) {
					t.Errorf("expected %q in result, got: %s", want, result)
				}
			}

			// Check that task-list is converted to ul
			if strings.Contains(tt.input, "<ac:task-list") {
				if strings.Contains(result, "<ac:task-list") {
					t.Error("expected ac:task-list to be converted to ul")
				}
				if !strings.Contains(result, "<ul>") {
					t.Error("expected <ul> in result")
				}
			}

			// Check that ac:task elements are gone
			if strings.Contains(result, "<ac:task>") {
				t.Error("expected ac:task to be converted")
			}
			if strings.Contains(result, "<ac:task-status>") {
				t.Error("expected ac:task-status to be removed")
			}
			if strings.Contains(result, "<ac:task-uuid>") {
				t.Error("expected ac:task-uuid to be removed")
			}
		})
	}
}

func TestPreprocessConfluenceTasks_NoTasks(t *testing.T) {
	input := `<p>This is just regular HTML with no tasks</p>`
	result := preprocessConfluenceTasks(input)

	if result != input {
		t.Errorf("expected unchanged input when no tasks present, got: %s", result)
	}
}
