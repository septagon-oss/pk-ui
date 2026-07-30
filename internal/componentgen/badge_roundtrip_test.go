// Validates: REQ-011.
// Per: ADR-0031, ADR-0076.
// Discipline: C-14.

package componentgen

// badge_roundtrip_test.go is the round-trip proof. It generates Badge from its
// authored document and compares the result against the Badge that actually
// ships, read from the repository rather than from a golden file written by
// the generator's own author.
//
// Comparison is on normalized Go tokens: comments are stripped and whitespace
// collapsed, because a fragment formatted standalone cannot match the
// alignment gofmt produces inside a shared var block. What must match is the
// code, not its indentation.

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func repoFile(t *testing.T, relative string) string {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("..", "..", relative))
	if err != nil {
		t.Fatalf("read %s: %v", relative, err)
	}
	return string(raw)
}

var (
	lineComment  = regexp.MustCompile(`//[^\n]*`)
	whitespaceRE = regexp.MustCompile(`\s+`)
	// Go source wraps freely around punctuation, so a chain broken across
	// lines must compare equal to the same chain written inline.
	punctuationSpace = regexp.MustCompile(`\s*([.,(){}\[\]:=])\s*`)
)

// normalize strips comments and makes formatting irrelevant: whitespace
// collapses to a single space, and whitespace adjacent to punctuation is
// removed entirely. Word boundaries survive, so two different expressions
// cannot normalize to the same string.
func normalize(source string) string {
	source = lineComment.ReplaceAllString(source, "")
	source = whitespaceRE.ReplaceAllString(source, " ")
	source = punctuationSpace.ReplaceAllString(source, "$1")
	return strings.TrimSpace(source)
}

// sliceBetween returns the source from the first occurrence of start up to and
// including the terminator line.
func sliceBetween(t *testing.T, source, start, terminator string) string {
	t.Helper()
	from := strings.Index(source, start)
	if from < 0 {
		t.Fatalf("could not find %q", start)
	}
	rest := source[from:]
	to := strings.Index(rest, terminator)
	if to < 0 {
		t.Fatalf("could not find terminator %q after %q", terminator, start)
	}
	return rest[:to+len(terminator)]
}

func generateBadge(t *testing.T) Artifacts {
	t.Helper()
	artifacts, err := Generate(BadgeSource(), BadgeWebStyle())
	if err != nil {
		t.Fatalf("generate Badge: %v", err)
	}
	return artifacts
}

// TestGeneratedBadgeContractMatchesShipped proves the typed contract is
// derivable. The struct body must match field for field, tag for tag.
func TestGeneratedBadgeContractMatchesShipped(t *testing.T) {
	t.Parallel()

	generated := generateBadge(t).Contract
	shipped := repoFile(t, "contracts/atoms/badge.go")

	got := normalize(sliceBetween(t, generated, "type BadgeProps struct", "}"))
	want := normalize(sliceBetween(t, shipped, "type BadgeProps struct", "}"))
	if got != want {
		t.Fatalf("generated contract differs\n got: %s\nwant: %s", got, want)
	}

	// The ToMap method is part of the contract surface every Props type has.
	if !strings.Contains(normalize(generated), normalize("func (p BadgeProps) ToMap() map[string]any { return propsToMap(p) }")) {
		t.Error("generated contract is missing the ToMap method")
	}
}

// TestGeneratedBadgeClassListsMatchShipped proves the tw class lists are
// derivable, including the variant table and the dot part.
func TestGeneratedBadgeClassListsMatchShipped(t *testing.T) {
	t.Parallel()

	generated := generateBadge(t).ClassLists
	shipped := repoFile(t, "render/web/classlists.go")

	for _, testCase := range []struct{ name, start, terminator string }{
		{"base", "clBadgeBase = tw.New().", "FontSize(tw.TextXS)"},
		{"dot part", "clBadgeDot = tw.New()", "Bg(tw.FgBrand)"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			want := normalize(sliceBetween(t, shipped, testCase.start, testCase.terminator))
			if !strings.Contains(normalize(generated), want) {
				t.Fatalf("generated class lists do not contain the shipped %s\n want: %s\n  got: %s",
					testCase.name, want, normalize(generated))
			}
		})
	}

	// The variant table is a map literal, so entry order carries no meaning.
	// The generator emits sorted keys for deterministic output; the shipped
	// table is in authored order. Compare the entries as a set.
	t.Run("variant table", func(t *testing.T) {
		shippedTable := sliceBetween(t, shipped, "clBadgeVariant = map[string]tw.ClassList{", "\n\t}")
		entries := 0
		for _, line := range strings.Split(shippedTable, "\n") {
			line = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(line), ","))
			if !strings.HasPrefix(line, `"`) {
				continue
			}
			entries++
			if !strings.Contains(normalize(generated), normalize(line)) {
				t.Errorf("generated variant table is missing the shipped entry %s", line)
			}
		}
		if entries != 8 {
			t.Fatalf("expected 8 shipped variant entries, parsed %d — the proof is reading the wrong region", entries)
		}
	})
}

// TestGeneratedBadgeRendererMatchesShipped proves the renderer of a
// presentational component is derivable — the artifact that decides the
// rendered HTML, and therefore the golden gallery.
func TestGeneratedBadgeRendererMatchesShipped(t *testing.T) {
	t.Parallel()

	generated := generateBadge(t).Renderer
	shipped := repoFile(t, "render/web/atoms.go")

	got := normalize(sliceBetween(t, generated, "func Badge(", "\n}"))
	want := normalize(sliceBetween(t, shipped, "func Badge(", "\n}"))
	if got != want {
		t.Fatalf("generated renderer differs\n got: %s\nwant: %s", got, want)
	}
}

// TestGeneratedBadgeDeliveryMatchesShipped proves the catalog entry is
// derivable. Two differences are expected and asserted explicitly rather than
// normalized away, because each is a defect the single source fixes.
func TestGeneratedBadgeDeliveryMatchesShipped(t *testing.T) {
	t.Parallel()

	generated := normalize(generateBadge(t).DeliveryDefinition)
	shipped := normalize(sliceBetween(
		t,
		repoFile(t, "render/web/delivery_catalog.go"),
		"func badgeDeliveryDefinition()",
		"\n}",
	))

	// Structure the generator must reproduce exactly.
	for _, fragment := range []string{
		`componentType := "Badge"`,
		`deliveryFrame(componentType, clBadgeBase.Compile()`,
		`deliveryText("Label", "", "Active", "text")`,
		`"variant", compileClassMap(clBadgeVariant)`,
		`uicomponent.TierAtom`,
		`stableDeliveryID(componentType)`,
		`"Compact semantic status label."`,
		`deliveryDesign("default", "02 Atoms", "Status", root)`,
		`canonicalDeliveryExample(componentType, "default"`,
		`deliveryExample(componentType, "success"`,
	} {
		if !strings.Contains(generated, normalize(fragment)) {
			t.Errorf("generated delivery entry is missing %q", fragment)
		}
		if !strings.Contains(shipped, normalize(fragment)) {
			t.Errorf("shipped delivery entry no longer contains %q — the proof is comparing against something that moved", fragment)
		}
	}

	// Deviation: the shipped contract advertises seven variant values and
	// omits "danger", which its own class table styles. Generating both from
	// one Enum makes that disagreement unrepresentable.
	if strings.Contains(shipped, `"error", "info"`) && !strings.Contains(shipped, `"danger"`) {
		t.Log("confirmed drift in shipped Badge: delivery advertises 7 variants, class table styles 8 (danger missing)")
	}
	if !strings.Contains(generated, `"danger"`) {
		t.Error("generated delivery entry should advertise every styled variant, including danger")
	}
}

// TestGeneratedBadgeRejectsIncoherentSource proves the generator fails closed
// on documents that would produce an unnormalizable contract, rather than
// emitting code that only breaks later in pk-ui's own gates.
func TestGeneratedBadgeRejectsIncoherentSource(t *testing.T) {
	t.Parallel()

	t.Run("default outside enum", func(t *testing.T) {
		source := BadgeSource()
		source.Props[1].Default = "chartreuse"
		if _, err := Generate(source, BadgeWebStyle()); err == nil {
			t.Fatal("expected rejection of a default that is not one of the enum values")
		}
	})

	t.Run("no canonical example", func(t *testing.T) {
		source := BadgeSource()
		source.Examples[0].Canonical = false
		if _, err := Generate(source, BadgeWebStyle()); err == nil {
			t.Fatal("expected rejection of a document with no canonical example")
		}
	})

	t.Run("two canonical examples", func(t *testing.T) {
		source := BadgeSource()
		source.Examples[1].Canonical = true
		if _, err := Generate(source, BadgeWebStyle()); err == nil {
			t.Fatal("expected rejection of a document with two canonical examples")
		}
	})

	t.Run("mismatched style", func(t *testing.T) {
		style := BadgeWebStyle()
		style.ID = "Chip"
		if _, err := Generate(BadgeSource(), style); err == nil {
			t.Fatal("expected rejection when the web style belongs to another component")
		}
	})
}
