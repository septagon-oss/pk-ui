// Validates: REQ-011.
// Per: ADR-0031, ADR-0076.
// Discipline: C-14.

package stylesheet_test

// resolve_test.go proves the index answers what a class declares, on the kind
// of stylesheet it exists for: hand-authored, with selector lists, specificity
// devices, media queries and context selectors mixed together.

import (
	"maps"
	"os"
	"path/filepath"
	"testing"

	"github.com/septagon-oss/pk-ui/stylesheet"
)

const sample = `
/* a section banner, with a { brace } inside the comment */
.card {
  padding: 1rem;
  border-radius: 8px;
}
.card, .tile {
  background: white;
}
.badge.badge {
  padding: 0;
}
.card .title {
  font-weight: 700;
}
.card:hover {
  background: grey;
}
@media (min-width: 40rem) {
  .card { padding: 2rem; }
}
.tile {
  padding: 1.5rem;
}
`

func newIndex(t *testing.T, css string) *stylesheet.Index {
	t.Helper()
	index, err := stylesheet.NewIndex([]byte(css))
	if err != nil {
		t.Fatalf("index sample: %v", err)
	}
	return index
}

// TestIndexResolvesWhatAClassDeclaresOnItsOwn pins the core question.
func TestIndexResolvesWhatAClassDeclaresOnItsOwn(t *testing.T) {
	t.Parallel()

	got := newIndex(t, sample).Declarations("card")
	want := map[string]string{
		"padding":       "1rem",
		"border-radius": "8px",
		"background":    "white",
	}
	if !maps.Equal(got, want) {
		t.Errorf("card resolves to %v, want %v", got, want)
	}
}

// TestSelectorListsStyleEverySubject proves ".card, .tile" reaches both. A
// stylesheet that groups shared declarations is normal authoring, and dropping
// grouped rules would silently under-report what a class declares.
func TestSelectorListsStyleEverySubject(t *testing.T) {
	t.Parallel()

	if got := newIndex(t, sample).Declarations("tile")["background"]; got != "white" {
		t.Errorf("tile takes no background from the grouped rule; got %q", got)
	}
}

// TestRepeatedClassIsASpecificityDeviceNotASecondSubject pins ".badge.badge".
// Collect uses exactly this to beat single-class utilities without resorting
// to an inline style; if the index read it as a compound of two subjects the
// declaration would vanish from the class it belongs to.
func TestRepeatedClassIsASpecificityDeviceNotASecondSubject(t *testing.T) {
	t.Parallel()

	if got := newIndex(t, sample).Declarations("badge")["padding"]; got != "0" {
		t.Errorf("badge loses the padding its doubled selector declares; got %q", got)
	}
}

// TestContextualAndConditionalRulesAreNotAttributedToTheClass is the boundary.
// ".card .title" styles a descendant, ".card:hover" a state, and the media
// query a viewport. None of them is what .card declares unconditionally, and
// folding them in would make the index confidently wrong.
func TestContextualAndConditionalRulesAreNotAttributedToTheClass(t *testing.T) {
	t.Parallel()

	index := newIndex(t, sample)
	if index.Defines("title") {
		t.Error("a descendant selector attributed styling to the descendant class")
	}
	if got := index.Declarations("card")["background"]; got == "grey" {
		t.Error("the :hover rule was folded into the unconditional declarations")
	}
	if got := index.Declarations("card")["padding"]; got != "1rem" {
		t.Errorf("the media query overrode the base padding; got %q", got)
	}
}

// TestComputeFollowsTheCascadeNotTheAttributeOrder pins the merge rule. Both
// classes set padding at equal specificity, so the later rule in the
// stylesheet wins however the class attribute happens to be written.
func TestComputeFollowsTheCascadeNotTheAttributeOrder(t *testing.T) {
	t.Parallel()

	index := newIndex(t, sample)
	for _, order := range [][]string{{"card", "tile"}, {"tile", "card"}} {
		if got := index.Compute(order...)["padding"]; got != "1.5rem" {
			t.Errorf("Compute(%v) resolved padding to %q, want the later rule's 1.5rem", order, got)
		}
	}
}

// TestIndexReadsARealProductStylesheet keeps the parser honest against the
// 13k-line hand-authored file it was written for, rather than only the sample.
func TestIndexReadsARealProductStylesheet(t *testing.T) {
	t.Parallel()

	css, err := os.ReadFile(collectStylesheet(t))
	if err != nil {
		t.Skipf("Collect stylesheet unavailable: %v", err)
	}
	index, err := stylesheet.NewIndex(css)
	if err != nil {
		t.Fatalf("index Collect stylesheet: %v", err)
	}
	classes := index.Classes()
	if len(classes) < 100 {
		t.Fatalf("only %d classes resolved from a 13k-line stylesheet; the parse is not working", len(classes))
	}

	var withDeclarations int
	for _, class := range classes {
		if len(index.Declarations(class)) > 0 {
			withDeclarations++
		}
	}
	if withDeclarations != len(classes) {
		t.Errorf("%d of %d indexed classes resolve to nothing", len(classes)-withDeclarations, len(classes))
	}
	t.Logf("resolved %d classes from the Collect stylesheet", len(classes))
}

// collectStylesheet locates the product stylesheet in the estate checkout.
func collectStylesheet(t *testing.T) string {
	t.Helper()
	return filepath.Join(
		"..", "..", "..", "..",
		"modules", "platformkit-business-modules",
		"collectibles_management", "browser", "css", "collect.css",
	)
}

// TestMissingContextNamesWhatTheMarkupLacks pins the diagnosis that a
// component rendered outside its page needs.
//
// A stylesheet that writes ".album .cover" and never ".cover" leaves markup
// carrying only .cover completely unstyled, and nothing about that markup says
// why. This names the ancestor it is waiting on.
func TestMissingContextNamesWhatTheMarkupLacks(t *testing.T) {
	t.Parallel()

	index := newIndex(t, `
.album .cover { border-radius: 12px; }
.album .cover img { width: 100%; }
.page .sidebar { width: 200px; }
.cover { display: block; }
`)

	missing := index.MissingContext("cover")
	if missing["album"] != 2 {
		t.Errorf("MissingContext reported .album enabling %d rules, want 2", missing["album"])
	}
	// A rule the document shares nothing with is not this document's missing
	// context; reporting it would bury the answer in every other page's rules.
	if _, reported := missing["page"]; reported {
		t.Error("an unrelated rule was reported as missing context")
	}
	// A rule that already matches is not missing anything.
	if _, reported := missing["cover"]; reported {
		t.Error("a satisfied rule was counted as missing its own subject")
	}
}

// TestSubjectStylesAttributeDeclarationsToTheElementTheyLandOn pins the
// difference between what a document paints and what each element paints. A
// design contract's nodes are elements, so it needs the second.
func TestSubjectStylesAttributeDeclarationsToTheElementTheyLandOn(t *testing.T) {
	t.Parallel()

	index := newIndex(t, `
.pack-3d .pack-front { width: 256px; }
.pack-face.pack-front { background: white; }
.pack-title { font-size: 34px; }
.book img { width: 100%; }
.other-page .pack-front { width: 999px; }
`)

	styles := index.SubjectStyles("pack-3d", "pack-face", "pack-front", "pack-title")

	find := func(requires ...string) map[string]string {
		for _, style := range styles {
			if maps.Equal(setFrom(style.Requires), setFrom(requires)) {
				return style.Declarations
			}
		}
		return nil
	}

	// The contextual rule lands on the node carrying pack-front, because the
	// document carries pack-3d.
	if got := find("pack-front"); got["width"] != "256px" {
		t.Errorf("pack-front resolved to %v; the contextual width did not land", got)
	}
	// A compound subject requires the whole compound: attributing white to
	// pack-face alone would paint the back face with the front's fill.
	if got := find("pack-face", "pack-front"); got["background"] != "white" {
		t.Errorf("the compound subject resolved to %v", got)
	}
	// A sole-class rule is the degenerate case of the same mechanism.
	if got := find("pack-title"); got["font-size"] != "34px" {
		t.Errorf("pack-title resolved to %v", got)
	}
	// A rule needing context the document lacks contributes nothing.
	for _, style := range styles {
		for property, value := range style.Declarations {
			if property == "width" && value == "999px" {
				t.Error("a rule whose context is absent was attributed anyway")
			}
		}
	}
	// ".book img" lands on an element no class identifies.
	for _, style := range styles {
		if maps.Equal(setFrom(style.Requires), setFrom([]string{"book"})) {
			t.Error("an element-subject rule was attributed to a class")
		}
	}
}

func setFrom(values []string) map[string]string {
	out := map[string]string{}
	for _, v := range values {
		out[v] = ""
	}
	return out
}
