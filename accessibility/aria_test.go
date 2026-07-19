package accessibility

// aria_test.go validates deterministic, injection-safe accessibility output.
//
// Validates: REQ-011.
// Per: ADR-0029.
// Discipline: C-14.

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestFocusTrapJavaScriptQuotesDynamicValues(t *testing.T) {
	t.Parallel()

	containerID := "dialog\");window.injected=true;//"
	selectors := []string{`button[data-label="close"]`, "input[name='search']"}
	trap := NewFocusTrap(containerID)
	trap.FocusableSelectors = selectors

	script := trap.GetJavaScript()
	encodedID, err := json.Marshal(containerID)
	if err != nil {
		t.Fatalf("json.Marshal(containerID) error = %v", err)
	}
	encodedSelectors, err := json.Marshal(strings.Join(selectors, ","))
	if err != nil {
		t.Fatalf("json.Marshal(selectors) error = %v", err)
	}
	wants := []string{
		"document.getElementById(" + string(encodedID) + ");",
		"container.querySelectorAll(" + string(encodedSelectors) + ");",
	}
	for _, want := range wants {
		if !strings.Contains(script, want) {
			t.Fatalf("GetJavaScript() missing %q:\n%s", want, script)
		}
	}
}
