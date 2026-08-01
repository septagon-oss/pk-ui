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
	return ExpandDeclaration(property, ResolveVars(value, vars))
}

// ExpandDeclaration is Normalize without the variable substitution: it reduces
// one declaration to longhand properties and comparable values, leaving every
// reference exactly as written.
//
// It is separate because resolving references and canonicalising a declaration
// are wanted independently. A caller deriving a design contract has already
// expanded the product's own tokens and needs the platform's left standing, so
// running the browser's resolution over the result would undo the distinction
// it just made.
func ExpandDeclaration(property, value string) map[string]string {
	property = strings.ToLower(strings.TrimSpace(property))
	if strings.HasPrefix(property, "--") {
		return nil
	}
	value = NormalizeValue(value)
	if value == "" {
		return nil
	}
	if longhand, ok := colorShorthand(property, value); ok {
		return map[string]string{longhand: value}
	}
	if expanded := expandShorthand(property, value); expanded != nil {
		return expanded
	}
	return map[string]string{property: value}
}

// boxSides are the shorthands whose one-to-four values map onto top, right,
// bottom and left in that order.
var boxSides = map[string][4]string{
	"padding": {"padding-top", "padding-right", "padding-bottom", "padding-left"},
	"margin":  {"margin-top", "margin-right", "margin-bottom", "margin-left"},
	"inset":   {"top", "right", "bottom", "left"},
}

// expandShorthand rewrites a shorthand as the longhands it sets, so that a
// contract stating "padding-left: 8px" and a component writing
// "padding: 18px 18px 16px 8px" agree about the left padding.
//
// Only shorthands whose parts can be assigned without guessing are expanded.
// The box shorthands are positional and fully determined; border is
// unordered but its three parts are distinguishable by what they are. A
// shorthand that cannot be decomposed with certainty is left as written,
// because a wrong split claims a component paints something it does not.
func expandShorthand(property, value string) map[string]string {
	parts := splitTopLevel(value)
	if property == "border" {
		// border is expanded even with a reference in it; see expandBorder.
		return expandBorder(parts)
	}
	if strings.Contains(value, "var(") {
		// In a positional shorthand an unresolved reference could stand for
		// any number of values, which moves every part after it. Where the
		// pieces land cannot be known, so nothing is claimed.
		return nil
	}
	if sides, ok := boxSides[property]; ok {
		return expandBox(sides, parts)
	}
	if property == "gap" {
		switch len(parts) {
		case 1:
			return map[string]string{"row-gap": parts[0], "column-gap": parts[0]}
		case 2:
			return map[string]string{"row-gap": parts[0], "column-gap": parts[1]}
		}
		return nil
	}
	return nil
}

// expandBox applies the one-to-four value rule.
func expandBox(sides [4]string, parts []string) map[string]string {
	var top, right, bottom, left string
	switch len(parts) {
	case 1:
		top, right, bottom, left = parts[0], parts[0], parts[0], parts[0]
	case 2:
		top, right, bottom, left = parts[0], parts[1], parts[0], parts[1]
	case 3:
		top, right, bottom, left = parts[0], parts[1], parts[2], parts[1]
	case 4:
		top, right, bottom, left = parts[0], parts[1], parts[2], parts[3]
	default:
		return nil
	}
	return map[string]string{
		sides[0]: top, sides[1]: right, sides[2]: bottom, sides[3]: left,
	}
}

// borderStyles is the closed set of line styles, which is what makes the
// border shorthand decomposable despite being unordered.
var borderStyles = map[string]bool{
	"none": true, "hidden": true, "dotted": true, "dashed": true, "solid": true,
	"double": true, "groove": true, "ridge": true, "inset": true, "outset": true,
}

// borderWidths are the named widths, alongside any length.
var borderWidths = map[string]bool{"thin": true, "medium": true, "thick": true}

// expandBorder assigns each part of a border shorthand by what it is: a style
// from the closed keyword set, a width if it is a length or a named width, and
// a colour otherwise.
//
// Unlike the positional shorthands, this survives an unresolved reference.
// border is unordered and takes at most one value per slot, so in
// "1px solid var(--sm-border)" the width and the style are certain whatever
// the variable holds — if it held a width there would be two, which is not a
// valid shorthand. The reference itself is placed only when exactly one slot
// is left for it; with more than one open slot it could be either, so nothing
// is claimed for it.
//
// This matters because the shipped stylesheet writes borders exactly that way,
// and refusing to expand them reported a 1px hairline as missing from a
// component that draws one.
func expandBorder(parts []string) map[string]string {
	out := map[string]string{}
	var references []string
	for _, part := range parts {
		var slot string
		switch {
		case borderStyles[part]:
			slot = "border-style"
		case borderWidths[part] || numeric.MatchString(part):
			slot = "border-width"
		case strings.Contains(part, "var("):
			references = append(references, part)
			continue
		case isColorToken(part):
			slot = "border-color"
		default:
			return nil
		}
		if _, taken := out[slot]; taken {
			// Two values for one slot is not a border shorthand; whatever this
			// is, it is not safe to decompose.
			return nil
		}
		out[slot] = part
	}

	var open []string
	for _, slot := range []string{"border-width", "border-style", "border-color"} {
		if _, taken := out[slot]; !taken {
			open = append(open, slot)
		}
	}
	if len(references) == 1 && len(open) == 1 {
		out[open[0]] = references[0]
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// splitTopLevel splits a value on spaces that are not inside brackets, so a
// function call stays one part.
func splitTopLevel(value string) []string {
	var parts []string
	depth, start := 0, 0
	for index := range len(value) {
		switch value[index] {
		case '(':
			depth++
		case ')':
			depth--
		case ' ':
			if depth == 0 {
				if start < index {
					parts = append(parts, value[start:index])
				}
				start = index + 1
			}
		}
	}
	if start < len(value) {
		parts = append(parts, value[start:])
	}
	return parts
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
func isColor(value string) bool {
	parts := splitTopLevel(value)
	return len(parts) == 1 && isColorToken(parts[0])
}

// isColorToken reports whether one value token is a colour. The caller has
// already established that it is a single token, so spaces inside a function
// call are not a reason to reject it.
//
// A single unresolved reference counts. "background: var(--sm-orange)" cannot
// be evaluated here, but whatever the variable holds, a one-token background
// that is not an image sets the background colour — and treating it as an
// unrelated property would report a colour difference between two components
// that may well paint the same one.
func isColorToken(token string) bool {
	for _, notAColor := range []string{"url(", "linear-gradient(", "radial-gradient(", "conic-gradient(", "image-set("} {
		if strings.HasPrefix(token, notAColor) {
			return false
		}
	}
	return token != "" && token != "none"
}

// ExpandDefined substitutes only the custom properties this stylesheet
// defines, leaving any other reference exactly as written.
//
// The difference from ResolveVars is the whole point: ResolveVars answers
// "what does this paint", taking a fallback when a variable is unknown, which
// is what a browser does. This answers "what does this refer to", and an
// unknown variable is the interesting part of the answer.
//
// A product composing on a design system declares its own tokens against the
// system's: "--sm-green: rgb(var(--surface-success, 69 159 80))". The product
// defines --sm-green and does not define --surface-success, because the theme
// does. Expanding only what the product defines turns var(--sm-green) into a
// value that still names the platform role, which is what a design target has
// to bind against. Resolving it the browser's way would collapse it to the
// literal standing by for an unthemed page, and the token would be lost.
func ExpandDefined(value string, vars map[string]string) string {
	for range 8 {
		expanded, changed := expandOnce(value, vars)
		if !changed {
			return expanded
		}
		value = expanded
	}
	return value
}

// expandOnce replaces every reference to a defined variable, reporting whether
// it replaced any.
func expandOnce(value string, vars map[string]string) (string, bool) {
	var b strings.Builder
	var changed bool
	for {
		start := strings.Index(value, "var(")
		if start < 0 {
			b.WriteString(value)
			return b.String(), changed
		}
		end := matchingParen(value, start+len("var(")-1)
		if end < 0 {
			b.WriteString(value)
			return b.String(), changed
		}
		name, _ := splitVarArgs(value[start+len("var(") : end])
		replacement, defined := vars[name]
		b.WriteString(value[:start])
		if defined {
			b.WriteString(replacement)
			changed = true
		} else {
			// Left whole, fallback included: the reference is the answer.
			b.WriteString(value[start : end+1])
		}
		value = value[end+1:]
	}
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
