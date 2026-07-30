// Validates: REQ-011.
// Per: ADR-0031, ADR-0076.
// Discipline: C-14.

package componentgen

// vocabulary_test.go proves the generator fails closed on utilities that name
// something tw does not have, or that compile to a class outside tw's
// enumerable universe. Without these guarantees a typo in an authored document
// would survive generation and only surface when someone built the generated
// tree — or worse, would render an element with no styling at all.

import (
	"strings"
	"testing"
)

func styleWithBase(utilities ...Utility) WebStyle {
	style := BadgeWebStyle()
	style.Base = utilities
	return style
}

func TestVocabularyRejectsUnusableUtilities(t *testing.T) {
	t.Parallel()

	for _, testCase := range []struct {
		name    string
		utility Utility
		wants   string
	}{
		{
			name:    "misspelled method",
			utility: Utility{Method: "Backgruond", Arg: "SurfaceTertiary"},
			wants:   "no method",
		},
		{
			name:    "misspelled constant",
			utility: Utility{Method: "Bg", Arg: "SurfaceTertiaryy"},
			wants:   "is not a constant",
		},
		{
			name: "constant of the wrong type",
			// S1 is a Spacing; Bg takes a Color. Both are real constants, so
			// only a type-aware check catches this.
			utility: Utility{Method: "Bg", Arg: "S1"},
			wants:   "is not a constant of type Color",
		},
		{
			name:    "argument passed to a zero-argument method",
			utility: Utility{Method: "Relative", Arg: "S1"},
			wants:   "takes no argument",
		},
		{
			name:    "missing required argument",
			utility: Utility{Method: "Bg"},
			wants:   "requires a Color argument",
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			_, err := Generate(BadgeSource(), styleWithBase(testCase.utility))
			if err == nil {
				t.Fatalf("expected rejection of %+v", testCase.utility)
			}
			if !strings.Contains(err.Error(), testCase.wants) {
				t.Fatalf("error %q does not explain the problem (wanted %q)", err, testCase.wants)
			}
		})
	}
}

// TestVocabularyRejectsEscapeHatchClasses proves the enumerable-universe check
// is real. WidthRaw is a genuine tw method taking a genuine string, so method
// and argument checks both pass — only compiling the class and resolving it
// catches that it would render unstyled. This is the exact failure that cost a
// round of rework when the skeleton primitives were built with w-[60%].
func TestVocabularyRejectsEscapeHatchClasses(t *testing.T) {
	t.Parallel()

	_, err := Generate(BadgeSource(), styleWithBase(Utility{Method: "WidthRaw", Arg: ""}))
	if err == nil {
		t.Fatal("expected rejection of an escape-hatch utility")
	}
}

// TestVocabularyRejectsBaseVariantCollision proves the generator refuses the
// mistake a Figma extractor is most likely to make: hoisting a property into
// the base that variants also set. Equal specificity means the winner would be
// decided by alphabetical sheet order — silently.
func TestVocabularyRejectsBaseVariantCollision(t *testing.T) {
	t.Parallel()

	style := BadgeWebStyle()
	// Every Badge variant sets a background; putting one in the base too is
	// the collision.
	style.Base = append(style.Base, Utility{Method: "Bg", Arg: "SurfaceSecondary"})

	_, err := Generate(BadgeSource(), style)
	if err == nil {
		t.Fatal("expected rejection of a base that shadows its variants")
	}
	if !strings.Contains(err.Error(), "hoist it") {
		t.Fatalf("error %q should tell the author how to fix it", err)
	}
}

// TestVocabularyCoversTheShippedBadge is the positive control: every utility
// the real Badge uses must be expressible and resolvable. A failure here means
// the vocabulary is too narrow to describe components that already exist.
func TestVocabularyCoversTheShippedBadge(t *testing.T) {
	t.Parallel()

	vocabulary, err := LoadVocabulary()
	if err != nil {
		t.Fatalf("load vocabulary: %v", err)
	}
	if err := vocabulary.ValidateStyle(BadgeWebStyle()); err != nil {
		t.Fatalf("shipped Badge is not expressible: %v", err)
	}

	// The compiled base must equal what the hand-written class list compiles
	// to, which is what makes the generated stylesheet identical.
	compiled, err := vocabulary.Compile(BadgeWebStyle().Base)
	if err != nil {
		t.Fatalf("compile base: %v", err)
	}
	want := "inline-flex items-center gap-1 rounded-full font-medium px-2.5 py-0.5 text-xs"
	if compiled != want {
		t.Fatalf("base compiles to %q, want %q", compiled, want)
	}
}
