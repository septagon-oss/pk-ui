// Validates: REQ-011.
// Per: ADR-0031, ADR-0076.
// Discipline: C-14.

package stylesheet_test

// computed_test.go pins the reductions that decide whether two ways of saying
// the same thing compare equal. Each case here was a false alarm found while
// comparing a real design contract against the component it describes.

import (
	"maps"
	"testing"

	"github.com/septagon-oss/pk-ui/stylesheet"
)

func TestNormalizeValueExpressesLengthsOneWay(t *testing.T) {
	t.Parallel()

	for _, testCase := range []struct{ in, want string }{
		// A design contract writes 20rem where the component writes 320px.
		{"20rem", "320px"},
		{"1.75rem", "28px"},
		{"12px", "12px"},
		// The unit on zero is optional, so 0 and 0px are one length.
		{"0px", "0"},
		{"0rem", "0"},
		{"0", "0"},
		// Case and spacing carry no meaning.
		{"  #2A9D5B ", "#2a9d5b"},
		{"0 2px 4px rgba(0, 0, 0, .15)", "0 2px 4px rgba(0, 0, 0, .15)"},
		// Percentages and keywords are lengths in no unit system and stay put.
		{"100%", "100%"},
		{"auto", "auto"},
	} {
		if got := stylesheet.NormalizeValue(testCase.in); got != testCase.want {
			t.Errorf("NormalizeValue(%q) = %q, want %q", testCase.in, got, testCase.want)
		}
	}
}

func TestResolveVarsFollowsTheDeclaredValueThenTheFallback(t *testing.T) {
	t.Parallel()

	vars := map[string]string{"--collect-badge-size": "28px", "--brand": "var(--green)", "--green": "#2a9d5b"}
	for _, testCase := range []struct{ in, want string }{
		// The component declares the variable inline; the browser paints 28px.
		{"var(--collect-badge-size, 26px)", "28px"},
		// Nothing declares it, so the fallback is what gets painted.
		{"var(--unset, 26px)", "26px"},
		// Neither: left as written so it reads as itself in a failure message.
		{"var(--unset)", "var(--unset)"},
		// A variable defined in terms of another resolves through.
		{"var(--brand)", "#2a9d5b"},
		// Only the reference is replaced, not the rest of the value.
		{"calc(var(--collect-badge-size, 26px) * 0.55)", "calc(28px * 0.55)"},
	} {
		if got := stylesheet.ResolveVars(testCase.in, vars); got != testCase.want {
			t.Errorf("ResolveVars(%q) = %q, want %q", testCase.in, got, testCase.want)
		}
	}
}

// TestResolveVarsTerminatesOnASelfReferentialVariable proves a cycle cannot
// hang the caller. Real stylesheets contain them by accident.
func TestResolveVarsTerminatesOnASelfReferentialVariable(t *testing.T) {
	t.Parallel()

	got := stylesheet.ResolveVars("var(--loop)", map[string]string{"--loop": "var(--loop)"})
	if got != "var(--loop)" {
		t.Errorf("a self-referential variable resolved to %q", got)
	}
}

func TestParseInlineReadsAStyleAttribute(t *testing.T) {
	t.Parallel()

	got := stylesheet.ParseInline("--collect-badge-size:28px;background:#2A9D5B;;width : 320px ")
	want := map[string]string{
		"--collect-badge-size": "28px",
		"background":           "#2A9D5B",
		"width":                "320px",
	}
	if !maps.Equal(got, want) {
		t.Errorf("ParseInline = %v, want %v", got, want)
	}
}

// TestNormalizeSplitsAColourShorthandIntoItsLonghand pins the one shorthand
// reduction. A design contract states background-color where a component
// writes background; with only a colour in it, those set the same thing.
func TestNormalizeSplitsAColourShorthandIntoItsLonghand(t *testing.T) {
	t.Parallel()

	got := stylesheet.Normalize("background", "#2A9D5B", nil)
	if want := map[string]string{"background-color": "#2a9d5b"}; !maps.Equal(got, want) {
		t.Errorf("Normalize = %v, want %v", got, want)
	}

	// A shorthand carrying more than a colour is left alone: deciding which
	// part is the colour would be guessing.
	got = stylesheet.Normalize("background", "#2A9D5B url(x.png) no-repeat", nil)
	if _, split := got["background-color"]; split {
		t.Errorf("a multi-part shorthand was split into longhands: %v", got)
	}
}

// TestNormalizeDropsCustomPropertyDeclarations pins that a variable is not
// itself a visual fact. It paints nothing; its effect is counted where it is
// used.
func TestNormalizeDropsCustomPropertyDeclarations(t *testing.T) {
	t.Parallel()

	if got := stylesheet.Normalize("--collect-badge-size", "28px", nil); len(got) != 0 {
		t.Errorf("a custom property normalized to %v, want nothing", got)
	}
}

// TestNormalizeValueWritesNumbersOneWay pins the notation reductions found
// while comparing a real contract: one target writes "0.1em" where the other
// writes ".1em", and they are the same letter-spacing.
func TestNormalizeValueWritesNumbersOneWay(t *testing.T) {
	t.Parallel()

	for _, testCase := range []struct{ in, want string }{
		{".1em", "0.1em"},
		{"-.035em", "-0.035em"},
		{"1.020", "1.02"},
		// A zero in a dimension that is not a length keeps its unit: "0s" is a
		// duration and "0" is not.
		{"0s", "0s"},
		{"0deg", "0deg"},
		{"0%", "0"},
	} {
		if got := stylesheet.NormalizeValue(testCase.in); got != testCase.want {
			t.Errorf("NormalizeValue(%q) = %q, want %q", testCase.in, got, testCase.want)
		}
	}
}

// TestNormalizeTreatsAOneTokenBackgroundAsAColour pins the case Collect
// actually ships: the component writes "background: var(--sm-orange)" where
// the design contract states background-color. The variable cannot be
// evaluated, but a one-token background that is not an image is a colour, and
// leaving it under a different property name would report a difference of
// notation as a difference of paint.
func TestNormalizeTreatsAOneTokenBackgroundAsAColour(t *testing.T) {
	t.Parallel()

	if got := stylesheet.Normalize("background", "var(--sm-orange)", nil); len(got["background-color"]) == 0 {
		t.Errorf("a one-token background did not resolve to a colour: %v", got)
	}
	// An image is not a colour, whatever else it is.
	if got := stylesheet.Normalize("background", "url(cover.png)", nil); len(got["background-color"]) != 0 {
		t.Errorf("a background image was reported as a colour: %v", got)
	}
}
