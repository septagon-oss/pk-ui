// Implements: REQ-011.
// Per: ADR-0031, ADR-0076.
// Discipline: C-14.

package componentgen

// stylesheet_edit.go keeps the round trip's vocabulary local while the parser
// itself lives in the public stylesheet package.
//
// The parse moved out because it answers a question two callers ask: this
// generator needs it to know what a class declares, and parity checking needs
// it to know whether a design contract and a rendered component resolve to the
// same CSS. One parser answering both is the point — two would drift, and the
// whole argument for generating from the stylesheet is that there is a single
// source of truth.

import "github.com/septagon-oss/pk-ui/stylesheet"

// Declaration is one property in one rule. See stylesheet.Declaration.
type Declaration = stylesheet.Declaration

// IndexDeclarations records every declaration a stylesheet makes.
func IndexDeclarations(css []byte) ([]Declaration, error) {
	return stylesheet.IndexDeclarations(css)
}

// SetValue replaces one declaration's value and returns the new bytes.
func SetValue(css []byte, selector, property, expected, replacement string) ([]byte, error) {
	return stylesheet.SetValue(css, selector, property, expected, replacement)
}
