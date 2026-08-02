// Validates: REQ-004, REQ-011.
// Per: ADR-0004, ADR-0031, ADR-0076.
// Discipline: C-14.

package web

import (
	"strings"
	"testing"

	"github.com/septagon-oss/styleengine"
)

func TestDefaultStylesheetOwnsTokensRolesAndRenderableCatalogRules(t *testing.T) {
	t.Parallel()

	css, err := DefaultStylesheet(styleengine.RenderOptions{Minify: true})
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{
		"--pk-color-surface-primary:",
		"--pk-color-accent-default:",
		"--pk-role-surface-primary:",
		".bg-surface-primary",
		".text-fg-primary",
		".rounded-lg",
		`[data-component=progress][data-progress-percent="0"] [data-progress-fill=true]{width:0%}`,
		`[data-component=progress][data-progress-percent="50"] [data-progress-fill=true]{width:50%}`,
		`[data-component=progress][data-progress-percent="100"] [data-progress-fill=true]{width:100%}`,
	} {
		if !strings.Contains(css, required) {
			t.Errorf("complete OSS stylesheet missing %q", required)
		}
	}
	if strings.Contains(css, "platformkit-dev") {
		t.Fatal("OSS stylesheet leaked a proprietary namespace")
	}
}
