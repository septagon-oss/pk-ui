// Validates: REQ-011.
// Per: ADR-0031, ADR-0076.
// Discipline: C-14.

package componentgen

// style_roundtrip_test.go proves the return path end to end against Collect's
// real stylesheet: snapshot the shipped file, simulate the edit a designer
// would make in the design tool, diff it into a change set, and apply that
// change set back to the original file.
//
// The simulation is honest about what it stands in for. The design tool's
// export produces the "edited" snapshot; everything after that point is the
// code side, and everything after that point is what these tests exercise.

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func testProvenance(reason string) StyleProvenance {
	return StyleProvenance{
		Source:      "figma-plugin",
		GeneratedAt: time.Date(2026, 7, 31, 0, 0, 0, 0, time.UTC),
		Author:      "designer@septagon.dev",
		Reason:      reason,
	}
}

// TestSnapshotStylesheetIsStableAndContentAddressed proves the digest depends
// on content and not on when or by whom it was taken — the property that makes
// it usable as a concurrency guard.
func TestSnapshotStylesheetIsStableAndContentAddressed(t *testing.T) {
	t.Parallel()

	css := collectCSS(t)
	first, err := SnapshotStylesheet(css, "collect.css", testProvenance("export"))
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	second, err := SnapshotStylesheet(css, "collect.css", StyleProvenance{
		Source:      "someone-else",
		GeneratedAt: time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC),
		Author:      "other@example.test",
	})
	if err != nil {
		t.Fatalf("snapshot again: %v", err)
	}

	firstDigest, err := first.Digest()
	if err != nil {
		t.Fatalf("digest: %v", err)
	}
	secondDigest, err := second.Digest()
	if err != nil {
		t.Fatalf("digest: %v", err)
	}
	if firstDigest != secondDigest {
		t.Fatal("digest depends on provenance; it must depend only on content")
	}
	if len(first.Styles) < 1000 {
		t.Fatalf("snapshot carries only %d declarations from a 13k-line stylesheet", len(first.Styles))
	}
	t.Logf("snapshot carries %d unambiguous declarations, digest %s", len(first.Styles), firstDigest[:19])
}

// TestStyleRoundTripAppliesADesignerEdit is the end-to-end proof: a value
// changed as the design tool would change it survives snapshot, diff, and
// apply, and lands in the real file as a one-value diff.
func TestStyleRoundTripAppliesADesignerEdit(t *testing.T) {
	t.Parallel()

	original := collectCSS(t)
	target := editableTarget(t, original)

	base, err := SnapshotStylesheet(original, "collect.css", testProvenance("baseline"))
	if err != nil {
		t.Fatalf("base snapshot: %v", err)
	}

	// Stand-in for the design tool: the same file with one value changed.
	editedCSS, err := SetValue(original, target.Selector, target.Property, target.Value, "3px")
	if err != nil {
		t.Fatalf("simulate designer edit: %v", err)
	}
	edited, err := SnapshotStylesheet(editedCSS, "collect.css", testProvenance("designer edit"))
	if err != nil {
		t.Fatalf("edited snapshot: %v", err)
	}

	changes, err := DiffStyles(base, edited, testProvenance("tightened spacing"))
	if err != nil {
		t.Fatalf("diff: %v", err)
	}
	if len(changes.Changes) != 1 {
		t.Fatalf("diff carries %d changes, want exactly the one that was made", len(changes.Changes))
	}
	if changes.Changes[0].After.Value != "3px" {
		t.Fatalf("change carries %q, want %q", changes.Changes[0].After.Value, "3px")
	}

	// Apply to the ORIGINAL file: this is the receive path.
	applied, err := ApplyStyleChanges(original, changes)
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	if !bytes.Equal(applied, editedCSS) {
		t.Fatal("applying the change set did not reproduce the edited stylesheet byte for byte")
	}
	// And the authored file is otherwise untouched.
	if bytes.Count(applied, []byte("\n")) != bytes.Count(original, []byte("\n")) {
		t.Fatal("apply reflowed the stylesheet")
	}
}

// TestApplyStyleChangesRefusesStaleChangeSets proves the concurrency guard:
// a change set exported against one state cannot be applied to another. This
// is what stops two designers, or a designer and a developer, from silently
// overwriting each other.
func TestApplyStyleChangesRefusesStaleChangeSets(t *testing.T) {
	t.Parallel()

	original := collectCSS(t)
	target := editableTarget(t, original)

	base, err := SnapshotStylesheet(original, "collect.css", testProvenance("baseline"))
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	editedCSS, err := SetValue(original, target.Selector, target.Property, target.Value, "4px")
	if err != nil {
		t.Fatalf("edit: %v", err)
	}
	edited, err := SnapshotStylesheet(editedCSS, "collect.css", testProvenance("edit"))
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	changes, err := DiffStyles(base, edited, testProvenance("edit"))
	if err != nil {
		t.Fatalf("diff: %v", err)
	}

	// Someone else lands a different edit first.
	moved, err := SetValue(original, target.Selector, target.Property, target.Value, "9px")
	if err != nil {
		t.Fatalf("competing edit: %v", err)
	}
	_, err = ApplyStyleChanges(moved, changes)
	if err == nil {
		t.Fatal("expected refusal: the change set was observed against a state that has moved on")
	}
	if !strings.Contains(err.Error(), "re-export and retry") {
		t.Fatalf("refusal should tell the operator what to do, got %q", err)
	}
}

// TestDiffStylesRefusesStructuralChange proves the protocol's boundary. A
// design tool that invented a rule, or dropped one, is proposing a structural
// change; those are code decisions and are reported rather than applied.
func TestDiffStylesRefusesStructuralChange(t *testing.T) {
	t.Parallel()

	base := StyleSnapshot{
		SchemaVersion: StyleSnapshotSchema,
		Stylesheet:    "x.css",
		Styles:        []StyleValue{{Selector: ".a", Property: "color", Value: "red"}},
	}

	t.Run("new rule", func(t *testing.T) {
		edited := StyleSnapshot{
			SchemaVersion: StyleSnapshotSchema,
			Stylesheet:    "x.css",
			Styles: []StyleValue{
				{Selector: ".a", Property: "color", Value: "red"},
				{Selector: ".b", Property: "color", Value: "blue"},
			},
		}
		_, err := DiffStyles(base, edited, testProvenance("r"))
		if err == nil || !strings.Contains(err.Error(), "cannot create rules") {
			t.Fatalf("expected refusal to create a rule, got %v", err)
		}
	})

	t.Run("deleted rule", func(t *testing.T) {
		edited := StyleSnapshot{
			SchemaVersion: StyleSnapshotSchema,
			Stylesheet:    "x.css",
			Styles:        []StyleValue{},
		}
		_, err := DiffStyles(base, edited, testProvenance("r"))
		if err == nil || !strings.Contains(err.Error(), "cannot delete rules") {
			t.Fatalf("expected refusal to delete a rule, got %v", err)
		}
	})

	t.Run("no change", func(t *testing.T) {
		if _, err := DiffStyles(base, base, testProvenance("r")); err == nil {
			t.Fatal("expected refusal when nothing changed")
		}
	})
}

// TestStyleChangeSetValidationFailsClosed proves a malformed change set never
// reaches a file.
func TestStyleChangeSetValidationFailsClosed(t *testing.T) {
	t.Parallel()

	valid := StyleChangeSet{
		SchemaVersion: StyleChangeSetSchema,
		Stylesheet:    "collect.css",
		BaseDigest:    "sha256:" + strings.Repeat("a", 64),
		Provenance:    testProvenance("r"),
		Changes: []StyleChange{{
			Before: StyleValue{Selector: ".a", Property: "color", Value: "red"},
			After:  StyleValue{Selector: ".a", Property: "color", Value: "blue"},
		}},
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("a well-formed change set was rejected: %v", err)
	}

	for name, mutate := range map[string]func(*StyleChangeSet){
		"wrong schema":    func(c *StyleChangeSet) { c.SchemaVersion = "something.else" },
		"bad digest":      func(c *StyleChangeSet) { c.BaseDigest = "nope" },
		"no source":       func(c *StyleChangeSet) { c.Provenance.Source = "" },
		"no changes":      func(c *StyleChangeSet) { c.Changes = nil },
		"empty new value": func(c *StyleChangeSet) { c.Changes[0].After.Value = "  " },
		"moved coordinate": func(c *StyleChangeSet) {
			c.Changes[0].After.Selector = ".different"
		},
	} {
		t.Run(name, func(t *testing.T) {
			candidate := valid
			candidate.Changes = append([]StyleChange(nil), valid.Changes...)
			mutate(&candidate)
			if err := candidate.Validate(); err == nil {
				t.Fatalf("expected rejection of a change set with %s", name)
			}
		})
	}
}

// TestSnapshotOmitsAmbiguousCoordinates proves a snapshot never offers a
// declaration the editor would refuse — otherwise a designer could propose a
// change that can only ever fail.
func TestSnapshotOmitsAmbiguousCoordinates(t *testing.T) {
	t.Parallel()

	css := []byte(".pill { padding: 4px; }\n.pill { padding: 8px; }\n.solo { color: red; }\n")
	snapshot, err := SnapshotStylesheet(css, "x.css", testProvenance("r"))
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	for _, value := range snapshot.Styles {
		if value.Selector == ".pill" && value.Property == "padding" {
			t.Fatal("snapshot offered a coordinate the editor refuses as ambiguous")
		}
	}
	if len(snapshot.Styles) != 1 {
		t.Fatalf("snapshot carries %d declarations, want only the unambiguous one", len(snapshot.Styles))
	}
}
