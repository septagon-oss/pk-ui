// Validates: REQ-011.
// Per: ADR-0031, ADR-0076.
// Discipline: C-14.

package componentgen

// oracle_test.go proves the stylesheet oracle against a real product
// stylesheet rather than a fixture, because the properties that matter — does
// it find the classes a product actually ships, does it reject a typo — are
// only meaningful at real scale. Collect's stylesheet is ~13k lines of
// hand-authored CSS carrying the class vocabulary of the estate's largest
// product surface.
//
// The test skips rather than fails when that stylesheet is absent, so pk-ui
// stays independently buildable outside the estate.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// collectStylesheets locates Collect's CSS within the estate checkout.
func collectStylesheets(t *testing.T) []string {
	t.Helper()
	root := filepath.Join(
		"..", "..", "..", "..", "..",
		"modules", "platformkit-business-modules", "collectibles_management",
		"browser", "css",
	)
	paths := []string{filepath.Join(root, "collect.css"), filepath.Join(root, "landing.css")}
	for _, path := range paths {
		if _, err := os.Stat(path); err != nil {
			t.Skipf("Collect stylesheet not present (%s); this test is estate-only", path)
		}
	}
	return paths
}

func collectOracle(t *testing.T) *StylesheetOracle {
	t.Helper()
	oracle, err := LoadStylesheetOracle("collect.css", collectStylesheets(t)...)
	if err != nil {
		t.Fatalf("load Collect stylesheet oracle: %v", err)
	}
	return oracle
}

// TestStylesheetOracleIndexesARealProductStylesheet proves the parser survives
// a genuine 13k-line stylesheet and finds a substantial vocabulary rather than
// a handful of selectors it happened to match.
func TestStylesheetOracleIndexesARealProductStylesheet(t *testing.T) {
	t.Parallel()

	oracle := collectOracle(t)
	classes := oracle.Classes()
	if len(classes) < 500 {
		t.Fatalf("indexed only %d classes from Collect's stylesheet; the parser is missing rules", len(classes))
	}
	t.Logf("indexed %d class selectors from Collect's stylesheet", len(classes))
}

// TestStylesheetOracleAcceptsClassesTheProductShips is the positive control:
// classes taken from Collect's own component renderers must validate. A
// failure means the oracle would reject the product's real markup.
func TestStylesheetOracleAcceptsClassesTheProductShips(t *testing.T) {
	t.Parallel()

	oracle := collectOracle(t)
	// Drawn from components/atoms/atoms.go and molecules/molecules.go.
	for _, class := range []string{
		"pill",
		"collect-icon",
		"num-tag",
		"collect-avatar",
		"sticker-art-frame",
	} {
		if err := oracle.Resolve(class); err != nil {
			t.Errorf("oracle rejects a class Collect ships: %v", err)
		}
	}
}

// TestStylesheetOracleRejectsUndefinedClasses proves the guarantee: a class the
// stylesheet does not define is refused, and the message names a near miss so
// the rejection is actionable rather than merely correct.
func TestStylesheetOracleRejectsUndefinedClasses(t *testing.T) {
	t.Parallel()

	oracle := collectOracle(t)

	err := oracle.Resolve("colect-icon") // transposed typo
	if err == nil {
		t.Fatal("expected rejection of an undefined class")
	}
	if !strings.Contains(err.Error(), "would render unstyled") {
		t.Errorf("message should say what goes wrong, got %q", err)
	}
	if !strings.Contains(err.Error(), "did you mean") {
		t.Errorf("message should suggest a near miss, got %q", err)
	}

	// A tw utility is not part of Collect's universe: the two oracles are
	// genuinely distinct, which is the whole point of parameterising them.
	if err := oracle.Resolve("inline-flex"); err == nil {
		t.Error("Collect's oracle should not accept a tw utility it does not define")
	}
}

// TestStylesheetOracleReportsDeclaredProperties proves the collision check has
// the information it needs on this substrate too — without it, base/variant
// shadowing could not be detected for a bespoke-CSS product.
func TestStylesheetOracleReportsDeclaredProperties(t *testing.T) {
	t.Parallel()

	oracle := collectOracle(t)
	properties, err := oracle.Properties("pill")
	if err != nil {
		t.Fatalf("properties for a shipped class: %v", err)
	}
	if len(properties) == 0 {
		t.Fatal("a shipped class declares no properties; the body parser is not recording declarations")
	}
	for _, property := range properties {
		if strings.HasPrefix(property, "--") {
			t.Errorf("custom property %q must be excluded: they compose rather than collide", property)
		}
	}
	t.Logf("%q declares %d properties", "pill", len(properties))
}

// TestBothOraclesSatisfyTheSameContract pins the abstraction: tw and a product
// stylesheet answer the same questions, which is what lets one generator serve
// both substrates.
func TestBothOraclesSatisfyTheSameContract(t *testing.T) {
	t.Parallel()

	vocabulary, err := LoadVocabulary()
	if err != nil {
		t.Fatalf("load tw vocabulary: %v", err)
	}
	var oracles []ClassOracle = []ClassOracle{vocabulary}
	if _, err := os.Stat(filepath.Join(
		"..", "..", "..", "..", "..",
		"modules", "platformkit-business-modules", "collectibles_management",
		"browser", "css", "collect.css",
	)); err == nil {
		oracles = append(oracles, collectOracle(t))
	}

	for _, oracle := range oracles {
		if strings.TrimSpace(oracle.Name()) == "" {
			t.Error("every oracle must name its universe for error messages")
		}
		// Each oracle rejects a class no universe defines.
		if err := oracle.Resolve("pk-definitely-not-a-real-class-xyz"); err == nil {
			t.Errorf("%s accepted a class no universe defines", oracle.Name())
		}
	}
}

// TestGenerateHonorsTheSuppliedOracle proves the abstraction is load-bearing:
// Badge generates fine under tw, and is refused under Collect's universe,
// because Collect's stylesheet does not define tw's utility classes. A
// component must be governed by the universe its product actually ships.
func TestGenerateHonorsTheSuppliedOracle(t *testing.T) {
	t.Parallel()

	if _, err := Generate(BadgeSource(), BadgeWebStyle()); err != nil {
		t.Fatalf("Badge should generate under the default tw oracle: %v", err)
	}

	oracle := collectOracle(t)
	_, err := Generate(BadgeSource(), BadgeWebStyle(), WithOracle(oracle))
	if err == nil {
		t.Fatal("Badge's tw utilities are not defined by Collect's stylesheet; generation should be refused")
	}
	if !strings.Contains(err.Error(), "would render unstyled") {
		t.Fatalf("rejection should explain the consequence, got %q", err)
	}
}
