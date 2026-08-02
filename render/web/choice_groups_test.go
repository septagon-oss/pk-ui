// Validates: REQ-011.
// Per: ADR-0031, ADR-0076.
// Discipline: C-14.

package web

import (
	"bytes"
	"strings"
	"testing"

	"github.com/septagon-oss/pk-ui/contracts"
	"github.com/septagon-oss/pk-ui/contracts/atoms"
	"github.com/septagon-oss/pk-ui/contracts/molecules"
	g "maragu.dev/gomponents"
)

func TestCheckboxGroupUsesGroupRequiredSemantics(t *testing.T) {
	t.Parallel()

	html := renderChoiceGroup(t, CheckboxGroup(molecules.CheckboxGroupProps{
		ComponentProps: contracts.ComponentProps{ID: "topics"},
		HTMXProps:      contracts.HTMXProps{Post: "/preferences", Target: "#summary"},
		Name:           "topics",
		Label:          "Topics",
		Description:    "Choose at least one.",
		Required:       true,
		Selected:       []string{"security"},
		Options: []molecules.Option{
			{Label: "Product", Value: "product"},
			{Label: "Security", Value: "security", Description: "Critical notices."},
		},
	}))

	for _, want := range []string{
		`<fieldset`, `data-component="checkbox-group"`, `aria-required="true"`,
		`aria-describedby="topics-description"`, `hx-post="/preferences"`,
		`hx-trigger="change"`, `value="security" checked`, `Critical notices.`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("checkbox group missing %q:\n%s", want, html)
		}
	}
	if strings.Contains(html, `type="checkbox" required`) {
		t.Fatalf("group required must not require every checkbox:\n%s", html)
	}
}

func TestRadioGroupOwnsNativeSelectionValidationAndErrors(t *testing.T) {
	t.Parallel()

	html := renderChoiceGroup(t, RadioGroup(molecules.RadioGroupProps{
		ComponentProps: contracts.ComponentProps{Disabled: true},
		Name:           "plan",
		Label:          "Plan",
		Description:    "Choose one.",
		Value:          "growth",
		Required:       true,
		Orientation:    "horizontal",
		Error:          "A plan is required.",
		Options: []molecules.Option{
			{Label: "Starter", Value: "starter"},
			{Label: "Growth", Value: "growth"},
		},
	}))

	for _, want := range []string{
		`role="radiogroup"`, `data-orientation="horizontal"`, `disabled`,
		`aria-invalid="true"`,
		`aria-describedby="pk-radio-group-plan-description pk-radio-group-plan-error"`,
		`type="radio"`, `name="plan"`, `value="growth" checked required`,
		`role="alert"`, `A plan is required.`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("radio group missing %q:\n%s", want, html)
		}
	}
}

func TestRadioHelpTextIsAssociatedWithoutExplicitID(t *testing.T) {
	t.Parallel()

	html := renderChoiceGroup(t, Radio(atoms.RadioProps{
		Name: "plan", Label: "Growth", Value: "growth", HelpText: "For growing teams.",
	}))
	if !strings.Contains(html, `aria-describedby="pk-radio-`) ||
		!strings.Contains(html, `-help"`) ||
		!strings.Contains(html, `For growing teams.`) {
		t.Fatalf("radio help text is not associated:\n%s", html)
	}
}

func renderChoiceGroup(t *testing.T, node g.Node) string {
	t.Helper()
	var output bytes.Buffer
	if err := node.Render(&output); err != nil {
		t.Fatal(err)
	}
	return output.String()
}
