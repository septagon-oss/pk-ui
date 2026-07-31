// Validates: REQ-011.
// Per: ADR-0031, ADR-0076.
// Discipline: C-14.

package layout

// layout_test.go proves the primitive that makes page arrangement editable:
// an arrangement expressed as data renders exactly what the hand-written tree
// renders, so lifting a page into a document changes nothing a user sees.
//
// The fixture is the real shape of Collect's stickify page — a landing section
// wrapping a hero band, then benefits, how-it-works, faq, a tool output, and a
// closing section — because a toy tree would prove the walker and not the
// decision.

import (
	"slices"
	"strings"
	"testing"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// stickifySections stands in for the page's real sub-view functions. Their
// bodies are irrelevant to the arrangement; what matters is that arrangement
// can place them and cannot invent them.
func stickifySections() Sections {
	text := func(class, body string) func() g.Node {
		return func() g.Node { return h.Div(h.Class(class), g.Text(body)) }
	}
	return Sections{
		"heroCopy":   text("hero-copy", "Turn anything into stickers"),
		"workbench":  text("workbench", "Upload"),
		"benefits":   text("benefits", "Why"),
		"howItWorks": text("how-it-works", "How"),
		"faq":        text("faq", "Questions"),
		"toolOutput": text("tool-output", "Working…"),
		"closing":    text("closing", "Start"),
	}
}

// stickifyLayout is the page arrangement as data.
func stickifyLayout() Layout {
	return Layout{
		SchemaVersion: SchemaVersion,
		Page:          "collect/stickify",
		Root: Node{
			Element: "section",
			Classes: "collect-landing collect-stickify",
			Attrs:   map[string]string{"data-collect-stickify": ""},
			Children: []Node{
				{
					Element: "div",
					Classes: "stickify-hero-band",
					Children: []Node{{
						Element: "section",
						Classes: "stickify-hero",
						Children: []Node{
							{Section: "heroCopy"},
							{Section: "workbench"},
						},
					}},
				},
				{Section: "benefits"},
				{Section: "howItWorks"},
				{Section: "faq"},
				{Section: "toolOutput"},
				{Section: "closing"},
			},
		},
	}
}

// handWrittenStickify is the same page as Go, in the shape Collect writes it.
func handWrittenStickify(sections Sections) g.Node {
	return h.Section(
		h.Class("collect-landing collect-stickify"),
		g.Attr("data-collect-stickify", ""),
		h.Div(
			h.Class("stickify-hero-band"),
			h.Section(
				h.Class("stickify-hero"),
				sections["heroCopy"](),
				sections["workbench"](),
			),
		),
		sections["benefits"](),
		sections["howItWorks"](),
		sections["faq"](),
		sections["toolOutput"](),
		sections["closing"](),
	)
}

func render(t *testing.T, node g.Node) string {
	t.Helper()
	var b strings.Builder
	if err := node.Render(&b); err != nil {
		t.Fatalf("render: %v", err)
	}
	return b.String()
}

// TestArrangementRendersIdenticallyToHandWrittenGo is the decision's proof:
// lifting a page into a document is not a rewrite, because the output does not
// move. Without this, adopting layout documents would be a visual risk on
// every screen.
func TestArrangementRendersIdenticallyToHandWrittenGo(t *testing.T) {
	t.Parallel()

	sections := stickifySections()
	fromData, err := Render(stickifyLayout(), sections)
	if err != nil {
		t.Fatalf("render arrangement: %v", err)
	}
	if got, want := render(t, fromData), render(t, handWrittenStickify(sections)); got != want {
		t.Fatalf("arrangement renders differently from the hand-written page\n got: %s\nwant: %s", got, want)
	}
}

// TestReorderingSectionsIsAnArrangementChange proves the edit a designer
// actually makes: moving faq above benefits changes the page and nothing else.
func TestReorderingSectionsIsAnArrangementChange(t *testing.T) {
	t.Parallel()

	sections := stickifySections()
	base := stickifyLayout()

	reordered := stickifyLayout()
	children := reordered.Root.Children
	// benefits, howItWorks, faq -> faq, benefits, howItWorks
	children[1], children[2], children[3] = children[3], children[1], children[2]

	changes, err := Diff(base, reordered, Provenance{Source: "figma-plugin"})
	if err != nil {
		t.Fatalf("diff: %v", err)
	}
	applied, err := Apply(base, changes, sections)
	if err != nil {
		t.Fatalf("apply: %v", err)
	}

	out, err := Render(applied, sections)
	if err != nil {
		t.Fatalf("render applied: %v", err)
	}
	html := render(t, out)
	faq, benefits := strings.Index(html, "faq"), strings.Index(html, "benefits")
	if faq < 0 || benefits < 0 {
		t.Fatal("both sections should still be present")
	}
	if faq > benefits {
		t.Fatal("the reorder did not take effect")
	}
	// The section set is unchanged: a reorder moves content, never invents it.
	if got := len(SectionNames(applied)); got != len(SectionNames(base)) {
		t.Fatalf("section count changed from %d to %d", len(SectionNames(base)), got)
	}
}

// TestArrangementCannotInventBehaviour is the boundary. A design tool may
// arrange what exists; it may not name something the product cannot draw.
func TestArrangementCannotInventBehaviour(t *testing.T) {
	t.Parallel()

	sections := stickifySections()
	base := stickifyLayout()

	invented := stickifyLayout()
	invented.Root.Children = append(invented.Root.Children, Node{Section: "pricingTable"})

	if _, err := Diff(base, invented, Provenance{Source: "figma-plugin"}); err == nil {
		t.Fatal("expected refusal to introduce a section with no renderer")
	} else if !strings.Contains(err.Error(), "no renderer behind it") {
		t.Fatalf("refusal should name the reason, got %q", err)
	}

	// And it is refused again at render time, so the guard does not depend on
	// the diff having been run.
	if _, err := Render(invented, sections); err == nil {
		t.Fatal("expected render to refuse an unknown section")
	}
}

// TestApplyRefusesStaleChangeSets proves arrangements share the conflict
// semantics of tokens and styles.
func TestApplyRefusesStaleChangeSets(t *testing.T) {
	t.Parallel()

	sections := stickifySections()
	base := stickifyLayout()

	proposed := stickifyLayout()
	proposed.Root.Children[1], proposed.Root.Children[2] = proposed.Root.Children[2], proposed.Root.Children[1]
	changes, err := Diff(base, proposed, Provenance{Source: "figma-plugin"})
	if err != nil {
		t.Fatalf("diff: %v", err)
	}

	// Someone else lands a different arrangement first.
	moved := stickifyLayout()
	moved.Root.Classes = "collect-landing collect-stickify is-compact"

	if _, err := Apply(moved, changes, sections); err == nil {
		t.Fatal("expected refusal: the change set was observed against an arrangement that has moved on")
	} else if !strings.Contains(err.Error(), "re-export and retry") {
		t.Fatalf("refusal should tell the operator what to do, got %q", err)
	}
}

// TestSectionRemovalIsPermitted proves the asymmetry is deliberate: dropping a
// section is a layout decision, while adding one is a code decision.
func TestSectionRemovalIsPermitted(t *testing.T) {
	t.Parallel()

	sections := stickifySections()
	base := stickifyLayout()

	trimmed := stickifyLayout()
	trimmed.Root.Children = trimmed.Root.Children[:len(trimmed.Root.Children)-1] // drop closing

	changes, err := Diff(base, trimmed, Provenance{Source: "figma-plugin"})
	if err != nil {
		t.Fatalf("removing a section should be allowed: %v", err)
	}
	applied, err := Apply(base, changes, sections)
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	if slices.Contains(SectionNames(applied), "closing") {
		t.Fatal("the removed section is still arranged")
	}
}

// TestMalformedArrangementsAreRefused proves the walker never guesses.
func TestMalformedArrangementsAreRefused(t *testing.T) {
	t.Parallel()

	sections := stickifySections()
	for name, page := range map[string]Layout{
		"node is neither element nor section": {
			Page: "p", Root: Node{Children: []Node{{Classes: "x"}}},
		},
		"section also declares children": {
			Page: "p",
			Root: Node{Element: "div", Children: []Node{
				{Section: "faq", Children: []Node{{Element: "span"}}},
			}},
		},
		"unknown schema": {
			SchemaVersion: "pk.ui.page-layout.v99", Page: "p",
			Root: Node{Element: "div"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := Render(page, sections); err == nil {
				t.Fatalf("expected refusal of an arrangement where %s", name)
			}
		})
	}
}
