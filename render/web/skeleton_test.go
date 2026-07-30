package web

// skeleton_test.go pins the loading-state contract: DeferredSlot's HTMX
// defaults and their overridability, the aria seam (busy slot, hidden
// placeholders), and the twins' geometry — the same class lists as the
// components they stand in for, so the swap-in cannot shift layout.
//
// Implements: REQ-011.
// Per: ADR-0031.
// Discipline: C-14.

import (
	"strings"
	"testing"

	g "maragu.dev/gomponents"

	"github.com/septagon-oss/pk-ui/contracts"
	"github.com/septagon-oss/pk-ui/contracts/atoms"
	"github.com/septagon-oss/pk-ui/contracts/molecules"
)

func renderNode(t *testing.T, n g.Node) string {
	t.Helper()
	var b strings.Builder
	if err := n.Render(&b); err != nil {
		t.Fatalf("render: %v", err)
	}
	return b.String()
}

func TestDeferredSlotDefaultsToLoadAndOuterHTML(t *testing.T) {
	t.Parallel()

	html := renderNode(t, DeferredSlot(atoms.DeferredSlotProps{
		HTMXProps: contracts.HTMXProps{Get: "/fragments/activity"},
	}, Skeleton(atoms.SkeletonProps{Shape: "text", Lines: 3})))

	for _, want := range []string{
		`hx-get="/fragments/activity"`,
		`hx-trigger="load"`,
		`hx-swap="outerHTML"`,
		`aria-busy="true"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("deferred slot missing %s in %s", want, html)
		}
	}
}

func TestDeferredSlotHonorsExplicitTriggerAndSwap(t *testing.T) {
	t.Parallel()

	html := renderNode(t, DeferredSlot(atoms.DeferredSlotProps{
		HTMXProps: contracts.HTMXProps{
			Get: "/fragments/chart", Trigger: "revealed", Swap: "innerHTML",
		},
	}))

	if !strings.Contains(html, `hx-trigger="revealed"`) {
		t.Errorf("explicit trigger overridden: %s", html)
	}
	if !strings.Contains(html, `hx-swap="innerHTML"`) {
		t.Errorf("explicit swap overridden: %s", html)
	}
}

func TestSkeletonPlaceholdersAreHiddenFromAssistiveTech(t *testing.T) {
	t.Parallel()

	for name, n := range map[string]g.Node{
		"block":  Skeleton(atoms.SkeletonProps{}),
		"text":   Skeleton(atoms.SkeletonProps{Shape: "text", Lines: 3}),
		"circle": Skeleton(atoms.SkeletonProps{Shape: "circle"}),
		"table":  TableSkeleton(molecules.TableSkeletonProps{}),
		"card":   CardSkeleton(molecules.CardSkeletonProps{}),
	} {
		if html := renderNode(t, n); !strings.Contains(html, `aria-hidden="true"`) {
			t.Errorf("%s placeholder is exposed to assistive tech: %s", name, html)
		}
	}
}

func TestSkeletonTextRendersRequestedLinesWithShortLast(t *testing.T) {
	t.Parallel()

	html := renderNode(t, Skeleton(atoms.SkeletonProps{Shape: "text", Lines: 3}))

	if got := strings.Count(html, clSkeleton.Compile()); got != 3 {
		t.Errorf("line count = %d, want 3: %s", got, html)
	}
	if got := strings.Count(html, clSkeletonLineLast.Compile()); got != 1 {
		t.Errorf("short-last-line count = %d, want exactly 1: %s", got, html)
	}
}

// The twins must be built from the mirrored component's own class lists —
// that identity is what guarantees zero layout shift at swap-in.
func TestTableSkeletonMirrorsTableGeometry(t *testing.T) {
	t.Parallel()

	html := renderNode(t, TableSkeleton(molecules.TableSkeletonProps{Columns: 2, Rows: 2}))

	for name, cl := range map[string]string{
		"wrap": clTableWrap.Compile(), "table": clTable.Compile(),
		"head": clTableHead.Compile(), "th": clTableTh.Compile(),
		"td": clTableTd.Compile(), "row": clTableRow.Compile(),
	} {
		if !strings.Contains(html, cl) {
			t.Errorf("table skeleton missing Table's %s classes %q", name, cl)
		}
	}
	if got := strings.Count(html, `<th `); got != 2 {
		t.Errorf("header cell count = %d, want 2", got)
	}
	if got := strings.Count(html, `<tr`); got != 3 {
		t.Errorf("row count = %d, want 3 (1 head + 2 body)", got)
	}
}

func TestCardSkeletonMirrorsCardFrame(t *testing.T) {
	t.Parallel()

	html := renderNode(t, CardSkeleton(molecules.CardSkeletonProps{}))

	if !strings.Contains(html, clCard.Compile()) {
		t.Errorf("card skeleton missing Card's frame classes %q: %s", clCard.Compile(), html)
	}
}
