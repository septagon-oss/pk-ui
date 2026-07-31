// Validates: REQ-011.
// Per: ADR-0031, ADR-0076.
// Discipline: C-14.

package icon

import (
	"strings"
	"testing"
)

type testProvider struct {
	name    string
	glyph   Glyph
	matches string
}

func (provider testProvider) Name() string {
	return provider.name
}

func (provider testProvider) Resolve(name string, _ string) (Glyph, bool) {
	if name != provider.matches {
		return Glyph{}, false
	}
	return provider.glyph, true
}

func TestBuiltinProviderReturnsEditableCurrentColorVectors(t *testing.T) {
	t.Parallel()

	glyph, found := BuiltinProvider().Resolve("search", "outline")
	if !found {
		t.Fatal("builtin provider did not resolve search")
	}
	if err := glyph.Validate(); err != nil {
		t.Fatalf("builtin search glyph is invalid: %v", err)
	}
	if !strings.Contains(glyph.Body, "<path") {
		t.Fatalf("search glyph does not contain editable path data: %s", glyph.Body)
	}

	markup := RenderSVG("search", "outline")
	for _, fragment := range []string{
		`<svg xmlns="http://www.w3.org/2000/svg"`,
		`viewBox="0 0 256 256"`,
		`fill="currentColor"`,
		`data-pk-icon="search"`,
		`<path`,
	} {
		if !strings.Contains(markup, fragment) {
			t.Errorf("rendered SVG is missing %q: %s", fragment, markup)
		}
	}
}

func TestExtensionPrependsVocabularyAndPreservesOSSFallback(t *testing.T) {
	extension := testProvider{
		name:    "client-icons",
		matches: "client-mark",
		glyph: Glyph{
			Name:    "client-mark",
			ViewBox: "0 0 24 24",
			Body:    `<path d="M2 2h20v20H2z" fill="currentColor"/>`,
		},
	}
	restore := InstallExtension(extension)
	t.Cleanup(restore)

	glyph, known := Resolve("client-mark", "outline")
	if !known || glyph.Name != "client-mark" {
		t.Fatalf("extension glyph = %#v, known = %t", glyph, known)
	}
	glyph, known = Resolve("search", "outline")
	if !known || glyph.Name != "magnifying-glass" {
		t.Fatalf("OSS fallback glyph = %#v, known = %t", glyph, known)
	}
}

func TestUnsafeExtensionGlyphIsRejectedAndVisiblyFallsBack(t *testing.T) {
	extension := testProvider{
		name:    "unsafe-client-icons",
		matches: "unsafe",
		glyph: Glyph{
			Name:    "unsafe",
			ViewBox: "0 0 24 24",
			Body:    `<script>alert("xss")</script>`,
		},
	}
	restore := InstallExtension(extension)
	t.Cleanup(restore)

	glyph, known := Resolve("unsafe", "outline")
	if known {
		t.Fatal("unsafe extension glyph was accepted")
	}
	if glyph.Name != "question" {
		t.Fatalf("unsafe glyph fallback = %q, want question", glyph.Name)
	}
}

func TestNamedExtensionsResolveClientThenProductThenOSS(t *testing.T) {
	product := testProvider{
		name:    "product-icons",
		matches: "shared-mark",
		glyph: Glyph{
			Name:    "product-mark",
			ViewBox: "0 0 24 24",
			Body:    `<path d="M1 1h22v22H1z" fill="currentColor"/>`,
		},
	}
	client := testProvider{
		name:    "client-icons",
		matches: "shared-mark",
		glyph: Glyph{
			Name:    "client-mark",
			ViewBox: "0 0 24 24",
			Body:    `<circle cx="12" cy="12" r="10" fill="currentColor"/>`,
		},
	}
	restoreProduct, err := RegisterExtension(
		"platformkit-pro-test",
		ExtensionPriorityProduct,
		product,
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(restoreProduct)
	restoreClient, err := RegisterExtension(
		"client/test",
		ExtensionPriorityClient,
		client,
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(restoreClient)

	glyph, known := Resolve("shared-mark", "outline")
	if !known || glyph.Name != "client-mark" {
		t.Fatalf("client override = %#v, known = %t", glyph, known)
	}
	restoreClient()
	glyph, known = Resolve("shared-mark", "outline")
	if !known || glyph.Name != "product-mark" {
		t.Fatalf("product fallback = %#v, known = %t", glyph, known)
	}
	glyph, known = Resolve("search", "outline")
	if !known || glyph.Name != "magnifying-glass" {
		t.Fatalf("OSS fallback = %#v, known = %t", glyph, known)
	}
}

func TestGlyphValidationRejectsActiveAndExternalSVGContent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		body string
	}{
		{name: "script", body: `<script>alert(1)</script>`},
		{name: "event", body: `<path d="M0 0" onload="alert(1)"/>`},
		{name: "external image", body: `<image href="https://example.test/x.svg"/>`},
		{name: "external paint", body: `<path d="M0 0" fill="url(https://example.test/x.svg)"/>`},
		{name: "foreign object", body: `<foreignObject><p>unsafe</p></foreignObject>`},
		{name: "text", body: `<g>unexpected text</g>`},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := (Glyph{
				Name:    "probe",
				ViewBox: "0 0 24 24",
				Body:    test.body,
			}).Validate()
			if err == nil {
				t.Fatalf("Validate() accepted %s", test.body)
			}
		})
	}
}
