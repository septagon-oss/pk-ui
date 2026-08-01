// Implements: REQ-011.
// Per: ADR-0031, ADR-0076.
// Discipline: C-14.

package componentgen

// stylesheet_edit.go is the write half of the style round trip: applying a
// designer's change to a product's own stylesheet without disturbing anything
// else in it.
//
// The requirement that shapes every decision here is that the target is
// hand-authored. Collect's stylesheet is ~13k lines carrying section banners,
// rationale comments, and deliberate ordering. A receive that re-serialised it
// would technically be correct and practically unusable — the diff would be
// the whole file and no reviewer could see what a designer actually changed.
// So an edit is a byte-level splice of one declaration's value, and every
// other byte is preserved exactly. This mirrors the discipline the client.yaml
// adapter already applies: parse for understanding, but write surgically.
//
// The coordinate is the full selector, never the bare class. Collect defines
// ".pill" in four different rules; "the padding of .pill" is not a question
// with one answer, while "the padding of `.pill.is-active`" is.

import (
	"fmt"
	"strings"
)

// Declaration is one property in one rule, with the byte range of its value so
// the value can be replaced without re-rendering the file.
type Declaration struct {
	// Selector is the rule's full selector list, whitespace-normalised.
	Selector string
	// Property is the declared CSS property.
	Property string
	// Value is the declared value, trimmed.
	Value string
	// AtRule is the enclosing at-rule prelude ("@media (min-width: 40rem)"),
	// empty at top level. Two rules with the same selector under different
	// media queries are different coordinates.
	AtRule string
	// ValueStart and ValueEnd bound the value in the source bytes.
	ValueStart, ValueEnd int
}

// Coordinate identifies one declaration uniquely.
func (d Declaration) Coordinate() string {
	if d.AtRule == "" {
		return d.Selector + " { " + d.Property + " }"
	}
	return d.AtRule + " { " + d.Selector + " { " + d.Property + " } }"
}

// IndexDeclarations records every declaration a stylesheet makes.
//
// The parse is structural, not semantic: it understands comments, nesting, and
// where a value begins and ends. It does not resolve the cascade, because the
// only questions the round trip asks are "what does this rule declare" and
// "where are those bytes".
func IndexDeclarations(css []byte) ([]Declaration, error) {
	source := string(css)
	var out []Declaration
	var selectors []string // rule nesting; at-rules push ""
	var atRules []string

	position := 0
	preludeStart := 0

	currentSelector := func() string {
		for index := len(selectors) - 1; index >= 0; index-- {
			if selectors[index] != "" {
				return selectors[index]
			}
		}
		return ""
	}
	currentAtRule := func() string {
		for index := len(atRules) - 1; index >= 0; index-- {
			if atRules[index] != "" {
				return atRules[index]
			}
		}
		return ""
	}

	for position < len(source) {
		// Comments are skipped wholesale so a "{" or ";" inside one cannot
		// be mistaken for structure.
		if strings.HasPrefix(source[position:], "/*") {
			end := strings.Index(source[position+2:], "*/")
			if end < 0 {
				break
			}
			position += 2 + end + 2
			continue
		}

		switch source[position] {
		case '{':
			prelude := strings.TrimSpace(source[preludeStart:position])
			if strings.HasPrefix(prelude, "@") {
				atRules = append(atRules, normalizeSelector(prelude))
				selectors = append(selectors, "")
			} else {
				atRules = append(atRules, "")
				selectors = append(selectors, normalizeSelector(prelude))
			}
			position++
			preludeStart = position

		case '}':
			if len(selectors) > 0 {
				selectors = selectors[:len(selectors)-1]
				atRules = atRules[:len(atRules)-1]
			}
			position++
			preludeStart = position

		case ';':
			if selector := currentSelector(); selector != "" {
				if declaration, ok := parseDeclaration(source, preludeStart, position); ok {
					declaration.Selector = selector
					declaration.AtRule = currentAtRule()
					out = append(out, declaration)
				}
			}
			position++
			preludeStart = position

		default:
			position++
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("stylesheet declares nothing parseable")
	}
	return out, nil
}

// parseDeclaration reads "property: value" from the given span.
func parseDeclaration(source string, start, end int) (Declaration, bool) {
	segment := source[start:end]
	colon := strings.Index(segment, ":")
	if colon < 0 {
		return Declaration{}, false
	}
	property := strings.TrimSpace(segment[:colon])
	if property == "" || strings.ContainsAny(property, "{}()") {
		return Declaration{}, false
	}
	rawValue := segment[colon+1:]
	trimmedLeft := strings.TrimLeft(rawValue, " \t\r\n")
	valueStart := start + colon + 1 + (len(rawValue) - len(trimmedLeft))
	value := strings.TrimRight(trimmedLeft, " \t\r\n")
	if value == "" {
		return Declaration{}, false
	}
	return Declaration{
		Property:   property,
		Value:      value,
		ValueStart: valueStart,
		ValueEnd:   valueStart + len(value),
	}, true
}

// normalizeSelector collapses whitespace so the same selector written across
// two lines yields one coordinate.
func normalizeSelector(selector string) string {
	return strings.Join(strings.Fields(selector), " ")
}

// SetValue replaces one declaration's value and returns the new bytes.
//
// It fails closed on every ambiguity rather than guessing: an unknown
// coordinate, a value that no longer matches what the change set observed, or
// a coordinate that the stylesheet declares more than once. The last case
// matters — a duplicated declaration is legal CSS where the last wins, so
// silently editing the first would change nothing visible and look like a
// successful write.
func SetValue(css []byte, selector, property, expected, replacement string) ([]byte, error) {
	declarations, err := IndexDeclarations(css)
	if err != nil {
		return nil, err
	}
	selector = normalizeSelector(selector)

	var matches []Declaration
	for _, declaration := range declarations {
		if declaration.Selector == selector && declaration.Property == property && declaration.AtRule == "" {
			matches = append(matches, declaration)
		}
	}
	switch len(matches) {
	case 0:
		return nil, fmt.Errorf(
			"stylesheet does not declare %q on %q — nothing to update", property, selector)
	case 1:
	default:
		return nil, fmt.Errorf(
			"stylesheet declares %q on %q %d times; the coordinate is ambiguous and the edit is refused",
			property, selector, len(matches))
	}

	target := matches[0]
	if target.Value != expected {
		return nil, fmt.Errorf(
			"%q on %q is now %q, not the %q the change set observed; re-export and retry",
			property, selector, target.Value, expected)
	}
	if replacement == "" {
		return nil, fmt.Errorf("refusing to write an empty value for %q on %q", property, selector)
	}

	out := make([]byte, 0, len(css)+len(replacement))
	out = append(out, css[:target.ValueStart]...)
	out = append(out, replacement...)
	out = append(out, css[target.ValueEnd:]...)
	return out, nil
}
