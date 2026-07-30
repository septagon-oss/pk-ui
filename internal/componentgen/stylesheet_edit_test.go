// Validates: REQ-011.
// Per: ADR-0031, ADR-0076.
// Discipline: C-14.

package componentgen

// stylesheet_edit_test.go proves the write half of the style round trip
// against Collect's real stylesheet. A fixture would prove the parser; only a
// 13k-line hand-authored file proves the property that actually matters —
// that applying a designer's change leaves a diff a human can review.

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func collectCSS(t *testing.T) []byte {
	t.Helper()
	paths := collectStylesheets(t) // skips when the estate is absent
	raw, err := os.ReadFile(paths[0])
	if err != nil {
		t.Fatalf("read Collect stylesheet: %v", err)
	}
	return raw
}

// TestIndexDeclarationsReadsARealStylesheet proves the parser survives real
// CSS: nested media queries, comment banners, and multi-line selectors.
func TestIndexDeclarationsReadsARealStylesheet(t *testing.T) {
	t.Parallel()

	declarations, err := IndexDeclarations(collectCSS(t))
	if err != nil {
		t.Fatalf("index: %v", err)
	}
	if len(declarations) < 2000 {
		t.Fatalf("indexed only %d declarations; the parser is skipping rules", len(declarations))
	}

	// Every recorded range must actually bound its value, or a splice would
	// corrupt the file.
	css := collectCSS(t)
	for _, declaration := range declarations {
		if declaration.ValueStart < 0 || declaration.ValueEnd > len(css) || declaration.ValueStart >= declaration.ValueEnd {
			t.Fatalf("declaration %s has an invalid range", declaration.Coordinate())
		}
		if got := string(css[declaration.ValueStart:declaration.ValueEnd]); got != declaration.Value {
			t.Fatalf("range for %s spans %q, not its value %q", declaration.Coordinate(), got, declaration.Value)
		}
	}
	t.Logf("indexed %d declarations", len(declarations))
}

// editableTarget finds a declaration that is unique by (selector, property) at
// the top level, which is the only kind the writer will touch.
func editableTarget(t *testing.T, css []byte) Declaration {
	t.Helper()
	declarations, err := IndexDeclarations(css)
	if err != nil {
		t.Fatalf("index: %v", err)
	}
	counts := map[string]int{}
	for _, declaration := range declarations {
		if declaration.AtRule == "" {
			counts[declaration.Selector+"|"+declaration.Property]++
		}
	}
	for _, declaration := range declarations {
		if declaration.AtRule == "" && counts[declaration.Selector+"|"+declaration.Property] == 1 {
			return declaration
		}
	}
	t.Fatal("no unambiguous declaration found")
	return Declaration{}
}

// TestSetValueIsSurgical is the property the whole design rests on: the edit
// changes the intended value and nothing else, so the diff a reviewer sees is
// exactly what the designer changed.
func TestSetValueIsSurgical(t *testing.T) {
	t.Parallel()

	css := collectCSS(t)
	target := editableTarget(t, css)

	updated, err := SetValue(css, target.Selector, target.Property, target.Value, "1px")
	if err != nil {
		t.Fatalf("set value: %v", err)
	}

	// Byte accounting: only the value's bytes may differ in length.
	wantDelta := len("1px") - len(target.Value)
	if delta := len(updated) - len(css); delta != wantDelta {
		t.Fatalf("file length changed by %d, want %d — the edit touched more than the value", delta, wantDelta)
	}
	if !bytes.Equal(css[:target.ValueStart], updated[:target.ValueStart]) {
		t.Fatal("bytes before the edited value changed")
	}
	if !bytes.Equal(css[target.ValueEnd:], updated[target.ValueStart+len("1px"):]) {
		t.Fatal("bytes after the edited value changed")
	}

	// The change is observable through the parser, not just the bytes.
	after := editableTargetValue(t, updated, target.Selector, target.Property)
	if after != "1px" {
		t.Fatalf("re-indexed value is %q, want %q", after, "1px")
	}
}

func editableTargetValue(t *testing.T, css []byte, selector, property string) string {
	t.Helper()
	declarations, err := IndexDeclarations(css)
	if err != nil {
		t.Fatalf("re-index: %v", err)
	}
	for _, declaration := range declarations {
		if declaration.Selector == selector && declaration.Property == property && declaration.AtRule == "" {
			return declaration.Value
		}
	}
	t.Fatalf("declaration %s { %s } vanished after the edit", selector, property)
	return ""
}

// TestSetValueIdentityIsByteIdentical proves an unchanged value writes back the
// file exactly, so a no-op receive can never produce a spurious diff.
func TestSetValueIdentityIsByteIdentical(t *testing.T) {
	t.Parallel()

	css := collectCSS(t)
	target := editableTarget(t, css)

	updated, err := SetValue(css, target.Selector, target.Property, target.Value, target.Value)
	if err != nil {
		t.Fatalf("identity write: %v", err)
	}
	if !bytes.Equal(css, updated) {
		t.Fatal("writing a value back unchanged did not reproduce the file byte for byte")
	}
}

// TestSetValuePreservesAuthoredStructure proves the thing that makes this
// usable on a hand-authored file: comments and rule count survive.
func TestSetValuePreservesAuthoredStructure(t *testing.T) {
	t.Parallel()

	css := collectCSS(t)
	target := editableTarget(t, css)
	updated, err := SetValue(css, target.Selector, target.Property, target.Value, "2px")
	if err != nil {
		t.Fatalf("set value: %v", err)
	}

	for _, marker := range []string{"/*", "*/"} {
		if bytes.Count(css, []byte(marker)) != bytes.Count(updated, []byte(marker)) {
			t.Fatalf("comment delimiter %q count changed", marker)
		}
	}
	if bytes.Count(css, []byte("{")) != bytes.Count(updated, []byte("{")) {
		t.Fatal("rule count changed")
	}
	if bytes.Count(css, []byte("\n")) != bytes.Count(updated, []byte("\n")) {
		t.Fatal("line count changed — the file was reflowed")
	}
}

// TestSetValueFailsClosed proves every ambiguity is refused rather than
// guessed. Each of these would otherwise be a silent wrong write.
func TestSetValueFailsClosed(t *testing.T) {
	t.Parallel()

	css := collectCSS(t)
	target := editableTarget(t, css)

	t.Run("stale expected value", func(t *testing.T) {
		_, err := SetValue(css, target.Selector, target.Property, "definitely-not-the-current-value", "1px")
		if err == nil {
			t.Fatal("expected refusal when the observed value has moved on")
		}
		if !strings.Contains(err.Error(), "re-export and retry") {
			t.Errorf("error should tell the operator what to do, got %q", err)
		}
	})

	t.Run("unknown selector", func(t *testing.T) {
		if _, err := SetValue(css, ".pk-no-such-rule-xyz", "padding", "0", "1px"); err == nil {
			t.Fatal("expected refusal for a selector the stylesheet does not define")
		}
	})

	t.Run("unknown property on a real selector", func(t *testing.T) {
		if _, err := SetValue(css, target.Selector, "definitely-not-a-property", "0", "1px"); err == nil {
			t.Fatal("expected refusal for a property the rule does not declare")
		}
	})

	t.Run("empty replacement", func(t *testing.T) {
		if _, err := SetValue(css, target.Selector, target.Property, target.Value, ""); err == nil {
			t.Fatal("expected refusal to write an empty value")
		}
	})
}

// TestSetValueRefusesAmbiguousCoordinates proves the duplicate case is caught.
// A repeated declaration is legal CSS where the last wins, so editing the
// first would change nothing visible while reporting success — the worst
// possible outcome for a round trip.
func TestSetValueRefusesAmbiguousCoordinates(t *testing.T) {
	t.Parallel()

	css := []byte(".pill {\n  padding: 4px;\n}\n\n.pill {\n  padding: 8px;\n}\n")
	_, err := SetValue(css, ".pill", "padding", "4px", "12px")
	if err == nil {
		t.Fatal("expected refusal when a coordinate is declared twice")
	}
	if !strings.Contains(err.Error(), "ambiguous") {
		t.Fatalf("error should name the ambiguity, got %q", err)
	}
}

// TestIndexDeclarationsSeparatesMediaQueries proves a rule inside a media
// query is a distinct coordinate from the same selector at top level —
// otherwise a responsive override and its base would collide.
func TestIndexDeclarationsSeparatesMediaQueries(t *testing.T) {
	t.Parallel()

	css := []byte(".pill {\n  padding: 4px;\n}\n@media (min-width: 40rem) {\n  .pill {\n    padding: 8px;\n  }\n}\n")
	declarations, err := IndexDeclarations(css)
	if err != nil {
		t.Fatalf("index: %v", err)
	}
	var top, scoped int
	for _, declaration := range declarations {
		if declaration.Selector != ".pill" || declaration.Property != "padding" {
			continue
		}
		if declaration.AtRule == "" {
			top++
		} else {
			scoped++
		}
	}
	if top != 1 || scoped != 1 {
		t.Fatalf("expected one top-level and one media-scoped declaration, got %d and %d", top, scoped)
	}

	// The top-level edit must succeed and must not touch the media override.
	updated, err := SetValue(css, ".pill", "padding", "4px", "6px")
	if err != nil {
		t.Fatalf("top-level edit: %v", err)
	}
	if !bytes.Contains(updated, []byte("padding: 8px")) {
		t.Fatal("editing the base rule disturbed the media-query override")
	}
}
