// Implements: REQ-011.
// Per: ADR-0031, ADR-0076.
// Discipline: C-14.

package stylesheet

// computed.go reduces a declaration to what the browser would act on, so that
// two ways of saying the same thing compare equal.
//
// This is what separates a real visual difference from a difference in
// notation. A component styled "width: 20rem" and one styled "width: 320px"
// paint the same pixels; so do "height: var(--size)" where --size is 28px and
// "height: 1.75rem"; so do "background: #2a9d5b" and "background-color:
// #2A9D5B". Comparing the declarations as written calls all three pairs
// different, and every one of those is a false alarm — the expensive kind,
// because it buries the true ones.
//
// The reductions here are only the ones that are always safe: a unit
// conversion with a fixed ratio, a custom property with a known value, a
// shorthand whose value can only mean one longhand. Anything requiring layout,
// inheritance, or the element's context is left alone, because a wrong
// reduction produces a confident false match, which is worse than a false
// alarm.

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// rootFontSize is the pixel size of one rem. Both targets render at the
// browser default; a product that overrides it would need this parameterised.
const rootFontSize = 16.0

// ParseInline reads a style attribute into its declarations.
//
// Inline styles cannot be ignored when comparing what two targets paint:
// markup that positions itself with style="position:relative;width:320px"
// is styled exactly as much as markup that uses classes, and treating it as
// unstyled reports differences that are not there.
func ParseInline(style string) map[string]string {
	out := map[string]string{}
	for _, statement := range strings.Split(style, ";") {
		colon := strings.Index(statement, ":")
		if colon < 0 {
			continue
		}
		property := strings.TrimSpace(statement[:colon])
		value := strings.TrimSpace(statement[colon+1:])
		if property == "" || value == "" {
			continue
		}
		out[property] = value
	}
	return out
}

// Normalize reduces one declaration to the longhand properties and comparable
// values it produces, resolving custom properties against vars.
//
// A custom property declaration itself normalizes to nothing: it paints
// nothing on its own, and its effect is already accounted for wherever it is
// used.
func Normalize(property, value string, vars map[string]string) map[string]string {
	property = strings.ToLower(strings.TrimSpace(property))
	if strings.HasPrefix(property, "--") {
		return nil
	}
	value = NormalizeValue(ResolveVars(value, vars))
	if value == "" {
		return nil
	}
	if longhand, ok := colorShorthand(property, value); ok {
		return map[string]string{longhand: value}
	}
	return map[string]string{property: value}
}

// colorShorthand reports the longhand a shorthand sets when its value can only
// be a colour.
//
// "background: #2a9d5b" sets exactly background-color and resets the other
// background longhands to their initial values, so comparing it against
// "background-color: #2a9d5b" is sound. A shorthand carrying anything more
// than a colour is left as written, because splitting it would require
// deciding which part is which.
func colorShorthand(property, value string) (string, bool) {
	if property != "background" || !isColor(value) {
		return "", false
	}
	return "background-color", true
}

// isColor reports whether a value sets a colour and nothing else.
//
// It must be the whole value: "#2a9d5b url(x.png) no-repeat" begins with a
// colour but also sets an image and a repeat, and calling that background-color
// would claim the component paints something it does not.
//
// A single unresolved reference counts. "background: var(--sm-orange)" cannot
// be evaluated here, but whatever the variable holds, a one-token background
// that is not an image sets the background colour — and treating it as an
// unrelated property would report a colour difference between two components
// that may well paint the same one.
func isColor(value string) bool {
	if strings.ContainsAny(value, " \t") {
		return false
	}
	for _, notAColor := range []string{"url(", "linear-gradient(", "radial-gradient(", "conic-gradient(", "image-set("} {
		if strings.HasPrefix(value, notAColor) {
			return false
		}
	}
	return value != "" && value != "none"
}

// ResolveVars substitutes custom properties into a value.
//
// A var() reference is resolved to the variable's value when it is known and
// to its declared fallback when it is not, which is what the browser does. A
// reference with neither is left as written rather than blanked, so it reads
// as itself in a failure message instead of as an empty value.
func ResolveVars(value string, vars map[string]string) string {
	// Nested references are resolved by repeating; the bound stops a variable
	// defined in terms of itself from looping.
	for range 8 {
		start := strings.Index(value, "var(")
		if start < 0 {
			return value
		}
		end := matchingParen(value, start+len("var(")-1)
		if end < 0 {
			return value
		}
		name, fallback := splitVarArgs(value[start+len("var(") : end])
		replacement, known := vars[name]
		if !known {
			if fallback == "" {
				return value
			}
			replacement = fallback
		}
		value = value[:start] + replacement + value[end+1:]
	}
	return value
}

// matchingParen finds the ")" closing the "(" at open.
func matchingParen(value string, open int) int {
	depth := 0
	for index := open; index < len(value); index++ {
		switch value[index] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return index
			}
		}
	}
	return -1
}

// splitVarArgs separates a var() reference's name from its fallback.
func splitVarArgs(args string) (name, fallback string) {
	depth := 0
	for index := range len(args) {
		switch args[index] {
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				return strings.TrimSpace(args[:index]), strings.TrimSpace(args[index+1:])
			}
		}
	}
	return strings.TrimSpace(args), ""
}

// NormalizeValue reduces a value to a comparable form: one space between
// tokens, lower case, lengths in pixels, and zero without a unit.
//
// Case and whitespace carry no meaning in a CSS value. Lengths do, but rem and
// px differ only by a constant, so expressing both in pixels compares what
// gets painted rather than which unit the author preferred.
func NormalizeValue(value string) string {
	fields := strings.Fields(strings.ToLower(strings.TrimSpace(value)))
	for index, field := range fields {
		fields[index] = normalizeQuantity(strings.TrimSuffix(field, ","))
		if strings.HasSuffix(field, ",") {
			fields[index] += ","
		}
	}
	return strings.Join(fields, " ")
}

// numeric splits a token into its number and unit.
var numeric = regexp.MustCompile(`^([+-]?(?:\d+\.?\d*|\.\d+))([a-z%]*)$`)

// lengthUnits are the units for which an author may drop the unit on zero.
var lengthUnits = map[string]bool{"px": true, "rem": true, "em": true, "%": true}

// normalizeQuantity writes one numeric token the same way whoever authored it
// would have: pixels where the unit converts, a leading zero where one was
// omitted, and no trailing zeros.
//
// The leading zero matters more than it looks. One target writes
// "letter-spacing: 0.1em" and the other ".1em"; those are the same
// letter-spacing, and reporting them as a difference is noise that hides the
// cases where the two really do disagree.
func normalizeQuantity(token string) string {
	match := numeric.FindStringSubmatch(token)
	if match == nil {
		return token
	}
	number, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		return token
	}
	unit := match[2]
	if unit == "rem" {
		number, unit = number*rootFontSize, "px"
	}
	if number == 0 && (unit == "" || lengthUnits[unit]) {
		// "0", "0px" and "0rem" are the same length; the unit on zero is
		// optional in CSS and carries no meaning. A zero in another dimension
		// keeps its unit, because "0s" and "0deg" are not lengths.
		return "0"
	}
	return formatNumber(number) + unit
}

// formatNumber renders a number without trailing zeros.
func formatNumber(number float64) string {
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.4f", number), "0"), ".")
}
